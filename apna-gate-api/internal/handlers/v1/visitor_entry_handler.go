package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	middleware "go-server/internal/middlewares"
	"go-server/internal/models"
	societysvc "go-server/internal/services/societySvc"
	visitorentrysvc "go-server/internal/services/visitorEntrySvc"
	"go-server/pkg/logger"
	"go-server/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VisitorEntryHandler struct {
	societySvc societysvc.SocietyService
	inviteSvc  visitorentrysvc.VisitorInviteService
	entrySvc   visitorentrysvc.VisitorEntryService
}

func NewVisitorEntryHandler(societySvc societysvc.SocietyService, inviteSvc visitorentrysvc.VisitorInviteService, entrySvc visitorentrysvc.VisitorEntryService) *VisitorEntryHandler {
	return &VisitorEntryHandler{societySvc: societySvc, inviteSvc: inviteSvc, entrySvc: entrySvc}
}

// GetEntryOptions godoc
// @Summary Get visitor entry options
// @Description [Public] Fetches allowed visitor purposes and flat/block options for the public visitor entry form.
// @Tags Visitor Entries
// @Produce json
// @Param societyCode path string true "Society code"
// @Success 200 {object} models.VisitorEntryOptionsAPIResponse "Visitor entry options fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid society code"
// @Failure 404 {object} models.ErrorResponseDoc "Society or flats not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/public/societies/{societyCode}/visitor-entry-options [get]
func (h *VisitorEntryHandler) GetEntryOptions(c *gin.Context) {
	societyID, ok := h.resolvePublicSocietyID(c)
	if !ok {
		return
	}
	result, err := h.entrySvc.GetEntryOptions(c.Request.Context(), societyID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Visitor entry options fetched successfully", gin.H{"options": result})
}

// CreateInvite godoc
// @Summary Create visitor invite
// @Description [Resident] Creates a pre-approved visitor invite for a flat and returns the invite plus QR token details. Service and Staff purposes are not allowed for resident-created invitations.
// @Tags Visitor Entries
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flatId path int true "Flat ID"
// @Param request body models.CreateVisitorInviteRequest true "Visitor invite request"
// @Success 201 {object} models.VisitorInviteTokenAPIResponse "Visitor invite created successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request, validation error, or path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Resident access required"
// @Failure 404 {object} models.ErrorResponseDoc "Society, flat, or visitor settings not found"
// @Failure 409 {object} models.ErrorResponseDoc "Visitor invite cannot be created in current state"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flats/{flatId}/visitor-invites [post]
func (h *VisitorEntryHandler) CreateInvite(c *gin.Context) {
	societyID, flatID, ok := visitorEntryFlatPath(c)
	if !ok {
		return
	}
	var req models.CreateVisitorInviteRequest
	if !bindJSON(c, &req) {
		return
	}
	if err := req.Validate(); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	token, invite, link, err := h.inviteSvc.CreateInvite(c.Request.Context(), societyID, flatID, req.Purpose, userID, req.ExpiresAt)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Visitor invite created successfully", gin.H{"invite": invite, "token": token.QR, "link": link, "short_link": link})
}

// GetInviteByToken godoc
// @Summary Get visitor invite by token
// @Description [Public] Fetches a visitor invite for the web form, or recovers entry/QR when the invite was already used.
// @Tags Visitor Entries
// @Produce json
// @Param token path string true "Visitor invite token"
// @Success 200 {object} models.VisitorInviteAPIResponse "Visitor invite fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid token"
// @Failure 404 {object} models.ErrorResponseDoc "Visitor invite not found"
// @Failure 409 {object} models.ErrorResponseDoc "Visitor invite is expired, used without entry, or cancelled"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/public/visitor-invites/{token} [get]
func (h *VisitorEntryHandler) GetInviteByToken(c *gin.Context) {
	page, err := h.inviteSvc.GetPublicInviteByToken(c.Request.Context(), c.Param("token"))
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Visitor invite fetched successfully", page)
}

// GetGuardInviteByToken godoc
// @Summary Get visitor invite by token (staff)
// @Description [Owner/Admin/Staff] Fetches an active invite only when the invite belongs to the authenticated staff member's society.
// @Tags Visitor Entries
// @Produce json
// @Param societyId path int true "Society ID"
// @Param token path string true "Visitor invite token"
// @Success 200 {object} models.VisitorInviteAPIResponse "Visitor invite fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid society ID or token"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Staff access required"
// @Failure 404 {object} models.ErrorResponseDoc "Visitor invite not found"
// @Failure 409 {object} models.ErrorResponseDoc "Visitor invite is expired, used, or cancelled"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/visitor-invites/{token} [get]
func (h *VisitorEntryHandler) GetGuardInviteByToken(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	invite, err := h.inviteSvc.GetInviteByToken(c.Request.Context(), societyID, userID, c.Param("token"))
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Visitor invite fetched successfully", gin.H{"invite": invite})
}

// CreateStaffInvite godoc
// @Summary Create visitor invite (staff)
// @Description [Staff] Creates a visitor invite for a flat and returns the invite plus shareable token details.
// @Tags Visitor Entries
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flatId path int true "Flat ID"
// @Param request body models.CreateVisitorInviteRequest true "Visitor invite request"
// @Success 201 {object} models.VisitorInviteTokenAPIResponse "Visitor invite created successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request, validation error, or path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Staff access required"
// @Failure 404 {object} models.ErrorResponseDoc "Society, flat, or visitor settings not found"
// @Failure 409 {object} models.ErrorResponseDoc "Visitor invite cannot be created in current state"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flats/{flatId}/visitor-invites/staff [post]
func (h *VisitorEntryHandler) CreateStaffInvite(c *gin.Context) {
	societyID, flatID, ok := visitorEntryFlatPath(c)
	if !ok {
		return
	}
	var req models.CreateVisitorInviteRequest
	if !bindJSON(c, &req) {
		return
	}
	if err := req.Validate(); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	token, invite, link, err := h.inviteSvc.CreateStaffInvite(c.Request.Context(), societyID, flatID, req.Purpose, userID, req.ExpiresAt)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Visitor invite created successfully", gin.H{"invite": invite, "token": token.QR, "link": link, "short_link": link})
}

// SubmitInviteForm godoc
// @Summary Submit visitor invite form
// @Description [Public] Submits visitor details for an active resident invite and creates the visitor entry.
// @Tags Visitor Entries
// @Accept json
// @Produce json
// @Param societyCode path string true "Society code"
// @Param token path string true "Visitor invite token"
// @Param request body models.VisitorFormRequest true "Visitor details"
// @Success 201 {object} models.VisitorEntryMutationAPIResponse "Visitor entry created successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request, token, or validation error"
// @Failure 404 {object} models.ErrorResponseDoc "Visitor invite not found"
// @Failure 409 {object} models.ErrorResponseDoc "Visitor invite is expired, used, or cancelled"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/public/societies/{societyCode}/visitor-invites/{token}/submit [post]
func (h *VisitorEntryHandler) SubmitInviteForm(c *gin.Context) {
	societyID, ok := h.resolvePublicSocietyID(c)
	if !ok {
		return
	}
	var req models.VisitorFormRequest
	if !bindJSON(c, &req) {
		return
	}
	result, err := h.inviteSvc.SubmitInviteForm(c.Request.Context(), societyID, c.Param("token"), req)
	if handleServiceError(c, err) {
		return
	}
	status := http.StatusCreated
	message := "Visitor entry created successfully"
	if result != nil && result.IdempotentReplay {
		status = http.StatusOK
		message = "Visitor entry retrieved successfully"
	}
	utils.SuccessResponse(c, status, message, result)
}

// CancelInvite godoc
// @Summary Cancel visitor invite
// @Description [Resident] Cancels a visitor invite created for a society before it is used.
// @Tags Visitor Entries
// @Produce json
// @Param societyId path int true "Society ID"
// @Param inviteId path int true "Visitor invite ID"
// @Success 200 {object} models.MessageAPIResponse "Visitor invite cancelled successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Invite owner or resident access required"
// @Failure 404 {object} models.ErrorResponseDoc "Visitor invite not found"
// @Failure 409 {object} models.ErrorResponseDoc "Visitor invite cannot be cancelled"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/visitor-invites/{inviteId}/cancel [post]
func (h *VisitorEntryHandler) CancelInvite(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	inviteID, ok := parsePathInt64(c, "inviteId")
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	if err := h.inviteSvc.CancelInvite(c.Request.Context(), societyID, inviteID, userID); handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Visitor invite cancelled successfully", gin.H{"message": "Visitor invite cancelled successfully"})
}

// ListVisitorInviteHistory godoc
// @Summary List visitor invite history
// @Description [Resident] Lists visitor invites for a resident flat with status and pagination filters.
// @Tags Visitor Entries
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flatId path int true "Flat ID"
// @Param status query string false "Visitor invite status" Enums(active, used, expired, cancelled)
// @Param search query string false "Search purpose, status, visitor, contact, or short code"
// @Param limit query int false "Maximum records to return" default(20)
// @Param offset query int false "Records to skip" default(0)
// @Success 200 {object} models.VisitorInviteHistoryAPIResponse "Visitor invites fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid query or path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Resident access required"
// @Failure 404 {object} models.ErrorResponseDoc "Flat not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flats/{flatId}/visitor-invites [get]
func (h *VisitorEntryHandler) ListVisitorInviteHistory(c *gin.Context) {
	societyID, flatID, ok := visitorEntryFlatPath(c)
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	filter, ok := visitorInviteHistoryFilterFromQuery(c)
	if !ok {
		return
	}
	result, err := h.inviteSvc.ListVisitorInviteHistoryForActor(c.Request.Context(), societyID, flatID, userID, filter)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Visitor invites fetched successfully", gin.H{
		"invites": result.Invites,
		"total":   result.Total,
		"limit":   result.Limit,
		"offset":  result.Offset,
	})
}

