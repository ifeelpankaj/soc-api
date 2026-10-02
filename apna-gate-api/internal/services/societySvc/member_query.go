package societysvc

import (
	"context"

	"go-server/internal/models"
	service "go-server/internal/services"
)

func (s *SocietySvc) GetSocietyMember(ctx context.Context, filter models.GetSocietyMemberFilter) (*models.SocietyMemberResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	member, err := s.memberRepo.Get(ctx, filter)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, ErrMemberNotFound
	}
	return member.ToResponse(), nil
}

func (s *SocietySvc) GetSocietyMemberDetail(ctx context.Context, filter models.GetSocietyMemberFilter) (*models.SocietyMemberDetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	member, err := s.memberRepo.Get(ctx, filter)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, ErrMemberNotFound
	}

	detail := &models.SocietyMemberDetailResponse{
		Member:     member.ToResponse(),
		OwnedFlats: []*models.FlatResidentResponse{},
		Residences: []*models.FlatResidentResponse{},
	}
	if s.flatResidentRepo == nil {
		return detail, nil
	}

	activeStatus := string(models.FlatResidentStatusActive)
	ownerRole := string(models.FlatResidentRoleOwner)
	residences, err := s.flatResidentRepo.List(ctx, &models.FlatResidentFilter{
		SocietyID: &member.SocietyID,
		UserID:    &member.UserID,
		Status:    &activeStatus,
		Limit:     100,
	})
	if err != nil {
		return nil, err
	}
	for _, residence := range residences {
		response := residence.ToResponse()
		detail.Residences = append(detail.Residences, response)
		if residence.Role == models.FlatResidentRoleOwner {
			detail.OwnedFlats = append(detail.OwnedFlats, response)
		}
	}

	owned, err := s.flatResidentRepo.List(ctx, &models.FlatResidentFilter{
		SocietyID: &member.SocietyID,
		UserID:    &member.UserID,
		Role:      &ownerRole,
		Status:    &activeStatus,
		Limit:     100,
	})
	if err == nil && len(owned) > len(detail.OwnedFlats) {
		detail.OwnedFlats = detail.OwnedFlats[:0]
		for _, residence := range owned {
			detail.OwnedFlats = append(detail.OwnedFlats, residence.ToResponse())
		}
	}
	return detail, nil
}

func (s *SocietySvc) GetSocietyMemberSummary(ctx context.Context, societyID int64) (*models.SocietyMemberSummaryResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	count := func(role *string, status *string) (int64, error) {
		return s.memberRepo.Count(ctx, models.ListSocietyMembersFilter{
			SocietyID: societyID,
			Role:      role,
			Status:    status,
		})
	}
	active, pending, suspended, removed := string(models.SocietyMemberStatusActive), string(models.SocietyMemberStatusPending), string(models.SocietyMemberStatusSuspended), string(models.SocietyMemberStatusRemoved)
	owner, admin, staff, resident := string(models.SocietyMemberRoleOwner), string(models.SocietyMemberRoleAdmin), string(models.SocietyMemberRoleStaff), string(models.SocietyMemberRoleResident)

	total, err := count(nil, nil)
	if err != nil {
		return nil, err
	}
	summary := &models.SocietyMemberSummaryResponse{TotalMembers: total}
	if summary.ActiveMembers, err = count(nil, &active); err != nil {
		return nil, err
	}
	if summary.PendingMembers, err = count(nil, &pending); err != nil {
		return nil, err
	}
	if summary.SuspendedMembers, err = count(nil, &suspended); err != nil {
		return nil, err
	}
	if summary.RemovedMembers, err = count(nil, &removed); err != nil {
		return nil, err
	}
	if summary.Owners, err = count(&owner, &active); err != nil {
		return nil, err
	}
	if summary.Admins, err = count(&admin, &active); err != nil {
		return nil, err
	}
	if summary.Staff, err = count(&staff, &active); err != nil {
		return nil, err
	}
	if summary.Residents, err = count(&resident, &active); err != nil {
		return nil, err
	}
	return summary, nil
}

func (s *SocietySvc) ListSocietyMembers(ctx context.Context, filter models.ListSocietyMembersFilter) (*models.PaginatedMembersResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	items, err := s.memberRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	total, err := s.memberRepo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]*models.SocietyMemberResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, item.ToResponse())
	}

	return &models.PaginatedMembersResponse{
		Items:  responses,
		Total:  total,
		Limit:  normalizeResponseLimit(filter.Limit),
		Offset: normalizeResponseOffset(filter.Offset),
	}, nil
}

func (s *SocietySvc) ListAllSocietyMember(ctx context.Context, filter models.ListSocietyMembersFilter) (*models.PaginatedMembersResponse, error) {
	return s.ListSocietyMembers(ctx, filter)
}

func (s *SocietySvc) ListMyMemberships(ctx context.Context, userID int64) ([]*models.SocietyMemberResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	items, err := s.memberRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*models.SocietyMemberResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, item.ToResponse())
	}
	return responses, nil
}
