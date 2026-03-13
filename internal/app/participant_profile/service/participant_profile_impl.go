package service

import (
	"context"
	"errors"

	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/BangNopall/paskihub-be/pkg/bcrypt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type participantProfileService struct {
	repo contracts.ParticipantProfileRepository
}

func NewParticipantProfileService(repo contracts.ParticipantProfileRepository) contracts.ParticipantProfileService {
	return &participantProfileService{
		repo: repo,
	}
}

func (s *participantProfileService) GetProfile(ctx context.Context, userID string) (*dto.ParticipantProfileResponse, error) {
	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	user, err := s.repo.GetUserWithInstitution(ctx, parsedUUID)
	if err != nil {
		return nil, err
	}

	response := &dto.ParticipantProfileResponse{
		Email: user.Email,
	}

	if len(user.Institutions) > 0 {
		inst := user.Institutions[0]
		response.Institution = &dto.InstitutionProfileResponse{
			Name:            inst.Name,
			Address:         inst.Address,
			InstitutionType: string(inst.InstitutionType),
			NamePj:          inst.NamePj,
			NoWaPj:          inst.NoWaPj,
		}
	}

	return response, nil
}

func (s *participantProfileService) UpdateInstitution(ctx context.Context, userID string, req dto.UpdateInstitutionRequest) error {
	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	instType := enums.InstitutionType(req.InstitutionType)

	institution, err := s.repo.GetInstitutionByUserID(ctx, parsedUUID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if institution != nil {
		institution.Name = req.Name
		institution.Address = req.Address
		institution.InstitutionType = instType
		institution.NamePj = req.NamePj
		institution.NoWaPj = req.NoWaPj
		return s.repo.UpdateInstitution(ctx, institution)
	}

	newInstitution := &entity.Institution{
		UserId:          parsedUUID,
		Name:            req.Name,
		Address:         req.Address,
		InstitutionType: instType,
		NamePj:          req.NamePj,
		NoWaPj:          req.NoWaPj,
	}

	return s.repo.CreateInstitution(ctx, newInstitution)
}

func (s *participantProfileService) UpdatePassword(ctx context.Context, userID string, req dto.UpdatePasswordRequest) error {
	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	user, err := s.repo.GetUserByID(ctx, parsedUUID)
	if err != nil {
		return err
	}

	if !bcrypt.Bcrypt.Compare(req.OldPassword, user.Password) {
		return errors.New("old password does not match")
	}

	hashedPassword, err := bcrypt.Bcrypt.Hash(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.repo.UpdateUserPassword(ctx, user)
}
