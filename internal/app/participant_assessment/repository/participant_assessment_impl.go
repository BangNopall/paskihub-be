package repository

import (
	"context"

	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type participantAssessmentRepository struct {
	db *gorm.DB
}

func NewParticipantAssessmentRepository(db *gorm.DB) contracts.ParticipantAssessmentRepository {
	return &participantAssessmentRepository{
		db: db,
	}
}

func (r *participantAssessmentRepository) GetScoresByRegisID(ctx context.Context, regisID uuid.UUID) ([]entity.Score, error) {
	var scores []entity.Score
	err := r.db.WithContext(ctx).
		Preload("Judge").
		Preload("ScoreSubCategory").
		Where("regis_id = ?", regisID).
		Find(&scores).Error
	return scores, err
}

func (r *participantAssessmentRepository) GetViolationsByRegisID(ctx context.Context, regisID uuid.UUID) ([]entity.TeamViolation, error) {
	var violations []entity.TeamViolation
	err := r.db.WithContext(ctx).
		Preload("Judge").
		Preload("ViolationType").
		Where("regis_id = ?", regisID).
		Find(&violations).Error
	return violations, err
}

func (r *participantAssessmentRepository) GetRegistrationOwnership(ctx context.Context, regisID uuid.UUID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Registration{}).
		Joins("JOIN teams ON teams.id = registrations.team_id").
		Joins("JOIN institutions ON institutions.id = teams.insti_id").
		Joins("JOIN users ON users.id = institutions.user_id").
		Where("registrations.id = ? AND users.id = ?", regisID, userID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *participantAssessmentRepository) GetRegistrationByID(ctx context.Context, regisID uuid.UUID) (*entity.Registration, error) {
	var regis entity.Registration
	err := r.db.WithContext(ctx).Preload("EventLevel.Event").First(&regis, "id = ?", regisID).Error
	if err != nil {
		return nil, err
	}
	return &regis, nil
}

func (r *participantAssessmentRepository) GetScoreCategoriesByEventID(ctx context.Context, eventID uuid.UUID) ([]entity.ScoreCategory, error) {
	var categories []entity.ScoreCategory
	err := r.db.WithContext(ctx).Preload("SubCategories").Where("event_id = ?", eventID).Find(&categories).Error
	return categories, err
}
