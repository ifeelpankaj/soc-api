package flatsvc

import (
	"context"

	"go-server/internal/models"
	service "go-server/internal/services"
)

func (s *FlatSvc) GetFlatResident(ctx context.Context, filter *models.FlatResidentFilter) (*models.FlatResidentResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	resident, err := s.residentRepo.Get(ctx, filter)
	if err != nil {
		return nil, err
	}
	if resident == nil {
		return nil, ErrResidentNotFound
	}
	return resident.ToResponse(), nil
}

func (s *FlatSvc) ListFlatResidents(ctx context.Context, filter *models.FlatResidentFilter) ([]*models.FlatResidentResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	items, err := s.residentRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	responses := make([]*models.FlatResidentResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, item.ToResponse())
	}
	return responses, nil
}

func (s *FlatSvc) ListMyResidences(ctx context.Context, userID int64, filter *models.FlatResidentFilter) ([]*models.FlatResidentResponse, error) {
	if filter == nil {
		filter = &models.FlatResidentFilter{}
	}
	filter.UserID = &userID
	return s.ListFlatResidents(ctx, filter)
}
