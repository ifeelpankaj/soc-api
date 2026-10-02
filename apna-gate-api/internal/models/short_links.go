package models

import (
	"strings"
	"time"
)

type ShortLinkResourceType string

const (
	ShortLinkResourceVisitorInvite ShortLinkResourceType = "visitor_invite"
	ShortLinkResourceMemberInvite  ShortLinkResourceType = "member_invite"
)

func (t ShortLinkResourceType) IsValid() bool {
	switch t {
	case ShortLinkResourceVisitorInvite, ShortLinkResourceMemberInvite:
		return true
	default:
		return false
	}
}

type ShortLinkState string

const (
	ShortLinkStateForm            ShortLinkState = "form"
	ShortLinkStateWaitingApproval ShortLinkState = "waiting_approval"
	ShortLinkStateApproved        ShortLinkState = "approved"
	ShortLinkStateCheckedIn       ShortLinkState = "checked_in"
	ShortLinkStateCompleted       ShortLinkState = "completed"
	ShortLinkStateJoin            ShortLinkState = "join"
	ShortLinkStateAccepted        ShortLinkState = "accepted"
	ShortLinkStateExpired         ShortLinkState = "expired"
	ShortLinkStateCancelled       ShortLinkState = "cancelled"
	ShortLinkStateUnavailable     ShortLinkState = "unavailable"
)

type ShortLink struct {
	ID           int64                 `json:"id"`
	ShortCode    string                `json:"short_code"`
	ResourceType ShortLinkResourceType `json:"resource_type"`
	ResourceID   int64                 `json:"resource_id"`
	ExpiresAt    *time.Time            `json:"expires_at,omitempty"`
	RevokedAt    *time.Time            `json:"revoked_at,omitempty"`
	CreatedBy    *int64                `json:"created_by,omitempty"`
	Metadata     map[string]any        `json:"metadata,omitempty"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
}

type ShortLinkResponse struct {
	ShortCode string     `json:"code"`
	URL       string     `json:"url"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type ShortLinkResolveResponse struct {
	Type    ShortLinkResourceType `json:"type"`
	State   ShortLinkState        `json:"state"`
	Invite  any                   `json:"invite,omitempty"`
	Visitor *VisitorSummary       `json:"visitor,omitempty"`
	Entry   *VisitorEntry         `json:"entry,omitempty"`
	QR      *QRTokenResponse      `json:"qr,omitempty"`
	Link    *ShortLinkResponse    `json:"link,omitempty"`
}

func NormalizeShortCode(value string) string {
	return strings.TrimSpace(value)
}
