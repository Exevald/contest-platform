package main

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	appmodel "contest-platform/pkg/contestplatform/app/model"
	appquery "contest-platform/pkg/contestplatform/app/query"
	appservice "contest-platform/pkg/contestplatform/app/service"
	"contest-platform/pkg/contestplatform/infrastructure"
)

const applicationID = "contestplatform"

type App struct {
	ctx       context.Context
	container *infrastructure.DependencyContainer
	initErr   error
}

type StartupData struct {
	Title          string                `json:"title"`
	Languages      []appmodel.UILanguage `json:"languages"`
	Tasks          []Task                `json:"tasks"`
	WorkspaceViews []WorkspaceView       `json:"workspaceViews"`
}

type Task struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Label string `json:"label"`
}

type WorkspaceView struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type SendFileResponse struct {
	SubmissionID string `json:"submission_id,omitempty"`
	Error        string `json:"error,omitempty"`
}

type SubmissionStatus struct {
	SubmissionID      string `json:"submissionId"`
	ProblemID         string `json:"problemId"`
	Language          string `json:"language"`
	Verdict           string `json:"verdict"`
	CreatedAt         string `json:"createdAt"`
	CompilationOutput string `json:"compilationOutput"`
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

	problems, err := a.container.PlatformQueryService().ListProblems()
	if err != nil {
		return StartupData{}, fmt.Errorf("load startup data: %w", err)
	}

	tasks := make([]Task, 0, len(problems))
	for _, problem := range problems {
		tasks = append(tasks, Task{
			ID:    problem.ID,
			Type:  "table",
			Label: problem.Title,
		})
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Label < tasks[j].Label
	})

	return StartupData{
		Title:     "ContestPlatform",
		Languages: appmodel.SupportedUILanguages(),
		Tasks:     tasks,
		WorkspaceViews: []WorkspaceView{
			{ID: "statement", Label: "Данные"},
			{ID: "submission_history", Label: "История посылок"},
		},
	}, nil
}

func (a *App) GetData(id string) (string, error) {
	description, err := a.container.PlatformQueryService().GetProblemDescription(id)
	if err != nil {
		return "", err
	}
	return description, nil
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
	history, err := a.GetSubmissionHistory(id)
	if err != nil || len(history) == 0 {
		return nil, err
	}

	return &history[0], nil
}

func (a *App) GetSubmissionStatus(id string) (*SubmissionStatus, error) {
	if a.initErr != nil {
		return nil, a.initErr
	}

	submission, err := a.container.PlatformQueryService().GetSubmissionStatus(id)
	if err != nil || submission == nil {
		return nil, err
	}

	status := makeSubmissionStatus(*submission)
	return &status, nil
}

func (a *App) GetSubmissionHistory(id string) ([]SubmissionStatus, error) {
	if a.initErr != nil {
		return nil, a.initErr
	}

	submissions, err := a.container.PlatformQueryService().ListSubmissionHistory(id)
	if err != nil {
		return nil, fmt.Errorf("load submission history for %s: %w", id, err)
	}

	history := make([]SubmissionStatus, 0, len(submissions))
	for _, submission := range submissions {
		history = append(history, makeSubmissionStatus(submission))
	}

	return history, nil
}

func normalizeLanguage(language string) string {
	for key, config := range appmodel.Languages {
		if config.DisplayName == language || string(key) == language {
			return string(key)
		}
	}

	return language
}

func makeSubmissionStatus(submission appquery.SubmissionView) SubmissionStatus {
	return SubmissionStatus{
		SubmissionID:      submission.ID,
		ProblemID:         submission.ProblemID,
		Language:          submission.Language,
		Verdict:           submission.Verdict,
		CreatedAt:         submission.CreatedAt.Format(time.RFC3339),
		CompilationOutput: submission.CompilationOutput,
	}
}
