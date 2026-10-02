package notificationsvc

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"firebase.google.com/go/v4/messaging"
	"go-server/internal/models"
	"go-server/internal/repositories/contracts"
)

func TestDeviceDeliveryPreservesVisitorActionCategory(t *testing.T) {
	fcm := &fcmClientFake{}
	s := &notificationService{fcmClient: fcm, enabled: true}
	d := &contracts.PushDelivery{Token: "android-token", Platform: models.DevicePlatformAndroid, Provider: "fcm", NotificationID: "inbox-id", Data: map[string]any{"type": "visitor.pending", "category_id": "visitor_decision"}}
	if _, _, _, err := s.sendPushDelivery(context.Background(), d); err != nil {
		t.Fatal(err)
	}
	if len(fcm.messages) != 1 || fcm.messages[0].Data["categoryId"] != "visitor_decision" {
		t.Fatal("visitor action category missing from Android push")
	}
}

func TestSendVisitorInviteAcceptedInsertsAndSendsOnlyForNewRows(t *testing.T) {
	ctx := context.Background()
	notifications := &notificationRepoFake{}
	fcm := &fcmClientFake{}
	service := &notificationService{
		deviceTokens: &deviceTokenRepoFake{tokensByUser: map[int64][]string{
			101: {"token-101", "token-101-second-device"},
			102: {"token-102"},
			103: {"token-103"},
		}},
		notifications: notifications,
		residents: &flatResidentRepoFake{items: []*models.FlatResident{
			{SocietyID: 1, FlatID: 2, UserID: 101, Status: models.FlatResidentStatusActive, IsPrimary: true},
			{SocietyID: 1, FlatID: 2, UserID: 102, Status: models.FlatResidentStatusActive},
			{SocietyID: 1, FlatID: 2, UserID: 103, Status: models.FlatResidentStatusActive},
			{SocietyID: 1, FlatID: 2, UserID: 102, Status: models.FlatResidentStatusActive},
			{SocietyID: 1, FlatID: 2, UserID: 104, Status: models.FlatResidentStatusInactive},
			{SocietyID: 1, FlatID: 3, UserID: 105, Status: models.FlatResidentStatusActive},
		}},
		fcmClient: fcm,
		enabled:   true,
	}

	inviteID := int64(55)
	entry := &models.VisitorEntry{
		ID:        77,
		SocietyID: 1,
		FlatID:    2,
		InviteID:  &inviteID,
		Visitor:   &models.VisitorSummary{FullName: "Rahul"},
		Flat:      &models.VisitorFlatSummary{ID: 2, FlatNumber: "005"},
	}

	if err := service.SendVisitorInviteAccepted(ctx, entry); err != nil {
		t.Fatalf("SendVisitorInviteAccepted returned error: %v", err)
	}
	if notifications.createCount != 3 {
		t.Fatalf("expected three notification inserts, got %d", notifications.createCount)
	}
	if fcm.sendCount != 3 {
		t.Fatalf("expected three FCM sends, got %d", fcm.sendCount)
	}
	assertCreatedUserIDs(t, notifications.created, []int64{101, 102, 103})
	created := notifications.created[0]
	if created.EventKey == nil || *created.EventKey != "visitor_invite.accepted:55:77" {
		t.Fatalf("unexpected event key: %#v", created.EventKey)
	}
	if created.Data["notification_id"] == "" || created.Data["event"] != "visitor_invite_accepted" {
		t.Fatalf("notification data was not normalized: %#v", created.Data)
	}
	if created.Title != "Rahul accepted your invite" || created.Body != "Rahul has completed the visitor details. The visitor pass is ready." {
		t.Fatalf("unexpected invite notification copy: %q / %q", created.Title, created.Body)
	}
	assertLastPushCopy(t, fcm, created.Title, created.Body)
	if len(multicastRegistrationTokens(fcm.messages[0])) != 2 {
		t.Fatal("expected both recipient devices to receive the same multicast payload")
	}
	for i, row := range notifications.created {
		message := fcm.messages[i]
		if row.Data["notification_id"] != row.ID || message.Data["notification_id"] != row.ID {
			t.Fatalf("durable and push IDs differ: row=%s data=%v push=%v", row.ID, row.Data, message.Data)
		}
		if message.Notification.Title != row.Title || message.Notification.Body != row.Body {
			t.Fatal("durable and push copy differ")
		}
	}

	notifications.returnDuplicate = true
	if err := service.SendVisitorInviteAccepted(ctx, entry); err != nil {
		t.Fatalf("duplicate SendVisitorInviteAccepted returned error: %v", err)
	}
	if fcm.sendCount != 3 {
		t.Fatalf("expected duplicate notification to skip FCM, got %d sends", fcm.sendCount)
	}
}

