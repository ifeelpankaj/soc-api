package visitorentrysvc_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-server/internal/models"
	repository "go-server/internal/repositories"
	visitorentrysvc "go-server/internal/services/visitorEntrySvc"
)

type qrStateEntryRepo struct {
	guardDeskEntryRepo
	entry *models.VisitorEntry
}

func (r *qrStateEntryRepo) GetByQRHash(context.Context, string) (*models.VisitorEntry, error) {
	return r.entry, nil
}

func TestValidateQRCriticalStates(t *testing.T) {
	validExpiry := time.Now().Add(time.Hour)
	expiredAt := time.Now().Add(-time.Minute)

	tests := []struct {
		name    string
		entry   *models.VisitorEntry
		wantErr error
	}{
		{
			name: "approved unexpired QR is valid",
			entry: &models.VisitorEntry{
				ID:          1,
				SocietyID:   10,
				FlatID:      20,
				Status:      models.VisitorStatusApproved,
				QRExpiresAt: &validExpiry,
			},
		},
		{
			name:    "unknown QR is invalid",
			entry:   nil,
			wantErr: visitorentrysvc.ErrVisitorQRInvalid,
		},
		{
			name: "checked in QR cannot be reused for validation",
			entry: &models.VisitorEntry{
				ID:          1,
				SocietyID:   10,
				FlatID:      20,
				Status:      models.VisitorStatusCheckedIn,
				QRExpiresAt: &validExpiry,
			},
			wantErr: visitorentrysvc.ErrVisitorAlreadyCheckedIn,
		},
		{
			name: "waiting approval QR is invalid state",
			entry: &models.VisitorEntry{
				ID:          1,
				SocietyID:   10,
				FlatID:      20,
				Status:      models.VisitorStatusWaitingApproval,
				QRExpiresAt: &validExpiry,
			},
			wantErr: visitorentrysvc.ErrVisitorInvalidState,
		},
		{
			name: "approved QR without expiry is unavailable",
			entry: &models.VisitorEntry{
				ID:        1,
				SocietyID: 10,
				FlatID:    20,
				Status:    models.VisitorStatusApproved,
			},
			wantErr: visitorentrysvc.ErrVisitorQRUnavailable,
		},
		{
			name: "approved expired QR is expired",
			entry: &models.VisitorEntry{
				ID:          1,
				SocietyID:   10,
				FlatID:      20,
				Status:      models.VisitorStatusApproved,
				QRExpiresAt: &expiredAt,
			},
			wantErr: visitorentrysvc.ErrVisitorQRExpired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, entrySvc := visitorentrysvc.NewVisitorService(
				nil,
				&inviteQueryInviteRepo{},
				&qrStateEntryRepo{entry: tt.entry},
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				noopTxManager{},
			)

			got, err := entrySvc.ValidateQR(context.Background(), "raw-token")
			if tt.wantErr != nil {
				var appErr *models.AppError
				var wantAppErr *models.AppError
				if !errors.As(err, &appErr) || !errors.As(tt.wantErr, &wantAppErr) || appErr.Code != wantAppErr.Code {
					t.Fatalf("expected %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("ValidateQR returned error: %v", err)
			}
			if got == nil || got.ID != tt.entry.ID {
				t.Fatalf("unexpected entry: %+v", got)
			}
		})
	}
}

var _ repository.VisitorEntryRepository = (*qrStateEntryRepo)(nil)
