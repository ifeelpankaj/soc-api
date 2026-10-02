package maintenancesvc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go-server/internal/models"
	"go-server/internal/repositories/contracts"

	service "go-server/internal/services"
	"strings"
	"time"
)

func pageLimit(n int32) int32 {
	if n < 1 || n > 100 {
		return 25
	}
	return n
}
func validStatus(value string, allowed ...string) error {
	if value == "" {
		return nil
	}
	for _, v := range allowed {
		if v == value {
			return nil
		}
	}
	return invalid("Invalid status filter")
}
func (s *PaymentService) SettingsHistory(ctx context.Context, f models.UPIListFilter) (models.UPISettingsPage, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	result := models.UPISettingsPage{Items: []models.MaintenancePaymentSettings{}}
	repo := s.repo
	if err := paymentAdmin(ctx, repo, f.SocietyID, f.UserID); err != nil {
		return result, err
	}
	limit := pageLimit(f.Limit)
	rows, err := repo.ListUPISettingsVersions(ctx, contracts.ListUPISettingsVersionsInput{SocietyID: f.SocietyID, BeforeVersion: f.BeforeID, Limit: limit + 1})
	if err != nil {
		return result, err
	}
	if len(rows) > int(limit) {
		rows = rows[:limit]
		v := rows[len(rows)-1].Version
		result.NextCursor = &v
	}
	for _, r := range rows {
		result.Items = append(result.Items, settingsDTO(r))
	}
	return result, nil
}
func (s *PaymentService) Claims(ctx context.Context, f models.UPIListFilter) (models.UPIClaimsPage, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	result := models.UPIClaimsPage{Items: []models.UPIClaim{}}
	repo := s.repo
	user := int64(0)
	if f.Resident {
		user = f.UserID
		if user <= 0 {
			return result, paymentError(contracts.ErrNotFound)
		}
	} else if err := paymentAdmin(ctx, repo, f.SocietyID, f.UserID); err != nil {
		return result, err
	}
	if err := validStatus(f.Status, "pending", "verified", "rejected", "cancelled"); err != nil {
		return result, err
	}
	month, err := paymentListMonth(f)
	if err != nil {
		return result, err
	}
	limit := pageLimit(f.Limit)
	if f.Reference != "" {
		var err error
		f.Reference, err = NormalizeUPIReference(f.Reference)
		if err != nil {
			return result, err
		}
	}
	block, flatNumber := nullableText(f.Block), nullableText(f.FlatNumber)
	rows, err := repo.ListUPIClaims(ctx, contracts.ListUPIClaimsInput{
		BillID: f.BillID, FlatID: f.FlatID, BillingMonth: month, SocietyID: f.SocietyID, UserID: user,
		Status: f.Status, Reference: f.Reference, BeforeID: f.BeforeID, Block: block, FlatNumber: flatNumber,
		Search: f.Search, Limit: limit + 1,
	})
	if err != nil {
		return result, err
	}
	if len(rows) > int(limit) {
		rows = rows[:limit]
		id := rows[len(rows)-1].ID
		result.NextCursor = &id
	}
	for _, r := range rows {
		item, err := claimDTO(ctx, repo, r)
		if err != nil {
			return result, err
		}
		if f.Resident {
			item.ReviewedBy = nil
		} else {
			b, err := repo.GetPaymentBill(ctx, contracts.GetPaymentBillInput{SocietyID: f.SocietyID, ID: item.BillID})
			if err != nil {
				return result, err
			}
			var snapshot struct {
				Flat models.MaintenanceFlat `json:"flat"`
			}
			if err := json.Unmarshal(b.Snapshot, &snapshot); err != nil {
				return result, err
			}
			snapshot.Flat.BilledParty = nil
			item.BillNumber = b.BillNumber
			item.Flat = &snapshot.Flat
		}
		result.Items = append(result.Items, item)
	}
	return result, nil
}
func (s *PaymentService) Payments(ctx context.Context, f models.UPIListFilter) (models.UPIPaymentsPage, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	result := models.UPIPaymentsPage{Items: []models.UPIPayment{}}
	repo := s.repo
	user := int64(0)
	if f.Resident {
		user = f.UserID
		if user <= 0 {
			return result, paymentError(contracts.ErrNotFound)
		}
		if f.Reference != "" {
			return result, invalid("Reference search is restricted to admins")
		}
	} else if err := paymentAdmin(ctx, repo, f.SocietyID, f.UserID); err != nil {
		return result, err
	}
	if err := validStatus(f.Status, "verified", "reversed"); err != nil {
		return result, err
	}
	month, err := paymentListMonth(f)
	if err != nil {
		return result, err
	}
	limit := pageLimit(f.Limit)
	if f.Reference != "" {
		var err error
		f.Reference, err = NormalizeUPIReference(f.Reference)
		if err != nil {
			return result, err
		}
	}
	block, flatNumber := nullableText(f.Block), nullableText(f.FlatNumber)
	rows, err := repo.ListUPIPayments(ctx, contracts.ListUPIPaymentsInput{
		BillID: f.BillID, FlatID: f.FlatID, BillingMonth: month, SocietyID: f.SocietyID, UserID: user,
		Status: f.Status, Reference: f.Reference, BeforeID: f.BeforeID, Block: block, FlatNumber: flatNumber,
		Search: f.Search, Limit: limit + 1,
	})
	if err != nil {
		return result, err
	}
	if len(rows) > int(limit) {
		rows = rows[:limit]
		id := rows[len(rows)-1].ID
		result.NextCursor = &id
	}
	for _, r := range rows {
		item, err := paymentDTO(ctx, repo, r, f.UserID, !f.Resident)
		if err != nil {
			return result, err
		}
		result.Items = append(result.Items, item)
	}
	return result, nil
}
func (s *PaymentService) Payment(ctx context.Context, society, user, id int64, resident bool) (models.UPIPayment, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	repo := s.repo
	if !resident {
		if err := paymentAdmin(ctx, repo, society, user); err != nil {
			return models.UPIPayment{}, err
		}
	}
	p, err := repo.GetUPIPayment(ctx, contracts.GetUPIPaymentInput{SocietyID: society, ID: id})
	if err != nil {
		return models.UPIPayment{}, paymentError(err)
	}
	if resident {
		if err = paymentResident(ctx, repo, society, user, p.BillID); err != nil {
			return models.UPIPayment{}, err
		}
	}
	result, err := paymentDTO(ctx, repo, p, user, !resident)
	if err != nil {
		return result, err
	}
	if p.VerifiedBy > 0 {
		if names, ok := repo.(interface {
			UserName(context.Context, int64) (string, error)
		}); ok {
			result.VerifiedByName, err = names.UserName(ctx, p.VerifiedBy)
			if err != nil {
				return models.UPIPayment{}, err
			}
		}
	}
	return result, nil
}
func (s *PaymentService) ReferenceHistory(ctx context.Context, society, user int64, reference string) (models.UPIReferenceHistory, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	var err error
	reference, err = NormalizeUPIReference(reference)
	if err != nil {
		return models.UPIReferenceHistory{}, err
	}
	f := models.UPIListFilter{SocietyID: society, UserID: user, Reference: reference, Limit: 100}
	claims, err := s.Claims(ctx, f)
	if err != nil {
		return models.UPIReferenceHistory{}, err
	}
	payments, err := s.Payments(ctx, f)
	return models.UPIReferenceHistory{Claims: claims, Payments: payments}, err
}
func (s *PaymentService) Claim(ctx context.Context, society, user, id int64) (models.UPIClaimDetail, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	var result models.UPIClaimDetail
	repo := s.repo
	if err := paymentAdmin(ctx, repo, society, user); err != nil {
		return result, err
	}
	c, err := repo.GetUPIClaim(ctx, contracts.GetUPIClaimInput{SocietyID: society, ID: id})
	if err != nil {
		return result, paymentError(err)
	}
	result.UPIClaim, err = claimDTO(ctx, repo, c)
	if err != nil {
		return result, err
	}
	result.ReferenceHistory, err = s.ReferenceHistory(ctx, society, user, c.Reference)
	return result, err
}
func (s *PaymentService) Summary(ctx context.Context, society, user int64, month string, flatIDs ...int64) (models.UPICollectionSummary, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	repo := s.repo
	if err := paymentAdmin(ctx, repo, society, user); err != nil {
		return models.UPICollectionSummary{}, err
	}
	date := time.Time{}
	if month != "" {
		t, err := time.Parse("2006-01", month)
		if err != nil {
			return models.UPICollectionSummary{}, invalid("billing_month must be YYYY-MM")
		}
		date = t
	}
	flatID := int64(0)
	if len(flatIDs) > 0 {
		flatID = flatIDs[0]
	}
	if flatID < 0 {
		return models.UPICollectionSummary{}, invalid("flat_id must be positive")
	}
	r, err := repo.UPICollectionSummary(ctx, contracts.UPICollectionSummaryInput{SocietyID: society, Month: date, FlatID: flatID})
	return models.UPICollectionSummary{BilledCount: r.BilledCount, BilledPaise: r.BilledPaise, CollectedPaise: r.CollectedPaise, OutstandingPaise: r.OutstandingPaise, OverduePaise: r.OverduePaise, PendingClaims: r.PendingClaims}, err
}
func reportDTO(r contracts.MaintenancePaymentReport, billNumber string) models.UPIReport {
	return models.UPIReport{
		ID: r.ID, SocietyID: r.SocietyID, BillNumber: billNumber,
		UPICreateReport: models.UPICreateReport{BillID: r.BillID, PaymentRequestID: uuidText(r.RequestID), Reference: r.Reference, AmountPaise: r.AmountPaise, PaymentDate: dateText(r.PaymentDate), Explanation: r.Explanation},
		ReportedBy:      r.ReportedBy, Status: r.Status, ResolutionNote: r.ResolutionNote, UpdatedBy: r.UpdatedBy, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
}
func (s *PaymentService) CreateReport(ctx context.Context, society, user int64, key string, req models.UPICreateReport, admin bool) (models.UPIReport, error) {
	var err error
	req.Reference, err = NormalizeUPIReference(req.Reference)
	if err != nil {
		return models.UPIReport{}, err
	}
	req.Explanation = strings.TrimSpace(req.Explanation)
	if req.BillID <= 0 || req.AmountPaise <= 0 {
		return models.UPIReport{}, invalid("Positive bill ID and transfer amount are required")
	}
	if err = validateReason(req.Explanation); err != nil {
		return models.UPIReport{}, err
	}
	return paymentWrite(s, ctx, society, user, paymentTarget{bill: req.BillID}, admin, key, "report/create", req, func(ctx context.Context, repo contracts.FinancialStore, b contracts.MaintenanceBill) (models.UPIReport, error) {
		var zero models.UPIReport
		date, err := s.paymentDate(req.PaymentDate, b.Timezone)
		if err != nil {
			return zero, err
		}
		requestID := uuid.UUID{}
		if req.PaymentRequestID != "" {
			requestID, err = paymentUUID(req.PaymentRequestID)
			if err != nil {
				return zero, err
			}
			r, err := repo.GetUPIRequest(ctx, contracts.GetUPIRequestInput{SocietyID: society, ID: requestID})
			if err != nil {
				return zero, err
			}
			if r.BillID != b.ID {
				return zero, paymentError(contracts.ErrNotFound)
			}
		}
		fingerprint, err := digest(req)
		if err != nil {
			return zero, err
		}
		r, err := repo.InsertUPIReport(ctx, contracts.InsertUPIReportInput{SocietyID: society, BillID: b.ID, RequestID: requestID, ReportedBy: user, Reference: req.Reference, AmountPaise: req.AmountPaise, PaymentDate: date, Explanation: req.Explanation, Fingerprint: fingerprint})
		if errors.Is(err, contracts.ErrNotFound) {
			r, err = repo.GetUPIReportByFingerprint(ctx, contracts.GetUPIReportByFingerprintInput{SocietyID: society, ReportedBy: user, Fingerprint: fingerprint})
			return reportDTO(r, b.BillNumber), err
		}
		if err != nil {
			return zero, err
		}
		if err = paymentAudit(ctx, repo, society, user, b.ID, "report_created", fmt.Sprint(r.ID), req); err != nil {
			return zero, err
		}
		if err = paymentEvent(ctx, repo, society, b.ID, "maintenance_reconciliation_opened", fmt.Sprint(r.ID), "admin"); err != nil {
			return zero, err
		}
		return reportDTO(r, b.BillNumber), nil
	})
}
func (s *PaymentService) UpdateReport(ctx context.Context, society, user, id int64, key string, req models.UPIUpdateReport) (models.UPIReport, error) {
	if req.Status != "investigating" && req.Status != "resolved" {
		return models.UPIReport{}, invalid("Report status must be investigating or resolved")
	}
	req.ResolutionNote = strings.TrimSpace(req.ResolutionNote)
	if req.Status == "resolved" {
		if err := validateReason(req.ResolutionNote); err != nil {
			return models.UPIReport{}, err
		}
	} else if len(req.ResolutionNote) > 1000 {
		return models.UPIReport{}, invalid("Note exceeds 1000 bytes")
	}
	return paymentWrite(s, ctx, society, user, paymentTarget{report: id}, true, key, fmt.Sprintf("report/update/%d", id), req, func(ctx context.Context, repo contracts.FinancialStore, b contracts.MaintenanceBill) (models.UPIReport, error) {
		var zero models.UPIReport
		row, err := repo.GetUPIReport(ctx, contracts.GetUPIReportInput{SocietyID: society, ID: id})
		if err != nil {
			return zero, err
		}
		if row.MaintenancePaymentReport.Status == "resolved" {
			return zero, paymentConflict("REPORT_RESOLVED", "Resolved reports cannot be changed")
		}
		if row.MaintenancePaymentReport.Status == req.Status {
			return reportDTO(row.MaintenancePaymentReport, row.BillNumber), nil
		}
		updated, err := repo.UpdateUPIReport(ctx, contracts.UpdateUPIReportInput{SocietyID: society, ID: id, Status: req.Status, ResolutionNote: nullableText(req.ResolutionNote), UpdatedBy: &user})
		if err != nil {
			return zero, err
		}
		if err = paymentAudit(ctx, repo, society, user, b.ID, "report_"+req.Status, fmt.Sprint(id), req); err != nil {
			return zero, err
		}
		if err = paymentEvent(ctx, repo, society, b.ID, "maintenance_reconciliation_"+req.Status, fmt.Sprint(id), "resident"); err != nil {
			return zero, err
		}
		return reportDTO(updated, row.BillNumber), nil
	})
}
func (s *PaymentService) Reports(ctx context.Context, f models.UPIListFilter) (models.UPIReportsPage, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	result := models.UPIReportsPage{Items: []models.UPIReport{}}
	repo := s.repo
	user := int64(0)
	if f.Resident {
		user = f.UserID
		if user <= 0 {
			return result, paymentError(contracts.ErrNotFound)
		}
	} else if err := paymentAdmin(ctx, repo, f.SocietyID, f.UserID); err != nil {
		return result, err
	}
	if err := validStatus(f.Status, "open", "investigating", "resolved"); err != nil {
		return result, err
	}
	limit := pageLimit(f.Limit)
	rows, err := repo.ListUPIReports(ctx, contracts.ListUPIReportsInput{SocietyID: f.SocietyID, UserID: user, Status: f.Status, BeforeID: f.BeforeID, Limit: limit + 1})
	if err != nil {
		return result, err
	}
	if len(rows) > int(limit) {
		rows = rows[:limit]
		id := rows[len(rows)-1].MaintenancePaymentReport.ID
		result.NextCursor = &id
	}
	for _, r := range rows {
		item := reportDTO(r.MaintenancePaymentReport, r.BillNumber)
		if f.Resident {
			item.UpdatedBy = nil
		}
		result.Items = append(result.Items, item)
	}
	return result, nil
}
func (s *PaymentService) Report(ctx context.Context, society, user, id int64) (models.UPIReportDetail, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	var result models.UPIReportDetail
	repo := s.repo
	if err := paymentAdmin(ctx, repo, society, user); err != nil {
		return result, err
	}
	row, err := repo.GetUPIReport(ctx, contracts.GetUPIReportInput{SocietyID: society, ID: id})
	if err != nil {
		return result, paymentError(err)
	}
	result.UPIReport = reportDTO(row.MaintenancePaymentReport, row.BillNumber)
	result.ReferenceHistory, err = s.ReferenceHistory(ctx, society, user, row.MaintenancePaymentReport.Reference)
	return result, err
}
func (s *PaymentService) Audit(ctx context.Context, f models.UPIListFilter, bill int64) (models.UPIAuditPage, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	result := models.UPIAuditPage{Items: []models.UPIAuditEvent{}}
	repo := s.repo
	if err := paymentAdmin(ctx, repo, f.SocietyID, f.UserID); err != nil {
		return result, err
	}
	limit := pageLimit(f.Limit)
	rows, err := repo.ListUPIAudit(ctx, contracts.ListUPIAuditInput{SocietyID: f.SocietyID, BillID: bill, BeforeID: f.BeforeID, Limit: limit + 1})
	if err != nil {
		return result, err
	}
	if len(rows) > int(limit) {
		rows = rows[:limit]
		id := rows[len(rows)-1].ID
		result.NextCursor = &id
	}
	for _, r := range rows {
		var details map[string]any
		decoder := json.NewDecoder(bytes.NewReader(r.Details))
		decoder.UseNumber()
		if err = decoder.Decode(&details); err != nil {
			return result, err
		}
		result.Items = append(result.Items, models.UPIAuditEvent{ID: r.ID, BillID: r.BillID, ActorID: r.ActorID, Action: r.Action, EntityID: r.EntityID, Details: details, CreatedAt: r.CreatedAt})
	}
	return result, nil
}

func paymentListMonth(f models.UPIListFilter) (time.Time, error) {
	if f.BillID < 0 || f.FlatID < 0 {
		return time.Time{}, invalid("IDs must be positive")
	}
	// Narrowing filters do not grant access: resident queries still bind UserID
	// and recheck active membership in SQL. Never accept a user override here.
	if f.BillingMonth == "" {
		return time.Time{}, nil
	}
	t, err := time.Parse("2006-01", f.BillingMonth)
	if err != nil {
		return time.Time{}, invalid("billing_month must be YYYY-MM")
	}
	return t, nil
}
