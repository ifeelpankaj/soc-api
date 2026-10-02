package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	// "go-server/internal/middleware"
	"go-server/internal/models"
	// "go-server/internal/service"
	"go-server/pkg/utils"
	// "strconv"

	"github.com/gin-gonic/gin"
)

// bindJSON deserializes the request body into dst.
// On failure it writes the correct 400 response and returns false so the
// caller can do a bare `return`.
//
// Handles two common cases:
//   - json.UnmarshalTypeError  → specific field + type message
//   - everything else          → generic "Invalid request payload"
func bindJSON(c *gin.Context, dst interface{}) bool {
	if err := c.ShouldBindJSON(dst); err != nil {
		var unmarshalErr *json.UnmarshalTypeError
		if errors.As(err, &unmarshalErr) {
			utils.BadRequestResponse(c,
				fmt.Sprintf("%s must be of type %s",
					unmarshalErr.Field,
					unmarshalErr.Type.String(),
				),
			)
		} else {
			message := "Invalid request payload"
			if strings.Contains(err.Error(), "date_of_birth") {
				message = err.Error()
			}
			utils.BadRequestResponse(c, message)
		}
		return false
	}
	return true
}

// handleServiceError writes the appropriate HTTP error response for a
// service-layer error and returns true so the caller can do a bare `return`.
// Returns false when err is nil (no error to handle).
//
// Error mapping:
//   - *models.AppError  → uses its own StatusCode / Code / Message
//   - anything else     → 500 Internal Server Error (avoids leaking internals)
func handleServiceError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, models.ErrVisitorPhotoUploadRequired) {
		err = models.ErrVisitorPhotoUploadRequired
	}
	if errors.Is(err, models.ErrAvatarUploadRequired) {
		err = models.ErrAvatarUploadRequired
	}
	var appErr *models.AppError
	if errors.As(err, &appErr) {
		utils.ErrorResponse(c, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Internal)
	} else {
		utils.InternalServerErrorResponse(c, err)
	}
	return true
}

// func actorScopeFromContext(c *gin.Context) (service.ActorScope, bool) {
// 	userID, ok := middleware.GetUserIDFromContext(c)
// 	if !ok || userID <= 0 {
// 		utils.UnauthorizedResponse(c, "Authentication required")
// 		return service.ActorScope{}, false
// 	}
// 	role, ok := middleware.GetUserRoleFromContext(c)
// 	if !ok || role == "" {
// 		utils.UnauthorizedResponse(c, "Authentication required")
// 		return service.ActorScope{}, false
// 	}
// 	societyID, _ := middleware.GetUserSocietyIDFromContext(c)
// 	return service.ActorScope{
// 		UserID:    userID,
// 		Role:      role,
// 		SocietyID: societyID,
// 	}, true
// }

// func paginationFromQuery(c *gin.Context) (models.PaginationRequest, bool) {
// 	page := models.DefaultPage
// 	if raw := c.Query("page"); raw != "" {
// 		value, err := strconv.Atoi(raw)
// 		if err != nil || value <= 0 {
// 			utils.BadRequestResponse(c, "page must be a positive integer")
// 			return models.PaginationRequest{}, false
// 		}
// 		page = value
// 	}

// 	pageSize := models.DefaultPageSize
// 	if raw := c.Query("page_size"); raw != "" {
// 		value, err := strconv.Atoi(raw)
// 		if err != nil || value <= 0 {
// 			utils.BadRequestResponse(c, "page_size must be a positive integer")
// 			return models.PaginationRequest{}, false
// 		}
// 		pageSize = value
// 	}

// 	return models.NewPaginationRequest(page, pageSize), true
// }
