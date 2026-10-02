package notificationsvc

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-server/internal/models"

	"github.com/google/uuid"
)

const (
	EventVisitorPending        = "visitor.pending"
	EventVisitorApproved       = "visitor.approved"
	EventVisitorRejected       = "visitor.rejected"
	EventVisitorCheckIn        = "visitor.checkin"
	EventVisitorCheckOut       = "visitor.checkout"
	EventVisitorInviteAccepted = "visitor_invite.accepted"
	EventMemberInviteAccepted  = "member_invite.accepted"

	InviteTypeVisitor = "visitor"
	InviteTypeMember  = "member"

	CategoryVisitorDecision  = "visitor_decision"
	CategoryNotificationInfo = "notification_info"
)

func (s *notificationService) SendVisitorApprovalRequested(ctx context.Context, entry *models.VisitorEntry) error {
	if entry == nil {
		return nil
	}
	residents, err := s.flatResidentUserIDs(ctx, entry.SocietyID, entry.FlatID)
	if err != nil {
		return err
	}
	staff, err := s.societyStaffUserIDs(ctx, entry.SocietyID)
	if err != nil {
		return err
	}
	copy := visitorApprovalRequestedCopy(entry)
	// Approval requests belong to residents; the initiating guard gets local UI feedback.
	return s.createAndSend(ctx, excludeUserIDs(residents, staff), notificationSpec{
		Type:       EventVisitorPending,
		Alias:      "visitor_pending_approval",
		Title:      copy.Title,
		Body:       copy.Body,
		SocietyID:  &entry.SocietyID,
		FlatID:     &entry.FlatID,
		CategoryID: CategoryVisitorDecision,
		Data:       visitorEntryData(entry, EventVisitorPending, "visitor_pending_approval"),
	})
}

func (s *notificationService) SendVisitorApproved(ctx context.Context, entry *models.VisitorEntry) error {
	if entry == nil || entry.Status != models.VisitorStatusApproved {
		return nil
	}
	recipients, err := s.societyStaffUserIDs(ctx, entry.SocietyID)
	if err != nil {
		return err
	}
	copy := visitorApprovedByFlatCopy(entry)
	var flatID *int64
	if entry.FlatID > 0 {
		flatID = &entry.FlatID
	}
	return s.createAndSend(ctx, recipients, notificationSpec{
		Type:       EventVisitorApproved,
		Alias:      "visitor_approved",
		Title:      copy.Title,
		Body:       copy.Body,
		SocietyID:  &entry.SocietyID,
		FlatID:     flatID,
		CategoryID: CategoryNotificationInfo,
		Data:       visitorEntryData(entry, EventVisitorApproved, "visitor_approved"),
	})
}

func (s *notificationService) SendVisitorRejected(ctx context.Context, entry *models.VisitorEntry) error {
	if entry == nil {
		return nil
	}
	staff, err := s.societyStaffUserIDs(ctx, entry.SocietyID)
	if err != nil {
		return err
	}
	copy := visitorRejectedByFlatCopy(entry)
	return s.createAndSend(ctx, staff, notificationSpec{
		Type:       EventVisitorRejected,
		Alias:      "visitor_rejected",
		Title:      copy.Title,
		Body:       copy.Body,
		SocietyID:  &entry.SocietyID,
		FlatID:     &entry.FlatID,
		CategoryID: CategoryNotificationInfo,
		Data:       visitorEntryData(entry, EventVisitorRejected, "visitor_rejected"),
	})
}

