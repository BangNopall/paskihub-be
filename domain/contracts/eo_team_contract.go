package contracts

import (
	"context"

	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/google/uuid"
)

type IEOTeamRepository interface {
	CheckEventOwnership(ctx context.Context, eventId, userId uuid.UUID) (bool, error)
	FindAllRegistrationsByEvent(ctx context.Context, eventId uuid.UUID, eventLevelId *uuid.UUID, institutionType *string) ([]entity.Registration, error)
	FindRegistrationByIdAndEvent(ctx context.Context, registrationId, eventId uuid.UUID) (*entity.Registration, error)
	UpdateRegistration(ctx context.Context, registration *entity.Registration) error
	ApproveRegistration(ctx context.Context, eventId, registrationId uuid.UUID, totalFee float64, status enums.RegistrationStatus) error
	GetStats(ctx context.Context, eventId uuid.UUID) (*dto.EOTeamStatsRes, error)
	GetAssessmentStatus(ctx context.Context, registrationId uuid.UUID) (string, error)
}

type IEOTeamService interface {
	GetListTeam(ctx context.Context, eventId, userId uuid.UUID, eventLevelId *uuid.UUID, institutionType *string) ([]dto.EOTeamListRes, error)
	GetDetailTeam(ctx context.Context, eventId, userId, registrationId uuid.UUID) (*dto.EOTeamDetailRes, error)
	ApproveTeam(ctx context.Context, eventId, userId, registrationId uuid.UUID, req dto.EOTeamApproveReq) error
	RejectTeam(ctx context.Context, eventId, userId, registrationId uuid.UUID, req dto.EOTeamRejectReq) error
	KickTeam(ctx context.Context, eventId, userId, registrationId uuid.UUID) error
	GetStats(ctx context.Context, eventId, userId uuid.UUID) (*dto.EOTeamStatsRes, error)
	StartAssessment(ctx context.Context, eventId, userId, registrationId uuid.UUID) error
}
