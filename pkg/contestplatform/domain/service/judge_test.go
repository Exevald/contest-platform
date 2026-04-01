package service

import (
	"testing"
	"time"

	"contest-platform/pkg/contestplatform/domain/model"
)

func TestJudgeService_Verify(t *testing.T) {
	svc := NewJudgeService()
	testCase := model.TestCase{ExpectedOutput: "42\n"}
	limits := model.Constraints{TimeLimit: time.Second, MemoryLimit: 1024 * 1024}

	t.Run("successful match", func(t *testing.T) {
		res := ExecutionResult{Stdout: "42", ExitCode: 0, TimeUsed: 10 * time.Millisecond}
		if v := svc.Verify(testCase, limits, res); v != model.VerdictOK {
			t.Errorf("expected OK, got %s", v)
		}
	})

	t.Run("trailing spaces should be ignored", func(t *testing.T) {
		res := ExecutionResult{Stdout: "42   \n", ExitCode: 0}
		if v := svc.Verify(testCase, limits, res); v != model.VerdictOK {
			t.Errorf("expected OK for output with spaces, got %s", v)
		}
	})

	t.Run("time limit exceeded", func(t *testing.T) {
		res := ExecutionResult{TimeUsed: 2 * time.Second, ExitCode: 0}
		if v := svc.Verify(testCase, limits, res); v != model.VerdictTLE {
			t.Errorf("expected TLE, got %s", v)
		}
	})

	t.Run("runtime error", func(t *testing.T) {
		res := ExecutionResult{ExitCode: 1, Stdout: "partial output"}
		if v := svc.Verify(testCase, limits, res); v != model.VerdictRE {
			t.Errorf("expected RE, got %s", v)
		}
	})

	t.Run("wrong answer", func(t *testing.T) {
		res := ExecutionResult{Stdout: "wrong", ExitCode: 0}
		if v := svc.Verify(testCase, limits, res); v != model.VerdictWA {
			t.Errorf("expected WA, got %s", v)
		}
	})
}
