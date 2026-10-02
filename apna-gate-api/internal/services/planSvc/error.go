package plansvc

import (
	"net/http"

	"go-server/internal/models"
)

var (
	ErrPlanNotFound       = models.NewAppError("PLAN_NOT_FOUND", "plan not found", http.StatusNotFound, nil)
	ErrPlanConflict       = models.NewAppError("PLAN_CONFLICT", "plan already exists", http.StatusConflict, nil)
	ErrInvalidPlanRequest = models.NewAppError("INVALID_PLAN_REQUEST", "invalid plan request", http.StatusBadRequest, nil)
)
