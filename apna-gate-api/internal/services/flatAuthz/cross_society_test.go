package flatauthz

import (
	"context"
	"errors"
	"testing"

	"go-server/internal/models"
)

type societyScopedFlatRepo struct {
	fakeFlatRepo
	flat *models.Flat
}

func (f *societyScopedFlatRepo) Get(_ context.Context, filter *models.FlatFilter) (*models.Flat, error) {
	if f.flat == nil || filter == nil || filter.ID == nil || filter.SocietyID == nil {
		return nil, nil
	}
	if f.flat.ID != *filter.ID || f.flat.SocietyID != *filter.SocietyID {
		return nil, nil
	}
	return f.flat, nil
}

func TestCanViewFlatVisitorsDeniesCrossSocietyResource(t *testing.T) {
	authz := New(
		&fakeMemberRepo{member: &models.SocietyMember{
			SocietyID: 1,
			UserID:    100,
			Role:      models.SocietyMemberRoleAdmin,
			Status:    models.SocietyMemberStatusActive,
		}},
		&fakeResidentRepo{},
		&societyScopedFlatRepo{flat: &models.Flat{ID: 10, SocietyID: 2}},
	)

	err := authz.CanViewFlatVisitors(context.Background(), 1, 10, 100)
	var appErr *models.AppError
	if !errors.As(err, &appErr) || appErr.Code != ErrFlatNotFound.Code {
		t.Fatalf("expected cross-society flat to be hidden as flat not found, got %v", err)
	}
}
