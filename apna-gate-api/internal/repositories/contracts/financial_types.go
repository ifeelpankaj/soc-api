package contracts

import (
	"github.com/google/uuid"
	"time"
)

type CloseUPIRequestsInput struct {
	SocietyID int64
	BillID    int64
}

type CompleteNotificationOutboxInput struct {
	Delivered  bool
	ID         int64
	LeaseToken uuid.UUID
}

type EnqueueMaintenanceReminderInput struct {
	EventKey  string
	EventData []byte
	BillID    int64
	SocietyID int64
}

type EnqueueUPIEventInput struct {
	BillID    int64
	EventType string
	EventKey  string
	Audience  string
	EventData []byte
	SocietyID int64
}

type GetActiveUPIPaymentInput struct {
	SocietyID int64
	BillID    int64
}

type GetActiveUPIRequestInput struct {
	SocietyID int64
	BillID    int64
}

type GetMaintenanceReminderDeliveryInput struct {
	UserID    int64
	EventKey  string
	SocietyID int64
}

type GetMaintenanceReminderDeliveryRecord struct {
	BillNumber             string
	BillingMonth           time.Time
	DueDate                time.Time
	OutstandingAmountPaise int64
}

type GetMaintenanceReviewInput struct {
	Token        uuid.UUID
	SocietyID    int64
	ActorID      int64
	BillingMonth time.Time
}

type GetMaintenanceReviewRecord struct {
	SnapshotHash string
	ExpiresAt    time.Time
}

type GetPaymentBillInput struct {
	SocietyID int64
	ID        int64
}

type GetPendingUPIClaimInput struct {
	SocietyID int64
	BillID    int64
}

type GetUPIClaimInput struct {
	SocietyID int64
	ID        int64
}

type GetUPIIdempotencyInput struct {
	SocietyID int64
	ActorID   int64
	Operation string
	Key       string
}

type GetUPIPaymentInput struct {
	SocietyID int64
	ID        int64
}

type GetUPIReportByFingerprintInput struct {
	SocietyID   int64
	ReportedBy  int64
	Fingerprint string
}

type GetUPIReportInput struct {
	SocietyID int64
	ID        int64
}

type GetUPIReportRecord struct {
	MaintenancePaymentReport MaintenancePaymentReport
	BillNumber               string
}

type GetUPIRequestInput struct {
	SocietyID int64
	ID        uuid.UUID
}

type GetUPISettingsVersionInput struct {
	SocietyID int64
	Version   int64
}

type HasReversedUPIReferenceInput struct {
	SocietyID int64
	Reference string
}

type InsertUPIAuditInput struct {
	SocietyID int64
	BillID    *int64
	ActorID   int64
	Action    string
	EntityID  string
	Details   []byte
}

type InsertUPIClaimInput struct {
	SocietyID   int64
	BillID      int64
	RequestID   uuid.UUID
	SubmittedBy int64
	Reference   string
	PaymentDate time.Time
}

type InsertUPIIdempotencyInput struct {
	SocietyID   int64
	ActorID     int64
	Operation   string
	Key         string
	RequestHash string
	Response    []byte
}

type InsertUPILedgerInput struct {
	SocietyID   int64
	BillID      int64
	PaymentID   int64
	Kind        string
	AmountPaise int64
	ActorID     int64
}

type InsertUPIPaymentInput struct {
	SocietyID         int64
	BillID            int64
	ClaimID           *int64
	PayerID           *int64
	SettingsVersion   int64
	AmountPaise       int64
	Reference         string
	CreditDate        time.Time
	ReceiptNumber     uuid.UUID
	VerifiedBy        int64
	EvidenceReference string
}

type InsertUPIReportInput struct {
	SocietyID   int64
	BillID      int64
	RequestID   uuid.UUID
	ReportedBy  int64
	Reference   string
	AmountPaise int64
	PaymentDate time.Time
	Explanation string
	Fingerprint string
}

type InsertUPIRequestInput struct {
	ID              uuid.UUID
	SocietyID       int64
	BillID          int64
	SettingsVersion int64
	AmountPaise     int64
	Reference       string
	CreatedBy       int64
}

type InsertUPISettingsVersionInput struct {
	SocietyID int64
	Version   int64
	Enabled   bool
	UpiID     string
	PayeeName string
	CreatedBy int64
}

type ListMaintenanceReminderCandidatesInput struct {
	SocietyID int64
	AfterID   int64
	AsOf      time.Time
}

type ListMaintenanceReminderCandidatesRecord struct {
	ID       int64
	DueDate  time.Time
	Timezone string
}

type ListUPIAuditInput struct {
	SocietyID int64
	BillID    int64
	BeforeID  int64
	Limit     int32
}

type ListUPIClaimsInput struct {
	SocietyID    int64
	UserID       int64
	Status       string
	Reference    string
	BeforeID     int64
	BillID       int64
	FlatID       int64
	BillingMonth time.Time
	Block        *string
	FlatNumber   *string
	Search       string
	Limit        int32
}

type ListUPIPaymentsInput struct {
	SocietyID    int64
	Status       string
	Reference    string
	BeforeID     int64
	BillID       int64
	FlatID       int64
	BillingMonth time.Time
	Block        *string
	FlatNumber   *string
	Search       string
	UserID       int64
	Limit        int32
}

type ListUPIReportsInput struct {
	SocietyID int64
	UserID    int64
	Status    string
	BeforeID  int64
	Limit     int32
}

type ListUPIReportsRecord struct {
	MaintenancePaymentReport MaintenancePaymentReport
	BillNumber               string
}

