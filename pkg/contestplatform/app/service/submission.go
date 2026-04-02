package service

import (
	"context"
	"fmt"

	appmodel "contest-platform/pkg/contestplatform/app/model"
	"contest-platform/pkg/contestplatform/domain/model"
	"contest-platform/pkg/contestplatform/domain/service"
)

type SubmitRequest struct {
	ProblemID       string
	ParticipantCode string
	Language        string
	Source          string
}

type SubmissionService interface {
	Submit(ctx context.Context, request SubmitRequest) (string, error)
}

func NewSubmissionService(
	submissionRepository model.SubmissionRepository,
	problemRepository model.ProblemRepository,
	domainService service.SubmissionService,
	gradingTasksChannel chan<- appmodel.GradingTask,
) SubmissionService {
	return &submissionService{
		submissionRepository: submissionRepository,
		problemRepository:    problemRepository,
		domainService:        domainService,
		gradingTasksChannel:  gradingTasksChannel,
	}
}

type submissionService struct {
	submissionRepository model.SubmissionRepository
	problemRepository    model.ProblemRepository
	domainService        service.SubmissionService
	gradingTasksChannel  chan<- appmodel.GradingTask
}

func (s *submissionService) Submit(ctx context.Context, request SubmitRequest) (string, error) {
	prob, err := s.problemRepository.Find(model.ProblemID(request.ProblemID))
	if err != nil {
		return "", fmt.Errorf("problem not found: %w", err)
	}

	if _, ok := appmodel.Languages[model.Language(request.Language)]; !ok {
		return "", fmt.Errorf("validation failed: language %q is not supported", request.Language)
	}

	err = s.domainService.CanSubmit(prob, model.Language(request.Language), request.Source)
	if err != nil {
		return "", fmt.Errorf("validation failed: %w", err)
	}

	subID := s.submissionRepository.NextID()
	sub := model.NewSubmission(
		subID,
		prob.ID(),
		request.ParticipantCode,
		model.Language(request.Language),
		request.Source,
	)

	if err = s.submissionRepository.Store(sub); err != nil {
		return "", fmt.Errorf("failed to store submission: %w", err)
	}

	s.gradingTasksChannel <- appmodel.GradingTask{
		SubmissionID: subID,
		ProblemID:    prob.ID(),
	}

	return string(subID), nil
}
