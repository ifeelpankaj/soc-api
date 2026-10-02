package visitorentrysvc_test

import (
	"context"
	"errors"
	"go-server/internal/repositories/contracts"
	"sync"
	"testing"
	"time"

	"go-server/internal/models"
	repository "go-server/internal/repositories"
	flatauthz "go-server/internal/services/flatAuthz"
	notificationsvc "go-server/internal/services/notificationSvc"
	visitorentrysvc "go-server/internal/services/visitorEntrySvc"
)

type inviteSecurityMemberRepo struct {
	members map[[2]int64]*models.SocietyMember
}

func (r *inviteSecurityMemberRepo) Add(context.Context, *models.SocietyMember) error { return nil }
func (r *inviteSecurityMemberRepo) Get(_ context.Context, filter models.GetSocietyMemberFilter) (*models.SocietyMember, error) {
	for _, member := range r.members {
		if filter.SocietyID != nil && member.SocietyID != *filter.SocietyID {
			continue
		}
		if filter.UserID != nil && member.UserID != *filter.UserID {
			continue
		}
		if filter.Status != nil && string(member.Status) != *filter.Status {
			continue
		}
		if filter.Role != nil && string(member.Role) != *filter.Role {
			continue
		}
		return member, nil
	}
	return nil, nil
}
func (r *inviteSecurityMemberRepo) List(context.Context, models.ListSocietyMembersFilter) ([]*models.SocietyMember, error) {
	return nil, nil
}
func (r *inviteSecurityMemberRepo) ListByUser(context.Context, int64) ([]*models.SocietyMember, error) {
	return nil, nil
}
func (r *inviteSecurityMemberRepo) ListMySocietiesByUser(context.Context, int64) ([]*models.MySocietyResponse, error) {
	return nil, nil
}
func (r *inviteSecurityMemberRepo) Count(context.Context, models.ListSocietyMembersFilter) (int64, error) {
	return 0, nil
}
func (r *inviteSecurityMemberRepo) ChangeRole(context.Context, int64, int64, models.SocietyMemberRole) (*models.SocietyMember, error) {
	return nil, nil
}
func (r *inviteSecurityMemberRepo) Suspend(context.Context, int64, int64) (*models.SocietyMember, error) {
	return nil, nil
}
func (r *inviteSecurityMemberRepo) Reactivate(context.Context, int64, int64) (*models.SocietyMember, error) {
	return nil, nil
}
func (r *inviteSecurityMemberRepo) Remove(context.Context, int64, int64, int64, string) error {
	return nil
}
func (r *inviteSecurityMemberRepo) CountActiveOwners(context.Context, int64) (int64, error) {
	return 0, nil
}
func (r *inviteSecurityMemberRepo) DemoteActiveOwners(context.Context, int64, int64) error {
	return nil
}
func (r *inviteSecurityMemberRepo) PromoteToOwner(context.Context, int64, int64) (*models.SocietyMember, error) {
	return nil, nil
}
func (r *inviteSecurityMemberRepo) UpsertResident(context.Context, int64, int64, int64) (*models.SocietyMember, error) {
	return nil, nil
}

type inviteSecurityResidentRepo struct {
	residents map[[3]int64]*models.FlatResident
}

func (r *inviteSecurityResidentRepo) Add(context.Context, *models.FlatResident) error { return nil }
func (r *inviteSecurityResidentRepo) Get(_ context.Context, filter *models.FlatResidentFilter) (*models.FlatResident, error) {
	for _, resident := range r.residents {
		if filter.SocietyID != nil && resident.SocietyID != *filter.SocietyID {
			continue
		}
		if filter.FlatID != nil && resident.FlatID != *filter.FlatID {
			continue
		}
		if filter.UserID != nil && resident.UserID != *filter.UserID {
			continue
		}
		if filter.Status != nil && string(resident.Status) != *filter.Status {
			continue
		}
		return resident, nil
	}
	return nil, nil
}
func (r *inviteSecurityResidentRepo) List(context.Context, *models.FlatResidentFilter) ([]*models.FlatResident, error) {
	return nil, nil
}
func (r *inviteSecurityResidentRepo) Remove(context.Context, *models.FlatResidentFilter) error {
	return nil
}
func (r *inviteSecurityResidentRepo) MoveOut(context.Context, *models.FlatResidentFilter) (*models.FlatResident, error) {
	return nil, nil
}
func (r *inviteSecurityResidentRepo) ClearPrimary(context.Context, int64, int64) error { return nil }
func (r *inviteSecurityResidentRepo) SetPrimary(context.Context, int64, int64, int64) (*models.FlatResident, error) {
	return nil, nil
}
func (r *inviteSecurityResidentRepo) UpdateRole(context.Context, *models.FlatResidentFilter, models.FlatResidentRole) (*models.FlatResident, error) {
	return nil, nil
}
func (r *inviteSecurityResidentRepo) CountActive(context.Context, int64, int64) (int64, error) {
	return 0, nil
}
func (r *inviteSecurityResidentRepo) CountPrimary(context.Context, int64, int64) (int64, error) {
	return 0, nil
}

type inviteSecurityFlatRepo struct {
	flats map[[2]int64]*models.Flat
}

