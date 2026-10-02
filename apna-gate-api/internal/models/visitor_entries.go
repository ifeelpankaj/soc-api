package models

import (
	"errors"
	"strings"
	"time"
	"unicode"
)

type VisitorStatus string

const (
	VisitorStatusWaitingApproval VisitorStatus = "waiting_approval"
	VisitorStatusApproved        VisitorStatus = "approved"
	VisitorStatusRejected        VisitorStatus = "rejected"
	VisitorStatusCheckedIn       VisitorStatus = "checked_in"
	VisitorStatusCheckedOut      VisitorStatus = "checked_out"
	VisitorStatusCancelled       VisitorStatus = "cancelled"
	VisitorStatusExpired         VisitorStatus = "expired"
	VisitorStatusAutoClosed      VisitorStatus = "auto_closed"
)

func (s VisitorStatus) IsValid() bool {
	switch s {
	case VisitorStatusWaitingApproval, VisitorStatusApproved, VisitorStatusRejected, VisitorStatusCheckedIn,
		VisitorStatusCheckedOut, VisitorStatusCancelled, VisitorStatusExpired, VisitorStatusAutoClosed:
		return true
	default:
		return false
	}
}

type VisitorEntryEventFilter string

const (
	VisitorEntryEventCreated    VisitorEntryEventFilter = "created"
	VisitorEntryEventCheckedIn  VisitorEntryEventFilter = "checked_in"
	VisitorEntryEventCheckedOut VisitorEntryEventFilter = "checked_out"
	VisitorEntryEventExpected   VisitorEntryEventFilter = "expected"
	VisitorEntryEventActivity   VisitorEntryEventFilter = "activity"
)

func (e VisitorEntryEventFilter) IsValid() bool {
	switch e {
	case VisitorEntryEventCreated, VisitorEntryEventCheckedIn, VisitorEntryEventCheckedOut,
		VisitorEntryEventExpected, VisitorEntryEventActivity:
		return true
	default:
		return false
	}
}

type VisitorVehicleType string

const (
	VisitorVehicleTypeBike  VisitorVehicleType = "bike"
	VisitorVehicleTypeCar   VisitorVehicleType = "car"
	VisitorVehicleTypeAuto  VisitorVehicleType = "auto"
	VisitorVehicleTypeCab   VisitorVehicleType = "cab"
	VisitorVehicleTypeTruck VisitorVehicleType = "truck"
	VisitorVehicleTypeOther VisitorVehicleType = "other"
)

func (t VisitorVehicleType) IsValid() bool {
	switch t {
	case VisitorVehicleTypeBike, VisitorVehicleTypeCar, VisitorVehicleTypeAuto, VisitorVehicleTypeCab,
		VisitorVehicleTypeTruck, VisitorVehicleTypeOther:
		return true
	default:
		return false
	}
}

type VisitorEventType string

const (
	VisitorEventTypeCreated               VisitorEventType = "created"
	VisitorEventTypeApproved              VisitorEventType = "approved"
	VisitorEventTypeRejected              VisitorEventType = "rejected"
	VisitorEventTypeCheckedIn             VisitorEventType = "checked_in"
	VisitorEventTypeCheckedOut            VisitorEventType = "checked_out"
	VisitorEventTypeCancelled             VisitorEventType = "cancelled"
	VisitorEventTypeExpired               VisitorEventType = "expired"
	VisitorEventTypeAutoClosed            VisitorEventType = "auto_closed"
	VisitorEventTypeQRGenerated           VisitorEventType = "qr_generated"
	VisitorEventTypeQRUsed                VisitorEventType = "qr_used"
	VisitorEventTypeGuardApprovedOnBehalf VisitorEventType = "guard_approved_on_behalf"
)

type VisitorInviteStatus string

const (
	VisitorInviteStatusActive    VisitorInviteStatus = "active"
	VisitorInviteStatusUsed      VisitorInviteStatus = "used"
	VisitorInviteStatusExpired   VisitorInviteStatus = "expired"
	VisitorInviteStatusCancelled VisitorInviteStatus = "cancelled"
)

func (s VisitorInviteStatus) IsValid() bool {
	switch s {
	case VisitorInviteStatusActive, VisitorInviteStatusUsed, VisitorInviteStatusExpired, VisitorInviteStatusCancelled:
		return true
	default:
		return false
	}
}

