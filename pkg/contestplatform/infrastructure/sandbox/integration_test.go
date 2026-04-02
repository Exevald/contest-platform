package sandbox

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	appmodel "contest-platform/pkg/contestplatform/app/model"
	appservice "contest-platform/pkg/contestplatform/app/service"
	domainmodel "contest-platform/pkg/contestplatform/domain/model"
	domainservice "contest-platform/pkg/contestplatform/domain/service"
	sqliterepo "contest-platform/pkg/contestplatform/infrastructure/sqlite/repo"
)

func TestSubmissionRunsThroughSandboxAndGetsOKVerdict(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, err := sqliterepo.OpenDatabase(filepath.Join(t.TempDir(), "contestplatform.sqlite"))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	problemRepository := sqliterepo.NewProblemRepository(db)
	submissionRepository := sqliterepo.NewSubmissionRepository(db)

	problem, err := domainmodel.NewProblem("sum", "A+B", "sum problem", domainmodel.Constraints{
		TimeLimit:   5 * time.Second,
		MemoryLimit: 256 * 1024 * 1024,
	})
	if err != nil {
		t.Fatalf("new problem: %v", err)
	}
	if err = problem.AddTestCase("1 2\n", "3\n", true); err != nil {
		t.Fatalf("add testcase: %v", err)
	}
	if err = problemRepository.Store(problem); err != nil {
		t.Fatalf("store problem: %v", err)
	}

	gradingTasks := make(chan appmodel.GradingTask, 1)
	sandboxRunner, err := NewSandbox()
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	judgeService := domainservice.NewJudgeService()
	submissionDomainService := domainservice.NewSubmissionService(judgeService)
	submissionAppService := appservice.NewSubmissionService(
		submissionRepository,
		problemRepository,
		submissionDomainService,
		gradingTasks,
	)
	worker := appservice.NewGradingWorker(
		submissionRepository,
		problemRepository,
		sandboxRunner,
		judgeService,
		submissionDomainService,
		gradingTasks,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- worker.Run(ctx)
	}()

	submissionID, err := submissionAppService.Submit(ctx, appservice.SubmitRequest{
		ProblemID: "sum",
		Language:  "cpp",
		Source: `#include <iostream>

		int main() {
			int a = 0;
			int b = 0;
			std::cin >> a >> b;
			std::cout << (a + b) << '\n';
			return 0;
		}`,
	})
	if err != nil {
		t.Fatalf("submit: %v", err)
	}

	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		submission, findErr := submissionRepository.Find(domainmodel.SubmissionID(submissionID))
		if findErr != nil {
			t.Fatalf("find submission: %v", findErr)
		}
		if submission.Verdict().IsTerminal() {
			if submission.Verdict() != domainmodel.VerdictOK {
				t.Fatalf("expected OK verdict, got %s", submission.Verdict())
			}
			cancel()
			<-done
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	cancel()
	<-done
	t.Fatal("submission verdict did not reach terminal state in time")
}
