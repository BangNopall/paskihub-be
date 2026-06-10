package helpers

import (
	"fmt"

	"github.com/google/uuid"
)

func PrivateFileURL(resourceType string, resourceID uuid.UUID) string {
	return fmt.Sprintf("/api/v1/files/%s/%s", resourceType, resourceID.String())
}