// GetVisitorInviteHistory godoc
// @Summary Get visitor invite detail
// @Description [Resident] Fetches one visitor invite for a resident flat with linked visitor, entry, and short-link details.
// @Tags Visitor Entries
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flatId path int true "Flat ID"
// @Param inviteId path int true "Visitor invite ID"
// @Success 200 {object} models.VisitorInviteHistoryItemAPIResponse "Visitor invite fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Resident access required"
// @Failure 404 {object} models.ErrorResponseDoc "Visitor invite not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flats/{flatId}/visitor-invites/{inviteId} [get]
func (h *VisitorEntryHandler) GetVisitorInviteHistory(c *gin.Context) {
	societyID, flatID, ok := visitorEntryFlatPath(c)
	if !ok {
		return
	}
	inviteID, ok := parsePathInt64(c, "inviteId")
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	invite, err := h.inviteSvc.GetVisitorInviteHistoryForActor(c.Request.Context(), societyID, flatID, inviteID, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Visitor invite fetched successfully", gin.H{"invite": invite})
}

// CreatePublicQREntry godoc
// @Summary Create public QR visitor entry
// @Description [Public] Creates a visitor entry from a society public QR flow. Approval may be required depending on visitor settings.
// @Tags Visitor Entries
// @Accept json
// @Produce json
// @Param societyCode path string true "Society code"
// @Param request body models.VisitorFormRequest true "Visitor details"
// @Success 201 {object} models.VisitorEntryMutationAPIResponse "Visitor entry created successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request, society code, or validation error"
// @Failure 403 {object} models.ErrorResponseDoc "Public QR visitor entry is disabled"
// @Failure 404 {object} models.ErrorResponseDoc "Society, flat, or visitor settings not found"
// @Failure 409 {object} models.ErrorResponseDoc "Visitor entry cannot be created in current state"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/public/societies/{societyCode}/visitor-entries/public-qr [post]
func (h *VisitorEntryHandler) CreatePublicQREntry(c *gin.Context) {
	h.createPublicEntry(c, h.entrySvc.CreatePublicQREntry)
}

// CreateQuickLinkEntry godoc
// @Summary Create quick-link visitor entry
// @Description [Public] Creates a visitor entry from a quick-link flow. Approval may be required depending on visitor settings.
// @Tags Visitor Entries
// @Accept json
// @Produce json
// @Param societyCode path string true "Society code"
// @Param request body models.VisitorFormRequest true "Visitor details"
// @Success 201 {object} models.VisitorEntryMutationAPIResponse "Visitor entry created successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request, society code, or validation error"
// @Failure 403 {object} models.ErrorResponseDoc "Quick-link visitor entry is disabled"
// @Failure 404 {object} models.ErrorResponseDoc "Society, flat, or visitor settings not found"
// @Failure 409 {object} models.ErrorResponseDoc "Visitor entry cannot be created in current state"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/public/societies/{societyCode}/visitor-entries/quick-link [post]
func (h *VisitorEntryHandler) CreateQuickLinkEntry(c *gin.Context) {
	h.createPublicEntry(c, h.entrySvc.CreateQuickLinkEntry)
}

// CreateGuardEntry godoc
// @Summary Create guard visitor entry
// @Description [Owner/Admin/Staff] Creates a visitor entry from the guard desk for a society. Approval may be required depending on visitor settings.
// @Tags Visitor Entries
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param request body models.VisitorFormRequest true "Visitor details"
// @Success 201 {object} models.VisitorEntryMutationAPIResponse "Visitor entry created successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request, validation error, or society ID"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Owner, admin, or staff access required"
// @Failure 404 {object} models.ErrorResponseDoc "Society, flat, or visitor settings not found"
// @Failure 409 {object} models.ErrorResponseDoc "Visitor entry cannot be created in current state"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/visitor-entries/guard [post]
func (h *VisitorEntryHandler) CreateGuardEntry(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	var req models.VisitorFormRequest
	if !bindJSON(c, &req) {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.entrySvc.CreateGuardEntry(c.Request.Context(), societyID, req, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Visitor entry created successfully", result)
}

// ApproveEntry godoc
// @Summary Approve visitor entry
// @Description [Resident] Approves a waiting visitor entry and returns QR token details for the visitor.
// @Tags Visitor Entries
// @Produce json
// @Param societyId path int true "Society ID"
// @Param entryId path int true "Visitor entry ID"
// @Success 200 {object} models.VisitorEntryMutationAPIResponse "Visitor entry approved successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Resident approval access required"
// @Failure 404 {object} models.ErrorResponseDoc "Visitor entry not found"
// @Failure 409 {object} models.ErrorResponseDoc "Visitor entry is not waiting for approval"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/visitor-entries/{entryId}/approve [post]
func (h *VisitorEntryHandler) ApproveEntry(c *gin.Context) {
	societyID, entryID, ok := visitorEntryPath(c)
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.entrySvc.ApproveEntry(c.Request.Context(), societyID, entryID, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Visitor entry approved successfully", result)
}

// RejectEntry godoc
// @Summary Reject visitor entry
// @Description [Resident] Rejects a waiting visitor entry with a reason.
// @Tags Visitor Entries
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param entryId path int true "Visitor entry ID"
// @Param request body models.RejectVisitorEntryRequest true "Visitor rejection reason"
// @Success 200 {object} models.VisitorEntryMutationAPIResponse "Visitor entry rejected successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request, validation error, or path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Resident rejection access required"
// @Failure 404 {object} models.ErrorResponseDoc "Visitor entry not found"
// @Failure 409 {object} models.ErrorResponseDoc "Visitor entry is not waiting for approval"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/visitor-entries/{entryId}/reject [post]
func (h *VisitorEntryHandler) RejectEntry(c *gin.Context) {
	societyID, entryID, ok := visitorEntryPath(c)
	if !ok {
		return
	}
	var req models.RejectVisitorEntryRequest
	if !bindJSON(c, &req) {
		return
	}
	if err := req.Validate(); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.entrySvc.RejectEntry(c.Request.Context(), societyID, entryID, req.Reason, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Visitor entry rejected successfully", result)
}

// ValidateQR godoc
// @Summary Validate visitor QR token
// @Description [Public] Validates a visitor QR token and returns the matching approved visitor entry.
// @Tags Visitor Entries
// @Accept json
// @Produce json
// @Param request body models.QRTokenRequest true "Visitor QR token"
// @Success 200 {object} models.VisitorEntryAPIResponse "Visitor QR validated successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request, token, or visitor entry state"
// @Failure 404 {object} models.ErrorResponseDoc "Visitor entry not found"
// @Failure 409 {object} models.ErrorResponseDoc "QR token is expired, unavailable, or already used"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/public/visitor-entries/qr/validate [post]
func (h *VisitorEntryHandler) ValidateQR(c *gin.Context) {
	var req models.QRTokenRequest
	if !bindJSON(c, &req) {
		return
	}
	if err := req.Validate(); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	entry, err := h.entrySvc.GetEntryForQRScanForViewer(c.Request.Context(), req.Token, viewerUserID(c))
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Visitor QR validated successfully", gin.H{"entry": entry})
}

// CheckIn godoc
// @Summary Check in visitor
// @Description [Owner/Admin/Staff] Checks in a visitor using a valid QR token.
// @Tags Visitor Entries
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param request body models.QRTokenRequest true "Visitor QR token"
// @Success 200 {object} models.VisitorEntryAPIResponse "Visitor checked in successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request, token, or visitor entry state"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Owner, admin, or staff access required"
// @Failure 404 {object} models.ErrorResponseDoc "Visitor entry not found"
// @Failure 409 {object} models.ErrorResponseDoc "QR token is expired, unavailable, or already used"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/visitor-entries/check-in [post]
func (h *VisitorEntryHandler) CheckIn(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	var req models.QRTokenRequest
	if !bindJSON(c, &req) {
		return
	}
	if err := req.Validate(); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	preview, err := h.entrySvc.ValidateQR(c.Request.Context(), req.Token)
	if handleServiceError(c, err) {
		return
	}
	if preview.SocietyID != societyID {
		utils.ErrorResponse(c, http.StatusForbidden, models.ErrCodeForbidden, "Visitor entry does not belong to this society", nil)
		return
	}
	entry, err := h.entrySvc.CheckIn(c.Request.Context(), req.Token, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Visitor checked in successfully", gin.H{"entry": entry})
}

// CheckOut godoc
// @Summary Check out visitor
// @Description [Owner/Admin/Staff] Checks out a checked-in visitor entry.
// @Tags Visitor Entries
// @Produce json
// @Param societyId path int true "Society ID"
// @Param entryId path int true "Visitor entry ID"
// @Success 200 {object} models.VisitorEntryAPIResponse "Visitor checked out successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter or visitor entry state"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Owner, admin, or staff access required"
// @Failure 404 {object} models.ErrorResponseDoc "Visitor entry not found"
// @Failure 409 {object} models.ErrorResponseDoc "Visitor is not checked in or already checked out"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/visitor-entries/{entryId}/check-out [post]
func (h *VisitorEntryHandler) CheckOut(c *gin.Context) {
	societyID, entryID, ok := visitorEntryPath(c)
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	entry, err := h.entrySvc.CheckOut(c.Request.Context(), societyID, entryID, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Visitor checked out successfully", gin.H{"entry": entry})
}

// GetEntry godoc
// @Summary Get visitor entry
// @Description [Owner/Admin/Staff] Fetches one visitor entry by ID.
// @Tags Visitor Entries
// @Produce json
// @Param societyId path int true "Society ID"
// @Param entryId path int true "Visitor entry ID"
// @Success 200 {object} models.VisitorEntryAPIResponse "Visitor entry fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Owner, admin, or staff access required"
// @Failure 404 {object} models.ErrorResponseDoc "Visitor entry not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/visitor-entries/{entryId} [get]
func (h *VisitorEntryHandler) GetEntry(c *gin.Context) {
	societyID, entryID, ok := visitorEntryPath(c)
	if !ok {
		return
	}
	entry, err := h.entrySvc.GetEntry(c.Request.Context(), societyID, entryID)
	if handleServiceError(c, err) {
		return
	}
	logger.Debug("Fetched visitor entry: ", zap.Any("entry", entry))
	// fmt.Println("Fetched visitor entry: ", entry)
	utils.SuccessResponse(c, http.StatusOK, "Visitor entry fetched successfully", gin.H{"entry": entry})
}

// UpdateGuardEntry godoc
// @Summary Update visitor entry details
// @Description [Owner/Admin/Staff] Updates visitor and entry details before check-in for approved or pending entries.
// @Tags Visitor Entries
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param entryId path int true "Visitor entry ID"
// @Param request body models.UpdateGuardVisitorEntryRequest true "Visitor entry update payload"
// @Success 200 {object} models.VisitorEntryMutationAPIResponse "Visitor entry updated successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request or validation error"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Owner, admin, or staff access required"
// @Failure 404 {object} models.ErrorResponseDoc "Visitor entry not found"
// @Failure 409 {object} models.ErrorResponseDoc "Visitor entry is not in an editable state"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/visitor-entries/{entryId} [patch]
func (h *VisitorEntryHandler) UpdateGuardEntry(c *gin.Context) {
	societyID, entryID, ok := visitorEntryPath(c)
	if !ok {
		return
	}
	guardUserID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		utils.UnauthorizedResponse(c, "Authentication required")
		return
	}

	var req models.UpdateGuardVisitorEntryRequest
	if !bindJSON(c, &req) {
		return
	}

	entry, err := h.entrySvc.UpdateGuardEntry(c.Request.Context(), societyID, entryID, guardUserID, req)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Visitor entry updated successfully", gin.H{"entry": entry})
}

// ListEntries godoc
// @Summary List visitor entries
// @Description [Owner/Admin/Staff] Lists visitor entries for a society with optional flat, status, source, purpose, limit, and offset filters.
// @Tags Visitor Entries
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flat_id query int false "Flat ID"
// @Param status query string false "Visitor status" Enums(waiting_approval, approved, rejected, checked_in, checked_out, cancelled, expired, auto_closed)
// @Param source query string false "Visitor entry source" Enums(resident_link, public_qr, guard_entry, quick_link)
// @Param purpose query string false "Visitor purpose" Enums(guest, delivery, cab, service, maintenance, staff, other)
// @Param limit query int false "Maximum records to return" default(50)
// @Param offset query int false "Records to skip" default(0)
// @Success 200 {object} models.VisitorEntriesAPIResponse "Visitor entries fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid query or path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Owner, admin, or staff access required"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/visitor-entries [get]
func (h *VisitorEntryHandler) ListEntries(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	filter, ok := visitorEntryFilterFromQuery(c, societyID)
	if !ok {
		return
	}
	entries, err := h.entrySvc.ListEntriesPaginated(c.Request.Context(), filter)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Visitor entries fetched successfully", gin.H{
		"entries": entries.Entries,
		"total":   entries.Total,
		"limit":   entries.Limit,
		"offset":  entries.Offset,
	})
}

// ListPendingApprovals godoc
// @Summary List pending visitor approvals
// @Description [Resident] Lists visitor entries waiting for approval for a resident flat.
// @Tags Visitor Entries
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flatId path int true "Flat ID"
// @Success 200 {object} models.VisitorEntriesAPIResponse "Pending visitor approvals fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Resident access required"
// @Failure 404 {object} models.ErrorResponseDoc "Flat not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flats/{flatId}/visitor-entries/pending [get]
func (h *VisitorEntryHandler) ListPendingApprovals(c *gin.Context) {
	societyID, flatID, ok := visitorEntryFlatPath(c)
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	entries, err := h.entrySvc.ListPendingApprovals(c.Request.Context(), societyID, flatID, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Pending visitor approvals fetched successfully", gin.H{"entries": entries})
}

// ListEvents godoc
// @Summary List visitor entry events
// @Description [Owner/Admin/Staff] Lists audit events recorded for a visitor entry.
// @Tags Visitor Entries
// @Produce json
// @Param societyId path int true "Society ID"
// @Param entryId path int true "Visitor entry ID"
// @Success 200 {object} models.VisitorEntryEventsAPIResponse "Visitor entry events fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Owner, admin, or staff access required"
// @Failure 404 {object} models.ErrorResponseDoc "Visitor entry not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/visitor-entries/{entryId}/events [get]
func (h *VisitorEntryHandler) ListEvents(c *gin.Context) {
	societyID, entryID, ok := visitorEntryPath(c)
	if !ok {
		return
	}
	events, err := h.entrySvc.ListEvents(c.Request.Context(), societyID, entryID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Visitor entry events fetched successfully", gin.H{"events": events})
}

// GetEntryStats godoc
// @Summary Get visitor entry stats
// @Description [Owner/Admin/Staff] Returns visitor dashboard statistics for a society.
// @Tags Visitor Entries
// @Produce json
// @Param societyId path int true "Society ID"
// @Success 200 {object} models.VisitorEntryStatsAPIResponse "Visitor entry stats fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid society ID"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Owner, admin, or staff access required"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/visitor-entries/stats [get]
func (h *VisitorEntryHandler) GetEntryStats(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}

	var rangeFrom, rangeTo *time.Time
	if !queryTimePtr(c, "event_from", &rangeFrom) || !queryTimePtr(c, "event_to", &rangeTo) {
		return
	}

	var stats *models.VisitorEntryStatsResponse
	var err error
	if rangeFrom != nil && rangeTo != nil {
		if !rangeFrom.Before(*rangeTo) {
			utils.BadRequestResponse(c, "event_from must be before event_to")
			return
		}
		stats, err = h.entrySvc.GetEntryStatsInRange(c.Request.Context(), societyID, *rangeFrom, *rangeTo)
	} else if rangeFrom != nil || rangeTo != nil {
		utils.BadRequestResponse(c, "event_from and event_to must both be provided")
		return
	} else {
		stats, err = h.entrySvc.GetEntryStats(c.Request.Context(), societyID)
	}
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Visitor entry stats fetched successfully", gin.H{"stats": stats})
}

// GetEntryDailyStats godoc
// @Summary Get visitor entry daily stats
// @Description [Owner/Admin/Staff] Returns daily visitor counts for charting.
// @Tags Visitor Entries
// @Produce json
// @Param societyId path int true "Society ID"
// @Param days query int false "Number of days (default 7, max 90)"
// @Success 200 {object} models.VisitorEntryDailyStatsAPIResponse "Visitor entry daily stats fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid society ID"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Owner, admin, or staff access required"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/visitor-entries/stats/daily [get]
func (h *VisitorEntryHandler) GetEntryDailyStats(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}

	days := int32(7)
	if rawDays := strings.TrimSpace(c.Query("days")); rawDays != "" {
		parsed, err := strconv.ParseInt(rawDays, 10, 32)
		if err != nil || parsed <= 0 {
			utils.BadRequestResponse(c, "days must be a positive integer")
			return
		}
		days = int32(parsed)
	}

	stats, err := h.entrySvc.GetDailyEntryStats(c.Request.Context(), societyID, days)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Visitor entry daily stats fetched successfully", gin.H{"stats": stats})
}

// GetGuardDeskBootstrap godoc
// @Summary Get guard desk bootstrap
// @Description [Owner/Admin/Staff] Returns aggregated guard desk dashboard data for a society.
// @Tags Visitor Entries
// @Produce json
// @Param societyId path int true "Society ID"
// @Success 200 {object} models.GuardDeskBootstrapAPIResponse "Guard desk bootstrap fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid society ID"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Owner, admin, or staff access required"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/guard-desk/bootstrap [get]
func (h *VisitorEntryHandler) GetGuardDeskBootstrap(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	result, err := h.entrySvc.GetGuardDeskBootstrap(c.Request.Context(), societyID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Guard desk bootstrap fetched successfully", gin.H{"desk": result})
}

// ListSocietyPendingApprovals godoc
// @Summary List society pending visitor approvals
// @Description [Owner/Admin/Staff] Lists visitor entries waiting for approval across the society.
// @Tags Visitor Entries
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flat_id query int false "Flat ID"
// @Param block query string false "Block"
// @Param limit query int false "Maximum records to return" default(50)
// @Param offset query int false "Records to skip" default(0)
// @Success 200 {object} models.VisitorPendingEntriesAPIResponse "Pending visitor approvals fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid query or path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Owner, admin, or staff access required"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/visitor-entries/pending [get]
func (h *VisitorEntryHandler) ListSocietyPendingApprovals(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	filter, ok := visitorPendingFilterFromQuery(c, societyID)
	if !ok {
		return
	}
	result, err := h.entrySvc.ListSocietyPendingApprovals(c.Request.Context(), filter)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Pending visitor approvals fetched successfully", gin.H{
		"entries": result.Entries,
		"total":   result.Total,
		"limit":   result.Limit,
		"offset":  result.Offset,
	})
}

// ListWaitingAtGate godoc
// @Summary List visitors waiting at gate
// @Description [Owner/Admin/Staff] Lists approved non-invite visitor entries awaiting manual check-in. Excludes resident-link QR guests (see Expected Guests).
// @Tags Visitor Entries
// @Produce json
// @Param societyId path int true "Society ID"
// @Param search query string false "Search by name, phone, flat, vehicle, delivery partner, or purpose"
// @Param limit query int false "Maximum records to return" default(50)
// @Param offset query int false "Records to skip" default(0)
// @Success 200 {object} models.VisitorEntriesAPIResponse "Waiting at gate entries fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid query or path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Owner, admin, or staff access required"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/visitor-entries/waiting-at-gate [get]
func (h *VisitorEntryHandler) ListWaitingAtGate(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	filter, ok := waitingAtGateFilterFromQuery(c, societyID)
	if !ok {
		return
	}
	result, err := h.entrySvc.ListWaitingAtGate(c.Request.Context(), filter)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Waiting at gate entries fetched successfully", gin.H{
		"entries": result.Entries,
		"total":   result.Total,
		"limit":   result.Limit,
		"offset":  result.Offset,
	})
}

// ListExpectedGuests godoc
// @Summary List expected guest entries
// @Description [Owner/Admin/Staff] Lists pre-invited resident-link visitor entries expected today and ready for arrival.
// @Tags Visitor Entries
// @Produce json
// @Param societyId path int true "Society ID"
// @Param search query string false "Search by name, phone, flat, vehicle, or purpose"
// @Param event_from query string false "Expected window start (RFC3339, defaults to IST today start)"
// @Param event_to query string false "Expected window end (RFC3339, defaults to IST tomorrow start)"
// @Param limit query int false "Maximum records to return" default(50)
// @Param offset query int false "Records to skip" default(0)
// @Success 200 {object} models.VisitorEntriesAPIResponse "Expected guest entries fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid query or path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Owner, admin, or staff access required"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/visitor-entries/expected-guests [get]
func (h *VisitorEntryHandler) ListExpectedGuests(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	filter, ok := expectedGuestFilterFromQuery(c, societyID)
	if !ok {
		return
	}
	result, err := h.entrySvc.ListExpectedGuests(c.Request.Context(), filter)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Expected guest entries fetched successfully", gin.H{
		"entries": result.Entries,
		"total":   result.Total,
		"limit":   result.Limit,
		"offset":  result.Offset,
	})
}

