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
	GetParticipantUpcomingEvents(ctx context.Context, userId uuid.UUID) ([]dto.ParticipantUpcomingEventRes, error)

	GetAdminDashboardStats(ctx context.Context) (*dto.AdminDashboardStats, error)
	GetAdminRecentTransactions(ctx context.Context, limit int) ([]dto.AdminDashboardTransactionRes, error)
	GetAdminEORegistrations(ctx context.Context, limit int) ([]dto.AdminDashboardEORegistrationRes, error)
}

type IDashboardService interface {
	GetOrganizerDashboard(ctx context.Context, userId uuid.UUID) (*dto.OrganizerDashboardRes, error)
	GetParticipantDashboard(ctx context.Context, userId uuid.UUID) (*dto.ParticipantDashboardRes, error)
	GetAdminDashboard(ctx context.Context) (*dto.AdminDashboardRes, error)
}
