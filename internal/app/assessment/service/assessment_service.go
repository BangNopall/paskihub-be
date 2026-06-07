package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/google/uuid"
)

type assessmentService struct {
	repo contracts.IAssessmentRepository
}

func NewAssessmentService(repo contracts.IAssessmentRepository) contracts.IAssessmentService {
	return &assessmentService{repo: repo}
}

func (s *assessmentService) ensureOwnership(ctx context.Context, eventId, userId uuid.UUID) error {
	isOwner, err := s.repo.CheckEventOwnership(ctx, eventId, userId)
	if err != nil {
		return err
	}
	if !isOwner {
		return domain.ErrForbidden // Assuming custom error
	}
	return nil
}

// Judge
func (s *assessmentService) CreateJudge(ctx context.Context, eventId, userId uuid.UUID, req dto.CreateJudgeReq) (*dto.JudgeRes, error) {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}
	judge := &entity.Judge{
		EventID: eventId,
		Name:    req.Name,
	}
	if err := s.repo.CreateJudge(ctx, judge); err != nil {
		return nil, domain.ErrInternalServer
	}
	return &dto.JudgeRes{
		ID:      judge.ID,
		EventID: judge.EventID,
		Name:    judge.Name,
	}, nil
}

func (s *assessmentService) GetJudges(ctx context.Context, eventId, userId uuid.UUID) ([]dto.JudgeRes, error) {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}
	judges, err := s.repo.GetJudgesByEvent(ctx, eventId)
	if err != nil {
		return nil, domain.ErrInternalServer
	}
	var res []dto.JudgeRes
	for _, j := range judges {
		res = append(res, dto.JudgeRes{
			ID:      j.ID,
			EventID: j.EventID,
			Name:    j.Name,
		})
	}
	return res, nil
}

func (s *assessmentService) UpdateJudge(ctx context.Context, eventId, userId, id uuid.UUID, req dto.UpdateJudgeReq) (*dto.JudgeRes, error) {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}
	judge, err := s.repo.FindJudgeById(ctx, id)
	if err != nil {
		return nil, domain.ErrInternalServer
	}
	if judge == nil || judge.EventID != eventId {
		return nil, domain.ErrNotFound
	}
	judge.Name = req.Name
	if err := s.repo.UpdateJudge(ctx, judge); err != nil {
		return nil, domain.ErrInternalServer
	}
	return &dto.JudgeRes{
		ID:      judge.ID,
		EventID: judge.EventID,
		Name:    judge.Name,
	}, nil
}

func (s *assessmentService) DeleteJudge(ctx context.Context, eventId, userId, id uuid.UUID) error {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return err
	}
	judge, err := s.repo.FindJudgeById(ctx, id)
	if err != nil {
		return domain.ErrInternalServer
	}
	if judge == nil || judge.EventID != eventId {
		return domain.ErrNotFound
	}
	if err := s.repo.DeleteJudge(ctx, id); err != nil {
		return domain.ErrInternalServer
	}
	return nil
}

// ViolationType
func (s *assessmentService) CreateViolationType(ctx context.Context, eventId, userId uuid.UUID, req dto.CreateViolationTypeReq) (*dto.ViolationTypeRes, error) {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}
	if err := s.ensureEventLevelRelation(ctx, eventId, req.EventLevelID); err != nil {
		return nil, err
	}
	vt := &entity.ViolationType{
		EventID:      eventId,
		EventLevelID: req.EventLevelID,
		Name:         req.Name,
		Point:        req.Point,
	}
	if err := s.repo.CreateViolationType(ctx, vt); err != nil {
		return nil, domain.ErrInternalServer
	}
	return &dto.ViolationTypeRes{
		ID:           vt.ID,
		EventID:      vt.EventID,
		EventLevelID: vt.EventLevelID,
		Name:         vt.Name,
		Point:        vt.Point,
	}, nil
}

func (s *assessmentService) GetViolationTypes(ctx context.Context, eventId, levelId, userId uuid.UUID) ([]dto.ViolationTypeRes, error) {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}
	vts, err := s.repo.GetViolationTypesByLevel(ctx, eventId, levelId)
	if err != nil {
		return nil, domain.ErrInternalServer
	}
	var res []dto.ViolationTypeRes
	for _, vt := range vts {
		res = append(res, dto.ViolationTypeRes{
			ID:           vt.ID,
			EventID:      vt.EventID,
			EventLevelID: vt.EventLevelID,
			Name:         vt.Name,
			Point:        vt.Point,
		})
	}
	return res, nil
}

