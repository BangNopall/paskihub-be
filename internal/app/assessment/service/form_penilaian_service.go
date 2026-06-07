package service

import (
	"context"
	"errors"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type formPenilaianService struct {
	repo     contracts.FormPenilaianRepository
	db       *gorm.DB
	validate *validator.Validate
}

func NewFormPenilaianService(repo contracts.FormPenilaianRepository, db *gorm.DB, validate *validator.Validate) contracts.FormPenilaianService {
	return &formPenilaianService{
		repo:     repo,
		db:       db,
		validate: validate,
	}
}

func (s *formPenilaianService) BulkInsertScores(ctx context.Context, organizerID uuid.UUID, req dto.BulkInsertScoresRequest) error {
	if s.validate != nil {
		if err := s.validate.Struct(req); err != nil {
			return err
		}
	}

	var subCatIDs []uuid.UUID
	for _, sc := range req.Scores {
		subCatIDs = append(subCatIDs, sc.SubCategoryID)
	}

	if err := validateAssessmentOwnership(ctx, s.repo, organizerID, req.RegisID, req.JudgesID, subCatIDs, nil); err != nil {
		return err
	}

	rules, err := s.repo.GetSubCategoryGradeRules(ctx, subCatIDs)
	if err != nil || len(rules) == 0 {
		return errors.New("grade rules not found for some sub categories")
	}

	var scores []entity.Score
	for _, sc := range req.Scores {
		grade := getGrade(sc.ScoreValue, sc.SubCategoryID, rules)
		scores = append(scores, entity.Score{
			RegisID:       req.RegisID,
			JudgesID:      req.JudgesID,
			SubCategoryID: sc.SubCategoryID,
			ScoreValue:    sc.ScoreValue,
			Grade:         grade,
		})
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		if err := validateAssessmentOwnership(ctx, txRepo, organizerID, req.RegisID, req.JudgesID, subCatIDs, nil); err != nil {
			return err
		}
		return txRepo.BulkUpsertScores(ctx, scores)
	})
}

func (s *formPenilaianService) BulkInsertTeamViolations(ctx context.Context, organizerID uuid.UUID, req dto.BulkInsertViolationsRequest) error {
	if s.validate != nil {
		if err := s.validate.Struct(req); err != nil {
			return err
		}
	}

	if err := validateAssessmentOwnership(ctx, s.repo, organizerID, req.RegisID, req.JudgesID, nil, req.ViolationTypeIDs); err != nil {
		return err
	}

	var violations []entity.TeamViolation
	for _, vID := range req.ViolationTypeIDs {
		violations = append(violations, entity.TeamViolation{
			RegisID:         req.RegisID,
			JudgesID:        req.JudgesID,
			ViolationTypeID: vID,
		})
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		if err := validateAssessmentOwnership(ctx, txRepo, organizerID, req.RegisID, req.JudgesID, nil, req.ViolationTypeIDs); err != nil {
			return err
		}
		return txRepo.BulkInsertTeamViolations(ctx, violations)
	})
}

func (s *formPenilaianService) FinalizeAssessment(ctx context.Context, organizerID uuid.UUID, req dto.FinalizeAssessmentRequest) error {
	if s.validate != nil {
		if err := s.validate.Struct(req); err != nil {
			return err
		}
	}

	var subCatIDs []uuid.UUID
	for _, sc := range req.Scores {
		subCatIDs = append(subCatIDs, sc.SubCategoryID)
	}

	if err := validateAssessmentOwnership(ctx, s.repo, organizerID, req.RegisID, req.JudgesID, subCatIDs, req.ViolationTypeIDs); err != nil {
		return err
	}

	rules, err := s.repo.GetSubCategoryGradeRules(ctx, subCatIDs)
	if err != nil || len(rules) == 0 {
		return errors.New("grade rules not found for some sub categories")
	}

	var scores []entity.Score
	for _, sc := range req.Scores {
		grade := getGrade(sc.ScoreValue, sc.SubCategoryID, rules)
		scores = append(scores, entity.Score{
			RegisID:       req.RegisID,
			JudgesID:      req.JudgesID,
			SubCategoryID: sc.SubCategoryID,
			ScoreValue:    sc.ScoreValue,
			Grade:         grade,
		})
	}

	var violations []entity.TeamViolation
	for _, vID := range req.ViolationTypeIDs {
		violations = append(violations, entity.TeamViolation{
			RegisID:         req.RegisID,
			JudgesID:        req.JudgesID,
			ViolationTypeID: vID,
		})
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		if err := validateAssessmentOwnership(ctx, txRepo, organizerID, req.RegisID, req.JudgesID, subCatIDs, req.ViolationTypeIDs); err != nil {
			return err
		}

		// 1. Bulk Upsert Scores
		if err := txRepo.BulkUpsertScores(ctx, scores); err != nil {
			return err
		}

		// 2. Bulk Insert Violations
		// Note: We might want to clear existing violations for this judge/team if it's an update,
		// but since it's "Finalize", we usually insert new ones or we should have a clear step.
		// For simplicity, we just insert. If the user wants to clear, they should do it.
		// Usually, "Finalize" means it's the final state.
		if len(violations) > 0 {
			if err := txRepo.BulkInsertTeamViolations(ctx, violations); err != nil {
				return err
			}
		}

		// 3. Update Status to COMPLETED
		return txRepo.UpdateAssessmentStatus(ctx, req.RegisID, "COMPLETED")
	})
}

func validateAssessmentOwnership(
	ctx context.Context,
	repo contracts.FormPenilaianRepository,
	organizerID uuid.UUID,
	regisID uuid.UUID,
	judgeID uuid.UUID,
	subCategoryIDs []uuid.UUID,
	violationTypeIDs []uuid.UUID,
) error {
	authorized, err := repo.ValidateAssessmentOwnership(
		ctx,
		organizerID,
		regisID,
		judgeID,
		subCategoryIDs,
		violationTypeIDs,
	)
	if err != nil {
		return err
	}
	if !authorized {
		return domain.ErrForbidden
	}
	return nil
}

func getGrade(score float64, subCatID uuid.UUID, rules []entity.GradeRule) string {
	for _, rule := range rules {
		if rule.ScoreSubCategoryID == subCatID && score >= rule.MinScore && score <= rule.MaxScore {
			return rule.GradeName
		}
	}
	return "Undefined"
}
