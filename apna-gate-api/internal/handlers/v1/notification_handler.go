package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"go-server/internal/models"
	"go-server/internal/repositories/contracts"
	notificationsvc "go-server/internal/services/notificationSvc"
	"go-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationSvc notificationsvc.NotificationService
}

// ListNotifications godoc
// @Summary List my notifications
// @Description Lists durable notifications for the authenticated user.
// @Tags Notifications
// @Produce json
// @Param limit query int false "Page size" default(20)
// @Param cursor query string false "Opaque pagination cursor"
// @Success 200 {object} models.NotificationsAPIResponse "Notifications fetched successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid notification request"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/me/notifications [get]
func (h *NotificationHandler) ListNotifications(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	limit := int32(20)
	if raw := strings.TrimSpace(c.Query("limit")); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 32)
		if err != nil || value <= 0 {
			utils.BadRequestResponse(c, "limit must be a positive integer")
			return
		}
		limit = int32(value)
	}
	result, err := h.notificationSvc.ListNotifications(c.Request.Context(), userID, limit, c.Query("cursor"))
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Notifications fetched successfully", result)
}

// GetUnreadCount godoc
// @Summary Get notification unread count
// @Description Returns the authenticated user's unread notification count.
// @Tags Notifications
// @Produce json
// @Success 200 {object} models.NotificationUnreadCountAPIResponse "Unread count fetched successfully"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/me/notifications/unread-count [get]
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.notificationSvc.GetUnreadCount(c.Request.Context(), userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Unread count fetched successfully", result)
}

// MarkNotificationRead godoc
// @Summary Mark notification read
// @Description Marks one owned notification as read.
// @Tags Notifications
// @Produce json
// @Param notificationId path string true "Notification ID"
// @Success 200 {object} models.NotificationReadAPIResponse "Notification marked read"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid notification request"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 404 {object} models.ErrorResponseDoc "Notification not found"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/me/notifications/{notificationId}/read [patch]
func (h *NotificationHandler) MarkNotificationRead(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.notificationSvc.MarkNotificationRead(c.Request.Context(), userID, c.Param("notificationId"))
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Notification marked read", result)
}

// MarkAllNotificationsRead godoc
// @Summary Mark all notifications read
// @Description Marks all notifications for the authenticated user as read.
// @Tags Notifications
// @Produce json
// @Success 200 {object} models.NotificationUnreadCountAPIResponse "Notifications marked read"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/me/notifications/read-all [post]
func (h *NotificationHandler) MarkAllNotificationsRead(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.notificationSvc.MarkAllNotificationsRead(c.Request.Context(), userID)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Notifications marked read", result)
}

func NewNotificationHandler(notificationSvc notificationsvc.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationSvc: notificationSvc}
}