func (s *assessmentService) UpdateViolationType(ctx context.Context, eventId, userId, id uuid.UUID, req dto.UpdateViolationTypeReq) (*dto.ViolationTypeRes, error) {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}
	vt, err := s.repo.FindViolationTypeById(ctx, id)
	if err != nil {
		return nil, domain.ErrInternalServer
	}
	if vt == nil || vt.EventID != eventId {
		return nil, domain.ErrNotFound
	}
	vt.Name = req.Name
	vt.Point = req.Point
	if err := s.repo.UpdateViolationType(ctx, vt); err != nil {
		return nil, domain.ErrInternalServer
	}
	return &dto.ViolationTypeRes{
		ID:      vt.ID,
		EventID: vt.EventID,
		Name:    vt.Name,
		Point:   vt.Point,
	}, nil
}

func (s *assessmentService) DeleteViolationType(ctx context.Context, eventId, userId, id uuid.UUID) error {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return err
	}
	vt, err := s.repo.FindViolationTypeById(ctx, id)
	if err != nil {
		return domain.ErrInternalServer
	}
	if vt == nil || vt.EventID != eventId {
		return domain.ErrNotFound
	}
	if err := s.repo.DeleteViolationType(ctx, id); err != nil {
		return domain.ErrInternalServer
	}
	return nil
}

// ScoreCategory
func (s *assessmentService) CreateScoreCategory(ctx context.Context, eventId, userId uuid.UUID, req dto.CreateScoreCategoryReq) (*dto.ScoreCategoryRes, error) {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}
	if err := s.ensureEventLevelRelation(ctx, eventId, req.EventLevelID); err != nil {
		return nil, err
	}
	sc := &entity.ScoreCategory{
		EventID:      eventId,
		EventLevelID: req.EventLevelID,
		Name:         req.Name,
	}
	if err := s.repo.CreateScoreCategory(ctx, sc); err != nil {
		return nil, domain.ErrInternalServer
	}
	return &dto.ScoreCategoryRes{
		ID:           sc.ID,
		EventID:      sc.EventID,
		EventLevelID: sc.EventLevelID,
		Name:         sc.Name,
	}, nil
}

func (s *assessmentService) GetScoreCategories(ctx context.Context, eventId, levelId, userId uuid.UUID) ([]dto.ScoreCategoryRes, error) {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}
	scs, err := s.repo.GetScoreCategoriesByLevel(ctx, eventId, levelId)
	if err != nil {
		return nil, domain.ErrInternalServer
	}
	var res []dto.ScoreCategoryRes
	for _, sc := range scs {
		var subcats []dto.ScoreSubCategoryRes
		for _, sub := range sc.SubCategories {
			grades := make(map[string][]string)
			gradeNumbers := make(map[string][]int)
			for _, gr := range sub.GradeRules {
				rangeStr := fmt.Sprintf("%v-%v", gr.MinScore, gr.MaxScore)
				if gr.MinScore == gr.MaxScore {
					rangeStr = fmt.Sprintf("%v", gr.MinScore)
				}
				grades[gr.GradeName] = append(grades[gr.GradeName], rangeStr)

				// Populate gradeNumbers with sequence of integers
				for i := int(gr.MinScore); i <= int(gr.MaxScore); i++ {
					gradeNumbers[gr.GradeName] = append(gradeNumbers[gr.GradeName], i)
				}
			}
			subcats = append(subcats, dto.ScoreSubCategoryRes{
				ID:                sub.ID,
				ScoreCategoriesID: sub.ScoreCategoriesID,
				Name:              sub.Name,
				MaxScore:          sub.MaxScore,
				Grades:            grades,
				GradeNumbers:      gradeNumbers,
			})
		}
		res = append(res, dto.ScoreCategoryRes{
			ID:            sc.ID,
			EventID:       sc.EventID,
			EventLevelID:  sc.EventLevelID,
			Name:          sc.Name,
			SubCategories: subcats,
		})
	}
	return res, nil
}

