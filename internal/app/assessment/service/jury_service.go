package service

import (
	"context"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
)

type juryService struct {
	repo contracts.JuryRepository
}

func NewJuryService(repo contracts.JuryRepository) contracts.JuryService {
	return &juryService{repo: repo}
}

func (s *juryService) GetAll(ctx context.Context) ([]dto.JuryResponse, error) {
	return s.repo.GetAll(ctx)
}

func (s *juryService) Create(ctx context.Context, req dto.JuryRequest) (dto.JuryResponse, error) {
	return s.repo.Create(ctx, req)
}

func (s *juryService) Update(ctx context.Context, id string, req dto.JuryRequest) (dto.JuryResponse, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *juryService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
