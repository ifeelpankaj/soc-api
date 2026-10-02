package models

import "time"

type SocietyStatus string

const (
	SocietyStatusPending   SocietyStatus = "pending"
	SocietyStatusActive    SocietyStatus = "active"
	SocietyStatusSuspended SocietyStatus = "suspended"
	SocietyStatusRejected  SocietyStatus = "rejected"
)

func (s SocietyStatus) IsValid() bool {
	switch s {
	case SocietyStatusPending, SocietyStatusActive, SocietyStatusSuspended, SocietyStatusRejected:
		return true
	default:
		return false
	}
}

type SocietyMemberRole string

const (
	SocietyMemberRoleOwner    SocietyMemberRole = "owner"
	SocietyMemberRoleAdmin    SocietyMemberRole = "admin"
	SocietyMemberRoleStaff    SocietyMemberRole = "staff"
	SocietyMemberRoleResident SocietyMemberRole = "resident"
)

func (r SocietyMemberRole) IsValid() bool {
	switch r {
	case SocietyMemberRoleOwner, SocietyMemberRoleAdmin, SocietyMemberRoleStaff, SocietyMemberRoleResident:
		return true
	default:
		return false
	}
}

type SocietyMemberStatus string

const (
	SocietyMemberStatusPending   SocietyMemberStatus = "pending"
	SocietyMemberStatusActive    SocietyMemberStatus = "active"
	SocietyMemberStatusSuspended SocietyMemberStatus = "suspended"
	SocietyMemberStatusRemoved   SocietyMemberStatus = "removed"
)

func (s SocietyMemberStatus) IsValid() bool {
	switch s {
	case SocietyMemberStatusPending, SocietyMemberStatusActive, SocietyMemberStatusSuspended, SocietyMemberStatusRemoved:
		return true
	default:
		return false
	}
}

type Society struct {
	ID               int64          `json:"id" db:"id"`
	Name             string         `json:"name" db:"name"`
	SocietyCode      string         `json:"society_code" db:"society_code"`
	Email            *string        `json:"email,omitempty" db:"email"`
	PhoneNumber      *string        `json:"phone_number,omitempty" db:"phone_number"`
	AddressLine1     *string        `json:"address_line1,omitempty" db:"address_line1"`
	AddressLine2     *string        `json:"address_line2,omitempty" db:"address_line2"`
	Landmark         *string        `json:"landmark,omitempty" db:"landmark"`
	City             *string        `json:"city,omitempty" db:"city"`
	State            *string        `json:"state,omitempty" db:"state"`
	Pincode          *string        `json:"pincode,omitempty" db:"pincode"`
	Country          string         `json:"country" db:"country"`
	TotalFlats       int32          `json:"total_flats" db:"total_flats"`
	TotalBlocks      int32          `json:"total_blocks" db:"total_blocks"`
	Status           SocietyStatus  `json:"status" db:"status"`
	CreatedBy        int64          `json:"created_by" db:"created_by"`
	ApprovedBy       *int64         `json:"approved_by,omitempty" db:"approved_by"`
	ApprovedAt       *time.Time     `json:"approved_at,omitempty" db:"approved_at"`
	RejectedBy       *int64         `json:"rejected_by,omitempty" db:"rejected_by"`
	RejectedAt       *time.Time     `json:"rejected_at,omitempty" db:"rejected_at"`
	RejectionReason  *string        `json:"rejection_reason,omitempty" db:"rejection_reason"`
	SuspendedBy      *int64         `json:"suspended_by,omitempty" db:"suspended_by"`
	SuspendedAt      *time.Time     `json:"suspended_at,omitempty" db:"suspended_at"`
	SuspensionReason *string        `json:"suspension_reason,omitempty" db:"suspension_reason"`
	Metadata         map[string]any `json:"metadata" db:"metadata"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" db:"updated_at"`
	DeletedAt        *time.Time     `json:"deleted_at,omitempty" db:"deleted_at"`
}

func (s *Society) ToResponse() *SocietyResponse {
	if s == nil {
		return nil
	}
	return &SocietyResponse{
		ID: s.ID, Name: s.Name, SocietyCode: s.SocietyCode, Email: s.Email,
		PhoneNumber: s.PhoneNumber, AddressLine1: s.AddressLine1, AddressLine2: s.AddressLine2,
		Landmark: s.Landmark, City: s.City, State: s.State, Pincode: s.Pincode, Country: s.Country,
		TotalFlats: s.TotalFlats, TotalBlocks: s.TotalBlocks, Status: s.Status, CreatedBy: s.CreatedBy,
		ApprovedBy: s.ApprovedBy, ApprovedAt: s.ApprovedAt, RejectedBy: s.RejectedBy, RejectedAt: s.RejectedAt,
		RejectionReason: s.RejectionReason, SuspendedBy: s.SuspendedBy, SuspendedAt: s.SuspendedAt,
		SuspensionReason: s.SuspensionReason, CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
	}
}

type SocietyMember struct {
	ID           int64               `json:"id" db:"id"`
	SocietyID    int64               `json:"society_id" db:"society_id"`
	UserID       int64               `json:"user_id" db:"user_id"`
	Role         SocietyMemberRole   `json:"role" db:"role"`
	Status       SocietyMemberStatus `json:"status" db:"status"`
	InvitedBy    *int64              `json:"invited_by,omitempty" db:"invited_by"`
	JoinedAt     time.Time           `json:"joined_at" db:"joined_at"`
	RemovedBy    *int64              `json:"removed_by,omitempty" db:"removed_by"`
	RemovedAt    *time.Time          `json:"removed_at,omitempty" db:"removed_at"`
	RemoveReason *string             `json:"remove_reason,omitempty" db:"remove_reason"`
	UserFullName *string             `json:"user_full_name,omitempty"`
	UserEmail    *string             `json:"user_email,omitempty"`
	UserPhone    *string             `json:"user_phone,omitempty"`
	Metadata     map[string]any      `json:"metadata" db:"metadata"`
	CreatedAt    time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at" db:"updated_at"`
}

func (m *SocietyMember) ToResponse() *SocietyMemberResponse {
	if m == nil {
		return nil
	}
	return &SocietyMemberResponse{
		ID: m.ID, SocietyID: m.SocietyID, UserID: m.UserID, Role: m.Role, Status: m.Status,
		InvitedBy: m.InvitedBy, JoinedAt: m.JoinedAt, RemovedBy: m.RemovedBy, RemovedAt: m.RemovedAt,
		RemoveReason: m.RemoveReason, UserFullName: m.UserFullName, UserEmail: m.UserEmail,
		UserPhone: m.UserPhone, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}
