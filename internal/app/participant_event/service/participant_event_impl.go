package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/google/uuid"
)

type participantEventService struct {
	repo contracts.ParticipantEventRepository
}

func NewParticipantEventService(repo contracts.ParticipantEventRepository) contracts.ParticipantEventService {
	return &participantEventService{
		repo: repo,
	}
}

func saveFile(fileHeader *multipart.FileHeader, folderPath string) (string, error) {
	if fileHeader == nil {
		return "", nil
	}
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileHeader.Filename)
	fullPath := filepath.Join(folderPath, filename)

	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return "/" + filepath.ToSlash(fullPath), nil
}

func (s *participantEventService) GetOpenEvents(ctx context.Context) ([]dto.OpenEventResponse, error) {
	events, err := s.repo.GetOpenEvents(ctx)
	if err != nil {
		return nil, err
	}

	var res []dto.OpenEventResponse
	for _, ev := range events {
		oev := dto.OpenEventResponse{
			Id:          ev.Id.String(),
			Name:        ev.Name,
			Description: ev.Description,
			LogoPath:    ev.LogoPath,
			PosterPath:  ev.PosterPath,
		}

		for _, lvl := range ev.EventLevels {
			oev.Levels = append(oev.Levels, dto.OpenEventLevelResponse{
				Id:       lvl.Id.String(),
				Name:     lvl.Name,
				RegisFee: lvl.RegisFee,
				DpFee:    lvl.DpFee,
			})
		}
		res = append(res, oev)
	}

	return res, nil
}

func (s *participantEventService) RegisterEvent(ctx context.Context, req dto.RegisterEventRequest) error {
	levelID, err := uuid.Parse(req.EventLevelId)
	if err != nil {
		return errors.New("invalid event level id")
	}

	teamID, err := uuid.Parse(req.TeamId)
	if err != nil {
		return errors.New("invalid team id")
	}

	proofPath, err := saveFile(req.PaymentProof, "public/uploads/payments")
	if err != nil {
		return err
	}

	paymentStatus := enums.Waiting
	// RegistrationStatus enum typically handles WAITING, DP_PAID, FULL_PAID
	// Here we just set WAITING since EO needs to verify the proof.

	regis := &entity.Registration{
		Id:               uuid.New(),
		TeamId:           teamID,
		EventLevelId:     levelID,
		PaymentStatus:    paymentStatus,
		PaymentProofPath: proofPath,
	}

	return s.repo.CreateRegistration(ctx, regis)
}

func (s *participantEventService) PelunasanEvent(ctx context.Context, regisID string, req dto.PelunasanEventRequest) error {
	parsedRegisID, err := uuid.Parse(regisID)
	if err != nil {
		return errors.New("invalid regis id")
	}

	regis, err := s.repo.GetRegistrationByID(ctx, parsedRegisID)
	if err != nil {
		return err
	}

	if regis.PaymentStatus == enums.FullPaid {
		return errors.New("registration is already fully paid")
	}

	proofPath, err := saveFile(req.PaymentProof, "public/uploads/payments_pelunasan")
	if err != nil {
		return err
	}

	regis.PaymentProofPath = proofPath
	regis.PaymentStatus = enums.Waiting // EO will verify this final payment
	
	return s.repo.UpdateRegistration(ctx, regis)
}

func (s *participantEventService) GetActiveEvents(ctx context.Context, userID string) ([]dto.ActiveEventResponse, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	registrations, err := s.repo.GetActiveRegistrationsByUserID(ctx, parsedUserID)
	if err != nil {
		return nil, err
	}

	var res []dto.ActiveEventResponse
	for _, r := range registrations {
		res = append(res, dto.ActiveEventResponse{
			RegistrationId: r.Id.String(),
			EventName:      r.EventLevel.Event.Name,
			EventLogoPath:  r.EventLevel.Event.LogoPath,
			TeamName:       r.Team.Name,
			PaymentStatus:  string(r.PaymentStatus),
		})
	}

	return res, nil
}
