package models

import "time"

type UPIConfig struct {
	Enabled   bool   `json:"enabled"`
	UPIID     string `json:"upi_id"`
	PayeeName string `json:"payee_name"`
}
type MaintenancePaymentMethods struct {
	UPI UPIConfig `json:"upi"`
}
type MaintenancePaymentSettings struct {
	PaymentMethods MaintenancePaymentMethods `json:"payment_methods"`
	Version        int64                     `json:"version"`
	UpdatedBy      int64                     `json:"updated_by,omitempty"`
	UpdatedAt      time.Time                 `json:"updated_at,omitempty"`
}
type UPIPaymentRequest struct {
	ID              string    `json:"id"`
	BillID          int64     `json:"bill_id"`
	Reference       string    `json:"reference"`
	State           string    `json:"state"`
	AmountPaise     int64     `json:"amount_paise"`
	SettingsVersion int64     `json:"settings_version"`
	Destination     UPIConfig `json:"destination"`
	URI             string    `json:"upi_uri"`
	QRURL           string    `json:"qr_url"`
}
type UPISubmitClaim struct {
	PaymentRequestID string `json:"payment_request_id"`
	Reference        string `json:"reference"`
	PaymentDate      string `json:"payment_date"`
}
type UPIVerifyCredit struct {
	CorrectionReason            string `json:"correction_reason,omitempty"`
	EvidenceReference           string `json:"evidence_reference,omitempty"`
	BankCreditConfirmed         bool   `json:"bank_credit_confirmed"`
	AmountPaise                 int64  `json:"amount_paise"`
	CreditDate                  string `json:"credit_date"`
	Reference                   string `json:"reference"`
	SettingsVersion             int64  `json:"settings_version"`
	UPIID                       string `json:"upi_id"`
	ReversalHistoryAcknowledged bool   `json:"reversal_history_acknowledged"`
}
type UPIDirectPayment struct {
	UPIVerifyCredit
	BillID          int64  `json:"bill_id"`
	PendingClaimID  int64  `json:"pending_claim_id,omitempty"`
	ClaimResolution string `json:"claim_resolution,omitempty"`
	Reason          string `json:"reason,omitempty"`
}
type UPIReason struct {
	Reason string `json:"reason"`
}
type UPIClaim struct {
	BillNumber       string           `json:"bill_number,omitempty"`
	Flat             *MaintenanceFlat `json:"flat,omitempty"`
	ID               int64            `json:"id"`
	SocietyID        int64            `json:"society_id"`
	BillID           int64            `json:"bill_id"`
	PaymentRequestID string           `json:"payment_request_id"`
	SubmittedBy      int64            `json:"submitted_by"`
	Reference        string           `json:"reference"`
	PaymentDate      string           `json:"payment_date"`
	Status           string           `json:"status"`
	AmountPaise      int64            `json:"amount_paise"`
	SettingsVersion  int64            `json:"settings_version"`
	Destination      UPIConfig        `json:"destination"`
	ReviewedBy       *int64           `json:"reviewed_by,omitempty"`
	ReviewedAt       *time.Time       `json:"reviewed_at,omitempty"`
	Reason           *string          `json:"reason,omitempty"`
	CreatedAt        time.Time        `json:"created_at"`
}
type UPIPayment struct {
	EvidenceReference string           `json:"evidence_reference,omitempty"`
	Flat              *MaintenanceFlat `json:"flat,omitempty"`
	ID                int64            `json:"id"`
	SocietyID         int64            `json:"society_id"`
	BillID            int64            `json:"bill_id"`
	BillNumber        string           `json:"bill_number"`
	ClaimID           *int64           `json:"claim_id,omitempty"`
	PayerID           *int64           `json:"payer_id,omitempty"`
	Reference         *string          `json:"reference,omitempty"`
	AmountPaise       int64            `json:"amount_paise"`
	Currency          string           `json:"currency"`
	CreditDate        string           `json:"credit_date"`
	ReceiptNumber     string           `json:"receipt_number"`
	SettingsVersion   int64            `json:"settings_version"`
	Destination       UPIConfig        `json:"destination"`
	Status            string           `json:"status"`
	VerifiedBy        *int64           `json:"verified_by,omitempty"`
	VerifiedByName    string           `json:"verified_by_name,omitempty"`
	VerifiedAt        time.Time        `json:"verified_at"`
	ReversedBy        *int64           `json:"reversed_by,omitempty"`
	ReversedAt        *time.Time       `json:"reversed_at,omitempty"`
	ReversalReason    *string          `json:"reversal_reason,omitempty"`
}
type UPICreateReport struct {
	BillID           int64  `json:"bill_id"`
	PaymentRequestID string `json:"payment_request_id,omitempty"`
	Reference        string `json:"reference"`
	AmountPaise      int64  `json:"amount_paise"`
	PaymentDate      string `json:"payment_date"`
	Explanation      string `json:"explanation"`
}
type UPIUpdateReport struct {
	Status         string `json:"status"`
	ResolutionNote string `json:"resolution_note"`
}
type UPIReport struct {
	ID        int64 `json:"id"`
	SocietyID int64 `json:"society_id"`
	UPICreateReport
	BillNumber     string    `json:"bill_number,omitempty"`
	ReportedBy     int64     `json:"reported_by"`
	Status         string    `json:"status"`
	ResolutionNote *string   `json:"resolution_note,omitempty"`
	UpdatedBy      *int64    `json:"updated_by,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
type UPIListFilter struct {
	BillID, FlatID              int64
	BillingMonth                string
	SocietyID, UserID, BeforeID int64
	Limit                       int32
	Status, Reference           string
	Block, FlatNumber, Search   string
	Resident                    bool
}
type UPIClaimsPage struct {
	Items      []UPIClaim `json:"items"`
	NextCursor *int64     `json:"next_cursor,omitempty"`
}
type UPIPaymentsPage struct {
	Items      []UPIPayment `json:"items"`
	NextCursor *int64       `json:"next_cursor,omitempty"`
}
type UPIReportsPage struct {
	Items      []UPIReport `json:"items"`
	NextCursor *int64      `json:"next_cursor,omitempty"`
}
type UPISettingsPage struct {
	Items      []MaintenancePaymentSettings `json:"items"`
	NextCursor *int64                       `json:"next_cursor,omitempty"`
}
type UPIReferenceHistory struct {
	Claims   UPIClaimsPage   `json:"claims"`
	Payments UPIPaymentsPage `json:"payments"`
}
type UPIClaimDetail struct {
	UPIClaim
	ReferenceHistory UPIReferenceHistory `json:"reference_history"`
}
type UPIReportDetail struct {
	UPIReport
	ReferenceHistory UPIReferenceHistory `json:"reference_history"`
}
type UPICollectionSummary struct {
	BilledCount      int64 `json:"billed_count"`
	BilledPaise      int64 `json:"billed_paise"`
	CollectedPaise   int64 `json:"collected_paise"`
	OutstandingPaise int64 `json:"outstanding_paise"`
	OverduePaise     int64 `json:"overdue_paise"`
	PendingClaims    int64 `json:"pending_claims"`
}
type UPIAuditEvent struct {
	ID        int64          `json:"id"`
	BillID    *int64         `json:"bill_id,omitempty"`
	ActorID   int64          `json:"actor_id"`
	Action    string         `json:"action"`
	EntityID  string         `json:"entity_id"`
	Details   map[string]any `json:"details"`
	CreatedAt time.Time      `json:"created_at"`
}
type UPIAuditPage struct {
	Items      []UPIAuditEvent `json:"items"`
	NextCursor *int64          `json:"next_cursor,omitempty"`
}
