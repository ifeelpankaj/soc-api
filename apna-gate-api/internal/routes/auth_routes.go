package routes

import (
	"go-server/internal/handlers/v1"
	"go-server/internal/middlewares/guards"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutesV1(rg *gin.RouterGroup, h *handlers.AuthHandler, g *guards.Guards) {
	auth := rg.Group("/auth")

	auth.POST("/register", h.Register)
	auth.POST("/resident/register", h.RegisterResident)
	auth.POST("/verify-otp", h.VerifyOTP)
	auth.POST("/resend-otp", h.ResendOTP)
	auth.POST("/login", h.Login)
	auth.POST("/forgot-password", h.ForgotPassword)
	auth.POST("/reset-password", h.ResetPassword)
	auth.POST("/logout", h.Logout)

	refresh := auth.Group("")
	refresh.Use(g.Refresh()...)
	{
		refresh.POST("/refresh", h.Refresh)
	}

	protected := auth.Group("")
	protected.Use(g.Authenticated()...)
	{
		protected.GET("/profile", h.GetProfile)
		protected.PATCH("/profile", h.UpdateProfile)
		protected.POST("/change-password", h.ChangePassword)
	}
}
