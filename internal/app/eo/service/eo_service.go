package service

import (
	"context"
	"time"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/BangNopall/paskihub-be/pkg/bcrypt"
	pkgUuid "github.com/BangNopall/paskihub-be/pkg/uuid"
	"github.com/google/uuid"
)

type eoService struct {
	eoRepo   contracts.IEORepository
	userRepo contracts.UserRepository
	bcrypt   bcrypt.BcryptInterface
	uuid     pkgUuid.UUIDInterface
	timeout  time.Duration
}

func NewEOService(
	eoRepo contracts.IEORepository,
	userRepo contracts.UserRepository,
	bcrypt bcrypt.BcryptInterface,
	uuid pkgUuid.UUIDInterface,
	timeout time.Duration,
) contracts.IEOService {
	return &eoService{
		eoRepo:   eoRepo,
		userRepo: userRepo,
		bcrypt:   bcrypt,
		uuid:     uuid,
		timeout:  timeout,
	}
}

func (s *eoService) GetProfile(ctx context.Context, userId uuid.UUID) (dto.EOProfileRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var user entity.User
	err := s.userRepo.FindUser(&user, &dto.UserParam{ID: userId})
	if err != nil {
		return dto.EOProfileRes{}, domain.ErrUserNotFound
	}

	return dto.EOProfileRes{
		Email: user.Email,
	}, nil
}

func (s *eoService) UpdatePassword(ctx context.Context, userId uuid.UUID, req dto.EOUpdatePasswordReq) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var user entity.User
	err := s.userRepo.FindUser(&user, &dto.UserParam{ID: userId})
	if err != nil {
		return domain.ErrUserNotFound
	}

	if !s.bcrypt.Compare(req.OldPassword, user.Password) {
		return domain.ErrWrongEmailOrPassword // Or a more specific error if available
	}

	hashPassword, err := s.bcrypt.Hash(req.NewPassword)
	if err != nil {
		return err
	}

	updateUser := dto.UserUpdate{
		Password: hashPassword,
	}

	return s.userRepo.UpdateUser(&updateUser, userId)
}

func (s *eoService) GetStaffs(ctx context.Context, parentId uuid.UUID) ([]dto.EOStaffRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	staffs, err := s.eoRepo.FindStaffsByParentId(parentId)
	if err != nil {
		return nil, err
	}

	var res []dto.EOStaffRes
	for _, staff := range staffs {
		status := "Unverified"
		if staff.EmailIsVerified {
			status = "Verified"
		}
		res = append(res, dto.EOStaffRes{
			Id:        staff.Id.String(),
			Email:     staff.Email,
			Status:    status,
			CreatedAt: staff.CreatedAt.Format("02 Jan 2006"),
		})
	}

	return res, nil
}

func (s *eoService) CreateStaff(ctx context.Context, parentId uuid.UUID, req dto.EOStaffCreateReq) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// Ensure email is not already registered
	var existingUser entity.User
	err := s.userRepo.FindUser(&existingUser, &dto.UserParam{Email: req.Email})
	if err == nil {
		return domain.ErrEmailRegistered
	}

	hashPassword, err := s.bcrypt.Hash(req.Password)
	if err != nil {
		return err
	}

	newId, err := s.uuid.New()
	if err != nil {
		return err
	}

	newStaff := entity.User{
		Id:              newId,
		Email:           req.Email,
		Role:            enums.Organizer,
		Password:        hashPassword,
		EmailIsVerified: true, // Automatically verified as requested
		ParentId:        &parentId,
	}

	return s.userRepo.CreateUser(&newStaff)
}

func (s *eoService) ResetStaffPassword(ctx context.Context, parentId uuid.UUID, staffId uuid.UUID, req dto.EOStaffResetPasswordReq) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// Verify staff belongs to this parent
	_, err := s.eoRepo.FindStaffById(staffId, parentId)
	if err != nil {
		return domain.ErrNotFound
	}

	hashPassword, err := s.bcrypt.Hash(req.Password)
	if err != nil {
		return err
	}

	updateUser := dto.UserUpdate{
		Password: hashPassword,
	}

	return s.userRepo.UpdateUser(&updateUser, staffId)
}

func (s *eoService) DeleteStaff(ctx context.Context, parentId uuid.UUID, staffId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// Verify staff belongs to this parent
	_, err := s.eoRepo.FindStaffById(staffId, parentId)
	if err != nil {
		return domain.ErrNotFound
	}

	return s.eoRepo.DeleteStaff(staffId, parentId)
}
