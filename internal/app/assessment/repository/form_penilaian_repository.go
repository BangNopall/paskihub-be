package repository

import (
	"context"

	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type formPenilaianRepository struct {
	db *gorm.DB
}

func NewFormPenilaianRepository(db *gorm.DB) contracts.FormPenilaianRepository {
	return &formPenilaianRepository{db: db}
}

func (r *formPenilaianRepository) WithTx(tx *gorm.DB) contracts.FormPenilaianRepository {
	return &formPenilaianRepository{db: tx}
}

func (r *formPenilaianRepository) GetEventIDByRegisID(ctx context.Context, regisID uuid.UUID) (uuid.UUID, error) {
	var regis entity.Registration
	err := r.db.WithContext(ctx).
		Preload("EventLevel").
		Where("id = ?", regisID).
		First(&regis).Error
	if err != nil {
		return uuid.Nil, err
	}
	return regis.EventLevel.EventId, nil
}

func (r *formPenilaianRepository) GetEventGradeRules(ctx context.Context, eventID uuid.UUID) ([]entity.GradeRule, error) {
	var rules []entity.GradeRule
	err := r.db.WithContext(ctx).
		Where("event_id = ?", eventID).
		Find(&rules).Error
	return rules, err
}

func (r *formPenilaianRepository) BulkUpsertScores(ctx context.Context, scores []entity.Score) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "regis_id"}, {Name: "judges_id"}, {Name: "sub_category_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"score_value", "grade", "updated_at"}),
		}).
		Create(&scores).Error
}

func (r *formPenilaianRepository) BulkInsertTeamViolations(ctx context.Context, violations []entity.TeamViolation) error {
	return r.db.WithContext(ctx).Create(&violations).Error
}
