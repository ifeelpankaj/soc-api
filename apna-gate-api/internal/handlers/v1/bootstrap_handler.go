package handlers

import (
	"net/http"

	middleware "go-server/internal/middlewares"
	bootstrapsvc "go-server/internal/services/bootstrapSvc"
	"go-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

type BootstrapHandler struct {
	bootstrapSvc bootstrapsvc.BootstrapService
}

func NewBootstrapHandler(bootstrapSvc bootstrapsvc.BootstrapService) *BootstrapHandler {
	return &BootstrapHandler{
		bootstrapSvc: bootstrapSvc,
	}
}

// GetBootstrap godoc
// @Summary Bootstrap authenticated session
// @Description Returns the authenticated user, memberships, residences, and the default dashboard destination.
// @Tags Bootstrap
// @Produce json
// @Success 200 {object} models.BootstrapAPIResponse "Bootstrap fetched successfully"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Account disabled or blocked"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/bootstrap [get]
func (h *BootstrapHandler) GetBootstrap(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		utils.UnauthorizedResponse(c, "Authentication required")
		return
	}

	result, err := h.bootstrapSvc.GetBootstrap(c.Request.Context(), userID)
	if handleServiceError(c, err) {
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Bootstrap fetched successfully", result)
}