func TestSendMemberInviteAcceptedExcludesJoinedUserNotResidentID(t *testing.T) {
	ctx := context.Background()
	notifications := &notificationRepoFake{}
	fcm := &fcmClientFake{}
	service := &notificationService{
		deviceTokens: &deviceTokenRepoFake{tokensByUser: map[int64][]string{
			101: {"token-101"},
			102: {"token-102"},
			103: {"token-103"},
		}},
		notifications: notifications,
		residents: &flatResidentRepoFake{items: []*models.FlatResident{
			{SocietyID: 1, FlatID: 2, UserID: 101, Status: models.FlatResidentStatusActive, IsPrimary: true},
			{SocietyID: 1, FlatID: 2, UserID: 102, Status: models.FlatResidentStatusActive},
			{ID: 101, SocietyID: 1, FlatID: 2, UserID: 103, Status: models.FlatResidentStatusActive},
			{SocietyID: 1, FlatID: 2, UserID: 102, Status: models.FlatResidentStatusActive},
			{SocietyID: 1, FlatID: 2, UserID: 104, Status: models.FlatResidentStatusMovedOut},
			{SocietyID: 1, FlatID: 3, UserID: 105, Status: models.FlatResidentStatusActive},
		}},
		fcmClient: fcm,
		enabled:   true,
	}

	invite := &models.FlatMemberInvite{
		ID:        56,
		SocietyID: 1,
		FlatID:    2,
		FullName:  "Amit",
		Role:      models.FlatMemberInviteRoleFamily,
	}

	if err := service.SendMemberInviteAccepted(ctx, invite, "", "Amit", 101); err != nil {
		t.Fatalf("SendMemberInviteAccepted returned error: %v", err)
	}
	if notifications.createCount != 2 {
		t.Fatalf("expected two notification inserts, got %d", notifications.createCount)
	}
	if fcm.sendCount != 2 {
		t.Fatalf("expected two FCM sends, got %d", fcm.sendCount)
	}
	assertCreatedUserIDs(t, notifications.created, []int64{101, 102})
	for _, created := range notifications.created {
		if created.EventKey == nil || *created.EventKey != "member_invite.accepted:56:101" {
			t.Fatalf("unexpected event key: %#v", created.EventKey)
		}
		if created.Title != "Amit joined your flat" || created.Body != "Amit accepted your invitation and is now connected to your flat." {
			t.Fatalf("unexpected member notification copy: %q / %q", created.Title, created.Body)
		}
	}
	assertLastPushCopy(t, fcm, "Amit joined your flat", "Amit accepted your invitation and is now connected to your flat.")
}

func TestApprovedDispatchRejectsNonApprovedEntries(t *testing.T) {
	service := &notificationService{}
	for _, status := range []models.VisitorStatus{"", models.VisitorStatusWaitingApproval, models.VisitorStatusCheckedIn, models.VisitorStatusCheckedOut, models.VisitorStatusRejected} {
		// No repositories are configured: these states must return before dispatch.
		if err := service.SendVisitorApproved(context.Background(), &models.VisitorEntry{Status: status}); err != nil {
			t.Fatalf("status %s: %v", status, err)
		}
	}
}

func TestMemberJoinedNameAndFlatFallbacks(t *testing.T) {
	for _, tc := range []struct{ joined, invited, flat, name, location string }{
		{"  Rahul  ", "Amit", " 302 ", "Rahul", "Flat 302"},
		{" ", "  Amit  ", "", "Amit", "your flat"},
		{"", " ", "", "A member", "your flat"},
	} {
		repo := &notificationRepoFake{}
		service := &notificationService{deviceTokens: &deviceTokenRepoFake{}, notifications: repo, residents: &flatResidentRepoFake{items: []*models.FlatResident{
			{ID: 77, SocietyID: 1, FlatID: 2, UserID: 101, Status: models.FlatResidentStatusActive},
			{ID: 78, SocietyID: 1, FlatID: 2, UserID: 102, Status: models.FlatResidentStatusActive},
		}}}
		if err := service.SendMemberInviteAccepted(context.Background(), &models.FlatMemberInvite{ID: 1, SocietyID: 1, FlatID: 2, FullName: tc.invited}, tc.flat, tc.joined, 77); err != nil {
			t.Fatal(err)
		}
		assertCreatedUserIDs(t, repo.created, []int64{102})
		if repo.created[0].Title != tc.name+" joined your flat" || repo.created[0].Body != tc.name+" accepted your invitation and is now connected to "+tc.location+"." {
			t.Fatalf("unexpected copy: %+v", repo.created[0])
		}
	}
}

func TestSendVisitorInviteAcceptedKeepsCreatedRowWhenFCMFails(t *testing.T) {
	ctx := context.Background()
	notifications := &notificationRepoFake{}
	service := &notificationService{
		deviceTokens:  &deviceTokenRepoFake{tokensByUser: map[int64][]string{101: {"token-101"}}},
		notifications: notifications,
		residents: &flatResidentRepoFake{items: []*models.FlatResident{
			{SocietyID: 1, FlatID: 2, UserID: 101, Status: models.FlatResidentStatusActive, IsPrimary: true},
		}},
		fcmClient: &fcmClientFake{err: errors.New("fcm unavailable")},
		enabled:   true,
	}

	inviteID := int64(55)
	err := service.SendVisitorInviteAccepted(ctx, &models.VisitorEntry{
		ID:        77,
		SocietyID: 1,
		FlatID:    2,
		InviteID:  &inviteID,
	})

	if err == nil {
		t.Fatal("expected FCM error")
	}
	if len(notifications.created) != 1 {
		t.Fatalf("expected DB row to remain created, got %d rows", len(notifications.created))
	}
}

