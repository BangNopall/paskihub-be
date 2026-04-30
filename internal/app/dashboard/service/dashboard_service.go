package service

import (
	"context"

	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/google/uuid"
)

type dashboardService struct {
	repo contracts.IDashboardRepository
}

func NewDashboardService(repo contracts.IDashboardRepository) contracts.IDashboardService {
	return &dashboardService{repo: repo}
}

func (s *dashboardService) GetOrganizerDashboard(ctx context.Context, userId uuid.UUID) (*dto.OrganizerDashboardRes, error) {
	stats, err := s.repo.GetOrganizerStats(ctx, userId)
	if err != nil {
		return nil, err
	}

	activities, err := s.repo.GetEORecentActivities(ctx, userId)
	if err != nil {
		return nil, err
	}

	upcoming, err := s.repo.GetEOUpcomingEvents(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &dto.OrganizerDashboardRes{
		Stats:            *stats,
		RecentActivities: activities,
		UpcomingEvents:   upcoming,
	}, nil
}

func (s *dashboardService) GetParticipantDashboard(ctx context.Context, userId uuid.UUID) (*dto.ParticipantDashboardRes, error) {
	stats, err := s.repo.GetParticipantStats(ctx, userId)
	if err != nil {
		return nil, err
	}

	activities, err := s.repo.GetParticipantRecentActivities(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &dto.ParticipantDashboardRes{
		Stats:            *stats,
		RecentActivities: activities,
	}, nil
}
