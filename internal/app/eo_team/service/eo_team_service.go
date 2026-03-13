package service

import (
	"context"
	"errors"

	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/google/uuid"
)

var (
	ErrUnauthorized = errors.New("unauthorized: you do not own this event")
	ErrNotFound     = errors.New("registration not found for this event")
)

type eoTeamService struct {
	repo contracts.IEOTeamRepository
}

func NewEOTeamService(repo contracts.IEOTeamRepository) contracts.IEOTeamService {
	return &eoTeamService{
		repo: repo,
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

func (s *eoTeamService) GetListTeam(ctx context.Context, eventId, userId uuid.UUID, eventLevelId *uuid.UUID) ([]dto.EOTeamListRes, error) {
	if err := s.checkOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}

	regs, err := s.repo.FindAllRegistrationsByEvent(ctx, eventId, eventLevelId)
	if err != nil {
		return nil, err
	}

	var res []dto.EOTeamListRes
	for _, r := range regs {
		res = append(res, dto.EOTeamListRes{
			RegistrationId: r.Id,
			TeamId:         r.TeamId,
			LogoPath:       r.Team.LogoPath,
			TeamName:       r.Team.Name,
			Institution:    r.Team.Institution.Name,
			EventLevel:     r.EventLevel.Name,
			PaymentStatus:  r.PaymentStatus,
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
			IdCardPath: m.IdCardPath,
		})
	}
	if members == nil {
		members = make([]dto.EOTeamMemberRes, 0)
	}

	return &dto.EOTeamDetailRes{
		RegistrationId:   reg.Id,
		TeamId:           reg.TeamId,
		TeamName:         reg.Team.Name,
		LogoPath:         reg.Team.LogoPath,
		Pelatih:          reg.Team.Pelatih,
		RecLetterPath:    reg.Team.RecLetterPath,
		Institution:      reg.Team.Institution.Name,
		EventLevel:       reg.EventLevel.Name,
		PaymentStatus:    reg.PaymentStatus,
		PaymentProofPath: reg.PaymentProofPath,
		RejectionReason:  reg.RejectionReason,
		IsKick:           reg.IsKick,
		Members:          members,
	}, nil
}

func (s *eoTeamService) ApproveTeam(ctx context.Context, eventId, userId, registrationId uuid.UUID, req dto.EOTeamApproveReq) error {
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

	reg.PaymentStatus = req.PaymentStatus // e.g., FULL_PAID or DP_PAID
	return s.repo.UpdateRegistration(ctx, reg)
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

	reg.IsKick = true
	return s.repo.UpdateRegistration(ctx, reg)
}