func TestSendVisitorCheckInNotifiesOnlyFlatResidents(t *testing.T) {
	ctx := context.Background()
	notifications := &notificationRepoFake{}
	fcm := &fcmClientFake{}
	service := &notificationService{
		deviceTokens: &deviceTokenRepoFake{tokensByUser: map[int64][]string{
			101: {"token-101"},
			102: {"token-102"},
			201: {"token-201"},
			202: {"token-202"},
		}},
		notifications: notifications,
		residents: &flatResidentRepoFake{items: []*models.FlatResident{
			{SocietyID: 1, FlatID: 2, UserID: 101, Status: models.FlatResidentStatusActive, IsPrimary: true},
			{SocietyID: 1, FlatID: 2, UserID: 102, Status: models.FlatResidentStatusActive},
			{SocietyID: 1, FlatID: 2, UserID: 102, Status: models.FlatResidentStatusActive},
			{SocietyID: 1, FlatID: 2, UserID: 103, Status: models.FlatResidentStatusMovedOut},
		}},
		members: &societyMemberRepoFake{items: []*models.SocietyMember{
			{SocietyID: 1, UserID: 201, Role: models.SocietyMemberRoleStaff, Status: models.SocietyMemberStatusActive},
			{SocietyID: 1, UserID: 202, Role: models.SocietyMemberRoleStaff, Status: models.SocietyMemberStatusActive},
		}},
		fcmClient: fcm,
		enabled:   true,
	}

	err := service.SendVisitorCheckIn(ctx, &models.VisitorEntry{
		ID:        77,
		SocietyID: 1,
		FlatID:    2,
		Visitor:   &models.VisitorSummary{FullName: "Rahul"},
		Flat:      &models.VisitorFlatSummary{ID: 2, FlatNumber: "005"},
	})

	if err != nil {
		t.Fatalf("SendVisitorCheckIn returned error: %v", err)
	}
	assertCreatedUserIDs(t, notifications.created, []int64{101, 102})
	if fcm.sendCount != 2 {
		t.Fatalf("expected two resident FCM sends, got %d", fcm.sendCount)
	}
	for _, created := range notifications.created {
		if created.Title != "Rahul has arrived" || created.Body != "Rahul just checked in at the society gate." {
			t.Fatalf("unexpected check-in copy: %q / %q", created.Title, created.Body)
		}
	}
}

func TestGuardWorkflowNotificationRecipients(t *testing.T) {
	tests := []struct {
		name              string
		send              func(context.Context, *notificationService, *models.VisitorEntry) error
		wantStaffTitle    string
		wantStaffBody     string
		wantResidentTitle string
		wantResidentBody  string
		wantUsers         []int64
		decisionUsers     map[int64]bool
	}{
		{
			name: "approval requested notifies residents only and excludes staff with residency",
			send: func(ctx context.Context, service *notificationService, entry *models.VisitorEntry) error {
				return service.SendVisitorApprovalRequested(ctx, entry)
			},
			wantResidentTitle: "Rahul is at the gate",
			wantResidentBody:  "Rahul is waiting for your approval. Approve or decline the request.",
			wantUsers:         []int64{101},
			decisionUsers:     map[int64]bool{101: true},
		},
		{
			name: "approved notifies staff",
			send: func(ctx context.Context, service *notificationService, entry *models.VisitorEntry) error {
				entry.Status = models.VisitorStatusApproved
				return service.SendVisitorApproved(ctx, entry)
			},
			wantStaffTitle: "Visitor approved",
			wantStaffBody:  "Rahul has been approved and can now enter the society.",
			wantUsers:      []int64{201, 202},
		},
		{
			name: "rejected notifies staff only",
			send: func(ctx context.Context, service *notificationService, entry *models.VisitorEntry) error {
				return service.SendVisitorRejected(ctx, entry)
			},
			wantStaffTitle:    "Visitor request declined",
			wantStaffBody:     "Rahul's entry request was declined.",
			wantResidentTitle: "Visitor request declined",
			wantResidentBody:  "Rahul's entry request was declined.",
			wantUsers:         []int64{201, 202},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			notifications := &notificationRepoFake{}
			fcm := &fcmClientFake{}
			service := &notificationService{
				deviceTokens: &deviceTokenRepoFake{tokensByUser: map[int64][]string{
					101: {"token-101"},
					201: {"token-201"},
					202: {"token-202"},
					203: {"token-203"},
					204: {"token-204"},
				}},
				notifications: notifications,
				residents: &flatResidentRepoFake{items: []*models.FlatResident{
					{SocietyID: 1, FlatID: 2, UserID: 101, Status: models.FlatResidentStatusActive, IsPrimary: true},
					{SocietyID: 1, FlatID: 2, UserID: 201, Status: models.FlatResidentStatusActive},
				}},
				members: &societyMemberRepoFake{items: []*models.SocietyMember{
					{SocietyID: 1, UserID: 201, Role: models.SocietyMemberRoleStaff, Status: models.SocietyMemberStatusActive},
					{SocietyID: 1, UserID: 202, Role: models.SocietyMemberRoleStaff, Status: models.SocietyMemberStatusActive},
					{SocietyID: 1, UserID: 202, Role: models.SocietyMemberRoleStaff, Status: models.SocietyMemberStatusActive},
					{SocietyID: 1, UserID: 203, Role: models.SocietyMemberRoleStaff, Status: models.SocietyMemberStatusSuspended},
					{SocietyID: 1, UserID: 204, Role: models.SocietyMemberRoleAdmin, Status: models.SocietyMemberStatusActive},
				}},
				fcmClient: fcm,
				enabled:   true,
			}
			entry := &models.VisitorEntry{
				ID:               77,
				SocietyID:        1,
				FlatID:           2,
				HandledByGuardID: int64Ptr(201),
				CreatedBy:        int64Ptr(202),
				ApproverName:     strPtr("Priya Shah"),
				Visitor:          &models.VisitorSummary{FullName: "Rahul"},
				Flat:             &models.VisitorFlatSummary{ID: 2, FlatNumber: "A-101", Block: strPtr("A")},
			}

			if err := tt.send(ctx, service, entry); err != nil {
				t.Fatalf("send returned error: %v", err)
			}
			assertCreatedUserIDs(t, notifications.created, tt.wantUsers)
			if len(fcm.messages) != len(tt.wantUsers) {
				t.Fatalf("expected %d recipient pushes, got %d", len(tt.wantUsers), len(fcm.messages))
			}
			for _, created := range notifications.created {
				if created.UserID == 201 || created.UserID == 202 {
					if created.Title != tt.wantStaffTitle || created.Body != tt.wantStaffBody {
						t.Fatalf("unexpected staff copy for user %d: %q / %q", created.UserID, created.Title, created.Body)
					}
					if created.Type == EventVisitorApproved || created.Type == EventVisitorRejected {
						if created.Data["approver_name"] != "Priya Shah" || created.Data["flat_number"] != "A-101" || created.Data["block"] != "A" {
							t.Fatalf("staff notification missing decision metadata: %#v", created.Data)
						}
					}
					continue
				}
				if tt.wantResidentTitle != "" && (created.Title != tt.wantResidentTitle || created.Body != tt.wantResidentBody) {
					t.Fatalf("unexpected resident copy: %q / %q", created.Title, created.Body)
				}
			}
			for _, message := range fcm.messages {
				tokens := multicastRegistrationTokens(message)
				if len(tokens) != 1 {
					t.Fatalf("expected one recipient per push: %#v", tokens)
				}
				userID := userIDForTestToken(tokens[0])
				allowed := false
				for _, wanted := range tt.wantUsers {
					allowed = allowed || userID == wanted
				}
				if !allowed {
					t.Fatalf("unexpected push recipient %d", userID)
				}
				wantCategory := CategoryNotificationInfo
				if tt.decisionUsers[userID] {
					wantCategory = CategoryVisitorDecision
				}
				if message.Data["categoryId"] != wantCategory {
					t.Fatalf("user %d category = %q, want %q", userID, message.Data["categoryId"], wantCategory)
				}
			}
		})
	}
}

