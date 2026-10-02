package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	middleware "go-server/internal/middlewares"
	"go-server/internal/models"
	societysvc "go-server/internal/services/societySvc"
	"go-server/pkg/utils"
	"go-server/pkg/validator"

	"github.com/gin-gonic/gin"
)

type SocietyHandler struct {
	societySvc societysvc.SocietyService
}

func NewSocietyHandler(societySvc societysvc.SocietyService) *SocietyHandler {
	return &SocietyHandler{societySvc: societySvc}
}

// CreateSocietyRequest godoc
// @Summary Create society request
// @Description [User] Creates a pending society request and assigns the requester as owner.
// @Tags Societies
// @Accept json
// @Produce json
// @Param request body models.CreateSocietyRequest true "Create society payload"
// @Success 201 {object} models.SocietyAPIResponse "Society request created successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request or validation error"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 409 {object} models.ErrorResponseDoc "Duplicate society request"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies [post]
func (h *SocietyHandler) CreateSocietyRequest(c *gin.Context) {
	var req models.CreateSocietyRequest
	if !bindJSON(c, &req) {
		return
	}

	req.Sanitize()
	if err := req.Validate(); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	if validationErrors := validator.ValidateStruct(&req); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(c, validationErrors.ToMap())
		return
	}

	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	result, err := h.societySvc.CreateSocietyRequest(c.Request.Context(), req, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Society request created successfully", gin.H{"society": result})
}

// ApproveSociety godoc
// @Summary Approve society
// @Description [Developer] Approves a pending society request. Transition allowed: pending -> active.
// @Tags Societies
// @Produce json
// @Param societyId path int true "Society ID"
// @Success 200 {object} models.SocietyAPIResponse "Society approved successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid society ID"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 404 {object} models.ErrorResponseDoc "Society not found"
// @Failure 409 {object} models.ErrorResponseDoc "Invalid society lifecycle transition"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/approve [post]
func (h *SocietyHandler) ApproveSociety(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.societySvc.ApproveSociety(c.Request.Context(), societyID, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Society approved successfully", gin.H{"society": result})
}

// RejectSociety godoc
// @Summary Reject society
// @Description [Developer] Rejects a pending society request. Transition allowed: pending -> rejected.
// @Tags Societies
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param request body models.SocietyReasonRequest true "Rejection reason"
// @Success 200 {object} models.SocietyAPIResponse "Society rejected successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request or validation error"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 404 {object} models.ErrorResponseDoc "Society not found"
// @Failure 409 {object} models.ErrorResponseDoc "Invalid society lifecycle transition"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/reject [post]
func (h *SocietyHandler) RejectSociety(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	var req models.SocietyReasonRequest
	if !bindJSON(c, &req) {
		return
	}
	req.Sanitize()
	if validationErrors := validator.ValidateStruct(&req); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(c, validationErrors.ToMap())
		return
	}
	result, err := h.societySvc.RejectSociety(c.Request.Context(), societyID, userID, req.Reason)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Society rejected successfully", gin.H{"society": result})
}

// SuspendSociety godoc
// @Summary Suspend society
// @Description [Developer] Suspends an active society. Transition allowed: active -> suspended.
// @Tags Societies
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param request body models.SocietyReasonRequest true "Suspension reason"
// @Success 200 {object} models.SocietyAPIResponse "Society suspended successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request or validation error"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 404 {object} models.ErrorResponseDoc "Society not found"
// @Failure 409 {object} models.ErrorResponseDoc "Invalid society lifecycle transition"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/suspend [post]
func (h *SocietyHandler) SuspendSociety(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	var req models.SocietyReasonRequest
	if !bindJSON(c, &req) {
		return
	}
	req.Sanitize()
	if validationErrors := validator.ValidateStruct(&req); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(c, validationErrors.ToMap())
		return
	}
	result, err := h.societySvc.SuspendSociety(c.Request.Context(), societyID, userID, req.Reason)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Society suspended successfully", gin.H{"society": result})
}

