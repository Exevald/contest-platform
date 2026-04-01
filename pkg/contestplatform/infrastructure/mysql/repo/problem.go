package repo

import (
	"errors"

	"contest-platform/pkg/contestplatform/domain/model"
)

func NewProblemRepository() model.ProblemRepository {
	return &problemRepository{}
}

type problemRepository struct {
}

func (repo *problemRepository) NextID() model.ProblemID {
	panic("mysql repository is not implemented")
}

func (repo *problemRepository) List() ([]model.Problem, error) {
	return nil, errors.New("mysql repository is not implemented")
}

func (repo *problemRepository) Find(id model.ProblemID) (model.Problem, error) {
	return nil, errors.New("mysql repository is not implemented")
}

func (repo *problemRepository) Store(problem model.Problem) error {
	return errors.New("mysql repository is not implemented")
}
