package repository

import (
	"context"
	"strings"

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

func (r *participantEventRepository) GetOpenEvents(ctx context.Context, location string, search string) ([]entity.Event, error) {
	var events []entity.Event
	query := r.db.WithContext(ctx).Preload("EventLevels").Where("status = ?", "OPEN")

	if location != "" && location != "all" {
		query = query.Where("LOWER(location) LIKE ?", "%"+strings.ToLower(location)+"%")
	}
	if search != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(search)+"%")
	}

	err := query.Find(&events).Error
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

func (r *participantEventRepository) GetRegistrationWithDetails(ctx context.Context, regisID uuid.UUID) (*entity.Registration, error) {
	var regis entity.Registration
	err := r.db.WithContext(ctx).
		Preload("Team.TeamMembers").
		Preload("EventLevel.Event").
		First(&regis, "id = ?", regisID).Error
	if err != nil {
		return nil, err
	}
	return &regis, nil
}

func (r *participantEventRepository) GetRegistrationOwnership(ctx context.Context, regisID uuid.UUID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Registration{}).
		Joins("JOIN teams ON teams.id = registrations.team_id").
		Joins("JOIN institutions ON institutions.id = teams.insti_id").
		Joins("JOIN users ON users.id = institutions.user_id").
		Where("registrations.id = ? AND users.id = ?", regisID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *participantEventRepository) GetEventLevelByID(ctx context.Context, levelID uuid.UUID) (*entity.EventLevel, error) {
	var level entity.EventLevel
	err := r.db.WithContext(ctx).Preload("Event").First(&level, "id = ?", levelID).Error
	if err != nil {
		return nil, err
	}
	return &level, nil
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

func (r *participantEventRepository) GetTeamWithMembers(ctx context.Context, teamID uuid.UUID) (*entity.Team, error) {
	var team entity.Team
	err := r.db.WithContext(ctx).
		Preload("TeamMembers").
		Preload("Institution").
		First(&team, "id = ?", teamID).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *participantEventRepository) CheckExistingRegistration(ctx context.Context, teamID, eventID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Registration{}).
		Joins("JOIN event_levels ON event_levels.id = registrations.event_level_id").
		Where("registrations.team_id = ? AND event_levels.event_id = ? AND registrations.is_kick = ?", teamID, eventID, false).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