func (s *notificationService) SendVisitorCheckIn(ctx context.Context, entry *models.VisitorEntry) error {
	if entry == nil {
		return nil
	}
	if entry.IsSocietyWideGuardService() {
		if entry.Status != models.VisitorStatusCheckedIn {
			return nil
		}
		residents, err := s.societyResidentUserIDs(ctx, entry.SocietyID)
		if err != nil {
			return err
		}
		data := visitorEntryData(entry, EventVisitorCheckIn, "visitor_checked_in")
		delete(data, "flat_id")
		delete(data, "flat_number")
		delete(data, "block")
		data["scope"] = "society"
		data["purpose"] = "service"
		title := "Service provider has arrived"
		if entry.ServiceProvider != nil && strings.TrimSpace(*entry.ServiceProvider) != "" {
			title = strings.TrimSpace(*entry.ServiceProvider) + " has arrived"
		}
		body := "A service provider has entered the society."
		if name := visitorDisplayName(entry); name != "" && entry.ServiceProvider != nil && strings.TrimSpace(*entry.ServiceProvider) != "" {
			body = fmt.Sprintf("%s from %s has entered the society.", name, strings.TrimSpace(*entry.ServiceProvider))
		}
		eventKey := fmt.Sprintf("visitor_service_checked_in:%d", entry.ID)
		return s.createAndSend(ctx, residents, notificationSpec{
			Type: EventVisitorCheckIn, Alias: "visitor_checked_in",
			Title: title, Body: body,
			SocietyID: &entry.SocietyID, EventKey: &eventKey,
			CategoryID: CategoryNotificationInfo, Data: data,
		})
	}
	residents, err := s.flatResidentUserIDs(ctx, entry.SocietyID, entry.FlatID)
	if err != nil {
		return err
	}
	copy := visitorCheckInCopy(entry)
	return s.createAndSend(ctx, residents, notificationSpec{
		Type:       EventVisitorCheckIn,
		Alias:      "visitor_checked_in",
		Title:      copy.Title,
		Body:       copy.Body,
		SocietyID:  &entry.SocietyID,
		FlatID:     &entry.FlatID,
		CategoryID: CategoryNotificationInfo,
		Data:       visitorEntryData(entry, EventVisitorCheckIn, "visitor_checked_in"),
	})
}

func (s *notificationService) SendVisitorCheckOut(ctx context.Context, entry *models.VisitorEntry) error {
	if entry == nil {
		return nil
	}
	residents, err := s.flatResidentUserIDs(ctx, entry.SocietyID, entry.FlatID)
	if err != nil {
		return err
	}
	copy := visitorCheckOutCopy(entry)
	return s.createAndSend(ctx, residents, notificationSpec{
		Type:       EventVisitorCheckOut,
		Alias:      "visitor_checked_out",
		Title:      copy.Title,
		Body:       copy.Body,
		SocietyID:  &entry.SocietyID,
		FlatID:     &entry.FlatID,
		CategoryID: CategoryNotificationInfo,
		Data:       visitorEntryData(entry, EventVisitorCheckOut, "visitor_checked_out"),
	})
}

func (s *notificationService) SendVisitorInviteAccepted(ctx context.Context, entry *models.VisitorEntry) error {
	if entry == nil || entry.InviteID == nil {
		return nil
	}
	recipients, err := s.flatResidentUserIDs(ctx, entry.SocietyID, entry.FlatID)
	if err != nil {
		return err
	}
	eventKey := fmt.Sprintf("%s:%d:%d", EventVisitorInviteAccepted, *entry.InviteID, entry.ID)
	copy := visitorInviteAcceptedCopy(entry)
	return s.createAndSend(ctx, recipients, notificationSpec{
		Type:       EventVisitorInviteAccepted,
		Alias:      "visitor_invite_accepted",
		Title:      copy.Title,
		Body:       copy.Body,
		SocietyID:  &entry.SocietyID,
		FlatID:     &entry.FlatID,
		EventKey:   &eventKey,
		CategoryID: CategoryNotificationInfo,
		Data:       visitorEntryData(entry, EventVisitorInviteAccepted, "visitor_invite_accepted"),
	})
}

