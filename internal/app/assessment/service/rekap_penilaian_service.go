package service

import (
	"context"
	"errors"

	"github.com/BangNopall/paskihub-be/domain"
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

func (s *rekapService) GetTeamAssessmentDetail(ctx context.Context, regisID, userID uuid.UUID) (dto.TeamAssessmentDetailResponse, error) {
	authorized, err := s.repo.RegistrationBelongsToOrganizer(ctx, regisID, userID)
	if err != nil {
		return dto.TeamAssessmentDetailResponse{}, err
	}
	if !authorized {
		return dto.TeamAssessmentDetailResponse{}, domain.ErrForbidden
	}
	return s.repo.GetTeamAssessmentDetail(ctx, regisID)
}

func (s *rekapService) GetScoreboard(ctx context.Context, eventLevelID, userID uuid.UUID) (dto.ScoreboardResponse, error) {
	if err := s.authorizeEventLevel(ctx, eventLevelID, userID); err != nil {
		return dto.ScoreboardResponse{}, err
	}

	items, err := s.repo.GetScoreboardByEventLevel(ctx, eventLevelID)
	if err != nil {
		return dto.ScoreboardResponse{}, err
	}
	return dto.ScoreboardResponse{
		EventLevelID: eventLevelID,
		Items:        items,
	}, nil
}

func (s *rekapService) GetLeaderboardCustom(ctx context.Context, req dto.CustomLeaderboardRequest, eventLevelID, userID uuid.UUID) (dto.ScoreboardResponse, error) {
	if len(req.ScoreCategoryIDs) == 0 {
		return dto.ScoreboardResponse{}, errors.New("at least one score category id is required")
	}

	if err := s.authorizeEventLevel(ctx, eventLevelID, userID); err != nil {
		return dto.ScoreboardResponse{}, err
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
	if err := s.authorizeEventLevel(ctx, eventLevelID, userID); err != nil {
		return err
	}
	return s.repo.UpdateScorePublishedStatus(ctx, eventLevelID, req.IsScorePublished)
}

func (s *rekapService) authorizeEventLevel(ctx context.Context, eventLevelID, userID uuid.UUID) error {
	authorized, err := s.repo.EventLevelBelongsToOrganizer(ctx, eventLevelID, userID)
	if err != nil {
		return err
	}
	if !authorized {
		return domain.ErrForbidden
	}
	return nil
}
