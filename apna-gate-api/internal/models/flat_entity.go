package models

import "time"

type FlatStatus string

const (
	FlatStatusVacant   FlatStatus = "vacant"
	FlatStatusOccupied FlatStatus = "occupied"
	FlatStatusBlocked  FlatStatus = "blocked"
)

func (s FlatStatus) IsValid() bool {
	switch s {
	case FlatStatusVacant, FlatStatusOccupied, FlatStatusBlocked:
		return true
	default:
		return false
	}
}

type FlatResidentRole string

const (
	FlatResidentRoleOwner  FlatResidentRole = "owner"
	FlatResidentRoleTenant FlatResidentRole = "tenant"
	FlatResidentRoleFamily FlatResidentRole = "family"
)

func (r FlatResidentRole) IsValid() bool {
	switch r {
	case FlatResidentRoleOwner, FlatResidentRoleTenant, FlatResidentRoleFamily:
		return true
	default:
		return false
	}
}

type FlatResidentStatus string

const (
	FlatResidentStatusActive   FlatResidentStatus = "active"
	FlatResidentStatusInactive FlatResidentStatus = "inactive"
	FlatResidentStatusMovedOut FlatResidentStatus = "moved_out"
)

func (s FlatResidentStatus) IsValid() bool {
	switch s {
	case FlatResidentStatusActive, FlatResidentStatusInactive, FlatResidentStatusMovedOut:
		return true
	default:
		return false
	}
}

type FlatClaimStatus string

const (
	FlatClaimStatusPending   FlatClaimStatus = "pending"
	FlatClaimStatusApproved  FlatClaimStatus = "approved"
	FlatClaimStatusRejected  FlatClaimStatus = "rejected"
	FlatClaimStatusCancelled FlatClaimStatus = "cancelled"
)

func (s FlatClaimStatus) IsValid() bool {
	switch s {
	case FlatClaimStatusPending, FlatClaimStatusApproved, FlatClaimStatusRejected, FlatClaimStatusCancelled:
		return true
	default:
		return false
	}
}

type Flat struct {
	FlatType            *string        `json:"flat_type,omitempty"`
	AreaSqft            *string        `json:"area_sqft,omitempty"`
	ID                  int64          `json:"id" db:"id"`
	SocietyID           int64          `json:"society_id" db:"society_id"`
	Block               *string        `json:"block,omitempty" db:"block"`
	Floor               *string        `json:"floor,omitempty" db:"floor"`
	FlatNumber          string         `json:"flat_number" db:"flat_number"`
	Status              FlatStatus     `json:"status" db:"status"`
	IsActive            bool           `json:"is_active" db:"is_active"`
	Metadata            map[string]any `json:"metadata" db:"metadata"`
	CreatedBy           *int64         `json:"created_by,omitempty" db:"created_by"`
	CreatedAt           time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at" db:"updated_at"`
	SocietyName         *string        `json:"society_name,omitempty"`
	SocietyCode         *string        `json:"society_code,omitempty"`
	PrimaryResidentName *string        `json:"primary_resident_name,omitempty"`
}

func (f *Flat) ToResponse() *FlatResponse {
	if f == nil {
		return nil
	}
	return &FlatResponse{
		FlatType: f.FlatType, AreaSqft: f.AreaSqft, ID: f.ID, SocietyID: f.SocietyID, Block: f.Block, Floor: f.Floor, FlatNumber: f.FlatNumber,
		Status: f.Status, IsActive: f.IsActive, CreatedBy: f.CreatedBy, CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt, SocietyName: f.SocietyName, SocietyCode: f.SocietyCode,
		PrimaryResidentName: f.PrimaryResidentName,
	}
}

type FlatResident struct {
	ID          int64              `json:"id" db:"id"`
	SocietyID   int64              `json:"society_id" db:"society_id"`
	FlatID      int64              `json:"flat_id" db:"flat_id"`
	UserID      int64              `json:"user_id" db:"user_id"`
	Role        FlatResidentRole   `json:"role" db:"role"`
	Status      FlatResidentStatus `json:"status" db:"status"`
	IsPrimary   bool               `json:"is_primary" db:"is_primary"`
	MovedInAt   time.Time          `json:"moved_in_at" db:"moved_in_at"`
	MovedOutAt  *time.Time         `json:"moved_out_at,omitempty" db:"moved_out_at"`
	Metadata    map[string]any     `json:"metadata" db:"metadata"`
	CreatedBy   *int64             `json:"created_by,omitempty" db:"created_by"`
	CreatedAt   time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" db:"updated_at"`
	UserName    *string            `json:"user_name,omitempty"`
	UserEmail   *string            `json:"user_email,omitempty"`
	UserPhone   *string            `json:"user_phone,omitempty"`
	FlatNumber  *string            `json:"flat_number,omitempty"`
	Block       *string            `json:"block,omitempty"`
	Floor       *string            `json:"floor,omitempty"`
	SocietyName *string            `json:"society_name,omitempty"`
	SocietyCode *string            `json:"society_code,omitempty"`
}

