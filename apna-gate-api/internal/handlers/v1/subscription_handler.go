package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-server/internal/models"
	subscriptionsvc "go-server/internal/services/subscriptionSvc"
	"go-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	subSvc subscriptionsvc.SubscriptionService
}

func NewSubscriptionHandler(subSvc subscriptionsvc.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{subSvc: subSvc}
}

// CreatePendingSubscription godoc
// @Summary Create pending subscription
// @Tags Subscriptions
// @Produce json
// @Param societyId path int true "Society ID"
// @Param planId path int true "Plan ID"
// @Success 201 {object} models.SubscriptionAPIResponse "Subscription created successfully"
// @Router /v1/societies/{societyId}/subscriptions/plans/{planId}/pending [post]
func (h *SubscriptionHandler) CreatePendingSubscription(c *gin.Context) {
	societyID, planID, ok := parseSocietyAndPlan(c)
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.subSvc.CreatePendingSubscription(c.Request.Context(), societyID, planID, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Subscription created successfully", gin.H{"subscription": result})
}

// CreateTrialSubscription godoc
// @Summary Create trial subscription
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param planId path int true "Plan ID"
// @Param request body models.CreateTrialSubscriptionRequest true "Trial subscription payload"
// @Success 201 {object} models.SubscriptionAPIResponse "Trial subscription created successfully"
// @Router /v1/societies/{societyId}/subscriptions/plans/{planId}/trial [post]
func (h *SubscriptionHandler) CreateTrialSubscription(c *gin.Context) {
	societyID, planID, ok := parseSocietyAndPlan(c)
	if !ok {
		return
	}
	var req models.CreateTrialSubscriptionRequest
	if !bindJSON(c, &req) {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.subSvc.CreateTrialSubscription(c.Request.Context(), societyID, planID, userID, &req)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Trial subscription created successfully", gin.H{"subscription": result})
}

// ActivateSubscription godoc
// @Summary Activate subscription
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param subscriptionId path int true "Subscription ID"
// @Param request body models.ActivateSubscriptionRequest true "Activate subscription payload"
// @Success 200 {object} models.SubscriptionAPIResponse "Subscription activated successfully"
// @Router /v1/subscriptions/{subscriptionId}/activate [post]
func (h *SubscriptionHandler) ActivateSubscription(c *gin.Context) {
	id, ok := parsePathInt64(c, "subscriptionId")
	if !ok {
		return
	}
	var req models.ActivateSubscriptionRequest
	if !bindJSON(c, &req) {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.subSvc.ActivateSubscription(c.Request.Context(), id, userID, &req)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Subscription activated successfully", gin.H{"subscription": result})
}

// RenewSubscription godoc
// @Summary Renew subscription
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param subscriptionId path int true "Subscription ID"
// @Param request body models.RenewSubscriptionRequest true "Renew subscription payload"
// @Success 200 {object} models.SubscriptionAPIResponse "Subscription renewed successfully"
// @Router /v1/subscriptions/{subscriptionId}/renew [post]
func (h *SubscriptionHandler) RenewSubscription(c *gin.Context) {
	id, ok := parsePathInt64(c, "subscriptionId")
	if !ok {
		return
	}
	var req models.RenewSubscriptionRequest
	if !bindJSON(c, &req) {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.subSvc.RenewSubscription(c.Request.Context(), id, userID, &req)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Subscription renewed successfully", gin.H{"subscription": result})
}

// CancelSubscription godoc
// @Summary Cancel subscription
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param subscriptionId path int true "Subscription ID"
// @Param request body models.CancelSubscriptionRequest true "Cancel subscription payload"
// @Success 200 {object} models.SubscriptionAPIResponse "Subscription cancelled successfully"
// @Router /v1/subscriptions/{subscriptionId}/cancel [post]
func (h *SubscriptionHandler) CancelSubscription(c *gin.Context) {
	id, ok := parsePathInt64(c, "subscriptionId")
	if !ok {
		return
	}
	var req models.CancelSubscriptionRequest
	if !bindJSON(c, &req) {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.subSvc.CancelSubscription(c.Request.Context(), id, userID, &req)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Subscription cancelled successfully", gin.H{"subscription": result})
}

// ExpireSubscription godoc
// @Summary Expire subscription
// @Tags Subscriptions
// @Produce json
// @Param subscriptionId path int true "Subscription ID"
// @Success 200 {object} models.SubscriptionAPIResponse "Subscription expired successfully"
// @Router /v1/subscriptions/{subscriptionId}/expire [post]
func (h *SubscriptionHandler) ExpireSubscription(c *gin.Context) {
	id, ok := parsePathInt64(c, "subscriptionId")
	if !ok {
		return
	}
	result, err := h.subSvc.ExpireSubscription(c.Request.Context(), id)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Subscription expired successfully", gin.H{"subscription": result})
}

// ChangeSubscriptionPlan godoc
// @Summary Change subscription plan
// @Tags Subscriptions
// @Produce json
// @Param subscriptionId path int true "Subscription ID"
// @Param planId path int true "New plan ID"
// @Success 200 {object} models.SubscriptionAPIResponse "Subscription plan changed successfully"
// @Router /v1/subscriptions/{subscriptionId}/plans/{planId} [post]
func (h *SubscriptionHandler) ChangeSubscriptionPlan(c *gin.Context) {
	id, ok := parsePathInt64(c, "subscriptionId")
	if !ok {
		return
	}
	planID, ok := parsePathInt64(c, "planId")
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.subSvc.ChangeSubscriptionPlan(c.Request.Context(), id, planID, userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Subscription plan changed successfully", gin.H{"subscription": result})
}

// GetSubscription godoc
// @Summary Get subscription
// @Tags Subscriptions
// @Produce json
// @Param id query int false "Subscription ID"
// @Param society_id query int false "Society ID"
// @Success 200 {object} models.SubscriptionAPIResponse "Subscription fetched successfully"
// @Router /v1/subscriptions/lookup [get]
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	filter, ok := subscriptionFilterFromQuery(c)
	if !ok {
		return
	}
	result, err := h.subSvc.GetSubscription(c.Request.Context(), filter)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Subscription fetched successfully", gin.H{"subscription": result})
}

// ListSubscriptions godoc
// @Summary List subscriptions
// @Tags Subscriptions
// @Produce json
// @Param society_id query int false "Society ID"
// @Param plan_id query int false "Plan ID"
// @Param status query string false "Subscription status"
// @Param search query string false "Search text"
// @Param search_mode query string false "Search mode"
// @Param is_active_only query bool false "Only active subscriptions"
// @Param expired_only query bool false "Only expired subscriptions"
// @Param expiring_before query string false "Subscriptions ending before RFC3339 timestamp"
// @Param expiring_days query int false "Subscriptions ending within N days"
// @Param limit query int false "Page size"
// @Param offset query int false "Page offset"
// @Success 200 {object} models.SubscriptionsAPIResponse "Subscriptions fetched successfully"
// @Router /v1/subscriptions [get]
func (h *SubscriptionHandler) ListSubscriptions(c *gin.Context) {
	filter, ok := subscriptionFilterFromQuery(c)
	if !ok {
		return
	}
	result, err := h.subSvc.ListSubscriptions(c.Request.Context(), filter)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Subscriptions fetched successfully", gin.H{"subscriptions": result})
}

// GetSubscriptionStats godoc
// @Summary Get subscription stats
// @Tags Subscriptions
// @Produce json
// @Success 200 {object} models.SubscriptionStatsAPIResponse "Subscription stats fetched successfully"
// @Router /v1/subscriptions/stats [get]
func (h *SubscriptionHandler) GetSubscriptionStats(c *gin.Context) {
	filter, ok := subscriptionFilterFromQuery(c)
	if !ok {
		return
	}
	result, err := h.subSvc.GetSubscriptionStats(c.Request.Context(), filter)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Subscription stats fetched successfully", gin.H{"stats": result})
}

func parseSocietyAndPlan(c *gin.Context) (int64, int64, bool) {
	societyID, ok := parsePathInt64(c, "societyId")
	if !ok {
		return 0, 0, false
	}
	planID, ok := parsePathInt64(c, "planId")
	return societyID, planID, ok
}

func subscriptionFilterFromQuery(c *gin.Context) (*models.SubscriptionFilter, bool) {
	filter := &models.SubscriptionFilter{Status: optionalString(c.Query("status")), PlanCode: optionalString(c.Query("plan_code")), BillingCycle: optionalString(c.Query("billing_cycle")), Search: optionalString(c.Query("search")), SearchMode: optionalString(c.Query("search_mode"))}
	if !queryInt64Ptr(c, "id", &filter.ID) || !queryInt64Ptr(c, "society_id", &filter.SocietyID) || !queryInt64Ptr(c, "plan_id", &filter.PlanID) {
		return nil, false
	}
	if raw := strings.TrimSpace(c.Query("is_active_only")); raw != "" {
		value := raw == "true" || raw == "1"
		filter.IsActiveOnly = &value
	}
	if raw := strings.TrimSpace(c.Query("expired_only")); raw != "" {
		value := raw == "true" || raw == "1"
		filter.ExpiredOnly = &value
	}
	if !queryTimePtr(c, "expiring_before", &filter.ExpiringBefore) {
		return nil, false
	}
	if raw := strings.TrimSpace(c.Query("expiring_days")); raw != "" {
		days, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || days <= 0 {
			utils.BadRequestResponse(c, "expiring_days must be a positive integer")
			return nil, false
		}
		expiringBefore := time.Now().UTC().Add(time.Duration(days) * 24 * time.Hour)
		filter.ExpiringBefore = &expiringBefore
	}
	if !queryTimePtr(c, "starts_after", &filter.StartsAfter) ||
		!queryTimePtr(c, "starts_before", &filter.StartsBefore) ||
		!queryTimePtr(c, "ends_after", &filter.EndsAfter) ||
		!queryTimePtr(c, "ends_before", &filter.EndsBefore) {
		return nil, false
	}
	limit, offset, ok := paginationQuery(c)
	if !ok {
		return nil, false
	}
	filter.Limit, filter.Offset = limit, offset
	return filter, true
}
