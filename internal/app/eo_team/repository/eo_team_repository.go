package repository

import (
	"context"

	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

func (r *eoTeamRepository) FindAllRegistrationsByEvent(ctx context.Context, eventId uuid.UUID, eventLevelId *uuid.UUID) ([]entity.Registration, error) {
	var registrations []entity.Registration

	query := r.db.WithContext(ctx).
		Preload("Team").
		Preload("Team.Institution").
		Preload("EventLevel").
		Joins("JOIN event_levels ON event_levels.id = registrations.event_level_id").
		Where("event_levels.event_id = ?", eventId).
		Where("registrations.is_kick = ?", false)

	if eventLevelId != nil {
		query = query.Where("registrations.event_level_id = ?", *eventLevelId)
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
		Preload("EventLevel").
		Joins("JOIN event_levels ON event_levels.id = registrations.event_level_id").
		Where("registrations.id = ?", registrationId).
		Where("event_levels.event_id = ?", eventId).
		First(&registration).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Or return a specific custom error if desired, but returning nil is okay if the service handles it
		}
		return nil, err
	}

	return &registration, nil
}

func (r *eoTeamRepository) UpdateRegistration(ctx context.Context, registration *entity.Registration) error {
	return r.db.WithContext(ctx).Save(registration).Error
}
