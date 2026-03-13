package dto

import "github.com/google/uuid"

type ScoreInput struct {
	SubCategoryID uuid.UUID `json:"sub_category_id" validate:"required"`
	ScoreValue    float64   `json:"score_value" validate:"required,min=0"`
}

type BulkInsertScoresRequest struct {
	RegisID  uuid.UUID    `json:"regis_id" validate:"required"`
	JudgesID uuid.UUID    `json:"judges_id" validate:"required"`
	Scores   []ScoreInput `json:"scores" validate:"required,dive"`
}

type BulkInsertViolationsRequest struct {
	RegisID          uuid.UUID   `json:"regis_id" validate:"required"`
	JudgesID         uuid.UUID   `json:"judges_id" validate:"required"`
	ViolationTypeIDs []uuid.UUID `json:"violation_type_ids" validate:"required,min=1"`
}
