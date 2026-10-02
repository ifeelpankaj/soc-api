package models

import "time"

type FlatResponse struct {
	FlatType            *string    `json:"flat_type,omitempty"`
	AreaSqft            *string    `json:"area_sqft,omitempty"`
	ID                  int64      `json:"id"`
	SocietyID           int64      `json:"society_id"`
	Block               *string    `json:"block,omitempty"`
	Floor               *string    `json:"floor,omitempty"`
	FlatNumber          string     `json:"flat_number"`
	Status              FlatStatus `json:"status"`
	IsActive            bool       `json:"is_active"`
	CreatedBy           *int64     `json:"created_by,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	SocietyName         *string    `json:"society_name,omitempty"`
	SocietyCode         *string    `json:"society_code,omitempty"`
	PrimaryResidentName *string    `json:"primary_resident_name,omitempty"`
}

type BulkCreateFlatsResponse struct {
	Items []*FlatResponse `json:"items"`
	Total int32           `json:"total"`
}

type PaginatedFlatsResponse struct {
	Items  []*FlatResponse `json:"items"`
	Total  int64           `json:"total"`
	Limit  int32           `json:"limit"`
	Offset int32           `json:"offset"`
}

type FlatStatsResponse struct {
	SocietyID     int64 `json:"society_id"`
	TotalFlats    int64 `json:"total_flats"`
	VacantFlats   int64 `json:"vacant_flats"`
	OccupiedFlats int64 `json:"occupied_flats"`
	BlockedFlats  int64 `json:"blocked_flats"`
	ActiveFlats   int64 `json:"active_flats"`
	InactiveFlats int64 `json:"inactive_flats"`
}

type FlatClaimStatsResponse struct {
	TotalClaims     int64 `json:"total_claims"`
	PendingClaims   int64 `json:"pending_claims"`
	ApprovedClaims  int64 `json:"approved_claims"`
	RejectedClaims  int64 `json:"rejected_claims"`
	CancelledClaims int64 `json:"cancelled_claims"`
}

type FlatResidentResponse struct {
	ID          int64              `json:"id"`
	SocietyID   int64              `json:"society_id"`
	FlatID      int64              `json:"flat_id"`
	UserID      int64              `json:"user_id"`
	Role        FlatResidentRole   `json:"role"`
	Status      FlatResidentStatus `json:"status"`
	IsPrimary   bool               `json:"is_primary"`
	MovedInAt   time.Time          `json:"moved_in_at"`
	MovedOutAt  *time.Time         `json:"moved_out_at,omitempty"`
	CreatedBy   *int64             `json:"created_by,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	UserName    *string            `json:"user_name,omitempty"`
	UserEmail   *string            `json:"user_email,omitempty"`
	UserPhone   *string            `json:"user_phone,omitempty"`
	FlatNumber  *string            `json:"flat_number,omitempty"`
	Block       *string            `json:"block,omitempty"`
	Floor       *string            `json:"floor,omitempty"`
	SocietyName *string            `json:"society_name,omitempty"`
	SocietyCode *string            `json:"society_code,omitempty"`
}

type FlatClaimResponse struct {
	ID               int64            `json:"id"`
	SocietyID        int64            `json:"society_id"`
	FlatID           int64            `json:"flat_id"`
	UserID           int64            `json:"user_id"`
	RequestedRole    FlatResidentRole `json:"requested_role"`
	RequestedPrimary bool             `json:"requested_primary"`
	Status           FlatClaimStatus  `json:"status"`
	Note             *string          `json:"note,omitempty"`
	RejectionReason  *string          `json:"rejection_reason,omitempty"`
	ReviewedBy       *int64           `json:"reviewed_by,omitempty"`
	ReviewedAt       *time.Time       `json:"reviewed_at,omitempty"`
	CancelledAt      *time.Time       `json:"cancelled_at,omitempty"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
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

type ApproveFlatClaimResponse struct {
	Claim    *FlatClaimResponse    `json:"claim"`
	Resident *FlatResidentResponse `json:"resident"`
	Flat     *FlatResponse         `json:"flat"`
}
