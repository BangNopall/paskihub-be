package dto

import "github.com/google/uuid"

// Scoreboard Response DTOs
type ScoreboardItem struct {
	RegisID              uuid.UUID `json:"regis_id"`
	TeamName             string    `json:"team_name"`
	InstitutionName      string    `json:"insti_name"`
	TotalScore           float64   `json:"total_score"`
	TotalViolationPoints float64   `json:"total_violation_points"`
	FinalScore           float64   `json:"final_score"`
	Rank                 int       `json:"rank"`
}

type ScoreboardResponse struct {
	EventLevelID uuid.UUID        `json:"event_level_id"`
	Items        []ScoreboardItem `json:"items"`
}

type PublishScoreboardRequest struct {
	IsScorePublished bool `json:"is_score_published"`
}

type CustomLeaderboardRequest struct {
	ScoreCategoryIDs []uuid.UUID `json:"score_category_ids" validate:"required,min=1"`
}

// Detail Assessment DTOs
type SubCategoryScoreDetail struct {
	SubCategoryID   uuid.UUID `json:"sub_category_id"`
	SubCategoryName string    `json:"sub_category_name"`
	JudgeName       string    `json:"judge_name"`
	ScoreValue      float64   `json:"score_value"`
	Grade           string    `json:"grade"`
}

type CategoryScoreDetail struct {
	CategoryID    uuid.UUID                `json:"category_id"`
	CategoryName  string                   `json:"category_name"`
	SubCategories []SubCategoryScoreDetail `json:"sub_categories"`
}

type TeamViolationDetail struct {
	ViolationID   uuid.UUID `json:"violation_id"`
	ViolationName string    `json:"violation_name"`
	Point         float64   `json:"point"`
	JudgeName     string    `json:"judge_name"`
}

type TeamAssessmentDetailResponse struct {
	RegisID         uuid.UUID             `json:"regis_id"`
	TeamName        string                `json:"team_name"`
	InstitutionName string                `json:"insti_name"`
	Categories      []CategoryScoreDetail `json:"categories"`
	Violations      []TeamViolationDetail `json:"violations"`
	TotalScore      float64               `json:"total_score"`
	TotalViolation  float64               `json:"total_violation"`
	FinalScore      float64               `json:"final_score"`
}
