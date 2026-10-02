package routes

import (
	"go-server/internal/app"
	"go-server/internal/middlewares/guards"

	"github.com/gin-gonic/gin"
)

func SetupPublicRoutesV1(rg *gin.RouterGroup, h *app.V1Handlers, g *guards.Guards) {
	public := rg.Group("/public")
	{
		public.GET("/societies/:societyCode/claim-options", h.Society.GetPublicClaimOptions)
		public.GET("/societies/:societyCode/visitor-entry-options", h.VisitorEntry.GetEntryOptions)
		public.GET("/short-links/:shortCode", h.ShortLink.Resolve)
		public.POST("/short-links/:shortCode/visitor-submit", h.ShortLink.SubmitVisitorInvite)
		public.POST("/short-links/:shortCode/member-join", h.ShortLink.JoinMemberInvite)
		acceptShortLinkInvite := append(g.Authenticated(), h.ShortLink.AcceptMemberInvite)
		public.POST("/short-links/:shortCode/member-accept", acceptShortLinkInvite...)
		public.GET("/visitor-invites/:token", h.VisitorEntry.GetInviteByToken)
		public.POST("/societies/:societyCode/visitor-invites/:token/submit", h.VisitorEntry.SubmitInviteForm)
		public.POST("/societies/:societyCode/visitor-entries/public-qr", h.VisitorEntry.CreatePublicQREntry)
		public.POST("/societies/:societyCode/visitor-entries/quick-link", h.VisitorEntry.CreateQuickLinkEntry)
		public.POST("/visitor-entries/qr/validate", h.VisitorEntry.ValidateQR)
		public.GET("/flat-member-invites/:token", h.MemberInvite.GetPublicMemberInviteByToken)
		public.POST("/flat-member-invites/:token/join", h.MemberInvite.JoinMemberInvite)
		acceptInvite := append(g.Authenticated(), h.MemberInvite.AcceptMemberInvite)
		public.POST("/flat-member-invites/:token/accept", acceptInvite...)
	}
}
