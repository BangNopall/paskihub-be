package contracts

import (
	"context"

	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/google/uuid"
)

type ParticipantAssessmentRepository interface {
	GetScoresByRegisID(ctx context.Context, regisID uuid.UUID) ([]entity.Score, error)
	GetViolationsByRegisID(ctx context.Context, regisID uuid.UUID) ([]entity.TeamViolation, error)
	GetRegistrationOwnership(ctx context.Context, regisID uuid.UUID, userID uuid.UUID) (bool, error)
	GetRegistrationByID(ctx context.Context, regisID uuid.UUID) (*entity.Registration, error)
	GetScoreCategoriesByEventID(ctx context.Context, eventID uuid.UUID) ([]entity.ScoreCategory, error)
}

type ParticipantAssessmentService interface {
	GetAssessmentRecap(ctx context.Context, userID string, regisID string) (*dto.AssessmentRecapResponse, error)
}