func (s *notificationService) SendMemberInviteAccepted(ctx context.Context, invite *models.FlatMemberInvite, flatNumber string, joinedName string, residentID int64) error {
	if invite == nil {
		return nil
	}
	recipients, err := s.flatResidentUserIDsExceptResident(ctx, invite.SocietyID, invite.FlatID, residentID)
	if err != nil {
		return err
	}
	name := strings.TrimSpace(joinedName)
	if name == "" {
		name = strings.TrimSpace(invite.FullName)
	}
	if name == "" {
		name = "A member"
	}
	title := name + " joined your flat"
	body := name + " accepted your invitation and is now connected to "
	if flat := strings.TrimSpace(flatNumber); flat != "" {
		body += "Flat " + flat + "."
	} else {
		body += "your flat."
	}
	eventKey := fmt.Sprintf("%s:%d:%d", EventMemberInviteAccepted, invite.ID, residentID)
	data := map[string]any{
		"type":        EventMemberInviteAccepted,
		"event":       "member_invite_accepted",
		"invite_type": InviteTypeMember,
		"society_id":  strconv.FormatInt(invite.SocietyID, 10),
		"flat_id":     strconv.FormatInt(invite.FlatID, 10),
		"invite_id":   strconv.FormatInt(invite.ID, 10),
	}
	if residentID > 0 {
		data["resident_id"] = strconv.FormatInt(residentID, 10)
	}
	data["member_name"] = name
	if flat := strings.TrimSpace(flatNumber); flat != "" {
		data["flat_number"] = flat
	}
	return s.createAndSend(ctx, recipients, notificationSpec{
		Type:       EventMemberInviteAccepted,
		Alias:      "member_invite_accepted",
		Title:      title,
		Body:       body,
		SocietyID:  &invite.SocietyID,
		FlatID:     &invite.FlatID,
		EventKey:   &eventKey,
		CategoryID: CategoryNotificationInfo,
		Data:       data,
	})
}

type notificationSpec struct {
	Type       string
	Alias      string
	Title      string
	Body       string
	SocietyID  *int64
	FlatID     *int64
	EventKey   *string
	CategoryID string
	Data       map[string]any
}

type NotificationCopy struct {
	Title string
	Body  string
}

func (s *notificationService) createAndSend(ctx context.Context, userIDs []int64, spec notificationSpec) error {
	userIDs = uniqueUserIDs(userIDs)
	if len(userIDs) == 0 {
		return nil
	}
	if s.DurableNotifications() {
		return s.enqueueSpec(ctx, userIDs, spec)
	}

	var sendErrors []error
	for _, userID := range userIDs {
		id := uuid.NewString()
		data := cloneData(spec.Data)
		data["notification_id"] = id
		data["type"] = spec.Type
		data["event"] = spec.Alias
		data["occurred_at"] = time.Now().UTC().Format(time.RFC3339)
		if spec.CategoryID != "" {
			data["category_id"] = spec.CategoryID
		}

		if s.notifications != nil {
			created, err := s.notifications.Create(ctx, models.NotificationCreate{
				ID:        id,
				UserID:    userID,
				SocietyID: spec.SocietyID,
				FlatID:    spec.FlatID,
				Type:      spec.Type,
				Title:     spec.Title,
				Body:      spec.Body,
				Data:      data,
				EventKey:  spec.EventKey,
			})
			if err != nil {
				return err
			}
			if created == nil {
				continue
			}
		}

		pushData := stringMap(data)
		if err := s.sendPushToUsers(ctx, []int64{userID}, models.NotificationPayload{
			Title:      spec.Title,
			Body:       spec.Body,
			Data:       pushData,
			CategoryID: spec.CategoryID,
		}); err != nil {
			// Continue creating durable inbox entries for the remaining recipients.
			sendErrors = append(sendErrors, err)
		}
	}
	return errors.Join(sendErrors...)
}

func (s *notificationService) flatResidentUserIDs(ctx context.Context, societyID int64, flatID int64) ([]int64, error) {
	return s.flatResidentUserIDsExceptResident(ctx, societyID, flatID, 0)
}

func (s *notificationService) flatResidentUserIDsExceptResident(ctx context.Context, societyID int64, flatID int64, excludedResidentID int64) ([]int64, error) {
	if s.residents == nil {
		return nil, nil
	}
	active := string(models.FlatResidentStatusActive)
	filter := &models.FlatResidentFilter{SocietyID: &societyID, FlatID: &flatID, Status: &active}
	residents, err := s.residents.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	userIDs := make([]int64, 0, len(residents))
	var excludedUserIDs []int64
	for _, resident := range residents {
		if excludedResidentID > 0 && resident.ID == excludedResidentID {
			excludedUserIDs = append(excludedUserIDs, resident.UserID)
		}
		userIDs = append(userIDs, resident.UserID)
	}
	return excludeUserIDs(userIDs, excludedUserIDs), nil
}