func userIDForTestToken(token string) int64 {
	switch token {
	case "token-101":
		return 101
	case "token-201":
		return 201
	case "token-202":
		return 202
	default:
		return 0
	}
}

func TestVisitorNotificationCopyUsesTrimmedNameAndFallbacks(t *testing.T) {
	tests := []struct {
		name      string
		entryName *string
		send      func(context.Context, *notificationService, *models.VisitorEntry) error
		wantTitle string
		wantBody  string
		flat      *models.VisitorFlatSummary
		purpose   models.VisitorPurpose
		partner   *string
	}{
		{
			name:      "invite accepted trims visitor name",
			entryName: strPtr("   Rahul   "),
			send: func(ctx context.Context, service *notificationService, entry *models.VisitorEntry) error {
				return service.SendVisitorInviteAccepted(ctx, entry)
			},
			wantTitle: "Rahul accepted your invite",
			wantBody:  "Rahul has completed the visitor details. The visitor pass is ready.",
		},
		{
			name:      "invite accepted falls back on blank visitor name",
			entryName: strPtr("   "),
			send: func(ctx context.Context, service *notificationService, entry *models.VisitorEntry) error {
				return service.SendVisitorInviteAccepted(ctx, entry)
			},
			wantTitle: "Guest accepted your invite",
			wantBody:  "The guest has completed the visitor details. The visitor pass is ready.",
		},
		{
			name:      "check in uses visitor name",
			entryName: strPtr("Rahul"),
			send: func(ctx context.Context, service *notificationService, entry *models.VisitorEntry) error {
				return service.SendVisitorCheckIn(ctx, entry)
			},
			wantTitle: "Rahul has arrived",
			wantBody:  "Rahul just checked in at the society gate.",
		},
		{
			name:      "check in falls back without visitor",
			entryName: nil,
			send: func(ctx context.Context, service *notificationService, entry *models.VisitorEntry) error {
				entry.Visitor = nil
				return service.SendVisitorCheckIn(ctx, entry)
			},
			wantTitle: "Visitor has arrived",
			wantBody:  "The visitor just checked in at the society gate.",
		},
		{
			name:      "approval requested uses visitor name",
			entryName: strPtr("Rahul"),
			send: func(ctx context.Context, service *notificationService, entry *models.VisitorEntry) error {
				return service.SendVisitorApprovalRequested(ctx, entry)
			},
			wantTitle: "Rahul is at the gate",
			wantBody:  "Rahul is waiting for your approval. Approve or decline the request.",
		},
		{
			name:      "approval requested falls back on blank visitor name",
			entryName: strPtr("   "),
			send: func(ctx context.Context, service *notificationService, entry *models.VisitorEntry) error {
				return service.SendVisitorApprovalRequested(ctx, entry)
			},
			wantTitle: "Visitor is at the gate",
			wantBody:  "Visitor is waiting for your approval. Approve or decline the request.",
		},
		{
			name:      "approval requested falls back to delivery partner",
			entryName: strPtr("   "),
			send: func(ctx context.Context, service *notificationService, entry *models.VisitorEntry) error {
				return service.SendVisitorApprovalRequested(ctx, entry)
			},
			wantTitle: "Swiggy is at the gate",
			wantBody:  "Swiggy is waiting for your approval. Approve or decline the request.",
			purpose:   models.VisitorPurposeDelivery,
			partner:   strPtr("  Swiggy  "),
		},
		{
			name:      "visitor name takes precedence over delivery partner",
			entryName: strPtr("  Rahul  "),
			send: func(ctx context.Context, service *notificationService, entry *models.VisitorEntry) error {
				return service.SendVisitorApprovalRequested(ctx, entry)
			},
			wantTitle: "Rahul is at the gate",
			wantBody:  "Rahul is waiting for your approval. Approve or decline the request.",
			purpose:   models.VisitorPurposeDelivery,
			partner:   strPtr("Zomato"),
		},
		{
			name:      "check in trims flat number",
			entryName: strPtr("Rahul"),
			send: func(ctx context.Context, service *notificationService, entry *models.VisitorEntry) error {
				return service.SendVisitorCheckIn(ctx, entry)
			},
			wantTitle: "Rahul has arrived",
			wantBody:  "Rahul just checked in at the society gate.",
			flat:      &models.VisitorFlatSummary{ID: 2, FlatNumber: " 005 ", Block: strPtr("East Side")},
		},
		{
			name:      "approval requested treats whitespace flat as missing",
			entryName: strPtr("Rahul"),
			send: func(ctx context.Context, service *notificationService, entry *models.VisitorEntry) error {
				return service.SendVisitorApprovalRequested(ctx, entry)
			},
			wantTitle: "Rahul is at the gate",
			wantBody:  "Rahul is waiting for your approval. Approve or decline the request.",
			flat:      &models.VisitorFlatSummary{ID: 2, FlatNumber: "   "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			notifications := &notificationRepoFake{}
			fcm := &fcmClientFake{}
			service := &notificationService{
				deviceTokens:  &deviceTokenRepoFake{tokensByUser: map[int64][]string{101: {"token-101"}}},
				notifications: notifications,
				residents: &flatResidentRepoFake{items: []*models.FlatResident{
					{SocietyID: 1, FlatID: 2, UserID: 101, Status: models.FlatResidentStatusActive, IsPrimary: true},
				}},
				fcmClient: fcm,
				enabled:   true,
			}

			inviteID := int64(55)
			entry := &models.VisitorEntry{
				ID:              77,
				SocietyID:       1,
				FlatID:          2,
				InviteID:        &inviteID,
				Flat:            tt.flat,
				Purpose:         tt.purpose,
				DeliveryPartner: tt.partner,
			}
			if tt.entryName != nil {
				entry.Visitor = &models.VisitorSummary{FullName: *tt.entryName}
			}

			if err := tt.send(ctx, service, entry); err != nil {
				t.Fatalf("send returned error: %v", err)
			}
			if len(notifications.created) != 1 {
				t.Fatalf("expected one created notification, got %d", len(notifications.created))
			}
			created := notifications.created[0]
			if created.Title != tt.wantTitle || created.Body != tt.wantBody {
				t.Fatalf("unexpected notification copy: got %q / %q, want %q / %q", created.Title, created.Body, tt.wantTitle, tt.wantBody)
			}
			if tt.partner != nil {
				if created.Data["delivery_partner"] != strings.TrimSpace(*tt.partner) {
					t.Fatalf("notification missing delivery partner: %#v", created.Data)
				}
			}
			if _, ok := created.Data["visitor_photo_url"]; ok {
				t.Fatalf("notification unexpectedly contains visitor photo: %#v", created.Data)
			}
			assertLastPushCopy(t, fcm, tt.wantTitle, tt.wantBody)
			message := fcm.messages[len(fcm.messages)-1]
			if message.Notification != nil && message.Notification.ImageURL != "" {
				t.Fatalf("notification unexpectedly renders an image: %q", message.Notification.ImageURL)
			}
		})
	}
}

