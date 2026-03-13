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
	UpdateRegistration(ctx context.Context, registration *entity.Registration) error
	GetActiveRegistrationsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Registration, error)
}

type ParticipantEventService interface {
	GetOpenEvents(ctx context.Context) ([]dto.OpenEventResponse, error)
	RegisterEvent(ctx context.Context, req dto.RegisterEventRequest) error
	PelunasanEvent(ctx context.Context, regisID string, req dto.PelunasanEventRequest) error
	GetActiveEvents(ctx context.Context, userID string) ([]dto.ActiveEventResponse, error)
}
