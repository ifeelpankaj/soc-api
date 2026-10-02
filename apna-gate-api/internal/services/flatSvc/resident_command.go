package flatsvc

import (
	"context"

	"go-server/internal/models"
	service "go-server/internal/services"
)

func (s *FlatSvc) AddResidentToFlat(ctx context.Context, societyID int64, flatID int64, userID int64, createdBy int64, req *models.AddFlatResidentRequest) (*models.FlatResidentResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if req == nil {
		return nil, ErrInvalidFlatRequest
	}
	if err := req.Validate(); err != nil {
		return nil, ErrInvalidResidentRole.WithCause(err)
	}
	if err := s.ensureFlatManager(ctx, societyID, createdBy); err != nil {
		return nil, err
	}
	if err := s.ensureFlatOperational(ctx, societyID); err != nil {
		return nil, err
	}

	var resident *models.FlatResident
	err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if s.subscriptionSvc != nil {
			if err := s.subscriptionSvc.CanAddResidentWithLock(txCtx, societyID, 1); err != nil {
				return err
			}
		}
		if err := s.ensureFlatAssignable(txCtx, societyID, flatID); err != nil {
			return err
		}
		if req.IsPrimary {
			count, err := s.residentRepo.CountPrimary(txCtx, societyID, flatID)
			if err != nil {
				return err
			}
			if count > 0 {
				return ErrPrimaryResidentExists
			}
		}
		if _, err := s.memberRepo.UpsertResident(txCtx, societyID, userID, createdBy); err != nil {
			return err
		}
		resident = &models.FlatResident{
			SocietyID: societyID, FlatID: flatID, UserID: userID, Role: req.Role,
			Status: models.FlatResidentStatusActive, IsPrimary: req.IsPrimary,
			CreatedBy: &createdBy, Metadata: req.Metadata,
		}
		if err := s.residentRepo.Add(txCtx, resident); err != nil {
			return ErrResidentConflict.WithCause(err)
		}
		_, err := s.flatRepo.MarkOccupied(txCtx, societyID, flatID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return s.GetFlatResident(ctx, &models.FlatResidentFilter{ID: &resident.ID, SocietyID: &societyID})
}

func (s *FlatSvc) RemoveResidentFromFlat(ctx context.Context, filter *models.FlatResidentFilter, removedBy int64) error {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	resident, err := s.requireResident(ctx, filter)
	if err != nil {
		return err
	}
	if err := s.ensureFlatManager(ctx, resident.SocietyID, removedBy); err != nil {
		return err
	}
	return s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := s.residentRepo.Remove(txCtx, filter); err != nil {
			return err
		}
		return s.markVacantIfEmpty(txCtx, resident.SocietyID, resident.FlatID)
	})
}

func (s *FlatSvc) ChangePrimaryResident(ctx context.Context, societyID int64, flatID int64, residentID int64, changedBy int64) (*models.FlatResidentResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := s.ensureFlatManager(ctx, societyID, changedBy); err != nil {
		return nil, err
	}
	var resident *models.FlatResident
	err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		current, err := s.residentRepo.Get(txCtx, &models.FlatResidentFilter{ID: &residentID, SocietyID: &societyID, FlatID: &flatID})
		if err != nil {
			return err
		}
		if current == nil || current.Status != models.FlatResidentStatusActive {
			return ErrResidentNotFound
		}
		resident, err = s.residentRepo.SetPrimary(txCtx, societyID, flatID, residentID)
		if err != nil {
			return err
		}
		if resident == nil {
			return ErrResidentNotFound
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.GetFlatResident(ctx, &models.FlatResidentFilter{ID: &resident.ID, SocietyID: &societyID})
}

func (s *FlatSvc) UpdateFlatResidentRole(ctx context.Context, filter *models.FlatResidentFilter, updatedBy int64, req *models.UpdateFlatResidentRoleRequest) (*models.FlatResidentResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if req == nil {
		return nil, ErrInvalidResidentRole
	}
	if err := req.Validate(); err != nil {
		return nil, ErrInvalidResidentRole.WithCause(err)
	}
	current, err := s.requireResident(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err := s.ensureFlatManager(ctx, current.SocietyID, updatedBy); err != nil {
		return nil, err
	}
	resident, err := s.residentRepo.UpdateRole(ctx, filter, req.Role)
	if err != nil {
		return nil, err
	}
	if resident == nil {
		return nil, ErrResidentNotFound
	}
	return s.GetFlatResident(ctx, &models.FlatResidentFilter{ID: &resident.ID, SocietyID: &resident.SocietyID})
}

func (s *FlatSvc) MoveOutResident(ctx context.Context, filter *models.FlatResidentFilter, movedOutBy int64) (*models.FlatResidentResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	current, err := s.requireResident(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err := s.ensureFlatManager(ctx, current.SocietyID, movedOutBy); err != nil {
		return nil, err
	}
	var moved *models.FlatResident
	err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		moved, err = s.residentRepo.MoveOut(txCtx, filter)
		if err != nil {
			return err
		}
		if moved == nil {
			return ErrResidentNotFound
		}
		return s.markVacantIfEmpty(txCtx, current.SocietyID, current.FlatID)
	})
	if err != nil {
		return nil, err
	}
	return s.GetFlatResident(ctx, &models.FlatResidentFilter{ID: &moved.ID, SocietyID: &moved.SocietyID})
}

func (s *FlatSvc) markVacantIfEmpty(ctx context.Context, societyID int64, flatID int64) error {
	count, err := s.residentRepo.CountActive(ctx, societyID, flatID)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = s.flatRepo.MarkVacant(ctx, societyID, flatID)
	}
	return err
}
