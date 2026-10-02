package subscriptionsvc

import (
	"context"
	"time"

	"go-server/internal/models"
	service "go-server/internal/services"
)

func (s *SubscriptionSvc) GetSubscription(ctx context.Context, filter *models.SubscriptionFilter) (*models.SocietySubscriptionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	sub, err := s.subRepo.Get(ctx, filter)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return nil, ErrSubscriptionNotFound
	}
	return sub.ToResponse(), nil
}

func (s *SubscriptionSvc) GetActiveSubscriptionBySocietyID(ctx context.Context, societyID int64) (*models.SocietySubscriptionResponse, error) {
	activeOnly := true
	sub, err := s.subRepo.Get(ctx, &models.SubscriptionFilter{SocietyID: &societyID, IsActiveOnly: &activeOnly})
	if err != nil {
		return nil, err
	}
	if sub != nil {
		return sub.ToResponse(), nil
	}

	latest, err := s.subRepo.List(ctx, &models.SubscriptionFilter{SocietyID: &societyID, Limit: 1})
	if err != nil {
		return nil, err
	}
	if len(latest) == 0 {
		return nil, ErrSubscriptionNotFound
	}

	item := latest[0]
	if item.Status == models.SubscriptionStatusExpired ||
		item.Status == models.SubscriptionStatusCancelled ||
		(item.EndsAt != nil && !item.EndsAt.After(time.Now().UTC())) {
		return nil, ErrSubscriptionExpired
	}

	return nil, ErrSubscriptionNotFound
}

func (s *SubscriptionSvc) ListSubscriptions(ctx context.Context, filter *models.SubscriptionFilter) ([]*models.SocietySubscriptionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	items, err := s.subRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	responses := make([]*models.SocietySubscriptionResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, item.ToResponse())
	}
	return responses, nil
}

func (s *SubscriptionSvc) GetSubscriptionStats(ctx context.Context, filter *models.SubscriptionFilter) (*models.SubscriptionStatsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	return s.subRepo.Stats(ctx, filter)
}
