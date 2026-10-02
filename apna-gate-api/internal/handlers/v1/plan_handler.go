package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"go-server/internal/models"
	plansvc "go-server/internal/services/planSvc"
	"go-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

type PlanHandler struct {
	planSvc plansvc.PlanService
}

func NewPlanHandler(planSvc plansvc.PlanService) *PlanHandler {
	return &PlanHandler{planSvc: planSvc}
}

// CreatePlan godoc
// @Summary Create plan
// @Description [Developer] Creates a subscription plan.
// @Tags Plans
// @Accept json
// @Produce json
// @Param request body models.CreatePlanRequest true "Create plan payload"
// @Success 201 {object} models.PlanAPIResponse "Plan created successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request"
// @Failure 409 {object} models.ErrorResponseDoc "Duplicate plan"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/plans [post]
func (h *PlanHandler) CreatePlan(c *gin.Context) {
	var req models.CreatePlanRequest
	if !bindJSON(c, &req) {
		return
	}
	result, err := h.planSvc.CreatePlan(c.Request.Context(), &req)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Plan created successfully", gin.H{"plan": result})
}

// UpdatePlan godoc
// @Summary Update plan
// @Description [Developer] Updates a subscription plan.
// @Tags Plans
// @Accept json
// @Produce json
// @Param planId path int true "Plan ID"
// @Param request body models.UpdatePlanRequest true "Update plan payload"
// @Success 200 {object} models.PlanAPIResponse "Plan updated successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid request"
// @Failure 404 {object} models.ErrorResponseDoc "Plan not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/plans/{planId} [patch]
func (h *PlanHandler) UpdatePlan(c *gin.Context) {
	planID, ok := parsePathInt64(c, "planId")
	if !ok {
		return
	}
	var req models.UpdatePlanRequest
	if !bindJSON(c, &req) {
		return
	}
	result, err := h.planSvc.UpdatePlan(c.Request.Context(), planID, &req)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Plan updated successfully", gin.H{"plan": result})
}

// ActivatePlan godoc
// @Summary Activate plan
// @Description [Developer] Activates a subscription plan.
// @Tags Plans
// @Produce json
// @Param planId path int true "Plan ID"
// @Success 200 {object} models.PlanAPIResponse "Plan activated successfully"
// @Failure 404 {object} models.ErrorResponseDoc "Plan not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/plans/{planId}/activate [post]
func (h *PlanHandler) ActivatePlan(c *gin.Context) {
	planID, ok := parsePathInt64(c, "planId")
	if !ok {
		return
	}
	result, err := h.planSvc.ActivatePlan(c.Request.Context(), planID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Plan activated successfully", gin.H{"plan": result})
}

// DeactivatePlan godoc
// @Summary Deactivate plan
// @Description [Developer] Deactivates a subscription plan.
// @Tags Plans
// @Produce json
// @Param planId path int true "Plan ID"
// @Success 200 {object} models.PlanAPIResponse "Plan deactivated successfully"
// @Failure 404 {object} models.ErrorResponseDoc "Plan not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Router /v1/plans/{planId}/deactivate [post]
func (h *PlanHandler) DeactivatePlan(c *gin.Context) {
	planID, ok := parsePathInt64(c, "planId")
	if !ok {
		return
	}
	result, err := h.planSvc.DeactivatePlan(c.Request.Context(), planID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Plan deactivated successfully", gin.H{"plan": result})
}

// GetPlan godoc
// @Summary Get plan
// @Description [Developer/User] Fetches a plan by flexible filter.
// @Tags Plans
// @Produce json
// @Param id query int false "Plan ID"
// @Param code query string false "Plan code"
// @Param name query string false "Plan name"
// @Success 200 {object} models.PlanAPIResponse "Plan fetched successfully"
// @Failure 404 {object} models.ErrorResponseDoc "Plan not found"
// @Router /v1/plans/lookup [get]
func (h *PlanHandler) GetPlan(c *gin.Context) {
	filter, ok := planFilterFromQuery(c)
	if !ok {
		return
	}
	result, err := h.planSvc.GetPlan(c.Request.Context(), filter)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Plan fetched successfully", gin.H{"plan": result})
}

// ListPlans godoc
// @Summary List plans
// @Description [Developer/User] Lists plans with flexible filters.
// @Tags Plans
// @Produce json
// @Param code query string false "Plan code"
// @Param billing_cycle query string false "Billing cycle"
// @Param is_active query bool false "Active state"
// @Param search query string false "Search text"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} models.PlansAPIResponse "Plans fetched successfully"
// @Router /v1/plans [get]
func (h *PlanHandler) ListPlans(c *gin.Context) {
	filter, ok := planFilterFromQuery(c)
	if !ok {
		return
	}
	result, err := h.planSvc.ListPlans(c.Request.Context(), filter)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Plans fetched successfully", gin.H{"plans": result})
}

func planFilterFromQuery(c *gin.Context) (*models.PlanFilter, bool) {
	filter := &models.PlanFilter{Code: optionalString(c.Query("code")), Name: optionalString(c.Query("name")), BillingCycle: optionalString(c.Query("billing_cycle")), Search: optionalString(c.Query("search"))}
	if !queryInt64Ptr(c, "id", &filter.ID) {
		return nil, false
	}
	if raw := strings.TrimSpace(c.Query("is_active")); raw != "" {
		value := raw == "true" || raw == "1"
		filter.IsActive = &value
	}
	if raw := strings.TrimSpace(c.Query("min_price_paise")); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			utils.BadRequestResponse(c, "min_price_paise must be an integer")
			return nil, false
		}
		filter.MinPricePaise = &value
	}
	if raw := strings.TrimSpace(c.Query("max_price_paise")); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			utils.BadRequestResponse(c, "max_price_paise must be an integer")
			return nil, false
		}
		filter.MaxPricePaise = &value
	}
	limit, offset, ok := paginationQuery(c)
	if !ok {
		return nil, false
	}
	filter.Limit, filter.Offset = limit, offset
	return filter, true
}
