package repo

import (
	"challenger/pkg/challenger/domain/model"
)

func NewProblemRepository() model.ProblemRepository {
	return &problemRepository{}
}

type problemRepository struct {
}

func (repo *problemRepository) NextID() model.ProblemID {
	//TODO implement me
	panic("implement me")
}

func (repo *problemRepository) Find(id model.ProblemID) (model.Problem, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *problemRepository) Store(problem model.Problem) error {
	//TODO implement me
	panic("implement me")
}
