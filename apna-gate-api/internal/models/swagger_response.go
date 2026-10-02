package models

type ErrorResponseDoc struct {
	Success bool      `json:"success" example:"false"`
	Error   ErrorData `json:"error"`
}

type MessageData struct {
	Message string `json:"message" example:"Request completed successfully"`
}

type DevOTPMessageData struct {
	Message string `json:"message" example:"Please verify your email using the OTP sent to your email address"`
	DevOTP  string `json:"dev_otp,omitempty" example:"123456"`
}

type UserData struct {
	User *UserResponse `json:"user"`
}

type AuthSessionData struct {
	User                  *UserResponse `json:"user"`
	AccessToken           string        `json:"access_token,omitempty"`
	RefreshToken          string        `json:"refresh_token,omitempty"`
	AccessTokenExpiresAt  string        `json:"access_token_expires_at,omitempty"`
	RefreshTokenExpiresAt string        `json:"refresh_token_expires_at,omitempty"`
}

type AuthRefreshData struct {
	Message              string `json:"message" example:"Access token refreshed successfully"`
	AccessToken          string `json:"access_token,omitempty"`
	AccessTokenExpiresAt string `json:"access_token_expires_at,omitempty"`
}

type BootstrapData struct {
	User             *UserResponse             `json:"user"`
	Memberships      []*SocietyMemberResponse  `json:"memberships"`
	Residences       []*FlatResidentResponse   `json:"residences"`
	DefaultDashboard *DefaultDashboardResponse `json:"defaultDashboard"`
}

type UserMessageData struct {
	User    *UserResponse `json:"user"`
	Message string        `json:"message" example:"Please verify your email using the OTP sent to your email address"`
	DevOTP  string        `json:"dev_otp,omitempty" example:"123456"`
}

type RegisterAPIResponse struct {
	Success bool            `json:"success" example:"true"`
	Message string          `json:"message" example:"Account created successfully"`
	Data    UserMessageData `json:"data"`
}

type VerifyOTPAPIResponse struct {
	Success bool     `json:"success" example:"true"`
	Message string   `json:"message" example:"Email verified successfully"`
	Data    UserData `json:"data"`
}

type ResendOTPAPIResponse struct {
	Success bool              `json:"success" example:"true"`
	Message string            `json:"message" example:"verification OTP sent successfully"`
	Data    DevOTPMessageData `json:"data"`
}

type LoginAPIResponse struct {
	Success bool            `json:"success" example:"true"`
	Message string          `json:"message" example:"Login successful"`
	Data    AuthSessionData `json:"data"`
}

type PublicClaimOptionsData struct {
	Society *PublicClaimSocietyResponse `json:"society"`
	Flats   []*PublicClaimFlatResponse  `json:"flats"`
}

type PublicClaimOptionsAPIResponse struct {
	Success bool                   `json:"success" example:"true"`
	Message string                 `json:"message" example:"Claim options fetched successfully"`
	Data    PublicClaimOptionsData `json:"data"`
}

type RefreshTokenAPIResponse struct {
	Success bool            `json:"success" example:"true"`
	Message string          `json:"message" example:"Access token refreshed successfully"`
	Data    AuthRefreshData `json:"data"`
}

type LogoutAPIResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Logout successful"`
	Data    MessageData `json:"data"`
}

type GetProfileAPIResponse struct {
	Success bool     `json:"success" example:"true"`
	Message string   `json:"message" example:"Profile fetched successfully"`
	Data    UserData `json:"data"`
}

type BootstrapAPIResponse struct {
	Success bool          `json:"success" example:"true"`
	Message string        `json:"message" example:"Bootstrap fetched successfully"`
	Data    BootstrapData `json:"data"`
}

type ForgotPasswordAPIResponse struct {
	Success bool              `json:"success" example:"true"`
	Message string            `json:"message" example:"if this email exists, password reset instructions have been sent"`
	Data    DevOTPMessageData `json:"data"`
}

type ResetPasswordAPIResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Password reset successfully"`
	Data    MessageData `json:"data"`
}

type ChangePasswordAPIResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Password changed successfully"`
	Data    MessageData `json:"data"`
}

type MessageAPIResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Request completed successfully"`
	Data    MessageData `json:"data"`
}

type DeviceTokenData struct {
	DeviceToken *DeviceToken `json:"device_token"`
}

