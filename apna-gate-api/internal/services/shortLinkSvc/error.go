package shortlinksvc

import (
	"go-server/internal/models"
	"net/http"
)

var (
	ErrShortLinkNotFound    = models.NewAppError("SHORT_LINK_NOT_FOUND", "short link not found", http.StatusNotFound, nil)
	ErrShortLinkUnavailable = models.NewAppError("SHORT_LINK_UNAVAILABLE", "short link is expired or revoked", http.StatusConflict, nil)
	ErrInvalidShortLink     = models.NewAppError("INVALID_SHORT_LINK", "invalid short link", http.StatusBadRequest, nil)
)
