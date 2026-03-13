package contracts

import (
	"context"

	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/google/uuid"
)

type IAssessmentRepository interface {
	CheckEventOwnership(ctx context.Context, eventId, userId uuid.UUID) (bool, error)

	CreateJudge(ctx context.Context, judge *entity.Judge) error
	GetJudgesByEvent(ctx context.Context, eventId uuid.UUID) ([]entity.Judge, error)
	FindJudgeById(ctx context.Context, id uuid.UUID) (*entity.Judge, error)
	UpdateJudge(ctx context.Context, judge *entity.Judge) error
	DeleteJudge(ctx context.Context, id uuid.UUID) error

	CreateViolationType(ctx context.Context, vt *entity.ViolationType) error
	GetViolationTypesByEvent(ctx context.Context, eventId uuid.UUID) ([]entity.ViolationType, error)
	FindViolationTypeById(ctx context.Context, id uuid.UUID) (*entity.ViolationType, error)
	UpdateViolationType(ctx context.Context, vt *entity.ViolationType) error
	DeleteViolationType(ctx context.Context, id uuid.UUID) error

	CreateScoreCategory(ctx context.Context, sc *entity.ScoreCategory) error
	GetScoreCategoriesByEvent(ctx context.Context, eventId uuid.UUID) ([]entity.ScoreCategory, error)
	FindScoreCategoryById(ctx context.Context, id uuid.UUID) (*entity.ScoreCategory, error)
	UpdateScoreCategory(ctx context.Context, sc *entity.ScoreCategory) error
	DeleteScoreCategory(ctx context.Context, id uuid.UUID) error

	CreateScoreSubCategory(ctx context.Context, ssc *entity.ScoreSubCategory) error
	FindScoreSubCategoryById(ctx context.Context, id uuid.UUID) (*entity.ScoreSubCategory, error)
	UpdateScoreSubCategory(ctx context.Context, ssc *entity.ScoreSubCategory) error
	DeleteScoreSubCategory(ctx context.Context, id uuid.UUID) error

	GetGradeRulesByEvent(ctx context.Context, eventId uuid.UUID) ([]entity.GradeRule, error)
	ReplaceGradeRules(ctx context.Context, eventId uuid.UUID, rules []entity.GradeRule) error

	CreateScore(ctx context.Context, score *entity.Score) error
}

type IAssessmentService interface {
	CreateJudge(ctx context.Context, eventId, userId uuid.UUID, req dto.CreateJudgeReq) (*dto.JudgeRes, error)
	GetJudges(ctx context.Context, eventId, userId uuid.UUID) ([]dto.JudgeRes, error)
	UpdateJudge(ctx context.Context, eventId, userId, id uuid.UUID, req dto.UpdateJudgeReq) (*dto.JudgeRes, error)
	DeleteJudge(ctx context.Context, eventId, userId, id uuid.UUID) error

	CreateViolationType(ctx context.Context, eventId, userId uuid.UUID, req dto.CreateViolationTypeReq) (*dto.ViolationTypeRes, error)
	GetViolationTypes(ctx context.Context, eventId, userId uuid.UUID) ([]dto.ViolationTypeRes, error)
	UpdateViolationType(ctx context.Context, eventId, userId, id uuid.UUID, req dto.UpdateViolationTypeReq) (*dto.ViolationTypeRes, error)
	DeleteViolationType(ctx context.Context, eventId, userId, id uuid.UUID) error

	CreateScoreCategory(ctx context.Context, eventId, userId uuid.UUID, req dto.CreateScoreCategoryReq) (*dto.ScoreCategoryRes, error)
	GetScoreCategories(ctx context.Context, eventId, userId uuid.UUID) ([]dto.ScoreCategoryRes, error)
	UpdateScoreCategory(ctx context.Context, eventId, userId, id uuid.UUID, req dto.UpdateScoreCategoryReq) (*dto.ScoreCategoryRes, error)
	DeleteScoreCategory(ctx context.Context, eventId, userId, id uuid.UUID) error

	CreateScoreSubCategory(ctx context.Context, eventId, userId uuid.UUID, req dto.CreateScoreSubCategoryReq) (*dto.ScoreSubCategoryRes, error)
	UpdateScoreSubCategory(ctx context.Context, eventId, userId, id uuid.UUID, req dto.UpdateScoreSubCategoryReq) (*dto.ScoreSubCategoryRes, error)
	DeleteScoreSubCategory(ctx context.Context, eventId, userId, id uuid.UUID) error

	SetupGradeRules(ctx context.Context, eventId, userId uuid.UUID, req dto.SetupGradeRulesReq) ([]dto.GradeRuleRes, error)
	GetGradeRules(ctx context.Context, eventId, userId uuid.UUID) ([]dto.GradeRuleRes, error)

	InputScore(ctx context.Context, eventId, userId uuid.UUID, req dto.InputScoreReq) (*dto.ScoreRes, error)
}
