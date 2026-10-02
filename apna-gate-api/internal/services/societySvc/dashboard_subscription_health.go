package societysvc

import (
	"time"

	"go-server/internal/constants"
	"go-server/internal/models"
)

func computeSubscriptionHealth(
	current *models.SocietySubscriptionResponse,
	latest *models.SocietySubscriptionResponse,
) *models.SocietyDashboardSubscriptionHealthResponse {
	source := current
	if source == nil {
		source = latest
	}
	if source == nil {
		return &models.SocietyDashboardSubscriptionHealthResponse{
			IsActive:       false,
			IsExpiringSoon: false,
			LifecycleLabel: "none",
		}
	}

	now := time.Now().UTC()
	isDateExpired := source.EndsAt != nil && !source.EndsAt.After(now)
	isTrialExpired := source.Status == models.SubscriptionStatusTrial &&
		source.TrialEndsAt != nil &&
		!source.TrialEndsAt.After(now)

	isActive := current != nil &&
		(source.Status == models.SubscriptionStatusActive || source.Status == models.SubscriptionStatusTrial) &&
		!isDateExpired &&
		!isTrialExpired

	var daysUntilExpiry *int64
	isExpiringSoon := false

	if source.EndsAt != nil && !isDateExpired {
		days := int64(source.EndsAt.Sub(now).Hours() / 24)
		if days < 0 {
			days = 0
		}
		daysUntilExpiry = &days
		if isActive && days <= int64(constants.SubscriptionExpiringSoonDays) {
			isExpiringSoon = true
		}
	}

	lifecycleLabel := deriveSubscriptionLifecycleLabel(source, isActive, isDateExpired, isTrialExpired, isExpiringSoon)

	return &models.SocietyDashboardSubscriptionHealthResponse{
		IsActive:        isActive,
		IsExpiringSoon:  isExpiringSoon,
		DaysUntilExpiry: daysUntilExpiry,
		LifecycleLabel:  lifecycleLabel,
	}
}

func deriveSubscriptionLifecycleLabel(
	sub *models.SocietySubscriptionResponse,
	isActive bool,
	isDateExpired bool,
	isTrialExpired bool,
	isExpiringSoon bool,
) string {
	if sub.Status == models.SubscriptionStatusCancelled {
		return "cancelled"
	}
	if sub.Status == models.SubscriptionStatusPending {
		return "pending"
	}
	if sub.Status == models.SubscriptionStatusExpired || isDateExpired || isTrialExpired {
		return "expired"
	}
	if isExpiringSoon {
		return "expiring_soon"
	}
	if sub.Status == models.SubscriptionStatusTrial {
		return "trial"
	}
	if isActive {
		return "active"
	}
	return "none"
}