// NotifyPendingEntry godoc
// @Summary Notify resident about pending visitor
// @Description [Owner/Admin/Staff] Re-sends a pending approval notification to flat residents.
// @Tags Visitor Entries
// @Produce json
// @Param societyId path int true "Society ID"
// @Param entryId path int true "Visitor entry ID"
// @Success 200 {object} models.APIResponse "Resident notified successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter or visitor entry state"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Owner, admin, or staff access required"
// @Failure 404 {object} models.ErrorResponseDoc "Visitor entry not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/visitor-entries/{entryId}/notify [post]
func (h *VisitorEntryHandler) NotifyPendingEntry(c *gin.Context) {
	societyID, entryID, ok := visitorEntryPath(c)
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	if err := h.entrySvc.NotifyPendingEntry(c.Request.Context(), societyID, entryID, userID); handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Resident notified successfully", nil)
}

// GuardApproveEntry godoc
// @Summary Guard approve visitor entry
// @Description [Owner/Admin/Staff] Approves a pending visitor entry without check-in.
// @Tags Visitor Entries
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param entryId path int true "Visitor entry ID"
// @Param request body models.GuardApproveEntryRequest false "Guard approval options"
// @Success 200 {object} models.VisitorEntryMutationAPIResponse "Visitor entry approved successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter or visitor entry state"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Owner, admin, or staff access required"
// @Failure 404 {object} models.ErrorResponseDoc "Visitor entry not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/visitor-entries/{entryId}/guard-approve [post]
func (h *VisitorEntryHandler) GuardApproveEntry(c *gin.Context) {
	societyID, entryID, ok := visitorEntryPath(c)
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	opts := guardApproveOptionsFromBody(c)
	result, err := h.entrySvc.GuardApproveEntry(c.Request.Context(), societyID, entryID, userID, opts)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Visitor entry approved successfully", result)
}

