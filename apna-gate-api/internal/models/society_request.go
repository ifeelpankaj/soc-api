package models

import (
	"errors"
	"strings"
	"time"
)

type CreateSocietyRequest struct {
	Name         string         `json:"name" validate:"required,min=2,max=200"`
	SocietyCode  *string        `json:"society_code,omitempty" validate:"omitempty,max=50"`
	Email        *string        `json:"email,omitempty" validate:"omitempty,email,max=255"`
	PhoneNumber  *string        `json:"phone_number,omitempty" validate:"omitempty,max=20"`
	AddressLine1 *string        `json:"address_line1,omitempty" validate:"omitempty,max=255"`
	AddressLine2 *string        `json:"address_line2,omitempty" validate:"omitempty,max=255"`
	Landmark     *string        `json:"landmark,omitempty" validate:"omitempty,max=255"`
	City         *string        `json:"city,omitempty" validate:"omitempty,max=100"`
	State        *string        `json:"state,omitempty" validate:"omitempty,max=100"`
	Pincode      *string        `json:"pincode,omitempty" validate:"omitempty,max=20"`
	Country      *string        `json:"country,omitempty" validate:"omitempty,max=100"`
	TotalFlats   int32          `json:"total_flats" validate:"gte=0"`
	TotalBlocks  int32          `json:"total_blocks" validate:"gte=0"`
	Metadata     map[string]any `json:"metadata,omitempty"`
}

func (r *CreateSocietyRequest) Sanitize() {
	r.Name = strings.TrimSpace(r.Name)
	r.SocietyCode = trimUpperPtr(r.SocietyCode)
	r.Email = trimLowerPtr(r.Email)
	r.PhoneNumber = trimPtr(r.PhoneNumber)
	r.AddressLine1 = trimPtr(r.AddressLine1)
	r.AddressLine2 = trimPtr(r.AddressLine2)
	r.Landmark = trimPtr(r.Landmark)
	r.City = trimPtr(r.City)
	r.State = trimPtr(r.State)
	r.Pincode = trimPtr(r.Pincode)
	r.Country = trimPtr(r.Country)
}

func (r *CreateSocietyRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("society name is required")
	}
	if r.TotalFlats < 0 || r.TotalBlocks < 0 {
		return errors.New("total_flats and total_blocks cannot be negative")
	}
	return nil
}

type UpdateSocietyRequest struct {
	Name         *string        `json:"name,omitempty" validate:"omitempty,min=2,max=200"`
	Email        *string        `json:"email,omitempty" validate:"omitempty,email,max=255"`
	PhoneNumber  *string        `json:"phone_number,omitempty" validate:"omitempty,max=20"`
	AddressLine1 *string        `json:"address_line1,omitempty" validate:"omitempty,max=255"`
	AddressLine2 *string        `json:"address_line2,omitempty" validate:"omitempty,max=255"`
	Landmark     *string        `json:"landmark,omitempty" validate:"omitempty,max=255"`
	City         *string        `json:"city,omitempty" validate:"omitempty,max=100"`
	State        *string        `json:"state,omitempty" validate:"omitempty,max=100"`
	Pincode      *string        `json:"pincode,omitempty" validate:"omitempty,max=20"`
	Country      *string        `json:"country,omitempty" validate:"omitempty,max=100"`
	TotalFlats   *int32         `json:"total_flats,omitempty" validate:"omitempty,gte=0"`
	TotalBlocks  *int32         `json:"total_blocks,omitempty" validate:"omitempty,gte=0"`
	Metadata     map[string]any `json:"metadata,omitempty"`
}

func (r *UpdateSocietyRequest) Sanitize() {
	r.Name = trimPtr(r.Name)
	r.Email = trimLowerPtr(r.Email)
	r.PhoneNumber = trimPtr(r.PhoneNumber)
	r.AddressLine1 = trimPtr(r.AddressLine1)
	r.AddressLine2 = trimPtr(r.AddressLine2)
	r.Landmark = trimPtr(r.Landmark)
	r.City = trimPtr(r.City)
	r.State = trimPtr(r.State)
	r.Pincode = trimPtr(r.Pincode)
	r.Country = trimPtr(r.Country)
}

