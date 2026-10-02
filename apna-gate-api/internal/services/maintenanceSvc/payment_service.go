package maintenancesvc

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"go-server/internal/models"
	"go-server/internal/repositories/contracts"

	"github.com/google/uuid"

	qrcode "github.com/skip2/go-qrcode"
)

type PaymentService struct {
	repo        contracts.PaymentStore
	operational OperationalGuard
	now         func() time.Time
}

func NewPaymentService(repo contracts.PaymentStore, guard OperationalGuard) *PaymentService {
	if repo == nil || guard == nil {
		panic("maintenanceSvc: required dependency is nil")
	}
	return &PaymentService{repo: repo, operational: guard, now: time.Now}
}

func paymentUUID(s string) (uuid.UUID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return uuid.UUID{}, invalid("Invalid payment request ID")
	}
	return u, nil
}
func uuidText(v uuid.UUID) string {
	if v == uuid.Nil {
		return ""
	}
	return v.String()
}
func dateText(v time.Time) string { return v.Format("2006-01-02") }
func timePointer(v time.Time) *time.Time {
	if v.IsZero() {
		return nil
	}
	return &v
}
func nullableText(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
func digest(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}

var upiPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{1,127}@[A-Za-z0-9][A-Za-z0-9.-]{1,63}$`)
var referencePattern = regexp.MustCompile(`^[A-Z0-9/-]{6,64}$`)

func NormalizeUPIReference(s string) (string, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	if !referencePattern.MatchString(s) {
		return "", invalid("Reference must be 6–64 letters, digits, slashes or hyphens")
	}
	return s, nil
}
func validateReason(s string) error {
	if strings.TrimSpace(s) == "" || len(s) > 1000 {
		return invalid("A reason of 1–1000 bytes is required")
	}
	return nil
}
func (s *PaymentService) paymentDate(value, zone string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, invalid("Payment date must be YYYY-MM-DD")
	}
	loc, err := time.LoadLocation(zone)
	if err != nil {
		return time.Time{}, err
	}
	if value > s.now().In(loc).Format("2006-01-02") {
		return time.Time{}, invalid("Payment date cannot be in the future")
	}
	return t, nil
}
func paymentAdmin(ctx context.Context, repo contracts.FinancialAccessRepository, society, user int64) error {
	ok, err := repo.MaintenanceAdmin(ctx, contracts.MaintenanceAdminInput{SocietyID: society, UserID: user})
	if err != nil {
		return err
	}
	if !ok {
		return models.NewAppError("PAYMENT_FORBIDDEN", "Active society owner or admin access required", 403, nil)
	}
	return nil
}
func paymentResident(ctx context.Context, repo contracts.FinancialAccessRepository, society, user, bill int64) error {
	ok, err := repo.MaintenanceResidentAccess(ctx, contracts.MaintenanceResidentAccessInput{SocietyID: society, ID: bill, UserID: user})
	if err != nil {
		return err
	}
	if !ok {
		return paymentError(contracts.ErrNotFound)
	}
	return nil
}

type paymentTarget struct {
	bill, claim, payment, report int64
	request                      string
}

func (t paymentTarget) resolve(ctx context.Context, repo contracts.FinancialStore, society int64) (int64, error) {
	if t.claim > 0 {
		r, e := repo.GetUPIClaim(ctx, contracts.GetUPIClaimInput{SocietyID: society, ID: t.claim})
		return r.BillID, e
	}
	if t.payment > 0 {
		r, e := repo.GetUPIPayment(ctx, contracts.GetUPIPaymentInput{SocietyID: society, ID: t.payment})
		return r.BillID, e
	}
	if t.report > 0 {
		r, e := repo.GetUPIReport(ctx, contracts.GetUPIReportInput{SocietyID: society, ID: t.report})
		return r.MaintenancePaymentReport.BillID, e
	}
	if t.request != "" {
		id, e := paymentUUID(t.request)
		if e != nil {
			return 0, e
		}
		r, e := repo.GetUPIRequest(ctx, contracts.GetUPIRequestInput{SocietyID: society, ID: id})
		return r.BillID, e
	}
	return t.bill, nil
}

// All financial commands use society lock -> bill row lock -> access recheck ->
// idempotency lookup -> writes. A replay never bypasses current authorization.
func paymentWrite[T any](s *PaymentService, ctx context.Context, society, user int64, t paymentTarget, admin bool, key, operation string, input any, fn func(context.Context, contracts.FinancialStore, contracts.MaintenanceBill) (T, error)) (T, error) {
	var result T
	if operation != "request" && (key == "" || len(key) > 128 || strings.TrimSpace(key) != key) {
		return result, invalid("Idempotency-Key is required (1–128 bytes, no surrounding whitespace)")
	}
	hash, err := digest(input)
	if err != nil {
		return result, err
	}
	ctx, cancel := context.WithTimeout(ctx, FinancialTimeout)
	defer cancel()
	err = s.repo.WithTransaction(ctx, func(ctx context.Context) error {
		repo := s.repo
		if err := repo.MaintenanceLock(ctx, society); err != nil {
			return err
		}
		if admin {
			if err := paymentAdmin(ctx, repo, society, user); err != nil {
				return err
			}
		}
		if err := s.operational.EnsureSocietyOperational(ctx, society); err != nil {
			return err
		}
		billID, err := t.resolve(ctx, repo, society)
		if err != nil {
			return err
		}
		var bill contracts.MaintenanceBill
		if billID > 0 {
			bill, err = repo.LockMaintenancePaymentBill(ctx, contracts.LockMaintenancePaymentBillInput{SocietyID: society, ID: billID})
			if err != nil {
				return err
			}
			if !admin {
				if err = paymentResident(ctx, repo, society, user, billID); err != nil {
					return err
				}
			}
		} else if !admin {
			return invalid("Bill ID is required")
		}
		if operation != "request" {
			prior, err := repo.GetUPIIdempotency(ctx, contracts.GetUPIIdempotencyInput{SocietyID: society, ActorID: user, Operation: operation, Key: key})
			if err == nil {
				if prior.RequestHash != hash {
					return paymentConflict("IDEMPOTENCY_CONFLICT", "Idempotency key was used with different request contents")
				}
				return json.Unmarshal(prior.Response, &result)
			}
			if !errors.Is(err, contracts.ErrNotFound) {
				return err
			}
		}
		result, err = fn(ctx, repo, bill)
		if err != nil {
			return err
		}
		if operation != "request" {
			response, err := json.Marshal(result)
			if err != nil {
				return err
			}
			return repo.InsertUPIIdempotency(ctx, contracts.InsertUPIIdempotencyInput{SocietyID: society, ActorID: user, Operation: operation, Key: key, RequestHash: hash, Response: response})
		}
		return nil
	})
	observeFinancial(operation, err)
	return result, paymentError(err)
}
func paymentAudit(ctx context.Context, repo contracts.FinancialEventsRepository, society, user, bill int64, action, entity string, details any) error {
	b, err := json.Marshal(details)
	if err != nil {
		return err
	}
	var billID *int64
	if bill > 0 {
		billID = &bill
	}
	return repo.InsertUPIAudit(ctx, contracts.InsertUPIAuditInput{SocietyID: society, BillID: billID, ActorID: user, Action: action, EntityID: entity, Details: b})
}
func paymentEvent(ctx context.Context, repo contracts.FinancialEventsRepository, society, bill int64, kind, entity, audience string) error {
	data, _ := json.Marshal(map[string]string{"entity_id": entity})
	return repo.EnqueueUPIEvent(ctx, contracts.EnqueueUPIEventInput{SocietyID: society, BillID: bill, EventType: kind, EventKey: kind + ":" + entity, Audience: audience, EventData: data})
}
func settingsDTO(v contracts.MaintenancePaymentSettingsVersion) models.MaintenancePaymentSettings {
	return models.MaintenancePaymentSettings{Version: v.Version, UpdatedBy: v.CreatedBy, UpdatedAt: v.CreatedAt, PaymentMethods: models.MaintenancePaymentMethods{UPI: models.UPIConfig{Enabled: v.Enabled, UPIID: v.UpiID, PayeeName: v.PayeeName}}}
}
func (s *PaymentService) Settings(ctx context.Context, society, user int64) (models.MaintenancePaymentSettings, error) {
	repo := s.repo
	if err := paymentAdmin(ctx, repo, society, user); err != nil {
		return models.MaintenancePaymentSettings{}, err
	}
	r, err := repo.GetUPISettings(ctx, society)
	if errors.Is(err, contracts.ErrNotFound) {
		return models.MaintenancePaymentSettings{}, nil
	}
	return settingsDTO(r), err
}
func (s *PaymentService) SaveSettings(ctx context.Context, society, user int64, key string, v models.MaintenancePaymentSettings) (models.MaintenancePaymentSettings, error) {
	cfg := v.PaymentMethods.UPI
	cfg.UPIID = strings.TrimSpace(cfg.UPIID)
	cfg.PayeeName = strings.TrimSpace(cfg.PayeeName)
	if (cfg.Enabled || cfg.UPIID != "") && !upiPattern.MatchString(cfg.UPIID) {
		return v, invalid("Invalid merchant UPI ID")
	}
	if len(cfg.PayeeName) > 120 || (cfg.Enabled && cfg.PayeeName == "") {
		return v, invalid("Payee name is required when UPI is enabled (maximum 120 bytes)")
	}
	return paymentWrite(s, ctx, society, user, paymentTarget{}, true, key, "settings", cfg, func(ctx context.Context, repo contracts.FinancialStore, _ contracts.MaintenanceBill) (models.MaintenancePaymentSettings, error) {
		current, err := repo.GetUPISettings(ctx, society)
		version := int64(1)
		if err == nil {
			version = current.Version + 1
			if current.Enabled == cfg.Enabled && current.UpiID == cfg.UPIID && current.PayeeName == cfg.PayeeName {
				return settingsDTO(current), nil
			}
		} else if !errors.Is(err, contracts.ErrNotFound) {
			return v, err
		}
		r, err := repo.InsertUPISettingsVersion(ctx, contracts.InsertUPISettingsVersionInput{SocietyID: society, Version: version, Enabled: cfg.Enabled, UpiID: cfg.UPIID, PayeeName: cfg.PayeeName, CreatedBy: user})
		if err != nil {
			return v, err
		}
		if err = repo.PointUPISettings(ctx, contracts.PointUPISettingsInput{SocietyID: society, Version: version}); err != nil {
			return v, err
		}
		if err = repo.SupersedeUPIRequests(ctx, society); err != nil {
			return v, err
		}
		result := settingsDTO(r)
		err = paymentAudit(ctx, repo, society, user, 0, "settings_updated", fmt.Sprint(version), result)
		return result, err
	})
}
func unsettled(ctx context.Context, repo contracts.CollectionsRepository, b contracts.MaintenanceBill) error {
	if b.TotalPaise <= 0 {
		return paymentConflict("BILL_NO_OUTSTANDING", "Bill has no outstanding amount to collect")
	}
	_, err := repo.GetActiveUPIPayment(ctx, contracts.GetActiveUPIPaymentInput{SocietyID: b.SocietyID, BillID: b.ID})
	if err == nil {
		return paymentConflict("BILL_ALREADY_PAID", "Bill is settled; report an additional transfer through /maintenance/my/payment-reports")
	}
	if !errors.Is(err, contracts.ErrNotFound) {
		return err
	}
	return nil
}
func noPending(ctx context.Context, repo contracts.PaymentClaimsRepository, b contracts.MaintenanceBill) error {
	_, err := repo.GetPendingUPIClaim(ctx, contracts.GetPendingUPIClaimInput{SocietyID: b.SocietyID, BillID: b.ID})
	if err == nil {
		return paymentConflict("PAYMENT_CLAIM_PENDING", "A payment claim is awaiting verification")
	}
	if !errors.Is(err, contracts.ErrNotFound) {
		return err
	}
	return nil
}
func UPIURI(config models.UPIConfig, amount int64, reference, note string) string {
	values := url.Values{"pa": {config.UPIID}, "pn": {config.PayeeName}, "am": {fmt.Sprintf("%d.%02d", amount/100, amount%100)}, "cu": {"INR"}, "tr": {reference}, "tn": {note}}
	return "upi://pay?" + values.Encode()
}
func requestDTO(ctx context.Context, repo contracts.FinancialStore, r contracts.MaintenancePaymentRequest, b contracts.MaintenanceBill) (models.UPIPaymentRequest, error) {
	v, err := repo.GetUPISettingsVersion(ctx, contracts.GetUPISettingsVersionInput{SocietyID: r.SocietyID, Version: r.SettingsVersion})
	if err != nil {
		return models.UPIPaymentRequest{}, err
	}
	cfg := settingsDTO(v).PaymentMethods.UPI
	id := uuidText(r.ID)
	return models.UPIPaymentRequest{ID: id, BillID: r.BillID, Reference: r.Reference, State: r.State, AmountPaise: r.AmountPaise, SettingsVersion: r.SettingsVersion, Destination: cfg, URI: UPIURI(cfg, r.AmountPaise, r.Reference, b.BillNumber), QRURL: fmt.Sprintf("/api/v1/societies/%d/maintenance/my/payment-requests/%s/qr", r.SocietyID, id)}, nil
}
func (s *PaymentService) Request(ctx context.Context, society, user, bill int64) (models.UPIPaymentRequest, error) {
	return s.requestWithKey(ctx, society, user, bill, "", "request")
}
func (s *PaymentService) RequestCommand(ctx context.Context, society, user, bill int64, key string) (models.UPIPaymentRequest, error) {
	return s.requestWithKey(ctx, society, user, bill, key, fmt.Sprintf("request/create/%d", bill))
}
func (s *PaymentService) requestWithKey(ctx context.Context, society, user, bill int64, key, operation string) (models.UPIPaymentRequest, error) {
	return paymentWrite(s, ctx, society, user, paymentTarget{bill: bill}, false, key, operation, bill, func(ctx context.Context, repo contracts.FinancialStore, b contracts.MaintenanceBill) (models.UPIPaymentRequest, error) {
		var zero models.UPIPaymentRequest
		if err := unsettled(ctx, repo, b); err != nil {
			return zero, err
		}
		if err := noPending(ctx, repo, b); err != nil {
			return zero, err
		}
		settings, err := repo.GetUPISettings(ctx, society)
		if errors.Is(err, contracts.ErrNotFound) || (err == nil && !settings.Enabled) {
			return zero, paymentConflict("UPI_DISABLED", "Society UPI payments are disabled")
		}
		if err != nil {
			return zero, err
		}
		r, err := repo.GetActiveUPIRequest(ctx, contracts.GetActiveUPIRequestInput{SocietyID: society, BillID: bill})
		if err == nil {
			return requestDTO(ctx, repo, r, b)
		}
		if !errors.Is(err, contracts.ErrNotFound) {
			return zero, err
		}
		id := uuid.New()
		r, err = repo.InsertUPIRequest(ctx, contracts.InsertUPIRequestInput{ID: id, SocietyID: society, BillID: bill, SettingsVersion: settings.Version, AmountPaise: b.TotalPaise, Reference: strings.ReplaceAll(id.String(), "-", ""), CreatedBy: user})
		if err != nil {
			return zero, err
		}
		if err = paymentAudit(ctx, repo, society, user, bill, "request_created", id.String(), map[string]any{"settings_version": settings.Version}); err != nil {
			return zero, err
		}
		return requestDTO(ctx, repo, r, b)
	})
}
func (s *PaymentService) QR(ctx context.Context, society, user int64, id string) ([]byte, error) {
	// Use the same lock order so settings changes and settlement cannot race the
	// authorization decision. A previously downloaded QR cannot be revoked.
	result, err := paymentWrite(s, ctx, society, user, paymentTarget{request: id}, false, "", "request", id, func(ctx context.Context, repo contracts.FinancialStore, b contracts.MaintenanceBill) (models.UPIPaymentRequest, error) {
		var zero models.UPIPaymentRequest
		uid, err := paymentUUID(id)
		if err != nil {
			return zero, err
		}
		r, err := repo.GetUPIRequest(ctx, contracts.GetUPIRequestInput{SocietyID: society, ID: uid})
		if err != nil {
			return zero, err
		}
		settings, err := repo.GetUPISettings(ctx, society)
		if err != nil {
			return zero, err
		}
		if !settings.Enabled || settings.Version != r.SettingsVersion || r.State != "active" {
			return zero, paymentConflict("PAYMENT_REQUEST_OBSOLETE", "Request is closed, superseded or disabled")
		}
		if err = unsettled(ctx, repo, b); err != nil {
			return zero, err
		}
		if err = noPending(ctx, repo, b); err != nil {
			return zero, err
		}
		return requestDTO(ctx, repo, r, b)
	})
	if err != nil {
		return nil, err
	}
	return qrcode.Encode(result.URI, qrcode.Medium, 384)
}
