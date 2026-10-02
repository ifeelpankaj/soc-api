package models

import (
	"strings"
	"testing"
	"time"
)

func ptr[T any](v T) *T { return &v }

func TestVisitorFormRequestValidateCleansAndRequiresContact(t *testing.T) {
	req := &VisitorFormRequest{
		FullName:      "  aSHA   rAO  ",
		PhoneNumber:   ptr("  9876543210  "),
		Email:         ptr("  "),
		PhotoURL:      ptr(" "),
		VehicleNumber: ptr(" MH12AB1234 "),
		Notes:         ptr("  call on arrival  "),
		FlatID:        42,
		Purpose:       VisitorPurposeGuest,
	}

	if err := req.Validate(true); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	if got := req.FullName; got != "Asha Rao" {
		t.Fatalf("full name was not normalized: %q", got)
	}
	if got := *req.PhoneNumber; got != "9876543210" {
		t.Fatalf("phone was not cleaned: %q", got)
	}
	if req.Email != nil {
		t.Fatalf("blank email should be nil, got %q", *req.Email)
	}
	if got := *req.Notes; got != "call on arrival" {
		t.Fatalf("notes were not cleaned: %q", got)
	}

	noContact := &VisitorFormRequest{FullName: "Asha"}
	if err := noContact.Validate(false); err == nil || !strings.Contains(err.Error(), "phone_number or email") {
		t.Fatalf("expected contact validation error, got %v", err)
	}
}

func TestVisitorFormRequestValidateRequiresTenDigitPhone(t *testing.T) {
	tests := []string{
		"987654321",
		"98765432101",
		"+919876543210",
		"98765 43210",
		"987654321a",
	}
	for _, phone := range tests {
		req := &VisitorFormRequest{FullName: "Asha Rao", PhoneNumber: &phone}
		if err := req.Validate(false); err == nil || !strings.Contains(err.Error(), "exactly 10 digits") {
			t.Errorf("expected phone %q to fail validation, got %v", phone, err)
		}
	}
}

func TestVisitorFormRequestValidateFlatPurposeAndTimingRules(t *testing.T) {
	phone := "9876543210"
	req := &VisitorFormRequest{FullName: "Asha", PhoneNumber: &phone, Purpose: VisitorPurposeGuest}
	if err := req.Validate(true); err == nil || !strings.Contains(err.Error(), "flat_id") {
		t.Fatalf("expected flat_id validation error, got %v", err)
	}

	req = &VisitorFormRequest{FullName: "Asha", PhoneNumber: &phone, FlatID: 1, Purpose: VisitorPurpose("bad")}
	if err := req.Validate(true); err == nil || !strings.Contains(err.Error(), "invalid visitor purpose") {
		t.Fatalf("expected purpose validation error, got %v", err)
	}

	start := time.Now().Add(time.Hour)
	beforeStart := start.Add(-time.Minute)
	req = &VisitorFormRequest{
		FullName:           "Asha",
		PhoneNumber:        &phone,
		FlatID:             1,
		Purpose:            VisitorPurposeGuest,
		ExpectedAt:         &start,
		ExpectedCheckoutAt: &beforeStart,
	}
	if err := req.Validate(true); err == nil || !strings.Contains(err.Error(), "expected_checkout_at") {
		t.Fatalf("expected checkout timing error, got %v", err)
	}
}

func TestVisitorFormRequestValidateForPurpose(t *testing.T) {
	phone := "9876543210"
	vehicle := VisitorVehicleTypeCab
	partner := "Porter"

	tests := []struct {
		name    string
		req     VisitorFormRequest
		wantErr string
		assert  func(*testing.T, *VisitorFormRequest)
	}{
		{
			name:    "guest needs phone",
			req:     VisitorFormRequest{Purpose: VisitorPurposeGuest},
			wantErr: "phone_number",
		},
		{
			name: "delivery needs partner and clears companions",
			req: VisitorFormRequest{
				Purpose:          VisitorPurposeDelivery,
				PhoneNumber:      &phone,
				DeliveryPartner:  ptr("  Porter  "),
				CompanionsCount:  2,
				CompanionDetails: []map[string]any{{"name": "helper"}},
			},
			assert: func(t *testing.T, req *VisitorFormRequest) {
				if req.CompanionsCount != 0 || *req.DeliveryPartner != "Porter" {
					t.Fatalf("delivery fields were not normalized: %#v", req)
				}
			},
		},
		{
			name:    "cab needs vehicle",
			req:     VisitorFormRequest{Purpose: VisitorPurposeCab},
			wantErr: "vehicle_number",
		},
		{
			name: "cab accepts vehicle and preserves companions",
			req:  VisitorFormRequest{Purpose: VisitorPurposeCab, VehicleNumber: ptr("MH12AB1234"), VehicleType: &vehicle, CompanionsCount: 1, CompanionDetails: []map[string]any{{"name": "Passenger"}}},
			assert: func(t *testing.T, req *VisitorFormRequest) {
				if req.CompanionsCount != 1 || len(req.CompanionDetails) != 1 || req.CompanionDetails[0]["name"] != "Passenger" {
					t.Fatalf("cab companions should be preserved, got %#v", req)
				}
			},
		},
		{
			name: "service needs provider",
			req:  VisitorFormRequest{Purpose: VisitorPurposeService, PhoneNumber: &phone, ServiceProvider: &partner},
		},
		{
			name: "staff clears flat",
			req:  VisitorFormRequest{Purpose: VisitorPurposeStaff, PhoneNumber: &phone, FlatID: 99},
			assert: func(t *testing.T, req *VisitorFormRequest) {
				if req.FlatID != 0 {
					t.Fatalf("staff entries should clear flat id, got %d", req.FlatID)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.req
			err := req.ValidateForPurpose()
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("ValidateForPurpose returned error: %v", err)
			}
			if tt.assert != nil {
				tt.assert(t, &req)
			}
		})
	}
}

