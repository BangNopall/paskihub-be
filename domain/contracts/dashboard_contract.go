package contracts

import (
	"context"

	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/google/uuid"
)

type IDashboardRepository interface {
	GetOrganizerStats(ctx context.Context, userId uuid.UUID) (*dto.OrganizerStats, error)
	GetEORecentActivities(ctx context.Context, userId uuid.UUID) ([]dto.EOActivityRes, error)
	GetEOUpcomingEvents(ctx context.Context, userId uuid.UUID) ([]dto.UpcomingEventRes, error)

	GetParticipantStats(ctx context.Context, userId uuid.UUID) (*dto.ParticipantStats, error)
	GetParticipantRecentActivities(ctx context.Context, userId uuid.UUID) ([]dto.ParticipantActivity, error)
}

type IDashboardService interface {
	GetOrganizerDashboard(ctx context.Context, userId uuid.UUID) (*dto.OrganizerDashboardRes, error)
	GetParticipantDashboard(ctx context.Context, userId uuid.UUID) (*dto.ParticipantDashboardRes, error)
}
