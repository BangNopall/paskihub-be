package service

import (
	"context"
	"errors"
	"testing"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/google/uuid"
)

type stubAssessmentRepository struct {
	eventLevelValid bool
	scoreInputValid bool
	awardValid      bool
	writeCalled     *bool
}

func (r stubAssessmentRepository) CheckEventOwnership(ctx context.Context, eventID, userID uuid.UUID) (bool, error) {
	return true, nil
}

func (r stubAssessmentRepository) EventLevelBelongsToEvent(ctx context.Context, eventID, eventLevelID uuid.UUID) (bool, error) {
	return r.eventLevelValid, nil
}

func (r stubAssessmentRepository) ValidateScoreInputRelations(ctx context.Context, eventID, regisID, judgeID, subCategoryID uuid.UUID) (bool, error) {
	return r.scoreInputValid, nil
}

func (r stubAssessmentRepository) ValidateAwardRelations(ctx context.Context, eventID uuid.UUID, eventLevelIDs, scoreCategoryIDs []uuid.UUID) (bool, error) {
	return r.awardValid, nil
}

func (r stubAssessmentRepository) CreateJudge(ctx context.Context, judge *entity.Judge) error {
	return nil
}
func (r stubAssessmentRepository) GetJudgesByEvent(ctx context.Context, eventID uuid.UUID) ([]entity.Judge, error) {
	return nil, nil
}
func (r stubAssessmentRepository) FindJudgeById(ctx context.Context, id uuid.UUID) (*entity.Judge, error) {
	return nil, nil
}
func (r stubAssessmentRepository) UpdateJudge(ctx context.Context, judge *entity.Judge) error {
	return nil
}
func (r stubAssessmentRepository) DeleteJudge(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (r stubAssessmentRepository) CreateViolationType(ctx context.Context, vt *entity.ViolationType) error {
	if r.writeCalled != nil {
		*r.writeCalled = true
	}
	return nil
}
func (r stubAssessmentRepository) GetViolationTypesByLevel(ctx context.Context, eventID, levelID uuid.UUID) ([]entity.ViolationType, error) {
	return nil, nil
}
func (r stubAssessmentRepository) FindViolationTypeById(ctx context.Context, id uuid.UUID) (*entity.ViolationType, error) {
	return nil, nil
}
func (r stubAssessmentRepository) UpdateViolationType(ctx context.Context, vt *entity.ViolationType) error {
	return nil
}
func (r stubAssessmentRepository) DeleteViolationType(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (r stubAssessmentRepository) CreateScoreCategory(ctx context.Context, sc *entity.ScoreCategory) error {
	if r.writeCalled != nil {
		*r.writeCalled = true
	}
	return nil
}
func (r stubAssessmentRepository) GetScoreCategoriesByLevel(ctx context.Context, eventID, levelID uuid.UUID) ([]entity.ScoreCategory, error) {
	return nil, nil
}
func (r stubAssessmentRepository) FindScoreCategoryById(ctx context.Context, id uuid.UUID) (*entity.ScoreCategory, error) {
	return nil, nil
}
func (r stubAssessmentRepository) UpdateScoreCategory(ctx context.Context, sc *entity.ScoreCategory) error {
	return nil
}
func (r stubAssessmentRepository) DeleteScoreCategory(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (r stubAssessmentRepository) CreateScoreSubCategory(ctx context.Context, ssc *entity.ScoreSubCategory) error {
	return nil
}
func (r stubAssessmentRepository) FindScoreSubCategoryById(ctx context.Context, id uuid.UUID) (*entity.ScoreSubCategory, error) {
	return nil, nil
}
func (r stubAssessmentRepository) UpdateScoreSubCategory(ctx context.Context, ssc *entity.ScoreSubCategory) error {
	return nil
}
func (r stubAssessmentRepository) DeleteScoreSubCategory(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (r stubAssessmentRepository) GetGradeRulesBySubCategory(ctx context.Context, subCategoryID uuid.UUID) ([]entity.GradeRule, error) {
	return []entity.GradeRule{{GradeName: "A", MinScore: 0, MaxScore: 100}}, nil
}
func (r stubAssessmentRepository) ReplaceSubCategoryGradeRules(ctx context.Context, subCategoryID uuid.UUID, rules []entity.GradeRule) error {
	return nil
}
func (r stubAssessmentRepository) CreateScore(ctx context.Context, score *entity.Score) error {
	if r.writeCalled != nil {
		*r.writeCalled = true
	}
	return nil
}
func (r stubAssessmentRepository) CreateAward(ctx context.Context, award *entity.EventAward) error {
	if r.writeCalled != nil {
		*r.writeCalled = true
	}
	return nil
}
func (r stubAssessmentRepository) GetAwardsByEvent(ctx context.Context, eventID uuid.UUID) ([]entity.EventAward, error) {
	return nil, nil
}
func (r stubAssessmentRepository) FindAwardById(ctx context.Context, id uuid.UUID) (*entity.EventAward, error) {
	return &entity.EventAward{ID: id}, nil
}
func (r stubAssessmentRepository) UpdateAward(ctx context.Context, award *entity.EventAward, levelIDs []uuid.UUID, categoryIDs []uuid.UUID) error {
	if r.writeCalled != nil {
		*r.writeCalled = true
	}
	return nil
}
func (r stubAssessmentRepository) DeleteAward(ctx context.Context, id uuid.UUID) error {
	return nil
}

func TestAssessmentCreateRejectsForeignEventLevelBeforeWrite(t *testing.T) {
	eventID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name string
		run  func(*assessmentService) error
	}{
		{
			name: "violation type",
			run: func(svc *assessmentService) error {
				_, err := svc.CreateViolationType(context.Background(), eventID, userID, dto.CreateViolationTypeReq{
					EventLevelID: uuid.New(),
					Name:         "Violation",
					Point:        1,
				})
				return err
			},
		},
		{
			name: "score category",
			run: func(svc *assessmentService) error {
				_, err := svc.CreateScoreCategory(context.Background(), eventID, userID, dto.CreateScoreCategoryReq{
					EventLevelID: uuid.New(),
					Name:         "Category",
				})
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeCalled := false
			svc := &assessmentService{repo: stubAssessmentRepository{eventLevelValid: false, writeCalled: &writeCalled}}
			err := tt.run(svc)
			if !errors.Is(err, domain.ErrBadRequest) {
				t.Fatalf("expected ErrBadRequest, got %v", err)
			}
			if writeCalled {
				t.Fatal("write called for foreign event level")
			}
		})
	}
}

func TestInputScoreRejectsMixedEventRelationsBeforeWrite(t *testing.T) {
	writeCalled := false
	svc := &assessmentService{repo: stubAssessmentRepository{scoreInputValid: false, writeCalled: &writeCalled}}

	_, err := svc.InputScore(context.Background(), uuid.New(), uuid.New(), dto.InputScoreReq{
		RegisID:       uuid.New(),
		JudgesID:      uuid.New(),
		SubCategoryID: uuid.New(),
		ScoreValue:    90,
	})
	if !errors.Is(err, domain.ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
	if writeCalled {
		t.Fatal("score write called for mixed event relations")
	}
}

func TestAwardRejectsForeignRelationsBeforeWrite(t *testing.T) {
	writeCalled := false
	svc := &assessmentService{repo: stubAssessmentRepository{awardValid: false, writeCalled: &writeCalled}}

	_, err := svc.CreateAward(context.Background(), uuid.New(), uuid.New(), dto.CreateAwardReq{
		EventLevelIDs:    []uuid.UUID{uuid.New()},
		Name:             "Champion",
		LimitRank:        3,
		ScoreCategoryIDs: []uuid.UUID{uuid.New()},
	})
	if !errors.Is(err, domain.ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
	if writeCalled {
		t.Fatal("award write called for foreign relations")
	}
}
