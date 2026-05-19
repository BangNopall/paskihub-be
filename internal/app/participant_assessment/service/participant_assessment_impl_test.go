package service

import (
	"context"
	"testing"

	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/google/uuid"
)

type stubParticipantAssessmentRepo struct {
	registration *entity.Registration
	categories   []entity.ScoreCategory
}

func (r stubParticipantAssessmentRepo) GetScoresByRegisID(ctx context.Context, regisID uuid.UUID) ([]entity.Score, error) {
	return nil, nil
}

func (r stubParticipantAssessmentRepo) GetViolationsByRegisID(ctx context.Context, regisID uuid.UUID) ([]entity.TeamViolation, error) {
	return nil, nil
}

func (r stubParticipantAssessmentRepo) GetRegistrationOwnership(ctx context.Context, regisID uuid.UUID, userID uuid.UUID) (bool, error) {
	return true, nil
}

func (r stubParticipantAssessmentRepo) GetRegistrationByID(ctx context.Context, regisID uuid.UUID) (*entity.Registration, error) {
	return r.registration, nil
}

func (r stubParticipantAssessmentRepo) GetScoreCategoriesByEventLevel(ctx context.Context, eventID uuid.UUID, eventLevelID uuid.UUID) ([]entity.ScoreCategory, error) {
	return r.categories, nil
}

func TestGetAssessmentRecapIncludesMaxScoreForEventLevel(t *testing.T) {
	eventID := uuid.New()
	eventLevelID := uuid.New()
	regisID := uuid.New()

	svc := &participantAssessmentService{
		repo: stubParticipantAssessmentRepo{
			registration: &entity.Registration{
				Id:           regisID,
				EventLevelId: eventLevelID,
				EventLevel: entity.EventLevel{
					Id:      eventLevelID,
					EventId: eventID,
					Event: entity.Event{
						Id: eventID,
					},
				},
			},
			categories: []entity.ScoreCategory{
				{
					EventID:      eventID,
					EventLevelID: eventLevelID,
					SubCategories: []entity.ScoreSubCategory{
						{MaxScore: 100},
						{MaxScore: 150},
					},
				},
				{
					EventID:      eventID,
					EventLevelID: eventLevelID,
					SubCategories: []entity.ScoreSubCategory{
						{MaxScore: 250},
					},
				},
			},
		},
	}

	res, err := svc.GetAssessmentRecap(context.Background(), uuid.NewString(), regisID.String())
	if err != nil {
		t.Fatalf("GetAssessmentRecap returned error: %v", err)
	}

	if res.MaxScore != 500 {
		t.Fatalf("expected max_score 500, got %v", res.MaxScore)
	}
}
