package routes

import (
	"go-server/internal/handlers/v1"
	"go-server/internal/middlewares/guards"

	"github.com/gin-gonic/gin"
)

func SetupBootstrapRoutesV1(rg *gin.RouterGroup, h *handlers.BootstrapHandler, g *guards.Guards) {
	bootstrap := rg.Group("")
	bootstrap.Use(g.Authenticated()...)
	{
		bootstrap.GET("/bootstrap", h.GetBootstrap)
	}
}
