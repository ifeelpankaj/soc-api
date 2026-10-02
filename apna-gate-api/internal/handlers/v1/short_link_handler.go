package handlers

import (
	"net/http"
	"strings"
	"time"

	"go-server/internal/models"
	flatsvc "go-server/internal/services/flatSvc"
	shortlinksvc "go-server/internal/services/shortLinkSvc"
	visitorentrysvc "go-server/internal/services/visitorEntrySvc"
	"go-server/pkg/utils"
	"go-server/pkg/validator"

	"github.com/gin-gonic/gin"
)

type ShortLinkHandler struct {
	shortLinks shortlinksvc.ShortLinkService
	inviteSvc  visitorentrysvc.VisitorInviteService
	flatSvc    flatsvc.FlatService
}

func NewShortLinkHandler(shortLinks shortlinksvc.ShortLinkService, inviteSvc visitorentrysvc.VisitorInviteService, flatSvc flatsvc.FlatService) *ShortLinkHandler {
	return &ShortLinkHandler{shortLinks: shortLinks, inviteSvc: inviteSvc, flatSvc: flatSvc}
}

// Resolve godoc
// @Summary Resolve short link
// @Description [Public] Resolves a reusable Apna Gate short link into the current invite workflow state.
// @Tags Short Links
// @Produce json
// @Param shortCode path string true "Short link code"
// @Success 200 {object} models.ShortLinkResolveAPIResponse "Short link resolved successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid short code"
// @Failure 404 {object} models.ErrorResponseDoc "Short link not found"
// @Failure 409 {object} models.ErrorResponseDoc "Short link unavailable"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/public/short-links/{shortCode} [get]
func (h *ShortLinkHandler) Resolve(c *gin.Context) {
	link, err := h.shortLinks.GetByCode(c.Request.Context(), c.Param("shortCode"))
	if handleServiceError(c, err) {
		return
	}
	response, err := h.resolveResponse(c, link)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Short link resolved successfully", gin.H{"link": response})
}

// SubmitVisitorInvite godoc
// @Summary Submit visitor invite via short link
// @Description [Public] Submits visitor details for an active visitor invite resolved by short code.
// @Tags Short Links
// @Accept json
// @Produce json
// @Param shortCode path string true "Short link code"
// @Param request body models.VisitorFormRequest true "Visitor details"
// @Success 201 {object} models.VisitorEntryMutationAPIResponse "Visitor entry created successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request or short code"
// @Failure 404 {object} models.ErrorResponseDoc "Short link or invite not found"
// @Failure 409 {object} models.ErrorResponseDoc "Visitor invite unavailable"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/public/short-links/{shortCode}/visitor-submit [post]
func (h *ShortLinkHandler) SubmitVisitorInvite(c *gin.Context) {
	link, ok := h.resolveTypedLink(c, models.ShortLinkResourceVisitorInvite)
	if !ok {
		return
	}
	var req models.VisitorFormRequest
	if !bindJSON(c, &req) {
		return
	}
	result, err := h.inviteSvc.SubmitInviteFormByID(c.Request.Context(), link.ResourceID, req)
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

// JoinMemberInvite godoc
// @Summary Join member invite via short link
// @Description [Public] Creates or authenticates a user and accepts a flat member invite resolved by short code.
// @Tags Short Links
// @Accept json
// @Produce json
// @Param shortCode path string true "Short link code"
// @Param request body models.JoinFlatMemberInviteRequest true "Join request"
// @Success 200 {object} models.JoinFlatMemberInviteAPIResponse "Member joined successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request or short code"
// @Failure 404 {object} models.ErrorResponseDoc "Short link or member invite not found"
// @Failure 409 {object} models.ErrorResponseDoc "Member invite unavailable or resident conflict"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/public/short-links/{shortCode}/member-join [post]
func (h *ShortLinkHandler) JoinMemberInvite(c *gin.Context) {
	link, ok := h.resolveTypedLink(c, models.ShortLinkResourceMemberInvite)
	if !ok {
		return
	}
	var req models.JoinFlatMemberInviteRequest
	if !bindJSON(c, &req) {
		return
	}
	req.Sanitize()
	if validationErrors := validator.ValidateStruct(&req); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(c, validationErrors.ToMap())
		return
	}
	if !req.IsRegisterFlow() && strings.TrimSpace(req.Identifier) == "" {
		utils.BadRequestResponse(c, "first_name or identifier is required")
		return
	}
	result, err := h.flatSvc.JoinMemberInviteByID(c.Request.Context(), link.ResourceID, &req)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Member joined successfully. Log in on the mobile app with your credentials.", gin.H{"join": result})
}

