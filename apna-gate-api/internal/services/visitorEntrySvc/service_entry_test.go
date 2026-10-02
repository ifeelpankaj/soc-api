package visitorentrysvc_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"go-server/internal/models"
	visitorentrysvc "go-server/internal/services/visitorEntrySvc"
)

func TestGuardServiceCreationClearsOnlyGuardFlat(t *testing.T) {
	for _, flatID := range []int64{0, 20, 999} {
		_, svc, _, repo, _, _ := inviteSecurityService()
		req := models.VisitorFormRequest{FullName: "Rajesh", PhoneNumber: strPtr("9876543210"), Purpose: models.VisitorPurposeService, ServiceProvider: strPtr("Newspaper Delivery"), FlatID: flatID}
		result, err := svc.CreateGuardEntry(context.Background(), 10, req, 200)
		if err != nil {
			t.Fatalf("guard Service flat=%d: %v", flatID, err)
		}
		if repo.lastCreateArgs.flatID != nil || repo.lastCreateArgs.req.FlatID != 0 || result.Entry.Purpose != models.VisitorPurposeService || result.Entry.Status != models.VisitorStatusApproved || result.QR == nil {
			t.Fatalf("incorrect guard Service result: %+v / %+v", result, repo.lastCreateArgs)
		}
	}
	for _, source := range []string{"public", "quick"} {
		for _, flatID := range []int64{0, 20} {
			_, svc, _, repo, _, _ := inviteSecurityService()
			req := models.VisitorFormRequest{FullName: "Rajesh", PhoneNumber: strPtr("9876543210"), Purpose: models.VisitorPurposeService, ServiceProvider: strPtr("Newspaper Delivery"), FlatID: flatID}
			var err error
			if source == "public" {
				_, err = svc.CreatePublicQREntry(context.Background(), 10, req)
			} else {
				_, err = svc.CreateQuickLinkEntry(context.Background(), 10, req)
			}
			if flatID == 0 {
				if err == nil || repo.createCalls != 0 {
					t.Fatal("non-guard Service accepted without a flat")
				}
			} else if err != nil || repo.lastCreateArgs.flatID == nil || *repo.lastCreateArgs.flatID != flatID {
				t.Fatalf("non-guard Service lost flat: %v / %+v", err, repo.lastCreateArgs)
			}
		}
	}
}

func TestResidentInviteRejectsServiceAndStaffOnly(t *testing.T) {
	for _, purpose := range []models.VisitorPurpose{models.VisitorPurposeService, models.VisitorPurposeStaff} {
		invites, _, repo, _, _, _ := inviteSecurityService()
		if _, _, _, err := invites.CreateInvite(context.Background(), 10, 20, purpose, 100, nil); err == nil {
			t.Fatalf("resident invite accepted %s", purpose)
		}
		if len(repo.invites) != 0 {
			t.Fatal("rejected invite was persisted")
		}
		if _, _, _, err := invites.CreateStaffInvite(context.Background(), 10, 20, purpose, 200, nil); err != nil {
			t.Fatalf("guard invite rejected %s: %v", purpose, err)
		}
	}
}

func TestSocietyServiceRejectsFlatAssignment(t *testing.T) {
	_, svc, _, repo, _, _ := inviteSecurityService()
	repo.entry = &models.VisitorEntry{ID: 703, SocietyID: 10, Purpose: models.VisitorPurposeService, Source: models.VisitorEntrySourceGuardEntry, Status: models.VisitorStatusApproved}
	flatID := int64(20)
	if _, err := svc.UpdateGuardEntry(context.Background(), 10, 703, 200, models.UpdateGuardVisitorEntryRequest{FlatID: &flatID}); err == nil {
		t.Fatal("society-wide Service accepted flat assignment")
	}
	repo.entry.FlatID = flatID
	if _, err := svc.UpdateGuardEntry(context.Background(), 10, 703, 200, models.UpdateGuardVisitorEntryRequest{FlatID: &flatID}); err != nil {
		t.Fatalf("historical flat-linked Service rejected: %v", err)
	}
}

type serviceCommitFailure struct{}

func (serviceCommitFailure) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	if err := fn(ctx); err != nil {
		return err
	}
	return errors.New("commit failed")
}

