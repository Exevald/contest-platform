package model

import "errors"

var (
	ErrEmptyTitle       = errors.New("problem title cannot be empty")
	ErrInvalidThreshold = errors.New("invalid limits: time or memory cannot be zero")
)

type TestCase struct {
	ID             int
	Input          string
	ExpectedOutput string
	IsSample       bool
}

type ProblemRepository interface {
	NextID() ProblemID
	Find(id ProblemID) (Problem, error)
	Store(problem Problem) error
}

type Problem interface {
	ID() ProblemID
	Title() Title
	Description() string
	Constraints() Constraints
	TestCases() []TestCase

	AddTestCase(input, output string, isSample bool) error
	UpdateDescription(desc string)
}

type Title string

type problem struct {
	id          ProblemID
	title       Title
	description string
	constraints Constraints
	testCases   []TestCase
}

func NewProblem(id ProblemID, title Title, desc string, limits Constraints) (Problem, error) {
	if title == "" {
		return nil, ErrEmptyTitle
	}
	if limits.TimeLimit <= 0 || limits.MemoryLimit <= 0 {
		return nil, ErrInvalidThreshold
	}

	return &problem{
		id:          id,
		title:       title,
		description: desc,
		constraints: limits,
		testCases:   make([]TestCase, 0),
	}, nil
}

func (p *problem) ID() ProblemID            { return p.id }
func (p *problem) Title() Title             { return p.title }
func (p *problem) Description() string      { return p.description }
func (p *problem) Constraints() Constraints { return p.constraints }
func (p *problem) TestCases() []TestCase    { return p.testCases }

func (p *problem) AddTestCase(input, output string, isSample bool) error {
	p.testCases = append(p.testCases, TestCase{
		ID:             len(p.testCases) + 1,
		Input:          input,
		ExpectedOutput: output,
		IsSample:       isSample,
	})
	return nil
}

func (p *problem) UpdateDescription(desc string) {
	p.description = desc
}
