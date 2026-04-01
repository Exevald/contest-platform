package ui

import (
	"encoding/json"
	"fmt"

	"github.com/energye/energy/v2/cef/ipc"
	"github.com/energye/energy/v2/cef/ipc/context"

	"challenger/pkg/challenger/app/service"
	"challenger/pkg/challenger/domain/model"
)

type SendFileArgs struct {
	ID       string `json:"id"`
	Language string `json:"language"`
	Text     string `json:"text"`
}

type GetDataArgs struct {
	ID string `json:"id"`
}

type ResetTaskArgs struct {
	ID string `json:"id"`
}

type EnergyBridge struct {
	subSvc service.SubmissionService
}

func NewEnergyBridge(subSvc service.SubmissionService) *EnergyBridge {
	return &EnergyBridge{
		subSvc: subSvc,
	}
}

func (b *EnergyBridge) RegisterHandlers() {
	// 1. Platform_getStartupData
	ipc.On("Platform_getStartupData", func(ctx context.IContext) {
		var problems []model.Problem

		data := StartupData{
			Title:     "iSpring Challenger",
			Languages: []string{"cpp", "python", "go", "pascal", "js", "php"},
			Tabs:      problems,
		}
		ctx.Result(data)
	})

	// 2. sendFile (отправка решения)
	ipc.On("sendFile", func(ctx context.Context) {
		argsJSON := ctx.ArgumentList().GetStringByIndex(0)
		var args SendFileArgs
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			ctx.Result(fmt.Errorf("invalid args: %v", err))
			return
		}

		// Маппим SendFileArgs -> SubmitRequest
		req := service.SubmitRequest{
			ProblemID: args.ID,
			Language:  args.Language,
			Source:    args.Text,
		}

		id, err := b.subSvc.Submit(nil, req)
		if err != nil {
			ctx.Result(map[string]string{"error": err.Error()})
			return
		}

		ctx.Result(map[string]string{"submission_id": id})
	})

	// 3. getData (получение статуса или данных задачи)
	ipc.On("getData", func(ctx context.IContext) {
		argsJSON := ctx.ArgumentList().GetStringByIndex(0)
		var args GetDataArgs
		json.Unmarshal([]byte(argsJSON), &args)

		// Логика получения данных...
		ctx.Result(map[string]interface{}{"status": "ok"})
	})

	// 4. resetTask
	ipc.On("resetTask", func(ctx context.Context) {
		argsJSON := ctx.ArgumentList().GetStringByIndex(0)
		var args ResetTaskArgs
		json.Unmarshal([]byte(argsJSON), &args)

		// Логика сброса...
		ctx.Result(true)
	})
}
