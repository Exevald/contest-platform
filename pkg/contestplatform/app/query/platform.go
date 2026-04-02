package query

import "time"

type ProblemListItem struct {
	ID    string
	Title string
}

type SubmissionView struct {
	ID                string
	ProblemID         string
	Language          string
	Verdict           string
	CompilationOutput string
	CreatedAt         time.Time
}

type PlatformQueryService interface {
	ListProblems() ([]ProblemListItem, error)
	GetProblemDescription(problemID string) (string, error)
	GetSubmissionStatus(submissionID string) (*SubmissionView, error)
	ListSubmissionHistory(problemID string) ([]SubmissionView, error)
}
