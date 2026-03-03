package repository

import (
	"context"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/pkg/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type eventRepository struct {
	conn *gorm.DB
}

func NewEventRepository(conn *gorm.DB) contracts.EventRepository {
	return &eventRepository{conn: conn}
}

func (r *eventRepository) CreateEvent(ctx context.Context, event *entity.Event) error {
	err := r.conn.Create(event).Error
	if err != nil {
		if err == gorm.ErrDuplicatedKey {
			return domain.ErrDuplicateEntry
		}

		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[EVENT REPOSITORY][CreateEvent] failed to create event")

		return domain.ErrInternalServer
	}
	return nil
}

func (r *eventRepository) UpdateEvent(ctx context.Context, updateEvent *dto.EventUpdate, eventId uuid.UUID) error {
	err := r.conn.Model(&entity.Event{}).Where("id = ?", eventId).Updates(updateEvent).Error
	if err != nil {
		if err == gorm.ErrDuplicatedKey {
			return domain.ErrDuplicateEntry
		}

		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[EVENT REPOSITORY][UpdateEvent] failed to update event")

		return domain.ErrInternalServer
	}
	return nil
}

func (r *eventRepository) FetchAllByConditionAndRelation(
	condition string,
	args []interface{},
	pageParms *dto.PaginationRequest,
	preload ...string,
) ([]entity.Event, dto.PaginationResponse, error) {
	var events []entity.Event
	var count int64

	query := r.conn.Model(&entity.Event{})

	if condition != "" {
		query = query.Where(condition, args...)
	}

	for _, p := range preload {
		query = query.Preload(p)
	}

	err := query.Count(&count).Error
	if err != nil {
		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[EVENT REPOSITORY][FetchAllByConditionAndRelation] failed to count events")
		return nil, dto.PaginationResponse{}, domain.ErrInternalServer
	}

	offset := (pageParms.Page - 1) * pageParms.Limit
	err = query.Limit(pageParms.Limit).Offset(offset).Find(&events).Error
	if err != nil {
		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[EVENT REPOSITORY][FetchAllByConditionAndRelation] failed to fetch events")
		return nil, dto.PaginationResponse{}, domain.ErrInternalServer
	}

	paginationResponse := dto.PaginationResponse{
		Page:       pageParms.Page,
		TotalPages: int((count + int64(pageParms.Limit) - 1) / int64(pageParms.Limit)),
	}

	return events, paginationResponse, nil
}

func (r *eventRepository) DeleteEvent(ctx context.Context, eventId uuid.UUID) error {
	err := r.conn.Where("id = ?", eventId).Delete(&entity.Event{}).Error
	if err != nil {
		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[EVENT REPOSITORY][DeleteEvent] failed to delete event")
		return domain.ErrInternalServer
	}
	return nil
}

func (r *eventRepository) FetchUserEvent(ctx context.Context, userId uuid.UUID) ([]entity.User, error) {
	var users []entity.User
	err := r.conn.Preload("Events").Where("id = ?", userId).Find(&users).Error
	if err != nil {
		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[EVENT REPOSITORY][FetchUserEvent] failed to fetch user events")
		return nil, domain.ErrInternalServer
	}
	return users, nil
}

func (r *eventRepository) FetchOneById(ctx context.Context, eventId uuid.UUID) (entity.Event, error) {
	var event entity.Event
	err := r.conn.Preload("User").Preload("EventLevels").Where("id = ?", eventId).First(&event).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.Event{}, domain.ErrNotFound
		}
		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[EVENT REPOSITORY][FetchOneById] failed to fetch event")
		return entity.Event{}, domain.ErrInternalServer
	}
	return event, nil
}

func (r *eventRepository) FetchOneByParams(ctx context.Context, event *entity.Event, eventParam *dto.EventParams, relations ...string) (entity.Event, error) {
	query := r.conn.Model(&entity.Event{})

	if eventParam != nil {
		if eventParam.Id != uuid.Nil {
			query = query.Where("id = ?", eventParam.Id)
		}
		if eventParam.Name != "" {
			query = query.Where("LOWER(name) = LOWER(?)", eventParam.Name)
		}
		if eventParam.Status != "" {
			query = query.Where("LOWER(status) = LOWER(?)", eventParam.Status)
		}
	}

	for _, relation := range relations {
		query = query.Preload(relation)
	}

	var result entity.Event
	err := query.First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.Event{}, domain.ErrNotFound
		}
		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[EVENT REPOSITORY][FetchOneByParams] failed to fetch event")
		return entity.Event{}, domain.ErrInternalServer
	}

	return result, nil
}

func (r *eventRepository) InsertEventLevel(ctx context.Context, eventLevel *entity.EventLevel) error {
	err := r.conn.Create(eventLevel).Error
	if err != nil {
		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[EVENT REPOSITORY][InsertEventLevel] failed to insert event level")
		return domain.ErrInternalServer
	}
	return nil
}

func (r *eventRepository) FetchAllEvent(ctx context.Context, params *dto.EventParams, relations ...string) ([]entity.Event, error) {
	var events []entity.Event
	query := r.conn.Model(&entity.Event{})

	if params != nil {
		if params.Name != "" {
			query = query.Where("name ILIKE ?", "%"+params.Name+"%")
		}
		if params.Status != "" {
			query = query.Where("status = ?", params.Status)
		}
	}

	for _, relation := range relations {
		query = query.Preload(relation)
	}

	err := query.Find(&events).Error
	if err != nil {
		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[EVENT REPOSITORY][FetchAllEvent] failed to fetch all events")
		return nil, domain.ErrInternalServer
	}

	return events, nil
}

func (r *eventRepository) UpdateEventLevel(ctx context.Context, eventLevel *entity.EventLevel) error {
	err := r.conn.Save(eventLevel).Error
	if err != nil {
		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[EVENT REPOSITORY][UpdateEventLevel] failed to update event level")
		return domain.ErrInternalServer
	}
	return nil
}

func (r *eventRepository) DeleteEventLevel(ctx context.Context, eventLevelId uuid.UUID) error {
	err := r.conn.Delete(&entity.EventLevel{}, eventLevelId).Error
	if err != nil {
		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[EVENT REPOSITORY][DeleteEventLevel] failed to delete event level")
		return domain.ErrInternalServer
	}
	return nil
}
