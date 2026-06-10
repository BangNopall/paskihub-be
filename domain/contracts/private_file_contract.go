package contracts

import (
	"context"

	"github.com/google/uuid"
)

type PrivateFileRecord struct {
	Path              string
	ParticipantUserID uuid.UUID
	OrganizerUserID   uuid.UUID
}

type PrivateFileRepository interface {
	FindPrivateFile(ctx context.Context, resourceType string, resourceID uuid.UUID) (PrivateFileRecord, error)
}

type PrivateFileService interface {
	ResolvePrivateFile(ctx context.Context, userID uuid.UUID, role, resourceType string, resourceID uuid.UUID) (string, error)
}
