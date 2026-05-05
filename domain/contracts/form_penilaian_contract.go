package contracts

import (
	"context"

	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FormPenilaianRepository interface {
	WithTx(tx *gorm.DB) FormPenilaianRepository
	GetSubCategoryGradeRules(ctx context.Context, subCategoryIDs []uuid.UUID) ([]entity.GradeRule, error)
	GetEventIDByRegisID(ctx context.Context, regisID uuid.UUID) (uuid.UUID, error)
	BulkUpsertScores(ctx context.Context, scores []entity.Score) error
	BulkInsertTeamViolations(ctx context.Context, violations []entity.TeamViolation) error
	UpdateAssessmentStatus(ctx context.Context, regisID uuid.UUID, status string) error
}

type FormPenilaianService interface {
	BulkInsertScores(ctx context.Context, req dto.BulkInsertScoresRequest) error
	BulkInsertTeamViolations(ctx context.Context, req dto.BulkInsertViolationsRequest) error
	FinalizeAssessment(ctx context.Context, req dto.FinalizeAssessmentRequest) error
}