func (r *inviteSecurityFlatRepo) Create(context.Context, *models.Flat) error { return nil }
func (r *inviteSecurityFlatRepo) Get(_ context.Context, filter *models.FlatFilter) (*models.Flat, error) {
	for _, flat := range r.flats {
		if filter.ID != nil && flat.ID != *filter.ID {
			continue
		}
		if filter.SocietyID != nil && flat.SocietyID != *filter.SocietyID {
			continue
		}
		if filter.Status != nil && string(flat.Status) != *filter.Status {
			continue
		}
		if filter.IsActive != nil && flat.IsActive != *filter.IsActive {
			continue
		}
		return flat, nil
	}
	return nil, nil
}
func (r *inviteSecurityFlatRepo) List(context.Context, *models.FlatFilter) ([]*models.Flat, error) {
	return nil, nil
}
func (r *inviteSecurityFlatRepo) Count(context.Context, *models.FlatFilter) (int64, error) {
	return 0, nil
}
func (r *inviteSecurityFlatRepo) Stats(context.Context, int64) (*models.FlatStatsResponse, error) {
	return nil, nil
}
func (r *inviteSecurityFlatRepo) Update(context.Context, *models.FlatFilter, *contracts.UpdateFlatInput) (*models.Flat, error) {
	return nil, nil
}
func (r *inviteSecurityFlatRepo) Deactivate(context.Context, *models.FlatFilter) error { return nil }
func (r *inviteSecurityFlatRepo) Block(context.Context, *models.FlatFilter) (*models.Flat, error) {
	return nil, nil
}
func (r *inviteSecurityFlatRepo) Unblock(context.Context, *models.FlatFilter) (*models.Flat, error) {
	return nil, nil
}
func (r *inviteSecurityFlatRepo) MarkOccupied(context.Context, int64, int64) (*models.Flat, error) {
	return nil, nil
}
func (r *inviteSecurityFlatRepo) MarkVacant(context.Context, int64, int64) (*models.Flat, error) {
	return nil, nil
}

type inviteSecuritySocietyRepo struct {
	societies map[int64]*models.Society
}

func (r *inviteSecuritySocietyRepo) Create(context.Context, *models.Society) error { return nil }
func (r *inviteSecuritySocietyRepo) Get(_ context.Context, filter models.GetSocietyFilter) (*models.Society, error) {
	for _, society := range r.societies {
		if filter.ID != nil && society.ID != *filter.ID {
			continue
		}
		if filter.Status != nil && string(society.Status) != *filter.Status {
			continue
		}
		return society, nil
	}
	return nil, nil
}
func (r *inviteSecuritySocietyRepo) List(context.Context, models.ListSocietiesFilter) ([]*models.Society, error) {
	return nil, nil
}
func (r *inviteSecuritySocietyRepo) Count(context.Context, models.ListSocietiesFilter) (int64, error) {
	return 0, nil
}
func (r *inviteSecuritySocietyRepo) Update(context.Context, int64, contracts.UpdateSocietyInput) (*models.Society, error) {
	return nil, nil
}
func (r *inviteSecuritySocietyRepo) Approve(context.Context, int64, int64) (*models.Society, error) {
	return nil, nil
}
func (r *inviteSecuritySocietyRepo) Reject(context.Context, int64, int64, string) (*models.Society, error) {
	return nil, nil
}
func (r *inviteSecuritySocietyRepo) Suspend(context.Context, int64, int64, string) (*models.Society, error) {
	return nil, nil
}
func (r *inviteSecuritySocietyRepo) Reactivate(context.Context, int64, int64) (*models.Society, error) {
	return nil, nil
}
func (r *inviteSecuritySocietyRepo) Restore(context.Context, int64) (*models.Society, error) {
	return nil, nil
}
func (r *inviteSecuritySocietyRepo) SoftDelete(context.Context, int64) error { return nil }
func (r *inviteSecuritySocietyRepo) CountPendingByCreator(context.Context, int64) (int64, error) {
	return 0, nil
}

type inviteSecurityInviteRepo struct {
	invites       map[int64]*models.VisitorInvite
	byHash        map[string]int64
	nextID        int64
	lastTokenHash string
}

func (r *inviteSecurityInviteRepo) Create(_ context.Context, societyID int64, flatID int64, purpose models.VisitorPurpose, tokenHash string, expiresAt time.Time, actorUserID int64) (*models.VisitorInvite, error) {
	r.nextID++
	invite := &models.VisitorInvite{ID: r.nextID, SocietyID: societyID, FlatID: flatID, Purpose: purpose, CreatedBy: actorUserID, Status: models.VisitorInviteStatusActive, ExpiresAt: expiresAt}
	if r.invites == nil {
		r.invites = map[int64]*models.VisitorInvite{}
	}
	if r.byHash == nil {
		r.byHash = map[string]int64{}
	}
	r.invites[invite.ID] = invite
	r.byHash[tokenHash] = invite.ID
	r.lastTokenHash = tokenHash
	return invite, nil
}
func (r *inviteSecurityInviteRepo) GetByID(_ context.Context, societyID int64, inviteID int64) (*models.VisitorInvite, error) {
	invite := r.invites[inviteID]
	if invite == nil || invite.SocietyID != societyID {
		return nil, nil
	}
	return invite, nil
}
func (r *inviteSecurityInviteRepo) GetHistoryByID(context.Context, int64, int64, int64, func(string, *time.Time) *models.ShortLinkResponse) (*models.VisitorInviteHistoryItem, error) {
	return nil, nil
}
func (r *inviteSecurityInviteRepo) GetByTokenHash(_ context.Context, tokenHash string) (*models.VisitorInvite, error) {
	id := r.byHash[tokenHash]
	if id == 0 {
		return nil, nil
	}
	return r.invites[id], nil
}
func (r *inviteSecurityInviteRepo) ListHistory(context.Context, int64, int64, models.VisitorInviteHistoryFilter, func(string, *time.Time) *models.ShortLinkResponse) ([]*models.VisitorInviteHistoryItem, error) {
	return nil, nil
}
func (r *inviteSecurityInviteRepo) CountHistory(context.Context, int64, int64, models.VisitorInviteHistoryFilter) (int64, error) {
	return 0, nil
}
func (r *inviteSecurityInviteRepo) MarkUsed(_ context.Context, inviteID int64) (*models.VisitorInvite, error) {
	invite := r.invites[inviteID]
	if invite == nil || invite.Status == models.VisitorInviteStatusUsed {
		return nil, nil
	}
	now := time.Now()
	invite.Status = models.VisitorInviteStatusUsed
	invite.UsedAt = &now
	return invite, nil
}
func (r *inviteSecurityInviteRepo) GetForUpdate(_ context.Context, inviteID int64) (*models.VisitorInvite, error) {
	return r.invites[inviteID], nil
}
func (r *inviteSecurityInviteRepo) Cancel(context.Context, int64, int64) (*models.VisitorInvite, error) {
	return nil, nil
}
func (r *inviteSecurityInviteRepo) ExpireOld(context.Context) error { return nil }