// GuardApproveAndCheckIn godoc
// @Summary Guard approve and check in visitor
// @Description [Owner/Admin/Staff] Atomically approves and checks in a pending visitor entry.
// @Tags Visitor Entries
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param entryId path int true "Visitor entry ID"
// @Param request body models.GuardApproveEntryRequest false "Guard approval options"
// @Success 200 {object} models.VisitorEntryAPIResponse "Visitor approved and checked in successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter or visitor entry state"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Owner, admin, or staff access required"
// @Failure 404 {object} models.ErrorResponseDoc "Visitor entry not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/visitor-entries/{entryId}/approve-and-check-in [post]
func (h *VisitorEntryHandler) GuardApproveAndCheckIn(c *gin.Context) {
	societyID, entryID, ok := visitorEntryPath(c)
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	opts := guardApproveOptionsFromBody(c)
	entry, err := h.entrySvc.GuardApproveAndCheckIn(c.Request.Context(), societyID, entryID, userID, opts)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Visitor approved and checked in successfully", gin.H{"entry": entry})
}

// CheckInByEntryID godoc
// @Summary Check in visitor by entry ID
// @Description [Owner/Admin/Staff] Checks in an approved visitor entry without scanning QR.
// @Tags Visitor Entries
// @Produce json
// @Param societyId path int true "Society ID"
// @Param entryId path int true "Visitor entry ID"
// @Success 200 {object} models.VisitorEntryAPIResponse "Visitor checked in successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter or visitor entry state"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Owner, admin, or staff access required"
// @Failure 404 {object} models.ErrorResponseDoc "Visitor entry not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/visitor-entries/{entryId}/check-in [post]
func (h *VisitorEntryHandler) CheckInByEntryID(c *gin.Context) {
	societyID, entryID, ok := visitorEntryPath(c)
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	entry, err := h.entrySvc.CheckInByEntryID(c.Request.Context(), societyID, entryID, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Visitor checked in successfully", gin.H{"entry": entry})
}

// ListFlatVisitorEntries godoc
// @Summary List flat visitor entries
// @Description [Resident] Lists visitor entries for a resident flat with optional status, purpose, and date filters. created_from/created_to are legacy creation-time filters; event/event_from/event_to filter by visitor lifecycle timestamp.
// @Tags Visitor Entries
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flatId path int true "Flat ID"
// @Param status query string false "Visitor status" Enums(waiting_approval, approved, rejected, checked_in, checked_out, cancelled, expired, auto_closed)
// @Param purpose query string false "Visitor purpose" Enums(guest, delivery, cab, service, maintenance, staff, other)
// @Param created_from query string false "Created from (RFC3339)"
// @Param created_to query string false "Created to (RFC3339)"
// @Param event query string false "Lifecycle timestamp to filter by" Enums(created, checked_in, checked_out, expected, activity)
// @Param event_from query string false "Lifecycle filter start (RFC3339)"
// @Param event_to query string false "Lifecycle filter end (RFC3339)"
// @Param limit query int false "Maximum records to return" default(50)
// @Param offset query int false "Records to skip" default(0)
// @Success 200 {object} models.VisitorEntriesAPIResponse "Visitor entries fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid query or path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Resident access required"
// @Failure 404 {object} models.ErrorResponseDoc "Flat not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flats/{flatId}/visitor-entries [get]
func (h *VisitorEntryHandler) ListFlatVisitorEntries(c *gin.Context) {
	societyID, flatID, ok := visitorEntryFlatPath(c)
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	filter, ok := visitorEntryFilterFromFlatQuery(c, societyID, flatID)
	if !ok {
		return
	}
	entries, err := h.entrySvc.ListFlatEntriesForActor(c.Request.Context(), societyID, flatID, userID, filter)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Visitor entries fetched successfully", gin.H{
		"entries": entries.Entries,
		"total":   entries.Total,
		"limit":   entries.Limit,
		"offset":  entries.Offset,
	})
}

