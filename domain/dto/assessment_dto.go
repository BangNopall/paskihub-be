package dto

import "github.com/google/uuid"

// Judge
type CreateJudgeReq struct {
	Name string `json:"name" validate:"required"`
}

type UpdateJudgeReq struct {
	Name string `json:"name" validate:"required"`
}

type JudgeRes struct {
	ID      uuid.UUID `json:"id"`
	EventID uuid.UUID `json:"event_id"`
	Name    string    `json:"name"`
}

// ViolationType
type CreateViolationTypeReq struct {
	Name         string    `json:"name" validate:"required"`
	Point        float64   `json:"point" validate:"required,min=0"`
	EventLevelID uuid.UUID `json:"event_level_id" validate:"required"`
}

type UpdateViolationTypeReq struct {
	Name  string  `json:"name" validate:"required"`
	Point float64 `json:"point" validate:"required,min=0"`
}

type ViolationTypeRes struct {
	ID           uuid.UUID `json:"id"`
	EventID      uuid.UUID `json:"event_id"`
	EventLevelID uuid.UUID `json:"event_level_id"`
	Name         string    `json:"name"`
	Point        float64   `json:"point"`
}

// ScoreCategory
type CreateScoreCategoryReq struct {
	Name         string    `json:"name" validate:"required"`
	EventLevelID uuid.UUID `json:"event_level_id" validate:"required"`
}

type UpdateScoreCategoryReq struct {
	Name string `json:"name" validate:"required"`
}

type ScoreCategoryRes struct {
	ID            uuid.UUID             `json:"id"`
	EventID       uuid.UUID             `json:"event_id"`
	EventLevelID  uuid.UUID             `json:"event_level_id"`
	Name          string                `json:"name"`
	SubCategories []ScoreSubCategoryRes `json:"sub_categories,omitempty"`
}

// ScoreSubCategory
type CreateScoreSubCategoryReq struct {
	ScoreCategoriesID uuid.UUID           `json:"score_categories_id" validate:"required"`
	Name              string              `json:"name" validate:"required"`
	MaxScore          float64             `json:"max_score" validate:"min=0"`
	Grades            map[string][]string `json:"grades"`
}

type UpdateScoreSubCategoryReq struct {
	Name     string              `json:"name" validate:"required"`
	MaxScore float64             `json:"max_score" validate:"min=0"`
	Grades   map[string][]string `json:"grades"`
}

type ScoreSubCategoryRes struct {
	ID                uuid.UUID           `json:"id"`
	ScoreCategoriesID uuid.UUID           `json:"score_categories_id"`
	Name              string              `json:"name"`
	MaxScore          float64             `json:"max_score"`
	Grades            map[string][]string `json:"grades,omitempty"`
	GradeNumbers      map[string][]int    `json:"grade_numbers,omitempty"`
}

// Unified Res
type UnifiedAssessmentRes struct {
	Violations []ViolationTypeRes `json:"violations"`
	Categories []ScoreCategoryRes `json:"categories"`
}

// GradeRule
type GradeRuleItemReq struct {
	GradeName string  `json:"grade_name" validate:"required,oneof=Kurang Cukup Baik 'Sangat Baik'"`
	MinScore  float64 `json:"min_score" validate:"min=0"`
	MaxScore  float64 `json:"max_score" validate:"gtefield=MinScore"`
}

type SetupGradeRulesReq struct {
	Rules []GradeRuleItemReq `json:"rules" validate:"required,dive,required"`
}

type GradeRuleRes struct {
	ID        uuid.UUID `json:"id"`
	EventID   uuid.UUID `json:"event_id"`
	GradeName string    `json:"grade_name"`
	MinScore  float64   `json:"min_score"`
	MaxScore  float64   `json:"max_score"`
}

// Score
type InputScoreReq struct {
	RegisID       uuid.UUID `json:"regis_id" validate:"required"`
	JudgesID      uuid.UUID `json:"judges_id" validate:"required"`
	SubCategoryID uuid.UUID `json:"sub_category_id" validate:"required"`
	ScoreValue    float64   `json:"score_value" validate:"required,min=0"`
}

type ScoreRes struct {
	ID            uuid.UUID `json:"id"`
	RegisID       uuid.UUID `json:"regis_id"`
	JudgesID      uuid.UUID `json:"judges_id"`
	SubCategoryID uuid.UUID `json:"sub_category_id"`
	ScoreValue    float64   `json:"score_value"`
	Grade         string    `json:"grade"`
}

// EventAward
type CreateAwardReq struct {
	EventLevelID     uuid.UUID   `json:"event_level_id" validate:"required"`
	Name             string      `json:"name" validate:"required"`
	LimitRank        int         `json:"limit_rank" validate:"required,min=1"`
	ScoreCategoryIDs []uuid.UUID `json:"score_category_ids" validate:"required,min=1"`
}

type UpdateAwardReq struct {
	Name             string      `json:"name" validate:"required"`
	LimitRank        int         `json:"limit_rank" validate:"required,min=1"`
	ScoreCategoryIDs []uuid.UUID `json:"score_category_ids" validate:"required,min=1"`
}

type AwardRes struct {
	ID              uuid.UUID          `json:"id"`
	EventID         uuid.UUID          `json:"event_id"`
	EventLevelID    uuid.UUID          `json:"event_level_id"`
	Name            string             `json:"name"`
	LimitRank       int                `json:"limit_rank"`
	ScoreCategories []ScoreCategoryRes `json:"score_categories,omitempty"`
}
