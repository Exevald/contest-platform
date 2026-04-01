package model

import (
	"errors"
	"testing"
	"time"
)

func TestNewProblem(t *testing.T) {
	limits := Constraints{TimeLimit: time.Second, MemoryLimit: 256 * 1024}

	t.Run("should create valid problem", func(t *testing.T) {
		p, err := NewProblem("PROB-1", "A+B", "Sum", limits)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.Title() != "A+B" {
			t.Errorf("expected title A+B, got %v", p.Title())
		}
	})

	t.Run("should fail on empty title", func(t *testing.T) {
		_, err := NewProblem("PROB-1", "", "Sum", limits)
		if !errors.Is(err, ErrEmptyTitle) {
			t.Errorf("expected ErrEmptyTitle, got %v", err)
		}
	})
}

func TestSubmission_VerdictLogic(t *testing.T) {
	sub := NewSubmission("SUB-1", "PROB-1", "go", "main...")

	t.Run("should start with pending", func(t *testing.T) {
		if sub.Verdict() != VerdictPending {
			t.Errorf("expected PENDING, got %s", sub.Verdict())
		}
	})

	t.Run("should set OK on first successful test", func(t *testing.T) {
		sub.AddTestResult(TestResult{TestCaseID: 1, Verdict: VerdictOK})
		if sub.Verdict() != VerdictOK {
			t.Errorf("expected OK, got %s", sub.Verdict())
		}
	})

	t.Run("should lock on first error", func(t *testing.T) {
		sub.AddTestResult(TestResult{TestCaseID: 2, Verdict: VerdictWA})
		if sub.Verdict() != VerdictWA {
			t.Errorf("expected WA, got %s", sub.Verdict())
		}

		sub.AddTestResult(TestResult{TestCaseID: 3, Verdict: VerdictOK})
		if sub.Verdict() != VerdictWA {
			t.Errorf("expected verdict to stay WA, but got %s", sub.Verdict())
		}
	})
}

func TestProblem_AddTestCase(t *testing.T) {
	p, _ := NewProblem("P1", "T", "D", Constraints{time.Second, 1024})

	err := p.AddTestCase("1 2", "3", true)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(p.TestCases()) != 1 {
		t.Errorf("expected 1 test case, got %d", len(p.TestCases()))
	}

	tc := p.TestCases()[0]
	if tc.Input != "1 2" || !tc.IsSample {
		t.Error("test case data mismatch")
	}
}
