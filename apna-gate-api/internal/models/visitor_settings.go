package models

import (
	"errors"
	"time"
)

type VisitorApprovalMode string

const (
	VisitorApprovalModeMandatory VisitorApprovalMode = "mandatory"
	VisitorApprovalModeOptional  VisitorApprovalMode = "optional"
	VisitorApprovalModeHybrid    VisitorApprovalMode = "hybrid"
)

func (m VisitorApprovalMode) IsValid() bool {
	switch m {
	case VisitorApprovalModeMandatory, VisitorApprovalModeOptional, VisitorApprovalModeHybrid:
		return true
	default:
		return false
	}
}

type VisitorPurpose string

const (
	VisitorPurposeGuest       VisitorPurpose = "guest"
	VisitorPurposeDelivery    VisitorPurpose = "delivery"
	VisitorPurposeCab         VisitorPurpose = "cab"
	VisitorPurposeService     VisitorPurpose = "service"
	VisitorPurposeMaintenance VisitorPurpose = "maintenance"
	VisitorPurposeStaff       VisitorPurpose = "staff"
	VisitorPurposeOther       VisitorPurpose = "other"
)

func (p VisitorPurpose) IsValid() bool {
	switch p {
	case VisitorPurposeGuest, VisitorPurposeDelivery, VisitorPurposeCab, VisitorPurposeService,
		VisitorPurposeMaintenance, VisitorPurposeStaff, VisitorPurposeOther:
		return true
	default:
		return false
	}
}

type VisitorEntrySource string

const (
	VisitorEntrySourceResidentLink VisitorEntrySource = "resident_link"
	VisitorEntrySourcePublicQR     VisitorEntrySource = "public_qr"
	VisitorEntrySourceGuardEntry   VisitorEntrySource = "guard_entry"
	VisitorEntrySourceQuickLink    VisitorEntrySource = "quick_link"
)

func (s VisitorEntrySource) IsValid() bool {
	switch s {
	case VisitorEntrySourceResidentLink, VisitorEntrySourcePublicQR, VisitorEntrySourceGuardEntry, VisitorEntrySourceQuickLink:
		return true
	default:
		return false
	}
}

type SocietyVisitorSettings struct {
	ID                          int64               `json:"id"`
	SocietyID                   int64               `json:"society_id"`
	ApprovalMode                VisitorApprovalMode `json:"approval_mode"`
	DefaultVisitDurationMinutes int32               `json:"default_visit_duration_minutes"`
	GracePeriodMinutes          int32               `json:"grace_period_minutes"`
	QRExpiryMinutes             int32               `json:"qr_expiry_minutes"`
	AllowResidentPreApproval    bool                `json:"allow_resident_pre_approval"`
	AllowPublicQREntry          bool                `json:"allow_public_qr_entry"`
	AllowGuardEntry             bool                `json:"allow_guard_entry"`
	AllowGuardOnBehalfApproval  bool                `json:"allow_guard_on_behalf_approval"`
	IsActive                    bool                `json:"is_active"`
	UpdatedBy                   *int64              `json:"updated_by,omitempty"`
	CreatedAt                   time.Time           `json:"created_at"`
	UpdatedAt                   time.Time           `json:"updated_at"`
}

func (s *SocietyVisitorSettings) ToResponse() *SocietyVisitorSettingsResponse {
	if s == nil {
		return nil
	}
	return &SocietyVisitorSettingsResponse{
		ID: s.ID, SocietyID: s.SocietyID, ApprovalMode: s.ApprovalMode,
		DefaultVisitDurationMinutes: s.DefaultVisitDurationMinutes, GracePeriodMinutes: s.GracePeriodMinutes,
		QRExpiryMinutes: s.QRExpiryMinutes, AllowResidentPreApproval: s.AllowResidentPreApproval,
		AllowPublicQREntry: s.AllowPublicQREntry, AllowGuardEntry: s.AllowGuardEntry,
		AllowGuardOnBehalfApproval: s.AllowGuardOnBehalfApproval, IsActive: s.IsActive,
		UpdatedBy: s.UpdatedBy, CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
	}
}

type FlatVisitorSettings struct {
	ID                          int64          `json:"id"`
	SocietyID                   int64          `json:"society_id"`
	FlatID                      int64          `json:"flat_id"`
	Purpose                     VisitorPurpose `json:"purpose"`
	ApprovalRequired            bool           `json:"approval_required"`
	DefaultVisitDurationMinutes *int32         `json:"default_visit_duration_minutes,omitempty"`
	IsEnabled                   bool           `json:"is_enabled"`
	UpdatedBy                   *int64         `json:"updated_by,omitempty"`
	CreatedAt                   time.Time      `json:"created_at"`
	UpdatedAt                   time.Time      `json:"updated_at"`
}