func TestVisitorDecisionCopyFallbackOrder(t *testing.T) {
	tests := []struct {
		name     string
		entry    *models.VisitorEntry
		decision string
		want     string
	}{
		{
			name: "visitor flat and approver",
			entry: &models.VisitorEntry{
				Visitor:      &models.VisitorSummary{FullName: "  Pankaj  "},
				Flat:         &models.VisitorFlatSummary{FlatNumber: " 005 ", Block: strPtr("East Side")},
				ApproverName: strPtr("  Mukesh Kholiya  "),
			},
			decision: "approved",
			want:     "Pankaj was approved for Flat 005 by Mukesh Kholiya.",
		},
		{
			name: "visitor and flat",
			entry: &models.VisitorEntry{
				Visitor: &models.VisitorSummary{FullName: "Pankaj"},
				Flat:    &models.VisitorFlatSummary{FlatNumber: "005"},
			},
			decision: "approved",
			want:     "Pankaj was approved for Flat 005.",
		},
		{
			name: "visitor only",
			entry: &models.VisitorEntry{
				Visitor: &models.VisitorSummary{FullName: "Pankaj"},
			},
			decision: "approved",
			want:     "Pankaj was approved.",
		},
		{
			name: "flat only",
			entry: &models.VisitorEntry{
				Flat: &models.VisitorFlatSummary{FlatNumber: "005"},
			},
			decision: "approved",
			want:     "The visitor was approved for Flat 005.",
		},
		{
			name:     "generic",
			entry:    &models.VisitorEntry{},
			decision: "declined",
			want:     "The visitor was declined.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := visitorDecisionBody(tt.entry, tt.decision); got != tt.want {
				t.Fatalf("visitorDecisionBody() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuildMulticastMessageUsesRegistrationTokens(t *testing.T) {
	message := buildMulticastMessage([]string{"fcm-token-1", "fcm-token-2"}, models.NotificationPayload{Title: "Title"})
	registrationTokens := multicastRegistrationTokens(message)

	if len(registrationTokens) != 2 || registrationTokens[0] != "fcm-token-1" || registrationTokens[1] != "fcm-token-2" {
		t.Fatalf("unexpected registration tokens: %#v", registrationTokens)
	}
	if len(message.Fids) != 0 {
		t.Fatalf("expected no firebase installation ids, got %#v", message.Fids)
	}
}

func TestBuildMulticastMessageMapsExpoCategoryForAndroidAndAPNS(t *testing.T) {
	payload := models.NotificationPayload{
		Title:      "Rahul is at the gate",
		Body:       "Rahul is waiting for your approval. Approve or decline the request.",
		Data:       map[string]string{"notification_id": "notification-1"},
		CategoryID: CategoryVisitorDecision,
	}
	message := buildMulticastMessage([]string{"fcm-token"}, payload)

	if message.Data["categoryId"] != CategoryVisitorDecision {
		t.Fatalf("Android categoryId = %q, want %q", message.Data["categoryId"], CategoryVisitorDecision)
	}
	if message.APNS == nil || message.APNS.Payload == nil || message.APNS.Payload.Aps == nil ||
		message.APNS.Payload.Aps.Category != CategoryVisitorDecision {
		t.Fatalf("APNS category was not mapped: %#v", message.APNS)
	}
	if _, exists := payload.Data["categoryId"]; exists {
		t.Fatal("buildMulticastMessage mutated the caller's data map")
	}
}

func TestSendPushDeletesExpoTokensAndSendsAndroidFCM(t *testing.T) {
	ctx := context.Background()
	fcm := &fcmClientFake{}
	deviceTokens := &deviceTokenRepoFake{tokensByUser: map[int64][]string{
		101: {
			expoToken("expo-1"),
			exponentToken("expo-2"),
			"fcm-token-1",
			"   ",
			"fcm-token-1",
		},
	}}
	service := &notificationService{
		deviceTokens: deviceTokens,
		fcmClient:    fcm,
		enabled:      true,
	}

	if err := service.sendPushToUsers(ctx, []int64{101}, models.NotificationPayload{
		Title: "Visitor checked in",
		Body:  "Pankaj checked in at the gate for Flat 005.",
		Data:  map[string]string{"type": EventVisitorCheckIn},
	}); err != nil {
		t.Fatalf("sendPushToUsers returned error: %v", err)
	}

	assertStringSlice(t, deviceTokens.deletedTokens, []string{"ExpoPushToken[expo-1]", "ExponentPushToken[expo-2]"})
	if fcm.sendCount != 1 {
		t.Fatalf("expected one fcm send, got %d", fcm.sendCount)
	}
	assertStringSlice(t, multicastRegistrationTokens(fcm.messages[0]), []string{"fcm-token-1"})
}

func TestSendPushWithFCMDisabledSkipsNativeFCM(t *testing.T) {
	ctx := context.Background()
	fcm := &fcmClientFake{}
	deviceTokens := &deviceTokenRepoFake{tokensByUser: map[int64][]string{
		101: {expoToken("expo-1"), "fcm-token-1"},
	}}
	service := &notificationService{
		deviceTokens: deviceTokens,
		fcmClient:    fcm,
		enabled:      false,
	}

	if err := service.sendPushToUsers(ctx, []int64{101}, models.NotificationPayload{Title: "Title"}); err != nil {
		t.Fatalf("sendPushToUsers returned error: %v", err)
	}
	if fcm.sendCount != 0 {
		t.Fatalf("expected fcm send to be skipped, got %d", fcm.sendCount)
	}
	assertStringSlice(t, deviceTokens.deletedTokens, []string{expoToken("expo-1")})
}

func TestSendPushSkipsIOSTokens(t *testing.T) {
	ctx := context.Background()
	fcm := &fcmClientFake{}
	service := &notificationService{
		deviceTokens: &deviceTokenRepoFake{tokensByUser: map[int64][]string{
			101: {"ios-token"},
		}},
		fcmClient: fcm,
		enabled:   true,
	}

	if err := service.sendPushToUsers(ctx, []int64{101}, models.NotificationPayload{Title: "Title"}); err != nil {
		t.Fatalf("sendPushToUsers returned error: %v", err)
	}
	if fcm.sendCount != 0 {
		t.Fatalf("expected iOS token to be skipped, got %d sends", fcm.sendCount)
	}
}

func TestRegisterDeviceTokenRejectsExpoTokenForAndroid(t *testing.T) {
	service := &notificationService{deviceTokens: &deviceTokenRepoFake{}}

	_, err := service.RegisterDeviceToken(context.Background(), 101, models.RegisterDeviceTokenRequest{
		Token:    expoToken("legacy"),
		Platform: models.DevicePlatformAndroid,
	})
	if !errors.Is(err, ErrInvalidDeviceToken) {
		t.Fatalf("expected invalid device token error, got %v", err)
	}
}

type notificationRepoFake struct {
	returnDuplicate bool
	createCount     int
	created         []*models.Notification
}

func (f *notificationRepoFake) Create(_ context.Context, input models.NotificationCreate) (*models.Notification, error) {
	f.createCount++
	if f.returnDuplicate {
		return nil, nil
	}
	if input.EventKey != nil {
		for _, existing := range f.created {
			if existing.UserID == input.UserID && existing.EventKey != nil && *existing.EventKey == *input.EventKey {
				return nil, nil
			}
		}
	}
	item := &models.Notification{
		ID:        input.ID,
		UserID:    input.UserID,
		SocietyID: input.SocietyID,
		FlatID:    input.FlatID,
		Type:      input.Type,
		Title:     input.Title,
		Body:      input.Body,
		Data:      input.Data,
		EventKey:  input.EventKey,
		CreatedAt: time.Now().UTC(),
	}
	f.created = append(f.created, item)
	return item, nil
}

func (f *notificationRepoFake) ListByUser(context.Context, models.NotificationListFilter) ([]*models.Notification, error) {
	return nil, nil
}

func (f *notificationRepoFake) CountUnreadByUser(context.Context, int64) (int64, error) {
	return 0, nil
}

func (f *notificationRepoFake) MarkRead(context.Context, int64, string) (*models.Notification, error) {
	return nil, nil
}

func (f *notificationRepoFake) MarkAllRead(context.Context, int64) error {
	return nil
}

type deviceTokenRepoFake struct {
	tokensByUser   map[int64][]string
	deletedTokens  []string
	upsertedTokens []*models.DeviceToken
}

func (f *deviceTokenRepoFake) Upsert(_ context.Context, userID int64, token string, platform models.DevicePlatform, deviceID *string) (*models.DeviceToken, error) {
	item := &models.DeviceToken{UserID: userID, Token: token, Platform: platform, DeviceID: deviceID}
	if deviceID != nil && *deviceID != "" {
		filtered := f.upsertedTokens[:0]
		for _, existing := range f.upsertedTokens {
			if existing.UserID == userID && existing.DeviceID != nil && *existing.DeviceID == *deviceID {
				continue
			}
			filtered = append(filtered, existing)
		}
		f.upsertedTokens = filtered
	}
	f.upsertedTokens = append(f.upsertedTokens, item)
	return item, nil
}

func (f *deviceTokenRepoFake) Delete(context.Context, int64, string) error {
	return nil
}

func (f *deviceTokenRepoFake) ListByUserID(_ context.Context, userID int64) ([]*models.DeviceToken, error) {
	tokens := f.tokensByUser[userID]
	items := make([]*models.DeviceToken, 0, len(tokens))
	for _, token := range tokens {
		platform := models.DevicePlatformAndroid
		if token == "ios-token" {
			platform = models.DevicePlatformIOS
		}
		items = append(items, &models.DeviceToken{UserID: userID, Token: token, Platform: platform})
	}
	return items, nil
}

func (f *deviceTokenRepoFake) DeleteByToken(_ context.Context, token string) error {
	f.deletedTokens = append(f.deletedTokens, token)
	return nil
}

type flatResidentRepoFake struct {
	items []*models.FlatResident
}

func (f *flatResidentRepoFake) Add(context.Context, *models.FlatResident) error {
	return nil
}

func (f *flatResidentRepoFake) Get(context.Context, *models.FlatResidentFilter) (*models.FlatResident, error) {
	return nil, nil
}

func (f *flatResidentRepoFake) List(_ context.Context, filter *models.FlatResidentFilter) ([]*models.FlatResident, error) {
	var result []*models.FlatResident
	for _, item := range f.items {
		if filter != nil && filter.SocietyID != nil && item.SocietyID != *filter.SocietyID {
			continue
		}
		if filter != nil && filter.FlatID != nil && item.FlatID != *filter.FlatID {
			continue
		}
		if filter != nil && filter.IsPrimary != nil && item.IsPrimary != *filter.IsPrimary {
			continue
		}
		if filter != nil && filter.Status != nil && string(item.Status) != *filter.Status {
			continue
		}
		result = append(result, item)
	}
	if filter != nil && filter.Limit > 0 {
		if int(filter.Offset) >= len(result) {
			return nil, nil
		}
		result = result[filter.Offset:]
		if len(result) > int(filter.Limit) {
			result = result[:filter.Limit]
		}
	}
	return result, nil
}

func (f *flatResidentRepoFake) Remove(context.Context, *models.FlatResidentFilter) error {
	return nil
}

func (f *flatResidentRepoFake) MoveOut(context.Context, *models.FlatResidentFilter) (*models.FlatResident, error) {
	return nil, nil
}

func (f *flatResidentRepoFake) ClearPrimary(context.Context, int64, int64) error {
	return nil
}

func (f *flatResidentRepoFake) SetPrimary(context.Context, int64, int64, int64) (*models.FlatResident, error) {
	return nil, nil
}

func (f *flatResidentRepoFake) UpdateRole(context.Context, *models.FlatResidentFilter, models.FlatResidentRole) (*models.FlatResident, error) {
	return nil, nil
}

func (f *flatResidentRepoFake) CountActive(context.Context, int64, int64) (int64, error) {
	return 0, nil
}

func (f *flatResidentRepoFake) CountPrimary(context.Context, int64, int64) (int64, error) {
	return 0, nil
}

type societyMemberRepoFake struct {
	items []*models.SocietyMember
}

func (f *societyMemberRepoFake) Add(context.Context, *models.SocietyMember) error {
	return nil
}

func (f *societyMemberRepoFake) Get(context.Context, models.GetSocietyMemberFilter) (*models.SocietyMember, error) {
	return nil, nil
}

func (f *societyMemberRepoFake) List(_ context.Context, filter models.ListSocietyMembersFilter) ([]*models.SocietyMember, error) {
	var result []*models.SocietyMember
	for _, item := range f.items {
		if item.SocietyID != filter.SocietyID {
			continue
		}
		if filter.Role != nil && string(item.Role) != *filter.Role {
			continue
		}
		if filter.Status != nil && string(item.Status) != *filter.Status {
			continue
		}
		if filter.UserID != nil && item.UserID != *filter.UserID {
			continue
		}
		result = append(result, item)
	}
	return result, nil
}

func (f *societyMemberRepoFake) ListByUser(context.Context, int64) ([]*models.SocietyMember, error) {
	return nil, nil
}

func (f *societyMemberRepoFake) ListMySocietiesByUser(context.Context, int64) ([]*models.MySocietyResponse, error) {
	return nil, nil
}

func (f *societyMemberRepoFake) Count(context.Context, models.ListSocietyMembersFilter) (int64, error) {
	return 0, nil
}

func (f *societyMemberRepoFake) ChangeRole(context.Context, int64, int64, models.SocietyMemberRole) (*models.SocietyMember, error) {
	return nil, nil
}

func (f *societyMemberRepoFake) Suspend(context.Context, int64, int64) (*models.SocietyMember, error) {
	return nil, nil
}

func (f *societyMemberRepoFake) Reactivate(context.Context, int64, int64) (*models.SocietyMember, error) {
	return nil, nil
}

func (f *societyMemberRepoFake) Remove(context.Context, int64, int64, int64, string) error {
	return nil
}

func (f *societyMemberRepoFake) CountActiveOwners(context.Context, int64) (int64, error) {
	return 0, nil
}

func (f *societyMemberRepoFake) DemoteActiveOwners(context.Context, int64, int64) error {
	return nil
}

func (f *societyMemberRepoFake) PromoteToOwner(context.Context, int64, int64) (*models.SocietyMember, error) {
	return nil, nil
}

func (f *societyMemberRepoFake) UpsertResident(context.Context, int64, int64, int64) (*models.SocietyMember, error) {
	return nil, nil
}

type fcmClientFake struct {
	err           error
	errorsByToken map[string]error
	sendCount     int
	messages      []*messaging.MulticastMessage
}

func (f *fcmClientFake) SendEachForMulticast(_ context.Context, message *messaging.MulticastMessage) (*messaging.BatchResponse, error) {
	f.sendCount++
	f.messages = append(f.messages, message)
	if f.err != nil {
		return nil, f.err
	}
	tokens := multicastRegistrationTokens(message)
	responses := make([]*messaging.SendResponse, 0, len(tokens))
	for _, token := range tokens {
		if sendErr := f.errorsByToken[token]; sendErr != nil {
			responses = append(responses, &messaging.SendResponse{Error: sendErr})
		} else {
			responses = append(responses, &messaging.SendResponse{Success: true})
		}
	}
	return &messaging.BatchResponse{
		SuccessCount: len(tokens),
		Responses:    responses,
	}, nil
}

func (f *fcmClientFake) Close() error {
	return nil
}

func assertLastPushCopy(t *testing.T, fcm *fcmClientFake, title string, body string) {
	t.Helper()
	if len(fcm.messages) == 0 {
		t.Fatal("expected one FCM message")
	}
	notification := fcm.messages[len(fcm.messages)-1].Notification
	if notification == nil {
		t.Fatal("expected FCM notification payload")
	}
	if notification.Title != title || notification.Body != body {
		t.Fatalf("unexpected FCM copy: got %q / %q, want %q / %q", notification.Title, notification.Body, title, body)
	}
}

func assertStringSlice(t *testing.T, got []string, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("slice length = %d, want %d; got %#v", len(got), len(want), got)
	}
	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("slice[%d] = %q, want %q; got %#v", index, got[index], want[index], got)
		}
	}
}

