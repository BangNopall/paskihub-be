package service

import (
	"context"
	"os"
	"strings"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/google/uuid"
)

var privateFileTypes = map[string]struct{}{
	"team-recommendation":  {},
	"member-id-card":       {},
	"member-photo":         {},
	"registration-payment": {},
	"wallet-topup-proof":   {},
}

type privateFileService struct {
	repo contracts.PrivateFileRepository
}

func NewPrivateFileService(repo contracts.PrivateFileRepository) contracts.PrivateFileService {
	return &privateFileService{repo: repo}
}

func (s *privateFileService) ResolvePrivateFile(
	ctx context.Context,
	userID uuid.UUID,
	role string,
	resourceType string,
	resourceID uuid.UUID,
) (string, error) {
	if _, valid := privateFileTypes[resourceType]; !valid {
		return "", domain.ErrBadRequest
	}

	record, err := s.repo.FindPrivateFile(ctx, resourceType, resourceID)
	if err != nil {
		return "", err
	}

	authorized := role == string(enums.Admin) ||
		(role == string(enums.Peserta) && record.ParticipantUserID == userID) ||
		(role == string(enums.Organizer) && record.OrganizerUserID == userID)
	if !authorized {
		return "", domain.ErrForbidden
	}

	if record.Path == "" {
		return "", domain.ErrNotFound
	}
	path := record.Path
	if strings.HasPrefix(path, "/public/") {
		path = strings.TrimPrefix(path, "/")
	}
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return "", domain.ErrNotFound
		}
		return "", domain.ErrInternalServer
	}

	return path, nil
}
