package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) contracts.IDashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetOrganizerStats(ctx context.Context, userId uuid.UUID) (*dto.OrganizerStats, error) {
	var totalEvent int64
	r.db.WithContext(ctx).Model(&entity.Event{}).Where("user_id = ?", userId).Count(&totalEvent)

	var totalTeam int64
	r.db.WithContext(ctx).Model(&entity.Registration{}).
		Joins("JOIN event_levels ON event_levels.id = registrations.event_level_id").
		Joins("JOIN events ON events.id = event_levels.event_id").
		Where("events.user_id = ?", userId).
		Count(&totalTeam)

	var wallet entity.Wallet
	r.db.WithContext(ctx).
		Joins("JOIN events ON events.id = wallets.event_id").
		Where("events.user_id = ?", userId).
		First(&wallet)

	// Revenue: Sum of all FULL_PAID registrations for events owned by user
	// This is a bit simplified. In reality, revenue might come from WalletTransaction as well.
	var revenue float64
	// Sum of FULL_PAID registrations * regis_fee
	// This needs complex join and conversion. For now, let's use a simpler way or mock it.
	// Let's just sum WalletTransactions with status SUCCESS for the user's wallet
	r.db.WithContext(ctx).Model(&entity.WalletTransaction{}).
		Joins("JOIN wallets ON wallets.id = wallet_transactions.wallet_id").
		Joins("JOIN events ON events.id = wallets.event_id").
		Where("events.user_id = ? AND wallet_transactions.status = ?", userId, enums.Approve).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&revenue)

	return &dto.OrganizerStats{
		TotalEvent: dto.StatValue{Value: totalEvent, Trend: "+0%"},
		TotalTeam:  dto.StatValue{Value: totalTeam, Trend: "+0%"},
		CoinBalance: dto.CoinValue{
			Value: wallet.Saldo,
			Coins: wallet.Saldo / 1000, // Example conversion
		},
		Revenue: dto.StatValue{Value: revenue, Trend: "+0%"},
	}, nil
}

func (r *dashboardRepository) GetEORecentActivities(ctx context.Context, userId uuid.UUID) ([]dto.EOActivityRes, error) {
	var regs []entity.Registration
	err := r.db.WithContext(ctx).
		Preload("Team").
		Preload("EventLevel.Event").
		Joins("JOIN event_levels ON event_levels.id = registrations.event_level_id").
		Joins("JOIN events ON events.id = event_levels.event_id").
		Where("events.user_id = ?", userId).
		Order("registrations.created_at DESC").
		Limit(5).
		Find(&regs).Error

	if err != nil {
		return nil, err
	}

	var res []dto.EOActivityRes
	for _, reg := range regs {
		res = append(res, dto.EOActivityRes{
			Id:        reg.Id,
			TeamName:  reg.Team.Name,
			EventName: reg.EventLevel.Event.Name,
			TimeAgo:   timeSince(reg.CreatedAt),
			Status:    string(reg.PaymentStatus),
		})
	}
	return res, nil
}

func (r *dashboardRepository) GetEOUpcomingEvents(ctx context.Context, userId uuid.UUID) ([]dto.UpcomingEventRes, error) {
	var events []entity.Event
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND compe_date > ?", userId, time.Now()).
		Order("compe_date ASC").
		Limit(5).
		Find(&events).Error

	if err != nil {
		return nil, err
	}

	var res []dto.UpcomingEventRes
	for _, e := range events {
		var count int64
		r.db.WithContext(ctx).Model(&entity.Registration{}).
			Joins("JOIN event_levels ON event_levels.id = registrations.event_level_id").
			Where("event_levels.event_id = ?", e.Id).
			Count(&count)

		res = append(res, dto.UpcomingEventRes{
			Id:              e.Id,
			Title:           e.Name,
			Date:            e.CompeDate.Format("2006-01-02"),
			RegisteredTeams: int(count),
			Status:          e.Status,
		})
	}
	return res, nil
}

func (r *dashboardRepository) GetParticipantStats(ctx context.Context, userId uuid.UUID) (*dto.ParticipantStats, error) {
	var totalTeam int64
	r.db.WithContext(ctx).Model(&entity.Team{}).
		Joins("JOIN participant_profiles ON participant_profiles.id = teams.participant_id").
		Where("participant_profiles.user_id = ?", userId).
		Count(&totalTeam)

	var activeEvent int64
	r.db.WithContext(ctx).Model(&entity.Registration{}).
		Joins("JOIN teams ON teams.id = registrations.team_id").
		Joins("JOIN participant_profiles ON participant_profiles.id = teams.participant_id").
		Joins("JOIN event_levels ON event_levels.id = registrations.event_level_id").
		Joins("JOIN events ON events.id = event_levels.event_id").
		Where("participant_profiles.user_id = ? AND events.compe_date > ?", userId, time.Now()).
		Count(&activeEvent)

	var pendingPayment int64
	r.db.WithContext(ctx).Model(&entity.Registration{}).
		Joins("JOIN teams ON teams.id = registrations.team_id").
		Joins("JOIN participant_profiles ON participant_profiles.id = teams.participant_id").
		Where("participant_profiles.user_id = ? AND registrations.payment_status = ?", userId, enums.Waiting).
		Count(&pendingPayment)

	return &dto.ParticipantStats{
		TotalTeam:      int(totalTeam),
		ActiveEvent:    int(activeEvent),
		FinishedEvent:  0, // Need logic for finished
		PendingPayment: int(pendingPayment),
	}, nil
}

func (r *dashboardRepository) GetParticipantRecentActivities(ctx context.Context, userId uuid.UUID) ([]dto.ParticipantActivity, error) {
	// For participant, maybe track status changes or new registrations
	var regs []entity.Registration
	r.db.WithContext(ctx).
		Preload("Team").
		Preload("EventLevel.Event").
		Joins("JOIN teams ON teams.id = registrations.team_id").
		Joins("JOIN participant_profiles ON participant_profiles.id = teams.participant_id").
		Where("participant_profiles.user_id = ?", userId).
		Order("registrations.updated_at DESC").
		Limit(5).
		Find(&regs)

	var res []dto.ParticipantActivity
	for _, reg := range regs {
		res = append(res, dto.ParticipantActivity{
			Title:       reg.EventLevel.Event.Name,
			Description: fmt.Sprintf("Status pendaftaran tim %s: %s", reg.Team.Name, reg.PaymentStatus),
			Time:        timeSince(reg.UpdatedAt),
		})
	}
	return res, nil
}

func timeSince(t time.Time) string {
	d := time.Since(t)
	if d.Hours() > 24 {
		return fmt.Sprintf("%d hari yang lalu", int(d.Hours()/24))
	}
	if d.Hours() > 1 {
		return fmt.Sprintf("%d jam yang lalu", int(d.Hours()))
	}
	return fmt.Sprintf("%d menit yang lalu", int(d.Minutes()))
}