func (s *notificationService) societyResidentUserIDs(ctx context.Context, societyID int64) ([]int64, error) {
	if s.residents == nil || societyID <= 0 {
		return nil, nil
	}
	active := string(models.FlatResidentStatusActive)
	filter := &models.FlatResidentFilter{SocietyID: &societyID, Status: &active, Limit: 100}
	var userIDs []int64
	for {
		residents, err := s.residents.List(ctx, filter)
		if err != nil {
			return nil, err
		}
		for _, resident := range residents {
			if resident != nil && resident.SocietyID == societyID && resident.Status == models.FlatResidentStatusActive {
				userIDs = append(userIDs, resident.UserID)
			}
		}
		if len(residents) < int(filter.Limit) {
			return uniqueUserIDs(userIDs), nil
		}
		filter.Offset += filter.Limit
	}
}

func (s *notificationService) societyStaffUserIDs(ctx context.Context, societyID int64) ([]int64, error) {
	if s.members == nil {
		return nil, nil
	}
	staffRole := string(models.SocietyMemberRoleStaff)
	activeStatus := string(models.SocietyMemberStatusActive)
	members, err := s.members.List(ctx, models.ListSocietyMembersFilter{
		SocietyID: societyID,
		Role:      &staffRole,
		Status:    &activeStatus,
		Limit:     500,
	})
	if err != nil {
		return nil, err
	}
	userIDs := make([]int64, 0, len(members))
	for _, member := range members {
		userIDs = append(userIDs, member.UserID)
	}
	return uniqueUserIDs(userIDs), nil
}

func visitorEntryData(entry *models.VisitorEntry, eventType string, alias string) map[string]any {
	data := map[string]any{
		"type":        eventType,
		"event":       alias,
		"invite_type": InviteTypeVisitor,
	}
	if entry == nil {
		return data
	}
	data["society_id"] = strconv.FormatInt(entry.SocietyID, 10)
	if entry.FlatID > 0 {
		data["flat_id"] = strconv.FormatInt(entry.FlatID, 10)
	}
	data["entry_id"] = strconv.FormatInt(entry.ID, 10)
	if entry.InviteID != nil && *entry.InviteID > 0 {
		data["invite_id"] = strconv.FormatInt(*entry.InviteID, 10)
	}
	if name := visitorDisplayName(entry); name != "" {
		data["visitor_name"] = name
	}
	if entry.Purpose == models.VisitorPurposeDelivery && entry.DeliveryPartner != nil {
		if partner := strings.TrimSpace(*entry.DeliveryPartner); partner != "" {
			data["delivery_partner"] = partner
		}
	}
	if name := approverDisplayName(entry); name != "" {
		data["approver_name"] = name
	}
	if entry.Flat != nil {
		if flatNumber := strings.TrimSpace(entry.Flat.FlatNumber); flatNumber != "" {
			data["flat_number"] = flatNumber
		}
		if entry.Flat.Block != nil {
			if block := strings.TrimSpace(*entry.Flat.Block); block != "" {
				data["block"] = block
			}
		}
	}
	return data
}

func visitorDisplayName(entry *models.VisitorEntry) string {
	if entry == nil {
		return ""
	}
	if entry.Visitor != nil {
		if name := strings.TrimSpace(entry.Visitor.FullName); name != "" {
			return name
		}
	}
	if entry.Purpose == models.VisitorPurposeDelivery && entry.DeliveryPartner != nil {
		return strings.TrimSpace(*entry.DeliveryPartner)
	}
	return ""
}

func visitorInviteAcceptedCopy(entry *models.VisitorEntry) NotificationCopy {
	name := visitorDisplayName(entry)
	return NotificationCopy{
		Title: personEventTitle(name, "Guest accepted your invite", "accepted your invite"),
		Body:  namedSentence(name, "The guest has completed the visitor details. The visitor pass is ready.", "has completed the visitor details. The visitor pass is ready."),
	}
}

func visitorCheckInCopy(entry *models.VisitorEntry) NotificationCopy {
	name := visitorDisplayName(entry)
	return NotificationCopy{
		Title: personEventTitle(name, "Visitor has arrived", "has arrived"),
		Body:  namedSentence(name, "The visitor just checked in at the society gate.", "just checked in at the society gate."),
	}
}

func visitorCheckOutCopy(entry *models.VisitorEntry) NotificationCopy {
	name := visitorDisplayName(entry)
	return NotificationCopy{
		Title: personEventTitle(name, "Visitor has left", "has left"),
		Body:  namedSentence(name, "The visitor has checked out and left the society.", "has checked out and left the society."),
	}
}

