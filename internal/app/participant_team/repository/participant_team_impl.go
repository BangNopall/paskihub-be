package repository

import (
	"context"

	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type participantTeamRepository struct {
	db *gorm.DB
}

func NewParticipantTeamRepository(db *gorm.DB) contracts.ParticipantTeamRepository {
	return &participantTeamRepository{
		db: db,
	}
}

func (r *participantTeamRepository) CreateTeamWithMembers(ctx context.Context, team *entity.Team, members []entity.TeamMember) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(team).Error; err != nil {
			return err
		}
		
		for i := range members {
			members[i].TeamId = team.Id
		}
		
		if len(members) > 0 {
			if err := tx.Create(&members).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *participantTeamRepository) GetTeamsByInstitutionID(ctx context.Context, instiID uuid.UUID) ([]entity.Team, error) {
	var teams []entity.Team
	err := r.db.WithContext(ctx).Preload("Registrations").Where("insti_id = ?", instiID).Find(&teams).Error
	return teams, err
}

func (r *participantTeamRepository) GetTeamByID(ctx context.Context, teamID uuid.UUID) (*entity.Team, error) {
	var team entity.Team
	err := r.db.WithContext(ctx).Preload("TeamMembers").Preload("Registrations").First(&team, "id = ?", teamID).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *participantTeamRepository) DeleteTeam(ctx context.Context, teamID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Team{}, "id = ?", teamID).Error
}