type inviteSecurityVisitorRepo struct {
	created int
	nextID  int64
}

func (r *inviteSecurityVisitorRepo) Create(_ context.Context, req contracts.VisitorFormInput) (*models.Visitor, error) {
	r.created++
	r.nextID++
	return &models.Visitor{ID: r.nextID, FullName: req.FullName, PhoneNumber: req.PhoneNumber, Email: req.Email}, nil
}
func (r *inviteSecurityVisitorRepo) Get(context.Context, int64) (*models.Visitor, error) {
	return nil, nil
}
func (r *inviteSecurityVisitorRepo) UpdateProfile(context.Context, int64, contracts.UpdateGuardVisitorEntryInput) (*models.Visitor, error) {
	return nil, nil
}

type inviteSecurityEntryRepo struct {
	*guardDeskEntryRepo
	mu             sync.Mutex
	checkInNil     bool
	checkInErr     error
	getErr         error
	entry          *models.VisitorEntry
	approveNil     bool
	createCalls    int
	approvedCalls  int
	rejectedCalls  int
	rejectNil      bool
	lockCalls      int
	lastCreateArgs struct {
		req         models.VisitorFormRequest
		societyID   int64
		flatID      *int64
		inviteID    *int64
		source      models.VisitorEntrySource
		purpose     models.VisitorPurpose
		status      models.VisitorStatus
		actorUserID *int64
	}
}

func (r *inviteSecurityEntryRepo) Create(_ context.Context, req contracts.VisitorFormInput, societyID int64, flatID *int64, visitorID int64, inviteID *int64, source models.VisitorEntrySource, purpose models.VisitorPurpose, status models.VisitorStatus, actorUserID *int64, guardUserID *int64, qrHash *string, qrExpiresAt *time.Time) (*models.VisitorEntry, error) {
	r.createCalls++
	r.lastCreateArgs.req = models.VisitorFormRequest(req)
	r.lastCreateArgs.societyID = societyID
	r.lastCreateArgs.flatID = flatID
	r.lastCreateArgs.inviteID = inviteID
	r.lastCreateArgs.source = source
	r.lastCreateArgs.purpose = purpose
	r.lastCreateArgs.status = status
	r.lastCreateArgs.actorUserID = actorUserID
	entry := &models.VisitorEntry{ID: 900 + int64(r.createCalls), SocietyID: societyID, VisitorID: visitorID, Source: source, Purpose: purpose, Status: status, InviteID: inviteID, CreatedBy: actorUserID, QRExpiresAt: qrExpiresAt, Metadata: map[string]any{}}
	if flatID != nil {
		entry.FlatID = *flatID
	}
	r.entry = entry
	return entry, nil
}
func (r *inviteSecurityEntryRepo) Get(_ context.Context, societyID int64, entryID int64) (*models.VisitorEntry, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.getErr != nil {
		return nil, r.getErr
	}
	if r.entry == nil || r.entry.SocietyID != societyID || r.entry.ID != entryID {
		return nil, nil
	}
	copy := *r.entry
	return &copy, nil
}
func (r *inviteSecurityEntryRepo) GetByQRHash(context.Context, string) (*models.VisitorEntry, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.entry == nil {
		return nil, nil
	}
	copy := *r.entry
	return &copy, nil
}
func (r *inviteSecurityEntryRepo) GetForUpdate(_ context.Context, societyID int64, entryID int64) (*models.VisitorEntry, error) {
	r.lockCalls++
	return r.Get(context.Background(), societyID, entryID)
}
func (r *inviteSecurityEntryRepo) GetByInviteID(_ context.Context, inviteID int64) (*models.VisitorEntry, error) {
	if r.entry == nil || r.entry.InviteID == nil || *r.entry.InviteID != inviteID {
		return nil, nil
	}
	return r.entry, nil
}
func (r *inviteSecurityEntryRepo) Approve(_ context.Context, societyID int64, entryID int64, actorUserID int64, qrHash string, qrExpiresAt time.Time) (*models.VisitorEntry, error) {
	r.approvedCalls++
	if r.approveNil {
		return nil, nil
	}
	if r.entry == nil || r.entry.SocietyID != societyID || r.entry.ID != entryID {
		return nil, nil
	}
	approved := *r.entry
	approved.Status = models.VisitorStatusApproved
	approved.ApprovedBy = &actorUserID
	approved.QRExpiresAt = &qrExpiresAt
	r.entry = &approved
	return r.entry, nil
}
func (r *inviteSecurityEntryRepo) Reject(_ context.Context, societyID int64, entryID int64, actorUserID int64, _ string) (*models.VisitorEntry, error) {
	r.rejectedCalls++
	if r.rejectNil {
		return nil, nil
	}
	if r.entry == nil || r.entry.SocietyID != societyID || r.entry.ID != entryID {
		return nil, nil
	}
	rejected := *r.entry
	rejected.Status = models.VisitorStatusRejected
	rejected.RejectedBy = &actorUserID
	r.entry = &rejected
	return r.entry, nil
}
func (r *inviteSecurityEntryRepo) CheckIn(_ context.Context, societyID int64, entryID int64, guardUserID int64) (*models.VisitorEntry, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.checkInErr != nil {
		return nil, r.checkInErr
	}
	if r.checkInNil || r.entry == nil || r.entry.SocietyID != societyID || r.entry.ID != entryID || r.entry.Status != models.VisitorStatusApproved {
		return nil, nil
	}
	checkedIn := *r.entry
	now := time.Now()
	checkedIn.Status = models.VisitorStatusCheckedIn
	checkedIn.CheckedInAt = &now
	checkedIn.HandledByGuardID = &guardUserID
	checkedIn.Visitor = nil
	stored := checkedIn
	r.entry = &stored
	return &checkedIn, nil
}

