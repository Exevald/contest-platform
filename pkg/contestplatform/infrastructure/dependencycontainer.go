package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	appmodel "contest-platform/pkg/contestplatform/app/model"
	appservice "contest-platform/pkg/contestplatform/app/service"
	domainmodel "contest-platform/pkg/contestplatform/domain/model"
	domainservice "contest-platform/pkg/contestplatform/domain/service"
	"contest-platform/pkg/contestplatform/infrastructure/sandbox"
	sqliterepo "contest-platform/pkg/contestplatform/infrastructure/sqlite/repo"
)

type DependencyContainer struct {
	db                   *sql.DB
	problemRepository    domainmodel.ProblemRepository
	submissionRepository domainmodel.SubmissionRepository
	submissionService    appservice.SubmissionService
	gradingWorker        appservice.GradingWorker
	gradingTasks         chan appmodel.GradingTask
	startOnce            sync.Once
	closeOnce            sync.Once
}

func NewDependencyContainer(appID string) (*DependencyContainer, error) {
	dataDir, err := resolveAppDataDir(appID)
	if err != nil {
		return nil, err
	}

	db, err := sqliterepo.OpenDatabase(filepath.Join(dataDir, appID+".sqlite"))
	if err != nil {
		return nil, err
	}

	problemRepository := sqliterepo.NewProblemRepository(db)
	submissionRepository := sqliterepo.NewSubmissionRepository(db)
	if err = seedProblems(problemRepository); err != nil {
		_ = db.Close()
		return nil, err
	}

	gradingTasks := make(chan appmodel.GradingTask, 32)
	sandboxRunner, err := sandbox.NewSandbox()
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	judgeService := domainservice.NewJudgeService()
	submissionDomainService := domainservice.NewSubmissionService(judgeService)

	container := &DependencyContainer{
		db:                   db,
		problemRepository:    problemRepository,
		submissionRepository: submissionRepository,
		submissionService: appservice.NewSubmissionService(
			submissionRepository,
			problemRepository,
			submissionDomainService,
			gradingTasks,
		),
		gradingWorker: appservice.NewGradingWorker(
			submissionRepository,
			problemRepository,
			sandboxRunner,
			judgeService,
			submissionDomainService,
			gradingTasks,
		),
		gradingTasks: gradingTasks,
	}

	return container, nil
}

func (c *DependencyContainer) Start(ctx context.Context, notify func(string)) {
	c.startOnce.Do(func() {
		c.gradingWorker.SetNotifyCallback(notify)
		go func() {
			_ = c.gradingWorker.Run(ctx)
		}()
	})
}

func (c *DependencyContainer) Close() error {
	var err error
	c.closeOnce.Do(func() {
		close(c.gradingTasks)
		err = c.db.Close()
	})
	return err
}

func (c *DependencyContainer) ProblemRepository() domainmodel.ProblemRepository {
	return c.problemRepository
}

func (c *DependencyContainer) SubmissionService() appservice.SubmissionService {
	return c.submissionService
}

func (c *DependencyContainer) SubmissionRepository() domainmodel.SubmissionRepository {
	return c.submissionRepository
}

func resolveAppDataDir(appID string) (string, error) {
	root, err := os.UserConfigDir()
	if err != nil {
		root = "."
	}

	path := filepath.Join(root, appID)
	if mkErr := os.MkdirAll(path, 0755); mkErr != nil {
		return "", fmt.Errorf("create app data dir: %w", mkErr)
	}

	return path, nil
}

func seedProblems(repo domainmodel.ProblemRepository) error {
	existing, err := repo.List()
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		return nil
	}

	problems, err := sampleProblems()
	if err != nil {
		return err
	}
	for _, problem := range problems {
		if err = repo.Store(problem); err != nil {
			return err
		}
	}

	return nil
}

func sampleProblems() ([]domainmodel.Problem, error) {
	limits := domainmodel.Constraints{
		TimeLimit:   2 * time.Second,
		MemoryLimit: 256 * 1024 * 1024,
	}

	orders, err := domainmodel.NewProblem(
		"task-orders",
		"Фильтр заказов",
		"Название,Цена\nМаргарита,500\nПепперони,400\n",
		limits,
	)
	if err != nil {
		return nil, err
	}
	if err = orders.AddTestCase("1 2\n", "3\n", true); err != nil {
		return nil, err
	}

	demand, err := domainmodel.NewProblem(
		"task-demand",
		"Карта спроса",
		"Город,Спрос\nМосква,120\nКазань,80\n",
		limits,
	)
	if err != nil {
		return nil, err
	}
	if err = demand.AddTestCase("2 5\n", "7\n", true); err != nil {
		return nil, err
	}

	routes, err := domainmodel.NewProblem(
		"task-routes",
		"Генератор маршрутов",
		"Маршрут,Длина\nA-B,14\nB-C,9\n",
		limits,
	)
	if err != nil {
		return nil, err
	}
	if err = routes.AddTestCase("10 32\n", "42\n", true); err != nil {
		return nil, err
	}

	return []domainmodel.Problem{orders, demand, routes}, nil
}