// AcceptMemberInvite godoc
// @Summary Accept member invite via short link
// @Description [User] Accepts a flat member invite resolved by short code.
// @Tags Short Links
// @Produce json
// @Param shortCode path string true "Short link code"
// @Success 200 {object} models.AcceptFlatMemberInviteAPIResponse "Member invite accepted successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid short code"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 404 {object} models.ErrorResponseDoc "Short link or member invite not found"
// @Failure 409 {object} models.ErrorResponseDoc "Member invite unavailable or resident conflict"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/public/short-links/{shortCode}/member-accept [post]
func (h *ShortLinkHandler) AcceptMemberInvite(c *gin.Context) {
	link, ok := h.resolveTypedLink(c, models.ShortLinkResourceMemberInvite)
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.flatSvc.AcceptMemberInviteByID(c.Request.Context(), link.ResourceID, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Member invite accepted successfully", gin.H{"acceptance": result})
}

func (h *ShortLinkHandler) resolveTypedLink(c *gin.Context, expected models.ShortLinkResourceType) (*models.ShortLink, bool) {
	link, err := h.shortLinks.GetActiveByCode(c.Request.Context(), c.Param("shortCode"))
	if handleServiceError(c, err) {
		return nil, false
	}
	if link.ResourceType != expected {
		utils.BadRequestResponse(c, "short link resource type mismatch")
		return nil, false
	}
	return link, true
}

func (h *ShortLinkHandler) resolveResponse(c *gin.Context, link *models.ShortLink) (*models.ShortLinkResolveResponse, error) {
	if link.RevokedAt != nil {
		return &models.ShortLinkResolveResponse{
			Type:  link.ResourceType,
			State: models.ShortLinkStateUnavailable,
			Link:  h.shortLinks.Response(link),
		}, nil
	}
	switch link.ResourceType {
	case models.ShortLinkResourceVisitorInvite:
		page, err := h.inviteSvc.GetPublicInviteByID(c.Request.Context(), link.ResourceID)
		if err != nil {
			return nil, err
		}
		return visitorShortLinkResponse(link, h.shortLinks.Response(link), page), nil
	case models.ShortLinkResourceMemberInvite:
		invite, err := h.flatSvc.GetPublicMemberInviteByID(c.Request.Context(), link.ResourceID)
		if err != nil {
			return nil, err
		}
		return memberShortLinkResponse(link, h.shortLinks.Response(link), invite), nil
	default:
		return nil, shortlinksvc.ErrInvalidShortLink
	}
}

func visitorShortLinkResponse(link *models.ShortLink, linkResponse *models.ShortLinkResponse, page *models.PublicVisitorInvitePageResponse) *models.ShortLinkResolveResponse {
	state := models.ShortLinkStateUnavailable
	if page != nil {
		switch page.View {
		case models.PublicVisitorInviteViewForm:
			state = models.ShortLinkStateForm
		case models.PublicVisitorInviteViewQR:
			state = models.ShortLinkStateApproved
		case models.PublicVisitorInviteViewCheckedIn:
			state = models.ShortLinkStateCheckedIn
		case models.PublicVisitorInviteViewCheckedOut:
			state = models.ShortLinkStateCompleted
		default:
			if page.Entry != nil && page.Entry.Status == models.VisitorStatusWaitingApproval {
				state = models.ShortLinkStateWaitingApproval
			}
		}
		if page.Invite != nil {
			switch page.Invite.Status {
			case models.VisitorInviteStatusExpired:
				state = models.ShortLinkStateExpired
			case models.VisitorInviteStatusCancelled:
				state = models.ShortLinkStateCancelled
			}
		}
	}
	if link.ExpiresAt != nil && !link.ExpiresAt.After(time.Now()) {
		state = models.ShortLinkStateExpired
	}
	return &models.ShortLinkResolveResponse{
		Type:   link.ResourceType,
		State:  state,
		Invite: page.Invite,
		Entry:  page.Entry,
		QR:     page.QR,
		Link:   linkResponse,
	}
}

func memberShortLinkResponse(link *models.ShortLink, linkResponse *models.ShortLinkResponse, invite *models.PublicFlatMemberInviteView) *models.ShortLinkResolveResponse {
	state := models.ShortLinkStateUnavailable
	if invite != nil {
		switch invite.Status {
		case models.FlatMemberInviteStatusPending:
			if invite.ExpiresAt.After(time.Now()) {
				state = models.ShortLinkStateJoin
			} else {
				state = models.ShortLinkStateExpired
			}
		case models.FlatMemberInviteStatusAccepted:
			state = models.ShortLinkStateAccepted
		case models.FlatMemberInviteStatusExpired:
			state = models.ShortLinkStateExpired
		case models.FlatMemberInviteStatusCancelled:
			state = models.ShortLinkStateCancelled
		}
	}
	if link.ExpiresAt != nil && !link.ExpiresAt.After(time.Now()) {
		state = models.ShortLinkStateExpired
	}
	return &models.ShortLinkResolveResponse{
		Type:   link.ResourceType,
		State:  state,
		Invite: invite,
		Link:   linkResponse,
	}
}
