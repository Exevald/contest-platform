package service

import (
	"errors"

	"challenger/pkg/challenger/domain/model"
)

var (
	ErrSubmissionFinished = errors.New("submission evaluation is already finished")
)

type SubmissionService interface {
	ProcessNextResult(submission model.Submission, problem model.Problem, result model.TestResult, strategy model.EvaluationStrategy) bool
	CanSubmit(problem model.Problem, language model.Language, code string) error
}

type submissionService struct {
	judge JudgeService
}

func NewSubmissionService(judgeService JudgeService) SubmissionService {
	return &submissionService{judge: judgeService}
}

func (svc *submissionService) CanSubmit(problem model.Problem, language model.Language, code string) error {
	if len(code) == 0 {
		return errors.New("source code is empty")
	}
	if len(problem.TestCases()) == 0 {
		return errors.New("problem has no test cases")
	}
	return nil
}

func (svc *submissionService) ProcessNextResult(
	submission model.Submission,
	problem model.Problem,
	result model.TestResult,
	strategy model.EvaluationStrategy,
) bool {
	if submission.Verdict().IsTerminal() {
		return false
	}

	submission.AddTestResult(result)
	if strategy == model.StrategyStopOnFirstFail && result.Verdict != model.VerdictOK {
		return false
	}
	if len(submission.TestResults()) >= len(problem.TestCases()) {
		return false
	}

	return true
}
