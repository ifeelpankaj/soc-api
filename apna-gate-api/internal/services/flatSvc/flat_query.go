package flatsvc

import (
	"context"

	"go-server/internal/models"
	service "go-server/internal/services"
)

func (s *FlatSvc) GetFlat(ctx context.Context, filter *models.FlatFilter) (*models.FlatResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	flat, err := s.flatRepo.Get(ctx, filter)
	if err != nil {
		return nil, err
	}
	if flat == nil {
		return nil, ErrFlatNotFound
	}
	return flat.ToResponse(), nil
}

func (s *FlatSvc) ListFlats(ctx context.Context, filter *models.FlatFilter) ([]*models.FlatResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	items, err := s.flatRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	responses := make([]*models.FlatResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, item.ToResponse())
	}
	return responses, nil
}

func (s *FlatSvc) ListFlatsPaginated(ctx context.Context, filter *models.FlatFilter) (*models.PaginatedFlatsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	items, err := s.flatRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	total, err := s.flatRepo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}
	responses := make([]*models.FlatResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, item.ToResponse())
	}
	return &models.PaginatedFlatsResponse{
		Items:  responses,
		Total:  total,
		Limit:  filter.Limit,
		Offset: filter.Offset,
	}, nil
}

func (s *FlatSvc) GetFlatStats(ctx context.Context, societyID int64) (*models.FlatStatsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	return s.flatRepo.Stats(ctx, societyID)
}
