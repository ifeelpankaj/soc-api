package repository

import (
	"github.com/google/uuid"
	"go-server/internal/db"
	"go-server/internal/repositories/contracts"
)

func mapFinancialGetMaintenanceReminderDeliveryRow(v db.GetMaintenanceReminderDeliveryRow) contracts.GetMaintenanceReminderDeliveryRecord {
	return contracts.GetMaintenanceReminderDeliveryRecord{
		BillNumber:             v.BillNumber,
		BillingMonth:           v.BillingMonth.Time,
		DueDate:                v.DueDate.Time,
		OutstandingAmountPaise: v.OutstandingAmountPaise,
	}
}

func mapFinancialGetMaintenanceReviewRow(v db.GetMaintenanceReviewRow) contracts.GetMaintenanceReviewRecord {
	return contracts.GetMaintenanceReviewRecord{
		SnapshotHash: v.SnapshotHash,
		ExpiresAt:    pgTimestamptzToTime(v.ExpiresAt),
	}
}

func mapFinancialGetUPIReportRow(v db.GetUPIReportRow) contracts.GetUPIReportRecord {
	return contracts.GetUPIReportRecord{
		MaintenancePaymentReport: mapFinancialMaintenancePaymentReport(v.MaintenancePaymentReport),
		BillNumber:               v.BillNumber,
	}
}

func mapFinancialListMaintenanceReminderCandidatesRow(v db.ListMaintenanceReminderCandidatesRow) contracts.ListMaintenanceReminderCandidatesRecord {
	return contracts.ListMaintenanceReminderCandidatesRecord{
		ID:       v.ID,
		DueDate:  v.DueDate.Time,
		Timezone: v.Timezone,
	}
}

func mapFinancialListUPIReportsRow(v db.ListUPIReportsRow) contracts.ListUPIReportsRecord {
	return contracts.ListUPIReportsRecord{
		MaintenancePaymentReport: mapFinancialMaintenancePaymentReport(v.MaintenancePaymentReport),
		BillNumber:               v.BillNumber,
	}
}

func mapFinancialMaintenanceBill(v db.MaintenanceBill) contracts.MaintenanceBill {
	return contracts.MaintenanceBill{
		ID:             v.ID,
		RunID:          v.RunID,
		SocietyID:      v.SocietyID,
		FlatID:         v.FlatID,
		BillingMonth:   v.BillingMonth.Time,
		BillNumber:     v.BillNumber,
		DueDate:        v.DueDate.Time,
		Timezone:       v.Timezone,
		TotalPaise:     v.TotalPaise,
		Snapshot:       v.Snapshot,
		BilledParty:    v.BilledParty,
		CreatedAt:      pgTimestamptzToTime(v.CreatedAt),
		IssuerSnapshot: v.IssuerSnapshot,
	}
}

func mapFinancialMaintenancePayment(v db.MaintenancePayment) contracts.MaintenancePayment {
	return contracts.MaintenancePayment{
		ID:                v.ID,
		SocietyID:         v.SocietyID,
		BillID:            v.BillID,
		ClaimID:           v.ClaimID,
		PayerID:           v.PayerID,
		SettingsVersion:   v.SettingsVersion,
		AmountPaise:       v.AmountPaise,
		Reference:         v.Reference,
		CreditDate:        v.CreditDate.Time,
		ReceiptNumber:     uuid.UUID(v.ReceiptNumber.Bytes),
		VerifiedBy:        v.VerifiedBy,
		VerifiedAt:        pgTimestamptzToTime(v.VerifiedAt),
		Status:            v.Status,
		ReversedBy:        v.ReversedBy,
		ReversedAt:        pgTimestamptzToTime(v.ReversedAt),
		ReversalReason:    v.ReversalReason,
		EvidenceReference: v.EvidenceReference,
	}
}

func mapFinancialMaintenancePaymentAuditEvent(v db.MaintenancePaymentAuditEvent) contracts.MaintenancePaymentAuditEvent {
	return contracts.MaintenancePaymentAuditEvent{
		ID:        v.ID,
		SocietyID: v.SocietyID,
		BillID:    v.BillID,
		ActorID:   v.ActorID,
		Action:    v.Action,
		EntityID:  v.EntityID,
		Details:   v.Details,
		CreatedAt: pgTimestamptzToTime(v.CreatedAt),
	}
}

