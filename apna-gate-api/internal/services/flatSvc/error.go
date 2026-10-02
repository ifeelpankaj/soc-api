package flatsvc

import (
	"net/http"

	"go-server/internal/models"
)

var (
	ErrFlatNotFound               = models.NewAppError("FLAT_NOT_FOUND", "flat not found", http.StatusNotFound, nil)
	ErrFlatConflict               = models.NewAppError("FLAT_CONFLICT", "flat already exists", http.StatusConflict, nil)
	ErrFlatInactive               = models.NewAppError("FLAT_INACTIVE", "flat is inactive", http.StatusConflict, nil)
	ErrFlatBlocked                = models.NewAppError("FLAT_BLOCKED", "flat is blocked", http.StatusConflict, nil)
	ErrFlatOccupied               = models.NewAppError("FLAT_OCCUPIED", "flat is already occupied", http.StatusConflict, nil)
	ErrInvalidFlatRequest         = models.NewAppError("INVALID_FLAT_REQUEST", "invalid flat request", http.StatusBadRequest, nil)
	ErrInvalidFlatStatus          = models.NewAppError("INVALID_FLAT_STATUS", "invalid flat status", http.StatusBadRequest, nil)
	ErrClaimNotFound              = models.NewAppError("FLAT_CLAIM_NOT_FOUND", "flat claim not found", http.StatusNotFound, nil)
	ErrClaimConflict              = models.NewAppError("FLAT_CLAIM_CONFLICT", "flat claim conflict", http.StatusConflict, nil)
	ErrInvalidClaimRequest        = models.NewAppError("INVALID_FLAT_CLAIM_REQUEST", "invalid flat claim request", http.StatusBadRequest, nil)
	ErrResidentNotFound           = models.NewAppError("FLAT_RESIDENT_NOT_FOUND", "flat resident not found", http.StatusNotFound, nil)
	ErrResidentConflict           = models.NewAppError("FLAT_RESIDENT_CONFLICT", "flat resident conflict", http.StatusConflict, nil)
	ErrPrimaryResidentExists      = models.NewAppError("PRIMARY_RESIDENT_EXISTS", "another active primary resident already exists", http.StatusConflict, nil)
	ErrInvalidResidentRole        = models.NewAppError("INVALID_FLAT_RESIDENT_ROLE", "invalid flat resident role", http.StatusBadRequest, nil)
	ErrMemberInviteNotFound       = models.NewAppError("FLAT_MEMBER_INVITE_NOT_FOUND", "flat member invite not found", http.StatusNotFound, nil)
	ErrMemberInviteUnavailable    = models.NewAppError("FLAT_MEMBER_INVITE_UNAVAILABLE", "flat member invite is expired, accepted, or cancelled", http.StatusConflict, nil)
	ErrMemberInviteForbidden      = models.NewAppError("FLAT_MEMBER_INVITE_FORBIDDEN", "you do not have permission to manage flat members", http.StatusForbidden, nil)
	ErrInvalidMemberInviteRequest = models.NewAppError("INVALID_FLAT_MEMBER_INVITE_REQUEST", "invalid flat member invite request", http.StatusBadRequest, nil)
)