func (r *UpdateSocietyRequest) Validate() error {
	if r.Name == nil && r.Email == nil && r.PhoneNumber == nil && r.AddressLine1 == nil &&
		r.AddressLine2 == nil && r.Landmark == nil && r.City == nil && r.State == nil &&
		r.Pincode == nil && r.Country == nil && r.TotalFlats == nil && r.TotalBlocks == nil && r.Metadata == nil {
		return errors.New("at least one field must be provided")
	}
	if r.Name != nil && *r.Name == "" {
		return errors.New("society name cannot be empty")
	}
	if r.TotalFlats != nil && *r.TotalFlats < 0 {
		return errors.New("total_flats cannot be negative")
	}
	if r.TotalBlocks != nil && *r.TotalBlocks < 0 {
		return errors.New("total_blocks cannot be negative")
	}
	return nil
}

type AddSocietyMemberRequest struct {
	SocietyID int64             `json:"society_id" validate:"required"`
	UserID    int64             `json:"user_id" validate:"required"`
	Role      SocietyMemberRole `json:"role" validate:"required"`
	Metadata  map[string]any    `json:"metadata,omitempty"`
}

type CreateGuardRequest struct {
	FirstName   string `json:"first_name" validate:"required,min=2,max=100,alphanumeric_space"`
	LastName    string `json:"last_name,omitempty" validate:"omitempty,max=100,alphanumeric_space"`
	Email       string `json:"email" validate:"required,email,max=255"`
	PhoneNumber string `json:"phone_number" validate:"required,phone_intl,max=20"`
	Password    string `json:"password" validate:"required,min=8,max=72"`
}

func (r *CreateGuardRequest) Sanitize() {
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	r.PhoneNumber = strings.TrimSpace(r.PhoneNumber)
}

type ChangeSocietyMemberRoleRequest struct {
	SocietyID int64             `json:"society_id" validate:"required"`
	UserID    int64             `json:"user_id" validate:"required"`
	Role      SocietyMemberRole `json:"role" validate:"required"`
}

type SuspendSocietyMemberRequest struct {
	SocietyID int64  `json:"society_id" validate:"required"`
	UserID    int64  `json:"user_id" validate:"required"`
	Reason    string `json:"reason,omitempty" validate:"max=500"`
}

type ReactivateSocietyMemberRequest struct {
	SocietyID int64 `json:"society_id" validate:"required"`
	UserID    int64 `json:"user_id" validate:"required"`
}

type RemoveSocietyMemberRequest struct {
	SocietyID int64  `json:"society_id" validate:"required"`
	UserID    int64  `json:"user_id" validate:"required"`
	Reason    string `json:"reason,omitempty" validate:"max=500"`
}

type SocietyReasonRequest struct {
	Reason string `json:"reason" validate:"required,max=500"`
}

func (r *SocietyReasonRequest) Sanitize() {
	r.Reason = strings.TrimSpace(r.Reason)
}

type TransferOwnershipRequest struct {
	NewOwnerUserID int64 `json:"new_owner_user_id" validate:"required,gt=0"`
}

type GetSocietyFilter struct {
	ID             *int64
	Code           *string
	CreatedBy      *int64
	Status         *string
	IncludeDeleted bool
}

type ListSocietiesFilter struct {
	ID          *int64
	Status      *string
	Search      string
	SearchMode  string
	Name        string
	Code        string
	City        string
	State       string
	Country     string
	Pincode     string
	CreatedBy   *int64
	ApprovedBy  *int64
	RejectedBy  *int64
	SuspendedBy *int64
	CreatedFrom *time.Time
	CreatedTo   *time.Time
	Limit       int32
	Offset      int32
	SortBy      string
	SortOrder   string
}

type GetSocietyMemberFilter struct {
	ID        *int64
	SocietyID *int64
	UserID    *int64
	Role      *string
	Status    *string
	Email     *string
	Phone     *string
}

type ListSocietyMembersFilter struct {
	SocietyID  int64
	Search     string
	SearchMode string
	Role       *string
	Status     *string
	UserID     *int64
	InvitedBy  *int64
	RemovedBy  *int64
	JoinedFrom *time.Time
	JoinedTo   *time.Time
	Limit      int32
	Offset     int32
	SortBy     string
	SortOrder  string
}

func trimPtr(v *string) *string {
	if v == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*v)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func trimLowerPtr(v *string) *string {
	if v = trimPtr(v); v != nil {
		lower := strings.ToLower(*v)
		return &lower
	}
	return nil
}

func trimUpperPtr(v *string) *string {
	if v = trimPtr(v); v != nil {
		upper := strings.ToUpper(*v)
		return &upper
	}
	return nil
}