// GetFlatVisitorEntry godoc
// @Summary Get flat visitor entry
// @Description [Resident] Fetches a single visitor entry for a resident flat.
// @Tags Visitor Entries
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flatId path int true "Flat ID"
// @Param entryId path int true "Visitor entry ID"
// @Success 200 {object} models.VisitorEntryAPIResponse "Visitor entry fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Resident access required"
// @Failure 404 {object} models.ErrorResponseDoc "Visitor entry or flat not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flats/{flatId}/visitor-entries/{entryId} [get]
func (h *VisitorEntryHandler) GetFlatVisitorEntry(c *gin.Context) {
	societyID, flatID, ok := visitorEntryFlatPath(c)
	if !ok {
		return
	}
	entryID, ok := parsePathInt64(c, "entryId")
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	entry, err := h.entrySvc.GetFlatEntryForActor(c.Request.Context(), societyID, flatID, entryID, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Visitor entry fetched successfully", gin.H{"entry": entry})
}

// GetFlatVisitorContextForResident godoc
// @Summary Get flat visitor context for resident
// @Description [Resident] Returns occupancy, visitor settings summary, and recent visitors for a resident flat.
// @Tags Visitor Entries
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flatId path int true "Flat ID"
// @Success 200 {object} models.FlatVisitorContextAPIResponse "Flat visitor context fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Resident access required"
// @Failure 404 {object} models.ErrorResponseDoc "Flat or visitor settings not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flats/{flatId}/visitor-context [get]
func (h *VisitorEntryHandler) GetFlatVisitorContextForResident(c *gin.Context) {
	societyID, flatID, ok := visitorEntryFlatPath(c)
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.entrySvc.GetFlatVisitorContextForActor(c.Request.Context(), societyID, flatID, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Flat visitor context fetched successfully", gin.H{"context": result})
}

// GetFlatVisitorContext is currently unregistered. The public route is documented
// on GetFlatVisitorContextForResident, which enforces flat viewing authorization.
func (h *VisitorEntryHandler) GetFlatVisitorContext(c *gin.Context) {
	societyID, flatID, ok := visitorEntryFlatPath(c)
	if !ok {
		return
	}
	result, err := h.entrySvc.GetFlatVisitorContext(c.Request.Context(), societyID, flatID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Flat visitor context fetched successfully", gin.H{"context": result})
}

// GetMemberVisitorApprovalStats godoc
// @Summary Get member visitor approval stats
// @Description [Owner/Admin] Returns how many visitor entries a member has approved or rejected.
// @Tags Visitor Entries
// @Produce json
// @Param societyId path int true "Society ID"
// @Param memberId path int true "Member ID"
// @Success 200 {object} models.MemberVisitorApprovalStatsAPIResponse "Member visitor approval stats fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 404 {object} models.ErrorResponseDoc "Member not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/members/{memberId}/visitor-approval-stats [get]
func (h *VisitorEntryHandler) GetMemberVisitorApprovalStats(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	memberID, ok := parsePathInt64(c, "memberId")
	if !ok {
		return
	}
	stats, err := h.entrySvc.GetMemberVisitorApprovalStats(c.Request.Context(), societyID, memberID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Member visitor approval stats fetched successfully", gin.H{"stats": stats})
}

func (h *VisitorEntryHandler) createPublicEntry(c *gin.Context, create func(context.Context, int64, models.VisitorFormRequest) (*models.VisitorEntryMutationResponse, error)) {
	societyID, ok := h.resolvePublicSocietyID(c)
	if !ok {
		return
	}
	var req models.VisitorFormRequest
	if !bindJSON(c, &req) {
		return
	}
	result, err := create(c.Request.Context(), societyID, req)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Visitor entry created successfully", result)
}

func (h *VisitorEntryHandler) resolvePublicSocietyID(c *gin.Context) (int64, bool) {
	if raw := strings.TrimSpace(c.Param("societyId")); raw != "" {
		return parsePathInt64(c, "societyId")
	}
	societyCode := strings.TrimSpace(c.Param("societyCode"))
	if societyCode == "" {
		utils.BadRequestResponse(c, "societyCode is required")
		return 0, false
	}
	societyID, err := h.societySvc.ResolveActiveSocietyIDByCode(c.Request.Context(), societyCode)
	if handleServiceError(c, err) {
		return 0, false
	}
	return societyID, true
}

func viewerUserID(c *gin.Context) *int64 {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok || userID <= 0 {
		return nil
	}
	return &userID
}

func visitorEntryPath(c *gin.Context) (int64, int64, bool) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return 0, 0, false
	}
	entryID, ok := parsePathInt64(c, "entryId")
	if !ok {
		return 0, 0, false
	}
	return societyID, entryID, true
}

func visitorEntryFlatPath(c *gin.Context) (int64, int64, bool) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return 0, 0, false
	}
	flatID, ok := parsePathInt64(c, "flatId")
	if !ok {
		return 0, 0, false
	}
	return societyID, flatID, true
}

