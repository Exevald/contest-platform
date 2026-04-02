package infrastructure

import (
	"contest-platform/pkg/contestplatform/app/storage"
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	appmodel "contest-platform/pkg/contestplatform/app/model"
	appquery "contest-platform/pkg/contestplatform/app/query"
	appservice "contest-platform/pkg/contestplatform/app/service"
	"contest-platform/pkg/contestplatform/config"
	domainmodel "contest-platform/pkg/contestplatform/domain/model"
	domainservice "contest-platform/pkg/contestplatform/domain/service"
	"contest-platform/pkg/contestplatform/infrastructure/sandbox"
	sqlitequery "contest-platform/pkg/contestplatform/infrastructure/sqlite/query"
	sqliterepo "contest-platform/pkg/contestplatform/infrastructure/sqlite/repo"
	sqlitestorage "contest-platform/pkg/contestplatform/infrastructure/sqlite/storage"
)

type DependencyContainer struct {
	db                   *sql.DB
	problemRepository    domainmodel.ProblemRepository
	submissionRepository domainmodel.SubmissionRepository
	sessionStorage       storage.SessionStorage
	platformQueryService appquery.PlatformQueryService
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
	sessionStorage := sqlitestorage.NewSessionStorage(db)
	platformQueryService := sqlitequery.NewPlatformQueryService(db)
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
		sessionStorage:       sessionStorage,
		platformQueryService: platformQueryService,
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

func (c *DependencyContainer) SessionStorage() storage.SessionStorage {
	return c.sessionStorage
}

func (c *DependencyContainer) PlatformQueryService() appquery.PlatformQueryService {
	return c.platformQueryService
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

	themes := config.AllThemes()
	problems := make([]domainmodel.Problem, 0, 6)

	for _, meta := range themes {
		for index, task := range meta.Tasks {
			problem, err := domainmodel.NewProblem(
				domainmodel.ProblemID(task.ID),
				domainmodel.Title(task.Label),
				fmt.Sprintf("statement://%s/%s", meta.Key, task.ID),
				limits,
			)
			if err != nil {
				return nil, err
			}

			sampleInput := fmt.Sprintf("%d\n", index+1)
			sampleOutput := fmt.Sprintf("%d\n", index+1)

			if err = problem.AddTestCase(sampleInput, sampleOutput, true); err != nil {
				return nil, err
			}
			if err = problem.AddTestCase(sampleInput, sampleOutput, false); err != nil {
				return nil, err
			}

			problems = append(problems, problem)
		}
	}

	return problems, nil
}