func (s *assessmentService) UpdateScoreCategory(ctx context.Context, eventId, userId, id uuid.UUID, req dto.UpdateScoreCategoryReq) (*dto.ScoreCategoryRes, error) {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}
	sc, err := s.repo.FindScoreCategoryById(ctx, id)
	if err != nil {
		return nil, domain.ErrInternalServer
	}
	if sc == nil || sc.EventID != eventId {
		return nil, domain.ErrNotFound
	}
	sc.Name = req.Name
	if err := s.repo.UpdateScoreCategory(ctx, sc); err != nil {
		return nil, domain.ErrInternalServer
	}
	return &dto.ScoreCategoryRes{
		ID:      sc.ID,
		EventID: sc.EventID,
		Name:    sc.Name,
	}, nil
}

func (s *assessmentService) DeleteScoreCategory(ctx context.Context, eventId, userId, id uuid.UUID) error {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return err
	}
	sc, err := s.repo.FindScoreCategoryById(ctx, id)
	if err != nil {
		return domain.ErrInternalServer
	}
	if sc == nil || sc.EventID != eventId {
		return domain.ErrNotFound
	}
	if err := s.repo.DeleteScoreCategory(ctx, id); err != nil {
		return domain.ErrInternalServer
	}
	return nil
}

// ScoreSubCategory
func (s *assessmentService) CreateScoreSubCategory(ctx context.Context, eventId, userId uuid.UUID, req dto.CreateScoreSubCategoryReq) (*dto.ScoreSubCategoryRes, error) {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}
	sc, err := s.repo.FindScoreCategoryById(ctx, req.ScoreCategoriesID)
	if err != nil || sc == nil || sc.EventID != eventId {
		return nil, domain.ErrNotFound
	}

	ssc := &entity.ScoreSubCategory{
		ScoreCategoriesID: req.ScoreCategoriesID,
		Name:              req.Name,
		MaxScore:          req.MaxScore,
	}
	if err := s.repo.CreateScoreSubCategory(ctx, ssc); err != nil {
		return nil, domain.ErrInternalServer
	}

	if len(req.Grades) > 0 {
		rules := s.parseGrades(eventId, sc.EventLevelID, ssc.ID, req.Grades)
		if err := s.repo.ReplaceSubCategoryGradeRules(ctx, ssc.ID, rules); err != nil {
			return nil, domain.ErrInternalServer
		}
	}

	return &dto.ScoreSubCategoryRes{
		ID:                ssc.ID,
		ScoreCategoriesID: ssc.ScoreCategoriesID,
		Name:              ssc.Name,
		MaxScore:          ssc.MaxScore,
		Grades:            req.Grades,
	}, nil
}

func (s *assessmentService) UpdateScoreSubCategory(ctx context.Context, eventId, userId, id uuid.UUID, req dto.UpdateScoreSubCategoryReq) (*dto.ScoreSubCategoryRes, error) {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}
	ssc, err := s.repo.FindScoreSubCategoryById(ctx, id)
	if err != nil || ssc == nil {
		return nil, domain.ErrNotFound
	}
	sc, err := s.repo.FindScoreCategoryById(ctx, ssc.ScoreCategoriesID)
	if err != nil || sc == nil || sc.EventID != eventId {
		return nil, domain.ErrNotFound
	}

	ssc.Name = req.Name
	ssc.MaxScore = req.MaxScore
	if err := s.repo.UpdateScoreSubCategory(ctx, ssc); err != nil {
		return nil, domain.ErrInternalServer
	}

	if len(req.Grades) > 0 {
		rules := s.parseGrades(eventId, sc.EventLevelID, ssc.ID, req.Grades)
		if err := s.repo.ReplaceSubCategoryGradeRules(ctx, ssc.ID, rules); err != nil {
			return nil, domain.ErrInternalServer
		}
	}

	return &dto.ScoreSubCategoryRes{
		ID:                ssc.ID,
		ScoreCategoriesID: ssc.ScoreCategoriesID,
		Name:              ssc.Name,
		MaxScore:          ssc.MaxScore,
		Grades:            req.Grades,
	}, nil
}

