package service

import (
	"context"

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
	vt := &entity.ViolationType{
		EventID: eventId,
		Name:    req.Name,
		Point:   req.Point,
	}
	if err := s.repo.CreateViolationType(ctx, vt); err != nil {
		return nil, domain.ErrInternalServer
	}
	return &dto.ViolationTypeRes{
		ID:      vt.ID,
		EventID: vt.EventID,
		Name:    vt.Name,
		Point:   vt.Point,
	}, nil
}

func (s *assessmentService) GetViolationTypes(ctx context.Context, eventId, userId uuid.UUID) ([]dto.ViolationTypeRes, error) {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}
	vts, err := s.repo.GetViolationTypesByEvent(ctx, eventId)
	if err != nil {
		return nil, domain.ErrInternalServer
	}
	var res []dto.ViolationTypeRes
	for _, vt := range vts {
		res = append(res, dto.ViolationTypeRes{
			ID:      vt.ID,
			EventID: vt.EventID,
			Name:    vt.Name,
			Point:   vt.Point,
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
	sc := &entity.ScoreCategory{
		EventID: eventId,
		Name:    req.Name,
	}
	if err := s.repo.CreateScoreCategory(ctx, sc); err != nil {
		return nil, domain.ErrInternalServer
	}
	return &dto.ScoreCategoryRes{
		ID:      sc.ID,
		EventID: sc.EventID,
		Name:    sc.Name,
	}, nil
}

func (s *assessmentService) GetScoreCategories(ctx context.Context, eventId, userId uuid.UUID) ([]dto.ScoreCategoryRes, error) {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}
	scs, err := s.repo.GetScoreCategoriesByEvent(ctx, eventId)
	if err != nil {
		return nil, domain.ErrInternalServer
	}
	var res []dto.ScoreCategoryRes
	for _, sc := range scs {
		var subcats []dto.ScoreSubCategoryRes
		for _, sub := range sc.SubCategories {
			subcats = append(subcats, dto.ScoreSubCategoryRes{
				ID:                sub.ID,
				ScoreCategoriesID: sub.ScoreCategoriesID,
				Name:              sub.Name,
			})
		}
		res = append(res, dto.ScoreCategoryRes{
			ID:            sc.ID,
			EventID:       sc.EventID,
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
	}
	if err := s.repo.CreateScoreSubCategory(ctx, ssc); err != nil {
		return nil, domain.ErrInternalServer
	}
	return &dto.ScoreSubCategoryRes{
		ID:                ssc.ID,
		ScoreCategoriesID: ssc.ScoreCategoriesID,
		Name:              ssc.Name,
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
	if err := s.repo.UpdateScoreSubCategory(ctx, ssc); err != nil {
		return nil, domain.ErrInternalServer
	}
	return &dto.ScoreSubCategoryRes{
		ID:                ssc.ID,
		ScoreCategoriesID: ssc.ScoreCategoriesID,
		Name:              ssc.Name,
	}, nil
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

// GradeRules
func (s *assessmentService) SetupGradeRules(ctx context.Context, eventId, userId uuid.UUID, req dto.SetupGradeRulesReq) ([]dto.GradeRuleRes, error) {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}

	var entities []entity.GradeRule
	for _, rule := range req.Rules {
		entities = append(entities, entity.GradeRule{
			EventID:   eventId,
			GradeName: rule.GradeName,
			MinScore:  rule.MinScore,
			MaxScore:  rule.MaxScore,
		})
	}

	if err := s.repo.ReplaceGradeRules(ctx, eventId, entities); err != nil {
		return nil, domain.ErrInternalServer
	}

	savedRules, err := s.repo.GetGradeRulesByEvent(ctx, eventId)
	if err != nil {
		return nil, domain.ErrInternalServer
	}

	var results []dto.GradeRuleRes
	for _, r := range savedRules {
		results = append(results, dto.GradeRuleRes{
			ID:        r.ID,
			EventID:   r.EventID,
			GradeName: r.GradeName,
			MinScore:  r.MinScore,
			MaxScore:  r.MaxScore,
		})
	}
	return results, nil
}

func (s *assessmentService) GetGradeRules(ctx context.Context, eventId, userId uuid.UUID) ([]dto.GradeRuleRes, error) {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}
	savedRules, err := s.repo.GetGradeRulesByEvent(ctx, eventId)
	if err != nil {
		return nil, domain.ErrInternalServer
	}

	var results []dto.GradeRuleRes
	for _, r := range savedRules {
		results = append(results, dto.GradeRuleRes{
			ID:        r.ID,
			EventID:   r.EventID,
			GradeName: r.GradeName,
			MinScore:  r.MinScore,
			MaxScore:  r.MaxScore,
		})
	}
	return results, nil
}

// Score
func (s *assessmentService) InputScore(ctx context.Context, eventId, userId uuid.UUID, req dto.InputScoreReq) (*dto.ScoreRes, error) {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}

	rules, err := s.repo.GetGradeRulesByEvent(ctx, eventId)
	if err != nil {
		return nil, domain.ErrInternalServer
	}

	if len(rules) == 0 {
		return nil, domain.ErrBadRequest // Or appropriate error to force setting grade rules
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
