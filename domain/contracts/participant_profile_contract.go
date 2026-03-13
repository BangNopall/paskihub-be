package contracts

import (
	"context"

	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/google/uuid"
)

type ParticipantProfileRepository interface {
	GetUserWithInstitution(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	GetInstitutionByUserID(ctx context.Context, userID uuid.UUID) (*entity.Institution, error)
	GetInstitutionByID(ctx context.Context, id uuid.UUID) (*entity.Institution, error)
	CreateInstitution(ctx context.Context, institution *entity.Institution) error
	UpdateInstitution(ctx context.Context, institution *entity.Institution) error
	
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	UpdateUserPassword(ctx context.Context, user *entity.User) error
}

type ParticipantProfileService interface {
	GetProfile(ctx context.Context, userID string) (*dto.ParticipantProfileResponse, error)
	UpdateInstitution(ctx context.Context, userID string, req dto.UpdateInstitutionRequest) error
	UpdatePassword(ctx context.Context, userID string, req dto.UpdatePasswordRequest) error
}
