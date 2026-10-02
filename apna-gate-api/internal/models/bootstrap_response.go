package models

type DashboardKind string

const (
	DashboardKindDeveloper     DashboardKind = "developer"
	DashboardKindSocietyAdmin  DashboardKind = "society_admin"
	DashboardKindSelectSociety DashboardKind = "select_society"
	DashboardKindOnboarding    DashboardKind = "onboarding"
)

type DefaultDashboardResponse struct {
	Kind      DashboardKind `json:"kind"`
	Path      string        `json:"path"`
	SocietyID *int64        `json:"society_id,omitempty"`
}

type BootstrapResponse struct {
	User             *UserResponse             `json:"user"`
	Memberships      []*SocietyMemberResponse  `json:"memberships"`
	Residences       []*FlatResidentResponse   `json:"residences"`
	DefaultDashboard *DefaultDashboardResponse `json:"defaultDashboard"`
}
