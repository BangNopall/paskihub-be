package contracts

import (
	"context"

	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/google/uuid"
)

type ParticipantTeamRepository interface {
	CreateTeamWithMembers(ctx context.Context, team *entity.Team, members []entity.TeamMember) error
	GetTeamsByInstitutionID(ctx context.Context, instiID uuid.UUID) ([]entity.Team, error)
	GetTeamByID(ctx context.Context, teamID uuid.UUID) (*entity.Team, error)
	DeleteTeam(ctx context.Context, teamID uuid.UUID) error
}

type ParticipantTeamService interface {
	CreateTeam(ctx context.Context, userID string, req dto.CreateTeamRequest) error
	UpdateTeam(ctx context.Context, userID string, teamID string, req dto.CreateTeamRequest) error // optional if identical struct
	GetTeams(ctx context.Context, userID string) ([]dto.ParticipantTeamResponse, error)
	GetTeamDetail(ctx context.Context, userID string, teamID string) (*dto.TeamDetailResponse, error)
	DeleteTeam(ctx context.Context, userID string, teamID string) error
}