func TestInviteAndQRRequestsValidate(t *testing.T) {
	future := time.Now().Add(time.Hour)
	req := &CreateVisitorInviteRequest{Purpose: VisitorPurposeGuest, ExpiresAt: &future}
	if err := req.Validate(); err != nil {
		t.Fatalf("valid invite request failed: %v", err)
	}

	past := time.Now().Add(-time.Minute)
	req.ExpiresAt = &past
	if err := req.Validate(); err == nil || !strings.Contains(err.Error(), "future") {
		t.Fatalf("expected past expiry error, got %v", err)
	}

	tooFar := time.Now().Add(8 * 24 * time.Hour)
	req.ExpiresAt = &tooFar
	if err := req.Validate(); err == nil || !strings.Contains(err.Error(), "7 days") {
		t.Fatalf("expected max duration error, got %v", err)
	}

	qr := &QRTokenRequest{Token: "  abc123  "}
	if err := qr.Validate(); err != nil {
		t.Fatalf("valid QR token failed: %v", err)
	}
	if qr.Token != "abc123" {
		t.Fatalf("QR token was not trimmed: %q", qr.Token)
	}
	if err := (&QRTokenRequest{Token: "  "}).Validate(); err == nil {
		t.Fatal("expected blank QR token to fail")
	}
}

func TestUpdateAndRejectVisitorRequestsValidate(t *testing.T) {
	if err := (&UpdateGuardVisitorEntryRequest{}).Validate(); err == nil {
		t.Fatal("expected empty update request to fail")
	}
	badFlatID := int64(0)
	if err := (&UpdateGuardVisitorEntryRequest{FlatID: &badFlatID}).Validate(); err == nil || !strings.Contains(err.Error(), "flat_id") {
		t.Fatalf("expected flat_id error, got %v", err)
	}
	badName := "  "
	if err := (&UpdateGuardVisitorEntryRequest{FullName: &badName}).Validate(); err == nil || !strings.Contains(err.Error(), "full_name") {
		t.Fatalf("expected full_name error, got %v", err)
	}
	goodName := "  aSHA   rAO "
	update := &UpdateGuardVisitorEntryRequest{FullName: &goodName}
	if err := update.Validate(); err != nil {
		t.Fatalf("valid update failed: %v", err)
	}
	if got := *update.FullName; got != "Asha Rao" {
		t.Fatalf("updated full name was not normalized: %q", got)
	}
	badPhone := "+919876543210"
	if err := (&UpdateGuardVisitorEntryRequest{PhoneNumber: &badPhone}).Validate(); err == nil || !strings.Contains(err.Error(), "exactly 10 digits") {
		t.Fatalf("expected invalid update phone error, got %v", err)
	}
	goodPhone := " 9876543210 "
	phoneUpdate := &UpdateGuardVisitorEntryRequest{PhoneNumber: &goodPhone}
	if err := phoneUpdate.Validate(); err != nil {
		t.Fatalf("valid phone update failed: %v", err)
	}
	if got := *phoneUpdate.PhoneNumber; got != "9876543210" {
		t.Fatalf("updated phone was not trimmed: %q", got)
	}

	reject := &RejectVisitorEntryRequest{Reason: "  duplicate visitor  "}
	if err := reject.Validate(); err != nil {
		t.Fatalf("valid rejection failed: %v", err)
	}
	if reject.Reason != "duplicate visitor" {
		t.Fatalf("reason was not trimmed: %q", reject.Reason)
	}
	if err := (&RejectVisitorEntryRequest{}).Validate(); err == nil {
		t.Fatal("expected empty rejection reason to fail")
	}
}

func TestVisitorEntryEnumsAndScanPreview(t *testing.T) {
	if !VisitorStatusCheckedIn.IsValid() || VisitorStatus("lost").IsValid() {
		t.Fatal("visitor status validation mismatch")
	}
	if !VisitorVehicleTypeBike.IsValid() || VisitorVehicleType("cycle").IsValid() {
		t.Fatal("vehicle type validation mismatch")
	}
	if !VisitorEntryEventExpected.IsValid() || VisitorEntryEventFilter("bad").IsValid() {
		t.Fatal("entry event filter validation mismatch")
	}

	block := "A"
	floor := "4"
	expires := time.Now().Add(time.Hour)
	entry := &VisitorEntry{
		ID:          10,
		SocietyID:   20,
		FlatID:      30,
		Status:      VisitorStatusApproved,
		Purpose:     VisitorPurposeGuest,
		QRExpiresAt: &expires,
		Visitor:     &VisitorSummary{FullName: "Asha", PhoneNumber: ptr("+919876543210")},
		Flat:        &VisitorFlatSummary{ID: 30, FlatNumber: "401", Block: &block, Floor: &floor},
	}
	preview := entry.ToScanPreview()
	if preview == entry || preview.Visitor.PhoneNumber != nil || preview.Flat.FlatNumber != "401" {
		t.Fatalf("unexpected scan preview: %#v", preview)
	}
	if (*VisitorEntry)(nil).ToScanPreview() != nil {
		t.Fatal("nil entry should return nil preview")
	}
}