func TestServiceCheckInDispatchRequiresCommittedTransition(t *testing.T) {
	for _, mode := range []string{"entry", "qr", "approve"} {
		for _, outcome := range []string{"success", "zero rows", "database failure", "commit failure"} {
			t.Run(mode+"/"+outcome, func(t *testing.T) {
				notifier := &checkInNotifier{entries: make(chan *models.VisitorEntry, 4)}
				_, svc, _, repo, _, _ := inviteSecurityServiceWithNotifier(notifier)
				if outcome == "commit failure" {
					_, svc, _, repo, _, _ = inviteSecurityServiceWithNotifier(notifier, serviceCommitFailure{})
				}
				expiry := time.Now().Add(time.Hour)
				repo.entry = &models.VisitorEntry{ID: 703, SocietyID: 10, Source: models.VisitorEntrySourceGuardEntry, Purpose: models.VisitorPurposeService, Status: models.VisitorStatusApproved, QRExpiresAt: &expiry, ServiceProvider: strPtr("Newspaper Delivery"), Visitor: &models.VisitorSummary{FullName: "Rajesh"}}
				if mode == "approve" {
					repo.entry.Status = models.VisitorStatusWaitingApproval
				}
				repo.checkInNil = outcome == "zero rows"
				if outcome == "database failure" {
					repo.checkInErr = errors.New("database failure")
				}
				checkIn := func() error {
					var err error
					switch mode {
					case "entry":
						_, err = svc.CheckInByEntryID(context.Background(), 10, 703, 200)
					case "qr":
						_, err = svc.CheckIn(context.Background(), "token", 200)
					case "approve":
						_, err = svc.GuardApproveAndCheckIn(context.Background(), 10, 703, 200, visitorentrysvc.GuardApproveOptions{})
					}
					return err
				}
				err := checkIn()
				if outcome == "success" {
					if err != nil {
						t.Fatal(err)
					}
					select {
					case entry := <-notifier.entries:
						if !entry.IsSocietyWideGuardService() || entry.Status != models.VisitorStatusCheckedIn || entry.Visitor == nil || entry.ServiceProvider == nil {
							t.Fatalf("incorrect notification entry: %+v", entry)
						}
					case <-time.After(time.Second):
						t.Fatal("missing notification")
					}
					if checkIn() == nil {
						t.Fatal("repeat check-in succeeded")
					}
				} else if err == nil {
					t.Fatal("failed transition returned success")
				}
				select {
				case <-notifier.entries:
					t.Fatal("notification without a new committed transition")
				case <-time.After(20 * time.Millisecond):
				}
			})
		}
	}
}

func TestConcurrentServiceCheckInHasOneWinner(t *testing.T) {
	for _, mode := range []string{"entry", "qr"} {
		t.Run(mode, func(t *testing.T) {
			notifier := &checkInNotifier{entries: make(chan *models.VisitorEntry, 8)}
			_, svc, _, repo, _, _ := inviteSecurityServiceWithNotifier(notifier)
			expiry := time.Now().Add(time.Hour)
			repo.entry = &models.VisitorEntry{ID: 703, SocietyID: 10, Source: models.VisitorEntrySourceGuardEntry, Purpose: models.VisitorPurposeService, Status: models.VisitorStatusApproved, QRExpiresAt: &expiry, Visitor: &models.VisitorSummary{FullName: "Rajesh"}}
			results := make(chan error, 8)
			var wg sync.WaitGroup
			for i := 0; i < 8; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var err error
					if mode == "qr" {
						_, err = svc.CheckIn(context.Background(), "token", 200)
					} else {
						_, err = svc.CheckInByEntryID(context.Background(), 10, 703, 200)
					}
					results <- err
				}()
			}
			wg.Wait()
			close(results)
			winners := 0
			for err := range results {
				if err == nil {
					winners++
				}
			}
			if winners != 1 {
				t.Fatalf("got %d winning transitions", winners)
			}
			select {
			case <-notifier.entries:
			case <-time.After(time.Second):
				t.Fatal("missing notification")
			}
			select {
			case <-notifier.entries:
				t.Fatal("duplicate broadcast")
			case <-time.After(20 * time.Millisecond):
			}
		})
	}
}
