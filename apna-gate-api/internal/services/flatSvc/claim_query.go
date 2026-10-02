package flatsvc

import (
	"context"

	"go-server/internal/models"
	service "go-server/internal/services"
)

func (s *FlatSvc) GetFlatClaim(ctx context.Context, filter *models.FlatClaimFilter) (*models.FlatClaimResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	claim, err := s.claimRepo.Get(ctx, filter)
	if err != nil {
		return nil, err
	}
	if claim == nil {
		return nil, ErrClaimNotFound
	}
	return claim.ToResponse(), nil
}

func (s *FlatSvc) ListFlatClaims(ctx context.Context, filter *models.FlatClaimFilter) ([]*models.FlatClaimResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	items, err := s.claimRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	responses := make([]*models.FlatClaimResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, item.ToResponse())
	}
	return responses, nil
}

func (s *FlatSvc) ListMyFlatClaims(ctx context.Context, userID int64, filter *models.FlatClaimFilter) ([]*models.FlatClaimResponse, error) {
	if filter == nil {
		filter = &models.FlatClaimFilter{}
	}
	filter.UserID = &userID
	return s.ListFlatClaims(ctx, filter)
}

func (s *FlatSvc) GetFlatClaimStats(ctx context.Context, societyID int64) (*models.FlatClaimStatsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	return s.claimRepo.Stats(ctx, societyID)
}
