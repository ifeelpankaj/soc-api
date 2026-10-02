package authsvc

import (
	"testing"

	"go-server/internal/models"
	"go-server/pkg/validator"
)

func TestUpdateProfileValidateStructGender(t *testing.T) {
	t.Parallel()

	gender := "invalid"
	req := models.UpdateUserRequest{
		Gender: &gender,
	}
	req.Sanitize()

	validationErrors := validator.ValidateStruct(&req)
	if len(validationErrors) == 0 {
		t.Fatalf("expected gender validation error")
	}
}
