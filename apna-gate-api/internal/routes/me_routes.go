package routes

import (
	"go-server/internal/handlers/v1"
	"go-server/internal/middlewares/guards"

	"github.com/gin-gonic/gin"
)

func SetupMeRoutesV1(rg *gin.RouterGroup, h *handlers.NotificationHandler, g *guards.Guards) {
	me := rg.Group("/me")
	me.Use(g.Authenticated()...)
	{
		me.POST("/device-tokens", h.RegisterDeviceToken)
		me.DELETE("/device-tokens", h.UnregisterDeviceToken)
		me.GET("/notifications", h.ListNotifications)
		me.GET("/notifications/unread-count", h.GetUnreadCount)
		me.GET("/notifications/preferences", h.GetPreferences)
		me.PUT("/notifications/preferences", h.SetPreferences)
		me.GET("/notifications/:notificationId", h.GetNotification)
		me.PATCH("/notifications/:notificationId/read", h.MarkNotificationRead)
		me.POST("/notifications/read-all", h.MarkAllNotificationsRead)
	}
}
