package contracts

import (
	"context"

	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/google/uuid"
)

type RekapRepository interface {
	RegistrationBelongsToOrganizer(ctx context.Context, regisID, organizerID uuid.UUID) (bool, error)
	EventLevelBelongsToOrganizer(ctx context.Context, eventLevelID, organizerID uuid.UUID) (bool, error)
	GetTeamAssessmentDetail(ctx context.Context, regisID uuid.UUID) (dto.TeamAssessmentDetailResponse, error)
	GetScoreboardByEventLevel(ctx context.Context, eventLevelID uuid.UUID) ([]dto.ScoreboardItem, error)
	GetLeaderboardCustom(ctx context.Context, eventLevelID uuid.UUID, categoryIDs []uuid.UUID) ([]dto.ScoreboardItem, error)
	UpdateScorePublishedStatus(ctx context.Context, eventLevelID uuid.UUID, isPublished bool) error
	GetEventIDByEventLevelID(ctx context.Context, eventLevelID uuid.UUID) (uuid.UUID, error)
}

type RekapService interface {
	GetTeamAssessmentDetail(ctx context.Context, regisID, userID uuid.UUID) (dto.TeamAssessmentDetailResponse, error)
	GetScoreboard(ctx context.Context, eventLevelID, userID uuid.UUID) (dto.ScoreboardResponse, error)
	GetLeaderboardCustom(ctx context.Context, req dto.CustomLeaderboardRequest, eventLevelID, userID uuid.UUID) (dto.ScoreboardResponse, error)
	PublishScoreboard(ctx context.Context, req dto.PublishScoreboardRequest, eventLevelID uuid.UUID, userID uuid.UUID) error
}
