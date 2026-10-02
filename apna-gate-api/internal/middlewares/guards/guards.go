package guards

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	middleware "go-server/internal/middlewares"
	"go-server/internal/models"
	"go-server/internal/requestctx"
	societysvc "go-server/internal/services/societySvc"
	subscriptionsvc "go-server/internal/services/subscriptionSvc"
	"go-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

type Option func(*Guards)

func WithSessionUsers(users middleware.SessionUserReader) Option {
	return func(g *Guards) { g.users = users }
}

type Guards struct {
	users           middleware.SessionUserReader
	jwtSecret       string
	jwtIssuer       string
	societySvc      societysvc.SocietyService
	subscriptionSvc subscriptionsvc.SubscriptionGuardService
}

func WithSubscriptionService(subscriptionSvc subscriptionsvc.SubscriptionGuardService) Option {
	return func(g *Guards) {
		g.subscriptionSvc = subscriptionSvc
	}
}

func New(jwtSecret, jwtIssuer string, args ...any) *Guards {
	g := &Guards{
		jwtSecret: jwtSecret,
		jwtIssuer: jwtIssuer,
	}

	for _, arg := range args {
		switch value := arg.(type) {
		case societysvc.SocietyService:
			g.societySvc = value
		case Option:
			value(g)
		}
	}

	return g
}

func (g *Guards) Authenticated() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.AccessAuthMiddleware(g.jwtSecret, g.jwtIssuer, g.users),
	}
}

func (g *Guards) Refresh() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.RefreshAuthMiddleware(g.jwtSecret, g.jwtIssuer, g.users),
	}
}

func (g *Guards) SocietyAdmin() []gin.HandlerFunc {
	return g.SocietyRoles(string(models.SocietyMemberRoleOwner), string(models.SocietyMemberRoleAdmin))
}

func (g *Guards) SocietyOwner() []gin.HandlerFunc {
	return g.SocietyRoles(string(models.SocietyMemberRoleOwner))
}

func (g *Guards) SocietyMember() []gin.HandlerFunc {
	return g.SocietyRoles()
}

func (g *Guards) SocietyRoles(roles ...string) []gin.HandlerFunc {
	return g.SocietyRolesForParam("societyId", roles...)
}

func (g *Guards) SocietyRolesForParam(societyIDParam string, roles ...string) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.AccessAuthMiddleware(g.jwtSecret, g.jwtIssuer, g.users),
		g.requireSocietyMembership(societyIDParam, roles...),
	}
}

func (g *Guards) SocietyMembershipRolesForParam(societyIDParam string, roles ...string) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.AccessAuthMiddleware(g.jwtSecret, g.jwtIssuer, g.users),
		g.requireSocietyMembershipRole(societyIDParam, roles...),
	}
}

func (g *Guards) GlobalRoles(roles ...models.GlobalRole) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.AccessAuthMiddleware(g.jwtSecret, g.jwtIssuer, g.users),
		g.requireGlobalRole(roles...),
	}
}

func (g *Guards) Developer() []gin.HandlerFunc {
	return g.GlobalRoles(models.GlobalRoleDeveloper, models.GlobalRoleSuperAdmin)
}

func (g *Guards) SocietyActiveForParam(societyIDParam string) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.AccessAuthMiddleware(g.jwtSecret, g.jwtIssuer, g.users),
		g.requireActiveSociety(societyIDParam),
	}
}

func (g *Guards) SocietyOperationalForParam(societyIDParam string) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.AccessAuthMiddleware(g.jwtSecret, g.jwtIssuer, g.users),
		g.requireOperationalSociety(societyIDParam),
	}
}

func (g *Guards) OperationalSocietyForParam(societyIDParam string) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		g.requireOperationalSociety(societyIDParam),
	}
}

func (g *Guards) requireGlobalRole(roles ...models.GlobalRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := middleware.GetUserRoleFromContext(c)
		if !ok || strings.TrimSpace(role) == "" {
			utils.UnauthorizedResponse(c, "Authentication required")
			c.Abort()
			return
		}

		for _, allowed := range roles {
			if role == string(allowed) {
				c.Set("global_guard_roles", roles)
				c.Next()
				return
			}
		}

		utils.ErrorResponse(c, http.StatusForbidden, models.ErrCodeForbidden, "Insufficient global role", nil)
		c.Abort()
	}
}

