package societysvc

import (
	"context"
	"fmt"
	"go-server/internal/repositories/contracts"
	"strings"

	"go-server/internal/models"
	service "go-server/internal/services"
	"go-server/pkg/utils"
)

func (s *SocietySvc) CreateSocietyRequest(ctx context.Context, req models.CreateSocietyRequest, requestedBy int64) (*models.SocietyResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	req.Sanitize()
	if err := s.ensureUserExists(ctx, requestedBy); err != nil {
		return nil, err
	}
	if err := s.ensureNoDuplicateSociety(ctx, req, requestedBy); err != nil {
		return nil, err
	}

	country := "India"
	if req.Country != nil {
		country = *req.Country
	}

	code := ""
	if req.SocietyCode != nil {
		code = *req.SocietyCode
	} else {
		generatedCode, err := utils.GenerateSocietyCode(
			req.Name,
			stringValue(req.City),
			stringValue(req.State),
			stringValue(req.Pincode),
		)
		if err != nil {
			return nil, fmt.Errorf("generate society code: %w", err)
		}
		code = generatedCode
	}

	society := &models.Society{
		Name:         req.Name,
		SocietyCode:  code,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		Landmark:     req.Landmark,
		City:         req.City,
		State:        req.State,
		Pincode:      req.Pincode,
		Country:      country,
		TotalFlats:   req.TotalFlats,
		TotalBlocks:  req.TotalBlocks,
		Status:       models.SocietyStatusPending,
		CreatedBy:    requestedBy,
		Metadata:     req.Metadata,
	}

	if err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := s.societyRepo.Create(txCtx, society); err != nil {
			return ErrSocietyConflict.WithCause(err)
		}

		creator := &models.SocietyMember{
			SocietyID: society.ID,
			UserID:    requestedBy,
			Role:      models.SocietyMemberRoleOwner,
			Status:    models.SocietyMemberStatusActive,
		}
		if err := s.memberRepo.Add(txCtx, creator); err != nil {
			return ErrMemberConflict.WithCause(err)
		}
		if s.visitorSettingSvc != nil {
			if err := s.visitorSettingSvc.CreateDefaultSocietySettings(txCtx, society.ID, requestedBy); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return society.ToResponse(), nil
}

func (s *SocietySvc) ApproveSociety(ctx context.Context, societyID int64, approvedBy int64) (*models.SocietyResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := s.ensureLifecycleStatus(ctx, societyID, models.SocietyStatusPending); err != nil {
		return nil, err
	}
	society, err := s.societyRepo.Approve(ctx, societyID, approvedBy)
	return societyResponseOrTransition(society, err)
}

func (s *SocietySvc) RejectSociety(ctx context.Context, societyID int64, rejectedBy int64, reason string) (*models.SocietyResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := s.ensureLifecycleStatus(ctx, societyID, models.SocietyStatusPending); err != nil {
		return nil, err
	}
	society, err := s.societyRepo.Reject(ctx, societyID, rejectedBy, strings.TrimSpace(reason))
	return societyResponseOrTransition(society, err)
}

func (s *SocietySvc) SuspendSociety(ctx context.Context, societyID int64, suspendedBy int64, reason string) (*models.SocietyResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := s.ensureLifecycleStatus(ctx, societyID, models.SocietyStatusActive); err != nil {
		return nil, err
	}
	society, err := s.societyRepo.Suspend(ctx, societyID, suspendedBy, strings.TrimSpace(reason))
	return societyResponseOrTransition(society, err)
}

func (s *SocietySvc) ReactivateSociety(ctx context.Context, societyID int64, reactivatedBy int64) (*models.SocietyResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := s.ensureLifecycleStatus(ctx, societyID, models.SocietyStatusSuspended); err != nil {
		return nil, err
	}
	society, err := s.societyRepo.Reactivate(ctx, societyID, reactivatedBy)
	return societyResponseOrTransition(society, err)
}

func (s *SocietySvc) RestoreSociety(ctx context.Context, societyID int64, restoredBy int64) (*models.SocietyResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := s.ensureUserExists(ctx, restoredBy); err != nil {
		return nil, err
	}

	existing, err := s.societyRepo.Get(ctx, models.GetSocietyFilter{ID: &societyID, IncludeDeleted: true})
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrSocietyNotFound
	}
	if existing.DeletedAt == nil {
		return nil, ErrInvalidTransition
	}

	society, err := s.societyRepo.Restore(ctx, societyID)
	if err != nil {
		return nil, ErrInvalidTransition.WithCause(err)
	}
	return society.ToResponse(), nil
}

func (s *SocietySvc) UpdateSociety(ctx context.Context, societyID int64, req models.UpdateSocietyRequest, updatedBy int64) (*models.SocietyResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	req.Sanitize()
	if err := s.EnsureRole(ctx, societyID, updatedBy, string(models.SocietyMemberRoleOwner), string(models.SocietyMemberRoleAdmin)); err != nil {
		return nil, err
	}

	society, err := s.societyRepo.Update(ctx, societyID, contracts.UpdateSocietyInput(req))
	if err != nil {
		return nil, err
	}
	if society == nil {
		return nil, ErrSocietyNotFound
	}
	return society.ToResponse(), nil
}

func (s *SocietySvc) DeleteSociety(ctx context.Context, societyID int64, deletedBy int64) error {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := s.EnsureRole(ctx, societyID, deletedBy, string(models.SocietyMemberRoleOwner)); err != nil {
		return err
	}
	return s.societyRepo.SoftDelete(ctx, societyID)
}

func societyResponseOrTransition(society *models.Society, err error) (*models.SocietyResponse, error) {
	if err != nil {
		return nil, err
	}
	if society == nil {
		return nil, ErrInvalidTransition
	}
	return society.ToResponse(), nil
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
