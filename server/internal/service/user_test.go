package service_test

import (
	"testing"

	"github.com/Jarmos-san/arthika/server/internal/service"
)

func TestNewUserService(t *testing.T) {
	t.Parallel()

	svc := service.NewUserService(nil, "test-secret")

	if svc == nil {
		t.Fatal("expected non-nil UserService, got nil")
	}
}

func TestUserService_GetUser(t *testing.T) {
	t.Parallel()

	svc := service.NewUserService(nil, "test-secret")

	user := svc.GetUser()

	if user.Name != "John Doe" {
		t.Errorf("expected Name to be 'John Doe', got '%s'", user.Name)
	}
}
