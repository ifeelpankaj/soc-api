package routes

import (
	"github.com/gin-gonic/gin"
	handlers "go-server/internal/handlers/v1"
	"go-server/internal/middlewares/guards"
)

func SetupHubRoutesV1(rg *gin.RouterGroup, h *handlers.HubHandler, g *guards.Guards) {
	if h == nil {
		return
	}
	hub := rg.Group("/societies/:societyId")
	// Do not use SocietyMember/Operational middleware here: they allow platform
	// bypass. The service enforces membership and subscription without that bypass.
	hub.Use(g.Authenticated()...)
	hub.Use(h.Timeout)
	hub.GET("/channels", h.Channels)
	hub.GET("/channels/categories", h.Categories)
	hub.GET("/channels/:channelId/posts", h.Posts)
	hub.POST("/channels/:channelId/posts", h.CreatePost)
	hub.POST("/channels/:channelId/read", h.Read)
	hub.GET("/posts/:postId", h.Post)
	hub.PATCH("/posts/:postId", h.UpdatePost)
	hub.DELETE("/posts/:postId", h.DeletePost)
	hub.GET("/posts/:postId/comments", h.Comments)
	hub.POST("/posts/:postId/comments", h.CreateComment)
	hub.PATCH("/comments/:commentId", h.UpdateComment)
	hub.DELETE("/comments/:commentId", h.DeleteComment)
	hub.POST("/posts/:postId/pin", h.Pin)
	hub.POST("/posts/:postId/lock", h.Lock)
	hub.POST("/posts/:postId/reactions", h.PostReaction)
	hub.DELETE("/posts/:postId/reactions/:reactionType", h.RemovePostReaction)
	hub.POST("/comments/:commentId/reactions", h.CommentReaction)
	hub.DELETE("/comments/:commentId/reactions/:reactionType", h.RemoveCommentReaction)
	hub.POST("/posts/:postId/reports", h.Report)
	hub.GET("/channel-reports", h.Reports)
	hub.PATCH("/channel-reports/:reportId", h.Resolve)
	hub.POST("/channel-uploads", h.Upload)
	hub.DELETE("/channel-uploads/:uploadId", h.DeleteUpload)
	hub.GET("/channel-attachments/:uploadId", h.Attachment)
}
