package config

import (
	"fmt"
	"os"
	"strings"
)

func maintenanceStatusesFromEnv() []string {
	value, ok := os.LookupEnv("MAINTENANCE_ELIGIBLE_FLAT_STATUSES")
	if !ok {
		return []string{"vacant", "occupied", "blocked"}
	}
	return strings.Split(value, ",")
}

func (c *Config) MaintenanceStatuses() ([]string, error) {
	values := c.MaintenanceEligibleFlatStatuses
	if values == nil {
		values = []string{"vacant", "occupied", "blocked"}
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("MAINTENANCE_ELIGIBLE_FLAT_STATUSES cannot be empty")
	}
	result := make([]string, 0, len(values))
	seen := map[string]bool{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		switch value {
		case "vacant", "occupied", "blocked":
		default:
			return nil, fmt.Errorf("MAINTENANCE_ELIGIBLE_FLAT_STATUSES contains invalid status %q", value)
		}
		if !seen[value] {
			result = append(result, value)
			seen[value] = true
		}
	}
	return result, nil
}