type RegisterDeviceTokenAPIResponse struct {
	Success bool            `json:"success" example:"true"`
	Message string          `json:"message" example:"Device token registered successfully"`
	Data    DeviceTokenData `json:"data"`
}

type NotificationsAPIResponse struct {
	Success bool                    `json:"success" example:"true"`
	Message string                  `json:"message" example:"Notifications fetched successfully"`
	Data    *NotificationListResult `json:"data"`
}

type NotificationUnreadCountAPIResponse struct {
	Success bool                             `json:"success" example:"true"`
	Message string                           `json:"message" example:"Unread count fetched successfully"`
	Data    *NotificationUnreadCountResponse `json:"data"`
}

type NotificationReadAPIResponse struct {
	Success bool                      `json:"success" example:"true"`
	Message string                    `json:"message" example:"Notification marked read"`
	Data    *NotificationReadResponse `json:"data"`
}

type SocietyData struct {
	Society *SocietyResponse `json:"society"`
}

type SocietyDetailData struct {
	Society *SocietyDetailResponse `json:"society"`
}

type SocietyMemberData struct {
	Member *SocietyMemberResponse `json:"member"`
}

type SocietyMemberDetailData struct {
	MemberDetail *SocietyMemberDetailResponse `json:"member_detail"`
}

type SocietyMemberSummaryData struct {
	Summary *SocietyMemberSummaryResponse `json:"summary"`
}

type PaginatedSocietiesData struct {
	Societies *PaginatedSocietiesResponse `json:"societies"`
}

type PaginatedMembersData struct {
	Members *PaginatedMembersResponse `json:"members"`
}

type MySocietiesData struct {
	Societies []*MySocietyResponse `json:"societies"`
}

type SocietyOnboardingBootstrapData struct {
	Onboarding *SocietyOnboardingBootstrapResponse `json:"onboarding"`
}

type SocietyDashboardBootstrapData struct {
	Dashboard *SocietyDashboardBootstrapResponse `json:"dashboard"`
}

type DeveloperDashboardBootstrapData struct {
	Dashboard *DeveloperDashboardBootstrapResponse `json:"dashboard"`
}

type GuardData struct {
	Guard *CreateGuardResponse `json:"guard"`
}

type SocietyAPIResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Society request created successfully"`
	Data    SocietyData `json:"data"`
}

type SocietyDetailAPIResponse struct {
	Success bool              `json:"success" example:"true"`
	Message string            `json:"message" example:"Society fetched successfully"`
	Data    SocietyDetailData `json:"data"`
}

type PaginatedSocietiesAPIResponse struct {
	Success bool                   `json:"success" example:"true"`
	Message string                 `json:"message" example:"Societies fetched successfully"`
	Data    PaginatedSocietiesData `json:"data"`
}

type MySocietiesAPIResponse struct {
	Success bool            `json:"success" example:"true"`
	Message string          `json:"message" example:"My societies fetched successfully"`
	Data    MySocietiesData `json:"data"`
}

type SocietyMemberAPIResponse struct {
	Success bool              `json:"success" example:"true"`
	Message string            `json:"message" example:"Member fetched successfully"`
	Data    SocietyMemberData `json:"data"`
}

type SocietyMemberDetailAPIResponse struct {
	Success bool                    `json:"success" example:"true"`
	Message string                  `json:"message" example:"Member fetched successfully"`
	Data    SocietyMemberDetailData `json:"data"`
}

type SocietyMemberSummaryAPIResponse struct {
	Success bool                     `json:"success" example:"true"`
	Message string                   `json:"message" example:"Member summary fetched successfully"`
	Data    SocietyMemberSummaryData `json:"data"`
}

type SocietyOnboardingBootstrapAPIResponse struct {
	Success bool                           `json:"success" example:"true"`
	Message string                         `json:"message" example:"Society onboarding bootstrap fetched successfully"`
	Data    SocietyOnboardingBootstrapData `json:"data"`
}

type SocietyDashboardBootstrapAPIResponse struct {
	Success bool                          `json:"success" example:"true"`
	Message string                        `json:"message" example:"Society dashboard bootstrap fetched successfully"`
	Data    SocietyDashboardBootstrapData `json:"data"`
}

