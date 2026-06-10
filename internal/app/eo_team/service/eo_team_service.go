package service

import (
	"context"
	"errors"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/BangNopall/paskihub-be/pkg/helpers"
	"github.com/google/uuid"
)

var (
	ErrUnauthorized        = errors.New("unauthorized: you do not own this event")
	ErrNotFound            = domain.ErrRegistrationNotFound
	ErrInsufficientBalance = domain.ErrInsufficientBalance
)

type eoTeamService struct {
	repo        contracts.IEOTeamRepository
	settingRepo contracts.SystemSettingRepository
}

func NewEOTeamService(
	repo contracts.IEOTeamRepository,
	settingRepo contracts.SystemSettingRepository,
) contracts.IEOTeamService {
	return &eoTeamService{
		repo:        repo,
		settingRepo: settingRepo,
	}
}

func (s *eoTeamService) checkOwnership(ctx context.Context, eventId, userId uuid.UUID) error {
	isOwner, err := s.repo.CheckEventOwnership(ctx, eventId, userId)
	if err != nil {
		return err
	}
	if !isOwner {
		return ErrUnauthorized
	}
	return nil
}

func (s *eoTeamService) GetListTeam(ctx context.Context, eventId, userId uuid.UUID, eventLevelId *uuid.UUID, institutionType *string) ([]dto.EOTeamListRes, error) {
	if err := s.checkOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}

	regs, err := s.repo.FindAllRegistrationsByEvent(ctx, eventId, eventLevelId, institutionType)
	if err != nil {
		return nil, err
	}

	var res []dto.EOTeamListRes
	for _, r := range regs {
		assessStatus, _ := s.repo.GetAssessmentStatus(ctx, r.Id)
		res = append(res, dto.EOTeamListRes{
			RegistrationId:   r.Id,
			TeamId:           r.TeamId,
			LogoPath:         r.Team.LogoPath,
			TeamName:         r.Team.Name,
			Institution:      r.Team.Institution.Name,
			InstitutionType:  string(r.Team.Institution.InstitutionType),
			EventLevel:       r.EventLevel.Name,
			PaymentStatus:    r.PaymentStatus,
			AssessmentStatus: assessStatus,
		})
	}
	if res == nil {
		res = make([]dto.EOTeamListRes, 0)
	}
	return res, nil
}

func (s *eoTeamService) GetDetailTeam(ctx context.Context, eventId, userId, registrationId uuid.UUID) (*dto.EOTeamDetailRes, error) {
	if err := s.checkOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}

	reg, err := s.repo.FindRegistrationByIdAndEvent(ctx, registrationId, eventId)
	if err != nil {
		return nil, err
	}
	if reg == nil {
		return nil, ErrNotFound
	}

	var members []dto.EOTeamMemberRes
	for _, m := range reg.Team.TeamMembers {
		members = append(members, dto.EOTeamMemberRes{
			Id:         m.Id,
			FullName:   m.FullName,
			Role:       m.Role,
			IdCardPath: helpers.PrivateFileURL("member-id-card", m.Id),
			PhotoPath:  helpers.PrivateFileURL("member-photo", m.Id),
		})
	}
	if members == nil {
		members = make([]dto.EOTeamMemberRes, 0)
	}

	return &dto.EOTeamDetailRes{
		RegistrationId:     reg.Id,
		TeamId:             reg.TeamId,
		TeamName:           reg.Team.Name,
		LogoPath:           reg.Team.LogoPath,
		Pelatih:            reg.Team.Pelatih,
		RecLetterPath:      helpers.PrivateFileURL("team-recommendation", reg.Team.Id),
		Institution:        reg.Team.Institution.Name,
		InstitutionAddress: reg.Team.Institution.Address,
		ContactEmail:       reg.Team.Institution.User.Email,
		EventLevel:         reg.EventLevel.Name,
		PaymentStatus:      reg.PaymentStatus,
		PaymentProofPath:   helpers.PrivateFileURL("registration-payment", reg.Id),
		RejectionReason:    reg.RejectionReason,
		IsKick:             reg.IsKick,
		Members:            members,
	}, nil
}

func (s *eoTeamService) ApproveTeam(ctx context.Context, eventId, userId, registrationId uuid.UUID, req dto.EOTeamApproveReq) error {
	if err := s.checkOwnership(ctx, eventId, userId); err != nil {
		return err
	}

	setting, err := s.settingRepo.Get(ctx)
	if err != nil {
		return err
	}

	totalFeeIDR := setting.ApprovalFee * setting.CoinRate

	return s.repo.ApproveRegistration(ctx, eventId, registrationId, totalFeeIDR, setting.ApprovalFee, setting.CoinRate, req.PaymentStatus)
}

func (s *eoTeamService) RejectTeam(ctx context.Context, eventId, userId, registrationId uuid.UUID, req dto.EOTeamRejectReq) error {
	if err := s.checkOwnership(ctx, eventId, userId); err != nil {
		return err
	}

	reg, err := s.repo.FindRegistrationByIdAndEvent(ctx, registrationId, eventId)
	if err != nil {
		return err
	}
	if reg == nil {
		return ErrNotFound
	}

	reg.PaymentStatus = enums.Rejected
	reg.RejectionReason = req.RejectionReason
	return s.repo.UpdateRegistration(ctx, reg)
}

func (s *eoTeamService) KickTeam(ctx context.Context, eventId, userId, registrationId uuid.UUID) error {
	if err := s.checkOwnership(ctx, eventId, userId); err != nil {
		return err
	}

	reg, err := s.repo.FindRegistrationByIdAndEvent(ctx, registrationId, eventId)
	if err != nil {
		return err
	}
	if reg == nil {
		return ErrNotFound
	}

	reg.PaymentStatus = enums.Kicked
	reg.IsKick = true
	return s.repo.UpdateRegistration(ctx, reg)
}

func (s *eoTeamService) GetStats(ctx context.Context, eventId, userId uuid.UUID) (*dto.EOTeamStatsRes, error) {
	if err := s.checkOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}

	return s.repo.GetStats(ctx, eventId)
}

func (s *eoTeamService) StartAssessment(ctx context.Context, eventId, userId, registrationId uuid.UUID) error {
	if err := s.checkOwnership(ctx, eventId, userId); err != nil {
		return err
	}

	reg, err := s.repo.FindRegistrationByIdAndEvent(ctx, registrationId, eventId)
	if err != nil {
		return err
	}
	if reg == nil {
		return ErrNotFound
	}

	// Only update to IN_PROGRESS if it's currently PENDING
	if reg.AssessmentStatus == enums.AssessPending {
		reg.AssessmentStatus = enums.AssessInProgress
		return s.repo.UpdateRegistration(ctx, reg)
	}

	return nil
}