func (s *assessmentService) parseGrades(eventId, levelId, subCategoryId uuid.UUID, grades map[string][]string) []entity.GradeRule {
	var rules []entity.GradeRule
	for gradeName, ranges := range grades {
		for _, r := range ranges {
			parts := strings.Split(r, "-")
			var min, max float64
			if len(parts) == 2 {
				min, _ = strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
				max, _ = strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			} else {
				min, _ = strconv.ParseFloat(strings.TrimSpace(r), 64)
				max = min
			}
			rules = append(rules, entity.GradeRule{
				EventID:            eventId,
				EventLevelID:       levelId,
				ScoreSubCategoryID: subCategoryId,
				GradeName:          gradeName,
				MinScore:           min,
				MaxScore:           max,
			})
		}
	}
	return rules
}

func (s *assessmentService) DeleteScoreSubCategory(ctx context.Context, eventId, userId, id uuid.UUID) error {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return err
	}
	ssc, err := s.repo.FindScoreSubCategoryById(ctx, id)
	if err != nil || ssc == nil {
		return domain.ErrNotFound
	}
	sc, err := s.repo.FindScoreCategoryById(ctx, ssc.ScoreCategoriesID)
	if err != nil || sc == nil || sc.EventID != eventId {
		return domain.ErrNotFound
	}

	if err := s.repo.DeleteScoreSubCategory(ctx, id); err != nil {
		return domain.ErrInternalServer
	}
	return nil
}

// Unified Endpoint
func (s *assessmentService) GetUnifiedAssessment(ctx context.Context, eventId, levelId, userId uuid.UUID) (*dto.UnifiedAssessmentRes, error) {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}

	violations, err := s.GetViolationTypes(ctx, eventId, levelId, userId)
	if err != nil {
		return nil, err
	}

	categories, err := s.GetScoreCategories(ctx, eventId, levelId, userId)
	if err != nil {
		return nil, err
	}

	return &dto.UnifiedAssessmentRes{
		Violations: violations,
		Categories: categories,
	}, nil
}

// Score
func (s *assessmentService) InputScore(ctx context.Context, eventId, userId uuid.UUID, req dto.InputScoreReq) (*dto.ScoreRes, error) {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}

	valid, err := s.repo.ValidateScoreInputRelations(ctx, eventId, req.RegisID, req.JudgesID, req.SubCategoryID)
	if err != nil {
		return nil, domain.ErrInternalServer
	}
	if !valid {
		return nil, domain.ErrBadRequest
	}

	rules, err := s.repo.GetGradeRulesBySubCategory(ctx, req.SubCategoryID)
	if err != nil {
		return nil, domain.ErrInternalServer
	}

	if len(rules) == 0 {
		return nil, domain.ErrBadRequest // Sub-category must have grade rules
	}

	var evaluatedGrade string
	for _, rule := range rules {
		if req.ScoreValue >= rule.MinScore && req.ScoreValue <= rule.MaxScore {
			evaluatedGrade = rule.GradeName
			break
		}
	}

	if evaluatedGrade == "" {
		return nil, domain.ErrBadRequest // Value doesn't hit any bracket
	}

	score := &entity.Score{
		RegisID:       req.RegisID,
		JudgesID:      req.JudgesID,
		SubCategoryID: req.SubCategoryID,
		ScoreValue:    req.ScoreValue,
		Grade:         evaluatedGrade,
	}

	if err := s.repo.CreateScore(ctx, score); err != nil {
		return nil, domain.ErrInternalServer
	}

	return &dto.ScoreRes{
		ID:            score.ID,
		RegisID:       score.RegisID,
		JudgesID:      score.JudgesID,
		SubCategoryID: score.SubCategoryID,
		ScoreValue:    score.ScoreValue,
		Grade:         score.Grade,
	}, nil
}

