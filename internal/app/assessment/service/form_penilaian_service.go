package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/go-playground/validator/v10"
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

func (s *formPenilaianService) BulkInsertScores(ctx context.Context, req dto.BulkInsertScoresRequest) error {
	if s.validate != nil {
		if err := s.validate.Struct(req); err != nil {
			return err
		}
	}

	eventID, err := s.repo.GetEventIDByRegisID(ctx, req.RegisID)
	if err != nil {
		return fmt.Errorf("failed to get event id from regis id: %w", err)
	}

	rules, err := s.repo.GetEventGradeRules(ctx, eventID)
	if err != nil || len(rules) == 0 {
		return errors.New("grade rules not found for this event")
	}

	var scores []entity.Score
	for _, sc := range req.Scores {
		grade := getGrade(sc.ScoreValue, rules)
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
		return txRepo.BulkUpsertScores(ctx, scores)
	})
}

func (s *formPenilaianService) BulkInsertTeamViolations(ctx context.Context, req dto.BulkInsertViolationsRequest) error {
	if s.validate != nil {
		if err := s.validate.Struct(req); err != nil {
			return err
		}
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
		return txRepo.BulkInsertTeamViolations(ctx, violations)
	})
}

func getGrade(score float64, rules []entity.GradeRule) string {
	for _, rule := range rules {
		if score >= rule.MinScore && score <= rule.MaxScore {
			return rule.GradeName
		}
	}
	return "Undefined"
}