func (g *Guards) requireActiveSociety(societyIDParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if g.societySvc == nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, models.ErrCodeInternalServer, "Society guard is not configured", nil)
			c.Abort()
			return
		}

		societyID, ok := societyIDFromParam(c, societyIDParam)
		if !ok {
			return
		}
		if err := g.societySvc.EnsureActiveSociety(c.Request.Context(), societyID); err != nil {
			writeGuardError(c, err)
			c.Abort()
			return
		}

		c.Set("society_id", societyID)
		c.Set("societyId", societyID)
		c.Next()
	}
}

func (g *Guards) requireOperationalSociety(societyIDParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if g.subscriptionSvc == nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, models.ErrCodeInternalServer, "Subscription guard is not configured", nil)
			c.Abort()
			return
		}

		societyID, ok := societyIDFromParam(c, societyIDParam)
		if !ok {
			return
		}
		if applyDeveloperBypass(c, societyID) {
			c.Next()
			return
		}
		if err := g.subscriptionSvc.EnsureSocietyOperational(c.Request.Context(), societyID); err != nil {
			writeGuardError(c, err)
			c.Abort()
			return
		}

		c.Set("society_id", societyID)
		c.Set("societyId", societyID)
		c.Next()
	}
}

func (g *Guards) requireSocietyMembership(societyIDParam string, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if g.societySvc == nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, models.ErrCodeInternalServer, "Society guard is not configured", nil)
			c.Abort()
			return
		}

		societyID, ok := societyIDFromParam(c, societyIDParam)
		if !ok {
			return
		}

		if applyDeveloperBypass(c, societyID) {
			c.Next()
			return
		}

		userID, ok := middleware.GetUserIDFromContext(c)
		if !ok {
			utils.UnauthorizedResponse(c, "Authentication required")
			c.Abort()
			return
		}

		var err error
		if len(roles) == 0 {
			_, err = g.societySvc.EnsureActiveMember(c.Request.Context(), societyID, userID)
		} else {
			err = g.societySvc.EnsureRole(c.Request.Context(), societyID, userID, roles...)
		}
		if err != nil {
			writeGuardError(c, err)
			c.Abort()
			return
		}

		c.Set("society_id", societyID)
		c.Set("societyId", societyID)
		c.Set("society_guard_roles", roles)
		c.Next()
	}
}

func (g *Guards) requireSocietyMembershipRole(societyIDParam string, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if g.societySvc == nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, models.ErrCodeInternalServer, "Society guard is not configured", nil)
			c.Abort()
			return
		}

		societyID, ok := societyIDFromParam(c, societyIDParam)
		if !ok {
			return
		}
		if applyDeveloperBypass(c, societyID) {
			c.Next()
			return
		}

		userID, ok := middleware.GetUserIDFromContext(c)
		if !ok {
			utils.UnauthorizedResponse(c, "Authentication required")
			c.Abort()
			return
		}

		if err := g.societySvc.EnsureMembershipRole(c.Request.Context(), societyID, userID, roles...); err != nil {
			writeGuardError(c, err)
			c.Abort()
			return
		}

		c.Set("society_id", societyID)
		c.Set("societyId", societyID)
		c.Set("society_guard_roles", roles)
		c.Next()
	}
}

func isDeveloperRole(c *gin.Context) bool {
	role, ok := middleware.GetUserRoleFromContext(c)
	if !ok {
		return false
	}

	return role == string(models.GlobalRoleDeveloper) || role == string(models.GlobalRoleSuperAdmin)
}

func applyDeveloperBypass(c *gin.Context, societyID int64) bool {
	if !isDeveloperRole(c) {
		return false
	}

	c.Request = c.Request.WithContext(requestctx.WithDeveloperGuardBypass(c.Request.Context()))
	c.Set("society_id", societyID)
	c.Set("societyId", societyID)
	c.Set("global_guard_roles", []models.GlobalRole{models.GlobalRoleDeveloper, models.GlobalRoleSuperAdmin})
	return true
}

func societyIDFromParam(c *gin.Context, name string) (int64, bool) {
	raw := strings.TrimSpace(c.Param(name))
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || value <= 0 {
		utils.BadRequestResponse(c, name+" must be a positive integer")
		c.Abort()
		return 0, false
	}
	return value, true
}

func writeGuardError(c *gin.Context, err error) {
	var appErr *models.AppError
	if errors.As(err, &appErr) {
		utils.ErrorResponse(c, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Internal)
		return
	}
	utils.InternalServerErrorResponse(c, err)
}
