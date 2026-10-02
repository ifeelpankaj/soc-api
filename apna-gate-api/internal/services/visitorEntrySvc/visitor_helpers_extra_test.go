package visitorentrysvc

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-server/internal/models"
)

type checkoutSettingSvc struct {
	duration int32
	err      error
}

func (s checkoutSettingSvc) GetSocietySettings(context.Context, int64) (*models.SocietyVisitorSettingsResponse, error) {
	return nil, nil
}
func (s checkoutSettingSvc) GetFlatSettings(context.Context, int64, int64) ([]models.FlatVisitorSettingsResponse, error) {
	return nil, nil
}
func (s checkoutSettingSvc) EnsureDefaultFlatSettingsIfMissing(context.Context, int64, int64, int64) error {
	return nil
}
func (s checkoutSettingSvc) ResolveApprovalRequirement(context.Context, int64, int64, models.VisitorPurpose, models.VisitorEntrySource) (bool, error) {
	return false, nil
}
func (s checkoutSettingSvc) ResolveVisitDurationMinutes(context.Context, int64, int64, models.VisitorPurpose) (int32, error) {
	if s.err != nil {
		return 0, s.err
	}
	return s.duration, nil
}

func TestParseVisitorEntryFilterValue(t *testing.T) {
	// if got, err := ParseVisitorEntryFilterValue[models.VisitorStatus]("", models.VisitorStatus.IsValid); err != nil || got != nil {
	// 	t.Fatalf("empty filter should be nil, got %v %v", got, err)
	// }

	got, err := ParseVisitorEntryFilterValue("checked_in", models.VisitorStatus.IsValid)
	if err != nil {
		t.Fatalf("valid status returned error: %v", err)
	}
	if got == nil || *got != models.VisitorStatusCheckedIn {
		t.Fatalf("unexpected status pointer: %v", got)
	}

	if _, err := ParseVisitorEntryFilterValue("bad", models.VisitorStatus.IsValid); !errors.Is(err, ErrInvalidVisitorRequest) {
		t.Fatalf("expected ErrInvalidVisitorRequest, got %v", err)
	}
}

func TestParsePositiveInt64(t *testing.T) {
	tests := []struct {
		raw     string
		wantNil bool
		want    int64
		wantErr bool
	}{
		{raw: "", wantNil: true},
		{raw: "42", want: 42},
		{raw: "0", wantErr: true},
		{raw: "-1", wantErr: true},
		{raw: "abc", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.raw, func(t *testing.T) {
			got, err := ParsePositiveInt64(tt.raw)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidVisitorRequest) {
					t.Fatalf("expected invalid request error, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantNil {
				if got != nil {
					t.Fatalf("expected nil, got %v", *got)
				}
				return
			}
			if got == nil || *got != tt.want {
				t.Fatalf("got %v, want %d", got, tt.want)
			}
		})
	}
}

func TestSmallVisitorHelpers(t *testing.T) {
	block := "A"
	if got := blockKey(&block); got != "A" {
		t.Fatalf("blockKey = %q", got)
	}
	if got := blockKey(nil); got != "" {
		t.Fatalf("nil blockKey = %q", got)
	}

	actor := int64(9)
	if got := guardActor(models.VisitorEntrySourceGuardEntry, &actor); got == nil || *got != actor {
		t.Fatalf("guard actor should pass through, got %v", got)
	}
	if got := guardActor(models.VisitorEntrySourcePublicQR, &actor); got != nil {
		t.Fatalf("public QR should not keep actor, got %v", *got)
	}

	if got := entryFlatIDPtr(12); got == nil || *got != 12 {
		t.Fatalf("entryFlatIDPtr positive mismatch: %v", got)
	}
	if got := entryFlatIDPtr(0); got != nil {
		t.Fatalf("entryFlatIDPtr zero should be nil, got %v", *got)
	}

	if !IsInvalidStateNoRows(ErrVisitorInvalidState) {
		t.Fatal("expected invalid state sentinel to match")
	}
	if IsInvalidStateNoRows(errors.New("other")) {
		t.Fatal("unexpected invalid state match")
	}
}

func TestQRDisplayTokenFromMetadata(t *testing.T) {
	if token, ok := qrDisplayTokenFromMetadata(nil); ok || token != "" {
		t.Fatalf("nil metadata should not return token: %q %v", token, ok)
	}
	if token, ok := qrDisplayTokenFromMetadata(map[string]any{qrDisplayTokenMetadataKey: ""}); ok || token != "" {
		t.Fatalf("empty token should not be accepted: %q %v", token, ok)
	}
	if token, ok := qrDisplayTokenFromMetadata(map[string]any{qrDisplayTokenMetadataKey: 123}); ok || token != "" {
		t.Fatalf("non-string token should not be accepted: %q %v", token, ok)
	}
	if token, ok := qrDisplayTokenFromMetadata(map[string]any{qrDisplayTokenMetadataKey: "visible"}); !ok || token != "visible" {
		t.Fatalf("token mismatch: %q %v", token, ok)
	}
}

func TestResolveExpectedCheckout(t *testing.T) {
	ctx := context.Background()
	expectedAt := time.Now().Add(30 * time.Minute)
	provided := expectedAt.Add(time.Hour)
	svc := &VisitorEntrySvc{settingSvc: checkoutSettingSvc{duration: 45}}

	got, err := svc.resolveExpectedCheckout(ctx, 1, 2, models.VisitorPurposeGuest, &expectedAt, &provided)
	if err != nil {
		t.Fatalf("provided checkout returned error: %v", err)
	}
	if got == nil || !got.Equal(provided) {
		t.Fatalf("provided checkout mismatch: %v", got)
	}

	before := expectedAt.Add(-time.Minute)
	if _, err := svc.resolveExpectedCheckout(ctx, 1, 2, models.VisitorPurposeGuest, &expectedAt, &before); err == nil {
		t.Fatal("expected invalid checkout error")
	} else {
		var appErr *models.AppError
		if !errors.As(err, &appErr) || appErr.Code != ErrInvalidVisitorRequest.Code {
			t.Fatalf("expected invalid visitor request code, got %v", err)
		}
	}

	got, err = svc.resolveExpectedCheckout(ctx, 1, 2, models.VisitorPurposeGuest, &expectedAt, nil)
	if err != nil {
		t.Fatalf("default checkout returned error: %v", err)
	}
	if got == nil || !got.Equal(expectedAt.Add(45*time.Minute)) {
		t.Fatalf("default checkout mismatch: got %v want %v", got, expectedAt.Add(45*time.Minute))
	}

	settingErr := errors.New("settings down")
	svc.settingSvc = checkoutSettingSvc{err: settingErr}
	if _, err := svc.resolveExpectedCheckout(ctx, 1, 2, models.VisitorPurposeGuest, nil, nil); !errors.Is(err, settingErr) {
		t.Fatalf("expected settings error, got %v", err)
	}
}

func TestApplyExpectedCheckout(t *testing.T) {
	svc := &VisitorEntrySvc{settingSvc: checkoutSettingSvc{duration: 30}}
	if err := svc.applyExpectedCheckout(context.Background(), 1, nil); !errors.Is(err, ErrInvalidVisitorRequest) {
		t.Fatalf("expected invalid request for nil form, got %v", err)
	}

	req := &models.VisitorFormRequest{FlatID: 2, Purpose: models.VisitorPurposeGuest}
	if err := svc.applyExpectedCheckout(context.Background(), 1, req); err != nil {
		t.Fatalf("applyExpectedCheckout returned error: %v", err)
	}
	if req.ExpectedCheckoutAt == nil {
		t.Fatal("expected checkout was not set")
	}
}
