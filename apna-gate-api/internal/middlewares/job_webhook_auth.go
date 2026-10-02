package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
	"strings"

	"go-server/internal/models"
	"go-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

// JobWebhookAuth authenticates internal cron requests with a shared bearer secret.
func JobWebhookAuth(secret string) gin.HandlerFunc {
	expectedHash := sha256.Sum256([]byte(strings.TrimSpace(secret)))
	return func(c *gin.Context) {
		parts := strings.Fields(c.GetHeader("Authorization"))
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			jobWebhookUnauthorized(c)
			return
		}
		actualHash := sha256.Sum256([]byte(parts[1]))
		if subtle.ConstantTimeCompare(expectedHash[:], actualHash[:]) != 1 {
			jobWebhookUnauthorized(c)
			return
		}
		c.Next()
	}
}

func jobWebhookUnauthorized(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusUnauthorized, models.ErrCodeUnauthorized, "Invalid job webhook credentials", nil)
	c.Abort()
}
