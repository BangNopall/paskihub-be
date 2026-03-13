package repository

import (
	"context"

	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type rekapRepository struct {
	db *gorm.DB
}

func NewRekapRepository(db *gorm.DB) contracts.RekapRepository {
	return &rekapRepository{db: db}
}

func (r *rekapRepository) GetTeamAssessmentDetail(ctx context.Context, regisID uuid.UUID) (dto.TeamAssessmentDetailResponse, error) {
	var resp dto.TeamAssessmentDetailResponse

	var regis entity.Registration
	if err := r.db.WithContext(ctx).
		Preload("Team").
		Preload("Team.Institution").
		Where("id = ?", regisID).First(&regis).Error; err != nil {
		return resp, err
	}

	resp.RegisID = regis.Id
	resp.TeamName = regis.Team.Name
	resp.InstitutionName = regis.Team.Institution.Name

	type ScoreFlat struct {
		CategoryID      uuid.UUID
		CategoryName    string
		SubCategoryID   uuid.UUID
		SubCategoryName string
		JudgeName       string
		ScoreValue      float64
		Grade           string
	}
	var flatScores []ScoreFlat
	r.db.WithContext(ctx).Table("scores").
		Select("sc.id as category_id, sc.name as category_name, ssc.id as sub_category_id, ssc.name as sub_category_name, j.name as judge_name, scores.score_value, scores.grade").
		Joins("JOIN score_sub_categories ssc ON scores.sub_category_id = ssc.id").
		Joins("JOIN score_categories sc ON ssc.score_categories_id = sc.id").
		Joins("JOIN judges j ON scores.judges_id = j.id").
		Where("scores.regis_id = ?", regisID).
		Scan(&flatScores)

	catMap := make(map[uuid.UUID]*dto.CategoryScoreDetail)
	var totalScore float64

	for _, fs := range flatScores {
		if _, ok := catMap[fs.CategoryID]; !ok {
			catMap[fs.CategoryID] = &dto.CategoryScoreDetail{
				CategoryID:    fs.CategoryID,
				CategoryName:  fs.CategoryName,
				SubCategories: []dto.SubCategoryScoreDetail{},
			}
		}
		catMap[fs.CategoryID].SubCategories = append(catMap[fs.CategoryID].SubCategories, dto.SubCategoryScoreDetail{
			SubCategoryID:   fs.SubCategoryID,
			SubCategoryName: fs.SubCategoryName,
			JudgeName:       fs.JudgeName,
			ScoreValue:      fs.ScoreValue,
			Grade:           fs.Grade,
		})
		totalScore += fs.ScoreValue
	}

	for _, c := range catMap {
		resp.Categories = append(resp.Categories, *c)
	}

	var tvs []entity.TeamViolation
	r.db.WithContext(ctx).
		Preload("ViolationType").
		Preload("Judge").
		Where("regis_id = ?", regisID).Find(&tvs)

	var totalViolation float64
	for _, tv := range tvs {
		resp.Violations = append(resp.Violations, dto.TeamViolationDetail{
			ViolationID:   tv.ID,
			ViolationName: tv.ViolationType.Name,
			Point:         tv.ViolationType.Point,
			JudgeName:     tv.Judge.Name,
		})
		totalViolation += tv.ViolationType.Point
	}

	resp.TotalScore = totalScore
	resp.TotalViolation = totalViolation
	resp.FinalScore = totalScore - totalViolation

	return resp, nil
}

func (r *rekapRepository) GetScoreboardByEventLevel(ctx context.Context, eventLevelID uuid.UUID) ([]dto.ScoreboardItem, error) {
	query := `
		WITH TeamScore AS (
			SELECT regis_id, SUM(score_value) as total_score
			FROM scores
			GROUP BY regis_id
		),
		TeamViolationPoint AS (
			SELECT tv.regis_id, SUM(vt.point) as total_violation_points
			FROM team_violations tv
			JOIN violation_types vt ON tv.violation_type_id = vt.id
			GROUP BY tv.regis_id
		)
		SELECT 
			r.id as regis_id, 
			t.name as team_name, 
			i.name as insti_name,
			COALESCE(ts.total_score, 0) as total_score,
			COALESCE(tvp.total_violation_points, 0) as total_violation_points,
			(COALESCE(ts.total_score, 0) - COALESCE(tvp.total_violation_points, 0)) as final_score,
			RANK() OVER (ORDER BY (COALESCE(ts.total_score, 0) - COALESCE(tvp.total_violation_points, 0)) DESC) as rank
		FROM registrations r
		JOIN teams t ON r.team_id = t.id
		JOIN institutions i ON t.insti_id = i.id
		LEFT JOIN TeamScore ts ON ts.regis_id = r.id
		LEFT JOIN TeamViolationPoint tvp ON tvp.regis_id = r.id
		WHERE r.event_level_id = ?
		ORDER BY final_score DESC
	`
	var items []dto.ScoreboardItem
	err := r.db.WithContext(ctx).Raw(query, eventLevelID).Scan(&items).Error
	return items, err
}

func (r *rekapRepository) GetLeaderboardCustom(ctx context.Context, eventLevelID uuid.UUID, categoryIDs []uuid.UUID) ([]dto.ScoreboardItem, error) {
	query := `
		WITH TeamScore AS (
			SELECT s.regis_id, SUM(s.score_value) as total_score
			FROM scores s
			JOIN score_sub_categories ssc ON s.sub_category_id = ssc.id
			WHERE ssc.score_categories_id IN ?
			GROUP BY s.regis_id
		),
		TeamViolationPoint AS (
			SELECT tv.regis_id, SUM(vt.point) as total_violation_points
			FROM team_violations tv
			JOIN violation_types vt ON tv.violation_type_id = vt.id
			GROUP BY tv.regis_id
		)
		SELECT 
			r.id as regis_id, 
			t.name as team_name, 
			i.name as insti_name,
			COALESCE(ts.total_score, 0) as total_score,
			COALESCE(tvp.total_violation_points, 0) as total_violation_points,
			(COALESCE(ts.total_score, 0) - COALESCE(tvp.total_violation_points, 0)) as final_score,
			RANK() OVER (ORDER BY (COALESCE(ts.total_score, 0) - COALESCE(tvp.total_violation_points, 0)) DESC) as rank
		FROM registrations r
		JOIN teams t ON r.team_id = t.id
		JOIN institutions i ON t.insti_id = i.id
		LEFT JOIN TeamScore ts ON ts.regis_id = r.id
		LEFT JOIN TeamViolationPoint tvp ON tvp.regis_id = r.id
		WHERE r.event_level_id = ?
		ORDER BY final_score DESC
	`
	var items []dto.ScoreboardItem
	err := r.db.WithContext(ctx).Raw(query, categoryIDs, eventLevelID).Scan(&items).Error
	return items, err
}

func (r *rekapRepository) UpdateScorePublishedStatus(ctx context.Context, eventLevelID uuid.UUID, isPublished bool) error {
	return r.db.WithContext(ctx).
		Model(&entity.EventLevel{}).
		Where("id = ?", eventLevelID).
		Update("is_score_published", isPublished).Error
}

func (r *rekapRepository) GetEventIDByEventLevelID(ctx context.Context, eventLevelID uuid.UUID) (uuid.UUID, error) {
	var eventLevel entity.EventLevel
	err := r.db.WithContext(ctx).
		Where("id = ?", eventLevelID).
		First(&eventLevel).Error
	return eventLevel.EventId, err
}
