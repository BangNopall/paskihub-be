package contracts

import (
	"context"

	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/google/uuid"
)

type IEORepository interface {
	FindStaffsByParentId(parentId uuid.UUID) ([]entity.User, error)
	FindStaffById(staffId uuid.UUID, parentId uuid.UUID) (*entity.User, error)
	DeleteStaff(staffId uuid.UUID, parentId uuid.UUID) error
}

type IEOService interface {
	GetProfile(ctx context.Context, userId uuid.UUID) (dto.EOProfileRes, error)
	UpdatePassword(ctx context.Context, userId uuid.UUID, req dto.EOUpdatePasswordReq) error
	GetStaffs(ctx context.Context, parentId uuid.UUID) ([]dto.EOStaffRes, error)
	CreateStaff(ctx context.Context, parentId uuid.UUID, req dto.EOStaffCreateReq) error
	ResetStaffPassword(ctx context.Context, parentId uuid.UUID, staffId uuid.UUID, req dto.EOStaffResetPasswordReq) error
	DeleteStaff(ctx context.Context, parentId uuid.UUID, staffId uuid.UUID) error
}
