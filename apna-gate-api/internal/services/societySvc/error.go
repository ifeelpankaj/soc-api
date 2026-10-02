package societysvc

import (
	"net/http"

	"go-server/internal/models"
)

var (
	ErrSocietyNotFound       = models.NewAppError("SOCIETY_NOT_FOUND", "society not found", http.StatusNotFound, nil)
	ErrSocietyConflict       = models.NewAppError("SOCIETY_CONFLICT", "society already exists", http.StatusConflict, nil)
	ErrInvalidTransition     = models.NewAppError("INVALID_SOCIETY_TRANSITION", "invalid society lifecycle transition", http.StatusConflict, nil)
	ErrMemberNotFound        = models.NewAppError("SOCIETY_MEMBER_NOT_FOUND", "society member not found", http.StatusNotFound, nil)
	ErrMemberConflict        = models.NewAppError("SOCIETY_MEMBER_CONFLICT", "society member already exists", http.StatusConflict, nil)
	ErrInvalidMemberRole     = models.NewAppError("INVALID_MEMBER_ROLE", "invalid society member role", http.StatusBadRequest, nil)
	ErrInvalidMemberStatus   = models.NewAppError("INVALID_MEMBER_STATUS", "invalid society member status", http.StatusBadRequest, nil)
	ErrOwnerProtection       = models.NewAppError("OWNER_PROTECTION", "operation would leave society without an active owner", http.StatusConflict, nil)
	ErrForbiddenSociety      = models.NewAppError("SOCIETY_FORBIDDEN", "you do not have permission for this society", http.StatusForbidden, nil)
	ErrInvalidSocietyRequest = models.NewAppError("INVALID_SOCIETY_REQUEST", "invalid society request", http.StatusBadRequest, nil)
	ErrSocietyInactive       = models.NewAppError("SOCIETY_INACTIVE", "society is not active", http.StatusConflict, nil)
	ErrMemberInactive        = models.NewAppError("SOCIETY_MEMBER_INACTIVE", "society member is not active", http.StatusConflict, nil)
	ErrDuplicateGuardEmail   = models.NewAppError("DUPLICATE_GUARD_EMAIL", "a user with this email already exists", http.StatusConflict, nil)
	ErrDuplicateGuardPhone   = models.NewAppError("DUPLICATE_GUARD_PHONE", "a user with this phone number already exists", http.StatusConflict, nil)
)