func (r *inviteSecurityEntryRepo) CheckOut(_ context.Context, societyID int64, entryID int64, guardUserID int64) (*models.VisitorEntry, error) {
	if r.entry == nil || r.entry.SocietyID != societyID || r.entry.ID != entryID || r.entry.Status != models.VisitorStatusCheckedIn {
		return nil, nil
	}
	r.entry.Status = models.VisitorStatusCheckedOut
	r.entry.HandledByGuardID = &guardUserID
	result := *r.entry
	// UPDATE RETURNING does not include the joined visitor summary.
	result.Visitor = nil
	return &result, nil
}

type checkOutNotifier struct {
	notificationsvc.NotificationService
	entries chan *models.VisitorEntry
}

func (n *checkOutNotifier) SendVisitorCheckOut(_ context.Context, entry *models.VisitorEntry) error {
	n.entries <- entry
	return nil
}

func TestCheckOutEnrichesNotificationWithoutFailingCommittedCheckout(t *testing.T) {
	for _, lookupFails := range []bool{false, true} {
		t.Run(map[bool]string{false: "named visitor", true: "lookup unavailable"}[lookupFails], func(t *testing.T) {
			notifier := &checkOutNotifier{entries: make(chan *models.VisitorEntry, 1)}
			_, svc, _, repo, _, events := inviteSecurityServiceWithNotifier(notifier)
			repo.entry = &models.VisitorEntry{ID: 703, SocietyID: 10, FlatID: 20, Status: models.VisitorStatusCheckedIn, Visitor: &models.VisitorSummary{FullName: "Pankaj"}}
			if lookupFails {
				repo.getErr = errors.New("lookup unavailable")
			}
			result, err := svc.CheckOut(context.Background(), 10, 703, 200)
			if err != nil || result == nil || result.Status != models.VisitorStatusCheckedOut {
				t.Fatalf("checkout failed: %v, %#v", err, result)
			}
			if len(events.events) != 1 || events.events[0] != models.VisitorEventTypeCheckedOut {
				t.Fatalf("events: %v", events.events)
			}
			select {
			case entry := <-notifier.entries:
				if !lookupFails && (entry.Visitor == nil || entry.Visitor.FullName != "Pankaj") {
					t.Fatalf("missing visitor name: %#v", entry)
				}
				if lookupFails && entry.Visitor != nil {
					t.Fatal("expected generic fallback")
				}
			case <-time.After(time.Second):
				t.Fatal("missing checkout notification")
			}
		})
	}
}

func TestGuardCabCreationPassesCompanionsToRepository(t *testing.T) {
	for _, count := range []int32{0, 2} {
		_, svc, _, repo, _, _ := inviteSecurityService()
		vehicle := models.VisitorVehicleTypeCab
		details := []map[string]any(nil)
		if count > 0 {
			details = []map[string]any{{"full_name": "Passenger"}, {"phone_number": "9876543210"}}
		}
		_, err := svc.CreateGuardEntry(context.Background(), 10, models.VisitorFormRequest{
			FullName: "Driver", PhoneNumber: strPtr("9876543210"), FlatID: 20, Purpose: models.VisitorPurposeCab,
			VehicleNumber: strPtr("MH12AB1234"), VehicleType: &vehicle, CompanionsCount: count, CompanionDetails: details,
		}, 200)
		if err != nil {
			t.Fatal(err)
		}
		if repo.lastCreateArgs.req.CompanionsCount != count || len(repo.lastCreateArgs.req.CompanionDetails) != int(count) {
			t.Fatalf("cab companions lost: %#v", repo.lastCreateArgs.req)
		}
	}
}
func (r *inviteSecurityEntryRepo) MergeMetadata(_ context.Context, societyID int64, entryID int64, metadata map[string]any) (*models.VisitorEntry, error) {
	if r.entry == nil || r.entry.SocietyID != societyID || r.entry.ID != entryID {
		return nil, nil
	}
	if r.entry.Metadata == nil {
		r.entry.Metadata = map[string]any{}
	}
	for key, value := range metadata {
		r.entry.Metadata[key] = value
	}
	return r.entry, nil
}

type inviteSecurityEventRepo struct {
	events []models.VisitorEventType
}

func (r *inviteSecurityEventRepo) Create(_ context.Context, _ int64, _ int64, _ *int64, eventType models.VisitorEventType, _ *string, _ map[string]any) (*models.VisitorEntryEvent, error) {
	r.events = append(r.events, eventType)
	return &models.VisitorEntryEvent{EventType: eventType}, nil
}
func (r *inviteSecurityEventRepo) List(context.Context, int64, int64) ([]*models.VisitorEntryEvent, error) {
	return nil, nil
}

type inviteAcceptedNotifier struct {
	notificationsvc.NotificationService
	entries chan *models.VisitorEntry
}

func (n *inviteAcceptedNotifier) SendVisitorInviteAccepted(_ context.Context, entry *models.VisitorEntry) error {
	n.entries <- entry
	return nil
}

type checkInNotifier struct {
	notificationsvc.NotificationService
	entries chan *models.VisitorEntry
}

func (n *checkInNotifier) SendVisitorApproved(context.Context, *models.VisitorEntry) error {
	return nil
}

func (n *checkInNotifier) SendVisitorCheckIn(_ context.Context, entry *models.VisitorEntry) error {
	n.entries <- entry
	return nil
}

type entryCreatedNotifier struct {
	notificationsvc.NotificationService
	entries chan *models.VisitorEntry
}

func (n *entryCreatedNotifier) SendVisitorApproved(_ context.Context, entry *models.VisitorEntry) error {
	n.entries <- entry
	return nil
}

func inviteSecurityService() (visitorentrysvc.VisitorInviteService, visitorentrysvc.VisitorEntryService, *inviteSecurityInviteRepo, *inviteSecurityEntryRepo, *inviteSecurityVisitorRepo, *inviteSecurityEventRepo) {
	return inviteSecurityServiceWithNotifier(nil)
}