type DeveloperDashboardBootstrapAPIResponse struct {
	Success bool                            `json:"success" example:"true"`
	Message string                          `json:"message" example:"Developer dashboard bootstrap fetched successfully"`
	Data    DeveloperDashboardBootstrapData `json:"data"`
}

type GuardAPIResponse struct {
	Success bool      `json:"success" example:"true"`
	Message string    `json:"message" example:"Guard created successfully"`
	Data    GuardData `json:"data"`
}

type PaginatedMembersAPIResponse struct {
	Success bool                 `json:"success" example:"true"`
	Message string               `json:"message" example:"Members fetched successfully"`
	Data    PaginatedMembersData `json:"data"`
}

type FlatData struct {
	Flat *FlatResponse `json:"flat"`
}

type FlatsData struct {
	Flats []*FlatResponse `json:"flats"`
}

type PaginatedFlatsData struct {
	Flats *PaginatedFlatsResponse `json:"flats"`
}

type BulkFlatsData struct {
	Flats *BulkCreateFlatsResponse `json:"flats"`
}

type FlatStatsData struct {
	Stats *FlatStatsResponse `json:"stats"`
}

type FlatClaimData struct {
	Claim *FlatClaimResponse `json:"claim"`
}

type FlatClaimsData struct {
	Claims []*FlatClaimResponse `json:"claims"`
}

type FlatApprovalData struct {
	Approval *ApproveFlatClaimResponse `json:"approval"`
}

type FlatResidentData struct {
	Resident *FlatResidentResponse `json:"resident"`
}

type FlatResidentsData struct {
	Residents []*FlatResidentResponse `json:"residents"`
}

type FlatMemberInvitesData struct {
	Invites []*FlatMemberInviteResponse `json:"invites"`
	Total   int64                       `json:"total,omitempty"`
	Limit   int32                       `json:"limit,omitempty"`
	Offset  int32                       `json:"offset,omitempty"`
}

type FlatMemberInviteData struct {
	Invite *FlatMemberInviteResponse `json:"invite"`
}

type PublicFlatMemberInviteData struct {
	Invite *PublicFlatMemberInviteView `json:"invite"`
}

type FlatMemberInviteTokenData struct {
	Invite    *FlatMemberInviteResponse      `json:"invite"`
	Token     *FlatMemberInviteTokenResponse `json:"token"`
	Link      *ShortLinkResponse             `json:"link,omitempty"`
	ShortLink *ShortLinkResponse             `json:"short_link,omitempty"`
}

type AcceptFlatMemberInviteData struct {
	Acceptance *AcceptFlatMemberInviteResponse `json:"acceptance"`
}

type JoinFlatMemberInviteData struct {
	Join *JoinFlatMemberInviteResponse `json:"join"`
}

type JoinFlatMemberInviteAPIResponse struct {
	Success bool                     `json:"success" example:"true"`
	Message string                   `json:"message" example:"Member joined successfully"`
	Data    JoinFlatMemberInviteData `json:"data"`
}

type MyResidencesData struct {
	Residences []*FlatResidentResponse `json:"residences"`
}

type SocietyVisitorSettingsData struct {
	VisitorSettings *SocietyVisitorSettingsResponse `json:"visitor_settings"`
}

type FlatVisitorSettingsData struct {
	VisitorSettings []FlatVisitorSettingsResponse `json:"visitor_settings"`
}

type FlatVisitorSettingData struct {
	VisitorSetting *FlatVisitorSettingsResponse `json:"visitor_setting"`
}

type VisitorEntryOptionsData struct {
	Options *VisitorEntryOptionsResponse `json:"options"`
}

type VisitorInviteData struct {
	Invite *PublicVisitorInviteView `json:"invite"`
	View   string                   `json:"view,omitempty" example:"form"`
	Entry  *VisitorEntry            `json:"entry,omitempty"`
	QR     *QRTokenResponse         `json:"qr,omitempty"`
}

type VisitorInviteTokenData struct {
	Invite    *VisitorInvite     `json:"invite"`
	Token     *QRTokenResponse   `json:"token"`
	Link      *ShortLinkResponse `json:"link,omitempty"`
	ShortLink *ShortLinkResponse `json:"short_link,omitempty"`
}

type VisitorInviteHistoryData struct {
	Invites []*VisitorInviteHistoryItem `json:"invites"`
	Total   int64                       `json:"total"`
	Limit   int32                       `json:"limit"`
	Offset  int32                       `json:"offset"`
}

