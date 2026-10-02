package contracts

import (
	"context"
	"go-server/internal/models"

	"github.com/google/uuid"
)

type PushDelivery struct {
	ID             int64
	NotificationID string
	UserID         int64
	SocietyID      *int64
	FlatID         *int64
	TokenID        *int64
	Token          string
	Platform       models.DevicePlatform
	Provider       string
	Audience       string
	Title          string
	Body           string
	Data           map[string]any
	AttemptCount   int
	LeaseToken     uuid.UUID
}

type HubFanoutJob struct {
	ID          int64
	SocietyID   int64
	PostID      int64
	ActorUserID int64
	Important   bool
	LeaseToken  uuid.UUID
}

type NotificationPreferences struct {
	VisitorUpdatesPush bool `json:"visitor_updates_push"`
	MaintenancePush    bool `json:"maintenance_push"`
	AnnouncementsPush  bool `json:"announcements_push"`
	HubRepliesPush     bool `json:"hub_replies_push"`
}

type PipelineRepository interface {
	MaterializeOutbox(context.Context, NotificationOutbox, models.NotificationCreate, bool, bool, bool, bool) error
	ClaimPushDelivery(context.Context, uuid.UUID) (*PushDelivery, error)
	FinishPushDelivery(context.Context, int64, uuid.UUID, string) error
	FailPushDelivery(context.Context, int64, uuid.UUID, string, string, bool) error
	ClaimHubFanout(context.Context, uuid.UUID) (*HubFanoutJob, error)
	FanoutHubBatch(context.Context, HubFanoutJob, int) (int, error)
	GetNotificationPreferences(context.Context, int64, int64) (NotificationPreferences, error)
	SetNotificationPreferences(context.Context, int64, int64, NotificationPreferences) error
	ReplayPushDelivery(context.Context, int64) error
	PipelineStats(context.Context) (map[string]float64, error)
	HubContentAvailable(context.Context, int64, int64) (bool, error)
}
