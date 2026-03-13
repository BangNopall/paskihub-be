package service

import (
	"context"
	"errors"

	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/google/uuid"
)

type rekapService struct {
	repo contracts.RekapRepository
}

func NewRekapService(repo contracts.RekapRepository) contracts.RekapService {
	return &rekapService{repo: repo}
}

func (s *rekapService) GetTeamAssessmentDetail(ctx context.Context, regisID uuid.UUID) (dto.TeamAssessmentDetailResponse, error) {
	return s.repo.GetTeamAssessmentDetail(ctx, regisID)
}

func (s *rekapService) GetScoreboard(ctx context.Context, eventLevelID uuid.UUID) (dto.ScoreboardResponse, error) {
	items, err := s.repo.GetScoreboardByEventLevel(ctx, eventLevelID)
	if err != nil {
		return dto.ScoreboardResponse{}, err
	}
	return dto.ScoreboardResponse{
		EventLevelID: eventLevelID,
		Items:        items,
	}, nil
}

func (s *rekapService) GetLeaderboardCustom(ctx context.Context, req dto.CustomLeaderboardRequest, eventLevelID uuid.UUID) (dto.ScoreboardResponse, error) {
	if len(req.ScoreCategoryIDs) == 0 {
		return dto.ScoreboardResponse{}, errors.New("at least one score category id is required")
	}

	items, err := s.repo.GetLeaderboardCustom(ctx, eventLevelID, req.ScoreCategoryIDs)
	if err != nil {
		return dto.ScoreboardResponse{}, err
	}
	return dto.ScoreboardResponse{
		EventLevelID: eventLevelID,
		Items:        items,
	}, nil
}

func (s *rekapService) PublishScoreboard(ctx context.Context, req dto.PublishScoreboardRequest, eventLevelID uuid.UUID, userID uuid.UUID) error {
	// Authorization checking (EO owns this event) should be placed in Controller by invoking EventService.
	return s.repo.UpdateScorePublishedStatus(ctx, eventLevelID, req.IsScorePublished)
}