// RegisterDeviceToken godoc
// @Summary Register or refresh a device push token
// @Description Stores an FCM token. Web requires an installation device_id and active resident access; the server binds it to the current authenticated session version. Unbound/revoked web tokens never receive delivery.
// @Tags Notifications
// @Accept json
// @Produce json
// @Param request body models.RegisterDeviceTokenRequest true "Device token payload"
// @Success 200 {object} models.RegisterDeviceTokenAPIResponse "Device token registered successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid device token request"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/me/device-tokens [post]
func (h *NotificationHandler) RegisterDeviceToken(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var req models.RegisterDeviceTokenRequest
	if !bindJSON(c, &req) {
		return
	}

	version := c.GetInt64("session_version")
	req.SessionVersion = &version
	result, err := h.notificationSvc.RegisterDeviceToken(c.Request.Context(), userID, req)
	if handleServiceError(c, err) {
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Device token registered successfully", gin.H{"device_token": result})
}

// UnregisterDeviceToken godoc
// @Summary Unregister a device push token
// @Description Removes the authenticated user's FCM device token, typically on logout.
// @Tags Notifications
// @Accept json
// @Produce json
// @Param request body models.UnregisterDeviceTokenRequest true "Device token payload"
// @Success 200 {object} models.MessageAPIResponse "Device token removed successfully"
// @Failure 400 {object} models.ErrorResponseDoc "Invalid device token request"
// @Failure 401 {object} models.ErrorResponseDoc "Missing, invalid, or expired access token"
// @Failure 500 {object} models.ErrorResponseDoc "Internal server error"
// @Security AccessToken
// @Router /v1/me/device-tokens [delete]
func (h *NotificationHandler) UnregisterDeviceToken(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var req models.UnregisterDeviceTokenRequest
	if !bindJSON(c, &req) {
		return
	}

	if handleServiceError(c, h.notificationSvc.UnregisterDeviceToken(c.Request.Context(), userID, req.Token)) {
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Device token removed successfully", nil)
}

// GetNotification godoc
// @Summary Get one owned notification for authenticated tap routing
// @Tags Notifications
// @Produce json
// @Param notificationId path string true "Notification ID"
// @Success 200 {object} models.NotificationDetailAPIResponse
// @Failure 400 {object} models.ErrorResponseDoc
// @Failure 401 {object} models.ErrorResponseDoc
// @Failure 404 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/me/notifications/{notificationId} [get]
func (h *NotificationHandler) GetNotification(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	svc, ok := h.notificationSvc.(interface {
		GetNotification(context.Context, int64, string) (*models.Notification, error)
	})
	if !ok {
		handleServiceError(c, notificationsvc.ErrNotificationNotFound)
		return
	}
	item, err := svc.GetNotification(c.Request.Context(), userID, c.Param("notificationId"))
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Notification fetched successfully", item)
}

func (h *NotificationHandler) GetPreferences(c *gin.Context) {
	user, ok := currentUserID(c)
	if !ok {
		return
	}
	society, err := strconv.ParseInt(c.Query("society_id"), 10, 64)
	if err != nil || society <= 0 {
		utils.BadRequestResponse(c, "society_id must be positive")
		return
	}
	svc, ok := h.notificationSvc.(interface {
		GetPreferences(context.Context, int64, int64) (contracts.NotificationPreferences, error)
	})
	if !ok {
		utils.InternalServerErrorResponse(c, errors.New("notification preferences unavailable"))
		return
	}
	p, err := svc.GetPreferences(c.Request.Context(), user, society)
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Notification preferences fetched", p)
}

func (h *NotificationHandler) SetPreferences(c *gin.Context) {
	user, ok := currentUserID(c)
	if !ok {
		return
	}
	society, err := strconv.ParseInt(c.Query("society_id"), 10, 64)
	if err != nil || society <= 0 {
		utils.BadRequestResponse(c, "society_id must be positive")
		return
	}
	var req struct {
		VisitorUpdatesPush *bool `json:"visitor_updates_push"`
		MaintenancePush    *bool `json:"maintenance_push"`
		AnnouncementsPush  *bool `json:"announcements_push"`
		HubRepliesPush     *bool `json:"hub_replies_push"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if req.VisitorUpdatesPush == nil || req.MaintenancePush == nil || req.AnnouncementsPush == nil || req.HubRepliesPush == nil {
		utils.BadRequestResponse(c, "all preference fields are required")
		return
	}
	p := contracts.NotificationPreferences{VisitorUpdatesPush: *req.VisitorUpdatesPush, MaintenancePush: *req.MaintenancePush, AnnouncementsPush: *req.AnnouncementsPush, HubRepliesPush: *req.HubRepliesPush}
	svc, ok := h.notificationSvc.(interface {
		SetPreferences(context.Context, int64, int64, contracts.NotificationPreferences) error
	})
	if !ok {
		utils.InternalServerErrorResponse(c, errors.New("notification preferences unavailable"))
		return
	}
	if handleServiceError(c, svc.SetPreferences(c.Request.Context(), user, society, p)) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Notification preferences updated", p)
}
