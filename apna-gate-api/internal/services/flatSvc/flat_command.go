package flatsvc

import (
	"context"
	"errors"
	"go-server/internal/repositories/contracts"
	"strings"

	"go-server/internal/models"
	service "go-server/internal/services"
)

func (s *FlatSvc) CreateFlat(ctx context.Context, societyID int64, createdBy int64, req *models.CreateFlatRequest) (*models.FlatResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if req == nil {
		return nil, ErrInvalidFlatRequest
	}
	req.Sanitize()
	if err := req.Validate(); err != nil {
		return nil, ErrInvalidFlatRequest.WithCause(err)
	}
	result, err := s.createManyFlats(ctx, societyID, createdBy, []flatCreateInput{{
		Block: req.Block, Floor: req.Floor, FlatNumber: req.FlatNumber, Metadata: req.Metadata, FlatType: req.FlatType, AreaSqft: req.AreaSqft,
	}})
	if err != nil {
		return nil, err
	}
	return result.Items[0], nil
}

func (s *FlatSvc) BulkCreateFlats(ctx context.Context, societyID int64, createdBy int64, req *models.BulkCreateFlatsRequest) (*models.BulkCreateFlatsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if req == nil {
		return nil, ErrInvalidFlatRequest
	}
	req.Sanitize()
	if err := req.Validate(); err != nil {
		return nil, ErrInvalidFlatRequest.WithCause(err)
	}
	inputs := make([]flatCreateInput, 0, len(req.Flats))
	for i := range req.Flats {
		item := req.Flats[i]
		inputs = append(inputs, flatCreateInput{
			Block: item.Block, Floor: item.Floor, FlatNumber: item.FlatNumber, Metadata: item.Metadata, FlatType: item.FlatType, AreaSqft: item.AreaSqft,
		})
	}
	return s.createManyFlats(ctx, societyID, createdBy, inputs)
}

func (s *FlatSvc) GenerateFlats(ctx context.Context, societyID int64, createdBy int64, req *models.GenerateFlatsRequest) (*models.BulkCreateFlatsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if req == nil {
		return nil, ErrInvalidFlatRequest
	}
	req.Sanitize()
	if err := req.Validate(); err != nil {
		return nil, ErrInvalidFlatRequest.WithCause(err)
	}
	if err := s.ensureFlatManager(ctx, societyID, createdBy); err != nil {
		return nil, err
	}
	if err := s.ensureFlatOperational(ctx, societyID); err != nil {
		return nil, err
	}
	inputs, err := generateFlatInputs(req)
	if err != nil {
		return nil, err
	}
	return s.createManyFlatsAuthorized(ctx, societyID, createdBy, inputs)
}

func (s *FlatSvc) createManyFlats(ctx context.Context, societyID int64, createdBy int64, inputs []flatCreateInput) (*models.BulkCreateFlatsResponse, error) {
	if err := s.ensureFlatManager(ctx, societyID, createdBy); err != nil {
		return nil, err
	}
	if err := s.ensureFlatOperational(ctx, societyID); err != nil {
		return nil, err
	}
	return s.createManyFlatsAuthorized(ctx, societyID, createdBy, inputs)
}

