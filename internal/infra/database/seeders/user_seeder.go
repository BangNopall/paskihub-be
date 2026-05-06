package seeders

import (
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/BangNopall/paskihub-be/pkg/bcrypt"
	"github.com/BangNopall/paskihub-be/pkg/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UserSeeder(db *gorm.DB) {
	var count int64
	db.Model(&entity.User{}).Where("email = ?", "admin@paskihub.com").Count(&count)

	if count > 0 {
		log.Info(nil, "[USER SEEDER] Admin user already exists, skipping...")
		return
	}

	hashedPassword, err := bcrypt.Bcrypt.Hash("admin123")
	if err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[USER SEEDER] Failed to hash admin password")
	}

	admin := entity.User{
		Id:              uuid.New(),
		Email:           "admin@paskihub.com",
		Password:        hashedPassword,
		Role:            enums.Admin,
		EmailIsVerified: true,
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[USER SEEDER] Failed to create admin user")
	}

	log.Info(nil, "[USER SEEDER] Admin user seeded successfully")
}