func visitorEntryFilterFromFlatQuery(c *gin.Context, societyID int64, flatID int64) (models.VisitorEntryFilter, bool) {
	filter := models.VisitorEntryFilter{SocietyID: societyID, FlatID: &flatID, Limit: 50}
	if raw := c.Query("status"); raw != "" {
		value := models.VisitorStatus(raw)
		if !value.IsValid() {
			utils.BadRequestResponse(c, "invalid status")
			return filter, false
		}
		filter.Status = &value
	}
	if raw := c.Query("purpose"); raw != "" {
		value := models.VisitorPurpose(raw)
		if !value.IsValid() {
			utils.BadRequestResponse(c, "invalid purpose")
			return filter, false
		}
		filter.Purpose = &value
	}
	if raw := c.Query("limit"); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 32)
		if err != nil || value <= 0 {
			utils.BadRequestResponse(c, "limit must be a positive integer")
			return filter, false
		}
		filter.Limit = int32(value)
	}
	if raw := c.Query("offset"); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 32)
		if err != nil || value < 0 {
			utils.BadRequestResponse(c, "offset must be zero or positive")
			return filter, false
		}
		filter.Offset = int32(value)
	}
	if !applyVisitorEntryDateFilters(c, &filter) {
		return filter, false
	}
	if raw := strings.TrimSpace(c.Query("search")); raw != "" {
		filter.Search = &raw
	}
	return filter, true
}