type VisitorInviteHistoryItemData struct {
	Invite *VisitorInviteHistoryItem `json:"invite"`
}

type ShortLinkResolveData struct {
	Link *ShortLinkResolveResponse `json:"link"`
}

type VisitorEntryMutationAPIResponse struct {
	Success bool                          `json:"success" example:"true"`
	Message string                        `json:"message" example:"Visitor entry created successfully"`
	Data    *VisitorEntryMutationResponse `json:"data"`
}

type VisitorEntryData struct {
	Entry *VisitorEntry `json:"entry"`
}

type VisitorEntriesData struct {
	Entries []*VisitorEntry `json:"entries"`
	Total   int64           `json:"total,omitempty"`
	Limit   int32           `json:"limit,omitempty"`
	Offset  int32           `json:"offset,omitempty"`
}

type VisitorEntryStatsData struct {
	Stats *VisitorEntryStatsResponse `json:"stats"`
}

type VisitorEntryStatsAPIResponse struct {
	Success bool                  `json:"success" example:"true"`
	Message string                `json:"message" example:"Visitor entry stats fetched successfully"`
	Data    VisitorEntryStatsData `json:"data"`
}

type VisitorEntryDailyStatsData struct {
	Stats *VisitorEntryDailyStatsResponse `json:"stats"`
}

type VisitorEntryDailyStatsAPIResponse struct {
	Success bool                       `json:"success" example:"true"`
	Message string                     `json:"message" example:"Visitor entry daily stats fetched successfully"`
	Data    VisitorEntryDailyStatsData `json:"data"`
}

type VisitorPendingEntriesData struct {
	Entries []*VisitorPendingEntry `json:"entries"`
	Total   int64                  `json:"total"`
	Limit   int32                  `json:"limit"`
	Offset  int32                  `json:"offset"`
}

type VisitorPendingEntriesAPIResponse struct {
	Success bool                      `json:"success" example:"true"`
	Message string                    `json:"message" example:"Pending visitor approvals fetched successfully"`
	Data    VisitorPendingEntriesData `json:"data"`
}

type GuardDeskBootstrapData struct {
	Desk *GuardDeskBootstrapResponse `json:"desk"`
}

type GuardDeskBootstrapAPIResponse struct {
	Success bool                   `json:"success" example:"true"`
	Message string                 `json:"message" example:"Guard desk bootstrap fetched successfully"`
	Data    GuardDeskBootstrapData `json:"data"`
}

type FlatVisitorContextData struct {
	Context *FlatVisitorContextResponse `json:"context"`
}

type FlatVisitorContextAPIResponse struct {
	Success bool                   `json:"success" example:"true"`
	Message string                 `json:"message" example:"Flat visitor context fetched successfully"`
	Data    FlatVisitorContextData `json:"data"`
}

type MemberVisitorApprovalStatsData struct {
	Stats *MemberVisitorApprovalStatsResponse `json:"stats"`
}

type MemberVisitorApprovalStatsAPIResponse struct {
	Success bool                           `json:"success" example:"true"`
	Message string                         `json:"message" example:"Member visitor approval stats fetched successfully"`
	Data    MemberVisitorApprovalStatsData `json:"data"`
}

type SocietyFlatVisitorSettingsData struct {
	Settings []*SocietyFlatVisitorSettingRow `json:"settings"`
	Total    int64                           `json:"total"`
	Limit    int32                           `json:"limit"`
	Offset   int32                           `json:"offset"`
}

type SocietyFlatVisitorSettingsAPIResponse struct {
	Success bool                           `json:"success" example:"true"`
	Message string                         `json:"message" example:"Society flat visitor settings fetched successfully"`
	Data    SocietyFlatVisitorSettingsData `json:"data"`
}

type VisitorEntryEventsData struct {
	Events []*VisitorEntryEvent `json:"events"`
}

type FlatAPIResponse struct {
	Success bool     `json:"success" example:"true"`
	Message string   `json:"message" example:"Flat fetched successfully"`
	Data    FlatData `json:"data"`
}

type FlatsAPIResponse struct {
	Success bool      `json:"success" example:"true"`
	Message string    `json:"message" example:"Flats fetched successfully"`
	Data    FlatsData `json:"data"`
}

type PaginatedFlatsAPIResponse struct {
	Success bool               `json:"success" example:"true"`
	Message string             `json:"message" example:"Flats fetched successfully"`
	Data    PaginatedFlatsData `json:"data"`
}

