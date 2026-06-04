package dto_test

import (
	"encoding/json"
	"testing"

	"github.com/Jarmos-san/arthika/server/internal/dto"
)

func TestCreateUser_JSON(t *testing.T) {
	t.Parallel()

	user := dto.CreateUser{
		ID:           "123",
		Username:     "jane",
		Email:        "jane@example.com",
		PasswordHash: "supersecret",
	}

	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if raw["id"] != "123" {
		t.Errorf("expected id '123', got %v", raw["id"])
	}
	if raw["username"] != "jane" {
		t.Errorf("expected username 'jane', got %v", raw["username"])
	}
	if raw["email"] != "jane@example.com" {
		t.Errorf("expected email 'jane@example.com', got %v", raw["email"])
	}
	if _, ok := raw["password_hash"]; ok {
		t.Error("password_hash should not appear in JSON output")
	}
}
