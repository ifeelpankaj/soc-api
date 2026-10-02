package visitorsettingsvc

import (
	"context"
	"errors"
	"go-server/internal/repositories/contracts"
	"testing"

	"go-server/internal/models"
)

func TestResolveApprovalRequirement(t *testing.T) {
	ctx := context.Background()
	flatSetting := &models.FlatVisitorSettings{ApprovalRequired: true, IsEnabled: true}

	tests := []struct {
		name     string
		society  *models.SocietyVisitorSettings
		flat     *models.FlatVisitorSettings
		source   models.VisitorEntrySource
		want     bool
		wantCode string
	}{
		{
			name: "resident link bypasses approval when enabled",
			society: &models.SocietyVisitorSettings{
				IsActive: true, ApprovalMode: models.VisitorApprovalModeMandatory, AllowResidentPreApproval: true,
			},
			source: models.VisitorEntrySourceResidentLink,
			want:   false,
		},
		{
			name: "disabled public qr returns source disabled",
			society: &models.SocietyVisitorSettings{
				IsActive: true, ApprovalMode: models.VisitorApprovalModeOptional, AllowPublicQREntry: false,
			},
			source:   models.VisitorEntrySourcePublicQR,
			wantCode: ErrVisitorSourceDisabled.Code,
		},
		{
			name: "mandatory requires approval",
			society: &models.SocietyVisitorSettings{
				IsActive: true, ApprovalMode: models.VisitorApprovalModeMandatory, AllowGuardEntry: true,
			},
			source: models.VisitorEntrySourceGuardEntry,
			want:   true,
		},
		{
			name: "optional skips approval",
			society: &models.SocietyVisitorSettings{
				IsActive: true, ApprovalMode: models.VisitorApprovalModeOptional, AllowGuardEntry: true,
			},
			source: models.VisitorEntrySourceGuardEntry,
			want:   false,
		},
		{
			name: "hybrid uses flat purpose",
			society: &models.SocietyVisitorSettings{
				IsActive: true, ApprovalMode: models.VisitorApprovalModeHybrid, AllowGuardEntry: true,
			},
			flat:   flatSetting,
			source: models.VisitorEntrySourceGuardEntry,
			want:   true,
		},
		{
			name: "hybrid disabled purpose errors",
			society: &models.SocietyVisitorSettings{
				IsActive: true, ApprovalMode: models.VisitorApprovalModeHybrid, AllowGuardEntry: true,
			},
			flat:     &models.FlatVisitorSettings{ApprovalRequired: true, IsEnabled: false},
			source:   models.VisitorEntrySourceGuardEntry,
			wantCode: ErrVisitorPurposeDisabled.Code,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &VisitorSettingSvc{
				settingsRepo: &fakeVisitorSettingRepository{society: tt.society, flat: tt.flat},
			}

			got, err := svc.ResolveApprovalRequirement(ctx, 1, 10, models.VisitorPurposeGuest, tt.source)
			if tt.wantCode != "" {
				var appErr *models.AppError
				if !errors.As(err, &appErr) || appErr.Code != tt.wantCode {
					t.Fatalf("expected app error %s, got %v", tt.wantCode, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("ResolveApprovalRequirement returned error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("approval required = %v, want %v", got, tt.want)
			}
		})
	}
}

type fakeVisitorSettingRepository struct {
	society *models.SocietyVisitorSettings
	flat    *models.FlatVisitorSettings
}

func (r *fakeVisitorSettingRepository) CreateDefaultSociety(context.Context, int64, int64) error {
	return nil
}

func (r *fakeVisitorSettingRepository) GetSociety(context.Context, int64) (*models.SocietyVisitorSettings, error) {
	return r.society, nil
}

func (r *fakeVisitorSettingRepository) UpdateSociety(context.Context, int64, contracts.UpdateSocietyVisitorSettingsInput, int64) (*models.SocietyVisitorSettings, error) {
	return nil, nil
}

func (r *fakeVisitorSettingRepository) CreateDefaultFlat(context.Context, int64, int64, int64) error {
	return nil
}

func (r *fakeVisitorSettingRepository) ListFlat(context.Context, int64, int64) ([]*models.FlatVisitorSettings, error) {
	return nil, nil
}

func (r *fakeVisitorSettingRepository) ListSocietyFlat(context.Context, models.SocietyFlatVisitorSettingsFilter) ([]*models.SocietyFlatVisitorSettingRow, error) {
	return nil, nil
}

func (r *fakeVisitorSettingRepository) CountSocietyFlat(context.Context, models.SocietyFlatVisitorSettingsFilter) (int64, error) {
	return 0, nil
}

func (r *fakeVisitorSettingRepository) GetFlatPurpose(context.Context, int64, int64, models.VisitorPurpose) (*models.FlatVisitorSettings, error) {
	return r.flat, nil
}

func (r *fakeVisitorSettingRepository) UpdateFlatPurpose(context.Context, int64, int64, models.VisitorPurpose, contracts.UpdateFlatVisitorSettingInput, int64) (*models.FlatVisitorSettings, error) {
	return nil, nil
}

func (r *fakeVisitorSettingRepository) DeleteFlat(context.Context, int64, int64) error {
	return nil
}
