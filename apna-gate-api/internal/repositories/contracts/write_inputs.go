package contracts

import (
	"go-server/internal/models"
	"time"
)

// ActivateSubscriptionInput contains validated values for persistence.
type ActivateSubscriptionInput struct {
	StartsAt time.Time
	EndsAt   time.Time
	Metadata map[string]any
}

// CancelSubscriptionInput contains validated values for persistence.
type CancelSubscriptionInput struct {
	Reason   string
	Metadata map[string]any
}

// CreateTrialSubscriptionInput contains validated values for persistence.
type CreateTrialSubscriptionInput struct {
	StartsAt    time.Time
	TrialEndsAt time.Time
	EndsAt      *time.Time
	Metadata    map[string]any
}

// RenewSubscriptionInput contains validated values for persistence.
type RenewSubscriptionInput struct {
	StartsAt time.Time
	EndsAt   time.Time
	Metadata map[string]any
}

// UpdateFlatInput contains validated values for persistence.
type UpdateFlatInput struct {
	FlatType   *string
	AreaSqft   *string
	Block      *string
	Floor      *string
	FlatNumber *string
	Status     *models.FlatStatus
	IsActive   *bool
	Metadata   map[string]any
}

// UpdateFlatVisitorSettingInput contains validated values for persistence.
type UpdateFlatVisitorSettingInput struct {
	ApprovalRequired            *bool
	DefaultVisitDurationMinutes *int32
	IsEnabled                   *bool
}

// UpdateGuardVisitorEntryInput contains validated values for persistence.
type UpdateGuardVisitorEntryInput struct {
	FlatID           *int64
	FullName         *string
	PhoneNumber      *string
	Email            *string
	PhotoURL         *string
	VehicleNumber    *string
	VehicleType      *models.VisitorVehicleType
	CompanionsCount  *int32
	CompanionDetails []map[string]any
	Notes            *string
}

// UpdatePlanInput contains validated values for persistence.
type UpdatePlanInput struct {
	Name             *string
	Code             *string
	Description      *string
	PriceAmountPaise *int64
	Currency         *string
	BillingCycle     *models.BillingCycle
	MaxFlats         *int32
	MaxAdmins        *int32
	MaxStaff         *int32
	MaxResidents     *int32
	Features         map[string]any
}

// UpdateSocietyInput contains validated values for persistence.
type UpdateSocietyInput struct {
	Name         *string
	Email        *string
	PhoneNumber  *string
	AddressLine1 *string
	AddressLine2 *string
	Landmark     *string
	City         *string
	State        *string
	Pincode      *string
	Country      *string
	TotalFlats   *int32
	TotalBlocks  *int32
	Metadata     map[string]any
}

// UpdateSocietyVisitorSettingsInput contains validated values for persistence.
type UpdateSocietyVisitorSettingsInput struct {
	ApprovalMode                *models.VisitorApprovalMode
	DefaultVisitDurationMinutes *int32
	GracePeriodMinutes          *int32
	QRExpiryMinutes             *int32
	AllowResidentPreApproval    *bool
	AllowPublicQREntry          *bool
	AllowGuardEntry             *bool
	AllowGuardOnBehalfApproval  *bool
	IsActive                    *bool
}

// UpdateUserInput contains validated values for persistence.
type UpdateUserInput struct {
	ID int64

	FirstName   *string
	LastName    *string
	PhoneNumber *string
	AvatarURL   *string
	DateOfBirth *time.Time
	Gender      *string
	Timezone    *string
	Language    *string
}

// VisitorFormInput contains validated values for persistence.
type VisitorFormInput struct {
	FullName           string
	PhoneNumber        *string
	Email              *string
	PhotoURL           *string
	FlatID             int64
	Purpose            models.VisitorPurpose
	VehicleNumber      *string
	VehicleType        *models.VisitorVehicleType
	CompanionsCount    int32
	CompanionDetails   []map[string]any
	ExpectedAt         *time.Time
	ExpectedCheckoutAt *time.Time
	Notes              *string
	Metadata           map[string]any
	DeliveryPartner    *string
	ServiceProvider    *string
}