func (r *FlatResident) ToResponse() *FlatResidentResponse {
	if r == nil {
		return nil
	}
	return &FlatResidentResponse{
		ID: r.ID, SocietyID: r.SocietyID, FlatID: r.FlatID, UserID: r.UserID, Role: r.Role,
		Status: r.Status, IsPrimary: r.IsPrimary, MovedInAt: r.MovedInAt, MovedOutAt: r.MovedOutAt,
		CreatedBy: r.CreatedBy, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt, UserName: r.UserName,
		UserEmail: r.UserEmail, UserPhone: r.UserPhone, FlatNumber: r.FlatNumber, Block: r.Block,
		Floor: r.Floor, SocietyName: r.SocietyName, SocietyCode: r.SocietyCode,
	}
}

type FlatClaim struct {
	ID               int64            `json:"id" db:"id"`
	SocietyID        int64            `json:"society_id" db:"society_id"`
	FlatID           int64            `json:"flat_id" db:"flat_id"`
	UserID           int64            `json:"user_id" db:"user_id"`
	RequestedRole    FlatResidentRole `json:"requested_role" db:"requested_role"`
	RequestedPrimary bool             `json:"requested_primary" db:"requested_primary"`
	Status           FlatClaimStatus  `json:"status" db:"status"`
	Note             *string          `json:"note,omitempty" db:"note"`
	RejectionReason  *string          `json:"rejection_reason,omitempty" db:"rejection_reason"`
	ReviewedBy       *int64           `json:"reviewed_by,omitempty" db:"reviewed_by"`
	ReviewedAt       *time.Time       `json:"reviewed_at,omitempty" db:"reviewed_at"`
	CancelledAt      *time.Time       `json:"cancelled_at,omitempty" db:"cancelled_at"`
	Metadata         map[string]any   `json:"metadata" db:"metadata"`
	CreatedAt        time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at" db:"updated_at"`
	UserName         *string          `json:"user_name,omitempty"`
	UserEmail        *string          `json:"user_email,omitempty"`
	UserPhone        *string          `json:"user_phone,omitempty"`
	FlatNumber       *string          `json:"flat_number,omitempty"`
	Block            *string          `json:"block,omitempty"`
	Floor            *string          `json:"floor,omitempty"`
	FlatStatus       *FlatStatus      `json:"flat_status,omitempty"`
	SocietyName      *string          `json:"society_name,omitempty"`
	SocietyCode      *string          `json:"society_code,omitempty"`
	ReviewerName     *string          `json:"reviewer_name,omitempty"`
	ReviewerEmail    *string          `json:"reviewer_email,omitempty"`
	ReviewerPhone    *string          `json:"reviewer_phone,omitempty"`
}

func (c *FlatClaim) ToResponse() *FlatClaimResponse {
	if c == nil {
		return nil
	}
	return &FlatClaimResponse{
		ID: c.ID, SocietyID: c.SocietyID, FlatID: c.FlatID, UserID: c.UserID,
		RequestedRole: c.RequestedRole, RequestedPrimary: c.RequestedPrimary, Status: c.Status,
		Note: c.Note, RejectionReason: c.RejectionReason, ReviewedBy: c.ReviewedBy,
		ReviewedAt: c.ReviewedAt, CancelledAt: c.CancelledAt, CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt, UserName: c.UserName, UserEmail: c.UserEmail, UserPhone: c.UserPhone,
		FlatNumber: c.FlatNumber, Block: c.Block, Floor: c.Floor, FlatStatus: c.FlatStatus,
		SocietyName: c.SocietyName, SocietyCode: c.SocietyCode, ReviewerName: c.ReviewerName,
		ReviewerEmail: c.ReviewerEmail, ReviewerPhone: c.ReviewerPhone,
	}
}