func expoToken(value string) string {
	return "ExpoPushToken[" + value + "]"
}

func exponentToken(value string) string {
	return "ExponentPushToken[" + value + "]"
}

func multicastRegistrationTokens(message *messaging.MulticastMessage) []string {
	if message == nil {
		return nil
	}
	value := reflect.ValueOf(message).Elem().FieldByName("Tokens")
	if !value.IsValid() || value.IsNil() {
		return nil
	}
	tokens := make([]string, 0, value.Len())
	for index := 0; index < value.Len(); index++ {
		tokens = append(tokens, value.Index(index).String())
	}
	return tokens
}

func assertCreatedUserIDs(t *testing.T, notifications []*models.Notification, want []int64) {
	t.Helper()
	if len(notifications) != len(want) {
		t.Fatalf("created notification count = %d, want %d", len(notifications), len(want))
	}
	got := make(map[int64]int, len(notifications))
	for _, notification := range notifications {
		got[notification.UserID]++
	}
	for _, userID := range want {
		if got[userID] != 1 {
			t.Fatalf("created notifications for user %d = %d, want 1; all counts: %#v", userID, got[userID], got)
		}
		delete(got, userID)
	}
	if len(got) > 0 {
		t.Fatalf("created notifications for unexpected users: %#v", got)
	}
}

func strPtr(value string) *string {
	return &value
}

func int64Ptr(value int64) *int64 {
	return &value
}
