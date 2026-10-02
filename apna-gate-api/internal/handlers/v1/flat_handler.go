package handlers

import (
	"net/http"
	"strings"

	"go-server/internal/models"
	flatsvc "go-server/internal/services/flatSvc"
	"go-server/pkg/utils"
	"go-server/pkg/validator"

	"github.com/gin-gonic/gin"
)

type FlatHandler struct {
	flatSvc flatsvc.FlatService
}

func NewFlatHandler(flatSvc flatsvc.FlatService) *FlatHandler {
	return &FlatHandler{flatSvc: flatSvc}
}

// CreateFlat godoc
// @Summary Create flat
// @Description [Owner/Admin/Staff] Creates a flat inside a society.
// @Tags Flats
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param request body models.CreateFlatRequest true "Create flat payload"
// @Success 201 {object} models.FlatAPIResponse "Flat created successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request or validation error"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 409 {object} models.ErrorResponseDoc "Duplicate flat"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flats [post]
func (h *FlatHandler) CreateFlat(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	var req models.CreateFlatRequest
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
	result, err := h.flatSvc.CreateFlat(c.Request.Context(), societyID, userID, &req)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Flat created successfully", gin.H{"flat": result})
}

// BulkCreateFlats godoc
// @Summary Bulk create flats
// @Description [Owner/Admin/Staff] Creates multiple flats inside a society in one request.
// @Tags Flats
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param request body models.BulkCreateFlatsRequest true "Bulk create flats payload"
// @Success 201 {object} models.BulkFlatsAPIResponse "Flats created successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request or validation error"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 409 {object} models.ErrorResponseDoc "Duplicate flat"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flats/bulk [post]
func (h *FlatHandler) BulkCreateFlats(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	var req models.BulkCreateFlatsRequest
	if !bindJSON(c, &req) {
		return
	}
	req.Sanitize()
	if err := req.Validate(); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.flatSvc.BulkCreateFlats(c.Request.Context(), societyID, userID, &req)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Flats created successfully", gin.H{"flats": result})
}

// GenerateFlats godoc
// @Summary Generate and create flats
// @Description [Owner/Admin] Generates and persists a batch of flats atomically from a continuous or floor-based numbering pattern.
// @Tags Flats
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param request body models.GenerateFlatsRequest true "Generate flats payload"
// @Success 201 {object} models.BulkFlatsAPIResponse "Flats generated and created successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request or generation configuration"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 402 {object} models.ErrorResponseDoc "Active subscription required or flat quota exceeded"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 409 {object} models.ErrorResponseDoc "Duplicate flat"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flats/generate [post]
func (h *FlatHandler) GenerateFlats(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	var req models.GenerateFlatsRequest
	if !bindJSON(c, &req) {
		return
	}
	req.Sanitize()
	if err := req.Validate(); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.flatSvc.GenerateFlats(c.Request.Context(), societyID, userID, &req)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Flats generated successfully", gin.H{"flats": result})
}

// UpdateFlat godoc
// @Summary Update flat
// @Description [Owner/Admin/Staff] Updates flat details, status, active state, or metadata.
// @Tags Flats
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flatId path int true "Flat ID"
// @Param request body models.UpdateFlatRequest true "Update flat payload"
// @Success 200 {object} models.FlatAPIResponse "Flat updated successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request or validation error"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 404 {object} models.ErrorResponseDoc "Flat not found"
// @Failure 409 {object} models.ErrorResponseDoc "Flat conflict"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flats/{flatId} [patch]
func (h *FlatHandler) UpdateFlat(c *gin.Context) {
	filter, ok := flatFilterFromPath(c)
	if !ok {
		return
	}
	var req models.UpdateFlatRequest
	if !bindJSON(c, &req) {
		return
	}
	req.Sanitize()
	if err := req.Validate(); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	result, err := h.flatSvc.UpdateFlat(c.Request.Context(), filter, &req)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Flat updated successfully", gin.H{"flat": result})
}