func visitorInviteHistoryFilterFromQuery(c *gin.Context) (models.VisitorInviteHistoryFilter, bool) {
	filter := models.VisitorInviteHistoryFilter{Limit: 20}
	if raw := c.Query("status"); raw != "" {
		value := models.VisitorInviteStatus(raw)
		if !value.IsValid() {
			utils.BadRequestResponse(c, "invalid status")
			return filter, false
		}
		filter.Status = &value
	}
	if raw := c.Query("limit"); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 32)
		if err != nil || value <= 0 {
			utils.BadRequestResponse(c, "limit must be a positive integer")
			return filter, false
		}
		if value > 100 {
			value = 100
		}
		filter.Limit = int32(value)
	}
	if raw := c.Query("offset"); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 32)
		if err != nil || value < 0 {
			utils.BadRequestResponse(c, "offset must be zero or positive")
			return filter, false
		}
		filter.Offset = int32(value)
	}
	if raw := strings.TrimSpace(c.Query("search")); raw != "" {
		filter.Search = &raw
	}
	return filter, true
}

func visitorEntryFilterFromQuery(c *gin.Context, societyID int64) (models.VisitorEntryFilter, bool) {
	filter := models.VisitorEntryFilter{SocietyID: societyID, Limit: 50}
	if raw := c.Query("flat_id"); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || value <= 0 {
			utils.BadRequestResponse(c, "flat_id must be a positive integer")
			return filter, false
		}
		filter.FlatID = &value
	}
	if raw := c.Query("status"); raw != "" {
		value := models.VisitorStatus(raw)
		if !value.IsValid() {
			utils.BadRequestResponse(c, "invalid status")
			return filter, false
		}
		filter.Status = &value
	}
	if raw := c.Query("source"); raw != "" {
		value := models.VisitorEntrySource(raw)
		if !value.IsValid() {
			utils.BadRequestResponse(c, "invalid source")
			return filter, false
		}
		filter.Source = &value
	}
	if raw := c.Query("purpose"); raw != "" {
		value := models.VisitorPurpose(raw)
		if !value.IsValid() {
			utils.BadRequestResponse(c, "invalid purpose")
			return filter, false
		}
		filter.Purpose = &value
	}
	if raw := c.Query("limit"); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 32)
		if err != nil || value <= 0 {
			utils.BadRequestResponse(c, "limit must be a positive integer")
			return filter, false
		}
		filter.Limit = int32(value)
	}
	if raw := c.Query("offset"); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 32)
		if err != nil || value < 0 {
			utils.BadRequestResponse(c, "offset must be zero or positive")
			return filter, false
		}
		filter.Offset = int32(value)
	}
	if raw := c.Query("block"); raw != "" {
		filter.Block = &raw
	}
	if !applyVisitorEntryDateFilters(c, &filter) {
		return filter, false
	}
	if raw := strings.TrimSpace(c.Query("search")); raw != "" {
		filter.Search = &raw
	}
	return filter, true
}

