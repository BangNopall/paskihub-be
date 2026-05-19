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
	repo      contracts.ParticipantEventRepository
	rekapRepo contracts.RekapRepository
}

func NewParticipantEventService(repo contracts.ParticipantEventRepository, rekapRepo contracts.RekapRepository) contracts.ParticipantEventService {
	return &participantEventService{
		repo:      repo,
		rekapRepo: rekapRepo,
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
			Id:             ev.Id.String(),
			Name:           ev.Name,
			Description:    ev.Description,
			LogoPath:       ev.LogoPath,
			PosterPath:     ev.PosterPath,
			Organizer:      ev.Organizer,
			Status:         string(ev.Status),
			OpenDate:       ev.OpenDate.Format("2006-01-02"),
			CloseDate:      ev.CloseDate.Format("2006-01-02"),
			Location:       ev.Location,
			MinTeamMembers: ev.MinTeamMembers,
			MaxTeamMembers: ev.MaxTeamMembers,
			BankName:       ev.BankName,
			BankNumber:     ev.BankNumber,
			NamePj:         ev.NamePj,
			NoWaPj:         ev.NoWaPj,
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

func (s *participantEventService) RegisterEvent(ctx context.Context, userID string, req dto.RegisterEventRequest) error {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	levelID, err := uuid.Parse(req.EventLevelId)
	if err != nil {
		return errors.New("invalid event level id")
	}

	teamID, err := uuid.Parse(req.TeamId)
	if err != nil {
		return errors.New("invalid team id")
	}

	// 1. Fetch event level with event
	level, err := s.repo.GetEventLevelByID(ctx, levelID)
	if err != nil {
		return errors.New("event level not found")
	}

	// 2. Fetch team with members and institution
	team, err := s.repo.GetTeamWithMembers(ctx, teamID)
	if err != nil {
		return errors.New("team not found")
	}

	// 3. Verify team belongs to the user's institution
	if team.Institution.UserId != parsedUserID {
		return errors.New("team does not belong to your institution")
	}

	// 4. Validate team members count
	memberCount := len(team.TeamMembers)
	if memberCount < level.Event.MinTeamMembers || memberCount > level.Event.MaxTeamMembers {
		return fmt.Errorf("team members count (%d) must be between %d and %d",
			memberCount, level.Event.MinTeamMembers, level.Event.MaxTeamMembers)
	}

	// 5. Check if already registered for this event
	exists, err := s.repo.CheckExistingRegistration(ctx, teamID, level.EventId)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("team is already registered for this event")
	}

	proofPath, err := saveFile(req.PaymentProof, "public/uploads/payments")
	if err != nil {
		return err
	}

	paymentStatus := enums.Waiting

	regis := &entity.Registration{
		Id:               uuid.New(),
		TeamId:           teamID,
		EventLevelId:     levelID,
		PaymentType:      req.PaymentType,
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
			RegistrationId:  r.Id.String(),
			EventName:       r.EventLevel.Event.Name,
			EventLogoPath:   r.EventLevel.Event.LogoPath,
			TeamName:        r.Team.Name,
			PaymentStatus:   string(r.PaymentStatus),
			PaymentType:     r.PaymentType,
			RejectionReason: r.RejectionReason,
			IsKick:          r.IsKick,
		})
	}

	return res, nil
}

func (s *participantEventService) GetRegistrationDetail(ctx context.Context, userID string, regisID string) (dto.RegistrationDetailResponse, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return dto.RegistrationDetailResponse{}, errors.New("invalid user id")
	}

	parsedID, err := uuid.Parse(regisID)
	if err != nil {
		return dto.RegistrationDetailResponse{}, errors.New("invalid registration id")
	}

	isOwner, err := s.repo.GetRegistrationOwnership(ctx, parsedID, parsedUserID)
	if err != nil {
		return dto.RegistrationDetailResponse{}, err
	}
	if !isOwner {
		return dto.RegistrationDetailResponse{}, errors.New("unauthorized to view this registration")
	}

	regis, err := s.repo.GetRegistrationWithDetails(ctx, parsedID)
	if err != nil {
		return dto.RegistrationDetailResponse{}, err
	}

	res := dto.RegistrationDetailResponse{}
	res.EventLevelId = regis.EventLevelId.String()
	res.Event.Id = regis.EventLevel.Event.Id.String()
	res.Event.Title = regis.EventLevel.Event.Name
	res.Event.Description = regis.EventLevel.Event.Description
	res.Event.Date = regis.EventLevel.Event.CompeDate.Format("02 January 2006")
	res.Event.Location = regis.EventLevel.Event.Location
	res.Event.Price = regis.EventLevel.RegisFee
	res.Event.TargetDate = regis.EventLevel.Event.CompeDate.Format(time.RFC3339)
	res.Event.LogoUrl = regis.EventLevel.Event.LogoPath

	res.Team.Id = regis.Team.Id.String()
	res.Team.Name = regis.Team.Name
	res.Team.LogoUrl = regis.Team.LogoPath

	officialCount := 0
	pasukanCount := 0
	for _, m := range regis.Team.TeamMembers {
		if m.Role == enums.Official {
			officialCount++
		} else if m.Role == enums.Pasukan {
			pasukanCount++
		}
	}
	res.Team.OfficialCount = officialCount
	res.Team.PasukanCount = pasukanCount

	res.Payment.Status = string(regis.PaymentStatus)
	res.Payment.ProofUrl = regis.PaymentProofPath
	res.Payment.TotalAmount = regis.EventLevel.RegisFee

	// Logic for amount paid and remaining amount
	if regis.PaymentStatus == enums.FullPaid {
		res.Payment.AmountPaid = regis.EventLevel.RegisFee
		res.Payment.RemainingAmount = "0"
	} else if regis.PaymentStatus == enums.DPPaid {
		res.Payment.AmountPaid = regis.EventLevel.DpFee
		// Simplified placeholder as RegisFee is string.
		res.Payment.RemainingAmount = "Pending Calculation"
	} else {
		res.Payment.AmountPaid = "0"
		res.Payment.RemainingAmount = regis.EventLevel.RegisFee
	}

	return res, nil
}

func (s *participantEventService) GetScoreboard(ctx context.Context, eventLevelID string) (dto.ScoreboardResponse, error) {
	parsedID, err := uuid.Parse(eventLevelID)
	if err != nil {
		return dto.ScoreboardResponse{}, errors.New("invalid event level id")
	}

	level, err := s.repo.GetEventLevelByID(ctx, parsedID)
	if err != nil {
		return dto.ScoreboardResponse{}, err
	}

	if !level.IsScorePublished {
		return dto.ScoreboardResponse{}, errors.New("scoreboard for this level is not published yet")
	}

	items, err := s.rekapRepo.GetScoreboardByEventLevel(ctx, parsedID)
	if err != nil {
		return dto.ScoreboardResponse{}, err
	}

	return dto.ScoreboardResponse{
		EventLevelID: parsedID,
		Items:        items,
	}, nil
}
