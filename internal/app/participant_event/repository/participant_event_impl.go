package repository

import (
	"context"

	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type participantEventRepository struct {
	db *gorm.DB
}

func NewParticipantEventRepository(db *gorm.DB) contracts.ParticipantEventRepository {
	return &participantEventRepository{
		db: db,
	}
}

func (r *participantEventRepository) GetOpenEvents(ctx context.Context) ([]entity.Event, error) {
	var events []entity.Event
	err := r.db.WithContext(ctx).Preload("EventLevels").Where("status = ?", "OPEN").Find(&events).Error
	return events, err
}

func (r *participantEventRepository) CreateRegistration(ctx context.Context, registration *entity.Registration) error {
	return r.db.WithContext(ctx).Create(registration).Error
}

func (r *participantEventRepository) GetRegistrationByID(ctx context.Context, regisID uuid.UUID) (*entity.Registration, error) {
	var regis entity.Registration
	err := r.db.WithContext(ctx).First(&regis, "id = ?", regisID).Error
	if err != nil {
		return nil, err
	}
	return &regis, nil
}

func (r *participantEventRepository) UpdateRegistration(ctx context.Context, registration *entity.Registration) error {
	return r.db.WithContext(ctx).Save(registration).Error
}

func (r *participantEventRepository) GetActiveRegistrationsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Registration, error) {
	var registrations []entity.Registration

	// Preload necessary fields
	err := r.db.WithContext(ctx).
		Joins("JOIN teams ON teams.id = registrations.team_id").
		Joins("JOIN institutions ON institutions.id = teams.insti_id").
		Joins("JOIN users ON users.id = institutions.user_id").
		Where("users.id = ?", userID).
		Preload("Team").
		Preload("EventLevel.Event").
		Find(&registrations).Error

	return registrations, err
}
