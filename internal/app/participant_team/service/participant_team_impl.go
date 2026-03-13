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

type participantTeamService struct {
	repo        contracts.ParticipantTeamRepository
	profileRepo contracts.ParticipantProfileRepository
}

func NewParticipantTeamService(repo contracts.ParticipantTeamRepository, profileRepo contracts.ParticipantProfileRepository) contracts.ParticipantTeamService {
	return &participantTeamService{
		repo:        repo,
		profileRepo: profileRepo,
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

func (s *participantTeamService) CreateTeam(ctx context.Context, userID string, req dto.CreateTeamRequest) error {
	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	institution, err := s.profileRepo.GetInstitutionByUserID(ctx, parsedUUID)
	if err != nil {
		return errors.New("institution not found, please complete profile first")
	}

	teamID := uuid.New()

	logoPath, err := saveFile(req.LogoTeam, "public/uploads/teams/logos")
	if err != nil {
		return err
	}

	recLetterPath, err := saveFile(req.SuratRekomendasi, "public/uploads/teams/rekomendasi")
	if err != nil {
		return err
	}

	team := &entity.Team{
		Id:            teamID,
		InstiId:       institution.Id,
		Name:          req.Name,
		Pelatih:       req.PelatihName,
		LogoPath:      logoPath,
		RecLetterPath: recLetterPath,
	}

	var members []entity.TeamMember
	for _, m := range req.Members {
		idCardPath, err := saveFile(m.IdCard, "public/uploads/teams/id_cards")
		if err != nil {
			return err
		}
		photoPath, err := saveFile(m.Photo, "public/uploads/teams/photos")
		if err != nil {
			return err
		}

		members = append(members, entity.TeamMember{
			Id:         uuid.New(),
			FullName:   m.FullName,
			Role:       enums.TeamType(m.Role),
			IdCardPath: idCardPath,
			PhotoPath:  photoPath,
		})
	}

	return s.repo.CreateTeamWithMembers(ctx, team, members)
}

func (s *participantTeamService) UpdateTeam(ctx context.Context, userID string, teamID string, req dto.CreateTeamRequest) error {
	// Optional feature to implement for PUT /teams/:id
	return errors.New("not implemented yet")
}

func (s *participantTeamService) GetTeams(ctx context.Context, userID string) ([]dto.ParticipantTeamResponse, error) {
	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	institution, err := s.profileRepo.GetInstitutionByUserID(ctx, parsedUUID)
	if err != nil {
		return nil, errors.New("institution not found")
	}

	teamsEntity, err := s.repo.GetTeamsByInstitutionID(ctx, institution.Id)
	if err != nil {
		return nil, err
	}

	var responses []dto.ParticipantTeamResponse
	for _, t := range teamsEntity {
		paymentStatus := "NONE"
		if len(t.Registrations) > 0 {
			paymentStatus = string(t.Registrations[0].PaymentStatus)
		}

		responses = append(responses, dto.ParticipantTeamResponse{
			Id:            t.Id.String(),
			Name:          t.Name,
			LogoPath:      t.LogoPath,
			PaymentStatus: paymentStatus,
		})
	}

	return responses, nil
}

func (s *participantTeamService) GetTeamDetail(ctx context.Context, userID string, teamID string) (*dto.TeamDetailResponse, error) {
	parsedTeamID, err := uuid.Parse(teamID)
	if err != nil {
		return nil, errors.New("invalid team id")
	}

	team, err := s.repo.GetTeamByID(ctx, parsedTeamID)
	if err != nil {
		return nil, err
	}

	groupedMembers := make(map[string][]dto.ParticipantTeamMemberResponse)
	for _, m := range team.TeamMembers {
		roleStr := string(m.Role)
		groupedMembers[roleStr] = append(groupedMembers[roleStr], dto.ParticipantTeamMemberResponse{
			Id:         m.Id.String(),
			FullName:   m.FullName,
			Role:       roleStr,
			IdCardPath: m.IdCardPath,
			PhotoPath:  m.PhotoPath,
		})
	}

	response := &dto.TeamDetailResponse{
		Id:             team.Id.String(),
		Name:           team.Name,
		LogoPath:       team.LogoPath,
		Pelatih:        team.Pelatih,
		RecLetterPath:  team.RecLetterPath,
		MembersGrouped: groupedMembers,
	}

	return response, nil
}

func (s *participantTeamService) DeleteTeam(ctx context.Context, userID string, teamID string) error {
	parsedTeamID, err := uuid.Parse(teamID)
	if err != nil {
		return errors.New("invalid team id")
	}

	team, err := s.repo.GetTeamByID(ctx, parsedTeamID)
	if err != nil {
		return err
	}

	if len(team.Registrations) > 0 && (team.Registrations[0].PaymentStatus == enums.FullPaid) {
		return errors.New("cannot delete team, payment is already full paid or approved")
	}

	return s.repo.DeleteTeam(ctx, parsedTeamID)
}