type BulkFlatsAPIResponse struct {
	Success bool          `json:"success" example:"true"`
	Message string        `json:"message" example:"Flats created successfully"`
	Data    BulkFlatsData `json:"data"`
}

type FlatStatsAPIResponse struct {
	Success bool          `json:"success" example:"true"`
	Message string        `json:"message" example:"Flat stats fetched successfully"`
	Data    FlatStatsData `json:"data"`
}

type FlatClaimStatsData struct {
	Stats *FlatClaimStatsResponse `json:"stats"`
}

type FlatClaimStatsAPIResponse struct {
	Success bool               `json:"success" example:"true"`
	Message string             `json:"message" example:"Flat claim stats fetched successfully"`
	Data    FlatClaimStatsData `json:"data"`
}

type FlatClaimAPIResponse struct {
	Success bool          `json:"success" example:"true"`
	Message string        `json:"message" example:"Flat claim fetched successfully"`
	Data    FlatClaimData `json:"data"`
}

type FlatClaimsAPIResponse struct {
	Success bool           `json:"success" example:"true"`
	Message string         `json:"message" example:"Flat claims fetched successfully"`
	Data    FlatClaimsData `json:"data"`
}

type FlatApprovalAPIResponse struct {
	Success bool             `json:"success" example:"true"`
	Message string           `json:"message" example:"Flat claim approved successfully"`
	Data    FlatApprovalData `json:"data"`
}

type FlatResidentAPIResponse struct {
	Success bool             `json:"success" example:"true"`
	Message string           `json:"message" example:"Resident fetched successfully"`
	Data    FlatResidentData `json:"data"`
}

type FlatResidentsAPIResponse struct {
	Success bool              `json:"success" example:"true"`
	Message string            `json:"message" example:"Residents fetched successfully"`
	Data    FlatResidentsData `json:"data"`
}

type FlatMemberInvitesAPIResponse struct {
	Success bool                  `json:"success" example:"true"`
	Message string                `json:"message" example:"Member invites fetched successfully"`
	Data    FlatMemberInvitesData `json:"data"`
}

type FlatMemberInviteAPIResponse struct {
	Success bool                 `json:"success" example:"true"`
	Message string               `json:"message" example:"Member invite fetched successfully"`
	Data    FlatMemberInviteData `json:"data"`
}

type FlatMemberInviteTokenAPIResponse struct {
	Success bool                      `json:"success" example:"true"`
	Message string                    `json:"message" example:"Member invite created successfully"`
	Data    FlatMemberInviteTokenData `json:"data"`
}

type VisitorInviteHistoryAPIResponse struct {
	Success bool                     `json:"success" example:"true"`
	Message string                   `json:"message" example:"Visitor invites fetched successfully"`
	Data    VisitorInviteHistoryData `json:"data"`
}

type VisitorInviteHistoryItemAPIResponse struct {
	Success bool                         `json:"success" example:"true"`
	Message string                       `json:"message" example:"Visitor invite fetched successfully"`
	Data    VisitorInviteHistoryItemData `json:"data"`
}

type ShortLinkResolveAPIResponse struct {
	Success bool                 `json:"success" example:"true"`
	Message string               `json:"message" example:"Short link resolved successfully"`
	Data    ShortLinkResolveData `json:"data"`
}

type PublicFlatMemberInviteAPIResponse struct {
	Success bool                       `json:"success" example:"true"`
	Message string                     `json:"message" example:"Member invite fetched successfully"`
	Data    PublicFlatMemberInviteData `json:"data"`
}

type AcceptFlatMemberInviteAPIResponse struct {
	Success bool                       `json:"success" example:"true"`
	Message string                     `json:"message" example:"Member invite accepted successfully"`
	Data    AcceptFlatMemberInviteData `json:"data"`
}

type MyResidencesAPIResponse struct {
	Success bool             `json:"success" example:"true"`
	Message string           `json:"message" example:"My residences fetched successfully"`
	Data    MyResidencesData `json:"data"`
}

type SocietyVisitorSettingsAPIResponse struct {
	Success bool                       `json:"success" example:"true"`
	Message string                     `json:"message" example:"Visitor settings fetched successfully"`
	Data    SocietyVisitorSettingsData `json:"data"`
}

