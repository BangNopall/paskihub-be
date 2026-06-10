package repository

import (
	"context"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type eoTeamRepository struct {
	db *gorm.DB
}

func NewEOTeamRepository(db *gorm.DB) contracts.IEOTeamRepository {
	return &eoTeamRepository{
		db: db,
	}
}

func (r *eoTeamRepository) CheckEventOwnership(ctx context.Context, eventId, userId uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Event{}).
		Where("id = ? AND user_id = ?", eventId, userId).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *eoTeamRepository) FindAllRegistrationsByEvent(ctx context.Context, eventId uuid.UUID, eventLevelId *uuid.UUID, institutionType *string) ([]entity.Registration, error) {
	var registrations []entity.Registration

	query := r.db.WithContext(ctx).
		Preload("Team").
		Preload("Team.Institution").
		Preload("EventLevel").
		Joins("JOIN event_levels ON event_levels.id = registrations.event_level_id").
		Joins("JOIN teams ON teams.id = registrations.team_id").
		Joins("JOIN institutions ON institutions.id = teams.insti_id").
		Where("event_levels.event_id = ?", eventId).
		Where("registrations.is_kick = ?", false)

	if eventLevelId != nil {
		query = query.Where("registrations.event_level_id = ?", *eventLevelId)
	}

	if institutionType != nil && *institutionType != "" {
		query = query.Where("institutions.institution_type = ?", *institutionType)
	}

	err := query.Find(&registrations).Error
	return registrations, err
}

func (r *eoTeamRepository) FindRegistrationByIdAndEvent(ctx context.Context, registrationId, eventId uuid.UUID) (*entity.Registration, error) {
	var registration entity.Registration

	err := r.db.WithContext(ctx).
		Preload("Team").
		Preload("Team.TeamMembers").
		Preload("Team.Institution").
		Preload("Team.Institution.User").
		Preload("EventLevel").
		Joins("JOIN event_levels ON event_levels.id = registrations.event_level_id").
		Where("registrations.id = ?", registrationId).
		Where("event_levels.event_id = ?", eventId).
		First(&registration).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &registration, nil
}

func (r *eoTeamRepository) UpdateRegistration(ctx context.Context, registration *entity.Registration) error {
	return r.db.WithContext(ctx).Save(registration).Error
}

func (r *eoTeamRepository) ApproveRegistration(
	ctx context.Context,
	eventId uuid.UUID,
	registrationId uuid.UUID,
	totalFee float64,
	approvalFee float64,
	coinRate float64,
	status enums.RegistrationStatus,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var registration entity.Registration
		eventLevelIds := tx.Model(&entity.EventLevel{}).
			Select("id").
			Where("event_id = ?", eventId)
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND event_level_id IN (?)", registrationId, eventLevelIds).
			First(&registration).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return domain.ErrRegistrationNotFound
			}
			return err
		}
		if registration.PaymentStatus != enums.Waiting {
			return domain.ErrBadRequest
		}

		var wallet entity.Wallet
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("event_id = ?", eventId).
			First(&wallet).Error; err != nil {
			return err
		}
		if wallet.CoinBalance < approvalFee {
			return domain.ErrInsufficientBalance
		}

		approvalFeeIDR := approvalFee * coinRate

		if err := tx.Model(&wallet).Updates(map[string]interface{}{
			"saldo":        gorm.Expr("saldo - ?", approvalFeeIDR),
			"coin_balance": gorm.Expr("coin_balance - ?", approvalFee),
		}).Error; err != nil {
			return err
		}
		if err := tx.Create(&entity.WalletTransaction{
			Id:               uuid.New(),
			WalletId:         wallet.Id,
			Type:             enums.Withdraw,
			Amount:           -approvalFeeIDR,
			CoinRateSnapshot: coinRate,
			Status:           enums.Approve,
		}).Error; err != nil {
			return err
		}
		return tx.Model(&registration).Update("payment_status", status).Error
	})
}

func (r *eoTeamRepository) GetStats(ctx context.Context, eventId uuid.UUID) (*dto.EOTeamStatsRes, error) {
	var stats dto.EOTeamStatsRes

	type result struct {
		PaymentStatus string
		Count         int64
	}
	var results []result

	err := r.db.WithContext(ctx).
		Model(&entity.Registration{}).
		Joins("JOIN event_levels ON event_levels.id = registrations.event_level_id").
		Where("event_levels.event_id = ?", eventId).
		Select("payment_status, count(*) as count").
		Group("payment_status").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	for _, res := range results {
		stats.TotalTeams += res.Count
		switch res.PaymentStatus {
		case string(enums.Waiting):
			stats.PendingApproval = res.Count
		case string(enums.DPPaid):
			stats.PaidDP = res.Count
			stats.Approved += res.Count
		case string(enums.FullPaid):
			stats.PaidFull = res.Count
			stats.Approved += res.Count
		case string(enums.Rejected):
			stats.Rejected = res.Count
		}
	}

	return &stats, nil
}

func (r *eoTeamRepository) GetAssessmentStatus(ctx context.Context, registrationId uuid.UUID) (string, error) {
	var registration entity.Registration
	err := r.db.WithContext(ctx).
		Select("assessment_status").
		Where("id = ?", registrationId).
		First(&registration).Error

	if err != nil {
		return "PENDING", err
	}

	return string(registration.AssessmentStatus), nil
}
