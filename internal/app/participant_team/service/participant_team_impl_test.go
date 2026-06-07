package service

import (
	"context"
	"errors"
	"testing"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/google/uuid"
)

type participantTeamRepoStub struct {
	team         *entity.Team
	deleteCalled bool
}

func (r *participantTeamRepoStub) CreateTeamWithMembers(context.Context, *entity.Team, []entity.TeamMember) error {
	return nil
}

func (r *participantTeamRepoStub) GetTeamsByInstitutionID(context.Context, uuid.UUID) ([]entity.Team, error) {
	return nil, nil
}

func (r *participantTeamRepoStub) GetTeamByID(context.Context, uuid.UUID) (*entity.Team, error) {
	return r.team, nil
}

func (r *participantTeamRepoStub) UpdateTeamWithMembers(context.Context, *entity.Team, []entity.TeamMember) error {
	return nil
}

func (r *participantTeamRepoStub) DeleteTeam(context.Context, uuid.UUID) error {
	r.deleteCalled = true
	return nil
}

type participantProfileRepoStub struct {
	institution *entity.Institution
}

func (r participantProfileRepoStub) GetUserWithInstitution(context.Context, uuid.UUID) (*entity.User, error) {
	return nil, nil
}

func (r participantProfileRepoStub) GetInstitutionByUserID(context.Context, uuid.UUID) (*entity.Institution, error) {
	return r.institution, nil
}

func (r participantProfileRepoStub) GetInstitutionByID(context.Context, uuid.UUID) (*entity.Institution, error) {
	return nil, nil
}

func (r participantProfileRepoStub) CreateInstitution(context.Context, *entity.Institution) error {
	return nil
}

func (r participantProfileRepoStub) UpdateInstitution(context.Context, *entity.Institution) error {
	return nil
}

func (r participantProfileRepoStub) GetUserByID(context.Context, uuid.UUID) (*entity.User, error) {
	return nil, nil
}

func (r participantProfileRepoStub) UpdateUserPassword(context.Context, *entity.User) error {
	return nil
}

func TestGetTeamDetailRejectsForeignTeam(t *testing.T) {
	userID := uuid.New()
	repo := &participantTeamRepoStub{
		team: &entity.Team{Id: uuid.New(), InstiId: uuid.New()},
	}
	svc := NewParticipantTeamService(repo, participantProfileRepoStub{
		institution: &entity.Institution{Id: uuid.New(), UserId: userID},
	})

	_, err := svc.GetTeamDetail(context.Background(), userID.String(), repo.team.Id.String())

	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestDeleteTeamRejectsForeignTeamBeforeDelete(t *testing.T) {
	userID := uuid.New()
	repo := &participantTeamRepoStub{
		team: &entity.Team{Id: uuid.New(), InstiId: uuid.New()},
	}
	svc := NewParticipantTeamService(repo, participantProfileRepoStub{
		institution: &entity.Institution{Id: uuid.New(), UserId: userID},
	})

	err := svc.DeleteTeam(context.Background(), userID.String(), repo.team.Id.String())

	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
	if repo.deleteCalled {
		t.Fatal("DeleteTeam repository method must not be called for a foreign team")
	}
}
