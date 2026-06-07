package service

import (
	"context"
	"testing"

	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/google/uuid"
)

type stubEventRepository struct {
	event             entity.Event
	updatedLevel      *entity.EventLevel
	updateParentEvent uuid.UUID
	deletedLevelID    uuid.UUID
	deleteParentEvent uuid.UUID
}

func (r *stubEventRepository) CreateEvent(ctx context.Context, event *entity.Event) error {
	return nil
}

func (r *stubEventRepository) UpdateEvent(ctx context.Context, updateEvent *entity.Event, eventID uuid.UUID) error {
	return nil
}

func (r *stubEventRepository) FetchAllByConditionAndRelation(
	condition string,
	args []interface{},
	pageParams *dto.PaginationRequest,
	preload ...string,
) ([]entity.Event, dto.PaginationResponse, error) {
	return nil, dto.PaginationResponse{}, nil
}

func (r *stubEventRepository) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	return nil
}

func (r *stubEventRepository) FetchUserEvent(ctx context.Context, userID uuid.UUID) ([]entity.User, error) {
	return nil, nil
}

func (r *stubEventRepository) FetchOneById(ctx context.Context, eventID uuid.UUID) (entity.Event, error) {
	return r.event, nil
}

func (r *stubEventRepository) FetchOneByParams(
	ctx context.Context,
	event *entity.Event,
	eventParam *dto.EventParams,
	relations ...string,
) (entity.Event, error) {
	return entity.Event{}, nil
}

func (r *stubEventRepository) InsertEventLevel(ctx context.Context, eventLevel *entity.EventLevel) error {
	return nil
}

func (r *stubEventRepository) FetchAllEvent(ctx context.Context, params *dto.EventParams, relations ...string) ([]entity.Event, error) {
	return nil, nil
}

func (r *stubEventRepository) UpdateEventLevel(ctx context.Context, eventID uuid.UUID, eventLevel *entity.EventLevel) error {
	r.updateParentEvent = eventID
	r.updatedLevel = eventLevel
	return nil
}

func (r *stubEventRepository) DeleteEventLevel(ctx context.Context, eventID, eventLevelID uuid.UUID) error {
	r.deleteParentEvent = eventID
	r.deletedLevelID = eventLevelID
	return nil
}

func (r *stubEventRepository) UpdateStatus(ctx context.Context, eventID uuid.UUID, status string) error {
	return nil
}

func TestUpdateEventLevelBindsLevelToURLParentEvent(t *testing.T) {
	ownerID := uuid.New()
	eventID := uuid.New()
	foreignEventID := uuid.New()
	levelID := uuid.New()
	repo := &stubEventRepository{event: entity.Event{Id: eventID, UserId: ownerID}}
	svc := &eventService{eventRepo: repo}

	err := svc.UpdateEventLevel(context.Background(), eventID.String(), ownerID.String(), &dto.EventLevelUpdate{
		Id:       levelID,
		EventId:  foreignEventID,
		Name:     "SMA",
		RegisFee: "100000",
		DpFee:    "50000",
	})
	if err != nil {
		t.Fatalf("UpdateEventLevel returned error: %v", err)
	}
	if repo.updateParentEvent != eventID {
		t.Fatalf("expected repository parent event %s, got %s", eventID, repo.updateParentEvent)
	}
	if repo.updatedLevel == nil || repo.updatedLevel.EventId != eventID {
		t.Fatalf("expected persisted event ID %s, got %#v", eventID, repo.updatedLevel)
	}
}

func TestDeleteEventLevelBindsDeleteToURLParentEvent(t *testing.T) {
	ownerID := uuid.New()
	eventID := uuid.New()
	levelID := uuid.New()
	repo := &stubEventRepository{event: entity.Event{Id: eventID, UserId: ownerID}}
	svc := &eventService{eventRepo: repo}

	err := svc.DeleteEventLevel(context.Background(), eventID.String(), levelID.String(), ownerID.String())
	if err != nil {
		t.Fatalf("DeleteEventLevel returned error: %v", err)
	}
	if repo.deleteParentEvent != eventID || repo.deletedLevelID != levelID {
		t.Fatalf(
			"expected delete predicate event=%s level=%s, got event=%s level=%s",
			eventID,
			levelID,
			repo.deleteParentEvent,
			repo.deletedLevelID,
		)
	}
}
