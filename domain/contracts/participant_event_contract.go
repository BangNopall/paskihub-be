package contracts

import (
	"context"

	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/google/uuid"
)

type ParticipantEventRepository interface {
	GetOpenEvents(ctx context.Context) ([]entity.Event, error)
	CreateRegistration(ctx context.Context, registration *entity.Registration) error
	GetRegistrationByID(ctx context.Context, regisID uuid.UUID) (*entity.Registration, error)
	GetRegistrationWithDetails(ctx context.Context, regisID uuid.UUID) (*entity.Registration, error)
	UpdateRegistration(ctx context.Context, registration *entity.Registration) error
	GetActiveRegistrationsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Registration, error)
	GetEventLevelByID(ctx context.Context, levelID uuid.UUID) (*entity.EventLevel, error)
	GetTeamWithMembers(ctx context.Context, teamID uuid.UUID) (*entity.Team, error)
	CheckExistingRegistration(ctx context.Context, teamID, eventID uuid.UUID) (bool, error)
}

type ParticipantEventService interface {
	GetOpenEvents(ctx context.Context) ([]dto.OpenEventResponse, error)
	RegisterEvent(ctx context.Context, userID string, req dto.RegisterEventRequest) error
	PelunasanEvent(ctx context.Context, regisID string, req dto.PelunasanEventRequest) error
	GetActiveEvents(ctx context.Context, userID string) ([]dto.ActiveEventResponse, error)
	GetRegistrationDetail(ctx context.Context, regisID string) (dto.RegistrationDetailResponse, error)
	GetScoreboard(ctx context.Context, eventLevelID string) (dto.ScoreboardResponse, error)
}
