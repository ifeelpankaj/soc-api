package maintenancesvc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"

	"go-server/internal/models"
	"go-server/internal/repositories/contracts"
	"strings"
	"unicode/utf8"
)

func claimDTO(ctx context.Context, repo contracts.FinancialStore, c contracts.MaintenancePaymentClaim) (models.UPIClaim, error) {
	r, err := repo.GetUPIRequest(ctx, contracts.GetUPIRequestInput{SocietyID: c.SocietyID, ID: c.RequestID})
	if err != nil {
		return models.UPIClaim{}, err
	}
	v, err := repo.GetUPISettingsVersion(ctx, contracts.GetUPISettingsVersionInput{SocietyID: c.SocietyID, Version: r.SettingsVersion})
	if err != nil {
		return models.UPIClaim{}, err
	}
	return models.UPIClaim{ID: c.ID, SocietyID: c.SocietyID, BillID: c.BillID, PaymentRequestID: uuidText(c.RequestID), SubmittedBy: c.SubmittedBy, Reference: c.Reference, PaymentDate: dateText(c.PaymentDate), Status: c.Status, AmountPaise: r.AmountPaise, SettingsVersion: r.SettingsVersion, Destination: settingsDTO(v).PaymentMethods.UPI, ReviewedBy: c.ReviewedBy, ReviewedAt: timePointer(c.ReviewedAt), Reason: c.Reason, CreatedAt: c.CreatedAt}, nil
}
func (s *PaymentService) SubmitClaim(ctx context.Context, society, user, bill int64, key string, req models.UPISubmitClaim) (models.UPIClaim, error) {
	var err error
	req.Reference, err = NormalizeUPIReference(req.Reference)
	if err != nil {
		return models.UPIClaim{}, err
	}
	requestID, err := paymentUUID(req.PaymentRequestID)
	if err != nil {
		return models.UPIClaim{}, err
	}
	return paymentWrite(s, ctx, society, user, paymentTarget{bill: bill}, false, key, fmt.Sprintf("claim/submit/%d", bill), req, func(ctx context.Context, repo contracts.FinancialStore, b contracts.MaintenanceBill) (models.UPIClaim, error) {
		var zero models.UPIClaim
		if err := unsettled(ctx, repo, b); err != nil {
			return zero, err
		}
		if err := noPending(ctx, repo, b); err != nil {
			return zero, err
		}
		r, err := repo.GetUPIRequest(ctx, contracts.GetUPIRequestInput{SocietyID: society, ID: requestID})
		if err != nil {
			return zero, err
		}
		if r.BillID != bill {
			return zero, paymentError(contracts.ErrNotFound)
		}
		if r.State == "closed" {
			return zero, paymentConflict("PAYMENT_REQUEST_CLOSED", "Report transfers against a closed request through payment-reports")
		}
		date, err := s.paymentDate(req.PaymentDate, b.Timezone)
		if err != nil {
			return zero, err
		}
		c, err := repo.InsertUPIClaim(ctx, contracts.InsertUPIClaimInput{SocietyID: society, BillID: bill, RequestID: requestID, SubmittedBy: user, Reference: req.Reference, PaymentDate: date})
		if err != nil {
			return zero, err
		}
		if err = repo.ReserveUPIReference(ctx, contracts.ReserveUPIReferenceInput{SocietyID: society, Reference: req.Reference, BillID: bill, ClaimID: &c.ID}); err != nil {
			return zero, err
		}
		entity := fmt.Sprint(c.ID)
		if err = paymentAudit(ctx, repo, society, user, bill, "claim_submitted", entity, req); err != nil {
			return zero, err
		}
		if err = paymentEvent(ctx, repo, society, bill, "maintenance_payment_submitted", entity, "admin"); err != nil {
			return zero, err
		}
		return claimDTO(ctx, repo, c)
	})
}
func (s *PaymentService) CloseClaim(ctx context.Context, society, user, id int64, key string, req models.UPIReason, admin bool) (models.UPIClaim, error) {
	action := "cancelled"
	if admin {
		action = "rejected"
		if err := validateReason(req.Reason); err != nil {
			return models.UPIClaim{}, err
		}
	}
	return paymentWrite(s, ctx, society, user, paymentTarget{claim: id}, admin, key, fmt.Sprintf("claim/%s/%d", action, id), req, func(ctx context.Context, repo contracts.FinancialStore, b contracts.MaintenanceBill) (models.UPIClaim, error) {
		var zero models.UPIClaim
		c, err := repo.GetUPIClaim(ctx, contracts.GetUPIClaimInput{SocietyID: society, ID: id})
		if err != nil {
			return zero, err
		}
		if !admin && c.SubmittedBy != user {
			return zero, paymentError(contracts.ErrNotFound)
		}
		if c.Status != "pending" {
			return zero, paymentConflict("CLAIM_ALREADY_REVIEWED", "Only a pending claim can be cancelled or rejected")
		}
		c, err = repo.ReviewUPIClaim(ctx, contracts.ReviewUPIClaimInput{SocietyID: society, ID: id, Status: action, ReviewedBy: &user, Reason: nullableText(strings.TrimSpace(req.Reason))})
		if err != nil {
			return zero, err
		}
		if err = repo.ReleaseUPIClaimReference(ctx, contracts.ReleaseUPIClaimReferenceInput{SocietyID: society, ClaimID: &id}); err != nil {
			return zero, err
		}
		if err = paymentAudit(ctx, repo, society, user, b.ID, "claim_"+action, fmt.Sprint(id), req); err != nil {
			return zero, err
		}
		audience := "resident"
		if !admin {
			audience = "admin"
		}
		if err = paymentEvent(ctx, repo, society, b.ID, "maintenance_payment_"+action, fmt.Sprint(id), audience); err != nil {
			return zero, err
		}
		return claimDTO(ctx, repo, c)
	})
}
func (s *PaymentService) validateCredit(ctx context.Context, repo contracts.FinancialStore, b contracts.MaintenanceBill, req models.UPIVerifyCredit) (time.Time, error) {
	if strings.TrimSpace(req.EvidenceReference) == "" {
		return time.Time{}, invalid("A bank statement or evidence reference is required")
	}
	if !req.BankCreditConfirmed {
		return time.Time{}, invalid("Explicit bank_credit_confirmed is required")
	}
	if req.AmountPaise != b.TotalPaise {
		return time.Time{}, invalid("Verified bank credit must exactly equal the full bill amount")
	}
	v, err := repo.GetUPISettingsVersion(ctx, contracts.GetUPISettingsVersionInput{SocietyID: b.SocietyID, Version: req.SettingsVersion})
	if err != nil {
		return time.Time{}, err
	}
	if !upiPattern.MatchString(req.UPIID) || !strings.EqualFold(v.UpiID, strings.TrimSpace(req.UPIID)) {
		return time.Time{}, invalid("Confirm the UPI destination belonging to the selected historical settings version")
	}
	reversed, err := repo.HasReversedUPIReference(ctx, contracts.HasReversedUPIReferenceInput{SocietyID: b.SocietyID, Reference: req.Reference})
	if err != nil {
		return time.Time{}, err
	}
	if reversed && !req.ReversalHistoryAcknowledged {
		return time.Time{}, paymentConflict("REVERSAL_HISTORY_ACKNOWLEDGMENT_REQUIRED", "Review prior payment/reversal history and explicitly acknowledge it")
	}
	if reversed {
		if err := validateReason(req.CorrectionReason); err != nil {
			return time.Time{}, invalid("A correction reason is required to reuse a reversed bank credit")
		}
	}
	return s.paymentDate(req.CreditDate, b.Timezone)
}
func paymentDTO(ctx context.Context, repo contracts.FinancialStore, p contracts.MaintenancePayment, user int64, admin bool) (models.UPIPayment, error) {
	v, err := repo.GetUPISettingsVersion(ctx, contracts.GetUPISettingsVersionInput{SocietyID: p.SocietyID, Version: p.SettingsVersion})
	if err != nil {
		return models.UPIPayment{}, err
	}
	b, err := repo.GetPaymentBill(ctx, contracts.GetPaymentBillInput{SocietyID: p.SocietyID, ID: p.BillID})
	if err != nil {
		return models.UPIPayment{}, err
	}
	result := models.UPIPayment{ID: p.ID, SocietyID: p.SocietyID, BillID: p.BillID, BillNumber: b.BillNumber, AmountPaise: p.AmountPaise, Currency: "INR", CreditDate: dateText(p.CreditDate), ReceiptNumber: "MPR-" + uuidText(p.ReceiptNumber), SettingsVersion: p.SettingsVersion, Destination: settingsDTO(v).PaymentMethods.UPI, Status: p.Status, VerifiedAt: p.VerifiedAt, ReversedAt: timePointer(p.ReversedAt)}
	if admin || (p.PayerID != nil && *p.PayerID == user) {
		result.ClaimID = p.ClaimID
		result.PayerID = p.PayerID
		result.Reference = &p.Reference
	}
	if admin {
		result.EvidenceReference = p.EvidenceReference
		var snapshot struct {
			Flat models.MaintenanceFlat `json:"flat"`
		}
		if err := json.Unmarshal(b.Snapshot, &snapshot); err != nil {
			return models.UPIPayment{}, err
		}
		snapshot.Flat.BilledParty = nil
		result.Flat = &snapshot.Flat
		result.VerifiedBy = &p.VerifiedBy
		result.ReversedBy = p.ReversedBy
		result.ReversalReason = p.ReversalReason
	}
	return result, nil
}
func (s *PaymentService) recordCredit(ctx context.Context, repo contracts.FinancialStore, b contracts.MaintenanceBill, user int64, req models.UPIVerifyCredit, claim *contracts.MaintenancePaymentClaim) (models.UPIPayment, error) {
	var zero models.UPIPayment
	if err := unsettled(ctx, repo, b); err != nil {
		return zero, err
	}
	date, err := s.validateCredit(ctx, repo, b, req)
	if err != nil {
		return zero, err
	}
	var claimID, payer *int64
	if claim != nil {
		if claim.Status != "pending" {
			return zero, paymentConflict("CLAIM_ALREADY_REVIEWED", "Claim is not pending")
		}
		r, err := repo.GetUPIRequest(ctx, contracts.GetUPIRequestInput{SocietyID: b.SocietyID, ID: claim.RequestID})
		if err != nil {
			return zero, err
		}
		if claim.Reference != req.Reference || r.SettingsVersion != req.SettingsVersion {
			return zero, invalid("Verified reference and destination version must match the claim snapshot")
		}
		claimID = &claim.ID
		payer = &claim.SubmittedBy
		if err = repo.ReleaseUPIClaimReference(ctx, contracts.ReleaseUPIClaimReferenceInput{SocietyID: b.SocietyID, ClaimID: claimID}); err != nil {
			return zero, err
		}
		if _, err = repo.ReviewUPIClaim(ctx, contracts.ReviewUPIClaimInput{SocietyID: b.SocietyID, ID: claim.ID, Status: "verified", ReviewedBy: &user}); err != nil {
			return zero, err
		}
	}
	p, err := repo.InsertUPIPayment(ctx, contracts.InsertUPIPaymentInput{SocietyID: b.SocietyID, BillID: b.ID, ClaimID: claimID, PayerID: payer, SettingsVersion: req.SettingsVersion, AmountPaise: b.TotalPaise, Reference: req.Reference, CreditDate: date, ReceiptNumber: uuid.New(), VerifiedBy: user, EvidenceReference: req.EvidenceReference})
	if err != nil {
		return zero, err
	}
	if err = repo.ReserveUPIReference(ctx, contracts.ReserveUPIReferenceInput{SocietyID: b.SocietyID, Reference: req.Reference, BillID: b.ID, PaymentID: &p.ID}); err != nil {
		return zero, err
	}
	if err = repo.InsertUPILedger(ctx, contracts.InsertUPILedgerInput{SocietyID: b.SocietyID, BillID: b.ID, PaymentID: p.ID, Kind: "collection", AmountPaise: b.TotalPaise, ActorID: user}); err != nil {
		return zero, err
	}
	if err = repo.CloseUPIRequests(ctx, contracts.CloseUPIRequestsInput{SocietyID: b.SocietyID, BillID: b.ID}); err != nil {
		return zero, err
	}
	if err = paymentAudit(ctx, repo, b.SocietyID, user, b.ID, "payment_verified", fmt.Sprint(p.ID), req); err != nil {
		return zero, err
	}
	if err = paymentEvent(ctx, repo, b.SocietyID, b.ID, "maintenance_payment_verified", fmt.Sprint(p.ID), "resident"); err != nil {
		return zero, err
	}
	return paymentDTO(ctx, repo, p, user, true)
}
func (s *PaymentService) Verify(ctx context.Context, society, user, claimID int64, key string, req models.UPIVerifyCredit) (models.UPIPayment, error) {
	req.EvidenceReference = strings.TrimSpace(req.EvidenceReference)
	if utf8.RuneCountInString(req.EvidenceReference) > 500 {
		return models.UPIPayment{}, invalid("Evidence reference must be at most 500 characters")
	}
	var err error
	req.Reference, err = NormalizeUPIReference(req.Reference)
	if err != nil {
		return models.UPIPayment{}, err
	}
	return paymentWrite(s, ctx, society, user, paymentTarget{claim: claimID}, true, key, fmt.Sprintf("claim/verify/%d", claimID), req, func(ctx context.Context, repo contracts.FinancialStore, b contracts.MaintenanceBill) (models.UPIPayment, error) {
		c, err := repo.GetUPIClaim(ctx, contracts.GetUPIClaimInput{SocietyID: society, ID: claimID})
		if err != nil {
			return models.UPIPayment{}, err
		}
		return s.recordCredit(ctx, repo, b, user, req, &c)
	})
}
func (s *PaymentService) Record(ctx context.Context, society, user int64, key string, req models.UPIDirectPayment) (models.UPIPayment, error) {
	req.EvidenceReference = strings.TrimSpace(req.EvidenceReference)
	if utf8.RuneCountInString(req.EvidenceReference) > 500 {
		return models.UPIPayment{}, invalid("Evidence reference must be at most 500 characters")
	}
	var err error
	req.Reference, err = NormalizeUPIReference(req.Reference)
	if err != nil {
		return models.UPIPayment{}, err
	}
	if req.BillID <= 0 {
		return models.UPIPayment{}, invalid("Bill ID is required")
	}
	return paymentWrite(s, ctx, society, user, paymentTarget{bill: req.BillID}, true, key, "payment/record", req, func(ctx context.Context, repo contracts.FinancialStore, b contracts.MaintenanceBill) (models.UPIPayment, error) {
		var zero models.UPIPayment
		c, err := repo.GetPendingUPIClaim(ctx, contracts.GetPendingUPIClaimInput{SocietyID: society, BillID: b.ID})
		if err == nil {
			if req.PendingClaimID != c.ID {
				return zero, paymentConflict("PENDING_CLAIM_RESOLUTION_REQUIRED", "Select and explicitly resolve the pending claim")
			}
			if req.ClaimResolution == "verify" {
				return s.recordCredit(ctx, repo, b, user, req.UPIVerifyCredit, &c)
			}
			if req.ClaimResolution != "reject" || req.Reference == c.Reference {
				return zero, invalid("Use verify for a matching credit; reject requires a distinct reference")
			}
			if err := validateReason(req.Reason); err != nil {
				return zero, err
			}
			if _, err = repo.ReviewUPIClaim(ctx, contracts.ReviewUPIClaimInput{SocietyID: society, ID: c.ID, Status: "rejected", ReviewedBy: &user, Reason: &req.Reason}); err != nil {
				return zero, err
			}
			if err = repo.ReleaseUPIClaimReference(ctx, contracts.ReleaseUPIClaimReferenceInput{SocietyID: society, ClaimID: &c.ID}); err != nil {
				return zero, err
			}
			if err = paymentAudit(ctx, repo, society, user, b.ID, "claim_rejected", fmt.Sprint(c.ID), models.UPIReason{Reason: req.Reason}); err != nil {
				return zero, err
			}
			if err = paymentEvent(ctx, repo, society, b.ID, "maintenance_payment_rejected", fmt.Sprint(c.ID), "resident"); err != nil {
				return zero, err
			}
		} else if !errors.Is(err, contracts.ErrNotFound) {
			return zero, err
		} else if req.PendingClaimID != 0 || req.ClaimResolution != "" {
			return zero, paymentConflict("PENDING_CLAIM_CHANGED", "No pending claim matches the requested resolution")
		}
		return s.recordCredit(ctx, repo, b, user, req.UPIVerifyCredit, nil)
	})
}
func (s *PaymentService) Reverse(ctx context.Context, society, user, id int64, key string, req models.UPIReason) (models.UPIPayment, error) {
	if err := validateReason(req.Reason); err != nil {
		return models.UPIPayment{}, err
	}
	return paymentWrite(s, ctx, society, user, paymentTarget{payment: id}, true, key, fmt.Sprintf("payment/reverse/%d", id), req, func(ctx context.Context, repo contracts.FinancialStore, b contracts.MaintenanceBill) (models.UPIPayment, error) {
		var zero models.UPIPayment
		p, err := repo.GetUPIPayment(ctx, contracts.GetUPIPaymentInput{SocietyID: society, ID: id})
		if err != nil {
			return zero, err
		}
		if p.Status != "verified" {
			return zero, paymentConflict("PAYMENT_ALREADY_REVERSED", "Payment has already been reversed")
		}
		p, err = repo.ReverseUPIPayment(ctx, contracts.ReverseUPIPaymentInput{SocietyID: society, ID: id, ReversedBy: &user, ReversalReason: &req.Reason})
		if err != nil {
			return zero, err
		}
		if err = repo.InsertUPILedger(ctx, contracts.InsertUPILedgerInput{SocietyID: society, BillID: b.ID, PaymentID: id, Kind: "reversal", AmountPaise: -p.AmountPaise, ActorID: user}); err != nil {
			return zero, err
		}
		if err = repo.ReleaseUPIPaymentReference(ctx, contracts.ReleaseUPIPaymentReferenceInput{SocietyID: society, PaymentID: &id}); err != nil {
			return zero, err
		}
		if err = repo.CloseUPIRequests(ctx, contracts.CloseUPIRequestsInput{SocietyID: society, BillID: b.ID}); err != nil {
			return zero, err
		}
		if err = paymentAudit(ctx, repo, society, user, b.ID, "payment_reversed", fmt.Sprint(id), req); err != nil {
			return zero, err
		}
		if err = paymentEvent(ctx, repo, society, b.ID, "maintenance_payment_reversed", fmt.Sprint(id), "resident"); err != nil {
			return zero, err
		}
		return paymentDTO(ctx, repo, p, user, true)
	})
}