type Visitor struct {
	ID          int64          `json:"id"`
	FullName    string         `json:"full_name"`
	PhoneNumber *string        `json:"phone_number,omitempty"`
	Email       *string        `json:"email,omitempty"`
	PhotoURL    *string        `json:"photo_url,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type VisitorInvite struct {
	ID        int64               `json:"id"`
	SocietyID int64               `json:"society_id"`
	FlatID    int64               `json:"flat_id"`
	CreatedBy int64               `json:"created_by"`
	Purpose   VisitorPurpose      `json:"purpose"`
	Status    VisitorInviteStatus `json:"status"`
	ExpiresAt time.Time           `json:"expires_at"`
	UsedAt    *time.Time          `json:"used_at,omitempty"`
	Metadata  map[string]any      `json:"metadata,omitempty"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

type VisitorInviteHistoryFilter struct {
	Status *VisitorInviteStatus
	Search *string
	Limit  int32
	Offset int32
}

type VisitorInviteHistoryItem struct {
	ID        int64               `json:"id"`
	SocietyID int64               `json:"society_id"`
	FlatID    int64               `json:"flat_id"`
	CreatedBy int64               `json:"created_by"`
	Purpose   VisitorPurpose      `json:"purpose"`
	Status    VisitorInviteStatus `json:"status"`
	ExpiresAt time.Time           `json:"expires_at"`
	UsedAt    *time.Time          `json:"used_at,omitempty"`
	Metadata  map[string]any      `json:"metadata,omitempty"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	ShortLink *ShortLinkResponse  `json:"short_link,omitempty"`
	Visitor   *VisitorSummary     `json:"visitor,omitempty"`
	Entry     *VisitorInviteEntry `json:"entry,omitempty"`
}

type VisitorInviteEntry struct {
	ID           int64         `json:"id"`
	Status       VisitorStatus `json:"status"`
	CheckedInAt  *time.Time    `json:"checked_in_at,omitempty"`
	CheckedOutAt *time.Time    `json:"checked_out_at,omitempty"`
	QRExpiresAt  *time.Time    `json:"qr_expires_at,omitempty"`
	QRUsedAt     *time.Time    `json:"qr_used_at,omitempty"`
}

type VisitorInviteHistoryResult struct {
	Invites []*VisitorInviteHistoryItem `json:"invites"`
	Total   int64                       `json:"total"`
	Limit   int32                       `json:"limit"`
	Offset  int32                       `json:"offset"`
}

// PublicVisitorInviteView is the public-safe invite payload for the visitor web form.
type PublicVisitorInviteView struct {
	ID          int64               `json:"id"`
	Purpose     VisitorPurpose      `json:"purpose"`
	Status      VisitorInviteStatus `json:"status"`
	ExpiresAt   time.Time           `json:"expires_at"`
	SocietyName string              `json:"society_name"`
	FlatNumber  string              `json:"flat_number"`
	Block       *string             `json:"block,omitempty"`
	Floor       *string             `json:"floor,omitempty"`
}

type PublicVisitorInviteViewKind string

const (
	PublicVisitorInviteViewForm       PublicVisitorInviteViewKind = "form"
	PublicVisitorInviteViewQR         PublicVisitorInviteViewKind = "qr"
	PublicVisitorInviteViewCheckedIn  PublicVisitorInviteViewKind = "checked_in"
	PublicVisitorInviteViewCheckedOut PublicVisitorInviteViewKind = "checked_out"
	PublicVisitorInviteViewClosed     PublicVisitorInviteViewKind = "closed"
)

// PublicVisitorInvitePageResponse is the public invite page payload for the visitor web app.
type PublicVisitorInvitePageResponse struct {
	Invite *PublicVisitorInviteView    `json:"invite"`
	View   PublicVisitorInviteViewKind `json:"view"`
	Entry  *VisitorEntry               `json:"entry,omitempty"`
	QR     *QRTokenResponse            `json:"qr,omitempty"`
}

type VisitorEntry struct {
	ID                        int64               `json:"id"`
	SocietyID                 int64               `json:"society_id"`
	FlatID                    int64               `json:"flat_id"`
	VisitorID                 int64               `json:"visitor_id"`
	InviteID                  *int64              `json:"invite_id,omitempty"`
	Source                    VisitorEntrySource  `json:"source"`
	Purpose                   VisitorPurpose      `json:"purpose"`
	Status                    VisitorStatus       `json:"status"`
	VehicleNumber             *string             `json:"vehicle_number,omitempty"`
	VehicleType               *VisitorVehicleType `json:"vehicle_type,omitempty"`
	CompanionsCount           int32               `json:"companions_count"`
	CompanionDetails          []map[string]any    `json:"companion_details,omitempty"`
	ExpectedAt                *time.Time          `json:"expected_at,omitempty"`
	ExpectedCheckoutAt        *time.Time          `json:"expected_checkout_at,omitempty"`
	CheckedInAt               *time.Time          `json:"checked_in_at,omitempty"`
	CheckedOutAt              *time.Time          `json:"checked_out_at,omitempty"`
	AutoClosedAt              *time.Time          `json:"auto_closed_at,omitempty"`
	ApprovedBy                *int64              `json:"approved_by,omitempty"`
	ApproverName              *string             `json:"approver_name,omitempty"`
	ApproverPhoneNumber       *string             `json:"approver_phone_number,omitempty"`
	ApprovedAt                *time.Time          `json:"approved_at,omitempty"`
	DeliveryPartner           *string             `json:"delivery_partner,omitempty"`
	ServiceProvider           *string             `json:"service_provider,omitempty"`
	RejectedBy                *int64              `json:"rejected_by,omitempty"`
	HandledByGuardID          *int64              `json:"handled_by_guard_id,omitempty"`
	HandledByGuardName        *string             `json:"handled_by_guard_name,omitempty"`
	HandledByGuardPhoneNumber *string             `json:"handled_by_guard_phone,omitempty"`
	CreatedBy                 *int64              `json:"created_by,omitempty"`
	QRExpiresAt               *time.Time          `json:"qr_expires_at,omitempty"`
	QRUsedAt                  *time.Time          `json:"qr_used_at,omitempty"`
	Notes                     *string             `json:"notes,omitempty"`
	RejectionReason           *string             `json:"rejection_reason,omitempty"`
	Metadata                  map[string]any      `json:"metadata,omitempty"`
	CreatedAt                 time.Time           `json:"created_at"`
	UpdatedAt                 time.Time           `json:"updated_at"`
	Visitor                   *VisitorSummary     `json:"visitor,omitempty"`
	Flat                      *VisitorFlatSummary `json:"flat,omitempty"`
}

type VisitorSummary struct {
	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Email       *string `json:"email,omitempty"`
	PhotoURL    *string `json:"photo_url,omitempty"`
}

type VisitorFlatSummary struct {
	ID         int64   `json:"id"`
	FlatNumber string  `json:"flat_number"`
	Block      *string `json:"block,omitempty"`
	Floor      *string `json:"floor,omitempty"`
}

type VisitorEntryEvent struct {
	ID             int64            `json:"id"`
	VisitorEntryID int64            `json:"visitor_entry_id"`
	SocietyID      int64            `json:"society_id"`
	ActorUserID    *int64           `json:"actor_user_id,omitempty"`
	EventType      VisitorEventType `json:"event_type"`
	Message        *string          `json:"message,omitempty"`
	Metadata       map[string]any   `json:"metadata,omitempty"`
	CreatedAt      time.Time        `json:"created_at"`
}

type QRTokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type VisitorEntryMutationResponse struct {
	Entry            *VisitorEntry    `json:"entry"`
	QR               *QRTokenResponse `json:"qr,omitempty"`
	IdempotentReplay bool             `json:"-"`
}

// IsSocietyWideGuardService identifies flatless guard Service entries without
// changing the behavior of historical flat-linked or non-guard Service entries.
func (e *VisitorEntry) IsSocietyWideGuardService() bool {
	return e != nil && e.Purpose == VisitorPurposeService && e.Source == VisitorEntrySourceGuardEntry && e.FlatID == 0
}

type CreateVisitorInviteRequest struct {
	Purpose   VisitorPurpose `json:"purpose"`
	ExpiresAt *time.Time     `json:"expires_at,omitempty"`
}

func (r *CreateVisitorInviteRequest) Validate() error {
	if r == nil {
		return errors.New("visitor invite request is required")
	}
	if !r.Purpose.IsValid() {
		return errors.New("invalid visitor purpose")
	}
	if r.ExpiresAt != nil && !r.ExpiresAt.After(time.Now()) {
		return errors.New("expires_at must be in the future")
	}
	if r.ExpiresAt != nil && r.ExpiresAt.After(time.Now().Add(maxInviteDuration)) {
		return errors.New("expires_at cannot be more than 7 days in the future")
	}
	return nil
}

type VisitorFormRequest struct {
	FullName           string              `json:"full_name"`
	PhoneNumber        *string             `json:"phone_number,omitempty"`
	Email              *string             `json:"email,omitempty"`
	PhotoURL           *string             `json:"photo_url,omitempty"`
	FlatID             int64               `json:"flat_id,omitempty"`
	Purpose            VisitorPurpose      `json:"purpose,omitempty"`
	VehicleNumber      *string             `json:"vehicle_number,omitempty"`
	VehicleType        *VisitorVehicleType `json:"vehicle_type,omitempty"`
	CompanionsCount    int32               `json:"companions_count,omitempty"`
	CompanionDetails   []map[string]any    `json:"companion_details,omitempty"`
	ExpectedAt         *time.Time          `json:"expected_at,omitempty"`
	ExpectedCheckoutAt *time.Time          `json:"expected_checkout_at,omitempty"`
	Notes              *string             `json:"notes,omitempty"`
	Metadata           map[string]any      `json:"metadata,omitempty"`
	DeliveryPartner    *string             `json:"delivery_partner,omitempty"`
	ServiceProvider    *string             `json:"service_provider,omitempty"`
}

func (r *VisitorFormRequest) Validate(requireFlatAndPurpose bool) error {
	if r == nil {
		return errors.New("visitor form request is required")
	}
	if r.PhotoURL != nil && strings.TrimSpace(*r.PhotoURL) != "" {
		return ErrVisitorPhotoUploadRequired
	}
	r.PhotoURL = nil
	r.FullName = normalizeVisitorFullName(r.FullName)
	if r.FullName == "" {
		return errors.New("full_name is required")
	}
	r.PhoneNumber = cleanPtr(r.PhoneNumber)
	r.Email = cleanPtr(r.Email)
	if r.PhoneNumber == nil && r.Email == nil {
		return errors.New("phone_number or email is required")
	}
	if r.PhoneNumber != nil && !isTenDigitPhone(*r.PhoneNumber) {
		return errors.New("phone_number must contain exactly 10 digits")
	}
	r.PhotoURL = cleanPtr(r.PhotoURL)
	r.VehicleNumber = cleanPtr(r.VehicleNumber)
	r.Notes = cleanPtr(r.Notes)
	if requireFlatAndPurpose {
		if !r.Purpose.IsValid() {
			return errors.New("invalid visitor purpose")
		}
		if r.Purpose != VisitorPurposeStaff && r.FlatID <= 0 {
			return errors.New("flat_id must be a positive integer")
		}
	}
	if r.VehicleType != nil && !r.VehicleType.IsValid() {
		return errors.New("invalid vehicle_type")
	}
	if r.CompanionsCount < 0 {
		return errors.New("companions_count must be zero or positive")
	}
	if r.ExpectedAt != nil && r.ExpectedCheckoutAt != nil && r.ExpectedCheckoutAt.Before(*r.ExpectedAt) {
		return errors.New("expected_checkout_at must be after expected_at")
	}
	return nil
}

func (r *VisitorFormRequest) ValidateForPurpose() error {
	if r == nil {
		return errors.New("visitor form request is required")
	}
	r.DeliveryPartner = cleanPtr(r.DeliveryPartner)
	r.ServiceProvider = cleanPtr(r.ServiceProvider)

	switch r.Purpose {
	case VisitorPurposeGuest:
		if r.PhoneNumber == nil {
			return errors.New("phone_number is required for guest entries")
		}
	case VisitorPurposeDelivery:
		if r.PhoneNumber == nil {
			return errors.New("phone_number is required for delivery entries")
		}
		if r.DeliveryPartner == nil {
			return errors.New("delivery_partner is required for delivery entries")
		}
		r.CompanionsCount = 0
	case VisitorPurposeCab:
		if r.VehicleNumber == nil {
			return errors.New("vehicle_number is required for cab entries")
		}
		if r.VehicleType == nil {
			return errors.New("vehicle_type is required for cab entries")
		}
	case VisitorPurposeService:
		if r.PhoneNumber == nil {
			return errors.New("phone_number is required for service entries")
		}
		if r.ServiceProvider == nil {
			return errors.New("service_provider is required for service entries")
		}
		r.CompanionsCount = 0
	case VisitorPurposeMaintenance:
		if r.PhoneNumber == nil {
			return errors.New("phone_number is required for maintenance entries")
		}
		if r.ServiceProvider == nil {
			return errors.New("service_provider is required for maintenance entries")
		}
		r.CompanionsCount = 0
	case VisitorPurposeStaff:
		if r.PhoneNumber == nil {
			return errors.New("phone_number is required for staff entries")
		}
		r.FlatID = 0
	}
	return nil
}

func (r *VisitorFormRequest) ValidateInviteSubmit() error {
	if r == nil {
		return errors.New("visitor form request is required")
	}
	if err := r.Validate(false); err != nil {
		return err
	}
	if r.ExpectedAt == nil {
		return errors.New("expected_at is required for invite submissions")
	}
	return nil
}

type GuardApproveEntryRequest struct {
	OnBehalf *bool   `json:"on_behalf,omitempty"`
	Reason   *string `json:"reason,omitempty"`
}

type UpdateGuardVisitorEntryRequest struct {
	FlatID           *int64              `json:"flat_id,omitempty"`
	FullName         *string             `json:"full_name,omitempty"`
	PhoneNumber      *string             `json:"phone_number,omitempty"`
	Email            *string             `json:"email,omitempty"`
	PhotoURL         *string             `json:"photo_url,omitempty"`
	VehicleNumber    *string             `json:"vehicle_number,omitempty"`
	VehicleType      *VisitorVehicleType `json:"vehicle_type,omitempty"`
	CompanionsCount  *int32              `json:"companions_count,omitempty"`
	CompanionDetails []map[string]any    `json:"companion_details,omitempty"`
	Notes            *string             `json:"notes,omitempty"`
}

func (r *UpdateGuardVisitorEntryRequest) Validate() error {
	if r == nil {
		return errors.New("update request is required")
	}
	if r.PhotoURL != nil && strings.TrimSpace(*r.PhotoURL) != "" {
		return ErrVisitorPhotoUploadRequired
	}
	r.PhotoURL = nil
	if r.FlatID == nil &&
		r.FullName == nil &&
		r.PhoneNumber == nil &&
		r.Email == nil &&
		r.PhotoURL == nil &&
		r.VehicleNumber == nil &&
		r.VehicleType == nil &&
		r.CompanionsCount == nil &&
		r.CompanionDetails == nil &&
		r.Notes == nil {
		return errors.New("at least one field must be provided for update")
	}
	if r.FlatID != nil && *r.FlatID <= 0 {
		return errors.New("flat_id must be positive")
	}
	if r.FullName != nil {
		fullName := normalizeVisitorFullName(*r.FullName)
		if fullName == "" {
			return errors.New("full_name cannot be empty")
		}
		r.FullName = &fullName
	}
	if r.PhoneNumber != nil {
		phoneNumber := strings.TrimSpace(*r.PhoneNumber)
		if !isTenDigitPhone(phoneNumber) {
			return errors.New("phone_number must contain exactly 10 digits")
		}
		r.PhoneNumber = &phoneNumber
	}
	if r.VehicleType != nil && !r.VehicleType.IsValid() {
		return errors.New("invalid vehicle_type")
	}
	if r.CompanionsCount != nil && *r.CompanionsCount < 0 {
		return errors.New("companions_count must be zero or positive")
	}
	return nil
}

type WaitingAtGateFilter struct {
	SocietyID int64
	Search    *string
	Limit     int32
	Offset    int32
}

type ExpectedGuestFilter struct {
	SocietyID int64
	Search    *string
	FromAt    time.Time
	ToAt      time.Time
	Limit     int32
	Offset    int32
}

type RejectVisitorEntryRequest struct {
	Reason string `json:"reason"`
}

func (r *RejectVisitorEntryRequest) Validate() error {
	if r == nil || strings.TrimSpace(r.Reason) == "" {
		return errors.New("reason is required")
	}
	r.Reason = strings.TrimSpace(r.Reason)
	return nil
}

type QRTokenRequest struct {
	Token string `json:"token"`
}

func (r *QRTokenRequest) Validate() error {
	if r == nil || strings.TrimSpace(r.Token) == "" {
		return errors.New("token is required")
	}
	r.Token = strings.TrimSpace(r.Token)
	return nil
}

type VisitorEntryFilter struct {
	SocietyID   int64
	FlatID      *int64
	Status      *VisitorStatus
	Source      *VisitorEntrySource
	Purpose     *VisitorPurpose
	Block       *string
	Event       *VisitorEntryEventFilter
	EventFrom   *time.Time
	EventTo     *time.Time
	CreatedFrom *time.Time
	CreatedTo   *time.Time
	Search      *string
	Limit       int32
	Offset      int32
}

type VisitorEntryListResult struct {
	Entries []*VisitorEntry `json:"entries"`
	Total   int64           `json:"total"`
	Limit   int32           `json:"limit"`
	Offset  int32           `json:"offset"`
}

type VisitorEntryStatsResponse struct {
	TodayVisitors     int64 `json:"today_visitors"`
	VisitorsInside    int64 `json:"visitors_inside"`
	PendingApprovals  int64 `json:"pending_approvals"`
	CheckedOutToday   int64 `json:"checked_out_today"`
	CheckedOutInRange int64 `json:"checked_out_in_range,omitempty"`
	RejectedToday     int64 `json:"rejected_today"`
	AutoClosedToday   int64 `json:"auto_closed_today"`
}

type VisitorDailyCountResponse struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type VisitorEntryDailyStatsResponse struct {
	Days     int32                       `json:"days"`
	Metric   string                      `json:"metric"`
	Timezone string                      `json:"timezone"`
	Daily    []VisitorDailyCountResponse `json:"daily"`
	Total    int64                       `json:"total"`
}

type VisitorPendingEntry struct {
	*VisitorEntry
	WaitingSince        time.Time `json:"waiting_since"`
	PrimaryResidentName *string   `json:"primary_resident_name,omitempty"`
	PrimaryResidentID   *int64    `json:"primary_resident_id,omitempty"`
}

type VisitorPendingListResult struct {
	Entries []*VisitorPendingEntry `json:"entries"`
	Total   int64                  `json:"total"`
	Limit   int32                  `json:"limit"`
	Offset  int32                  `json:"offset"`
}

type VisitorPendingFilter struct {
	SocietyID int64
	FlatID    *int64
	Block     *string
	Limit     int32
	Offset    int32
}

type MemberVisitorApprovalStatsResponse struct {
	ApprovedCount int64 `json:"approved_count"`
	RejectedCount int64 `json:"rejected_count"`
}

type FlatVisitorContextResponse struct {
	OccupancyStatus     FlatStatus                    `json:"occupancy_status"`
	PrimaryResident     *FlatVisitorContextResident   `json:"primary_resident,omitempty"`
	TotalResidents      int64                         `json:"total_residents"`
	InheritsSocietyMode bool                          `json:"inherits_society_mode"`
	SocietyApprovalMode VisitorApprovalMode           `json:"society_approval_mode"`
	VisitorSettings     []FlatVisitorSettingsResponse `json:"visitor_settings"`
	RecentVisitors      []*FlatRecentVisitorSummary   `json:"recent_visitors"`
}

type FlatVisitorContextResident struct {
	ID       int64  `json:"id"`
	UserID   int64  `json:"user_id"`
	FullName string `json:"full_name"`
}

type FlatRecentVisitorSummary struct {
	EntryID   int64          `json:"entry_id"`
	FullName  string         `json:"full_name"`
	Purpose   VisitorPurpose `json:"purpose"`
	Status    VisitorStatus  `json:"status"`
	VisitedOn time.Time      `json:"visited_on"`
}

type GuardDeskVisitorSettingsSummary struct {
	AllowGuardEntry            bool `json:"allow_guard_entry"`
	AllowGuardOnBehalfApproval bool `json:"allow_guard_on_behalf_approval"`
}

type GuardDeskBootstrapResponse struct {
	Society             *SocietyResponse                 `json:"society"`
	Stats               *VisitorEntryStatsResponse       `json:"stats"`
	ExpectedGuestsCount int64                            `json:"expected_guests_count"`
	WaitingAtGateCount  int64                            `json:"waiting_at_gate_count"`
	PendingPreview      []*VisitorPendingEntry           `json:"pending_preview"`
	VisitorSettings     *GuardDeskVisitorSettingsSummary `json:"visitor_settings,omitempty"`
}

type SocietyFlatVisitorSettingRow struct {
	FlatID                      int64          `json:"flat_id"`
	FlatNumber                  string         `json:"flat_number"`
	Block                       *string        `json:"block,omitempty"`
	Purpose                     VisitorPurpose `json:"purpose"`
	ApprovalRequired            bool           `json:"approval_required"`
	IsEnabled                   bool           `json:"is_enabled"`
	DefaultVisitDurationMinutes *int32         `json:"default_visit_duration_minutes,omitempty"`
}

type SocietyFlatVisitorSettingsFilter struct {
	SocietyID int64
	FlatID    *int64
	Block     *string
	Purpose   *VisitorPurpose
	Limit     int32
	Offset    int32
}

type SocietyFlatVisitorSettingsListResult struct {
	Settings []*SocietyFlatVisitorSettingRow `json:"settings"`
	Total    int64                           `json:"total"`
	Limit    int32                           `json:"limit"`
	Offset   int32                           `json:"offset"`
}

type VisitorEntryOptionsResponse struct {
	Purposes []VisitorPurpose           `json:"purposes"`
	Blocks   []VisitorEntryOptionsBlock `json:"blocks"`
	Flats    []VisitorEntryOptionsFlat  `json:"flats"`
	HasMore  bool                       `json:"has_more,omitempty"`
}

const maxInviteDuration = 7 * 24 * time.Hour

func (e *VisitorEntry) ToScanPreview() *VisitorEntry {
	if e == nil {
		return nil
	}
	preview := &VisitorEntry{
		ID:          e.ID,
		SocietyID:   e.SocietyID,
		FlatID:      e.FlatID,
		Status:      e.Status,
		Purpose:     e.Purpose,
		QRExpiresAt: e.QRExpiresAt,
	}
	if e.Visitor != nil {
		preview.Visitor = &VisitorSummary{FullName: e.Visitor.FullName}
	}
	if e.Flat != nil {
		preview.Flat = &VisitorFlatSummary{
			ID:         e.Flat.ID,
			FlatNumber: e.Flat.FlatNumber,
			Block:      e.Flat.Block,
			Floor:      e.Flat.Floor,
		}
	}
	return preview
}

type VisitorEntryOptionsBlock struct {
	Block *string                   `json:"block,omitempty"`
	Flats []VisitorEntryOptionsFlat `json:"flats"`
}

type VisitorEntryOptionsFlat struct {
	ID         int64   `json:"id"`
	Block      *string `json:"block,omitempty"`
	Floor      *string `json:"floor,omitempty"`
	FlatNumber string  `json:"flat_number"`
}

func cleanPtr(value *string) *string {
	if value == nil {
		return nil
	}
	cleaned := strings.TrimSpace(*value)
	if cleaned == "" {
		return nil
	}
	return &cleaned
}

func normalizeVisitorFullName(value string) string {
	words := strings.Fields(value)
	for index, word := range words {
		runes := []rune(word)
		for runeIndex := range runes {
			if runeIndex == 0 {
				runes[runeIndex] = unicode.ToUpper(runes[runeIndex])
			} else {
				runes[runeIndex] = unicode.ToLower(runes[runeIndex])
			}
		}
		words[index] = string(runes)
	}
	return strings.Join(words, " ")
}

func isTenDigitPhone(value string) bool {
	if len(value) != 10 {
		return false
	}
	for _, digit := range value {
		if digit < '0' || digit > '9' {
			return false
		}
	}
	return true
}
