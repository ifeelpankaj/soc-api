package flatauthz

import (
	"net/http"

	"go-server/internal/models"
)

var (
	ErrForbidden = models.NewAppError(
		"FLAT_VISITOR_FORBIDDEN",
		"only an active flat resident or society admin can perform this action",
		http.StatusForbidden,
		nil,
	)
	ErrViewForbidden = models.NewAppError(
		"FLAT_VISITOR_VIEW_FORBIDDEN",
		"you do not have permission to view visitor data for this flat",
		http.StatusForbidden,
		nil,
	)
	ErrManageForbidden = models.NewAppError(
		"FLAT_MEMBER_MANAGE_FORBIDDEN",
		"only the flat owner or primary resident can manage flat members",
		http.StatusForbidden,
		nil,
	)
	ErrFlatNotFound = models.NewAppError(
		"FLAT_NOT_FOUND",
		"flat not found",
		http.StatusNotFound,
		nil,
	)
)
