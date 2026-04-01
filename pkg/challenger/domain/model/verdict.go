package model

import "time"

type Verdict string

const (
	VerdictPending   Verdict = "PENDING"
	VerdictCompiling Verdict = "COMPILING"
	VerdictRunning   Verdict = "RUNNING"
	VerdictOK        Verdict = "OK"
	VerdictWA        Verdict = "WA"
	VerdictTLE       Verdict = "TLE"
	VerdictMLE       Verdict = "MLE"
	VerdictRE        Verdict = "RE"
	VerdictCE        Verdict = "CE"
	VerdictInternal  Verdict = "SE"
)

type Language string
type ProblemID string
type SubmissionID string

type Constraints struct {
	TimeLimit   time.Duration
	MemoryLimit uint64
}

func (v Verdict) IsTerminal() bool {
	return v != VerdictPending && v != VerdictCompiling && v != VerdictRunning
}
