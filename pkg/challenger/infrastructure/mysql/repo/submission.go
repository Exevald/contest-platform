package repo

import (
	"challenger/pkg/challenger/domain/model"
)

func NewSubmissionRepository() model.SubmissionRepository {
	return &submissionRepository{}
}

type submissionRepository struct {
}

func (repo *submissionRepository) NextID() model.SubmissionID {
	panic("implement me")
}

func (repo *submissionRepository) Find(id model.SubmissionID) (model.Submission, error) {
	panic("implement me")
}

func (repo *submissionRepository) Store(s model.Submission) error {
	panic("implement me")
}
