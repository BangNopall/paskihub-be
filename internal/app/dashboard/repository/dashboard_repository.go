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
		TotalEvent: dto.StatValue{Value: float64(totalEvent), Trend: "+0%"},
		TotalTeam:  dto.StatValue{Value: float64(totalTeam), Trend: "+0%"},
		CoinBalance: dto.CoinValue{
			Value: wallet.Saldo,
			Coins: wallet.CoinBalance,
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
		Joins("JOIN institutions ON institutions.id = teams.insti_id").
		Where("institutions.user_id = ?", userId).
		Count(&totalTeam)

	var activeEvent int64
	r.db.WithContext(ctx).Model(&entity.Registration{}).
		Joins("JOIN teams ON teams.id = registrations.team_id").
		Joins("JOIN institutions ON institutions.id = teams.insti_id").
		Joins("JOIN event_levels ON event_levels.id = registrations.event_level_id").
		Joins("JOIN events ON events.id = event_levels.event_id").
		Where("institutions.user_id = ? AND events.compe_date >= ?", userId, time.Now()).
		Count(&activeEvent)

	var finishedEvent int64
	r.db.WithContext(ctx).Model(&entity.Registration{}).
		Joins("JOIN teams ON teams.id = registrations.team_id").
		Joins("JOIN institutions ON institutions.id = teams.insti_id").
		Joins("JOIN event_levels ON event_levels.id = registrations.event_level_id").
		Joins("JOIN events ON events.id = event_levels.event_id").
		Where("institutions.user_id = ? AND events.compe_date < ?", userId, time.Now()).
		Count(&finishedEvent)

	var pendingPayment int64
	r.db.WithContext(ctx).Model(&entity.Registration{}).
		Joins("JOIN teams ON teams.id = registrations.team_id").
		Joins("JOIN institutions ON institutions.id = teams.insti_id").
		Where("institutions.user_id = ? AND registrations.payment_status = ?", userId, enums.Waiting).
		Count(&pendingPayment)

	return &dto.ParticipantStats{
		TotalTeam:      int(totalTeam),
		ActiveEvent:    int(activeEvent),
		FinishedEvent:  int(finishedEvent),
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
		Joins("JOIN institutions ON institutions.id = teams.insti_id").
		Where("institutions.user_id = ?", userId).
		Order("registrations.updated_at DESC").
		Limit(5).
		Find(&regs)

	var res []dto.ParticipantActivity
	for _, reg := range regs {
		res = append(res, dto.ParticipantActivity{
			Title:       reg.EventLevel.Event.Name,
			Description: fmt.Sprintf("Status pendaftaran tim %s: %s", reg.Team.Name, registrationStatusLabel(reg.PaymentStatus)),
			Time:        timeSince(reg.UpdatedAt),
		})
	}
	return res, nil
}

func registrationStatusLabel(status enums.RegistrationStatus) string {
	switch status {
	case enums.Waiting:
		return "Pending"
	case enums.DPPaid:
		return "Belum Lunas/DP"
	case enums.FullPaid:
		return "Lunas"
	case enums.Rejected:
		return "Ditolak"
	case enums.Kicked:
		return "Kicked"
	default:
		return string(status)
	}
}

func (r *dashboardRepository) GetParticipantUpcomingEvents(ctx context.Context, userId uuid.UUID) ([]dto.ParticipantUpcomingEventRes, error) {
	var regs []entity.Registration
	err := r.db.WithContext(ctx).
		Preload("EventLevel.Event").
		Joins("JOIN teams ON teams.id = registrations.team_id").
		Joins("JOIN institutions ON institutions.id = teams.insti_id").
		Joins("JOIN event_levels ON event_levels.id = registrations.event_level_id").
		Joins("JOIN events ON events.id = event_levels.event_id").
		Where("institutions.user_id = ? AND events.compe_date >= ?", userId, time.Now()).
		Order("events.compe_date ASC").
		Limit(5).
		Find(&regs).Error
	if err != nil {
		return nil, err
	}

	res := make([]dto.ParticipantUpcomingEventRes, 0, len(regs))
	for _, reg := range regs {
		event := reg.EventLevel.Event

		var count int64
		r.db.WithContext(ctx).Model(&entity.Registration{}).
			Joins("JOIN event_levels ON event_levels.id = registrations.event_level_id").
			Where("event_levels.event_id = ?", event.Id).
			Count(&count)

		res = append(res, dto.ParticipantUpcomingEventRes{
			Id:              reg.Id,
			Title:           event.Name,
			Date:            event.CompeDate.Format("2006-01-02"),
			RegisteredTeams: int(count),
			Status:          event.Status,
			DetailURLId:     reg.Id,
		})
	}
	return res, nil
}

func (r *dashboardRepository) GetAdminDashboardStats(ctx context.Context) (*dto.AdminDashboardStats, error) {
	now := time.Now()
	currentMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	nextMonthStart := currentMonthStart.AddDate(0, 1, 0)
	previousMonthStart := currentMonthStart.AddDate(0, -1, 0)

	totalRevenue, err := r.sumTopupAmount(ctx, enums.Approve, time.Time{}, time.Time{})
	if err != nil {
		return nil, err
	}

	currentRevenue, err := r.sumTopupAmount(ctx, enums.Approve, currentMonthStart, nextMonthStart)
	if err != nil {
		return nil, err
	}

	previousRevenue, err := r.sumTopupAmount(ctx, enums.Approve, previousMonthStart, currentMonthStart)
	if err != nil {
		return nil, err
	}

	totalEO, currentEO, err := r.countUsersByRole(ctx, enums.Organizer, currentMonthStart, nextMonthStart)
	if err != nil {
		return nil, err
	}

	totalParticipants, currentParticipants, err := r.countUsersByRole(ctx, enums.Peserta, currentMonthStart, nextMonthStart)
	if err != nil {
		return nil, err
	}

	var pendingTopups int64
	if err := r.db.WithContext(ctx).Model(&entity.WalletTransaction{}).
		Where("type = ? AND status = ?", enums.TopUp, enums.Pending).
		Count(&pendingTopups).Error; err != nil {
		return nil, err
	}

	return &dto.AdminDashboardStats{
		TotalRevenue: dto.AdminDashboardStatValue{
			Value: totalRevenue,
			Trend: formatRevenueTrend(currentRevenue, previousRevenue),
		},
		TotalEO: dto.AdminDashboardStatValue{
			Value: float64(totalEO),
			Trend: fmt.Sprintf("+%d baru bulan ini", currentEO),
		},
		TotalParticipants: dto.AdminDashboardStatValue{
			Value: float64(totalParticipants),
			Trend: fmt.Sprintf("+%d baru bulan ini", currentParticipants),
		},
		PendingTopups: dto.AdminDashboardStatValue{
			Value: float64(pendingTopups),
			Trend: "Perlu approval",
		},
	}, nil
}

func (r *dashboardRepository) GetAdminRecentTransactions(ctx context.Context, limit int) ([]dto.AdminDashboardTransactionRes, error) {
	if limit <= 0 {
		limit = 5
	}

	coinRate, err := r.getCoinRate(ctx)
	if err != nil {
		return nil, err
	}

	var transactions []entity.WalletTransaction
	err = r.db.WithContext(ctx).
		Preload("Wallet.Event").
		Where("type = ?", enums.TopUp).
		Order("created_at DESC").
		Limit(limit).
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	res := make([]dto.AdminDashboardTransactionRes, 0, len(transactions))
	for _, tx := range transactions {
		eoName := tx.Wallet.Event.Organizer
		if eoName == "" {
			eoName = tx.Wallet.Event.Name
		}

		res = append(res, dto.AdminDashboardTransactionRes{
			Id:         tx.Id,
			EOName:     eoName,
			Amount:     tx.Amount,
			AmountKoin: tx.Amount / coinRate,
			TimeAgo:    timeSince(tx.CreatedAt),
			Status:     tx.Status,
		})
	}

	return res, nil
}

func (r *dashboardRepository) GetAdminEORegistrations(ctx context.Context, limit int) ([]dto.AdminDashboardEORegistrationRes, error) {
	if limit <= 0 {
		limit = 5
	}

	var rows []dto.AdminEORegistrationRaw
	err := r.db.WithContext(ctx).Model(&entity.User{}).
		Select(`
			users.id,
			COALESCE(NULLIF(MAX(events.organizer), ''), NULLIF(MAX(events.name), ''), users.email) AS name,
			users.email,
			users.created_at
		`).
		Joins("LEFT JOIN events ON events.user_id = users.id").
		Where("users.role = ?", enums.Organizer).
		Group("users.id, users.email, users.created_at").
		Order("users.created_at DESC").
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	res := make([]dto.AdminDashboardEORegistrationRes, 0, len(rows))
	for _, row := range rows {
		res = append(res, dto.AdminDashboardEORegistrationRes{
			Id:           row.Id,
			Name:         row.Name,
			Email:        row.Email,
			RegisteredAt: row.CreatedAt.UTC().Format(time.RFC3339),
		})
	}

	return res, nil
}

func (r *dashboardRepository) GetHomeStats(ctx context.Context) (*dto.HomeStatsResponse, error) {
	var stats dto.HomeStatsResponse

	if err := r.db.WithContext(ctx).Model(&entity.Event{}).
		Where("status <> ?", enums.Archived).
		Count(&stats.TotalEvents).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Model(&entity.User{}).
		Where("role = ? AND is_banned = ? AND email_is_verified = ? AND parent_id IS NULL", enums.Organizer, false, true).
		Count(&stats.TotalOrganizers).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Model(&entity.User{}).
		Where("role = ? AND is_banned = ? AND email_is_verified = ?", enums.Peserta, false, true).
		Count(&stats.TotalParticipants).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Model(&entity.Team{}).
		Joins("JOIN institutions ON institutions.id = teams.insti_id").
		Joins("JOIN users ON users.id = institutions.user_id").
		Where("users.role = ? AND users.is_banned = ? AND users.email_is_verified = ?", enums.Peserta, false, true).
		Count(&stats.TotalTeams).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *dashboardRepository) sumTopupAmount(ctx context.Context, status enums.TransactionStatus, start, end time.Time) (float64, error) {
	var amount float64
	query := r.db.WithContext(ctx).Model(&entity.WalletTransaction{}).
		Where("type = ? AND status = ?", enums.TopUp, status)

	if !start.IsZero() {
		query = query.Where("created_at >= ?", start)
	}
	if !end.IsZero() {
		query = query.Where("created_at < ?", end)
	}

	err := query.Select("COALESCE(SUM(amount), 0)").Scan(&amount).Error
	return amount, err
}

func (r *dashboardRepository) countUsersByRole(ctx context.Context, role enums.Role, start, end time.Time) (int64, int64, error) {
	totalQuery := r.db.WithContext(ctx).Model(&entity.User{}).
		Where("role = ? AND is_banned = ? AND email_is_verified = ?", role, false, true)

	var total int64
	if err := totalQuery.Count(&total).Error; err != nil {
		return 0, 0, err
	}

	var current int64
	if err := r.db.WithContext(ctx).Model(&entity.User{}).
		Where("role = ? AND is_banned = ? AND email_is_verified = ?", role, false, true).
		Where("created_at >= ? AND created_at < ?", start, end).
		Count(&current).Error; err != nil {
		return 0, 0, err
	}

	return total, current, nil
}

func (r *dashboardRepository) getCoinRate(ctx context.Context) (float64, error) {
	var setting entity.SystemSetting
	err := r.db.WithContext(ctx).Order("created_at ASC").First(&setting).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 1000, nil
		}
		return 0, err
	}
	if setting.CoinRate <= 0 {
		return 1000, nil
	}
	return setting.CoinRate, nil
}

func formatRevenueTrend(current, previous float64) string {
	if previous == 0 {
		if current > 0 {
			return "+100% bulan ini"
		}
		return "0% bulan ini"
	}

	percentage := ((current - previous) / previous) * 100
	if percentage >= 0 {
		return fmt.Sprintf("+%.0f%% bulan ini", percentage)
	}
	return fmt.Sprintf("%.0f%% bulan ini", percentage)
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
