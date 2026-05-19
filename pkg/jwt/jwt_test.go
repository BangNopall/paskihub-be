package jwt

import (
	"testing"
	"time"

	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/google/uuid"
)

func TestValidateTokenReturnsParentIdForStaffAccount(t *testing.T) {
	parentId := uuid.New()
	staffId := uuid.New()
	j := &JwtStruct{
		SecretKey:   "test-secret",
		ExpiredTime: time.Hour,
	}

	token, err := j.GenerateToken(staffId, entity.User{
		Id:       staffId,
		Email:    "staff@example.com",
		Role:     enums.Organizer,
		ParentId: &parentId,
	})
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	_, _, _, gotParentId, err := j.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}
	if gotParentId == nil {
		t.Fatal("ValidateToken() parentId = nil, want staff parent id")
	}
	if *gotParentId != parentId {
		t.Fatalf("ValidateToken() parentId = %s, want %s", *gotParentId, parentId)
	}
}