// Award
func (s *assessmentService) CreateAward(ctx context.Context, eventId, userId uuid.UUID, req dto.CreateAwardReq) (*dto.AwardRes, error) {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}
	if err := s.ensureAwardRelations(ctx, eventId, req.EventLevelIDs, req.ScoreCategoryIDs); err != nil {
		return nil, err
	}

	var categories []entity.ScoreCategory
	for _, catId := range req.ScoreCategoryIDs {
		categories = append(categories, entity.ScoreCategory{ID: catId})
	}

	var levels []entity.EventLevel
	for _, levelId := range req.EventLevelIDs {
		levels = append(levels, entity.EventLevel{Id: levelId})
	}

	award := &entity.EventAward{
		EventID:         eventId,
		Name:            req.Name,
		LimitRank:       req.LimitRank,
		EventLevels:     levels,
		ScoreCategories: categories,
	}

	if err := s.repo.CreateAward(ctx, award); err != nil {
		return nil, domain.ErrInternalServer
	}

	// Re-fetch to get full details for response
	newAward, _ := s.repo.FindAwardById(ctx, award.ID)
	return s.mapToAwardRes(newAward), nil
}

func (s *assessmentService) GetAwards(ctx context.Context, eventId, userId uuid.UUID) ([]dto.AwardRes, error) {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}
	awards, err := s.repo.GetAwardsByEvent(ctx, eventId)
	if err != nil {
		return nil, domain.ErrInternalServer
	}
	var res []dto.AwardRes
	for _, a := range awards {
		res = append(res, *s.mapToAwardRes(&a))
	}
	return res, nil
}

func (s *assessmentService) UpdateAward(ctx context.Context, eventId, userId, id uuid.UUID, req dto.UpdateAwardReq) (*dto.AwardRes, error) {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}
	award, err := s.repo.FindAwardById(ctx, id)
	if err != nil {
		return nil, domain.ErrInternalServer
	}
	if award == nil || award.EventID != eventId {
		return nil, domain.ErrNotFound
	}
	if err := s.ensureAwardRelations(ctx, eventId, req.EventLevelIDs, req.ScoreCategoryIDs); err != nil {
		return nil, err
	}

	award.Name = req.Name
	award.LimitRank = req.LimitRank

	if err := s.repo.UpdateAward(ctx, award, req.EventLevelIDs, req.ScoreCategoryIDs); err != nil {
		return nil, domain.ErrInternalServer
	}

	// Re-fetch to get updated relations
	updatedAward, err := s.repo.FindAwardById(ctx, id)
	if err != nil {
		return nil, domain.ErrInternalServer
	}
	return s.mapToAwardRes(updatedAward), nil
}

func (s *assessmentService) ensureEventLevelRelation(ctx context.Context, eventId, eventLevelId uuid.UUID) error {
	valid, err := s.repo.EventLevelBelongsToEvent(ctx, eventId, eventLevelId)
	if err != nil {
		return domain.ErrInternalServer
	}
	if !valid {
		return domain.ErrBadRequest
	}
	return nil
}

func (s *assessmentService) ensureAwardRelations(
	ctx context.Context,
	eventId uuid.UUID,
	eventLevelIds []uuid.UUID,
	scoreCategoryIds []uuid.UUID,
) error {
	valid, err := s.repo.ValidateAwardRelations(ctx, eventId, eventLevelIds, scoreCategoryIds)
	if err != nil {
		return domain.ErrInternalServer
	}
	if !valid {
		return domain.ErrBadRequest
	}
	return nil
}

func (s *assessmentService) DeleteAward(ctx context.Context, eventId, userId, id uuid.UUID) error {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return err
	}
	award, err := s.repo.FindAwardById(ctx, id)
	if err != nil {
		return domain.ErrInternalServer
	}
	if award == nil || award.EventID != eventId {
		return domain.ErrNotFound
	}
	return s.repo.DeleteAward(ctx, id)
}

func (s *assessmentService) mapToAwardRes(a *entity.EventAward) *dto.AwardRes {
	var cats []dto.ScoreCategoryRes
	for _, c := range a.ScoreCategories {
		cats = append(cats, dto.ScoreCategoryRes{
			ID:   c.ID,
			Name: c.Name,
		})
	}

	var levelIDs []uuid.UUID
	var levels []dto.EventLevelResponse
	for _, l := range a.EventLevels {
		levelIDs = append(levelIDs, l.Id)
		levels = append(levels, dto.EventLevelResponse{
			Id:   l.Id,
			Name: l.Name,
		})
	}

	return &dto.AwardRes{
		ID:              a.ID,
		EventID:         a.EventID,
		EventLevelIDs:   levelIDs,
		Levels:          levels,
		Name:            a.Name,
		LimitRank:       a.LimitRank,
		ScoreCategories: cats,
	}
}
