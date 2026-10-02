package models

import (
	"errors"
	"strings"
	"time"
)

type FlatMemberInviteRole string

const (
	FlatMemberInviteRoleFamily FlatMemberInviteRole = "family"
	FlatMemberInviteRoleTenant FlatMemberInviteRole = "tenant"
)

func (r FlatMemberInviteRole) IsValid() bool {
	switch r {
	case FlatMemberInviteRoleFamily, FlatMemberInviteRoleTenant:
		return true
	default:
		return false
	}
}

func (r FlatMemberInviteRole) ToResidentRole() FlatResidentRole {
	return FlatResidentRole(r)
}

type FlatMemberInviteStatus string

const (
	FlatMemberInviteStatusPending   FlatMemberInviteStatus = "pending"
	FlatMemberInviteStatusAccepted  FlatMemberInviteStatus = "accepted"
	FlatMemberInviteStatusExpired   FlatMemberInviteStatus = "expired"
	FlatMemberInviteStatusCancelled FlatMemberInviteStatus = "cancelled"
)

func (s FlatMemberInviteStatus) IsValid() bool {
	switch s {
	case FlatMemberInviteStatusPending, FlatMemberInviteStatusAccepted, FlatMemberInviteStatusExpired, FlatMemberInviteStatusCancelled:
		return true
	default:
		return false
	}
}

type FlatMemberInvite struct {
	ID        int64                  `json:"id"`
	SocietyID int64                  `json:"society_id"`
	FlatID    int64                  `json:"flat_id"`
	InvitedBy int64                  `json:"invited_by"`
	Role      FlatMemberInviteRole   `json:"role"`
	Phone     *string                `json:"phone,omitempty"`
	Email     *string                `json:"email,omitempty"`
	FullName  string                 `json:"full_name"`
	Status    FlatMemberInviteStatus `json:"status"`
	ExpiresAt time.Time              `json:"expires_at"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

func (i *FlatMemberInvite) ToResponse() *FlatMemberInviteResponse {
	if i == nil {
		return nil
	}
	return &FlatMemberInviteResponse{
		ID:        i.ID,
		SocietyID: i.SocietyID,
		FlatID:    i.FlatID,
		InvitedBy: i.InvitedBy,
		Role:      i.Role,
		Phone:     i.Phone,
		Email:     i.Email,
		FullName:  i.FullName,
		Status:    i.Status,
		ExpiresAt: i.ExpiresAt,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
	}
}

type FlatMemberInviteResponse struct {
	ID        int64                  `json:"id"`
	SocietyID int64                  `json:"society_id"`
	FlatID    int64                  `json:"flat_id"`
	InvitedBy int64                  `json:"invited_by"`
	Role      FlatMemberInviteRole   `json:"role"`
	Phone     *string                `json:"phone,omitempty"`
	Email     *string                `json:"email,omitempty"`
	FullName  string                 `json:"full_name"`
	Status    FlatMemberInviteStatus `json:"status"`
	ExpiresAt time.Time              `json:"expires_at"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	ShortLink *ShortLinkResponse     `json:"short_link,omitempty"`
}

type PublicFlatMemberInviteView struct {
	ID          int64                  `json:"id"`
	Role        FlatMemberInviteRole   `json:"role"`
	FullName    string                 `json:"full_name"`
	Phone       *string                `json:"phone,omitempty"`
	Email       *string                `json:"email,omitempty"`
	Status      FlatMemberInviteStatus `json:"status"`
	ExpiresAt   time.Time              `json:"expires_at"`
	SocietyName string                 `json:"society_name"`
	FlatNumber  string                 `json:"flat_number"`
	Block       *string                `json:"block,omitempty"`
	Floor       *string                `json:"floor,omitempty"`
}

type FlatMemberInviteHistoryFilter struct {
	Status *FlatMemberInviteStatus
	Search *string
	Limit  int32
	Offset int32
}

type FlatMemberInviteHistoryResult struct {
	Invites []*FlatMemberInviteResponse `json:"invites"`
	Total   int64                       `json:"total"`
	Limit   int32                       `json:"limit"`
	Offset  int32                       `json:"offset"`
}

type FlatMemberInviteTokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type CreateFlatMemberInviteRequest struct {
	Role     FlatMemberInviteRole `json:"role" validate:"required"`
	Phone    *string              `json:"phone,omitempty" validate:"omitempty,max=20"`
	Email    *string              `json:"email,omitempty" validate:"omitempty,email,max=255"`
	FullName string               `json:"full_name" validate:"required,max=200"`
}

func (r *CreateFlatMemberInviteRequest) Sanitize() {
	r.Phone = trimPtr(r.Phone)
	r.Email = trimPtr(r.Email)
	r.FullName = strings.TrimSpace(r.FullName)
}

func (r *CreateFlatMemberInviteRequest) Validate() error {
	if r == nil {
		return errors.New("member invite request is required")
	}
	if !r.Role.IsValid() {
		return errors.New("invalid member invite role")
	}
	if strings.TrimSpace(r.FullName) == "" {
		return errors.New("full_name is required")
	}
	phone := ""
	if r.Phone != nil {
		phone = strings.TrimSpace(*r.Phone)
	}
	email := ""
	if r.Email != nil {
		email = strings.TrimSpace(*r.Email)
	}
	if phone == "" && email == "" {
		return errors.New("phone or email is required")
	}
	return nil
}

type AcceptFlatMemberInviteResponse struct {
	Invite   *FlatMemberInviteResponse `json:"invite"`
	Resident *FlatResidentResponse     `json:"resident"`
}

type JoinFlatMemberInviteRequest struct {
	FirstName  string `json:"first_name,omitempty" validate:"omitempty,min=2,max=100,alphanumeric_space"`
	LastName   string `json:"last_name,omitempty" validate:"omitempty,max=100,alphanumeric_space"`
	Email      string `json:"email,omitempty" validate:"omitempty,email,max=255"`
	Identifier string `json:"identifier,omitempty" validate:"omitempty,max=255"`
	Password   string `json:"password" validate:"required,min=8,max=72"`
}

func (r *JoinFlatMemberInviteRequest) Sanitize() {
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	r.Identifier = strings.TrimSpace(r.Identifier)
}

func (r *JoinFlatMemberInviteRequest) IsRegisterFlow() bool {
	return strings.TrimSpace(r.FirstName) != ""
}

type JoinFlatMemberInviteResponse struct {
	User       *UserResponse                   `json:"user"`
	Acceptance *AcceptFlatMemberInviteResponse `json:"acceptance"`
}
