package repository

import (
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type eoRepository struct {
	db *gorm.DB
}

func NewEORepository(db *gorm.DB) contracts.IEORepository {
	return &eoRepository{db: db}
}

func (r *eoRepository) FindStaffsByParentId(parentId uuid.UUID) ([]entity.User, error) {
	var staffs []entity.User
	err := r.db.Where("parent_id = ?", parentId).Find(&staffs).Error
	return staffs, err
}

func (r *eoRepository) FindStaffById(staffId uuid.UUID, parentId uuid.UUID) (*entity.User, error) {
	var staff entity.User
	err := r.db.Where("id = ? AND parent_id = ?", staffId, parentId).First(&staff).Error
	return &staff, err
}

func (r *eoRepository) DeleteStaff(staffId uuid.UUID, parentId uuid.UUID) error {
	return r.db.Where("id = ? AND parent_id = ?", staffId, parentId).Delete(&entity.User{}).Error
}
