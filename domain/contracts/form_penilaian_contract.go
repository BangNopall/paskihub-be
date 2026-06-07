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
	ValidateAssessmentOwnership(ctx context.Context, organizerID, regisID, judgeID uuid.UUID, subCategoryIDs, violationTypeIDs []uuid.UUID) (bool, error)
	GetSubCategoryGradeRules(ctx context.Context, subCategoryIDs []uuid.UUID) ([]entity.GradeRule, error)
	GetEventIDByRegisID(ctx context.Context, regisID uuid.UUID) (uuid.UUID, error)
	BulkUpsertScores(ctx context.Context, scores []entity.Score) error
	BulkInsertTeamViolations(ctx context.Context, violations []entity.TeamViolation) error
	UpdateAssessmentStatus(ctx context.Context, regisID uuid.UUID, status string) error
}

type FormPenilaianService interface {
	BulkInsertScores(ctx context.Context, organizerID uuid.UUID, req dto.BulkInsertScoresRequest) error
	BulkInsertTeamViolations(ctx context.Context, organizerID uuid.UUID, req dto.BulkInsertViolationsRequest) error
	FinalizeAssessment(ctx context.Context, organizerID uuid.UUID, req dto.FinalizeAssessmentRequest) error
}
