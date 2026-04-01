package model

import (
	"context"
	"time"

	"challenger/pkg/challenger/domain/model"
)

type SandboxRequest struct {
	ID          string
	SourceCode  string
	Language    model.Language
	Input       string
	TimeLimit   time.Duration
	MemoryLimit uint64
}

type SandboxResponse struct {
	Stdout     string
	Stderr     string
	ExitCode   int
	TimeUsed   time.Duration
	MemoryUsed uint64
	IsTimeout  bool
}

type GradingTask struct {
	SubmissionID model.SubmissionID
	ProblemID    model.ProblemID
}

type Sandbox interface {
	Prepare(ctx context.Context, lang model.Language, sourceCode string) (workingDir string, exePath string, err error)
	Execute(ctx context.Context, workingDir string, exePath string, input string, limits model.Constraints) (SandboxResponse, error)
	Cleanup(workingDir string) error
}