// ReactivateSociety godoc
// @Summary Reactivate society
// @Description [Developer] Reactivates a suspended society. Transition allowed: suspended -> active.
// @Tags Societies
// @Produce json
// @Param societyId path int true "Society ID"
// @Success 200 {object} models.SocietyAPIResponse "Society reactivated successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid society ID"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 404 {object} models.ErrorResponseDoc "Society not found"
// @Failure 409 {object} models.ErrorResponseDoc "Invalid society lifecycle transition"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/reactivate [post]
func (h *SocietyHandler) ReactivateSociety(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.societySvc.ReactivateSociety(c.Request.Context(), societyID, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Society reactivated successfully", gin.H{"society": result})
}

// RestoreSociety godoc
// @Summary Restore soft-deleted society
// @Description [Developer] Restores a soft-deleted society to pending status.
// @Tags Societies
// @Produce json
// @Param societyId path int true "Society ID"
// @Success 200 {object} models.SocietyAPIResponse "Society restored successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid society ID"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 404 {object} models.ErrorResponseDoc "Society not found"
// @Failure 409 {object} models.ErrorResponseDoc "Society is not soft-deleted"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/restore [post]
func (h *SocietyHandler) RestoreSociety(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.societySvc.RestoreSociety(c.Request.Context(), societyID, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Society restored successfully", gin.H{"society": result})
}

// UpdateSociety godoc
// @Summary Update society
// @Description [User] Updates society profile fields. Requires active owner/admin membership.
// @Tags Societies
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param request body models.UpdateSocietyRequest true "Update society payload"
// @Success 200 {object} models.SocietyAPIResponse "Society updated successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request or validation error"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 404 {object} models.ErrorResponseDoc "Society not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId} [patch]
func (h *SocietyHandler) UpdateSociety(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	var req models.UpdateSocietyRequest
	if !bindJSON(c, &req) {
		return
	}
	req.Sanitize()
	if err := req.Validate(); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	if validationErrors := validator.ValidateStruct(&req); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(c, validationErrors.ToMap())
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.societySvc.UpdateSociety(c.Request.Context(), societyID, req, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Society updated successfully", gin.H{"society": result})
}

// DeleteSociety godoc
// @Summary Soft delete society
// @Description [User] Soft-deletes a society. Requires active owner membership.
// @Tags Societies
// @Produce json
// @Param societyId path int true "Society ID"
// @Success 200 {object} models.MessageAPIResponse "Society deleted successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid society ID"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 404 {object} models.ErrorResponseDoc "Society not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId} [delete]
func (h *SocietyHandler) DeleteSociety(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	if err := h.societySvc.DeleteSociety(c.Request.Context(), societyID, userID); handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Society deleted successfully", gin.H{"message": "Society deleted successfully"})
}

// GetSociety godoc
// @Summary Get society
// @Description [User/Developer] Fetches one society by ID. Developer routes may use the same handler with broader guards.
// @Tags Societies
// @Produce json
// @Param societyId path int true "Society ID"
// @Success 200 {object} models.SocietyDetailAPIResponse "Society fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid society ID"
// @Failure 404 {object} models.ErrorResponseDoc "Society not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/societies/{societyId} [get]
func (h *SocietyHandler) GetSociety(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	result, err := h.societySvc.GetSociety(c.Request.Context(), models.GetSocietyFilter{ID: &societyID})
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Society fetched successfully", gin.H{"society": result})
}

// GetPublicClaimOptions godoc
// @Summary Get public resident claim options
// @Description Public QR flow endpoint that resolves an active society code and returns safe society details plus active flats for resident claims.
// @Tags Public
// @Produce json
// @Param societyCode path string true "Society code"
// @Success 200 {object} models.PublicClaimOptionsAPIResponse "Claim options fetched successfully"
// @Failure 404 {object} models.ErrorResponseDoc "Society not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/public/societies/{societyCode}/claim-options [get]
func (h *SocietyHandler) GetPublicClaimOptions(c *gin.Context) {
	societyCode := strings.TrimSpace(c.Param("societyCode"))
	result, err := h.societySvc.GetPublicClaimOptions(c.Request.Context(), societyCode)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Claim options fetched successfully", gin.H{
		"society": result.Society,
		"flats":   result.Flats,
	})
}

// ListSocieties godoc
// @Summary List societies
// @Description [Developer] Lists societies with flexible filters for admin/developer panels.
// @Tags Societies
// @Produce json
// @Param status query string false "Society status"
// @Param search query string false "Search text"
// @Param search_mode query string false "Search mode"
// @Param name query string false "Society name"
// @Param code query string false "Society code"
// @Param city query string false "City"
// @Param state query string false "State"
// @Param country query string false "Country"
// @Param pincode query string false "Pincode"
// @Param created_by query int false "Created by user ID"
// @Param approved_by query int false "Approved by user ID"
// @Param rejected_by query int false "Rejected by user ID"
// @Param suspended_by query int false "Suspended by user ID"
// @Param created_from query string false "Created from RFC3339 timestamp"
// @Param created_to query string false "Created to RFC3339 timestamp"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Param sort_by query string false "Sort by: created_at, updated_at, name, city, status"
// @Param sort_order query string false "Sort order: asc, desc"
// @Success 200 {object} models.PaginatedSocietiesAPIResponse "Societies fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid query parameter"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/societies [get]
func (h *SocietyHandler) ListSocieties(c *gin.Context) {
	filter, ok := listSocietiesFilterFromQuery(c)
	if !ok {
		return
	}
	result, err := h.societySvc.ListSocieties(c.Request.Context(), filter)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Societies fetched successfully", gin.H{"societies": result})
}

// ListMySocieties godoc
// @Summary List my societies
// @Description [User] Lists societies where the authenticated user has a membership.
// @Tags Societies
// @Produce json
// @Success 200 {object} models.MySocietiesAPIResponse "My societies fetched successfully"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/my [get]
func (h *SocietyHandler) ListMySocieties(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.societySvc.ListMySocieties(c.Request.Context(), userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "My societies fetched successfully", gin.H{"societies": result})
}

// GetOnboardingBootstrap godoc
// @Summary Get society onboarding bootstrap
// @Description [Owner/Admin] Returns lifecycle, subscription, flat, and staff readiness before dashboard access.
// @Tags Societies
// @Produce json
// @Param societyId path int true "Society ID"
// @Success 200 {object} models.SocietyOnboardingBootstrapAPIResponse "Society onboarding bootstrap fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid society ID"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 404 {object} models.ErrorResponseDoc "Society not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/onboarding/bootstrap [get]
func (h *SocietyHandler) GetOnboardingBootstrap(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}

	result, err := h.societySvc.GetOnboardingBootstrap(c.Request.Context(), societyID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Society onboarding bootstrap fetched successfully", gin.H{"onboarding": result})
}

// GetDashboardBootstrap godoc
// @Summary Get society dashboard bootstrap
// @Description [Owner/Admin] Returns dashboard summary data for a society.
// @Tags Societies
// @Produce json
// @Param societyId path int true "Society ID"
// @Success 200 {object} models.SocietyDashboardBootstrapAPIResponse "Society dashboard bootstrap fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid society ID"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 404 {object} models.ErrorResponseDoc "Society not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/dashboard/bootstrap [get]
func (h *SocietyHandler) GetDashboardBootstrap(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}

	result, err := h.societySvc.GetDashboardBootstrap(c.Request.Context(), societyID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Society dashboard bootstrap fetched successfully", gin.H{"dashboard": result})
}

// GetDeveloperDashboardBootstrap godoc
// @Summary Get developer dashboard bootstrap
// @Description [Developer] Returns platform summary data for the developer dashboard.
// @Tags Developer
// @Produce json
// @Success 200 {object} models.DeveloperDashboardBootstrapAPIResponse "Developer dashboard bootstrap fetched successfully"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Developer access required"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/developer/dashboard/bootstrap [get]
func (h *SocietyHandler) GetDeveloperDashboardBootstrap(c *gin.Context) {
	result, err := h.societySvc.GetDeveloperDashboardBootstrap(c.Request.Context())
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Developer dashboard bootstrap fetched successfully", gin.H{"dashboard": result})
}

// CreateGuard godoc
// @Summary Create guard
// @Description [Owner/Admin] Creates a verified staff user for guard access without OTP or invite token.
// @Tags Society Members
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param request body models.CreateGuardRequest true "Create guard payload"
// @Success 201 {object} models.GuardAPIResponse "Guard created successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request or validation error"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 409 {object} models.ErrorResponseDoc "Duplicate email, phone, or member"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/guards [post]
func (h *SocietyHandler) CreateGuard(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	var req models.CreateGuardRequest
	if !bindJSON(c, &req) {
		return
	}
	req.Sanitize()
	if validationErrors := validator.ValidateStruct(&req); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(c, validationErrors.ToMap())
		return
	}
	actorID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.societySvc.CreateGuard(c.Request.Context(), societyID, actorID, req)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Guard created successfully", gin.H{"guard": result})
}

