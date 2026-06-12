package service_test

import (
	"testing"

	"github.com/Jarmos-san/arthika/server/internal/service"
)

// TestNewUserService verifies that the constructor returns a non-nil service.
//
// This ensures that the service is properly initialized and ready for use.
func TestNewUserService(t *testing.T) {
	t.Parallel()

	svc := service.NewUserService(nil)

	if svc == nil {
		t.Fatal("expected non-nil UserService, got nil")
	}
}