func inviteSecurityServiceWithNotifier(notifier notificationsvc.NotificationService, txManagers ...repository.TransactionManager) (visitorentrysvc.VisitorInviteService, visitorentrysvc.VisitorEntryService, *inviteSecurityInviteRepo, *inviteSecurityEntryRepo, *inviteSecurityVisitorRepo, *inviteSecurityEventRepo) {
	var txManager repository.TransactionManager = noopTxManager{}
	if len(txManagers) > 0 {
		txManager = txManagers[0]
	}
	societyID := int64(10)
	flatID := int64(20)
	ownerID := int64(100)
	staffID := int64(200)
	residentID := int64(300)
	memberRepo := &inviteSecurityMemberRepo{members: map[[2]int64]*models.SocietyMember{
		{societyID, ownerID}:    {SocietyID: societyID, UserID: ownerID, Role: models.SocietyMemberRoleOwner, Status: models.SocietyMemberStatusActive},
		{societyID, staffID}:    {SocietyID: societyID, UserID: staffID, Role: models.SocietyMemberRoleStaff, Status: models.SocietyMemberStatusActive},
		{societyID, residentID}: {SocietyID: societyID, UserID: residentID, Role: models.SocietyMemberRoleResident, Status: models.SocietyMemberStatusActive},
		{99, staffID}:           {SocietyID: 99, UserID: staffID, Role: models.SocietyMemberRoleStaff, Status: models.SocietyMemberStatusActive},
	}}
	residentRepo := &inviteSecurityResidentRepo{residents: map[[3]int64]*models.FlatResident{
		{societyID, flatID, ownerID}: {SocietyID: societyID, FlatID: flatID, UserID: ownerID, Role: models.FlatResidentRoleOwner, Status: models.FlatResidentStatusActive, IsPrimary: true},
	}}
	flatRepo := &inviteSecurityFlatRepo{flats: map[[2]int64]*models.Flat{
		{societyID, flatID}: {ID: flatID, SocietyID: societyID, FlatNumber: "A-101", Status: models.FlatStatusOccupied, IsActive: true},
	}}
	societyRepo := &inviteSecuritySocietyRepo{societies: map[int64]*models.Society{
		societyID: {ID: societyID, Name: "Apna Gate", SocietyCode: "AG001", Status: models.SocietyStatusActive},
		99:        {ID: 99, Name: "Other Society", SocietyCode: "OS001", Status: models.SocietyStatusActive},
	}}
	inviteRepo := &inviteSecurityInviteRepo{invites: map[int64]*models.VisitorInvite{}, byHash: map[string]int64{}, nextID: 300}
	entryRepo := &inviteSecurityEntryRepo{guardDeskEntryRepo: &guardDeskEntryRepo{}}
	visitorRepo := &inviteSecurityVisitorRepo{}
	eventRepo := &inviteSecurityEventRepo{}
	inviteSvc, entrySvc := visitorentrysvc.NewVisitorService(
		visitorRepo,
		inviteRepo,
		entryRepo,
		eventRepo,
		&inviteQuerySettingSvc{settings: &models.SocietyVisitorSettingsResponse{SocietyID: societyID, IsActive: true, AllowResidentPreApproval: true, QRExpiryMinutes: 60}},
		memberRepo,
		residentRepo,
		flatRepo,
		societyRepo,
		flatauthz.New(memberRepo, residentRepo, flatRepo),
		notifier,
		txManager,
	)
	return inviteSvc, entrySvc, inviteRepo, entryRepo, visitorRepo, eventRepo
}

func TestCreateInviteAllowsFlatManagerAndStoresInviteScope(t *testing.T) {
	inviteSvc, _, inviteRepo, _, _, _ := inviteSecurityService()
	expiresAt := time.Now().Add(2 * time.Hour)

	result, invite, _, err := inviteSvc.CreateInvite(context.Background(), 10, 20, models.VisitorPurposeGuest, 100, &expiresAt)
	if err != nil {
		t.Fatalf("CreateInvite() error = %v", err)
	}
	if result == nil || result.QR == nil || result.QR.Token == "" {
		t.Fatalf("expected share token in response, got %+v", result)
	}
	if invite.SocietyID != 10 || invite.FlatID != 20 || invite.CreatedBy != 100 || invite.Purpose != models.VisitorPurposeGuest {
		t.Fatalf("invite stored wrong scope: %+v", invite)
	}
	if inviteRepo.lastTokenHash == "" || inviteRepo.lastTokenHash != hashInviteToken(result.QR.Token) {
		t.Fatalf("token hash was not stored from response token")
	}
	if !invite.ExpiresAt.Equal(expiresAt) {
		t.Fatalf("expiresAt = %v, want %v", invite.ExpiresAt, expiresAt)
	}
}

func TestCreateInviteRejectsInvalidAndUnauthorizedActors(t *testing.T) {
	inviteSvc, _, _, _, _, _ := inviteSecurityService()

	_, _, _, err := inviteSvc.CreateInvite(context.Background(), 0, 20, models.VisitorPurposeGuest, 100, nil)
	requireAppErrCode(t, err, visitorentrysvc.ErrInvalidVisitorRequest.Code)

	_, _, _, err = inviteSvc.CreateInvite(context.Background(), 10, 20, models.VisitorPurpose(""), 100, nil)
	requireAppErrCode(t, err, visitorentrysvc.ErrInvalidVisitorRequest.Code)

	_, _, _, err = inviteSvc.CreateInvite(context.Background(), 10, 20, models.VisitorPurposeGuest, 999, nil)
	requireAppErrCode(t, err, visitorentrysvc.ErrVisitorForbidden.Code)

	_, _, _, err = inviteSvc.CreateInvite(context.Background(), 10, 999, models.VisitorPurposeGuest, 100, nil)
	requireAppErrCode(t, err, visitorentrysvc.ErrVisitorFlatNotFound.Code)
}

