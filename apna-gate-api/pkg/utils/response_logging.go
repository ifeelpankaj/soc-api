package utils

import (
	"github.com/gin-gonic/gin"
	"go-server/pkg/logger"
)

func recordResponse(c *gin.Context, message, code string, err error) {
	c.Set("response_message", logger.Sanitize(message))
	c.Set("response_message_source", "response_helper")
	c.Set("error_code", logger.Sanitize(code))
	if err != nil {
		c.Set("internal_error", logger.Sanitize(err.Error()))
	}
}