type FlatVisitorSettingsAPIResponse struct {
	Success bool                    `json:"success" example:"true"`
	Message string                  `json:"message" example:"Flat visitor settings fetched successfully"`
	Data    FlatVisitorSettingsData `json:"data"`
}

type FlatVisitorSettingAPIResponse struct {
	Success bool                   `json:"success" example:"true"`
	Message string                 `json:"message" example:"Flat visitor setting updated successfully"`
	Data    FlatVisitorSettingData `json:"data"`
}

type VisitorEntryOptionsAPIResponse struct {
	Success bool                    `json:"success" example:"true"`
	Message string                  `json:"message" example:"Visitor entry options fetched successfully"`
	Data    VisitorEntryOptionsData `json:"data"`
}

type VisitorInviteAPIResponse struct {
	Success bool              `json:"success" example:"true"`
	Message string            `json:"message" example:"Visitor invite fetched successfully"`
	Data    VisitorInviteData `json:"data"`
}

type VisitorInviteTokenAPIResponse struct {
	Success bool                   `json:"success" example:"true"`
	Message string                 `json:"message" example:"Visitor invite created successfully"`
	Data    VisitorInviteTokenData `json:"data"`
}

type VisitorEntryAPIResponse struct {
	Success bool             `json:"success" example:"true"`
	Message string           `json:"message" example:"Visitor entry fetched successfully"`
	Data    VisitorEntryData `json:"data"`
}

type VisitorEntriesAPIResponse struct {
	Success bool               `json:"success" example:"true"`
	Message string             `json:"message" example:"Visitor entries fetched successfully"`
	Data    VisitorEntriesData `json:"data"`
}

type VisitorEntryEventsAPIResponse struct {
	Success bool                   `json:"success" example:"true"`
	Message string                 `json:"message" example:"Visitor entry events fetched successfully"`
	Data    VisitorEntryEventsData `json:"data"`
}

type PlanData struct {
	Plan *PlanResponse `json:"plan"`
}

type PlansData struct {
	Plans []*PlanResponse `json:"plans"`
}

type SubscriptionData struct {
	Subscription *SocietySubscriptionResponse `json:"subscription"`
}

type SubscriptionsData struct {
	Subscriptions []*SocietySubscriptionResponse `json:"subscriptions"`
}

type SubscriptionStatsData struct {
	Stats *SubscriptionStatsResponse `json:"stats"`
}

type PlanAPIResponse struct {
	Success bool     `json:"success" example:"true"`
	Message string   `json:"message" example:"Plan fetched successfully"`
	Data    PlanData `json:"data"`
}

type PlansAPIResponse struct {
	Success bool      `json:"success" example:"true"`
	Message string    `json:"message" example:"Plans fetched successfully"`
	Data    PlansData `json:"data"`
}

type SubscriptionAPIResponse struct {
	Success bool             `json:"success" example:"true"`
	Message string           `json:"message" example:"Subscription fetched successfully"`
	Data    SubscriptionData `json:"data"`
}

type SubscriptionsAPIResponse struct {
	Success bool              `json:"success" example:"true"`
	Message string            `json:"message" example:"Subscriptions fetched successfully"`
	Data    SubscriptionsData `json:"data"`
}

type SubscriptionStatsAPIResponse struct {
	Success bool                  `json:"success" example:"true"`
	Message string                `json:"message" example:"Subscription stats fetched successfully"`
	Data    SubscriptionStatsData `json:"data"`
}

type HealthCheckResponseDoc struct {
	Status      string `json:"status" example:"healthy"`
	App         string `json:"app" example:"go-server"`
	Version     string `json:"version" example:"1.0.0"`
	Environment string `json:"environment" example:"development"`
	Timestamp   string `json:"timestamp" example:"2026-05-26T02:55:50Z"`
}

type ReadinessResponseDoc struct {
	Status    string `json:"status" example:"ready"`
	Database  string `json:"database" example:"connected"`
	Timestamp string `json:"timestamp" example:"2026-05-26T02:55:50Z"`
}

type HealthErrorResponseDoc struct {
	Status string `json:"status" example:"not_ready"`
	Error  string `json:"error" example:"database connection failed"`
}

type LivenessResponseDoc struct {
	Status    string `json:"status" example:"alive"`
	Timestamp string `json:"timestamp" example:"2026-05-26T02:55:50Z"`
}
