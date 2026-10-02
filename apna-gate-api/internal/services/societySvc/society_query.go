package societysvc

import (
	"context"
	"strings"

	"go-server/internal/models"
	service "go-server/internal/services"
)

func (s *SocietySvc) GetSociety(ctx context.Context, filter models.GetSocietyFilter) (*models.SocietyDetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	society, err := s.societyRepo.Get(ctx, filter)
	if err != nil {
		return nil, err
	}
	if society == nil {
		return nil, ErrSocietyNotFound
	}

	resp := &models.SocietyDetailResponse{SocietyResponse: society.ToResponse()}
	if filter.ID != nil {
		total, err := s.memberRepo.Count(ctx, models.ListSocietyMembersFilter{SocietyID: *filter.ID})
		if err == nil {
			resp.MembersCount = total
		}
	}
	return resp, nil
}

func (s *SocietySvc) ResolveActiveSocietyIDByCode(ctx context.Context, societyCode string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	code := strings.ToUpper(strings.TrimSpace(societyCode))
	if code == "" {
		return 0, ErrSocietyNotFound
	}
	active := string(models.SocietyStatusActive)
	society, err := s.societyRepo.Get(ctx, models.GetSocietyFilter{Code: &code, Status: &active})
	if err != nil {
		return 0, err
	}
	if society == nil {
		return 0, ErrSocietyNotFound
	}
	return society.ID, nil
}

func (s *SocietySvc) GetPublicClaimOptions(ctx context.Context, societyCode string) (*models.PublicClaimOptionsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	code := strings.ToUpper(strings.TrimSpace(societyCode))
	if code == "" {
		return nil, ErrSocietyNotFound
	}
	active := string(models.SocietyStatusActive)
	society, err := s.societyRepo.Get(ctx, models.GetSocietyFilter{Code: &code, Status: &active})
	if err != nil {
		return nil, err
	}
	if society == nil {
		return nil, ErrSocietyNotFound
	}

	isActive := true
	flats := []*models.Flat{}
	for offset := int32(0); ; offset += 100 {
		page, err := s.flatRepo.List(ctx, &models.FlatFilter{
			SocietyID: &society.ID,
			IsActive:  &isActive,
			Limit:     100,
			Offset:    offset,
		})
		if err != nil {
			return nil, err
		}
		flats = append(flats, page...)
		if len(page) < 100 {
			break
		}
	}

	publicFlats := make([]*models.PublicClaimFlatResponse, 0, len(flats))
	for _, flat := range flats {
		publicFlats = append(publicFlats, &models.PublicClaimFlatResponse{
			ID:         flat.ID,
			Block:      flat.Block,
			Floor:      flat.Floor,
			FlatNumber: flat.FlatNumber,
			Status:     flat.Status,
		})
	}

	return &models.PublicClaimOptionsResponse{
		Society: &models.PublicClaimSocietyResponse{
			ID:          society.ID,
			Name:        society.Name,
			SocietyCode: society.SocietyCode,
			City:        society.City,
			State:       society.State,
			Pincode:     society.Pincode,
			Country:     society.Country,
			TotalFlats:  society.TotalFlats,
		},
		Flats: publicFlats,
	}, nil
}

func (s *SocietySvc) ListSocieties(ctx context.Context, filter models.ListSocietiesFilter) (*models.PaginatedSocietiesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	items, err := s.societyRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	total, err := s.societyRepo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]*models.SocietyResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, item.ToResponse())
	}

	return &models.PaginatedSocietiesResponse{
		Items:  responses,
		Total:  total,
		Limit:  normalizeResponseLimit(filter.Limit),
		Offset: normalizeResponseOffset(filter.Offset),
	}, nil
}

func (s *SocietySvc) ListMySocieties(ctx context.Context, userID int64) ([]*models.MySocietyResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	return s.memberRepo.ListMySocietiesByUser(ctx, userID)
}

func normalizeResponseLimit(limit int32) int32 {
	if limit <= 0 {
		return 20
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func normalizeResponseOffset(offset int32) int32 {
	if offset < 0 {
		return 0
	}
	return offset
}
