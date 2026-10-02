package visitorentrysvc

import (
	"net/http"

	"go-server/internal/models"
)

var (
	ErrVisitorEntryNotFound           = models.NewAppError("VISITOR_ENTRY_NOT_FOUND", "visitor entry not found", http.StatusNotFound, nil)
	ErrVisitorInviteNotFound          = models.NewAppError("VISITOR_INVITE_NOT_FOUND", "visitor invite not found", http.StatusNotFound, nil)
	ErrInvalidVisitorRequest          = models.NewAppError("INVALID_VISITOR_REQUEST", "invalid visitor request", http.StatusBadRequest, nil)
	ErrVisitorForbidden               = models.NewAppError("VISITOR_FORBIDDEN", "only an active flat resident or society admin can perform this action", http.StatusForbidden, nil)
	ErrVisitorInviteUnavailable       = models.NewAppError("VISITOR_INVITE_UNAVAILABLE", "visitor invite is not active or has expired", http.StatusConflict, nil)
	ErrVisitorInvalidState            = models.NewAppError("VISITOR_INVALID_STATE", "visitor entry is not in a valid state for this action", http.StatusConflict, nil)
	ErrVisitorAlreadyCheckedIn        = models.NewAppError("VISITOR_ALREADY_CHECKED_IN", "visitor is already checked in", http.StatusConflict, nil)
	ErrVisitorQRUnavailable           = models.NewAppError("VISITOR_QR_UNAVAILABLE", "visitor QR is unavailable", http.StatusConflict, nil)
	ErrVisitorQRInvalid               = models.NewAppError("VISITOR_QR_INVALID", "invalid visitor QR token", http.StatusBadRequest, nil)
	ErrVisitorQRExpired               = models.NewAppError("VISITOR_QR_EXPIRED", "visitor QR token has expired", http.StatusConflict, nil)
	ErrVisitorFlatNotFound            = models.NewAppError("VISITOR_FLAT_NOT_FOUND", "flat not found", http.StatusNotFound, nil)
	ErrVisitorSettingsNotFound        = models.NewAppError("VISITOR_SETTINGS_NOT_FOUND", "visitor settings not found", http.StatusNotFound, nil)
	ErrVisitorGuardOnBehalfNotAllowed = models.NewAppError("VISITOR_GUARD_ON_BEHALF_NOT_ALLOWED", "guard on-behalf approval is not allowed for this society", http.StatusForbidden, nil)
)