// DeleteFlat godoc
// @Summary Deactivate flat
// @Description [Owner/Admin/Staff] Deactivates a flat. This is not a hard delete.
// @Tags Flats
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flatId path int true "Flat ID"
// @Success 200 {object} models.MessageAPIResponse "Flat deleted successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 404 {object} models.ErrorResponseDoc "Flat not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flats/{flatId} [delete]
func (h *FlatHandler) DeleteFlat(c *gin.Context) {
	filter, ok := flatFilterFromPath(c)
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	if err := h.flatSvc.DeleteFlat(c.Request.Context(), filter, userID); handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Flat deleted successfully", gin.H{"message": "Flat deleted successfully"})
}

// BlockFlat godoc
// @Summary Block flat
// @Description [Owner/Admin/Staff] Blocks an active flat so it cannot receive claims or residents.
// @Tags Flats
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flatId path int true "Flat ID"
// @Success 200 {object} models.FlatAPIResponse "Flat blocked successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 404 {object} models.ErrorResponseDoc "Flat not found"
// @Failure 409 {object} models.ErrorResponseDoc "Invalid flat transition"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flats/{flatId}/block [post]
func (h *FlatHandler) BlockFlat(c *gin.Context) {
	filter, ok := flatFilterFromPath(c)
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.flatSvc.BlockFlat(c.Request.Context(), filter, userID, "")
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Flat blocked successfully", gin.H{"flat": result})
}

