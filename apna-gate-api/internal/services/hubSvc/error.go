package hubsvc

import "go-server/internal/models"

var (
	ErrForbidden = models.NewAppError("HUB_FORBIDDEN", "Society Hub access is not permitted", 403, nil)
	ErrNotFound  = models.NewAppError("HUB_NOT_FOUND", "Society Hub resource not found", 404, nil)
	ErrInvalid   = models.NewAppError("HUB_INVALID", "Invalid Society Hub request", 400, nil)
	ErrConflict  = models.NewAppError("HUB_CONFLICT", "Content is locked, deleted, or attachment is unavailable", 409, nil)
)

type RateError struct{ RetryAfter int }

func (e *RateError) Error() string { return "Society Hub posting limit reached" }
