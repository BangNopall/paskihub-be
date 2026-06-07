package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/google/uuid"
)

type stubParticipantEventRepo struct {
	registration *entity.Registration
	isOwner      bool
	updateCalled *bool
}

func (r stubParticipantEventRepo) GetOpenEvents(ctx context.Context) ([]entity.Event, error) {
	return nil, nil
}

func (r stubParticipantEventRepo) CreateRegistration(ctx context.Context, registration *entity.Registration) error {
	return nil
}

func (r stubParticipantEventRepo) GetRegistrationByID(ctx context.Context, regisID uuid.UUID) (*entity.Registration, error) {
	return r.registration, nil
}

func (r stubParticipantEventRepo) GetRegistrationWithDetails(ctx context.Context, regisID uuid.UUID) (*entity.Registration, error) {
	return r.registration, nil
}

func (r stubParticipantEventRepo) GetRegistrationOwnership(ctx context.Context, regisID uuid.UUID, userID uuid.UUID) (bool, error) {
	return r.isOwner, nil
}

func (r stubParticipantEventRepo) UpdateRegistration(ctx context.Context, registration *entity.Registration) error {
	if r.updateCalled != nil {
		*r.updateCalled = true
	}
	return nil
}

func (r stubParticipantEventRepo) GetActiveRegistrationsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Registration, error) {
	return nil, nil
}

func (r stubParticipantEventRepo) GetEventLevelByID(ctx context.Context, levelID uuid.UUID) (*entity.EventLevel, error) {
	return nil, nil
}

func (r stubParticipantEventRepo) GetTeamWithMembers(ctx context.Context, teamID uuid.UUID) (*entity.Team, error) {
	return nil, nil
}

func (r stubParticipantEventRepo) CheckExistingRegistration(ctx context.Context, teamID, eventID uuid.UUID) (bool, error) {
	return false, nil
}

type stubRekapRepo struct{}

func (r stubRekapRepo) RegistrationBelongsToOrganizer(ctx context.Context, regisID, organizerID uuid.UUID) (bool, error) {
	return false, nil
}

func (r stubRekapRepo) EventLevelBelongsToOrganizer(ctx context.Context, eventLevelID, organizerID uuid.UUID) (bool, error) {
	return false, nil
}

func (r stubRekapRepo) GetTeamAssessmentDetail(ctx context.Context, regisID uuid.UUID) (dto.TeamAssessmentDetailResponse, error) {
	return dto.TeamAssessmentDetailResponse{}, nil
}

func (r stubRekapRepo) GetScoreboardByEventLevel(ctx context.Context, eventLevelID uuid.UUID) ([]dto.ScoreboardItem, error) {
	return nil, nil
}

func (r stubRekapRepo) GetLeaderboardCustom(ctx context.Context, eventLevelID uuid.UUID, categoryIDs []uuid.UUID) ([]dto.ScoreboardItem, error) {
	return nil, nil
}

func (r stubRekapRepo) UpdateScorePublishedStatus(ctx context.Context, eventLevelID uuid.UUID, isPublished bool) error {
	return nil
}

func (r stubRekapRepo) GetEventIDByEventLevelID(ctx context.Context, eventLevelID uuid.UUID) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func TestGetRegistrationDetailIncludesFrontendOverviewFields(t *testing.T) {
	userID := uuid.New()
	eventID := uuid.New()
	eventLevelID := uuid.New()
	teamID := uuid.New()
	regisID := uuid.New()

	svc := &participantEventService{
		repo: stubParticipantEventRepo{
			isOwner: true,
			registration: &entity.Registration{
				Id:               regisID,
				TeamId:           teamID,
				EventLevelId:     eventLevelID,
				PaymentStatus:    enums.DPPaid,
				PaymentProofPath: "/public/proof.png",
				Team: entity.Team{
					Id:       teamID,
					Name:     "Garuda",
					LogoPath: "/public/team.png",
					TeamMembers: []entity.TeamMember{
						{Role: enums.Official},
						{Role: enums.Pasukan},
					},
				},
				EventLevel: entity.EventLevel{
					Id:       eventLevelID,
					EventId:  eventID,
					RegisFee: "250000",
					DpFee:    "100000",
					Event: entity.Event{
						Id:          eventID,
						Name:        "Paskibra Cup",
						Description: "Event desc",
						LogoPath:    "/public/event-logo.png",
						CompeDate:   time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
						Location:    "Jakarta",
					},
				},
			},
		},
		rekapRepo: stubRekapRepo{},
	}

	res, err := svc.GetRegistrationDetail(context.Background(), userID.String(), regisID.String())
	if err != nil {
		t.Fatalf("GetRegistrationDetail returned error: %v", err)
	}

	if res.EventLevelId != eventLevelID.String() {
		t.Fatalf("expected event_level_id %q, got %q", eventLevelID.String(), res.EventLevelId)
	}
	if res.Event.LogoUrl != "/public/event-logo.png" {
		t.Fatalf("expected event logo url to be mapped, got %q", res.Event.LogoUrl)
	}
}

func TestPelunasanEventRejectsForeignRegistrationBeforeUpdate(t *testing.T) {
	updateCalled := false
	svc := &participantEventService{
		repo: stubParticipantEventRepo{
			isOwner:      false,
			registration: &entity.Registration{Id: uuid.New()},
			updateCalled: &updateCalled,
		},
	}

	err := svc.PelunasanEvent(
		context.Background(),
		uuid.NewString(),
		uuid.NewString(),
		dto.PelunasanEventRequest{},
	)

	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
	if updateCalled {
		t.Fatal("UpdateRegistration was called for a foreign registration")
	}
}
