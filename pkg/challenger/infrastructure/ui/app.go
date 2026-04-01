package ui

import (
	"context"

	"challenger/pkg/challenger/app/service"
	"challenger/pkg/challenger/domain/model"
	"github.com/wailsapp/wails/v2/pkg/runtime" // Важно для событий
)

type App struct {
	ctx    context.Context
	subSvc service.SubmissionService
}

func NewApp(subSvc service.SubmissionService) *App {
	return &App{
		subSvc: subSvc,
	}
}

// Startup вызывается Wails при старте
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

type StartupData struct {
	Title     string          `json:"title"`
	Languages []string        `json:"languages"`
	Tabs      []model.Problem `json:"tabs"`
}

// Переименуем для соответствия CEF-вызову на фронте
func (a *App) Platform_getStartupData() StartupData {
	return StartupData{
		Title:     "iSpring Challenger",
		Languages: []string{"cpp", "python", "go", "pascal", "js", "php"},
		Tabs:      nil, // Тут можно вызвать получение задач
	}
}

func (a *App) SendFile(id string, language string, text string) map[string]string {
	req := service.SubmitRequest{
		ProblemID: id,
		Language:  language,
		Source:    text,
	}

	// Важно: используем контекст Wails
	subID, err := a.subSvc.Submit(a.ctx, req)
	if err != nil {
		return map[string]string{"error": err.Error()}
	}

	return map[string]string{"submission_id": subID}
}

func (a *App) GetData(id string) map[string]interface{} {
	return map[string]interface{}{"status": "ok", "data": ""}
}

// NotifyVerdict отправляет событие во фронтенд (вызывается воркером)
func (a *App) NotifyVerdict(subID string) {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "on-verdict-changed", subID)
	}
}
