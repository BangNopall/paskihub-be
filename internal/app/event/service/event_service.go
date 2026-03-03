package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	timePkg "github.com/BangNopall/paskihub-be/pkg/time"
	uuidPkg "github.com/BangNopall/paskihub-be/pkg/uuid"
	"github.com/google/uuid"
)

type eventService struct {
	eventRepo  contracts.EventRepository
	walletRepo contracts.WalletRepository
	uuid       uuidPkg.UUIDInterface
	time       timePkg.TimeInterface
	timeout    time.Duration
}

func NewEventService(
	eventRepo contracts.EventRepository,
	walletRepo contracts.WalletRepository,
	uuid uuidPkg.UUIDInterface,
	time timePkg.TimeInterface,
	timeout time.Duration,
) contracts.EventService {
	return &eventService{
		eventRepo,
		walletRepo,
		uuid,
		time,
		timeout,
	}
}

func (s *eventService) CreateEvent(ctx context.Context, UserId uuid.UUID, event dto.EventCreate) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	_, err := s.eventRepo.FetchOneByParams(ctx, &entity.Event{}, &dto.EventParams{Name: event.Name})
	if err == nil {
		return domain.ErrDuplicateEntry
	}
	if err != domain.ErrNotFound {
		return err
	}

	eventId, err := s.uuid.New()
	if err != nil {
		return err
	}

	layout := "2006-01-02 15:04:05"
	openDate, _ := time.Parse(layout, event.OpenDate)
	closeDate, _ := time.Parse(layout, event.CloseDate)

	newEvent := entity.Event{
		Id:         eventId,
		UserId:     UserId,
		Name:       event.Name,
		OpenDate:   openDate,
		CloseDate:  closeDate,
		Location:   event.Location,
		NamePj:     event.NamaPj,
		NoWaPj:     event.NoWaPj,
		BankName:   event.BankName,
		BankNumber: event.BankNumber,
		Status:     "DRAFT",
	}

	err = s.eventRepo.CreateEvent(ctx, &newEvent)
	if err != nil {
		return err
	}

	walletId, err := s.uuid.New()
	if err != nil {
		return err
	}

	newWallet := entity.Wallet{
		Id:      walletId,
		EventId: eventId,
		Saldo:   0,
	}

	err = s.walletRepo.CreateWallet(ctx, &newWallet)
	if err != nil {
		return err
	}

	return nil
}

func (s *eventService) ShowEventData(ctx context.Context, eventId uuid.UUID) (dto.EventResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	event, err := s.eventRepo.FetchOneById(ctx, eventId)
	if err != nil {
		return dto.EventResponse{}, err
	}

	return *dto.EventEntityToResponse(&event), nil
}

func (s *eventService) ShowUserEvent(ctx context.Context, userId uuid.UUID) ([]dto.EventResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	users, err := s.eventRepo.FetchUserEvent(ctx, userId)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return []dto.EventResponse{}, nil
	}

	user := users[0]

	var responses []dto.EventResponse
	for _, evt := range user.Events {

		evt.User = user
		responses = append(responses, *dto.EventEntityToResponse(&evt))
	}

	return responses, nil
}

func (s *eventService) UploadLogo(ctx context.Context, Id string, UserId string, logoFile *multipart.FileHeader) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	eventId, err := uuid.Parse(Id)
	if err != nil {
		return domain.ErrInternalServer // Invalid ID
	}

	uId, err := uuid.Parse(UserId)
	if err != nil {
		return domain.ErrInternalServer
	}

	event, err := s.eventRepo.FetchOneById(ctx, eventId)
	if err != nil {
		return err
	}

	if event.UserId != uId {
		return domain.ErrForbiddenUpdate // Or unauthorized
	}

	filename := fmt.Sprintf("%s-%s-%s", Id, "logo", logoFile.Filename)
	path := filepath.Join("public", "uploads", "events", filename)

	if err := s.saveFile(logoFile, path); err != nil {
		return domain.ErrInternalServer
	}

	update := dto.EventUpdate{
		LogoPath: path,
	}

	return s.eventRepo.UpdateEvent(ctx, &update, eventId)
}

func (s *eventService) UploadPoster(ctx context.Context, Id string, UserId string, posterFile *multipart.FileHeader) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	eventId, err := uuid.Parse(Id)
	if err != nil {
		return domain.ErrInternalServer
	}

	uId, err := uuid.Parse(UserId)
	if err != nil {
		return domain.ErrInternalServer
	}

	event, err := s.eventRepo.FetchOneById(ctx, eventId)
	if err != nil {
		return err
	}

	if event.UserId != uId {
		return domain.ErrForbiddenUpdate
	}

	filename := fmt.Sprintf("%s-%s-%s", Id, "poster", posterFile.Filename)
	path := filepath.Join("public", "uploads", "events", filename)

	if err := s.saveFile(posterFile, path); err != nil {
		return domain.ErrInternalServer
	}

	update := dto.EventUpdate{
		PosterPath: path,
	}

	return s.eventRepo.UpdateEvent(ctx, &update, eventId)
}