// AddMember godoc
// @Summary Add society member
// @Description [User] Adds a member to a society. Requires active owner/admin membership.
// @Tags Society Members
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param request body models.AddSocietyMemberRequest true "Add member payload"
// @Success 201 {object} models.SocietyMemberAPIResponse "Member added successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request or validation error"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 409 {object} models.ErrorResponseDoc "Duplicate member or owner protection"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/members [post]
func (h *SocietyHandler) AddMember(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	var req models.AddSocietyMemberRequest
	if !bindJSON(c, &req) {
		return
	}
	req.SocietyID = societyID
	if validationErrors := validator.ValidateStruct(&req); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(c, validationErrors.ToMap())
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.societySvc.AddMember(c.Request.Context(), req, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Member added successfully", gin.H{"member": result})
}

// ChangeMemberRole godoc
// @Summary Change member role
// @Description [User] Changes a member role. Requires active owner/admin membership and protects the last owner.
// @Tags Society Members
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param userId path int true "User ID"
// @Param request body models.ChangeSocietyMemberRoleRequest true "Change role payload"
// @Success 200 {object} models.SocietyMemberAPIResponse "Member role changed successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request or validation error"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 404 {object} models.ErrorResponseDoc "Member not found"
// @Failure 409 {object} models.ErrorResponseDoc "Owner protection conflict"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/members/{userId}/role [patch]
func (h *SocietyHandler) ChangeMemberRole(c *gin.Context) {
	societyID, userIDParam, ok := parseSocietyAndUser(c)
	if !ok {
		return
	}
	var req models.ChangeSocietyMemberRoleRequest
	if !bindJSON(c, &req) {
		return
	}
	req.SocietyID = societyID
	req.UserID = userIDParam
	if validationErrors := validator.ValidateStruct(&req); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(c, validationErrors.ToMap())
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.societySvc.ChangeMemberRole(c.Request.Context(), req, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Member role changed successfully", gin.H{"member": result})
}

// SuspendMember godoc
// @Summary Suspend member
// @Description [User] Suspends a member. Requires active owner/admin membership and protects the last owner.
// @Tags Society Members
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param userId path int true "User ID"
// @Param request body models.SocietyReasonRequest false "Suspension reason"
// @Success 200 {object} models.SocietyMemberAPIResponse "Member suspended successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request or validation error"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 404 {object} models.ErrorResponseDoc "Member not found"
// @Failure 409 {object} models.ErrorResponseDoc "Owner protection conflict"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/members/{userId}/suspend [post]
func (h *SocietyHandler) SuspendMember(c *gin.Context) {
	societyID, memberUserID, ok := parseSocietyAndUser(c)
	if !ok {
		return
	}
	req := models.SuspendSocietyMemberRequest{SocietyID: societyID, UserID: memberUserID}
	if c.Request.Body != nil && c.Request.ContentLength != 0 {
		var reason models.SocietyReasonRequest
		if !bindJSON(c, &reason) {
			return
		}
		reason.Sanitize()
		req.Reason = reason.Reason
	}
	actorID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.societySvc.SuspendMember(c.Request.Context(), req, actorID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Member suspended successfully", gin.H{"member": result})
}

// ReactivateMember godoc
// @Summary Reactivate member
// @Description [User] Reactivates a suspended or removed member. Requires active owner/admin membership.
// @Tags Society Members
// @Produce json
// @Param societyId path int true "Society ID"
// @Param userId path int true "User ID"
// @Success 200 {object} models.SocietyMemberAPIResponse "Member reactivated successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 404 {object} models.ErrorResponseDoc "Member not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/members/{userId}/reactivate [post]
func (h *SocietyHandler) ReactivateMember(c *gin.Context) {
	societyID, memberUserID, ok := parseSocietyAndUser(c)
	if !ok {
		return
	}
	actorID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.societySvc.ReactivateMember(c.Request.Context(), models.ReactivateSocietyMemberRequest{SocietyID: societyID, UserID: memberUserID}, actorID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Member reactivated successfully", gin.H{"member": result})
}

// RemoveMember godoc
// @Summary Remove member
// @Description [User] Removes a member from the society. Requires active owner/admin membership and protects the last owner.
// @Tags Society Members
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param userId path int true "User ID"
// @Param request body models.SocietyReasonRequest false "Removal reason"
// @Success 200 {object} models.MessageAPIResponse "Member removed successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request or validation error"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 404 {object} models.ErrorResponseDoc "Member not found"
// @Failure 409 {object} models.ErrorResponseDoc "Owner protection conflict"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/members/{userId} [delete]
func (h *SocietyHandler) RemoveMember(c *gin.Context) {
	societyID, memberUserID, ok := parseSocietyAndUser(c)
	if !ok {
		return
	}
	req := models.RemoveSocietyMemberRequest{SocietyID: societyID, UserID: memberUserID}
	if c.Request.Body != nil && c.Request.ContentLength != 0 {
		var reason models.SocietyReasonRequest
		if !bindJSON(c, &reason) {
			return
		}
		reason.Sanitize()
		req.Reason = reason.Reason
	}
	actorID, ok := currentUserID(c)
	if !ok {
		return
	}
	if err := h.societySvc.RemoveMember(c.Request.Context(), req, actorID); handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Member removed successfully", gin.H{"message": "Member removed successfully"})
}

// TransferOwnership godoc
// @Summary Transfer society ownership
// @Description [User] Transfers ownership to another user and leaves exactly one active owner.
// @Tags Society Members
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param request body models.TransferOwnershipRequest true "Transfer ownership payload"
// @Success 200 {object} models.SocietyMemberAPIResponse "Ownership transferred successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request or validation error"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 409 {object} models.ErrorResponseDoc "Owner protection conflict"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/transfer-ownership [post]
func (h *SocietyHandler) TransferOwnership(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	var req models.TransferOwnershipRequest
	if !bindJSON(c, &req) {
		return
	}
	if validationErrors := validator.ValidateStruct(&req); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(c, validationErrors.ToMap())
		return
	}
	actorID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.societySvc.TransferOwnership(c.Request.Context(), societyID, req.NewOwnerUserID, actorID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Ownership transferred successfully", gin.H{"member": result})
}

// GetSocietyMember godoc
// @Summary Get society member
// @Description [User] Fetches one society member by member ID inside a society.
// @Tags Society Members
// @Produce json
// @Param societyId path int true "Society ID"
// @Param memberId path int true "Member ID"
// @Success 200 {object} models.SocietyMemberAPIResponse "Member fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter"
// @Failure 404 {object} models.ErrorResponseDoc "Member not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/societies/{societyId}/members/{memberId} [get]
func (h *SocietyHandler) GetSocietyMember(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	memberID, ok := parsePathInt64(c, "memberId")
	if !ok {
		return
	}
	result, err := h.societySvc.GetSocietyMemberDetail(c.Request.Context(), models.GetSocietyMemberFilter{ID: &memberID, SocietyID: &societyID})
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Member fetched successfully", gin.H{"member_detail": result})
}

// GetSocietyMemberSummary godoc
// @Summary Get society member summary
// @Description [User] Returns member counts for a society.
// @Tags Society Members
// @Produce json
// @Param societyId path int true "Society ID"
// @Success 200 {object} models.SocietyMemberSummaryAPIResponse "Member summary fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid society ID"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/societies/{societyId}/members/summary [get]
func (h *SocietyHandler) GetSocietyMemberSummary(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	result, err := h.societySvc.GetSocietyMemberSummary(c.Request.Context(), societyID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Member summary fetched successfully", gin.H{"summary": result})
}

// ListSocietyMembers godoc
// @Summary List society members
// @Description [User] Lists society members with filters for owner/admin member management screens.
// @Tags Society Members
// @Produce json
// @Param societyId path int true "Society ID"
// @Param search query string false "Search full name, email, phone, role, or status"
// @Param role query string false "Member role"
// @Param status query string false "Member status"
// @Param user_id query int false "User ID"
// @Param invited_by query int false "Invited by user ID"
// @Param removed_by query int false "Removed by user ID"
// @Param joined_from query string false "Joined from RFC3339 timestamp"
// @Param joined_to query string false "Joined to RFC3339 timestamp"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Param sort_by query string false "Sort by: joined_at, role, status"
// @Param sort_order query string false "Sort order: asc, desc"
// @Success 200 {object} models.PaginatedMembersAPIResponse "Members fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid query parameter"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/societies/{societyId}/members [get]
func (h *SocietyHandler) ListSocietyMembers(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	filter, ok := listMembersFilterFromQuery(c, societyID)
	// fmt.Printf("ListSocietyMembers filter: %+v\n", filter) // Debug log to check filter values
	if !ok {
		return
	}
	result, err := h.societySvc.ListSocietyMembers(c.Request.Context(), filter)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Members fetched successfully", gin.H{"members": result})
}

// ListAllSocietyMember godoc
// @Summary List all society members for developer
// @Description [Developer] Lists society members with filters without requiring owner/admin society membership.
// @Tags Developer
// @Produce json
// @Param societyId path int true "Society ID"
// @Param search query string false "Search full name, email, phone, role, or status"
// @Param role query string false "Member role"
// @Param status query string false "Member status"
// @Param user_id query int false "User ID"
// @Param invited_by query int false "Invited by user ID"
// @Param removed_by query int false "Removed by user ID"
// @Param joined_from query string false "Joined from RFC3339 timestamp"
// @Param joined_to query string false "Joined to RFC3339 timestamp"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Param sort_by query string false "Sort by: joined_at, role, status"
// @Param sort_order query string false "Sort order: asc, desc"
// @Success 200 {object} models.PaginatedMembersAPIResponse "Members fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid query parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Unauthorized"
// @Failure 403 {object} models.ErrorResponseDoc "Forbidden"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/societies/{societyId}/allmember [get]
func (h *SocietyHandler) ListAllSocietyMember(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	filter, ok := listMembersFilterFromQuery(c, societyID)
	if !ok {
		return
	}
	result, err := h.societySvc.ListAllSocietyMember(c.Request.Context(), filter)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Members fetched successfully", gin.H{"members": result})
}

func currentUserID(c *gin.Context) (int64, bool) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		utils.UnauthorizedResponse(c, "Authentication required")
		return 0, false
	}
	return userID, true
}

func parsePathInt64(c *gin.Context, name string) (int64, bool) {
	raw := strings.TrimSpace(c.Param(name))
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || value <= 0 {
		utils.BadRequestResponse(c, name+" must be a positive integer")
		return 0, false
	}
	return value, true
}

func parseSocietyAndUser(c *gin.Context) (int64, int64, bool) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return 0, 0, false
	}
	userID, ok := parsePathInt64(c, "userId")
	if !ok {
		return 0, 0, false
	}
	return societyID, userID, true
}

func listSocietiesFilterFromQuery(c *gin.Context) (models.ListSocietiesFilter, bool) {
	filter := models.ListSocietiesFilter{
		Status:     optionalString(c.Query("status")),
		Search:     strings.TrimSpace(c.Query("search")),
		SearchMode: strings.TrimSpace(c.Query("search_mode")),
		Name:       strings.TrimSpace(c.Query("name")),
		Code:       strings.TrimSpace(c.Query("code")),
		City:       strings.TrimSpace(c.Query("city")),
		State:      strings.TrimSpace(c.Query("state")),
		Country:    strings.TrimSpace(c.Query("country")),
		Pincode:    strings.TrimSpace(c.Query("pincode")),
		SortBy:     strings.TrimSpace(c.Query("sort_by")),
		SortOrder:  strings.ToLower(strings.TrimSpace(c.Query("sort_order"))),
	}
	if !queryInt64Ptr(c, "id", &filter.ID) || !queryInt64Ptr(c, "created_by", &filter.CreatedBy) ||
		!queryInt64Ptr(c, "approved_by", &filter.ApprovedBy) || !queryInt64Ptr(c, "rejected_by", &filter.RejectedBy) ||
		!queryInt64Ptr(c, "suspended_by", &filter.SuspendedBy) {
		return models.ListSocietiesFilter{}, false
	}
	if !queryTimePtr(c, "created_from", &filter.CreatedFrom) || !queryTimePtr(c, "created_to", &filter.CreatedTo) {
		return models.ListSocietiesFilter{}, false
	}
	limit, offset, ok := paginationQuery(c)
	if !ok {
		return models.ListSocietiesFilter{}, false
	}
	filter.Limit = limit
	filter.Offset = offset
	return filter, true
}

func listMembersFilterFromQuery(c *gin.Context, societyID int64) (models.ListSocietyMembersFilter, bool) {
	filter := models.ListSocietyMembersFilter{
		SocietyID:  societyID,
		Search:     strings.TrimSpace(c.Query("search")),
		SearchMode: strings.TrimSpace(c.Query("search_mode")),
		Role:       optionalString(c.Query("role")),
		Status:     optionalString(c.Query("status")),
		SortBy:     strings.TrimSpace(c.Query("sort_by")),
		SortOrder:  strings.ToLower(strings.TrimSpace(c.Query("sort_order"))),
	}
	if !queryInt64Ptr(c, "user_id", &filter.UserID) || !queryInt64Ptr(c, "invited_by", &filter.InvitedBy) ||
		!queryInt64Ptr(c, "removed_by", &filter.RemovedBy) {
		return models.ListSocietyMembersFilter{}, false
	}
	if !queryTimePtr(c, "joined_from", &filter.JoinedFrom) || !queryTimePtr(c, "joined_to", &filter.JoinedTo) {
		return models.ListSocietyMembersFilter{}, false
	}
	limit, offset, ok := paginationQuery(c)
	if !ok {
		return models.ListSocietyMembersFilter{}, false
	}
	filter.Limit = limit
	filter.Offset = offset
	return filter, true
}

func optionalString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func queryInt64Ptr(c *gin.Context, name string, out **int64) bool {
	raw := strings.TrimSpace(c.Query(name))
	if raw == "" {
		return true
	}
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || value <= 0 {
		utils.BadRequestResponse(c, name+" must be a positive integer")
		return false
	}
	*out = &value
	return true
}

func queryTimePtr(c *gin.Context, name string, out **time.Time) bool {
	raw := strings.TrimSpace(c.Query(name))
	if raw == "" {
		return true
	}
	value, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		utils.BadRequestResponse(c, name+" must be a valid RFC3339 timestamp")
		return false
	}
	*out = &value
	return true
}

func paginationQuery(c *gin.Context) (int32, int32, bool) {
	limit := int32(20)
	offset := int32(0)
	if raw := strings.TrimSpace(c.Query("limit")); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 32)
		if err != nil || value <= 0 {
			utils.BadRequestResponse(c, "limit must be a positive integer")
			return 0, 0, false
		}
		limit = int32(value)
	}
	if raw := strings.TrimSpace(c.Query("offset")); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 32)
		if err != nil || value < 0 {
			utils.BadRequestResponse(c, "offset must be zero or a positive integer")
			return 0, 0, false
		}
		offset = int32(value)
	}
	return limit, offset, true
}
