package service

import (
	"context"

	appmodel "contest-platform/pkg/contestplatform/app/model"
	"contest-platform/pkg/contestplatform/domain/model"
	"contest-platform/pkg/contestplatform/domain/service"

	"github.com/sirupsen/logrus"
)

type GradingWorker interface {
	Run(ctx context.Context) error
	SetNotifyCallback(fn func(submissionID string))
}

func NewGradingWorker(
	submissionRepository model.SubmissionRepository,
	problemRepository model.ProblemRepository,
	sandbox appmodel.Sandbox,
	judgeService service.JudgeService,
	submissionService service.SubmissionService,
	gradingTasksChannel <-chan appmodel.GradingTask,
) GradingWorker {
	return &worker{
		submissionRepository: submissionRepository,
		problemRepository:    problemRepository,
		sandbox:              sandbox,
		judgeService:         judgeService,
		submissionService:    submissionService,
		gradingTasksChannel:  gradingTasksChannel,
	}
}

type worker struct {
	submissionRepository model.SubmissionRepository
	problemRepository    model.ProblemRepository
	sandbox              appmodel.Sandbox
	judgeService         service.JudgeService
	submissionService    service.SubmissionService
	gradingTasksChannel  <-chan appmodel.GradingTask
	notifyUI             func(string)
}

func (w *worker) SetNotifyCallback(fn func(submissionID string)) {
	w.notifyUI = fn
}

func (w *worker) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case task, ok := <-w.gradingTasksChannel:
			if !ok {
				return nil
			}
			if err := w.processTask(ctx, task); err != nil {
				return err
			}
		}
	}
}

func (w *worker) processTask(ctx context.Context, task appmodel.GradingTask) error {
	if w.notifyUI != nil {
		defer w.notifyUI(string(task.SubmissionID))
	}

	submission, err := w.submissionRepository.Find(task.SubmissionID)
	if err != nil {
		return err
	}
	problem, err := w.problemRepository.Find(task.ProblemID)
	if err != nil {
		return err
	}

	submission.UpdateVerdict(model.VerdictCompiling)
	if err = w.submissionRepository.Store(submission); err != nil {
		return err
	}

	workingDir, exePath, err := w.sandbox.Prepare(ctx, submission.Language(), submission.SourceCode())
	if err != nil {
		submission.SetCompilationOutput(err.Error())
		submission.UpdateVerdict(model.VerdictCE)
		return w.submissionRepository.Store(submission)
	}

	submission.SetCompilationOutput("Compilation finished successfully.")

	submission.UpdateVerdict(model.VerdictRunning)
	if err = w.submissionRepository.Store(submission); err != nil {
		return err
	}

	defer func() {
		cleanupErr := w.sandbox.Cleanup(workingDir)
		if cleanupErr != nil {
			logrus.Errorf("Error cleaning up submission dir %s: %v", workingDir, cleanupErr)
		}
	}()

	for _, tc := range problem.TestCases() {
		sandboxResponse, err2 := w.sandbox.Execute(ctx, workingDir, exePath, tc.Input, problem.Constraints())
		if err2 != nil {
			submission.UpdateVerdict(model.VerdictInternal)
			break
		}

		execRes := service.ExecutionResult{
			Stdout:     sandboxResponse.Stdout,
			Stderr:     sandboxResponse.Stderr,
			TimeUsed:   sandboxResponse.TimeUsed,
			MemoryUsed: sandboxResponse.MemoryUsed,
			ExitCode:   sandboxResponse.ExitCode,
			TimedOut:   sandboxResponse.IsTimeout,
		}

		verdict := w.judgeService.Verify(tc, problem.Constraints(), execRes)

		testResult := model.TestResult{
			TestCaseID: tc.ID,
			Verdict:    verdict,
			TimeUsed:   sandboxResponse.TimeUsed,
			MemoryUsed: sandboxResponse.MemoryUsed,
		}
		shouldContinue := w.submissionService.ProcessNextResult(
			submission,
			problem,
			testResult,
			model.StrategyStopOnFirstFail,
		)

		err = w.submissionRepository.Store(submission)
		if err != nil {
			return err
		}

		if !shouldContinue {
			break
		}
	}
	return w.submissionRepository.Store(submission)
}
