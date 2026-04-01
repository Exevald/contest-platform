package repo

import (
	"errors"

	"contest-platform/pkg/contestplatform/domain/model"
)

func NewSubmissionRepository() model.SubmissionRepository {
	return &submissionRepository{}
}

type submissionRepository struct {
}

func (repo *submissionRepository) NextID() model.SubmissionID {
	panic("mysql repository is not implemented")
}

func (repo *submissionRepository) Find(id model.SubmissionID) (model.Submission, error) {
	return nil, errors.New("mysql repository is not implemented")
}

func (repo *submissionRepository) Store(s model.Submission) error {
	return errors.New("mysql repository is not implemented")
}