// UnblockFlat godoc
// @Summary Unblock flat
// @Description [Owner/Admin/Staff] Unblocks a blocked flat and returns it to vacant state.
// @Tags Flats
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flatId path int true "Flat ID"
// @Success 200 {object} models.FlatAPIResponse "Flat unblocked successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 404 {object} models.ErrorResponseDoc "Flat not found"
// @Failure 409 {object} models.ErrorResponseDoc "Invalid flat transition"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flats/{flatId}/unblock [post]
func (h *FlatHandler) UnblockFlat(c *gin.Context) {
	filter, ok := flatFilterFromPath(c)
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.flatSvc.UnblockFlat(c.Request.Context(), filter, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Flat unblocked successfully", gin.H{"flat": result})
}

// GetFlat godoc
// @Summary Get flat
// @Description [User/Developer] Fetches one flat with society context.
// @Tags Flats
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flatId path int true "Flat ID"
// @Success 200 {object} models.FlatAPIResponse "Flat fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter"
// @Failure 404 {object} models.ErrorResponseDoc "Flat not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/societies/{societyId}/flats/{flatId} [get]
func (h *FlatHandler) GetFlat(c *gin.Context) {
	filter, ok := flatFilterFromPath(c)
	if !ok {
		return
	}
	result, err := h.flatSvc.GetFlat(c.Request.Context(), filter)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Flat fetched successfully", gin.H{"flat": result})
}

// ListSocietyFlats godoc
// @Summary List society flats
// @Description [Owner/Admin/Staff] Lists flats inside one society with paginated metadata.
// @Tags Flats
// @Produce json
// @Param societyId path int true "Society ID"
// @Param block query string false "Block"
// @Param floor query string false "Floor"
// @Param flat_number query string false "Flat number"
// @Param status query string false "Flat status: vacant, occupied, blocked"
// @Param is_active query bool false "Active state"
// @Param search query string false "Search flat number, block, floor, or status"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} models.PaginatedFlatsAPIResponse "Flats fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid query parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flats [get]
func (h *FlatHandler) ListSocietyFlats(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	filter, ok := flatFilterFromQuery(c)
	if !ok {
		return
	}
	filter.SocietyID = &societyID
	result, err := h.flatSvc.ListFlatsPaginated(c.Request.Context(), filter)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Flats fetched successfully", gin.H{"flats": result})
}

// ListFlats godoc
// @Summary List flats
// @Description [User/Developer] Lists flats with flexible admin/developer filters.
// @Tags Flats
// @Produce json
// @Param id query int false "Flat ID"
// @Param society_id query int false "Society ID"
// @Param block query string false "Block"
// @Param floor query string false "Floor"
// @Param flat_number query string false "Flat number"
// @Param status query string false "Flat status: vacant, occupied, blocked"
// @Param is_active query bool false "Active state"
// @Param search query string false "Search flat number, block, floor, status, society name/code"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} models.FlatsAPIResponse "Flats fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid query parameter"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/flats [get]
func (h *FlatHandler) ListFlats(c *gin.Context) {
	filter, ok := flatFilterFromQuery(c)
	if !ok {
		return
	}
	result, err := h.flatSvc.ListFlats(c.Request.Context(), filter)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Flats fetched successfully", gin.H{"flats": result})
}

// GetFlatStats godoc
// @Summary Get flat stats
// @Description [Owner/Admin/Staff/Developer] Returns flat counts by status and active state for a society.
// @Tags Flats
// @Produce json
// @Param societyId path int true "Society ID"
// @Success 200 {object} models.FlatStatsAPIResponse "Flat stats fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid society ID"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/societies/{societyId}/flats/stats [get]
func (h *FlatHandler) GetFlatStats(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	result, err := h.flatSvc.GetFlatStats(c.Request.Context(), societyID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Flat stats fetched successfully", gin.H{"stats": result})
}

// GetFlatClaimStats godoc
// @Summary Get flat claim stats
// @Description [Owner/Admin/Staff] Returns flat claim counts by status for a society.
// @Tags Flat Claims
// @Produce json
// @Param societyId path int true "Society ID"
// @Success 200 {object} models.FlatClaimStatsAPIResponse "Flat claim stats fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid society ID"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/societies/{societyId}/flat-claims/stats [get]
func (h *FlatHandler) GetFlatClaimStats(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	result, err := h.flatSvc.GetFlatClaimStats(c.Request.Context(), societyID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Flat claim stats fetched successfully", gin.H{"stats": result})
}

// SubmitFlatClaim godoc
// @Summary Submit flat claim
// @Description [User] Submits a pending claim for a flat.
// @Tags Flat Claims
// @Accept json
// @Produce json
// @Param request body models.SubmitFlatClaimRequest true "Submit flat claim payload"
// @Success 201 {object} models.FlatClaimAPIResponse "Flat claim submitted successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request or validation error"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 404 {object} models.ErrorResponseDoc "Flat not found"
// @Failure 409 {object} models.ErrorResponseDoc "Duplicate pending claim or flat unavailable"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/flat-claims [post]
func (h *FlatHandler) SubmitFlatClaim(c *gin.Context) {
	var req models.SubmitFlatClaimRequest
	if !bindJSON(c, &req) {
		return
	}
	req.Sanitize()
	if err := req.Validate(); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.flatSvc.SubmitFlatClaim(c.Request.Context(), userID, &req)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Flat claim submitted successfully", gin.H{"claim": result})
}

// ApproveFlatClaim godoc
// @Summary Approve flat claim
// @Description [Owner/Admin/Staff] Transactionally approves a pending flat claim, activates society membership, creates flat resident, marks flat occupied, and returns rich joined data.
// @Tags Flat Claims
// @Produce json
// @Param societyId path int true "Society ID"
// @Param claimId path int true "Claim ID"
// @Success 200 {object} models.FlatApprovalAPIResponse "Flat claim approved successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 404 {object} models.ErrorResponseDoc "Flat claim or flat not found"
// @Failure 409 {object} models.ErrorResponseDoc "Flat unavailable, duplicate resident, or primary resident conflict"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flat-claims/{claimId}/approve [post]
func (h *FlatHandler) ApproveFlatClaim(c *gin.Context) {
	societyID, claimID, ok := parseSocietyAndClaim(c)
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.flatSvc.ApproveFlatClaim(c.Request.Context(), societyID, claimID, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Flat claim approved successfully", gin.H{"approval": result})
}

// RejectFlatClaim godoc
// @Summary Reject flat claim
// @Description [Owner/Admin/Staff] Rejects a pending flat claim with a reason.
// @Tags Flat Claims
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param claimId path int true "Claim ID"
// @Param request body models.RejectFlatClaimRequest true "Reject flat claim payload"
// @Success 200 {object} models.FlatClaimAPIResponse "Flat claim rejected successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request or path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 404 {object} models.ErrorResponseDoc "Flat claim not found"
// @Failure 409 {object} models.ErrorResponseDoc "Invalid claim transition"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flat-claims/{claimId}/reject [post]
func (h *FlatHandler) RejectFlatClaim(c *gin.Context) {
	societyID, claimID, ok := parseSocietyAndClaim(c)
	if !ok {
		return
	}
	var req models.RejectFlatClaimRequest
	if !bindJSON(c, &req) {
		return
	}
	req.Sanitize()
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.flatSvc.RejectFlatClaim(c.Request.Context(), societyID, claimID, userID, &req)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Flat claim rejected successfully", gin.H{"claim": result})
}

// CancelMyFlatClaim godoc
// @Summary Cancel my flat claim
// @Description [User] Cancels the authenticated user's pending flat claim.
// @Tags Flat Claims
// @Produce json
// @Param claimId path int true "Claim ID"
// @Success 200 {object} models.FlatClaimAPIResponse "Flat claim cancelled successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid claim ID"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 404 {object} models.ErrorResponseDoc "Flat claim not found"
// @Failure 409 {object} models.ErrorResponseDoc "Invalid claim transition"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/flat-claims/{claimId}/cancel [post]
func (h *FlatHandler) CancelMyFlatClaim(c *gin.Context) {
	claimID, ok := parsePathInt64(c, "claimId")
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.flatSvc.CancelMyFlatClaim(c.Request.Context(), claimID, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Flat claim cancelled successfully", gin.H{"claim": result})
}

// GetFlatClaim godoc
// @Summary Get flat claim
// @Description [User/Developer] Fetches one flat claim with joined user, flat, and society data.
// @Tags Flat Claims
// @Produce json
// @Param claimId path int true "Claim ID"
// @Success 200 {object} models.FlatClaimAPIResponse "Flat claim fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid claim ID"
// @Failure 404 {object} models.ErrorResponseDoc "Flat claim not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/flat-claims/{claimId} [get]
func (h *FlatHandler) GetFlatClaim(c *gin.Context) {
	claimID, ok := parsePathInt64(c, "claimId")
	if !ok {
		return
	}
	result, err := h.flatSvc.GetFlatClaim(c.Request.Context(), &models.FlatClaimFilter{ID: &claimID})
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Flat claim fetched successfully", gin.H{"claim": result})
}

// GetSocietyFlatClaim godoc
// @Summary Get society flat claim
// @Description [Owner/Admin/Staff] Fetches one flat claim inside a society with joined user, flat, and society data.
// @Tags Flat Claims
// @Produce json
// @Param societyId path int true "Society ID"
// @Param claimId path int true "Claim ID"
// @Success 200 {object} models.FlatClaimAPIResponse "Flat claim fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 404 {object} models.ErrorResponseDoc "Flat claim not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flat-claims/{claimId} [get]
func (h *FlatHandler) GetSocietyFlatClaim(c *gin.Context) {
	societyID, claimID, ok := parseSocietyAndClaim(c)
	if !ok {
		return
	}
	result, err := h.flatSvc.GetFlatClaim(c.Request.Context(), &models.FlatClaimFilter{ID: &claimID, SocietyID: &societyID})
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Flat claim fetched successfully", gin.H{"claim": result})
}

// ListFlatClaims godoc
// @Summary List flat claims
// @Description [Owner/Admin/Staff/Developer] Lists flat claims with flexible filters and rich joined data.
// @Tags Flat Claims
// @Produce json
// @Param id query int false "Claim ID"
// @Param society_id query int false "Society ID"
// @Param flat_id query int false "Flat ID"
// @Param user_id query int false "User ID"
// @Param status query string false "Claim status: pending, approved, rejected, cancelled"
// @Param search query string false "Search user, contact, flat, block, or status"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} models.FlatClaimsAPIResponse "Flat claims fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid query parameter"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/flat-claims [get]
func (h *FlatHandler) ListFlatClaims(c *gin.Context) {
	filter, ok := flatClaimFilterFromQuery(c)
	if !ok {
		return
	}
	result, err := h.flatSvc.ListFlatClaims(c.Request.Context(), filter)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Flat claims fetched successfully", gin.H{"claims": result})
}

// ListSocietyFlatClaims godoc
// @Summary List society flat claims
// @Description [Owner/Admin/Staff] Lists flat claims inside one society with flexible filters and rich joined data.
// @Tags Flat Claims
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flat_id query int false "Flat ID"
// @Param user_id query int false "User ID"
// @Param status query string false "Claim status: pending, approved, rejected, cancelled"
// @Param search query string false "Search user, contact, flat, block, or status"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} models.FlatClaimsAPIResponse "Flat claims fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid query parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flat-claims [get]
func (h *FlatHandler) ListSocietyFlatClaims(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	filter, ok := flatClaimFilterFromQuery(c)
	if !ok {
		return
	}
	filter.SocietyID = &societyID
	result, err := h.flatSvc.ListFlatClaims(c.Request.Context(), filter)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Flat claims fetched successfully", gin.H{"claims": result})
}

// ListMyFlatClaims godoc
// @Summary List my flat claims
// @Description [User] Lists flat claims submitted by the authenticated user.
// @Tags Flat Claims
// @Produce json
// @Param society_id query int false "Society ID"
// @Param flat_id query int false "Flat ID"
// @Param status query string false "Claim status"
// @Param search query string false "Search text"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} models.FlatClaimsAPIResponse "My flat claims fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid query parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/me/flat-claims [get]
func (h *FlatHandler) ListMyFlatClaims(c *gin.Context) {
	filter, ok := flatClaimFilterFromQuery(c)
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.flatSvc.ListMyFlatClaims(c.Request.Context(), userID, filter)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "My flat claims fetched successfully", gin.H{"claims": result})
}

// AddResidentToFlat godoc
// @Summary Add resident to flat
// @Description [Owner/Admin/Staff] Manually adds a user as an active resident of a flat.
// @Tags Flat Residents
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flatId path int true "Flat ID"
// @Param userId path int true "User ID"
// @Param request body models.AddFlatResidentRequest true "Add flat resident payload"
// @Success 201 {object} models.FlatResidentAPIResponse "Resident added successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request or path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 404 {object} models.ErrorResponseDoc "Flat or user not found"
// @Failure 409 {object} models.ErrorResponseDoc "Duplicate resident, blocked flat, or primary resident conflict"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flats/{flatId}/residents/users/{userId} [post]
func (h *FlatHandler) AddResidentToFlat(c *gin.Context) {
	societyID, flatID, userID, ok := parseSocietyFlatUser(c)
	if !ok {
		return
	}
	var req models.AddFlatResidentRequest
	if !bindJSON(c, &req) {
		return
	}
	if err := req.Validate(); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}
	actorID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.flatSvc.AddResidentToFlat(c.Request.Context(), societyID, flatID, userID, actorID, &req)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Resident added successfully", gin.H{"resident": result})
}

// RemoveResidentFromFlat godoc
// @Summary Remove resident from flat
// @Description [Owner/Admin/Staff] Marks an active resident inactive and marks the flat vacant if this was the last active resident.
// @Tags Flat Residents
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flatId path int true "Flat ID"
// @Param residentId path int true "Resident ID"
// @Success 200 {object} models.MessageAPIResponse "Resident removed successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 404 {object} models.ErrorResponseDoc "Resident not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flats/{flatId}/residents/{residentId} [delete]
func (h *FlatHandler) RemoveResidentFromFlat(c *gin.Context) {
	filter, ok := flatResidentFilterFromPath(c)
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	if err := h.flatSvc.RemoveResidentFromFlat(c.Request.Context(), filter, userID); handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Resident removed successfully", gin.H{"message": "Resident removed successfully"})
}

// ChangePrimaryResident godoc
// @Summary Change primary resident
// @Description [Owner/Admin/Staff] Transactionally changes the active primary resident for a flat.
// @Tags Flat Residents
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flatId path int true "Flat ID"
// @Param residentId path int true "Resident ID"
// @Success 200 {object} models.FlatResidentAPIResponse "Primary resident changed successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 404 {object} models.ErrorResponseDoc "Resident not found"
// @Failure 409 {object} models.ErrorResponseDoc "Primary resident conflict"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flats/{flatId}/residents/{residentId}/primary [post]
func (h *FlatHandler) ChangePrimaryResident(c *gin.Context) {
	societyID, flatID, residentID, ok := parseSocietyFlatResident(c)
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.flatSvc.ChangePrimaryResident(c.Request.Context(), societyID, flatID, residentID, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Primary resident changed successfully", gin.H{"resident": result})
}

// UpdateFlatResidentRole godoc
// @Summary Update flat resident role
// @Description [Owner/Admin/Staff] Updates an active flat resident role.
// @Tags Flat Residents
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flatId path int true "Flat ID"
// @Param residentId path int true "Resident ID"
// @Param request body models.UpdateFlatResidentRoleRequest true "Update flat resident role payload"
// @Success 200 {object} models.FlatResidentAPIResponse "Resident role updated successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request or path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 404 {object} models.ErrorResponseDoc "Resident not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flats/{flatId}/residents/{residentId}/role [patch]
func (h *FlatHandler) UpdateFlatResidentRole(c *gin.Context) {
	filter, ok := flatResidentFilterFromPath(c)
	if !ok {
		return
	}
	var req models.UpdateFlatResidentRoleRequest
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
	result, err := h.flatSvc.UpdateFlatResidentRole(c.Request.Context(), filter, userID, &req)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Resident role updated successfully", gin.H{"resident": result})
}

// MoveOutResident godoc
// @Summary Move out resident
// @Description [Owner/Admin/Staff] Marks a resident moved out and marks the flat vacant if this was the last active resident.
// @Tags Flat Residents
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flatId path int true "Flat ID"
// @Param residentId path int true "Resident ID"
// @Success 200 {object} models.FlatResidentAPIResponse "Resident moved out successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 404 {object} models.ErrorResponseDoc "Resident not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flats/{flatId}/residents/{residentId}/move-out [post]
func (h *FlatHandler) MoveOutResident(c *gin.Context) {
	filter, ok := flatResidentFilterFromPath(c)
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.flatSvc.MoveOutResident(c.Request.Context(), filter, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Resident moved out successfully", gin.H{"resident": result})
}

// GetFlatResident godoc
// @Summary Get flat resident
// @Description [User/Developer] Fetches one flat resident with joined user, flat, and society data.
// @Tags Flat Residents
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flatId path int true "Flat ID"
// @Param residentId path int true "Resident ID"
// @Success 200 {object} models.FlatResidentAPIResponse "Resident fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter"
// @Failure 404 {object} models.ErrorResponseDoc "Resident not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/societies/{societyId}/flats/{flatId}/residents/{residentId} [get]
func (h *FlatHandler) GetFlatResident(c *gin.Context) {
	filter, ok := flatResidentFilterFromPath(c)
	if !ok {
		return
	}
	result, err := h.flatSvc.GetFlatResident(c.Request.Context(), filter)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Resident fetched successfully", gin.H{"resident": result})
}

// ListFlatResidents godoc
// @Summary List flat residents
// @Description [Owner/Admin/Staff/Developer] Lists flat residents with flexible filters and joined user, flat, and society data.
// @Tags Flat Residents
// @Produce json
// @Param id query int false "Resident ID"
// @Param society_id query int false "Society ID"
// @Param flat_id query int false "Flat ID"
// @Param user_id query int false "User ID"
// @Param role query string false "Resident role: owner, tenant, family"
// @Param status query string false "Resident status: active, inactive, moved_out"
// @Param is_primary query bool false "Primary resident flag"
// @Param search query string false "Search user, contact, flat, block, role, or status"
// @Param search_mode query string false "Search mode"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} models.FlatResidentsAPIResponse "Residents fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid query parameter"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/flat-residents [get]
func (h *FlatHandler) ListFlatResidents(c *gin.Context) {
	filter, ok := flatResidentFilterFromQuery(c)
	if !ok {
		return
	}
	result, err := h.flatSvc.ListFlatResidents(c.Request.Context(), filter)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Residents fetched successfully", gin.H{"residents": result})
}

// ListSocietyFlatResidents godoc
// @Summary List flat residents for a society flat
// @Description [Owner/Admin/Staff] Lists residents for the flat identified in the path.
// @Tags Flat Residents
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flatId path int true "Flat ID"
// @Param role query string false "Resident role: owner, tenant, family"
// @Param status query string false "Resident status: active, inactive, moved_out"
// @Param is_primary query bool false "Primary resident flag"
// @Param search query string false "Search user, contact, flat, block, role, or status"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} models.FlatResidentsAPIResponse "Residents fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path or query parameter"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/societies/{societyId}/flats/{flatId}/residents [get]
func (h *FlatHandler) ListSocietyFlatResidents(c *gin.Context) {
	filter, ok := flatResidentFilterFromSocietyFlatPath(c)
	if !ok {
		return
	}
	result, err := h.flatSvc.ListFlatResidents(c.Request.Context(), filter)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Residents fetched successfully", gin.H{"residents": result})
}

// ListMyResidences godoc
// @Summary List my residences
// @Description [User] Lists flat residences for the authenticated user with joined flat and society data.
// @Tags Flat Residents
// @Produce json
// @Param society_id query int false "Society ID"
// @Param flat_id query int false "Flat ID"
// @Param role query string false "Resident role"
// @Param status query string false "Resident status"
// @Param is_primary query bool false "Primary resident flag"
// @Param search query string false "Search text"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} models.MyResidencesAPIResponse "My residences fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid query parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/me/residences [get]
func (h *FlatHandler) ListMyResidences(c *gin.Context) {
	filter, ok := flatResidentFilterFromQuery(c)
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.flatSvc.ListMyResidences(c.Request.Context(), userID, filter)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "My residences fetched successfully", gin.H{"residences": result})
}

func flatFilterFromPath(c *gin.Context) (*models.FlatFilter, bool) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return nil, false
	}
	flatID, ok := parsePathInt64(c, "flatId")
	if !ok {
		return nil, false
	}
	return &models.FlatFilter{ID: &flatID, SocietyID: &societyID}, true
}

func flatFilterFromQuery(c *gin.Context) (*models.FlatFilter, bool) {
	filter := &models.FlatFilter{
		Block: optionalString(c.Query("block")), Floor: optionalString(c.Query("floor")),
		FlatNumber: optionalString(c.Query("flat_number")), Status: optionalString(c.Query("status")),
		Search: strings.TrimSpace(c.Query("search")),
	}
	if !queryInt64Ptr(c, "id", &filter.ID) || !queryInt64Ptr(c, "society_id", &filter.SocietyID) {
		return nil, false
	}
	if raw := strings.TrimSpace(c.Query("is_active")); raw != "" {
		value := raw == "true" || raw == "1"
		filter.IsActive = &value
	}
	limit, offset, ok := paginationQuery(c)
	if !ok {
		return nil, false
	}
	filter.Limit = limit
	filter.Offset = offset
	return filter, true
}

func flatClaimFilterFromQuery(c *gin.Context) (*models.FlatClaimFilter, bool) {
	filter := &models.FlatClaimFilter{Status: optionalString(c.Query("status")), Search: strings.TrimSpace(c.Query("search")), SearchMode: strings.TrimSpace(c.Query("search_mode"))}
	if !queryInt64Ptr(c, "id", &filter.ID) || !queryInt64Ptr(c, "society_id", &filter.SocietyID) ||
		!queryInt64Ptr(c, "flat_id", &filter.FlatID) || !queryInt64Ptr(c, "user_id", &filter.UserID) {
		return nil, false
	}
	limit, offset, ok := paginationQuery(c)
	if !ok {
		return nil, false
	}
	filter.Limit = limit
	filter.Offset = offset
	return filter, true
}

func flatResidentFilterFromQuery(c *gin.Context) (*models.FlatResidentFilter, bool) {
	filter := &models.FlatResidentFilter{Role: optionalString(c.Query("role")), Status: optionalString(c.Query("status")), Search: strings.TrimSpace(c.Query("search")), SearchMode: strings.TrimSpace(c.Query("search_mode"))}
	if !queryInt64Ptr(c, "id", &filter.ID) || !queryInt64Ptr(c, "society_id", &filter.SocietyID) ||
		!queryInt64Ptr(c, "flat_id", &filter.FlatID) || !queryInt64Ptr(c, "user_id", &filter.UserID) {
		return nil, false
	}
	if raw := strings.TrimSpace(c.Query("is_primary")); raw != "" {
		value := raw == "true" || raw == "1"
		filter.IsPrimary = &value
	}
	limit, offset, ok := paginationQuery(c)
	if !ok {
		return nil, false
	}
	filter.Limit = limit
	filter.Offset = offset
	return filter, true
}

func flatResidentFilterFromPath(c *gin.Context) (*models.FlatResidentFilter, bool) {
	societyID, flatID, residentID, ok := parseSocietyFlatResident(c)
	if !ok {
		return nil, false
	}
	return &models.FlatResidentFilter{ID: &residentID, SocietyID: &societyID, FlatID: &flatID}, true
}

func parseSocietyAndClaim(c *gin.Context) (int64, int64, bool) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return 0, 0, false
	}
	claimID, ok := parsePathInt64(c, "claimId")
	return societyID, claimID, ok
}

