package models

import (
	"errors"
	"strings"
	"time"
)

type CreatePlanRequest struct {
	Name             string         `json:"name" validate:"required,max=100"`
	Code             string         `json:"code" validate:"required,max=50"`
	Description      *string        `json:"description,omitempty"`
	PriceAmountPaise int64          `json:"price_amount_paise" validate:"gte=0"`
	Currency         string         `json:"currency" validate:"required,max=10"`
	BillingCycle     BillingCycle   `json:"billing_cycle" validate:"required"`
	MaxFlats         int32          `json:"max_flats" validate:"required,gt=0"`
	MaxAdmins        int32          `json:"max_admins" validate:"gte=0"`
	MaxStaff         int32          `json:"max_staff" validate:"gte=0"`
	MaxResidents     int32          `json:"max_residents" validate:"required,gt=0"`
	Features         map[string]any `json:"features,omitempty"`
}

func (r *CreatePlanRequest) Sanitize() {
	r.Name = strings.TrimSpace(r.Name)
	r.Code = strings.ToUpper(strings.TrimSpace(r.Code))
	r.Currency = strings.ToUpper(strings.TrimSpace(r.Currency))
	r.Description = trimPtr(r.Description)
	if r.Currency == "" {
		r.Currency = "INR"
	}
}

func (r *CreatePlanRequest) Validate() error {
	if r.Name == "" || r.Code == "" {
		return errors.New("name and code are required")
	}
	if !r.BillingCycle.IsValid() {
		return errors.New("invalid billing_cycle")
	}
	if r.PriceAmountPaise < 0 || r.MaxFlats <= 0 || r.MaxAdmins < 0 || r.MaxStaff < 0 || r.MaxResidents <= 0 {
		return errors.New("invalid plan limits or price")
	}
	return nil
}

type UpdatePlanRequest struct {
	Name             *string        `json:"name,omitempty" validate:"omitempty,max=100"`
	Code             *string        `json:"code,omitempty" validate:"omitempty,max=50"`
	Description      *string        `json:"description,omitempty"`
	PriceAmountPaise *int64         `json:"price_amount_paise,omitempty" validate:"omitempty,gte=0"`
	Currency         *string        `json:"currency,omitempty" validate:"omitempty,max=10"`
	BillingCycle     *BillingCycle  `json:"billing_cycle,omitempty"`
	MaxFlats         *int32         `json:"max_flats,omitempty" validate:"omitempty,gt=0"`
	MaxAdmins        *int32         `json:"max_admins,omitempty" validate:"omitempty,gte=0"`
	MaxStaff         *int32         `json:"max_staff,omitempty" validate:"omitempty,gte=0"`
	MaxResidents     *int32         `json:"max_residents,omitempty" validate:"omitempty,gt=0"`
	Features         map[string]any `json:"features,omitempty"`
}

func (r *UpdatePlanRequest) Sanitize() {
	r.Name = trimPtr(r.Name)
	r.Code = trimUpperPtr(r.Code)
	r.Currency = trimUpperPtr(r.Currency)
	r.Description = trimPtr(r.Description)
}

func (r *UpdatePlanRequest) Validate() error {
	if r.Name == nil && r.Code == nil && r.Description == nil && r.PriceAmountPaise == nil &&
		r.Currency == nil && r.BillingCycle == nil && r.MaxFlats == nil && r.MaxAdmins == nil &&
		r.MaxStaff == nil && r.MaxResidents == nil && r.Features == nil {
		return errors.New("at least one field must be provided")
	}
	if r.MaxResidents != nil && *r.MaxResidents <= 0 {
		return errors.New("invalid max_residents")
	}
	if r.BillingCycle != nil && !r.BillingCycle.IsValid() {
		return errors.New("invalid billing_cycle")
	}
	return nil
}

type PlanFilter struct {
	ID            *int64
	Code          *string
	Name          *string
	BillingCycle  *string
	IsActive      *bool
	Search        *string
	MinPricePaise *int64
	MaxPricePaise *int64
	Limit         int32
	Offset        int32
}

type CreateTrialSubscriptionRequest struct {
	StartsAt    time.Time      `json:"starts_at" validate:"required"`
	TrialEndsAt time.Time      `json:"trial_ends_at" validate:"required"`
	EndsAt      *time.Time     `json:"ends_at,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

type ActivateSubscriptionRequest struct {
	StartsAt time.Time      `json:"starts_at" validate:"required"`
	EndsAt   time.Time      `json:"ends_at" validate:"required"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type RenewSubscriptionRequest struct {
	StartsAt time.Time      `json:"starts_at" validate:"required"`
	EndsAt   time.Time      `json:"ends_at" validate:"required"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type CancelSubscriptionRequest struct {
	Reason   string         `json:"reason" validate:"required,max=500"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

func (r *CancelSubscriptionRequest) Sanitize() {
	r.Reason = strings.TrimSpace(r.Reason)
}

type SubscriptionFilter struct {
	ID        *int64
	SocietyID *int64
	PlanID    *int64

	Status       *string
	PlanCode     *string
	BillingCycle *string
	IsActiveOnly *bool

	StartsAfter  *time.Time
	StartsBefore *time.Time
	EndsAfter    *time.Time
	EndsBefore   *time.Time

	ExpiringBefore *time.Time
	ExpiredOnly    *bool

	Search     *string
	SearchMode *string

	Limit  int32
	Offset int32
}
