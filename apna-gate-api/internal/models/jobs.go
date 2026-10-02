package models

import "time"

type JobTriggerResponse struct {
	Job       string `json:"job"`
	TriggerID string `json:"trigger_id"`
	Status    string `json:"status"`
}

type JobExpiryResult struct {
	WaitingVisitorEntries  int64
	ApprovedVisitorEntries int64
	VisitorInvites         int64
	FlatMemberInvites      int64
	Subscriptions          int64
}

type JobCleanupResult struct {
	Notifications     int64
	Verifications     int64
	VisitorEntries    int64
	VisitorInvites    int64
	FlatMemberInvites int64
}

type VisitorReportSociety struct {
	ID    int64
	Name  string
	Email *string
}

type VisitorReportDelivery struct {
	SocietyID       int64
	SocietyName     string
	SocietyEmail    *string
	ReportMonth     time.Time
	AttemptCount    int32
	ProcessingUntil *time.Time
	SentAt          *time.Time
}

type VisitorReportCleanupWindow struct {
	SocietyID   int64
	ReportMonth time.Time
}

type MonthlyVisitorReportRow struct {
	ID                 int64
	VisitorName        string
	VisitorPhone       *string
	VisitorEmail       *string
	FlatNumber         *string
	Block              *string
	Floor              *string
	Source             string
	Purpose            string
	Status             string
	ExpectedAt         *time.Time
	ExpectedCheckoutAt *time.Time
	CheckedInAt        *time.Time
	CheckedOutAt       *time.Time
	AutoClosedAt       *time.Time
	VehicleNumber      *string
	VehicleType        *string
	CompanionsCount    int32
	CompanionDetails   string
	ApproverName       *string
	GuardName          *string
	RejectionReason    *string
	Notes              *string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