func parseSocietyFlatUser(c *gin.Context) (int64, int64, int64, bool) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return 0, 0, 0, false
	}
	flatID, ok := parsePathInt64(c, "flatId")
	if !ok {
		return 0, 0, 0, false
	}
	userID, ok := parsePathInt64(c, "userId")
	return societyID, flatID, userID, ok
}

func parseSocietyFlatResident(c *gin.Context) (int64, int64, int64, bool) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return 0, 0, 0, false
	}
	flatID, ok := parsePathInt64(c, "flatId")
	if !ok {
		return 0, 0, 0, false
	}
	residentID, ok := parsePathInt64(c, "residentId")
	return societyID, flatID, residentID, ok
}
func flatResidentFilterFromSocietyFlatPath(c *gin.Context) (*models.FlatResidentFilter, bool) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return nil, false
	}
	flatID, ok := parsePathInt64(c, "flatId")
	if !ok {
		return nil, false
	}
	filter := &models.FlatResidentFilter{
		SocietyID:  &societyID,
		FlatID:     &flatID,
		Role:       optionalString(c.Query("role")),
		Status:     optionalString(c.Query("status")),
		Search:     strings.TrimSpace(c.Query("search")),
		SearchMode: strings.TrimSpace(c.Query("search_mode")),
	}
	if raw := strings.TrimSpace(c.Query("is_primary")); raw != "" {
		value := raw == "true" || raw == "1"
		filter.IsPrimary = &value
	}
	limit, offset, ok := paginationQuery(c)
	if !ok {
		return nil, false
	}
	filter.Limit = limit
	filter.Offset = offset
	return filter, true
}
