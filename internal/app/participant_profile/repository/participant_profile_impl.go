package repository

import (
	"context"

	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type participantProfileRepository struct {
	db *gorm.DB
}

func NewParticipantProfileRepository(db *gorm.DB) contracts.ParticipantProfileRepository {
	return &participantProfileRepository{
		db: db,
	}
}

func (r *participantProfileRepository) GetUserWithInstitution(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Preload("Institutions").First(&user, "id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *participantProfileRepository) GetInstitutionByUserID(ctx context.Context, userID uuid.UUID) (*entity.Institution, error) {
	var institution entity.Institution
	err := r.db.WithContext(ctx).First(&institution, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &institution, nil
}

func (r *participantProfileRepository) GetInstitutionByID(ctx context.Context, id uuid.UUID) (*entity.Institution, error) {
	var institution entity.Institution
	err := r.db.WithContext(ctx).First(&institution, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &institution, nil
}

func (r *participantProfileRepository) CreateInstitution(ctx context.Context, institution *entity.Institution) error {
	return r.db.WithContext(ctx).Create(institution).Error
}

func (r *participantProfileRepository) UpdateInstitution(ctx context.Context, institution *entity.Institution) error {
	return r.db.WithContext(ctx).Save(institution).Error
}

func (r *participantProfileRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *participantProfileRepository) UpdateUserPassword(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Model(user).Update("password", user.Password).Error
}
