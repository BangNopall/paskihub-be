package contracts

import (
	"context"
	"mime/multipart"

	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/google/uuid"
)

type EventRepository interface {
	CreateEvent(ctx context.Context, event *entity.Event) error
	UpdateEvent(ctx context.Context, updateEvent *dto.EventUpdate, eventId uuid.UUID) error
	FetchAllByConditionAndRelation(
		condition string,
		args []interface{},
		pageParms *dto.PaginationRequest,
		preload ...string,
	) ([]entity.Event, dto.PaginationResponse, error)
	DeleteEvent(ctx context.Context, eventId uuid.UUID) error
	FetchUserEvent(ctx context.Context, userId uuid.UUID) ([]entity.User, error)
	FetchOneById(ctx context.Context, eventId uuid.UUID) (entity.Event, error)
	FetchOneByParams(ctx context.Context, event *entity.Event, eventParam *dto.EventParams, relations ...string) (entity.Event, error)
	InsertEventLevel(ctx context.Context, eventLevel *entity.EventLevel) error
	FetchAllEvent(ctx context.Context, params *dto.EventParams, relations ...string) ([]entity.Event, error)
	UpdateEventLevel(ctx context.Context, eventLevel *entity.EventLevel) error
	DeleteEventLevel(ctx context.Context, eventLevelId uuid.UUID) error
}

type EventService interface {
	CreateEvent(ctx context.Context, UserId uuid.UUID, event dto.EventCreate) error
	ShowEventData(ctx context.Context, eventId uuid.UUID) (dto.EventResponse, error)
	ShowUserEvent(ctx context.Context, userId uuid.UUID) ([]dto.EventResponse, error)
	UploadLogo(ctx context.Context, Id string, UserId string, logoFile *multipart.FileHeader) error
	UploadPoster(ctx context.Context, Id string, UserId string, posterFile *multipart.FileHeader) error
	UpdateEvent(ctx context.Context, Id string, UserId string, event *dto.EventUpdate) error
	RemoveUserEvent(ctx context.Context, User *entity.User) error
	AddUserEvent(ctx context.Context, User *entity.User) error
	CreateEventLevel(ctx context.Context, EventId string, UserId string, EventLevel *dto.EventLevelCreate) error
	UpdateEventLevel(ctx context.Context, EventId string, UserId string, EventLevel *dto.EventLevelUpdate) error
	DeleteEventLevel(ctx context.Context, EventId string, EventLevelId string, UserId string) error
	DeleteEvent(ctx context.Context, EventId string, UserId string) error
}
