package service

import (
	"context"
	"strings"
	"testing"
	"time"

	appmodel "contest-platform/pkg/contestplatform/app/model"
	domainmodel "contest-platform/pkg/contestplatform/domain/model"
	domainservice "contest-platform/pkg/contestplatform/domain/service"
)

func TestSubmissionService_SubmitRejectsUnsupportedLanguage(t *testing.T) {
	problem, err := domainmodel.NewProblem(
		"P1",
		"Title",
		"Desc",
		domainmodel.Constraints{TimeLimit: time.Second, MemoryLimit: 1024},
	)
	if err != nil {
		t.Fatalf("new problem: %v", err)
	}
	if err = problem.AddTestCase("1", "1", true); err != nil {
		t.Fatalf("add testcase: %v", err)
	}

	problemRepo := &mockProblemRepository{problem: problem}
	submissionRepo := &mockSubmissionRepository{}
	svc := NewSubmissionService(
		submissionRepo,
		problemRepo,
		domainservice.NewSubmissionService(domainservice.NewJudgeService()),
		make(chan appmodel.GradingTask, 1),
	)

	_, err = svc.Submit(context.Background(), SubmitRequest{
		ProblemID: string(problem.ID()),
		Language:  "ruby",
		Source:    "puts 1",
	})
	if err == nil {
		t.Fatal("expected unsupported language error")
	}
	if !strings.Contains(err.Error(), "not supported") {
		t.Fatalf("expected not supported error, got %v", err)
	}
}

type mockProblemRepository struct {
	problem domainmodel.Problem
}

func (repo *mockProblemRepository) NextID() domainmodel.ProblemID {
	return "problem-1"
}

func (repo *mockProblemRepository) List() ([]domainmodel.Problem, error) {
	return []domainmodel.Problem{repo.problem}, nil
}

func (repo *mockProblemRepository) Find(_ domainmodel.ProblemID) (domainmodel.Problem, error) {
	return repo.problem, nil
}

func (repo *mockProblemRepository) Store(problem domainmodel.Problem) error {
	repo.problem = problem
	return nil
}

type mockSubmissionRepository struct{}

func (repo *mockSubmissionRepository) NextID() domainmodel.SubmissionID {
	return "submission-1"
}

func (repo *mockSubmissionRepository) Find(_ domainmodel.SubmissionID) (domainmodel.Submission, error) {
	return nil, nil
}

func (repo *mockSubmissionRepository) FindLatest(_ domainmodel.ProblemID) (domainmodel.Submission, error) {
	return nil, nil
}

func (repo *mockSubmissionRepository) ListByProblem(_ domainmodel.ProblemID) ([]domainmodel.Submission, error) {
	return nil, nil
}

func (repo *mockSubmissionRepository) Store(s domainmodel.Submission) error {
	return nil
}