func mapFinancialMaintenancePaymentClaim(v db.MaintenancePaymentClaim) contracts.MaintenancePaymentClaim {
	return contracts.MaintenancePaymentClaim{
		ID:          v.ID,
		SocietyID:   v.SocietyID,
		BillID:      v.BillID,
		RequestID:   uuid.UUID(v.RequestID.Bytes),
		SubmittedBy: v.SubmittedBy,
		Reference:   v.Reference,
		PaymentDate: v.PaymentDate.Time,
		Status:      v.Status,
		ReviewedBy:  v.ReviewedBy,
		ReviewedAt:  pgTimestamptzToTime(v.ReviewedAt),
		Reason:      v.Reason,
		CreatedAt:   pgTimestamptzToTime(v.CreatedAt),
	}
}

func mapFinancialMaintenancePaymentIdempotency(v db.MaintenancePaymentIdempotency) contracts.MaintenancePaymentIdempotency {
	return contracts.MaintenancePaymentIdempotency{
		SocietyID:   v.SocietyID,
		ActorID:     v.ActorID,
		Operation:   v.Operation,
		Key:         v.Key,
		RequestHash: v.RequestHash,
		Response:    v.Response,
		CreatedAt:   pgTimestamptzToTime(v.CreatedAt),
	}
}

func mapFinancialMaintenancePaymentReport(v db.MaintenancePaymentReport) contracts.MaintenancePaymentReport {
	return contracts.MaintenancePaymentReport{
		ID:             v.ID,
		SocietyID:      v.SocietyID,
		BillID:         v.BillID,
		RequestID:      uuid.UUID(v.RequestID.Bytes),
		ReportedBy:     v.ReportedBy,
		Reference:      v.Reference,
		AmountPaise:    v.AmountPaise,
		PaymentDate:    v.PaymentDate.Time,
		Explanation:    v.Explanation,
		Fingerprint:    v.Fingerprint,
		Status:         v.Status,
		ResolutionNote: v.ResolutionNote,
		UpdatedBy:      v.UpdatedBy,
		UpdatedAt:      pgTimestamptzToTime(v.UpdatedAt),
		CreatedAt:      pgTimestamptzToTime(v.CreatedAt),
	}
}

func mapFinancialMaintenancePaymentRequest(v db.MaintenancePaymentRequest) contracts.MaintenancePaymentRequest {
	return contracts.MaintenancePaymentRequest{
		ID:              uuid.UUID(v.ID.Bytes),
		SocietyID:       v.SocietyID,
		BillID:          v.BillID,
		SettingsVersion: v.SettingsVersion,
		AmountPaise:     v.AmountPaise,
		Reference:       v.Reference,
		State:           v.State,
		CreatedBy:       v.CreatedBy,
		CreatedAt:       pgTimestamptzToTime(v.CreatedAt),
	}
}

func mapFinancialMaintenancePaymentSettingsVersion(v db.MaintenancePaymentSettingsVersion) contracts.MaintenancePaymentSettingsVersion {
	return contracts.MaintenancePaymentSettingsVersion{
		SocietyID: v.SocietyID,
		Version:   v.Version,
		Enabled:   v.Enabled,
		UpiID:     v.UpiID,
		PayeeName: v.PayeeName,
		CreatedBy: v.CreatedBy,
		CreatedAt: pgTimestamptzToTime(v.CreatedAt),
	}
}

func mapFinancialNotificationBacklogRow(v db.NotificationBacklogRow) contracts.NotificationBacklogRecord {
	return contracts.NotificationBacklogRecord{
		Pending:       v.Pending,
		OldestSeconds: v.OldestSeconds,
	}
}

func mapFinancialNotificationOutbox(v db.NotificationOutbox) contracts.NotificationOutbox {
	return contracts.NotificationOutbox{
		ID:               v.ID,
		UserID:           v.UserID,
		SocietyID:        v.SocietyID,
		FlatID:           v.FlatID,
		Audience:         v.Audience,
		EventKey:         v.EventKey,
		Payload:          v.Payload,
		Attempts:         v.Attempts,
		AvailableAt:      pgTimestamptzToTime(v.AvailableAt),
		LeaseToken:       uuid.UUID(v.LeaseToken.Bytes),
		InboxCompletedAt: pgTimestamptzToTime(v.InboxCompletedAt),
		PushCompletedAt:  pgTimestamptzToTime(v.PushCompletedAt),
		CompletedAt:      pgTimestamptzToTime(v.CompletedAt),
		LastError:        v.LastError,
		PushEnabled:      v.PushEnabled,
	}
}

func mapFinancialUPICollectionSummaryRow(v db.UPICollectionSummaryRow) contracts.UPICollectionSummaryRecord {
	return contracts.UPICollectionSummaryRecord{
		BilledCount:      v.BilledCount,
		BilledPaise:      v.BilledPaise,
		CollectedPaise:   v.CollectedPaise,
		OutstandingPaise: v.OutstandingPaise,
		OverduePaise:     v.OverduePaise,
		PendingClaims:    v.PendingClaims,
	}
}
