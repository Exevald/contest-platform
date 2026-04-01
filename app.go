package main

import (
	"context"
	"fmt"
	"sort"

	appmodel "contest-platform/pkg/contestplatform/app/model"
	appservice "contest-platform/pkg/contestplatform/app/service"
	domainmodel "contest-platform/pkg/contestplatform/domain/model"
	"contest-platform/pkg/contestplatform/infrastructure"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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
