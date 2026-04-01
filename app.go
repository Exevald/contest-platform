package main

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	appmodel "contest-platform/pkg/contestplatform/app/model"
	appservice "contest-platform/pkg/contestplatform/app/service"
	domainmodel "contest-platform/pkg/contestplatform/domain/model"
	"contest-platform/pkg/contestplatform/infrastructure"
)

const applicationID = "contestplatform"

type App struct {
	ctx       context.Context
	container *infrastructure.DependencyContainer
	initErr   error
}

type StartupData struct {
	Title     string                `json:"title"`
	Languages []appmodel.UILanguage `json:"languages"`
	Tasks     []Task                `json:"tasks"`
}

type Task struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Label string `json:"label"`
}

type SendFileResponse struct {
	SubmissionID string `json:"submission_id,omitempty"`
	Error        string `json:"error,omitempty"`
}

type SubmissionStatus struct {
	SubmissionID string `json:"submissionId"`
	ProblemID    string `json:"problemId"`
	Language     string `json:"language"`
	Verdict      string `json:"verdict"`
	CreatedAt    string `json:"createdAt"`
	TestsPassed  int    `json:"testsPassed"`
	TestsTotal   int    `json:"testsTotal"`
}

func NewApp() *App {
	container, err := infrastructure.NewDependencyContainer(applicationID)
	return &App{
		container: container,
		initErr:   err,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	if a.container != nil {
		a.container.Start(ctx, func(submissionID string) {
			runtime.EventsEmit(ctx, "on-verdict-changed", submissionID)
		})
	}
}

func (a *App) shutdown(context.Context) {
	if a.container != nil {
		_ = a.container.Close()
	}
}

func (a *App) GetStartupData() (StartupData, error) {
	if a.initErr != nil {
		return StartupData{}, a.initErr
	}

	problems, err := a.container.ProblemRepository().List()
	if err != nil {
		return StartupData{}, fmt.Errorf("load startup data: %w", err)
	}

	tasks := make([]Task, 0, len(problems))
	for _, problem := range problems {
		tasks = append(tasks, Task{
			ID:    string(problem.ID()),
			Type:  "table",
			Label: string(problem.Title()),
		})
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Label < tasks[j].Label
	})

	return StartupData{
		Title:     "ContestPlatform",
		Languages: appmodel.SupportedUILanguages(),
		Tasks:     tasks,
	}, nil
}

func (a *App) GetData(id string) (string, error) {
	problem, err := a.findProblem(id)
	if err != nil {
		return "", err
	}
	return problem.Description(), nil
}

func (a *App) ResetTask(id string) (string, error) {
	return a.GetData(id)
}

func (a *App) SendFile(id string, language string, text string) (SendFileResponse, error) {
	if a.initErr != nil {
		return SendFileResponse{}, a.initErr
	}

	submissionID, err := a.container.SubmissionService().Submit(a.ctx, appservice.SubmitRequest{
		ProblemID: id,
		Language:  normalizeLanguage(language),
		Source:    text,
	})
	if err != nil {
		return SendFileResponse{
			Error: err.Error(),
		}, nil
	}

	return SendFileResponse{
		SubmissionID: submissionID,
	}, nil
}

func (a *App) GetLatestSubmission(id string) (*SubmissionStatus, error) {
	if a.initErr != nil {
		return nil, a.initErr
	}

	submission, err := a.container.SubmissionRepository().FindLatest(domainmodel.ProblemID(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("load latest submission for %s: %w", id, err)
	}

	status := makeSubmissionStatus(submission)
	return &status, nil
}

func (a *App) GetSubmissionStatus(id string) (*SubmissionStatus, error) {
	if a.initErr != nil {
		return nil, a.initErr
	}

	submission, err := a.container.SubmissionRepository().Find(domainmodel.SubmissionID(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("load submission %s: %w", id, err)
	}

	status := makeSubmissionStatus(submission)
	return &status, nil
}

func (a *App) findProblem(id string) (domainmodel.Problem, error) {
	if a.initErr != nil {
		return nil, a.initErr
	}

	problem, err := a.container.ProblemRepository().Find(domainmodel.ProblemID(id))
	if err != nil {
		return nil, fmt.Errorf("load problem %s: %w", id, err)
	}
	return problem, nil
}

func normalizeLanguage(language string) string {
	for key, config := range appmodel.Languages {
		if config.DisplayName == language || string(key) == language {
			return string(key)
		}
	}

	return language
}

func makeSubmissionStatus(submission domainmodel.Submission) SubmissionStatus {
	testsPassed := 0
	for _, result := range submission.TestResults() {
		if result.Verdict == domainmodel.VerdictOK {
			testsPassed++
		}
	}

	return SubmissionStatus{
		SubmissionID: string(submission.ID()),
		ProblemID:    string(submission.ProblemID()),
		Language:     string(submission.Language()),
		Verdict:      string(submission.Verdict()),
		CreatedAt:    submission.CreatedAt().Format(time.RFC3339),
		TestsPassed:  testsPassed,
		TestsTotal:   len(submission.TestResults()),
	}
}
