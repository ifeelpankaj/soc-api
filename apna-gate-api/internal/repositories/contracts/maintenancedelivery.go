package contracts

import (
	"github.com/google/uuid"
	"go-server/internal/models"
)

type MaintenanceDelivery struct {
	ID, UserID int64
	Token      uuid.UUID
	Bill       *models.MaintenanceBill
	EventType  string
	EventKey   string
	EventData  map[string]string
}
