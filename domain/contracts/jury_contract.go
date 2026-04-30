package contracts

import (
	"context"
	"github.com/BangNopall/paskihub-be/domain/dto"
)

type JuryService interface {
	GetAll(ctx context.Context) ([]dto.JuryResponse, error)
	Create(ctx context.Context, req dto.JuryRequest) (dto.JuryResponse, error)
	Update(ctx context.Context, id string, req dto.JuryRequest) (dto.JuryResponse, error)
	Delete(ctx context.Context, id string) error
}

type JuryRepository interface {
	GetAll(ctx context.Context) ([]dto.JuryResponse, error)
	Create(ctx context.Context, req dto.JuryRequest) (dto.JuryResponse, error)
	Update(ctx context.Context, id string, req dto.JuryRequest) (dto.JuryResponse, error)
	Delete(ctx context.Context, id string) error
}
