package service

import (
	"strings"
	"time"

	"challenger/pkg/challenger/domain/model"
)

type ExecutionResult struct {
	Stdout     string
	Stderr     string
	TimeUsed   time.Duration
	MemoryUsed uint64
	ExitCode   int
	TimedOut   bool
}

type JudgeService interface {
	Verify(test model.TestCase, constraints model.Constraints, run ExecutionResult) model.Verdict
}

type judgeService struct {
}

func NewJudgeService() JudgeService {
	return &judgeService{}
}

func (j *judgeService) Verify(
	test model.TestCase,
	constraints model.Constraints,
	result ExecutionResult,
) model.Verdict {
	if result.TimedOut || result.TimeUsed > constraints.TimeLimit {
		return model.VerdictTLE
	}
	if result.MemoryUsed > constraints.MemoryLimit {
		return model.VerdictMLE
	}
	if result.ExitCode != 0 {
		return model.VerdictRE
	}

	if j.compareStrings(result.Stdout, test.ExpectedOutput) {
		return model.VerdictOK
	}

	return model.VerdictWA
}

func (j *judgeService) compareStrings(actual, expected string) bool {
	actualLines := strings.Split(strings.TrimSpace(actual), "\n")
	expectedLines := strings.Split(strings.TrimSpace(expected), "\n")

	if len(actualLines) != len(expectedLines) {
		return false
	}

	for i := range actualLines {
		if strings.TrimRight(actualLines[i], " \r\t") != strings.TrimRight(expectedLines[i], " \r\t") {
			return false
		}
	}

	return true
}
