package service

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/google/uuid"
)

type stubPrivateFileRepository struct {
	record contracts.PrivateFileRecord
}

func (r stubPrivateFileRepository) FindPrivateFile(ctx context.Context, resourceType string, resourceID uuid.UUID) (contracts.PrivateFileRecord, error) {
	return r.record, nil
}

func TestResolvePrivateFileEnforcesResourceOwnership(t *testing.T) {
	tempDir := t.TempDir()
	privatePath := filepath.Join(tempDir, "proof.png")
	if err := os.WriteFile(privatePath, []byte("proof"), 0600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	participantID := uuid.New()
	organizerID := uuid.New()
	foreignID := uuid.New()
	record := contracts.PrivateFileRecord{
		Path:              privatePath,
		ParticipantUserID: participantID,
		OrganizerUserID:   organizerID,
	}

	tests := []struct {
		name    string
		userID  uuid.UUID
		role    string
		wantErr error
	}{
		{name: "participant owner", userID: participantID, role: string(enums.Peserta)},
		{name: "organizer owner", userID: organizerID, role: string(enums.Organizer)},
		{name: "admin", userID: foreignID, role: string(enums.Admin)},
		{name: "foreign participant", userID: foreignID, role: string(enums.Peserta), wantErr: domain.ErrForbidden},
		{name: "foreign organizer", userID: foreignID, role: string(enums.Organizer), wantErr: domain.ErrForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewPrivateFileService(stubPrivateFileRepository{record: record})
			path, err := svc.ResolvePrivateFile(context.Background(), tt.userID, tt.role, "registration-payment", uuid.New())
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
			if tt.wantErr == nil && path != privatePath {
				t.Fatalf("expected %q, got %q", privatePath, path)
			}
		})
	}
}

func TestResolvePrivateFileRejectsUnknownResourceType(t *testing.T) {
	svc := NewPrivateFileService(stubPrivateFileRepository{})
	_, err := svc.ResolvePrivateFile(context.Background(), uuid.New(), string(enums.Admin), "arbitrary-path", uuid.New())
	if !errors.Is(err, domain.ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
}
