package flatsvc

import (
	"context"

	"go-server/internal/models"
	service "go-server/internal/services"
)

func (s *FlatSvc) SubmitFlatClaim(ctx context.Context, userID int64, req *models.SubmitFlatClaimRequest) (*models.FlatClaimResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if req == nil {
		return nil, ErrInvalidClaimRequest
	}
	req.Sanitize()
	if err := req.Validate(); err != nil {
		return nil, ErrInvalidClaimRequest.WithCause(err)
	}
	flat, err := s.flatRepo.Get(ctx, &models.FlatFilter{ID: &req.FlatID, SocietyID: &req.SocietyID})
	if err != nil {
		return nil, err
	}
	if flat == nil {
		return nil, ErrFlatNotFound
	}
	if !flat.IsActive {
		return nil, ErrFlatInactive
	}
	if flat.Status == models.FlatStatusBlocked {
		return nil, ErrFlatBlocked
	}

	claim := &models.FlatClaim{
		SocietyID: req.SocietyID, FlatID: req.FlatID, UserID: userID,
		RequestedRole: req.RequestedRole, RequestedPrimary: req.RequestedPrimary,
		Status: models.FlatClaimStatusPending, Note: req.Note, Metadata: req.Metadata,
	}
	if err := s.claimRepo.Submit(ctx, claim); err != nil {
		return nil, ErrClaimConflict.WithCause(err)
	}
	return s.GetFlatClaim(ctx, &models.FlatClaimFilter{ID: &claim.ID, SocietyID: &claim.SocietyID})
}

func (s *FlatSvc) ApproveFlatClaim(ctx context.Context, societyID int64, claimID int64, reviewedBy int64) (*models.ApproveFlatClaimResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	var approved *models.FlatClaim
	var resident *models.FlatResident
	var flat *models.Flat

	err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := s.ensureFlatOperational(txCtx, societyID); err != nil {
			return err
		}
		if s.subscriptionSvc != nil {
			if err := s.subscriptionSvc.CanAddResidentWithLock(txCtx, societyID, 1); err != nil {
				return err
			}
		}
		pending := string(models.FlatClaimStatusPending)
		claim, err := s.claimRepo.Get(txCtx, &models.FlatClaimFilter{ID: &claimID, SocietyID: &societyID, Status: &pending})
		if err != nil {
			return err
		}
		if claim == nil {
			return ErrClaimNotFound
		}

		currentFlat, err := s.flatRepo.Get(txCtx, &models.FlatFilter{ID: &claim.FlatID, SocietyID: &claim.SocietyID})
		if err != nil {
			return err
		}
		if currentFlat == nil {
			return ErrFlatNotFound
		}
		if currentFlat.SocietyID != claim.SocietyID {
			return ErrInvalidClaimRequest
		}
		if !currentFlat.IsActive {
			return ErrFlatInactive
		}
		if currentFlat.Status == models.FlatStatusBlocked {
			return ErrFlatBlocked
		}
		if err := s.ensureFlatManager(txCtx, societyID, reviewedBy); err != nil {
			return err
		}
		if claim.RequestedPrimary {
			count, err := s.residentRepo.CountPrimary(txCtx, claim.SocietyID, claim.FlatID)
			if err != nil {
				return err
			}
			if count > 0 {
				return ErrPrimaryResidentExists
			}
		}
		if _, err := s.memberRepo.UpsertResident(txCtx, claim.SocietyID, claim.UserID, reviewedBy); err != nil {
			return err
		}
		resident = &models.FlatResident{
			SocietyID: claim.SocietyID, FlatID: claim.FlatID, UserID: claim.UserID,
			Role: claim.RequestedRole, Status: models.FlatResidentStatusActive,
			IsPrimary: claim.RequestedPrimary, CreatedBy: &reviewedBy,
		}
		if err := s.residentRepo.Add(txCtx, resident); err != nil {
			return ErrResidentConflict.WithCause(err)
		}
		if s.visitorSettingSvc != nil {
			if err := s.visitorSettingSvc.CreateDefaultFlatSettings(txCtx, claim.SocietyID, claim.FlatID, reviewedBy); err != nil {
				return err
			}
		}
		var errMark error
		flat, errMark = s.flatRepo.MarkOccupied(txCtx, claim.SocietyID, claim.FlatID)
		if errMark != nil {
			return errMark
		}
		approved, err = s.claimRepo.Approve(txCtx, societyID, claimID, reviewedBy)
		if err != nil {
			return err
		}
		if approved == nil {
			return ErrClaimNotFound
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	claimResp, err := s.GetFlatClaim(ctx, &models.FlatClaimFilter{ID: &approved.ID, SocietyID: &societyID})
	if err != nil {
		return nil, err
	}
	residentResp, err := s.GetFlatResident(ctx, &models.FlatResidentFilter{ID: &resident.ID, SocietyID: &societyID})
	if err != nil {
		return nil, err
	}
	flatResp, err := s.GetFlat(ctx, &models.FlatFilter{ID: &flat.ID, SocietyID: &societyID})
	if err != nil {
		return nil, err
	}
	return &models.ApproveFlatClaimResponse{Claim: claimResp, Resident: residentResp, Flat: flatResp}, nil
}

func (s *FlatSvc) RejectFlatClaim(ctx context.Context, societyID int64, claimID int64, reviewedBy int64, req *models.RejectFlatClaimRequest) (*models.FlatClaimResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if req == nil {
		return nil, ErrInvalidClaimRequest
	}
	req.Sanitize()
	if req.Reason == "" {
		return nil, ErrInvalidClaimRequest
	}
	if err := s.ensureFlatManager(ctx, societyID, reviewedBy); err != nil {
		return nil, err
	}
	claim, err := s.claimRepo.Reject(ctx, societyID, claimID, reviewedBy, req.Reason)
	if err != nil {
		return nil, err
	}
	if claim == nil {
		return nil, ErrClaimNotFound
	}
	return s.GetFlatClaim(ctx, &models.FlatClaimFilter{ID: &claim.ID, SocietyID: &societyID})
}

func (s *FlatSvc) CancelMyFlatClaim(ctx context.Context, claimID int64, userID int64) (*models.FlatClaimResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	claim, err := s.claimRepo.Cancel(ctx, claimID, userID)
	if err != nil {
		return nil, err
	}
	if claim == nil {
		return nil, ErrClaimNotFound
	}
	return s.GetFlatClaim(ctx, &models.FlatClaimFilter{ID: &claim.ID, UserID: &userID})
}
