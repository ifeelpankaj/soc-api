package routes

import (
	"github.com/gin-gonic/gin"
	handlers "go-server/internal/handlers/v1"
	"go-server/internal/middlewares/guards"
)

func SetupImageRoutesV1(rg *gin.RouterGroup, h *handlers.ImageHandler, g *guards.Guards) {
	images := rg.Group("")
	images.Use(func(c *gin.Context) { c.Header("Cache-Control", "no-store"); c.Next() })
	images.Use(g.Authenticated()...)
	images.PUT("/auth/profile/avatar", h.PutAvatar)
	images.GET("/auth/profile/avatar", h.GetAvatar)
	images.DELETE("/auth/profile/avatar", h.DeleteAvatar)
	images.PUT("/societies/:societyId/visitor-entries/:entryId/photo", h.PutVisitorPhoto)
	images.GET("/societies/:societyId/visitor-entries/:entryId/photo", h.GetVisitorPhoto)
	images.DELETE("/societies/:societyId/visitor-entries/:entryId/photo", h.DeleteVisitorPhoto)
	images.GET("/societies/:societyId/flats/:flatId/visitor-entries/:entryId/photo", h.GetResidentVisitorPhoto)
}
