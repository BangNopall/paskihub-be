package repository

import (
	"context"

	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type assessmentRepository struct {
	db *gorm.DB
}

func NewAssessmentRepository(db *gorm.DB) contracts.IAssessmentRepository {
	return &assessmentRepository{db: db}
}

func (r *assessmentRepository) CheckEventOwnership(ctx context.Context, eventId, userId uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Event{}).
		Where("id = ? AND user_id = ?", eventId, userId).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *assessmentRepository) CreateJudge(ctx context.Context, judge *entity.Judge) error {
	return r.db.WithContext(ctx).Create(judge).Error
}

func (r *assessmentRepository) GetJudgesByEvent(ctx context.Context, eventId uuid.UUID) ([]entity.Judge, error) {
	var judges []entity.Judge
	err := r.db.WithContext(ctx).Where("event_id = ?", eventId).Find(&judges).Error
	return judges, err
}

func (r *assessmentRepository) FindJudgeById(ctx context.Context, id uuid.UUID) (*entity.Judge, error) {
	var judge entity.Judge
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&judge).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &judge, nil
}

func (r *assessmentRepository) UpdateJudge(ctx context.Context, judge *entity.Judge) error {
	return r.db.WithContext(ctx).Save(judge).Error
}

func (r *assessmentRepository) DeleteJudge(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Judge{}).Error
}

func (r *assessmentRepository) CreateViolationType(ctx context.Context, vt *entity.ViolationType) error {
	return r.db.WithContext(ctx).Create(vt).Error
}

func (r *assessmentRepository) GetViolationTypesByEvent(ctx context.Context, eventId uuid.UUID) ([]entity.ViolationType, error) {
	var vts []entity.ViolationType
	err := r.db.WithContext(ctx).Where("event_id = ?", eventId).Find(&vts).Error
	return vts, err
}

func (r *assessmentRepository) FindViolationTypeById(ctx context.Context, id uuid.UUID) (*entity.ViolationType, error) {
	var vt entity.ViolationType
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&vt).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &vt, nil
}

func (r *assessmentRepository) UpdateViolationType(ctx context.Context, vt *entity.ViolationType) error {
	return r.db.WithContext(ctx).Save(vt).Error
}

func (r *assessmentRepository) DeleteViolationType(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.ViolationType{}).Error
}

func (r *assessmentRepository) CreateScoreCategory(ctx context.Context, sc *entity.ScoreCategory) error {
	return r.db.WithContext(ctx).Create(sc).Error
}

func (r *assessmentRepository) GetScoreCategoriesByEvent(ctx context.Context, eventId uuid.UUID) ([]entity.ScoreCategory, error) {
	var scs []entity.ScoreCategory
	err := r.db.WithContext(ctx).Preload("SubCategories").Where("event_id = ?", eventId).Find(&scs).Error
	return scs, err
}

func (r *assessmentRepository) FindScoreCategoryById(ctx context.Context, id uuid.UUID) (*entity.ScoreCategory, error) {
	var sc entity.ScoreCategory
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&sc).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &sc, nil
}

func (r *assessmentRepository) UpdateScoreCategory(ctx context.Context, sc *entity.ScoreCategory) error {
	return r.db.WithContext(ctx).Save(sc).Error
}

func (r *assessmentRepository) DeleteScoreCategory(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.ScoreCategory{}).Error
}

func (r *assessmentRepository) CreateScoreSubCategory(ctx context.Context, ssc *entity.ScoreSubCategory) error {
	return r.db.WithContext(ctx).Create(ssc).Error
}

func (r *assessmentRepository) FindScoreSubCategoryById(ctx context.Context, id uuid.UUID) (*entity.ScoreSubCategory, error) {
	var ssc entity.ScoreSubCategory
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&ssc).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &ssc, nil
}

func (r *assessmentRepository) UpdateScoreSubCategory(ctx context.Context, ssc *entity.ScoreSubCategory) error {
	return r.db.WithContext(ctx).Save(ssc).Error
}

func (r *assessmentRepository) DeleteScoreSubCategory(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.ScoreSubCategory{}).Error
}

func (r *assessmentRepository) GetGradeRulesByEvent(ctx context.Context, eventId uuid.UUID) ([]entity.GradeRule, error) {
	var rules []entity.GradeRule
	err := r.db.WithContext(ctx).Where("event_id = ?", eventId).Order("min_score asc").Find(&rules).Error
	return rules, err
}

func (r *assessmentRepository) ReplaceGradeRules(ctx context.Context, eventId uuid.UUID, rules []entity.GradeRule) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("event_id = ?", eventId).Delete(&entity.GradeRule{}).Error; err != nil {
			return err
		}
		if len(rules) > 0 {
			if err := tx.Create(&rules).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *assessmentRepository) CreateScore(ctx context.Context, score *entity.Score) error {
	return r.db.WithContext(ctx).Create(score).Error
}