func TestCreateStaffInviteRequiresSameSocietyStaffAndFlat(t *testing.T) {
	inviteSvc, _, _, _, _, _ := inviteSecurityService()

	_, invite, _, err := inviteSvc.CreateStaffInvite(context.Background(), 10, 20, models.VisitorPurposeDelivery, 200, nil)
	if err != nil {
		t.Fatalf("CreateStaffInvite() error = %v", err)
	}
	if invite.SocietyID != 10 || invite.FlatID != 20 || invite.CreatedBy != 200 || invite.Purpose != models.VisitorPurposeDelivery {
		t.Fatalf("staff invite stored wrong scope: %+v", invite)
	}

	_, _, _, err = inviteSvc.CreateStaffInvite(context.Background(), 10, 20, models.VisitorPurposeGuest, 300, nil)
	requireAppErrCode(t, err, visitorentrysvc.ErrVisitorForbidden.Code)

	_, _, _, err = inviteSvc.CreateStaffInvite(context.Background(), 10, 999, models.VisitorPurposeGuest, 200, nil)
	requireAppErrCode(t, err, visitorentrysvc.ErrVisitorFlatNotFound.Code)

	_, _, _, err = inviteSvc.CreateStaffInvite(context.Background(), 10, 20, models.VisitorPurpose("bad"), 200, nil)
	requireAppErrCode(t, err, visitorentrysvc.ErrInvalidVisitorRequest.Code)
}

func TestGetInviteByTokenRequiresGuardScopeAndUsableInvite(t *testing.T) {
	inviteSvc, _, inviteRepo, _, _, _ := inviteSecurityService()
	rawToken := "invite-token"
	invite := &models.VisitorInvite{ID: 501, SocietyID: 10, FlatID: 20, Purpose: models.VisitorPurposeGuest, CreatedBy: 100, Status: models.VisitorInviteStatusActive, ExpiresAt: time.Now().Add(time.Hour)}
	inviteRepo.invites[invite.ID] = invite
	inviteRepo.byHash[hashInviteToken(rawToken)] = invite.ID

	got, err := inviteSvc.GetInviteByToken(context.Background(), 10, 200, rawToken)
	if err != nil {
		t.Fatalf("GetInviteByToken() error = %v", err)
	}
	if got.ID != invite.ID {
		t.Fatalf("invite ID = %d, want %d", got.ID, invite.ID)
	}

	_, err = inviteSvc.GetInviteByToken(context.Background(), 99, 200, rawToken)
	requireAppErrCode(t, err, visitorentrysvc.ErrVisitorInviteNotFound.Code)

	_, err = inviteSvc.GetInviteByToken(context.Background(), 10, 300, rawToken)
	requireAppErrCode(t, err, visitorentrysvc.ErrVisitorForbidden.Code)

	_, err = inviteSvc.GetInviteByToken(context.Background(), 10, 200, "missing")
	requireAppErrCode(t, err, visitorentrysvc.ErrVisitorInviteNotFound.Code)

	invite.Status = models.VisitorInviteStatusCancelled
	_, err = inviteSvc.GetInviteByToken(context.Background(), 10, 200, rawToken)
	requireAppErrCode(t, err, visitorentrysvc.ErrVisitorInviteUnavailable.Code)
}

func TestSubmitInviteFormIsSocietyScopedAndUsesInviteOwnedFields(t *testing.T) {
	inviteSvc, _, inviteRepo, entryRepo, visitorRepo, eventRepo := inviteSecurityService()
	rawToken := "submit-token"
	invite := &models.VisitorInvite{ID: 601, SocietyID: 10, FlatID: 20, Purpose: models.VisitorPurposeGuest, CreatedBy: 100, Status: models.VisitorInviteStatusActive, ExpiresAt: time.Now().Add(time.Hour)}
	inviteRepo.invites[invite.ID] = invite
	inviteRepo.byHash[hashInviteToken(rawToken)] = invite.ID

	req := inviteSubmitRequest()
	req.FlatID = 999
	req.Purpose = models.VisitorPurposeDelivery

	_, err := inviteSvc.SubmitInviteForm(context.Background(), 99, rawToken, req)
	requireAppErrCode(t, err, visitorentrysvc.ErrVisitorInviteNotFound.Code)
	if visitorRepo.created != 0 || entryRepo.createCalls != 0 {
		t.Fatalf("society mismatch created visitor=%d entry=%d", visitorRepo.created, entryRepo.createCalls)
	}

	result, err := inviteSvc.SubmitInviteForm(context.Background(), 10, rawToken, req)
	if err != nil {
		t.Fatalf("SubmitInviteForm() error = %v", err)
	}
	if result == nil || result.Entry == nil || result.QR == nil || result.QR.Token == "" {
		t.Fatalf("expected entry and QR, got %+v", result)
	}
	if entryRepo.lastCreateArgs.societyID != invite.SocietyID || entryRepo.lastCreateArgs.flatID == nil || *entryRepo.lastCreateArgs.flatID != invite.FlatID {
		t.Fatalf("entry used request scope instead of invite scope: %+v", entryRepo.lastCreateArgs)
	}
	if entryRepo.lastCreateArgs.purpose != invite.Purpose || entryRepo.lastCreateArgs.actorUserID == nil || *entryRepo.lastCreateArgs.actorUserID != invite.CreatedBy {
		t.Fatalf("entry used request purpose/creator instead of invite owner fields: %+v", entryRepo.lastCreateArgs)
	}
	if invite.Status != models.VisitorInviteStatusUsed {
		t.Fatalf("invite status = %q, want used", invite.Status)
	}
	if len(eventRepo.events) != 3 {
		t.Fatalf("events = %v, want created/approved/qr_generated", eventRepo.events)
	}
}

