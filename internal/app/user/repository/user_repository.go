package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/domain/enums"
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
	result := r.conn.Where("email_is_verified = ? AND created_at < ?", false, time.Now().Add(-24*time.Hour)).Delete(&entity.User{})
	err := result.Error

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
	result := r.conn.Model(&entity.User{}).Where("id = ?", userId).Updates(updateUser)
	err := result.Error
	if err != nil {

		if err == gorm.ErrDuplicatedKey {
			return domain.ErrDuplicateEntry
		}

		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[USER REPOSITORY][UpdateUser] failed to update user")

		return domain.ErrInternalServer
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *userRepository) FetchAllUsers(ctx context.Context, role *string) ([]entity.User, error) {
	var users []entity.User
	query := r.conn.WithContext(ctx).Model(&entity.User{}).Preload("Institutions").Preload("Events.EventLevels")
	if role != nil {
		query = query.Where("role = ?", *role)
	}
	err := query.Find(&users).Error
	return users, err
}

func (r *userRepository) UpdateUserStatus(ctx context.Context, userId uuid.UUID, isBanned bool) error {
	result := r.conn.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userId).Update("is_banned", isBanned)
	err := result.Error
	if err != nil {
		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[USER REPOSITORY][UpdateUserStatus] failed to update user status")
		return domain.ErrInternalServer
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *userRepository) VerifyUserEmail(ctx context.Context, userId uuid.UUID) error {
	result := r.conn.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userId).Update("email_is_verified", true)
	err := result.Error
	if err != nil {
		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[USER REPOSITORY][VerifyUserEmail] failed to verify user email")
		return domain.ErrInternalServer
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *userRepository) FetchUserDetailAdmin(ctx context.Context, userId uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.conn.WithContext(ctx).
		Preload("Institutions.Teams.TeamMembers").
		Preload("Institutions.Teams.Registrations.EventLevel.Event").
		Preload("Events.EventLevels").
		Where("id = ?", userId).
		First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[USER REPOSITORY][FetchUserDetailAdmin] failed to fetch user detail")
		return nil, domain.ErrInternalServer
	}

	// For Organizer (EO), we also need the staff/panitia
	if user.Role == enums.Organizer {
		var staffs []entity.User
		r.conn.WithContext(ctx).Where("parent_id = ?", userId).Find(&staffs)
	}

	return &user, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, userId uuid.UUID) error {
	result := r.conn.WithContext(ctx).Delete(&entity.User{}, userId)
	err := result.Error
	if err != nil {
		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[USER REPOSITORY][DeleteUser] failed to delete user")
		return domain.ErrInternalServer
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
