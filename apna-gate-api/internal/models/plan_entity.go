package models

import "time"

type BillingCycle string

const (
	BillingCycleMonthly BillingCycle = "monthly"
	BillingCycleYearly  BillingCycle = "yearly"
)

func (b BillingCycle) IsValid() bool {
	switch b {
	case BillingCycleMonthly, BillingCycleYearly:
		return true
	default:
		return false
	}
}

type SubscriptionStatus string

const (
	SubscriptionStatusPending   SubscriptionStatus = "pending"
	SubscriptionStatusTrial     SubscriptionStatus = "trial"
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusExpired   SubscriptionStatus = "expired"
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
)

func (s SubscriptionStatus) IsOpen() bool {
	return s == SubscriptionStatusPending || s == SubscriptionStatusTrial || s == SubscriptionStatusActive
}

type Plan struct {
	ID               int64          `json:"id"`
	Name             string         `json:"name"`
	Code             string         `json:"code"`
	Description      *string        `json:"description,omitempty"`
	PriceAmountPaise int64          `json:"price_amount_paise"`
	Currency         string         `json:"currency"`
	BillingCycle     BillingCycle   `json:"billing_cycle"`
	MaxFlats         int32          `json:"max_flats"`
	MaxAdmins        int32          `json:"max_admins"`
	MaxStaff         int32          `json:"max_staff"`
	MaxResidents     int32          `json:"max_residents"`
	Features         map[string]any `json:"features"`
	IsActive         bool           `json:"is_active"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

func (p *Plan) ToResponse() *PlanResponse {
	if p == nil {
		return nil
	}
	return &PlanResponse{
		ID: p.ID, Name: p.Name, Code: p.Code, Description: p.Description,
		PriceAmountPaise: p.PriceAmountPaise, Currency: p.Currency, BillingCycle: p.BillingCycle,
		MaxFlats: p.MaxFlats, MaxAdmins: p.MaxAdmins, MaxStaff: p.MaxStaff, MaxResidents: p.MaxResidents,
		Features: p.Features, IsActive: p.IsActive, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

type SocietySubscription struct {
	ID                 int64              `json:"id"`
	SocietyID          int64              `json:"society_id"`
	PlanID             *int64             `json:"plan_id,omitempty"`
	Status             SubscriptionStatus `json:"status"`
	StartsAt           *time.Time         `json:"starts_at,omitempty"`
	EndsAt             *time.Time         `json:"ends_at,omitempty"`
	TrialEndsAt        *time.Time         `json:"trial_ends_at,omitempty"`
	PlanName           string             `json:"plan_name"`
	PlanCode           string             `json:"plan_code"`
	PriceAmountPaise   int64              `json:"price_amount_paise"`
	Currency           string             `json:"currency"`
	BillingCycle       BillingCycle       `json:"billing_cycle"`
	MaxFlats           int32              `json:"max_flats"`
	MaxAdmins          int32              `json:"max_admins"`
	MaxStaff           int32              `json:"max_staff"`
	MaxResidents       int32              `json:"max_residents"`
	Features           map[string]any     `json:"features"`
	ActivatedAt        *time.Time         `json:"activated_at,omitempty"`
	ActivatedBy        *int64             `json:"activated_by,omitempty"`
	ExpiredAt          *time.Time         `json:"expired_at,omitempty"`
	CancelledAt        *time.Time         `json:"cancelled_at,omitempty"`
	CancelledBy        *int64             `json:"cancelled_by,omitempty"`
	CancellationReason *string            `json:"cancellation_reason,omitempty"`
	Metadata           map[string]any     `json:"metadata"`
	CreatedBy          *int64             `json:"created_by,omitempty"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
	SocietyName        *string            `json:"society_name,omitempty"`
	SocietyCode        *string            `json:"society_code,omitempty"`
	CurrentPlanName    *string            `json:"current_plan_name,omitempty"`
	CurrentPlanCode    *string            `json:"current_plan_code,omitempty"`
}

func (s *SocietySubscription) ToResponse() *SocietySubscriptionResponse {
	if s == nil {
		return nil
	}
	return &SocietySubscriptionResponse{
		ID: s.ID, SocietyID: s.SocietyID, PlanID: s.PlanID, Status: s.Status,
		StartsAt: s.StartsAt, EndsAt: s.EndsAt, TrialEndsAt: s.TrialEndsAt,
		PlanName: s.PlanName, PlanCode: s.PlanCode, PriceAmountPaise: s.PriceAmountPaise,
		Currency: s.Currency, BillingCycle: s.BillingCycle, MaxFlats: s.MaxFlats,
		MaxAdmins: s.MaxAdmins, MaxStaff: s.MaxStaff, MaxResidents: s.MaxResidents, Features: s.Features,
		ActivatedAt: s.ActivatedAt, ActivatedBy: s.ActivatedBy, ExpiredAt: s.ExpiredAt,
		CancelledAt: s.CancelledAt, CancelledBy: s.CancelledBy, CancellationReason: s.CancellationReason,
		CreatedBy: s.CreatedBy, CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
		SocietyName: s.SocietyName, SocietyCode: s.SocietyCode,
		CurrentPlanName: s.CurrentPlanName, CurrentPlanCode: s.CurrentPlanCode,
	}
}