func (s *FlatVisitorSettings) ToResponse() *FlatVisitorSettingsResponse {
	if s == nil {
		return nil
	}
	return &FlatVisitorSettingsResponse{
		ID: s.ID, SocietyID: s.SocietyID, FlatID: s.FlatID, Purpose: s.Purpose,
		ApprovalRequired: s.ApprovalRequired, DefaultVisitDurationMinutes: s.DefaultVisitDurationMinutes,
		IsEnabled: s.IsEnabled, UpdatedBy: s.UpdatedBy, CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
	}
}

type SocietyVisitorSettingsResponse struct {
	ID                          int64               `json:"id"`
	SocietyID                   int64               `json:"society_id"`
	ApprovalMode                VisitorApprovalMode `json:"approval_mode"`
	DefaultVisitDurationMinutes int32               `json:"default_visit_duration_minutes"`
	GracePeriodMinutes          int32               `json:"grace_period_minutes"`
	QRExpiryMinutes             int32               `json:"qr_expiry_minutes"`
	AllowResidentPreApproval    bool                `json:"allow_resident_pre_approval"`
	AllowPublicQREntry          bool                `json:"allow_public_qr_entry"`
	AllowGuardEntry             bool                `json:"allow_guard_entry"`
	AllowGuardOnBehalfApproval  bool                `json:"allow_guard_on_behalf_approval"`
	IsActive                    bool                `json:"is_active"`
	UpdatedBy                   *int64              `json:"updated_by,omitempty"`
	CreatedAt                   time.Time           `json:"created_at"`
	UpdatedAt                   time.Time           `json:"updated_at"`
}

type FlatVisitorSettingsResponse struct {
	ID                          int64          `json:"id"`
	SocietyID                   int64          `json:"society_id"`
	FlatID                      int64          `json:"flat_id"`
	Purpose                     VisitorPurpose `json:"purpose"`
	ApprovalRequired            bool           `json:"approval_required"`
	DefaultVisitDurationMinutes *int32         `json:"default_visit_duration_minutes,omitempty"`
	IsEnabled                   bool           `json:"is_enabled"`
	UpdatedBy                   *int64         `json:"updated_by,omitempty"`
	CreatedAt                   time.Time      `json:"created_at"`
	UpdatedAt                   time.Time      `json:"updated_at"`
}

type UpdateSocietyVisitorSettingsRequest struct {
	ApprovalMode                *VisitorApprovalMode `json:"approval_mode,omitempty"`
	DefaultVisitDurationMinutes *int32               `json:"default_visit_duration_minutes,omitempty"`
	GracePeriodMinutes          *int32               `json:"grace_period_minutes,omitempty"`
	QRExpiryMinutes             *int32               `json:"qr_expiry_minutes,omitempty"`
	AllowResidentPreApproval    *bool                `json:"allow_resident_pre_approval,omitempty"`
	AllowPublicQREntry          *bool                `json:"allow_public_qr_entry,omitempty"`
	AllowGuardEntry             *bool                `json:"allow_guard_entry,omitempty"`
	AllowGuardOnBehalfApproval  *bool                `json:"allow_guard_on_behalf_approval,omitempty"`
	IsActive                    *bool                `json:"is_active,omitempty"`
}

func (r *UpdateSocietyVisitorSettingsRequest) Validate() error {
	if r == nil {
		return errors.New("visitor settings request is required")
	}
	if r.ApprovalMode == nil && r.DefaultVisitDurationMinutes == nil && r.GracePeriodMinutes == nil &&
		r.QRExpiryMinutes == nil && r.AllowResidentPreApproval == nil && r.AllowPublicQREntry == nil &&
		r.AllowGuardEntry == nil && r.AllowGuardOnBehalfApproval == nil && r.IsActive == nil {
		return errors.New("at least one field must be provided")
	}
	if r.ApprovalMode != nil && !r.ApprovalMode.IsValid() {
		return errors.New("invalid approval_mode")
	}
	if r.DefaultVisitDurationMinutes != nil && *r.DefaultVisitDurationMinutes <= 0 {
		return errors.New("default_visit_duration_minutes must be positive")
	}
	if r.GracePeriodMinutes != nil && *r.GracePeriodMinutes < 0 {
		return errors.New("grace_period_minutes must be zero or positive")
	}
	if r.QRExpiryMinutes != nil && *r.QRExpiryMinutes <= 0 {
		return errors.New("qr_expiry_minutes must be positive")
	}
	return nil
}

type UpdateFlatVisitorSettingRequest struct {
	ApprovalRequired            *bool  `json:"approval_required,omitempty"`
	DefaultVisitDurationMinutes *int32 `json:"default_visit_duration_minutes,omitempty"`
	IsEnabled                   *bool  `json:"is_enabled,omitempty"`
}

func (r *UpdateFlatVisitorSettingRequest) Validate() error {
	if r == nil {
		return errors.New("flat visitor setting request is required")
	}
	if r.ApprovalRequired == nil && r.DefaultVisitDurationMinutes == nil && r.IsEnabled == nil {
		return errors.New("at least one field must be provided")
	}
	if r.DefaultVisitDurationMinutes != nil && *r.DefaultVisitDurationMinutes <= 0 {
		return errors.New("default_visit_duration_minutes must be positive")
	}
	return nil
}
