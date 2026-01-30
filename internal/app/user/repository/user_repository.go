package repository

import (
	// "math"

	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/pkg/log"
)

type userRepository struct {
	conn *gorm.DB
}

func NewUserRepository(conn *gorm.DB) contracts.UserRepository {
	return &userRepository{conn}
}

func (r *userRepository) CreateUser(user *entity.User) error {
	err := r.conn.Create(user).Error

	if err != nil {
		if err == gorm.ErrDuplicatedKey {
			return domain.ErrDuplicateEntry
		}

		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[USER REPOSITORY][Registration] failed to register user")

		return domain.ErrInternalServer
	}

	return nil
}

func (r *userRepository) DeleteUnverifiedUser() error {
	err := r.conn.Where("email_is_verified = ?", false).Delete(&entity.User{}).Error

	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[USER REPOSITORY][DeleteUnverifiedUser] failed to delete unverified user")

		return domain.ErrInternalServer
	}

	return nil
}

func (r *userRepository) FindUser(user *entity.User, userParam *dto.UserParam, relations ...string) error {
	preloadConn := r.conn

	for _, relation := range relations {
		preloadConn = preloadConn.Preload(relation)
	}

	err := preloadConn.First(user, userParam).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.ErrNotFound
		}

		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[USER REPOSITORY][FindUser] failed to find user")
		return domain.ErrInternalServer
	}

	return nil
}

func (r *userRepository) UpdateUser(updateUser *dto.UserUpdate, userId uuid.UUID) error {
	err := r.conn.Model(&entity.User{}).Where("id = ?", userId).Updates(updateUser).Error
	if err != nil {

		if err == gorm.ErrDuplicatedKey {
			return domain.ErrDuplicateEntry
		}

		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[USER REPOSITORY][UpdateUser] failed to update user")

		return domain.ErrInternalServer
	}

	return nil
}