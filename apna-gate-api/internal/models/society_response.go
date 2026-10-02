package models

import "time"

type SocietyResponse struct {
	ID               int64         `json:"id"`
	Name             string        `json:"name"`
	SocietyCode      string        `json:"society_code"`
	Email            *string       `json:"email,omitempty"`
	PhoneNumber      *string       `json:"phone_number,omitempty"`
	AddressLine1     *string       `json:"address_line1,omitempty"`
	AddressLine2     *string       `json:"address_line2,omitempty"`
	Landmark         *string       `json:"landmark,omitempty"`
	City             *string       `json:"city,omitempty"`
	State            *string       `json:"state,omitempty"`
	Pincode          *string       `json:"pincode,omitempty"`
	Country          string        `json:"country"`
	TotalFlats       int32         `json:"total_flats"`
	TotalBlocks      int32         `json:"total_blocks"`
	Status           SocietyStatus `json:"status"`
	CreatedBy        int64         `json:"created_by"`
	ApprovedBy       *int64        `json:"approved_by,omitempty"`
	ApprovedAt       *time.Time    `json:"approved_at,omitempty"`
	RejectedBy       *int64        `json:"rejected_by,omitempty"`
	RejectedAt       *time.Time    `json:"rejected_at,omitempty"`
	RejectionReason  *string       `json:"rejection_reason,omitempty"`
	SuspendedBy      *int64        `json:"suspended_by,omitempty"`
	SuspendedAt      *time.Time    `json:"suspended_at,omitempty"`
	SuspensionReason *string       `json:"suspension_reason,omitempty"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

type SocietyDetailResponse struct {
	*SocietyResponse
	MembersCount int64 `json:"members_count,omitempty"`
}

type PublicClaimSocietyResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	SocietyCode string  `json:"society_code"`
	City        *string `json:"city,omitempty"`
	State       *string `json:"state,omitempty"`
	Pincode     *string `json:"pincode,omitempty"`
	Country     string  `json:"country"`
	TotalFlats  int32   `json:"total_flats"`
}

type PublicClaimFlatResponse struct {
	ID         int64      `json:"id"`
	Block      *string    `json:"block,omitempty"`
	Floor      *string    `json:"floor,omitempty"`
	FlatNumber string     `json:"flat_number"`
	Status     FlatStatus `json:"status"`
}

type PublicClaimOptionsResponse struct {
	Society *PublicClaimSocietyResponse `json:"society"`
	Flats   []*PublicClaimFlatResponse  `json:"flats"`
}

type PaginatedSocietiesResponse struct {
	Items  []*SocietyResponse `json:"items"`
	Total  int64              `json:"total"`
	Limit  int32              `json:"limit"`
	Offset int32              `json:"offset"`
}

type MySocietyResponse struct {
	Society *SocietyResponse       `json:"society"`
	Member  *SocietyMemberResponse `json:"member"`
}

type SocietyOnboardingBootstrapResponse struct {
	Society               *SocietyResponse `json:"society"`
	HasActiveSubscription bool             `json:"has_active_subscription"`
	HasFlats              bool             `json:"has_flats"`
	HasStaff              bool             `json:"has_staff"`
	IsOnboarded           bool             `json:"is_onboarded"`
	FlatCount             int64            `json:"flat_count"`
	StaffCount            int64            `json:"staff_count"`
	MissingSteps          []string         `json:"missing_steps"`
	NextPath              string           `json:"next_path"`
}

type SocietyDashboardMemberStatsResponse struct {
	TotalActiveMembers int64 `json:"total_active_members"`
	Owners             int64 `json:"owners"`
	Admins             int64 `json:"admins"`
	Staff              int64 `json:"staff"`
	Residents          int64 `json:"residents"`
}

type SocietyDashboardQuotaUsageResponse struct {
	Used      int64 `json:"used"`
	Limit     int64 `json:"limit"`
	Remaining int64 `json:"remaining"`
	Percent   int64 `json:"percent"`
}

type SocietyDashboardSubscriptionUsageResponse struct {
	Flats     SocietyDashboardQuotaUsageResponse `json:"flats"`
	Admins    SocietyDashboardQuotaUsageResponse `json:"admins"`
	Staff     SocietyDashboardQuotaUsageResponse `json:"staff"`
	Residents SocietyDashboardQuotaUsageResponse `json:"residents"`
}

type SocietyDashboardSubscriptionHealthResponse struct {
	IsActive        bool   `json:"is_active"`
	IsExpiringSoon  bool   `json:"is_expiring_soon"`
	DaysUntilExpiry *int64 `json:"days_until_expiry,omitempty"`
	LifecycleLabel  string `json:"lifecycle_label"`
}

type SocietyDashboardBootstrapResponse struct {
	Society             *SocietyResponse                            `json:"society"`
	FlatStats           *FlatStatsResponse                          `json:"flat_stats"`
	ClaimStats          *FlatClaimStatsResponse                     `json:"claim_stats"`
	RecentPendingClaims []*FlatClaimResponse                        `json:"recent_pending_claims"`
	MemberStats         *SocietyDashboardMemberStatsResponse        `json:"member_stats"`
	CurrentSubscription *SocietySubscriptionResponse                `json:"current_subscription,omitempty"`
	SubscriptionUsage   *SocietyDashboardSubscriptionUsageResponse  `json:"subscription_usage,omitempty"`
	SubscriptionHealth  *SocietyDashboardSubscriptionHealthResponse `json:"subscription_health,omitempty"`
	PlanAds             []*PlanResponse                             `json:"plan_ads"`
	VisitorStats        *VisitorEntryStatsResponse                  `json:"visitor_stats,omitempty"`
	VisitorDailyLast7   []VisitorDailyCountResponse                 `json:"visitor_daily_last_7_days,omitempty"`
}

type DeveloperDashboardSocietyStatsResponse struct {
	Total     int64 `json:"total"`
	Pending   int64 `json:"pending"`
	Active    int64 `json:"active"`
	Suspended int64 `json:"suspended"`
	Rejected  int64 `json:"rejected"`
}

type DeveloperDashboardPlanStatsResponse struct {
	Total    int64 `json:"total"`
	Active   int64 `json:"active"`
	Inactive int64 `json:"inactive"`
}

type DeveloperDashboardResidenceStatsResponse struct {
	TotalResidents  int64 `json:"total_residents"`
	ActiveResidents int64 `json:"active_residents"`
}

type DeveloperDashboardBootstrapResponse struct {
	SocietyStats           *DeveloperDashboardSocietyStatsResponse   `json:"society_stats"`
	PlanStats              *DeveloperDashboardPlanStatsResponse      `json:"plan_stats"`
	SubscriptionStats      *SubscriptionStatsResponse                `json:"subscription_stats"`
	ResidenceStats         *DeveloperDashboardResidenceStatsResponse `json:"residence_stats"`
	RecentPendingSocieties []*SocietyResponse                        `json:"recent_pending_societies"`
	RecentSubscriptions    []*SocietySubscriptionResponse            `json:"recent_subscriptions"`
}

type CreateGuardResponse struct {
	User   *UserResponse          `json:"user"`
	Member *SocietyMemberResponse `json:"member"`
}

type SocietyMemberResponse struct {
	ID           int64               `json:"id"`
	SocietyID    int64               `json:"society_id"`
	UserID       int64               `json:"user_id"`
	Role         SocietyMemberRole   `json:"role"`
	Status       SocietyMemberStatus `json:"status"`
	InvitedBy    *int64              `json:"invited_by,omitempty"`
	JoinedAt     time.Time           `json:"joined_at"`
	RemovedBy    *int64              `json:"removed_by,omitempty"`
	RemovedAt    *time.Time          `json:"removed_at,omitempty"`
	RemoveReason *string             `json:"remove_reason,omitempty"`
	UserFullName *string             `json:"user_full_name,omitempty"`
	UserEmail    *string             `json:"user_email,omitempty"`
	UserPhone    *string             `json:"user_phone,omitempty"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

type SocietyMemberDetailResponse struct {
	Member     *SocietyMemberResponse  `json:"member"`
	OwnedFlats []*FlatResidentResponse `json:"owned_flats"`
	Residences []*FlatResidentResponse `json:"residences"`
}

type SocietyMemberSummaryResponse struct {
	TotalMembers     int64 `json:"total_members"`
	ActiveMembers    int64 `json:"active_members"`
	PendingMembers   int64 `json:"pending_members"`
	SuspendedMembers int64 `json:"suspended_members"`
	RemovedMembers   int64 `json:"removed_members"`
	Owners           int64 `json:"owners"`
	Admins           int64 `json:"admins"`
	Staff            int64 `json:"staff"`
	Residents        int64 `json:"residents"`
}

type PaginatedMembersResponse struct {
	Items  []*SocietyMemberResponse `json:"items"`
	Total  int64                    `json:"total"`
	Limit  int32                    `json:"limit"`
	Offset int32                    `json:"offset"`
}