func (s *eventService) UpdateEvent(ctx context.Context, Id string, UserId string, event *dto.EventUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	eventId, err := uuid.Parse(Id)
	if err != nil {
		return domain.ErrInternalServer
	}

	uId, err := uuid.Parse(UserId)
	if err != nil {
		return domain.ErrInternalServer
	}

	existingEvent, err := s.eventRepo.FetchOneById(ctx, eventId)
	if err != nil {
		return err
	}

	if existingEvent.UserId != uId {
		return domain.ErrForbiddenUpdate
	}

	return s.eventRepo.UpdateEvent(ctx, event, eventId)
}

func (s *eventService) RemoveUserEvent(ctx context.Context, User *entity.User) error {
	// Implementation ambiguous in contract. Returning nil.
	return nil
}

func (s *eventService) AddUserEvent(ctx context.Context, User *entity.User) error {
	// Implementation ambiguous in contract. Returning nil.
	return nil
}

func (s *eventService) CreateEventLevel(ctx context.Context, EventId string, UserId string, EventLevel *dto.EventLevelCreate) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	eId, err := uuid.Parse(EventId)
	if err != nil {
		return domain.ErrInternalServer
	}

	uId, err := uuid.Parse(UserId)
	if err != nil {
		return domain.ErrInternalServer
	}

	event, err := s.eventRepo.FetchOneById(ctx, eId)
	if err != nil {
		return err
	}

	if event.UserId != uId {
		return domain.ErrForbiddenUpdate
	}

	id, err := s.uuid.New()
	if err != nil {
		return err
	}

	newLevel := entity.EventLevel{
		Id:       id,
		EventId:  eId,
		Name:     EventLevel.Name,
		RegisFee: EventLevel.RegisFee,
		DpFee:    EventLevel.DpFee,
	}

	return s.eventRepo.InsertEventLevel(ctx, &newLevel)
}

func (s *eventService) UpdateEventLevel(ctx context.Context, EventId string, UserId string, EventLevel *dto.EventLevelUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	eId, err := uuid.Parse(EventId)
	if err != nil {
		return domain.ErrInternalServer
	}

	uId, err := uuid.Parse(UserId)
	if err != nil {
		return domain.ErrInternalServer
	}

	event, err := s.eventRepo.FetchOneById(ctx, eId)
	if err != nil {
		return err
	}

	if event.UserId != uId {
		return domain.ErrForbiddenUpdate
	}

	id := EventLevel.Id
	// Ideally we verify EventId belongs to user or similar, but contract is simple.

	// Try to find the event level or just update using ID
	// Since we are updating, we constructed the entity.
	// EventLevelUpdate contains Id.

	updatedLevel := entity.EventLevel{
		Id:       id,
		EventId:  EventLevel.EventId,
		Name:     EventLevel.Name,
		RegisFee: EventLevel.RegisFee,
		DpFee:    EventLevel.DpFee,
	}

	return s.eventRepo.UpdateEventLevel(ctx, &updatedLevel)
}

func (s *eventService) DeleteEventLevel(ctx context.Context, EventId string, EventLevelId string, UserId string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	eId, err := uuid.Parse(EventId)
	if err != nil {
		return domain.ErrInternalServer
	}

	uId, err := uuid.Parse(UserId)
	if err != nil {
		return domain.ErrInternalServer
	}

	event, err := s.eventRepo.FetchOneById(ctx, eId)
	if err != nil {
		return err
	}

	if event.UserId != uId {
		return domain.ErrForbiddenUpdate
	}

	id, err := uuid.Parse(EventLevelId)
	if err != nil {
		return domain.ErrInternalServer
	}

	return s.eventRepo.DeleteEventLevel(ctx, id)
}

func (s *eventService) DeleteEvent(ctx context.Context, EventId string, UserId string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	id, err := uuid.Parse(EventId)
	if err != nil {
		return domain.ErrInternalServer
	}

	uId, err := uuid.Parse(UserId)
	if err != nil {
		return domain.ErrInternalServer
	}

	event, err := s.eventRepo.FetchOneById(ctx, id)
	if err != nil {
		return err
	}

	if event.UserId != uId {
		return domain.ErrForbiddenUpdate
	}

	return s.eventRepo.DeleteEvent(ctx, id)
}

// Helper to save file
func (s *eventService) saveFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
