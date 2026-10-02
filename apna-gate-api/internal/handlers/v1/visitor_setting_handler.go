package handlers

import (
	"net/http"
	"strconv"

	"go-server/internal/models"
	visitorsettingsvc "go-server/internal/services/visitorSettingSvc"
	"go-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

type VisitorSettingHandler struct {
	visitorSettingSvc visitorsettingsvc.VisitorSettingService
}

func NewVisitorSettingHandler(visitorSettingSvc visitorsettingsvc.VisitorSettingService) *VisitorSettingHandler {
	return &VisitorSettingHandler{visitorSettingSvc: visitorSettingSvc}
}

// GetSocietySettings godoc
// @Summary Get society visitor settings
// @Description [Owner/Admin] Fetches society-level visitor configuration for approval mode, QR entry, guard entry, pre-approval, durations, and active state.
// @Tags Visitor Settings
// @Produce json
// @Param societyId path int true "Society ID"
// @Success 200 {object} models.SocietyVisitorSettingsAPIResponse "Visitor settings fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid society ID"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 404 {object} models.ErrorResponseDoc "Visitor settings not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/visitor-settings [get]
func (h *VisitorSettingHandler) GetSocietySettings(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	result, err := h.visitorSettingSvc.GetSocietySettings(c.Request.Context(), societyID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Visitor settings fetched successfully", gin.H{"visitor_settings": result})
}

// UpdateSocietySettings godoc
// @Summary Update society visitor settings
// @Description [Owner/Admin] Updates society-level visitor configuration. Boolean false values are accepted and persisted.
// @Tags Visitor Settings
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param request body models.UpdateSocietyVisitorSettingsRequest true "Society visitor settings update payload"
// @Success 200 {object} models.SocietyVisitorSettingsAPIResponse "Visitor settings updated successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request, validation error, or invalid society ID"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 404 {object} models.ErrorResponseDoc "Visitor settings not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/visitor-settings [patch]
func (h *VisitorSettingHandler) UpdateSocietySettings(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	var req models.UpdateSocietyVisitorSettingsRequest
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
	result, err := h.visitorSettingSvc.UpdateSocietySettings(c.Request.Context(), societyID, req, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Visitor settings updated successfully", gin.H{"visitor_settings": result})
}

// GetFlatSettings godoc
// @Summary Get flat visitor settings
// @Description [Owner/Admin/Flat Resident] Fetches all visitor purpose settings for a flat.
// @Tags Visitor Settings
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flatId path int true "Flat ID"
// @Success 200 {object} models.FlatVisitorSettingsAPIResponse "Flat visitor settings fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Owner/admin or flat resident access required"
// @Failure 404 {object} models.ErrorResponseDoc "Flat or visitor settings not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flats/{flatId}/visitor-settings [get]
func (h *VisitorSettingHandler) GetFlatSettings(c *gin.Context) {
	societyID, flatID, ok := visitorSettingsPath(c)
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.visitorSettingSvc.GetFlatSettingsForActor(c.Request.Context(), societyID, flatID, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Flat visitor settings fetched successfully", gin.H{"visitor_settings": result})
}

// UpdateFlatPurposeSetting godoc
// @Summary Update flat visitor purpose setting
// @Description [Owner/Admin/Flat Resident] Updates approval, duration override, or enabled state for one visitor purpose on a flat.
// @Tags Visitor Settings
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flatId path int true "Flat ID"
// @Param purpose path string true "Visitor purpose" Enums(guest, delivery, cab, service, maintenance, staff, other)
// @Param request body models.UpdateFlatVisitorSettingRequest true "Flat visitor purpose setting update payload"
// @Success 200 {object} models.FlatVisitorSettingAPIResponse "Flat visitor setting updated successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request, visitor purpose, or path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Owner/admin or flat resident access required"
// @Failure 404 {object} models.ErrorResponseDoc "Flat or visitor setting not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flats/{flatId}/visitor-settings/{purpose} [patch]
func (h *VisitorSettingHandler) UpdateFlatPurposeSetting(c *gin.Context) {
	societyID, flatID, ok := visitorSettingsPath(c)
	if !ok {
		return
	}
	purpose := models.VisitorPurpose(c.Param("purpose"))
	if !purpose.IsValid() {
		utils.BadRequestResponse(c, "invalid visitor purpose")
		return
	}
	var req models.UpdateFlatVisitorSettingRequest
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
	result, err := h.visitorSettingSvc.UpdateFlatPurposeSetting(c.Request.Context(), societyID, flatID, purpose, req, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Flat visitor setting updated successfully", gin.H{"visitor_setting": result})
}

// ResetFlatSettingsToDefault godoc
// @Summary Reset flat visitor settings
// @Description [Owner/Admin/Flat Resident] Resets a flat's visitor purpose settings back to the production defaults.
// @Tags Visitor Settings
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flatId path int true "Flat ID"
// @Success 200 {object} models.MessageAPIResponse "Flat visitor settings reset successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Owner/admin or flat resident access required"
// @Failure 404 {object} models.ErrorResponseDoc "Flat or visitor settings not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/flats/{flatId}/visitor-settings/reset [post]
func (h *VisitorSettingHandler) ResetFlatSettingsToDefault(c *gin.Context) {
	societyID, flatID, ok := visitorSettingsPath(c)
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	if err := h.visitorSettingSvc.ResetFlatSettingsToDefault(c.Request.Context(), societyID, flatID, userID); handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Flat visitor settings reset successfully", gin.H{"message": "Flat visitor settings reset successfully"})
}

// ListSocietyFlatSettings godoc
// @Summary List society flat visitor settings
// @Description [Owner/Admin] Lists flat-level visitor purpose settings across the society with pagination.
// @Tags Visitor Settings
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flat_id query int false "Flat ID"
// @Param block query string false "Block"
// @Param purpose query string false "Visitor purpose" Enums(guest, delivery, cab, service, maintenance, staff, other)
// @Param limit query int false "Maximum records to return" default(50)
// @Param offset query int false "Records to skip" default(0)
// @Success 200 {object} models.SocietyFlatVisitorSettingsAPIResponse "Society flat visitor settings fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid query or path parameter"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 403 {object} models.ErrorResponseDoc "Insufficient society role"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/societies/{societyId}/visitor-settings/flats [get]
func (h *VisitorSettingHandler) ListSocietyFlatSettings(c *gin.Context) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return
	}
	filter, ok := societyFlatVisitorSettingsFilterFromQuery(c, societyID)
	if !ok {
		return
	}
	result, err := h.visitorSettingSvc.ListSocietyFlatSettings(c.Request.Context(), filter)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Society flat visitor settings fetched successfully", gin.H{
		"settings": result.Settings,
		"total":    result.Total,
		"limit":    result.Limit,
		"offset":   result.Offset,
	})
}

func visitorSettingsPath(c *gin.Context) (int64, int64, bool) {
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

func societyFlatVisitorSettingsFilterFromQuery(c *gin.Context, societyID int64) (models.SocietyFlatVisitorSettingsFilter, bool) {
	filter := models.SocietyFlatVisitorSettingsFilter{SocietyID: societyID, Limit: 50}
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
	return filter, true
}
