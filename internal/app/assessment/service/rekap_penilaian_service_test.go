package service

import (
	"context"
	"errors"
	"testing"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/google/uuid"
)

type stubRekapRepository struct {
	registrationOwned bool
	eventLevelOwned   bool
	readCalled        *bool
	updateCalled      *bool
}

func (r stubRekapRepository) RegistrationBelongsToOrganizer(ctx context.Context, regisID, organizerID uuid.UUID) (bool, error) {
	return r.registrationOwned, nil
}

func (r stubRekapRepository) EventLevelBelongsToOrganizer(ctx context.Context, eventLevelID, organizerID uuid.UUID) (bool, error) {
	return r.eventLevelOwned, nil
}

func (r stubRekapRepository) GetTeamAssessmentDetail(ctx context.Context, regisID uuid.UUID) (dto.TeamAssessmentDetailResponse, error) {
	if r.readCalled != nil {
		*r.readCalled = true
	}
	return dto.TeamAssessmentDetailResponse{}, nil
}

func (r stubRekapRepository) GetScoreboardByEventLevel(ctx context.Context, eventLevelID uuid.UUID) ([]dto.ScoreboardItem, error) {
	if r.readCalled != nil {
		*r.readCalled = true
	}
	return nil, nil
}

func (r stubRekapRepository) GetLeaderboardCustom(ctx context.Context, eventLevelID uuid.UUID, categoryIDs []uuid.UUID) ([]dto.ScoreboardItem, error) {
	if r.readCalled != nil {
		*r.readCalled = true
	}
	return nil, nil
}

func (r stubRekapRepository) UpdateScorePublishedStatus(ctx context.Context, eventLevelID uuid.UUID, isPublished bool) error {
	if r.updateCalled != nil {
		*r.updateCalled = true
	}
	return nil
}

func (r stubRekapRepository) GetEventIDByEventLevelID(ctx context.Context, eventLevelID uuid.UUID) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func TestRekapOperationsRejectForeignOrganizerObjectsBeforeAccess(t *testing.T) {
	organizerID := uuid.New()
	regisID := uuid.New()
	eventLevelID := uuid.New()
	categoryID := uuid.New()

	tests := []struct {
		name string
		run  func(*rekapService) error
	}{
		{
			name: "team detail",
			run: func(svc *rekapService) error {
				_, err := svc.GetTeamAssessmentDetail(context.Background(), regisID, organizerID)
				return err
			},
		},
		{
			name: "scoreboard",
			run: func(svc *rekapService) error {
				_, err := svc.GetScoreboard(context.Background(), eventLevelID, organizerID)
				return err
			},
		},
		{
			name: "custom leaderboard",
			run: func(svc *rekapService) error {
				_, err := svc.GetLeaderboardCustom(
					context.Background(),
					dto.CustomLeaderboardRequest{ScoreCategoryIDs: []uuid.UUID{categoryID}},
					eventLevelID,
					organizerID,
				)
				return err
			},
		},
		{
			name: "publish scoreboard",
			run: func(svc *rekapService) error {
				return svc.PublishScoreboard(
					context.Background(),
					dto.PublishScoreboardRequest{IsScorePublished: true},
					eventLevelID,
					organizerID,
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readCalled := false
			updateCalled := false
			svc := &rekapService{
				repo: stubRekapRepository{
					registrationOwned: false,
					eventLevelOwned:   false,
					readCalled:        &readCalled,
					updateCalled:      &updateCalled,
				},
			}

			err := tt.run(svc)
			if !errors.Is(err, domain.ErrForbidden) {
				t.Fatalf("expected ErrForbidden, got %v", err)
			}
			if readCalled || updateCalled {
				t.Fatal("rekap data was accessed for a foreign organizer object")
			}
		})
	}
}
