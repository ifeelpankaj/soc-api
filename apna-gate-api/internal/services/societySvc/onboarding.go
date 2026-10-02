package societysvc

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go-server/internal/models"
	service "go-server/internal/services"
	"go-server/pkg/utils"
)

const (
	missingStepFlats = "flats"
	missingStepStaff = "staff"
)

func (s *SocietySvc) GetOnboardingBootstrap(ctx context.Context, societyID int64) (*models.SocietyOnboardingBootstrapResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	society, err := s.societyRepo.Get(ctx, models.GetSocietyFilter{ID: &societyID})
	if err != nil {
		return nil, err
	}
	if society == nil {
		return nil, ErrSocietyNotFound
	}

	stats, err := s.flatRepo.Stats(ctx, societyID)
	if err != nil {
		return nil, err
	}

	staffRole := string(models.SocietyMemberRoleStaff)
	activeStatus := string(models.SocietyMemberStatusActive)
	staffCount, err := s.memberRepo.Count(ctx, models.ListSocietyMembersFilter{
		SocietyID: societyID,
		Role:      &staffRole,
		Status:    &activeStatus,
	})
	if err != nil {
		return nil, err
	}

	flatCount := int64(0)
	if stats != nil {
		flatCount = stats.ActiveFlats
	}

	hasFlats := flatCount > 0
	hasStaff := staffCount > 0
	currentSubscription, err := s.dashboardCurrentSubscription(ctx, societyID)
	if err != nil {
		return nil, err
	}
	hasActiveSubscription := currentSubscription != nil &&
		(currentSubscription.Status == models.SubscriptionStatusActive ||
			currentSubscription.Status == models.SubscriptionStatusTrial)
	isOnboarded := society.Status == models.SocietyStatusActive &&
		hasActiveSubscription && hasFlats && hasStaff
	missingSteps := make([]string, 0, 2)
	if !hasFlats {
		missingSteps = append(missingSteps, missingStepFlats)
	}
	if !hasStaff {
		missingSteps = append(missingSteps, missingStepStaff)
	}

	encodedID := encodeSocietyRouteID(societyID)
	nextPath := fmt.Sprintf("/onboarding/%s", encodedID)
	if isOnboarded {
		nextPath = fmt.Sprintf("/dashboard/%s", encodedID)
	}

	return &models.SocietyOnboardingBootstrapResponse{
		Society:               society.ToResponse(),
		HasActiveSubscription: hasActiveSubscription,
		HasFlats:              hasFlats,
		HasStaff:              hasStaff,
		IsOnboarded:           isOnboarded,
		FlatCount:             flatCount,
		StaffCount:            staffCount,
		MissingSteps:          missingSteps,
		NextPath:              nextPath,
	}, nil
}

func (s *SocietySvc) CreateGuard(ctx context.Context, societyID int64, actorID int64, req models.CreateGuardRequest) (*models.CreateGuardResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	req.Sanitize()
	if err := s.EnsureRole(ctx, societyID, actorID, string(models.SocietyMemberRoleOwner), string(models.SocietyMemberRoleAdmin)); err != nil {
		return nil, err
	}
	if s.subscriptionSvc != nil {
		if err := s.subscriptionSvc.CanAddStaff(ctx, societyID, 1); err != nil {
			return nil, err
		}
	}

	emailExists, err := s.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, ErrDuplicateGuardEmail
	}

	phoneExists, err := s.userRepo.PhoneExists(ctx, req.PhoneNumber)
	if err != nil {
		return nil, err
	}
	if phoneExists {
		return nil, ErrDuplicateGuardPhone
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	var user *models.User
	var member *models.SocietyMember
	err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		emailExists, err := s.userRepo.EmailExists(txCtx, req.Email)
		if err != nil {
			return err
		}
		if emailExists {
			return ErrDuplicateGuardEmail
		}

		phoneExists, err := s.userRepo.PhoneExists(txCtx, req.PhoneNumber)
		if err != nil {
			return err
		}
		if phoneExists {
			return ErrDuplicateGuardPhone
		}

		firstName := req.FirstName
		lastName := req.LastName
		email := req.Email
		phone := req.PhoneNumber
		fullName := strings.TrimSpace(strings.Join([]string{firstName, lastName}, " "))
		password := passwordHash

		user = &models.User{
			FirstName:     &firstName,
			Email:         &email,
			PhoneNumber:   &phone,
			FullName:      fullName,
			PasswordHash:  &password,
			AuthProvider:  models.AuthProviderEmail,
			GlobalRole:    models.GlobalRoleUser,
			EmailVerified: true,
			PhoneVerified: false,
			IsActive:      true,
			IsBlocked:     false,
			Timezone:      "Asia/Kolkata",
			Language:      "en",
			Metadata:      map[string]any{"created_as": "guard"},
		}
		if lastName != "" {
			user.LastName = &lastName
		}

		if err := s.userRepo.Create(txCtx, user); err != nil {
			return err
		}

		member = &models.SocietyMember{
			SocietyID: societyID,
			UserID:    user.ID,
			Role:      models.SocietyMemberRoleStaff,
			Status:    models.SocietyMemberStatusActive,
			InvitedBy: &actorID,
			Metadata:  map[string]any{"created_as": "guard"},
		}
		if err := s.memberRepo.Add(txCtx, member); err != nil {
			return ErrMemberConflict.WithCause(err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	loaded, _ := s.memberRepo.Get(ctx, models.GetSocietyMemberFilter{SocietyID: &societyID, UserID: &user.ID})
	if loaded != nil {
		member = loaded
	}

	return &models.CreateGuardResponse{
		User:   user.ToResponse(),
		Member: member.ToResponse(),
	}, nil
}

func encodeSocietyRouteID(id int64) string {
	return fmt.Sprintf("soc%s", strings.ToLower(strconv.FormatInt(id, 36)))
}
