package service

import (
	"context"
	"testing"

	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/google/uuid"
)

type stubDashboardRepo struct {
	participantStats    *dto.ParticipantStats
	participantActivity []dto.ParticipantActivity
	participantUpcoming []dto.ParticipantUpcomingEventRes
	adminStats          *dto.AdminDashboardStats
	adminTransactions   []dto.AdminDashboardTransactionRes
	adminRegistrations  []dto.AdminDashboardEORegistrationRes
	homeStats           *dto.HomeStatsResponse
}

func (r stubDashboardRepo) GetOrganizerStats(ctx context.Context, userId uuid.UUID) (*dto.OrganizerStats, error) {
	return nil, nil
}

func (r stubDashboardRepo) GetEORecentActivities(ctx context.Context, userId uuid.UUID) ([]dto.EOActivityRes, error) {
	return nil, nil
}

func (r stubDashboardRepo) GetEOUpcomingEvents(ctx context.Context, userId uuid.UUID) ([]dto.UpcomingEventRes, error) {
	return nil, nil
}

func (r stubDashboardRepo) GetParticipantStats(ctx context.Context, userId uuid.UUID) (*dto.ParticipantStats, error) {
	return r.participantStats, nil
}

func (r stubDashboardRepo) GetParticipantRecentActivities(ctx context.Context, userId uuid.UUID) ([]dto.ParticipantActivity, error) {
	return r.participantActivity, nil
}

func (r stubDashboardRepo) GetParticipantUpcomingEvents(ctx context.Context, userId uuid.UUID) ([]dto.ParticipantUpcomingEventRes, error) {
	return r.participantUpcoming, nil
}

func (r stubDashboardRepo) GetAdminDashboardStats(ctx context.Context) (*dto.AdminDashboardStats, error) {
	return r.adminStats, nil
}

func (r stubDashboardRepo) GetAdminRecentTransactions(ctx context.Context, limit int) ([]dto.AdminDashboardTransactionRes, error) {
	return r.adminTransactions, nil
}

func (r stubDashboardRepo) GetAdminEORegistrations(ctx context.Context, limit int) ([]dto.AdminDashboardEORegistrationRes, error) {
	return r.adminRegistrations, nil
}

func (r stubDashboardRepo) GetHomeStats(ctx context.Context) (*dto.HomeStatsResponse, error) {
	return r.homeStats, nil
}

func TestGetParticipantDashboardIncludesUpcomingEvents(t *testing.T) {
	userID := uuid.New()
	registrationID := uuid.New()

	svc := NewDashboardService(stubDashboardRepo{
		participantStats: &dto.ParticipantStats{
			TotalTeam:      2,
			ActiveEvent:    1,
			FinishedEvent:  1,
			PendingPayment: 1,
		},
		participantActivity: []dto.ParticipantActivity{
			{
				Title:       "Lomba Paskibra Nasional 2026",
				Description: "Status pendaftaran tim Garuda Muda: WAITING",
				Time:        "5 jam yang lalu",
			},
		},
		participantUpcoming: []dto.ParticipantUpcomingEventRes{
			{
				Id:              registrationID,
				Title:           "Lomba Paskibra Nasional 2026",
				Date:            "2026-03-15",
				RegisteredTeams: 18,
				Status:          "OPEN",
				DetailURLId:     registrationID,
			},
		},
	})

	res, err := svc.GetParticipantDashboard(context.Background(), userID)
	if err != nil {
		t.Fatalf("GetParticipantDashboard returned error: %v", err)
	}

	if len(res.UpcomingEvents) != 1 {
		t.Fatalf("expected 1 upcoming event, got %d", len(res.UpcomingEvents))
	}

	upcoming := res.UpcomingEvents[0]
	if upcoming.DetailURLId != registrationID {
		t.Fatalf("expected detail_url_id %q, got %q", registrationID, upcoming.DetailURLId)
	}
	if upcoming.RegisteredTeams != 18 {
		t.Fatalf("expected registered_teams 18, got %d", upcoming.RegisteredTeams)
	}
}

func TestGetHomeStatsReturnsPublicAggregateStats(t *testing.T) {
	svc := NewDashboardService(stubDashboardRepo{
		homeStats: &dto.HomeStatsResponse{
			TotalEvents:       120,
			TotalOrganizers:   45,
			TotalParticipants: 980,
			TotalTeams:        210,
		},
	})

	res, err := svc.GetHomeStats(context.Background())
	if err != nil {
		t.Fatalf("GetHomeStats returned error: %v", err)
	}

	if res.TotalEvents != 120 {
		t.Fatalf("expected total_events 120, got %d", res.TotalEvents)
	}
	if res.TotalOrganizers != 45 {
		t.Fatalf("expected total_organizers 45, got %d", res.TotalOrganizers)
	}
	if res.TotalParticipants != 980 {
		t.Fatalf("expected total_participants 980, got %d", res.TotalParticipants)
	}
	if res.TotalTeams != 210 {
		t.Fatalf("expected total_teams 210, got %d", res.TotalTeams)
	}
}

func TestGetAdminDashboardCombinesStatsTransactionsAndRegistrations(t *testing.T) {
	transactionID := uuid.New()
	eoID := uuid.New()

	svc := NewDashboardService(stubDashboardRepo{
		adminStats: &dto.AdminDashboardStats{
			TotalRevenue: dto.AdminDashboardStatValue{
				Value: 125000000,
				Trend: "+12% bulan ini",
			},
			TotalEO: dto.AdminDashboardStatValue{
				Value: 42,
				Trend: "+5 baru",
			},
			TotalParticipants: dto.AdminDashboardStatValue{
				Value: 1240,
				Trend: "+156 baru",
			},
			PendingTopups: dto.AdminDashboardStatValue{
				Value: 8,
				Trend: "Perlu approval",
			},
		},
		adminTransactions: []dto.AdminDashboardTransactionRes{
			{
				Id:         transactionID,
				EOName:     "SMA 1 Jakarta",
				Amount:     500000,
				AmountKoin: 500,
				TimeAgo:    "10 menit yang lalu",
				Status:     "PENDING",
			},
		},
		adminRegistrations: []dto.AdminDashboardEORegistrationRes{
			{
				Id:           eoID,
				Name:         "Lomba Jaya Abadi",
				Email:        "contact@lombajaya.com",
				RegisteredAt: "2026-04-24T00:00:00Z",
			},
		},
	})

	res, err := svc.GetAdminDashboard(context.Background())
	if err != nil {
		t.Fatalf("GetAdminDashboard returned error: %v", err)
	}

	if res.Stats.TotalRevenue.Value != 125000000 {
		t.Fatalf("expected total revenue 125000000, got %v", res.Stats.TotalRevenue.Value)
	}
	if len(res.RecentTransactions) != 1 {
		t.Fatalf("expected 1 recent transaction, got %d", len(res.RecentTransactions))
	}
	if res.RecentTransactions[0].Id != transactionID {
		t.Fatalf("expected transaction id %q, got %q", transactionID, res.RecentTransactions[0].Id)
	}
	if len(res.EORegistrations) != 1 {
		t.Fatalf("expected 1 EO registration, got %d", len(res.EORegistrations))
	}
	if res.EORegistrations[0].Id != eoID {
		t.Fatalf("expected EO registration id %q, got %q", eoID, res.EORegistrations[0].Id)
	}
}