func (s *FlatSvc) createManyFlatsAuthorized(ctx context.Context, societyID int64, createdBy int64, inputs []flatCreateInput) (*models.BulkCreateFlatsResponse, error) {
	if len(inputs) == 0 {
		return nil, ErrInvalidFlatRequest.WithCause(errors.New("at least one flat is required"))
	}
	if err := validateUniqueFlatInputs(inputs); err != nil {
		return nil, err
	}

	items := make([]*models.FlatResponse, 0, len(inputs))
	err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := s.ensureCanAddFlatsWithLock(txCtx, societyID, int64(len(inputs))); err != nil {
			return err
		}
		for _, input := range inputs {
			flat := &models.Flat{
				SocietyID: societyID, Block: input.Block, Floor: input.Floor, FlatNumber: input.FlatNumber,
				Status: models.FlatStatusVacant, IsActive: true, CreatedBy: &createdBy, Metadata: input.Metadata, FlatType: input.FlatType, AreaSqft: input.AreaSqft,
			}
			if err := s.flatRepo.Create(txCtx, flat); err != nil {
				return ErrFlatConflict.WithCause(err)
			}
			if s.visitorSettingSvc != nil {
				if err := s.visitorSettingSvc.CreateDefaultFlatSettings(txCtx, societyID, flat.ID, createdBy); err != nil {
					return err
				}
			}
			loaded, err := s.flatRepo.Get(txCtx, &models.FlatFilter{ID: &flat.ID, SocietyID: &societyID})
			if err != nil {
				return err
			}
			if loaded == nil {
				return ErrFlatNotFound
			}
			items = append(items, loaded.ToResponse())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &models.BulkCreateFlatsResponse{Items: items, Total: int32(len(items))}, nil //nolint:gosec // Input is bounded by request validation or request body size.
}

func validateUniqueFlatInputs(inputs []flatCreateInput) error {
	seen := make(map[string]struct{}, len(inputs))
	for _, input := range inputs {
		block := ""
		if input.Block != nil {
			block = strings.TrimSpace(*input.Block)
		}
		key := block + "\x00" + input.FlatNumber
		if _, exists := seen[key]; exists {
			return ErrFlatConflict
		}
		seen[key] = struct{}{}
	}
	return nil
}

func (s *FlatSvc) UpdateFlat(ctx context.Context, filter *models.FlatFilter, req *models.UpdateFlatRequest) (*models.FlatResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if req == nil || filter == nil {
		return nil, ErrInvalidFlatRequest
	}
	req.Sanitize()
	if err := req.Validate(); err != nil {
		return nil, ErrInvalidFlatRequest.WithCause(err)
	}
	flat, err := s.flatRepo.Update(ctx, filter, (*contracts.UpdateFlatInput)(req))
	if err != nil {
		return nil, err
	}
	if flat == nil {
		return nil, ErrFlatNotFound
	}
	return s.GetFlat(ctx, &models.FlatFilter{ID: &flat.ID, SocietyID: &flat.SocietyID})
}

func (s *FlatSvc) DeleteFlat(ctx context.Context, filter *models.FlatFilter, deletedBy int64) error {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if filter == nil {
		return ErrInvalidFlatRequest
	}
	if filter.SocietyID != nil {
		if err := s.ensureFlatManager(ctx, *filter.SocietyID, deletedBy); err != nil {
			return err
		}
	}
	return s.flatRepo.Deactivate(ctx, filter)
}

func (s *FlatSvc) BlockFlat(ctx context.Context, filter *models.FlatFilter, blockedBy int64, reason string) (*models.FlatResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if filter == nil {
		return nil, ErrInvalidFlatRequest
	}
	if filter.SocietyID != nil {
		if err := s.ensureFlatManager(ctx, *filter.SocietyID, blockedBy); err != nil {
			return nil, err
		}
	}
	flat, err := s.flatRepo.Block(ctx, filter)
	if err != nil {
		return nil, err
	}
	if flat == nil {
		return nil, ErrFlatNotFound
	}
	return s.GetFlat(ctx, &models.FlatFilter{ID: &flat.ID, SocietyID: &flat.SocietyID})
}

func (s *FlatSvc) UnblockFlat(ctx context.Context, filter *models.FlatFilter, unblockedBy int64) (*models.FlatResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if filter == nil {
		return nil, ErrInvalidFlatRequest
	}
	if filter.SocietyID != nil {
		if err := s.ensureFlatManager(ctx, *filter.SocietyID, unblockedBy); err != nil {
			return nil, err
		}
	}
	flat, err := s.flatRepo.Unblock(ctx, filter)
	if err != nil {
		return nil, err
	}
	if flat == nil {
		return nil, ErrFlatNotFound
	}
	return s.GetFlat(ctx, &models.FlatFilter{ID: &flat.ID, SocietyID: &flat.SocietyID})
}
