package repository

import (
	"context"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type juryRepository struct {
	db *gorm.DB
}

func NewJuryRepository(db *gorm.DB) contracts.JuryRepository {
	return &juryRepository{db: db}
}

func (r *juryRepository) GetAll(ctx context.Context) ([]dto.JuryResponse, error) {
	var juries []entity.Jury
	if err := r.db.WithContext(ctx).Find(&juries).Error; err != nil {
		return nil, err
	}

	res := make([]dto.JuryResponse, 0)
	for _, j := range juries {
		res = append(res, dto.JuryResponse{
			Id:   j.Id,
			Name: j.Name,
		})
	}
	return res, nil
}

func (r *juryRepository) Create(ctx context.Context, req dto.JuryRequest) (dto.JuryResponse, error) {
	jury := entity.Jury{
		Id:   uuid.New(),
		Name: req.Name,
	}

	if err := r.db.WithContext(ctx).Create(&jury).Error; err != nil {
		return dto.JuryResponse{}, err
	}

	return dto.JuryResponse{Id: jury.Id, Name: jury.Name}, nil
}

func (r *juryRepository) Update(ctx context.Context, id string, req dto.JuryRequest) (dto.JuryResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return dto.JuryResponse{}, err
	}

	if err := r.db.WithContext(ctx).Model(&entity.Jury{}).Where("id = ?", uid).Update("name", req.Name).Error; err != nil {
		return dto.JuryResponse{}, err
	}

	return dto.JuryResponse{Id: uid, Name: req.Name}, nil
}

func (r *juryRepository) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Delete(&entity.Jury{}, "id = ?", uid).Error
}
