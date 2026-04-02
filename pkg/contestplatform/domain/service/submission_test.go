package service

import (
	"testing"
	"time"

	"contest-platform/pkg/contestplatform/domain/model"
)

func TestSubmissionService_ProcessNextResult(t *testing.T) {
	judge := NewJudgeService()
	svc := NewSubmissionService(judge)

	p, _ := model.NewProblem("P1", "Title", "Desc", model.Constraints{TimeLimit: time.Second, MemoryLimit: 1024})
	_ = p.AddTestCase("1", "1", false)
	_ = p.AddTestCase("2", "2", false)
	_ = p.AddTestCase("3", "3", false)

	t.Run("should stop on first fail with STOP_ON_FIRST_FAIL strategy", func(t *testing.T) {
		s := model.NewSubmission("S1", p.ID(), "participant-1", "go", "code")

		res := model.TestResult{TestCaseID: 1, Verdict: model.VerdictWA}
		continueTesting := svc.ProcessNextResult(s, p, res, model.StrategyStopOnFirstFail)

		if continueTesting {
			t.Error("expected to stop testing after WA")
		}
		if s.Verdict() != model.VerdictWA {
			t.Errorf("expected verdict WA, got %s", s.Verdict())
		}
	})

	t.Run("should continue on fail with RUN_ALL_TESTS strategy", func(t *testing.T) {
		s := model.NewSubmission("S2", p.ID(), "participant-1", "go", "code")

		res := model.TestResult{TestCaseID: 1, Verdict: model.VerdictWA}
		continueTesting := svc.ProcessNextResult(s, p, res, model.StrategyRunAllTests)

		if !continueTesting {
			t.Error("expected to continue testing after WA in RUN_ALL_TESTS mode")
		}
	})

	t.Run("should stop when all tests finished", func(t *testing.T) {
		s := model.NewSubmission("S3", p.ID(), "participant-1", "go", "code")

		svc.ProcessNextResult(s, p, model.TestResult{TestCaseID: 1, Verdict: model.VerdictOK}, model.StrategyRunAllTests)
		svc.ProcessNextResult(s, p, model.TestResult{TestCaseID: 2, Verdict: model.VerdictOK}, model.StrategyRunAllTests)
		last := svc.ProcessNextResult(s, p, model.TestResult{TestCaseID: 3, Verdict: model.VerdictOK}, model.StrategyRunAllTests)

		if last {
			t.Error("expected to stop after the last test case")
		}
	})
}

func TestSubmissionService_CanSubmit(t *testing.T) {
	judge := NewJudgeService()
	svc := NewSubmissionService(judge)

	p, _ := model.NewProblem("P1", "Title", "Desc", model.Constraints{TimeLimit: time.Second, MemoryLimit: 1024})
	_ = p.AddTestCase("1", "1", true)

	t.Run("accepts cpp", func(t *testing.T) {
		if err := svc.CanSubmit(p, "cpp", "int main(){}"); err != nil {
			t.Fatalf("expected cpp to be accepted, got %v", err)
		}
	})

	t.Run("accepts other languages too", func(t *testing.T) {
		if err := svc.CanSubmit(p, "python", "print(1)"); err != nil {
			t.Fatalf("expected python to be accepted, got %v", err)
		}
	})
}