func TestSubmitInviteFormNotifiesWithSubmittedVisitorName(t *testing.T) {
	notifier := &inviteAcceptedNotifier{entries: make(chan *models.VisitorEntry, 1)}
	inviteSvc, _, inviteRepo, _, _, _ := inviteSecurityServiceWithNotifier(notifier)
	rawToken := "submit-token"
	invite := &models.VisitorInvite{ID: 602, SocietyID: 10, FlatID: 20, Purpose: models.VisitorPurposeGuest, CreatedBy: 100, Status: models.VisitorInviteStatusActive, ExpiresAt: time.Now().Add(time.Hour)}
	inviteRepo.invites[invite.ID] = invite
	inviteRepo.byHash[hashInviteToken(rawToken)] = invite.ID
	req := inviteSubmitRequest()
	req.FullName = "Pankaj Kholiya"

	result, err := inviteSvc.SubmitInviteForm(context.Background(), 10, rawToken, req)
	if err != nil {
		t.Fatalf("SubmitInviteForm() error = %v", err)
	}
	if result == nil || result.Entry == nil || result.Entry.Visitor == nil || result.Entry.Visitor.FullName != "Pankaj Kholiya" {
		t.Fatalf("response entry missing submitted visitor name: %+v", result)
	}

	select {
	case entry := <-notifier.entries:
		if entry == nil || entry.Visitor == nil || entry.Visitor.FullName != "Pankaj Kholiya" {
			t.Fatalf("notification entry missing submitted visitor name: %+v", entry)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for invite accepted notification")
	}
}

func TestCreateEntryReturnsAndNotifiesWithCreatedVisitor(t *testing.T) {
	notifier := &entryCreatedNotifier{entries: make(chan *models.VisitorEntry, 1)}
	_, entrySvc, _, _, _, _ := inviteSecurityServiceWithNotifier(notifier)
	phone := "9999999999"
	email := "pankaj@example.com"
	req := models.VisitorFormRequest{
		FullName:    "Pankaj Kholiya",
		PhoneNumber: &phone,
		Email:       &email,
		FlatID:      20,
		Purpose:     models.VisitorPurposeGuest,
	}

	result, err := entrySvc.CreateQuickLinkEntry(context.Background(), 10, req)
	if err != nil {
		t.Fatalf("CreateQuickLinkEntry() error = %v", err)
	}
	if result == nil || result.Entry == nil || result.Entry.Visitor == nil {
		t.Fatalf("response entry missing visitor summary: %+v", result)
	}
	visitor := result.Entry.Visitor
	if visitor.FullName != req.FullName || visitor.PhoneNumber == nil || *visitor.PhoneNumber != phone || visitor.Email == nil || *visitor.Email != email {
		t.Fatalf("response visitor = %+v, want submitted visitor details", visitor)
	}

	select {
	case entry := <-notifier.entries:
		if entry == nil || entry.Visitor == nil || entry.Visitor.FullName != req.FullName {
			t.Fatalf("notification entry missing created visitor: %+v", entry)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for entry-created notification")
	}
}

func TestCheckInByEntryIDNotifiesWithVisitorName(t *testing.T) {
	notifier := &checkInNotifier{entries: make(chan *models.VisitorEntry, 1)}
	_, entrySvc, _, entryRepo, _, eventRepo := inviteSecurityServiceWithNotifier(notifier)
	entryRepo.entry = &models.VisitorEntry{
		ID:        703,
		SocietyID: 10,
		FlatID:    20,
		Status:    models.VisitorStatusApproved,
		Purpose:   models.VisitorPurposeGuest,
		Visitor:   &models.VisitorSummary{FullName: "Pankaj Kholiya"},
		Metadata:  map[string]any{},
	}

	result, err := entrySvc.CheckInByEntryID(context.Background(), 10, 703, 200)
	if err != nil {
		t.Fatalf("CheckInByEntryID() error = %v", err)
	}
	if result == nil || result.Visitor == nil || result.Visitor.FullName != "Pankaj Kholiya" {
		t.Fatalf("checked-in entry missing visitor name: %+v", result)
	}
	if len(eventRepo.events) != 1 || eventRepo.events[0] != models.VisitorEventTypeCheckedIn {
		t.Fatalf("events = %v, want checked_in", eventRepo.events)
	}

	select {
	case entry := <-notifier.entries:
		if entry == nil || entry.Visitor == nil || entry.Visitor.FullName != "Pankaj Kholiya" {
			t.Fatalf("check-in notification entry missing visitor name: %+v", entry)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for check-in notification")
	}
}

func TestGuardApproveAndCheckInNotifiesWithVisitorName(t *testing.T) {
	notifier := &checkInNotifier{entries: make(chan *models.VisitorEntry, 1)}
	_, entrySvc, _, entryRepo, _, _ := inviteSecurityServiceWithNotifier(notifier)
	entryRepo.entry = &models.VisitorEntry{
		ID:        704,
		SocietyID: 10,
		FlatID:    20,
		Status:    models.VisitorStatusWaitingApproval,
		Purpose:   models.VisitorPurposeGuest,
		Visitor:   &models.VisitorSummary{FullName: "Pankaj Kholiya"},
		Metadata:  map[string]any{},
	}

	result, err := entrySvc.GuardApproveAndCheckIn(context.Background(), 10, 704, 200, visitorentrysvc.GuardApproveOptions{})
	if err != nil {
		t.Fatalf("GuardApproveAndCheckIn() error = %v", err)
	}
	if result == nil || result.Visitor == nil || result.Visitor.FullName != "Pankaj Kholiya" {
		t.Fatalf("checked-in entry missing visitor name: %+v", result)
	}

	select {
	case entry := <-notifier.entries:
		if entry == nil || entry.Visitor == nil || entry.Visitor.FullName != "Pankaj Kholiya" {
			t.Fatalf("check-in notification entry missing visitor name: %+v", entry)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for check-in notification")
	}
}

func TestApproveEntryChecksAuthStateAndRecordsQR(t *testing.T) {
	_, entrySvc, _, entryRepo, _, eventRepo := inviteSecurityService()
	entryRepo.entry = &models.VisitorEntry{ID: 701, SocietyID: 10, FlatID: 20, Status: models.VisitorStatusWaitingApproval, Purpose: models.VisitorPurposeGuest, Metadata: map[string]any{}}

	result, err := entrySvc.ApproveEntry(context.Background(), 10, 701, 100)
	if err != nil {
		t.Fatalf("ApproveEntry() error = %v", err)
	}
	if result == nil || result.Entry == nil || result.QR == nil || result.QR.Token == "" {
		t.Fatalf("expected approved entry and QR, got %+v", result)
	}
	if result.Entry.Status != models.VisitorStatusApproved || entryRepo.lockCalls != 1 || entryRepo.approvedCalls != 1 {
		t.Fatalf("approval did not lock/update correctly: entry=%+v locks=%d approves=%d", result.Entry, entryRepo.lockCalls, entryRepo.approvedCalls)
	}
	if len(eventRepo.events) != 2 {
		t.Fatalf("events = %v, want approved/qr_generated", eventRepo.events)
	}

	for _, status := range []models.VisitorStatus{
		models.VisitorStatusApproved,
		models.VisitorStatusCheckedIn,
		models.VisitorStatusCheckedOut,
	} {
		entryRepo.entry.Status = status
		idempotent, idempotentErr := entrySvc.ApproveEntry(context.Background(), 10, 701, 100)
		if idempotentErr != nil || idempotent == nil || idempotent.Entry == nil ||
			idempotent.Entry.Status != status || idempotent.QR != nil {
			t.Fatalf("idempotent approval for %s = %+v, err = %v", status, idempotent, idempotentErr)
		}
	}
	if entryRepo.approvedCalls != 1 || len(eventRepo.events) != 2 {
		t.Fatalf("idempotent approval added side effects: approves=%d events=%v", entryRepo.approvedCalls, eventRepo.events)
	}

	_, err = entrySvc.ApproveEntry(context.Background(), 10, 701, 999)
	requireAppErrCode(t, err, visitorentrysvc.ErrVisitorForbidden.Code)

	entryRepo.entry.Status = models.VisitorStatusWaitingApproval
	_, err = entrySvc.ApproveEntry(context.Background(), 10, 701, 999)
	requireAppErrCode(t, err, visitorentrysvc.ErrVisitorForbidden.Code)

	_, err = entrySvc.ApproveEntry(context.Background(), 99, 701, 100)
	requireAppErrCode(t, err, visitorentrysvc.ErrVisitorEntryNotFound.Code)

	entryRepo.approveNil = true
	_, err = entrySvc.ApproveEntry(context.Background(), 10, 701, 100)
	requireAppErrCode(t, err, visitorentrysvc.ErrVisitorInvalidState.Code)
}

func TestRejectEntryAuthorizesBeforeIdempotentResult(t *testing.T) {
	_, entrySvc, _, entryRepo, _, eventRepo := inviteSecurityService()
	entryRepo.entry = &models.VisitorEntry{
		ID: 702, SocietyID: 10, FlatID: 20,
		Status: models.VisitorStatusWaitingApproval,
	}

	result, err := entrySvc.RejectEntry(context.Background(), 10, 702, "No longer expected", 100)
	if err != nil || result == nil || result.Entry == nil ||
		result.Entry.Status != models.VisitorStatusRejected {
		t.Fatalf("RejectEntry() = %+v, err = %v", result, err)
	}
	if entryRepo.rejectedCalls != 1 || len(eventRepo.events) != 1 {
		t.Fatalf("rejection side effects: rejects=%d events=%v", entryRepo.rejectedCalls, eventRepo.events)
	}

	idempotent, err := entrySvc.RejectEntry(context.Background(), 10, 702, "Duplicate", 100)
	if err != nil || idempotent == nil || idempotent.Entry == nil ||
		idempotent.Entry.Status != models.VisitorStatusRejected {
		t.Fatalf("idempotent rejection = %+v, err = %v", idempotent, err)
	}
	if entryRepo.rejectedCalls != 1 || len(eventRepo.events) != 1 {
		t.Fatalf("idempotent rejection added side effects: rejects=%d events=%v", entryRepo.rejectedCalls, eventRepo.events)
	}

	_, err = entrySvc.RejectEntry(context.Background(), 10, 702, "Unauthorized", 999)
	requireAppErrCode(t, err, visitorentrysvc.ErrVisitorForbidden.Code)

	entryRepo.entry.Status = models.VisitorStatusApproved
	_, err = entrySvc.RejectEntry(context.Background(), 10, 702, "Opposite decision", 100)
	requireAppErrCode(t, err, visitorentrysvc.ErrVisitorInvalidState.Code)
}

func inviteSubmitRequest() models.VisitorFormRequest {
	phone := "9999999999"
	expectedAt := time.Now().Add(30 * time.Minute)
	return models.VisitorFormRequest{
		FullName:    "Test Visitor",
		PhoneNumber: &phone,
		ExpectedAt:  &expectedAt,
	}
}

func requireAppErrCode(t *testing.T, err error, want string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error code %s, got nil", want)
	}
	var appErr *models.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError code %s, got %T: %v", want, err, err)
	}
	if appErr.Code != want {
		t.Fatalf("error code = %s, want %s", appErr.Code, want)
	}
}

var (
	_ repository.SocietyMemberRepository     = (*inviteSecurityMemberRepo)(nil)
	_ repository.FlatResidentRepository      = (*inviteSecurityResidentRepo)(nil)
	_ repository.FlatRepository              = (*inviteSecurityFlatRepo)(nil)
	_ repository.SocietyRepository           = (*inviteSecuritySocietyRepo)(nil)
	_ repository.VisitorInviteRepository     = (*inviteSecurityInviteRepo)(nil)
	_ repository.VisitorRepository           = (*inviteSecurityVisitorRepo)(nil)
	_ repository.VisitorEntryRepository      = (*inviteSecurityEntryRepo)(nil)
	_ repository.VisitorEntryEventRepository = (*inviteSecurityEventRepo)(nil)
)
