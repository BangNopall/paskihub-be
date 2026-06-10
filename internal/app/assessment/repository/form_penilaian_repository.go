package repository

import (
	"context"
	"errors"

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

func (r *formPenilaianRepository) ValidateAssessmentOwnership(
	ctx context.Context,
	organizerID uuid.UUID,
	regisID uuid.UUID,
	judgeID uuid.UUID,
	subCategoryIDs []uuid.UUID,
	violationTypeIDs []uuid.UUID,
) (bool, error) {
	var scope struct {
		EventID      uuid.UUID
		EventLevelID uuid.UUID
	}

	err := r.db.WithContext(ctx).
		Table("registrations AS r").
		Select("el.event_id, r.event_level_id").
		Joins("JOIN event_levels AS el ON el.id = r.event_level_id").
		Joins("JOIN events AS e ON e.id = el.event_id").
		Where("r.id = ? AND e.user_id = ?", regisID, organizerID).
		Take(&scope).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.Judge{}).
		Where("id = ? AND event_id = ?", judgeID, scope.EventID).
		Count(&count).Error; err != nil {
		return false, err
	}
	if count != 1 {
		return false, nil
	}

	uniqueSubCategoryIDs := uniqueUUIDs(subCategoryIDs)
	if len(uniqueSubCategoryIDs) > 0 {
		count = 0
		if err := r.db.WithContext(ctx).
			Table("score_sub_categories AS ssc").
			Joins("JOIN score_categories AS sc ON sc.id = ssc.score_categories_id").
			Where("ssc.id IN ? AND sc.event_id = ? AND sc.event_level_id = ?", uniqueSubCategoryIDs, scope.EventID, scope.EventLevelID).
			Distinct("ssc.id").
			Count(&count).Error; err != nil {
			return false, err
		}
		if count != int64(len(uniqueSubCategoryIDs)) {
			return false, nil
		}
	}

	uniqueViolationTypeIDs := uniqueUUIDs(violationTypeIDs)
	if len(uniqueViolationTypeIDs) > 0 {
		count = 0
		if err := r.db.WithContext(ctx).
			Model(&entity.ViolationType{}).
			Where("id IN ? AND event_id = ? AND event_level_id = ?", uniqueViolationTypeIDs, scope.EventID, scope.EventLevelID).
			Distinct("id").
			Count(&count).Error; err != nil {
			return false, err
		}
		if count != int64(len(uniqueViolationTypeIDs)) {
			return false, nil
		}
	}

	return true, nil
}

func uniqueUUIDs(ids []uuid.UUID) []uuid.UUID {
	seen := make(map[uuid.UUID]struct{}, len(ids))
	unique := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		unique = append(unique, id)
	}
	return unique
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

func (r *formPenilaianRepository) GetSubCategoryGradeRules(ctx context.Context, subCategoryIDs []uuid.UUID) ([]entity.GradeRule, error) {
	var rules []entity.GradeRule
	err := r.db.WithContext(ctx).
		Where("score_sub_category_id IN ?", subCategoryIDs).
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

func (r *formPenilaianRepository) UpdateAssessmentStatus(ctx context.Context, regisID uuid.UUID, status string) error {
	return r.db.WithContext(ctx).
		Model(&entity.Registration{}).
		Where("id = ?", regisID).
		Update("assessment_status", status).Error
}
