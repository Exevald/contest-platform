package model

import "time"

type EvaluationStrategy string

const (
	StrategyStopOnFirstFail EvaluationStrategy = "STOP_ON_FIRST_FAIL"
	StrategyRunAllTests     EvaluationStrategy = "RUN_ALL_TESTS"
)

type TestResult struct {
	TestCaseID int
	Verdict    Verdict
	TimeUsed   time.Duration
	MemoryUsed uint64
}

type SubmissionRepository interface {
	NextID() SubmissionID
	Find(id SubmissionID) (Submission, error)
	Store(s Submission) error
}

type Submission interface {
	ID() SubmissionID
	ProblemID() ProblemID
	Language() Language
	SourceCode() string
	Verdict() Verdict
	TestResults() []TestResult
	CreatedAt() time.Time

	AddTestResult(res TestResult)
	UpdateVerdict(v Verdict)
}

type submission struct {
	id          SubmissionID
	problemID   ProblemID
	language    Language
	sourceCode  string
	verdict     Verdict
	testResults []TestResult
	createdAt   time.Time
}

func NewSubmission(id SubmissionID, pID ProblemID, lang Language, code string) Submission {
	return &submission{
		id:          id,
		problemID:   pID,
		language:    lang,
		sourceCode:  code,
		verdict:     VerdictPending,
		testResults: make([]TestResult, 0),
		createdAt:   time.Now(),
	}
}

func (s *submission) ID() SubmissionID          { return s.id }
func (s *submission) ProblemID() ProblemID      { return s.problemID }
func (s *submission) Language() Language        { return s.language }
func (s *submission) SourceCode() string        { return s.sourceCode }
func (s *submission) Verdict() Verdict          { return s.verdict }
func (s *submission) TestResults() []TestResult { return s.testResults }
func (s *submission) CreatedAt() time.Time      { return s.createdAt }

func (s *submission) AddTestResult(res TestResult) {
	s.testResults = append(s.testResults, res)

	if s.verdict == VerdictPending || s.verdict == VerdictRunning || s.verdict == VerdictOK {
		s.verdict = res.Verdict
	}
}

func (s *submission) UpdateVerdict(v Verdict) {
	s.verdict = v
}
