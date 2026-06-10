package service

import (
	"context"
	"errors"
	"testing"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type stubFormPenilaianRepo struct {
	authorized      bool
	gradeRulesRead  *bool
	scoresWritten   *bool
	violationsWrote *bool
	statusUpdated   *bool
}

func (r stubFormPenilaianRepo) WithTx(tx *gorm.DB) contracts.FormPenilaianRepository {
	return r
}

func (r stubFormPenilaianRepo) ValidateAssessmentOwnership(
	ctx context.Context,
	organizerID uuid.UUID,
	regisID uuid.UUID,
	judgeID uuid.UUID,
	subCategoryIDs []uuid.UUID,
	violationTypeIDs []uuid.UUID,
) (bool, error) {
	return r.authorized, nil
}

func (r stubFormPenilaianRepo) GetSubCategoryGradeRules(ctx context.Context, subCategoryIDs []uuid.UUID) ([]entity.GradeRule, error) {
	if r.gradeRulesRead != nil {
		*r.gradeRulesRead = true
	}
	return nil, nil
}

func (r stubFormPenilaianRepo) GetEventIDByRegisID(ctx context.Context, regisID uuid.UUID) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (r stubFormPenilaianRepo) BulkUpsertScores(ctx context.Context, scores []entity.Score) error {
	if r.scoresWritten != nil {
		*r.scoresWritten = true
	}
	return nil
}

func (r stubFormPenilaianRepo) BulkInsertTeamViolations(ctx context.Context, violations []entity.TeamViolation) error {
	if r.violationsWrote != nil {
		*r.violationsWrote = true
	}
	return nil
}

func (r stubFormPenilaianRepo) UpdateAssessmentStatus(ctx context.Context, regisID uuid.UUID, status string) error {
	if r.statusUpdated != nil {
		*r.statusUpdated = true
	}
	return nil
}

func TestAssessmentWritesRejectForeignObjectsBeforePersistence(t *testing.T) {
	organizerID := uuid.New()
	regisID := uuid.New()
	judgeID := uuid.New()
	subCategoryID := uuid.New()
	violationTypeID := uuid.New()

	tests := []struct {
		name string
		run  func(*formPenilaianService) error
	}{
		{
			name: "bulk scores",
			run: func(svc *formPenilaianService) error {
				return svc.BulkInsertScores(context.Background(), organizerID, dto.BulkInsertScoresRequest{
					RegisID:  regisID,
					JudgesID: judgeID,
					Scores: []dto.ScoreInput{
						{SubCategoryID: subCategoryID, ScoreValue: 80},
					},
				})
			},
		},
		{
			name: "bulk violations",
			run: func(svc *formPenilaianService) error {
				return svc.BulkInsertTeamViolations(context.Background(), organizerID, dto.BulkInsertViolationsRequest{
					RegisID:          regisID,
					JudgesID:         judgeID,
					ViolationTypeIDs: []uuid.UUID{violationTypeID},
				})
			},
		},
		{
			name: "finalize",
			run: func(svc *formPenilaianService) error {
				return svc.FinalizeAssessment(context.Background(), organizerID, dto.FinalizeAssessmentRequest{
					RegisID:  regisID,
					JudgesID: judgeID,
					Scores: []dto.ScoreInput{
						{SubCategoryID: subCategoryID, ScoreValue: 80},
					},
					ViolationTypeIDs: []uuid.UUID{violationTypeID},
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gradeRulesRead := false
			scoresWritten := false
			violationsWrote := false
			statusUpdated := false
			svc := &formPenilaianService{
				repo: stubFormPenilaianRepo{
					authorized:      false,
					gradeRulesRead:  &gradeRulesRead,
					scoresWritten:   &scoresWritten,
					violationsWrote: &violationsWrote,
					statusUpdated:   &statusUpdated,
				},
			}

			err := tt.run(svc)
			if !errors.Is(err, domain.ErrForbidden) {
				t.Fatalf("expected ErrForbidden, got %v", err)
			}
			if gradeRulesRead || scoresWritten || violationsWrote || statusUpdated {
				t.Fatal("assessment persistence was reached for foreign objects")
			}
		})
	}
}