type ListUPISettingsVersionsInput struct {
	SocietyID     int64
	BeforeVersion int64
	Limit         int32
}

type LockMaintenancePaymentBillInput struct {
	SocietyID int64
	ID        int64
}

type MaintenanceAdminInput struct {
	SocietyID int64
	UserID    int64
}

type MaintenanceBill struct {
	ID             int64
	RunID          int64
	SocietyID      int64
	FlatID         int64
	BillingMonth   time.Time
	BillNumber     string
	DueDate        time.Time
	Timezone       string
	TotalPaise     int64
	Snapshot       []byte
	BilledParty    []byte
	CreatedAt      time.Time
	IssuerSnapshot []byte
}

type MaintenancePayment struct {
	ID                int64
	SocietyID         int64
	BillID            int64
	ClaimID           *int64
	PayerID           *int64
	SettingsVersion   int64
	AmountPaise       int64
	Reference         string
	CreditDate        time.Time
	ReceiptNumber     uuid.UUID
	VerifiedBy        int64
	VerifiedAt        time.Time
	Status            string
	ReversedBy        *int64
	ReversedAt        time.Time
	ReversalReason    *string
	EvidenceReference string
}

type MaintenancePaymentAuditEvent struct {
	ID        int64
	SocietyID int64
	BillID    *int64
	ActorID   int64
	Action    string
	EntityID  string
	Details   []byte
	CreatedAt time.Time
}

type MaintenancePaymentClaim struct {
	ID          int64
	SocietyID   int64
	BillID      int64
	RequestID   uuid.UUID
	SubmittedBy int64
	Reference   string
	PaymentDate time.Time
	Status      string
	ReviewedBy  *int64
	ReviewedAt  time.Time
	Reason      *string
	CreatedAt   time.Time
}

type MaintenancePaymentIdempotency struct {
	SocietyID   int64
	ActorID     int64
	Operation   string
	Key         string
	RequestHash string
	Response    []byte
	CreatedAt   time.Time
}

type MaintenancePaymentReport struct {
	ID             int64
	SocietyID      int64
	BillID         int64
	RequestID      uuid.UUID
	ReportedBy     int64
	Reference      string
	AmountPaise    int64
	PaymentDate    time.Time
	Explanation    string
	Fingerprint    string
	Status         string
	ResolutionNote *string
	UpdatedBy      *int64
	UpdatedAt      time.Time
	CreatedAt      time.Time
}

type MaintenancePaymentRequest struct {
	ID              uuid.UUID
	SocietyID       int64
	BillID          int64
	SettingsVersion int64
	AmountPaise     int64
	Reference       string
	State           string
	CreatedBy       int64
	CreatedAt       time.Time
}

type MaintenancePaymentSettingsVersion struct {
	SocietyID int64
	Version   int64
	Enabled   bool
	UpiID     string
	PayeeName string
	CreatedBy int64
	CreatedAt time.Time
}

type MaintenanceResidentAccessInput struct {
	SocietyID int64
	ID        int64
	UserID    int64
}

type MarkOutboxInboxInput struct {
	ID         int64
	LeaseToken uuid.UUID
}

type NotificationBacklogRecord struct {
	Pending       int64
	OldestSeconds float64
}

type NotificationOutbox struct {
	ID               int64
	UserID           int64
	SocietyID        int64
	FlatID           *int64
	Audience         string
	EventKey         string
	Payload          []byte
	Attempts         int32
	AvailableAt      time.Time
	LeaseToken       uuid.UUID
	InboxCompletedAt time.Time
	PushCompletedAt  time.Time
	CompletedAt      time.Time
	LastError        *string
	PushEnabled      bool
}

type NotificationOutboxAccessInput struct {
	SocietyID int64
	UserID    int64
	Audience  string
	FlatID    *int64
}

type OutboxInboxIDInput struct {
	UserID   int64
	EventKey *string
}

type PointUPISettingsInput struct {
	SocietyID int64
	Version   int64
}

type ReleaseUPIClaimReferenceInput struct {
	SocietyID int64
	ClaimID   *int64
}

type ReleaseUPIPaymentReferenceInput struct {
	SocietyID int64
	PaymentID *int64
}

type ReserveUPIReferenceInput struct {
	SocietyID int64
	Reference string
	BillID    int64
	ClaimID   *int64
	PaymentID *int64
}

type RetryNotificationOutboxInput struct {
	ID         int64
	LeaseToken uuid.UUID
}

type ReverseUPIPaymentInput struct {
	SocietyID      int64
	ID             int64
	ReversedBy     *int64
	ReversalReason *string
}

type ReviewUPIClaimInput struct {
	SocietyID  int64
	ID         int64
	Status     string
	ReviewedBy *int64
	Reason     *string
}

type SaveMaintenanceReviewInput struct {
	Token        uuid.UUID
	SocietyID    int64
	ActorID      int64
	BillingMonth time.Time
	SnapshotHash string
	ExpiresAt    time.Time
}

type UPICollectionSummaryInput struct {
	SocietyID int64
	FlatID    int64
	Month     time.Time
}

type UPICollectionSummaryRecord struct {
	BilledCount      int64
	BilledPaise      int64
	CollectedPaise   int64
	OutstandingPaise int64
	OverduePaise     int64
	PendingClaims    int64
}

type UpdateUPIReportInput struct {
	SocietyID      int64
	ID             int64
	Status         string
	ResolutionNote *string
	UpdatedBy      *int64
}