func visitorPendingFilterFromQuery(c *gin.Context, societyID int64) (models.VisitorPendingFilter, bool) {
	filter := models.VisitorPendingFilter{SocietyID: societyID, Limit: 50}
	if raw := c.Query("flat_id"); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || value <= 0 {
			utils.BadRequestResponse(c, "flat_id must be a positive integer")
			return filter, false
		}
		filter.FlatID = &value
	}
	if raw := c.Query("block"); raw != "" {
		filter.Block = &raw
	}
	if raw := c.Query("limit"); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 32)
		if err != nil || value <= 0 {
			utils.BadRequestResponse(c, "limit must be a positive integer")
			return filter, false
		}
		filter.Limit = int32(value)
	}
	if raw := c.Query("offset"); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 32)
		if err != nil || value < 0 {
			utils.BadRequestResponse(c, "offset must be zero or positive")
			return filter, false
		}
		filter.Offset = int32(value)
	}
	return filter, true
}

func waitingAtGateFilterFromQuery(c *gin.Context, societyID int64) (models.WaitingAtGateFilter, bool) {
	filter := models.WaitingAtGateFilter{SocietyID: societyID, Limit: 50}
	if raw := c.Query("search"); raw != "" {
		filter.Search = &raw
	}
	if raw := c.Query("limit"); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 32)
		if err != nil || value <= 0 {
			utils.BadRequestResponse(c, "limit must be a positive integer")
			return filter, false
		}
		filter.Limit = int32(value)
	}
	if raw := c.Query("offset"); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 32)
		if err != nil || value < 0 {
			utils.BadRequestResponse(c, "offset must be zero or positive")
			return filter, false
		}
		filter.Offset = int32(value)
	}
	return filter, true
}

func expectedGuestFilterFromQuery(c *gin.Context, societyID int64) (models.ExpectedGuestFilter, bool) {
	filter := models.ExpectedGuestFilter{SocietyID: societyID, Limit: 50}
	if raw := c.Query("search"); raw != "" {
		filter.Search = &raw
	}
	if raw := c.Query("event_from"); raw != "" {
		value, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			utils.BadRequestResponse(c, "event_from must be a valid RFC3339 timestamp")
			return filter, false
		}
		filter.FromAt = value
	}
	if raw := c.Query("event_to"); raw != "" {
		value, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			utils.BadRequestResponse(c, "event_to must be a valid RFC3339 timestamp")
			return filter, false
		}
		filter.ToAt = value
	}
	if raw := c.Query("limit"); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 32)
		if err != nil || value <= 0 {
			utils.BadRequestResponse(c, "limit must be a positive integer")
			return filter, false
		}
		filter.Limit = int32(value)
	}
	if raw := c.Query("offset"); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 32)
		if err != nil || value < 0 {
			utils.BadRequestResponse(c, "offset must be zero or positive")
			return filter, false
		}
		filter.Offset = int32(value)
	}
	return filter, true
}

func applyVisitorEntryDateFilters(c *gin.Context, filter *models.VisitorEntryFilter) bool {
	eventRaw := strings.TrimSpace(c.Query("event"))
	hasLegacy := strings.TrimSpace(c.Query("created_from")) != "" || strings.TrimSpace(c.Query("created_to")) != ""

	if eventRaw != "" && hasLegacy {
		utils.BadRequestResponse(c, "use either event with event_from/event_to or created_from/created_to, not both")
		return false
	}

	if eventRaw != "" {
		event := models.VisitorEntryEventFilter(eventRaw)
		if !event.IsValid() {
			utils.BadRequestResponse(c, "invalid event")
			return false
		}
		var from, to *time.Time
		if !queryTimePtr(c, "event_from", &from) || !queryTimePtr(c, "event_to", &to) {
			return false
		}
		if from == nil || to == nil {
			utils.BadRequestResponse(c, "event_from and event_to are required when event is set")
			return false
		}
		if !from.Before(*to) {
			utils.BadRequestResponse(c, "event_from must be before event_to")
			return false
		}
		filter.Event = &event
		filter.EventFrom = from
		filter.EventTo = to
		return true
	}

	if hasLegacy {
		var from, to *time.Time
		if !queryTimePtr(c, "created_from", &from) || !queryTimePtr(c, "created_to", &to) {
			return false
		}
		filter.CreatedFrom = from
		filter.CreatedTo = to
	}
	return true
}

func guardApproveOptionsFromBody(c *gin.Context) visitorentrysvc.GuardApproveOptions {
	var req models.GuardApproveEntryRequest
	_ = c.ShouldBindJSON(&req)
	opts := visitorentrysvc.GuardApproveOptions{Reason: req.Reason}
	if req.OnBehalf != nil {
		opts.OnBehalf = *req.OnBehalf
	}
	return opts
}