func visitorApprovalRequestedCopy(entry *models.VisitorEntry) NotificationCopy {
	name := visitorDisplayName(entry)
	return NotificationCopy{
		Title: personEventTitle(name, "Visitor is at the gate", "is at the gate"),
		Body:  namedSentence(name, "Visitor is waiting for your approval. Approve or decline the request.", "is waiting for your approval. Approve or decline the request."),
	}
}

func visitorApprovedByFlatCopy(entry *models.VisitorEntry) NotificationCopy {
	name := visitorDisplayName(entry)
	return NotificationCopy{
		Title: "Visitor approved",
		Body:  namedSentence(name, "The visitor has been approved and can now enter the society.", "has been approved and can now enter the society."),
	}
}

func visitorRejectedByFlatCopy(entry *models.VisitorEntry) NotificationCopy {
	name := visitorDisplayName(entry)
	return NotificationCopy{
		Title: "Visitor request declined",
		Body:  possessiveSentence(name, "The visitor's entry request was declined.", "entry request was declined."),
	}
}

func personEventTitle(name string, fallback string, event string) string {
	if name == "" {
		return fallback
	}
	return name + " " + event
}

func namedSentence(name string, fallback string, sentence string) string {
	if name == "" {
		return fallback
	}
	return name + " " + sentence
}

func possessiveSentence(name string, fallback string, sentence string) string {
	if name == "" {
		return fallback
	}
	return name + "'s " + sentence
}

func approverDisplayName(entry *models.VisitorEntry) string {
	if entry == nil || entry.ApproverName == nil {
		return ""
	}
	return strings.TrimSpace(*entry.ApproverName)
}

func visitorFlatNumberLabel(entry *models.VisitorEntry) string {
	if entry == nil || entry.Flat == nil {
		return ""
	}
	flatNumber := strings.TrimSpace(entry.Flat.FlatNumber)
	if flatNumber != "" {
		return fmt.Sprintf("Flat %s", flatNumber)
	}
	return ""
}

func visitorDecisionBody(entry *models.VisitorEntry, decision string) string {
	name := visitorDisplayName(entry)
	flat := visitorFlatNumberLabel(entry)
	approver := approverDisplayName(entry)
	if name != "" && flat != "" && approver != "" {
		return fmt.Sprintf("%s was %s for %s by %s.", name, decision, flat, approver)
	}
	if name != "" && flat != "" {
		return fmt.Sprintf("%s was %s for %s.", name, decision, flat)
	}
	if name != "" {
		return fmt.Sprintf("%s was %s.", name, decision)
	}
	if flat != "" {
		return fmt.Sprintf("The visitor was %s for %s.", decision, flat)
	}
	return fmt.Sprintf("The visitor was %s.", decision)
}

func cloneData(data map[string]any) map[string]any {
	clone := make(map[string]any, len(data)+3)
	for key, value := range data {
		clone[key] = value
	}
	return clone
}

func stringMap(data map[string]any) map[string]string {
	result := make(map[string]string, len(data))
	for key, value := range data {
		switch typed := value.(type) {
		case string:
			result[key] = typed
		case fmt.Stringer:
			result[key] = typed.String()
		default:
			result[key] = fmt.Sprint(typed)
		}
	}
	return result
}

func uniqueUserIDs(userIDs []int64) []int64 {
	seen := make(map[int64]struct{}, len(userIDs))
	unique := make([]int64, 0, len(userIDs))
	for _, userID := range userIDs {
		if userID <= 0 {
			continue
		}
		if _, exists := seen[userID]; exists {
			continue
		}
		seen[userID] = struct{}{}
		unique = append(unique, userID)
	}
	return unique
}

func excludeUserIDs(userIDs []int64, excluded []int64) []int64 {
	if len(userIDs) == 0 || len(excluded) == 0 {
		return uniqueUserIDs(userIDs)
	}
	blocked := make(map[int64]struct{}, len(excluded))
	for _, userID := range excluded {
		blocked[userID] = struct{}{}
	}
	filtered := make([]int64, 0, len(userIDs))
	for _, userID := range userIDs {
		if _, exists := blocked[userID]; exists {
			continue
		}
		filtered = append(filtered, userID)
	}
	return uniqueUserIDs(filtered)
}
