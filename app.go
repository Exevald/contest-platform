package main

import (
	"context"
	"fmt"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	appmodel "contest-platform/pkg/contestplatform/app/model"
	appquery "contest-platform/pkg/contestplatform/app/query"
	appservice "contest-platform/pkg/contestplatform/app/service"
	"contest-platform/pkg/contestplatform/app/storage"
	"contest-platform/pkg/contestplatform/config"
	"contest-platform/pkg/contestplatform/infrastructure"
)

const applicationID = "contestplatform"

type App struct {
	ctx       context.Context
	container *infrastructure.DependencyContainer
	initErr   error
}

type StartupData struct {
	Title           string                `json:"title"`
	Languages       []appmodel.UILanguage `json:"languages"`
	Tasks           []Task                `json:"tasks"`
	WorkspaceViews  []WorkspaceView       `json:"workspaceViews"`
	ParticipantCode string                `json:"participantCode"`
	SelectedTheme   string                `json:"selectedTheme"`
}

type Task struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Label string `json:"label"`
	Theme string `json:"theme"`
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

type StartSessionRequest struct {
	ParticipantCode string `json:"participantCode"`
	Theme           string `json:"theme"`
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

	session, err := a.container.SessionStorage().Load()
	if err != nil {
		return StartupData{}, fmt.Errorf("load participant session: %w", err)
	}

	themes := config.AllThemes()
	tasks := make([]Task, 0, 6)
	for _, theme := range themes {
		for _, problem := range theme.Tasks {
			tasks = append(tasks, Task{
				ID:    problem.ID,
				Type:  "statement",
				Label: problem.Label,
				Theme: string(theme.Key),
			})
		}
	}

	var participantCode, selectedTheme string
	if session != nil {
		participantCode = session.ParticipantCode
		selectedTheme = session.Theme
	}

	return StartupData{
		Title:           "ContestPlatform",
		Languages:       appmodel.SupportedUILanguages(),
		Tasks:           tasks,
		ParticipantCode: participantCode,
		SelectedTheme:   selectedTheme,
		WorkspaceViews: []WorkspaceView{
			{ID: "statement", Label: "Задание"},
			{ID: "submission_history", Label: "История посылок"},
		},
	}, nil
}

func (a *App) StartSession(request StartSessionRequest) (StartupData, error) {
	if a.initErr != nil {
		return StartupData{}, a.initErr
	}

	participantCode := request.ParticipantCode
	if participantCode == "" {
		return StartupData{}, fmt.Errorf("participant code is required")
	}
	if request.Theme != "pizza" && request.Theme != "soc" {
		return StartupData{}, fmt.Errorf("invalid theme")
	}

	if err := a.container.SessionStorage().Save(storage.ParticipantSession{
		ParticipantCode: participantCode,
		Theme:           request.Theme,
	}); err != nil {
		return StartupData{}, err
	}

	return a.GetStartupData()
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

	session, err := a.container.SessionStorage().Load()
	if err != nil {
		return SendFileResponse{}, err
	}
	if session == nil || session.ParticipantCode == "" {
		return SendFileResponse{
			Error: "participant session is not initialized",
		}, nil
	}

	submissionID, err := a.container.SubmissionService().Submit(a.ctx, appservice.SubmitRequest{
		ProblemID:       id,
		ParticipantCode: session.ParticipantCode,
		Language:        normalizeLanguage(language),
		Source:          text,
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
