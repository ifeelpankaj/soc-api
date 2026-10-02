package models

import "time"

type PlanResponse struct {
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

type SocietySubscriptionResponse struct {
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
	CreatedBy          *int64             `json:"created_by,omitempty"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
	SocietyName        *string            `json:"society_name,omitempty"`
	SocietyCode        *string            `json:"society_code,omitempty"`
	CurrentPlanName    *string            `json:"current_plan_name,omitempty"`
	CurrentPlanCode    *string            `json:"current_plan_code,omitempty"`
}

type SubscriptionStatsResponse struct {
	TotalSubscriptions     int64 `json:"total_subscriptions"`
	PendingSubscriptions   int64 `json:"pending_subscriptions"`
	TrialSubscriptions     int64 `json:"trial_subscriptions"`
	ActiveSubscriptions    int64 `json:"active_subscriptions"`
	ExpiredSubscriptions   int64 `json:"expired_subscriptions"`
	CancelledSubscriptions int64 `json:"cancelled_subscriptions"`
}
