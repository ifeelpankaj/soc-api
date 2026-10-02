package handlers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"go-server/pkg/utils"
)

const maintenanceFlatQueryMaxLen = 50

type maintenanceFlatQueryParams struct {
	Block      string
	FlatNumber string
	Search     string
}

func maintenanceFlatQuery(c *gin.Context) (maintenanceFlatQueryParams, bool) {
	block := strings.TrimSpace(c.Query("block"))
	flatNumber := strings.TrimSpace(c.Query("flat_number"))
	search := strings.TrimSpace(c.Query("search"))
	for name, value := range map[string]string{
		"block":       block,
		"flat_number": flatNumber,
		"search":      search,
	} {
		if len(value) > maintenanceFlatQueryMaxLen {
			utils.BadRequestResponse(c, name+" must be at most 50 characters")
			return maintenanceFlatQueryParams{}, false
		}
	}
	return maintenanceFlatQueryParams{
		Block:      block,
		FlatNumber: flatNumber,
		Search:     search,
	}, true
}
