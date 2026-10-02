package notificationsvc

import (
	"context"
	"fmt"
	"testing"

	"go-server/internal/models"
)

func TestSocietyServiceCheckInRecipientsAndIdempotency(t *testing.T) {
	repo := &notificationRepoFake{}
	residents := &flatResidentRepoFake{}
	tokens := &deviceTokenRepoFake{tokensByUser: map[int64][]string{}}
	// More than one repository page, plus duplicate membership and excluded rows.
	for i := int64(1); i <= 105; i++ {
		residents.items = append(residents.items, &models.FlatResident{SocietyID: 10, FlatID: i, UserID: i, Status: models.FlatResidentStatusActive})
		tokens.tokensByUser[i] = []string{fmt.Sprintf("token-%d", i)}
	}
	residents.items = append(residents.items,
		&models.FlatResident{SocietyID: 10, FlatID: 999, UserID: 1, Status: models.FlatResidentStatusActive},
		&models.FlatResident{SocietyID: 10, UserID: 200, Status: models.FlatResidentStatusInactive},
		&models.FlatResident{SocietyID: 99, UserID: 300, Status: models.FlatResidentStatusActive})
	fcm := &fcmClientFake{}
	svc := &notificationService{notifications: repo, residents: residents, deviceTokens: tokens, fcmClient: fcm, enabled: true}
	provider := "Newspaper Delivery"
	entry := &models.VisitorEntry{ID: 42, SocietyID: 10, Source: models.VisitorEntrySourceGuardEntry, Purpose: models.VisitorPurposeService, Status: models.VisitorStatusCheckedIn, ServiceProvider: &provider, Visitor: &models.VisitorSummary{FullName: "Rajesh"}}
	for i := 0; i < 2; i++ {
		if err := svc.SendVisitorCheckIn(context.Background(), entry); err != nil {
			t.Fatal(err)
		}
	}
	if len(repo.created) != 105 || fcm.sendCount != 105 {
		t.Fatalf("notifications=%d pushes=%d, want 105 each", len(repo.created), fcm.sendCount)
	}
	seen := map[int64]bool{}
	for _, item := range repo.created {
		if seen[item.UserID] {
			t.Fatal("duplicate user notification")
		}
		seen[item.UserID] = true
		if item.FlatID != nil || item.Data["flat_id"] != nil || item.Data["scope"] != "society" || item.Data["category_id"] != CategoryNotificationInfo || item.EventKey == nil || *item.EventKey != "visitor_service_checked_in:42" {
			t.Fatalf("incorrect metadata: %+v", item)
		}
		if item.Title != "Newspaper Delivery has arrived" || item.Body != "Rajesh from Newspaper Delivery has entered the society." {
			t.Fatalf("incorrect copy: %+v", item)
		}
	}
}

func TestServiceCheckInTitleUsesProvider(t *testing.T) {
	for _, provider := range []string{"Newspaper wala", "Milkman", "Kachre wala"} {
		t.Run(provider, func(t *testing.T) {
			repo := &notificationRepoFake{}
			svc := &notificationService{notifications: repo, deviceTokens: &deviceTokenRepoFake{}, residents: &flatResidentRepoFake{items: []*models.FlatResident{
				{SocietyID: 10, FlatID: 20, UserID: 1, Status: models.FlatResidentStatusActive},
			}}}
			padded := " " + provider + " "
			entry := &models.VisitorEntry{ID: 42, SocietyID: 10, Source: models.VisitorEntrySourceGuardEntry, Purpose: models.VisitorPurposeService, Status: models.VisitorStatusCheckedIn, ServiceProvider: &padded}
			if err := svc.SendVisitorCheckIn(context.Background(), entry); err != nil {
				t.Fatal(err)
			}
			if len(repo.created) != 1 || repo.created[0].Title != provider+" has arrived" {
				t.Fatalf("incorrect provider title: %+v", repo.created)
			}
		})
	}
}

func TestServiceBroadcastScopeAndState(t *testing.T) {
	for _, tc := range []struct {
		name    string
		source  models.VisitorEntrySource
		purpose models.VisitorPurpose
		flatID  int64
		status  models.VisitorStatus
		want    int
	}{
		{"not yet checked in", models.VisitorEntrySourceGuardEntry, models.VisitorPurposeService, 0, models.VisitorStatusApproved, 0},
		{"historical", models.VisitorEntrySourceGuardEntry, models.VisitorPurposeService, 20, models.VisitorStatusCheckedIn, 1},
		{"public", models.VisitorEntrySourcePublicQR, models.VisitorPurposeService, 20, models.VisitorStatusCheckedIn, 1},
		{"staff", models.VisitorEntrySourceGuardEntry, models.VisitorPurposeStaff, 0, models.VisitorStatusCheckedIn, 0},
		{"fallback", models.VisitorEntrySourceGuardEntry, models.VisitorPurposeService, 0, models.VisitorStatusCheckedIn, 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := &notificationRepoFake{}
			svc := &notificationService{notifications: repo, deviceTokens: &deviceTokenRepoFake{}, residents: &flatResidentRepoFake{items: []*models.FlatResident{
				{SocietyID: 10, FlatID: 20, UserID: 1, Status: models.FlatResidentStatusActive},
				{SocietyID: 10, FlatID: 30, UserID: 2, Status: models.FlatResidentStatusActive},
			}}}
			err := svc.SendVisitorCheckIn(context.Background(), &models.VisitorEntry{ID: 42, SocietyID: 10, Source: tc.source, Purpose: tc.purpose, FlatID: tc.flatID, Status: tc.status})
			if err != nil || len(repo.created) != tc.want {
				t.Fatalf("got %d notifications, err=%v", len(repo.created), err)
			}
			for _, n := range repo.created {
				if tc.name == "fallback" {
					if n.Body != "A service provider has entered the society." {
						t.Fatal(n.Body)
					}
				} else if n.Data["scope"] == "society" {
					t.Fatal("unexpected broadcast")
				}
			}
		})
	}
}
