package seeders

import (
	"testing"

	"github.com/BangNopall/paskihub-be/domain/enums"
)

func TestBuildDemoUsers(t *testing.T) {
	users := buildDemoUsers("hashed-password")

	if len(users) != 3 {
		t.Fatalf("expected 3 demo users, got %d", len(users))
	}

	expectedRoles := map[string]enums.Role{
		"admin.demo@paskihub.local":     enums.Admin,
		"organizer.demo@paskihub.local": enums.Organizer,
		"peserta.demo@paskihub.local":   enums.Peserta,
	}

	for _, user := range users {
		role, ok := expectedRoles[user.Email]
		if !ok {
			t.Fatalf("unexpected demo user email %q", user.Email)
		}

		if user.Role != role {
			t.Fatalf("expected %s role %s, got %s", user.Email, role, user.Role)
		}

		if user.Password != "hashed-password" {
			t.Fatalf("expected %s to use supplied hashed password", user.Email)
		}

		if !user.EmailIsVerified {
			t.Fatalf("expected %s to be verified", user.Email)
		}

		if user.ParentId != nil {
			t.Fatalf("expected %s parent id to be nil", user.Email)
		}
	}
}
