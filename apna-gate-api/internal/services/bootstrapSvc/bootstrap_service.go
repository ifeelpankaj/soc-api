package bootstrapsvc

import (
	"context"

	"go-server/internal/models"
	service "go-server/internal/services"
)

type BootstrapService interface {
	GetBootstrap(ctx context.Context, userID int64) (*models.BootstrapResponse, error)
}

type ProfileReader interface {
	GetProfile(context.Context, int64) (*models.UserResponse, error)
}
type SocietyReader interface {
	ListMySocieties(context.Context, int64) ([]*models.MySocietyResponse, error)
	GetOnboardingBootstrap(context.Context, int64) (*models.SocietyOnboardingBootstrapResponse, error)
}
type ResidenceReader interface {
	ListMyResidences(context.Context, int64, *models.FlatResidentFilter) ([]*models.FlatResidentResponse, error)
}

type bootstrapService struct {
	sessionSvc ProfileReader
	societySvc SocietyReader
	flatSvc    ResidenceReader
}

func NewBootstrapService(
	sessionSvc ProfileReader,
	societySvc SocietyReader,
	flatSvc ResidenceReader,
) BootstrapService {
	return &bootstrapService{
		sessionSvc: sessionSvc,
		societySvc: societySvc,
		flatSvc:    flatSvc,
	}
}

func (s *bootstrapService) GetBootstrap(ctx context.Context, userID int64) (*models.BootstrapResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	user, err := s.sessionSvc.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	mySocieties, err := s.societySvc.ListMySocieties(ctx, userID)
	if err != nil {
		return nil, err
	}
	memberships := membershipsFromSocieties(mySocieties)

	residences, err := s.flatSvc.ListMyResidences(ctx, userID, &models.FlatResidentFilter{})
	if err != nil {
		return nil, err
	}

	defaultDashboard, err := s.resolveDefaultDashboard(ctx, user, mySocieties)
	if err != nil {
		return nil, err
	}

	return &models.BootstrapResponse{
		User:             user,
		Memberships:      memberships,
		Residences:       residences,
		DefaultDashboard: defaultDashboard,
	}, nil
}

func (s *bootstrapService) resolveDefaultDashboard(
	ctx context.Context,
	user *models.UserResponse,
	mySocieties []*models.MySocietyResponse,
) (*models.DefaultDashboardResponse, error) {
	if user != nil {
		switch user.GlobalRole {
		case models.GlobalRoleDeveloper, models.GlobalRoleSuperAdmin:
			return &models.DefaultDashboardResponse{
				Kind: models.DashboardKindDeveloper,
				Path: "/developer",
			}, nil
		}
	}

	adminSocieties := activeAdminSocieties(mySocieties)
	if len(adminSocieties) == 0 {
		return &models.DefaultDashboardResponse{
			Kind: models.DashboardKindOnboarding,
			Path: "/onboarding",
		}, nil
	}

	if len(adminSocieties) > 1 {
		return &models.DefaultDashboardResponse{
			Kind: models.DashboardKindSelectSociety,
			Path: "/select-society",
		}, nil
	}

	selected := adminSocieties[0]
	if selected.Society.Status == models.SocietyStatusRejected ||
		selected.Society.Status == models.SocietyStatusSuspended {
		return &models.DefaultDashboardResponse{
			Kind: models.DashboardKindSelectSociety,
			Path: "/select-society",
		}, nil
	}

	onboarding, err := s.societySvc.GetOnboardingBootstrap(ctx, selected.Society.ID)
	if err != nil {
		return nil, err
	}

	kind := models.DashboardKindOnboarding
	if onboarding.IsOnboarded {
		kind = models.DashboardKindSocietyAdmin
	}
	societyID := selected.Society.ID
	return &models.DefaultDashboardResponse{
		Kind:      kind,
		Path:      onboarding.NextPath,
		SocietyID: &societyID,
	}, nil
}

func membershipsFromSocieties(mySocieties []*models.MySocietyResponse) []*models.SocietyMemberResponse {
	memberships := make([]*models.SocietyMemberResponse, 0, len(mySocieties))
	for _, item := range mySocieties {
		if item == nil || item.Member == nil {
			continue
		}
		memberships = append(memberships, item.Member)
	}
	return memberships
}

func activeAdminSocieties(mySocieties []*models.MySocietyResponse) []*models.MySocietyResponse {
	adminSocieties := make([]*models.MySocietyResponse, 0, len(mySocieties))
	for _, item := range mySocieties {
		if item == nil || item.Member == nil || item.Society == nil ||
			item.Member.Status != models.SocietyMemberStatusActive {
			continue
		}
		if item.Member.Role == models.SocietyMemberRoleOwner ||
			item.Member.Role == models.SocietyMemberRoleAdmin {
			adminSocieties = append(adminSocieties, item)
		}
	}
	return adminSocieties
}

var _ BootstrapService = (*bootstrapService)(nil)
