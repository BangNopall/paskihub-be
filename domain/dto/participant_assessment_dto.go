package dto

type AssessmentRecapResponse struct {
	TotalScore           float64                     `json:"total_score"`
	TotalViolationPoints float64                     `json:"total_violation_points"`
	FinalScore           float64                     `json:"final_score"`
	Violations           []AssessmentViolationDetail `json:"violations"`
	Categories           []AssessmentCategoryDetail  `json:"categories"`
}

type AssessmentViolationDetail struct {
	Name      string  `json:"name"`
	Point     float64 `json:"point"`
	JudgeName string  `json:"judge_name"`
}

type AssessmentCategoryDetail struct {
	CategoryName  string                        `json:"category_name"`
	SubCategories []AssessmentSubCategoryDetail `json:"sub_categories"`
}

type AssessmentSubCategoryDetail struct {
	Name   string             `json:"name"`
	Scores []JudgeScoreDetail `json:"scores"`
}

type JudgeScoreDetail struct {
	JudgeName string  `json:"judge_name"`
	Value     float64 `json:"value"`
}
