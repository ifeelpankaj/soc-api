package maintenancesvc

import (
	"context"
	"errors"
	"fmt"
	"go-server/internal/models"
	repository "go-server/internal/repositories/contracts"
	service "go-server/internal/services"
	"math"
	"math/big"
	"time"
)

type Store interface {
	Locked(context.Context, int64, func(context.Context) error) error
	Admin(context.Context, int64, int64) (bool, error)
	Settings(context.Context, int64) (models.MaintenanceSettings, error)
	SaveSettings(context.Context, int64, int64, models.MaintenanceSettings) error
	Flats(context.Context, int64, []string) ([]models.MaintenanceFlat, error)
	Run(context.Context, int64, string) (*models.MaintenanceRunResult, error)
	RunTerms(context.Context, int64, string) (models.MaintenanceSettings, []int64, error)
	Issue(context.Context, int64, int64, string, models.MaintenanceSettings, []models.MaintenanceBill) (models.MaintenanceRunResult, error)
	IssueMissing(context.Context, int64, int64, models.MaintenanceRunResult, []models.MaintenanceBill) (models.MaintenanceRunResult, error)
	Bills(context.Context, models.MaintenanceBillFilter) ([]models.MaintenanceBill, error)
	Societies(context.Context) ([]int64, error)
	ClaimDelivery(context.Context) (*repository.MaintenanceDelivery, error)
	FinishDelivery(context.Context, *repository.MaintenanceDelivery, error) error
}
type OperationalGuard interface {
	EnsureSocietyOperational(context.Context, int64) error
}
type PushSender interface {
	SendToUser(context.Context, int64, models.NotificationPayload) error
}
type Inbox interface {
	Create(context.Context, models.NotificationCreate) (*models.Notification, error)
}
type Service struct {
	store       Store
	operational OperationalGuard
	statuses    []string
	now         func() time.Time
	inbox       Inbox
	push        PushSender
}

func New(store Store, guard OperationalGuard, statuses []string, inbox Inbox, push PushSender) *Service {
	if store == nil {
		panic("maintenanceSvc: required dependency is nil")
	}
	return &Service{store: store, operational: guard, statuses: append([]string(nil), statuses...), now: time.Now, inbox: inbox, push: push}
}

// EnsureMaintenancePushDeliverer requires the shared notification outbox so maintenance
// jobs do not fall back to the legacy maintenance_notification_deliveries claim loop.
func EnsureMaintenancePushDeliverer(push PushSender) error {
	if _, ok := push.(interface{ DeliverOutbox(context.Context) error }); !ok {
		return errors.New("maintenance module requires a push sender that implements DeliverOutbox (wire notificationService)")
	}
	return nil
}

func (s *Service) admin(ctx context.Context, id, user int64) error {
	ok, err := s.store.Admin(ctx, id, user)
	if err != nil {
		return err
	}
	if !ok {
		return models.NewAppError("MAINTENANCE_FORBIDDEN", "Active society owner or admin access required", 403, nil)
	}
	return nil
}
func (s *Service) Settings(ctx context.Context, id, user int64) (models.MaintenanceSettings, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	if err := s.admin(ctx, id, user); err != nil {
		return models.MaintenanceSettings{}, err
	}
	v, err := s.store.Settings(ctx, id)
	v.EligibleStatuses = append([]string(nil), s.statuses...)
	return v, err
}
func (s *Service) SaveSettings(ctx context.Context, id, user int64, v models.MaintenanceSettings) (models.MaintenanceSettings, error) {
	ctx, cancel := context.WithTimeout(ctx, FinancialTimeout)
	defer cancel()
	if err := v.Validate(); err != nil {
		return v, invalid(err.Error())
	}
	v.EligibleStatuses = nil
	err := s.store.Locked(ctx, id, func(ctx context.Context) error {
		if err := s.admin(ctx, id, user); err != nil {
			return err
		}
		previous, err := s.store.Settings(ctx, id)
		if err != nil {
			return err
		}
		v.FirstEnabledMonth = previous.FirstEnabledMonth
		if v.Enabled && v.FirstEnabledMonth == "" {
			loc, _ := time.LoadLocation(v.Timezone)
			v.FirstEnabledMonth = s.now().In(loc).Format("2006-01")
		}
		return s.store.SaveSettings(ctx, id, user, v)
	})
	v.EligibleStatuses = append([]string(nil), s.statuses...)
	return v, err
}

type ValidationIssues struct {
	Issues []models.MaintenanceFlatIssue `json:"issues"`
}

func (e *ValidationIssues) Error() string {
	return "Eligible flats have missing or invalid billing data"
}

// AreaCharge rounds half-up to one paise after exact multiplication.
func AreaCharge(hundredths, rate int64) (int64, error) {
	if hundredths <= 0 || rate <= 0 {
		return 0, errors.New("positive area and rate are required")
	}
	v := new(big.Int).Mul(big.NewInt(hundredths), big.NewInt(rate))
	v.Add(v, big.NewInt(50))
	v.Quo(v, big.NewInt(100))
	if !v.IsInt64() {
		return 0, errors.New("area charge exceeds supported monetary range")
	}
	return v.Int64(), nil
}
func Charges(settings models.MaintenanceSettings, f models.MaintenanceFlat) ([]models.MaintenanceItem, int64, error) {
	items := make([]models.MaintenanceItem, 0, 2)
	switch settings.PricingModel {
	case "fixed", "hybrid":
		items = append(items, models.MaintenanceItem{Description: "Fixed maintenance", AmountPaise: settings.FixedPaise})
	case "flat_type":
		if f.FlatType == nil {
			return nil, 0, errors.New("flat_type is required")
		}
		rate, ok := settings.TypeRates[*f.FlatType]
		if !ok || rate <= 0 {
			return nil, 0, fmt.Errorf("no positive rate configured for flat type %q", *f.FlatType)
		}
		items = append(items, models.MaintenanceItem{Description: "Maintenance: " + *f.FlatType, AmountPaise: rate})
	case "per_sqft":
	default:
		return nil, 0, errors.New("unsupported pricing model")
	}
	if settings.PricingModel == "per_sqft" || settings.PricingModel == "hybrid" {
		if f.AreaHundredths == nil {
			return nil, 0, errors.New("area_sqft is required")
		}
		amount, err := AreaCharge(*f.AreaHundredths, settings.AreaRatePaise)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, models.MaintenanceItem{Description: fmt.Sprintf("Area maintenance: %s sq ft at %d paise/sq ft", *models.FlatAreaString(f.AreaHundredths), settings.AreaRatePaise), AmountPaise: amount})
	}
	var total int64
	for _, i := range items {
		if i.AmountPaise < 0 || i.AmountPaise > math.MaxInt64-total {
			return nil, 0, errors.New("bill amount exceeds supported monetary range")
		}
		total += i.AmountPaise
	}
	if total <= 0 {
		return nil, 0, errors.New("bill total must be positive")
	}
	return items, total, nil
}
func CurrentMonth(month string, now time.Time, zone string) (time.Time, error) {
	loc, err := time.LoadLocation(zone)
	if err != nil {
		return time.Time{}, err
	}
	parsed, err := time.ParseInLocation("2006-01", month, loc)
	if err != nil || month != now.In(loc).Format("2006-01") {
		return time.Time{}, invalid("billing_month must be the current month in the society timezone (YYYY-MM)")
	}
	return parsed, nil
}
func DueDate(month time.Time, billingDay, dueDay int) time.Time {
	offset := 0
	if dueDay < billingDay {
		offset = 1
	}
	return time.Date(month.Year(), month.Month()+time.Month(offset), dueDay, 0, 0, 0, 0, month.Location())
}
func (s *Service) preview(ctx context.Context, id int64, month string, v models.MaintenanceSettings, catchUp ...bool) (models.MaintenancePreview, error) {
	p := models.MaintenancePreview{Bills: []models.MaintenanceBill{}, Issues: []models.MaintenanceFlatIssue{}}
	if !v.Enabled {
		return p, models.NewAppError("MAINTENANCE_DISABLED", "Maintenance billing is disabled", 409, nil)
	}
	m, err := CurrentMonth(month, s.now(), v.Timezone)
	if len(catchUp) > 0 && catchUp[0] {
		m, err = catchUpMonth(month, s.now(), v)
	}
	if err != nil {
		return p, err
	}
	flats, err := s.store.Flats(ctx, id, s.statuses)
	if err != nil {
		return p, err
	}
	return previewFlats(id, m, v, flats)
}

func previewFlats(id int64, m time.Time, v models.MaintenanceSettings, flats []models.MaintenanceFlat) (models.MaintenancePreview, error) {
	p := models.MaintenancePreview{Bills: []models.MaintenanceBill{}, Issues: []models.MaintenanceFlatIssue{}}
	month := m.Format("2006-01")
	for _, f := range flats {
		items, total, err := Charges(v, f)
		if err != nil {
			p.Issues = append(p.Issues, models.MaintenanceFlatIssue{FlatID: f.ID, Reason: err.Error()})
			continue
		}
		if total > math.MaxInt64-p.TotalPaise {
			return p, invalid("monthly total exceeds supported monetary range")
		}
		p.TotalPaise += total
		party := f.BilledParty
		f.BilledParty = nil
		p.Bills = append(p.Bills, models.MaintenanceBill{SocietyID: id, FlatID: f.ID, BillNumber: "Assigned when generated", BillingMonth: month, DueDate: DueDate(m, v.BillingDay, v.DueDay).Format("2006-01-02"), Timezone: v.Timezone, Currency: "INR", TotalPaise: total, Flat: f, Items: items, BilledParty: party})
	}
	return p, nil
}

// A run snapshot describes its original issuance. Additional bills use those
// same terms, but select currently eligible flats and omit every issued flat.
func (s *Service) previewMissing(ctx context.Context, id int64, month string) (models.MaintenancePreview, error) {
	v, issued, err := s.store.RunTerms(ctx, id, month)
	if err != nil {
		return models.MaintenancePreview{}, err
	}
	m, err := CurrentMonth(month, s.now(), v.Timezone)
	if err != nil {
		return models.MaintenancePreview{}, err
	}
	flats, err := s.store.Flats(ctx, id, v.EligibleStatuses)
	if err != nil {
		return models.MaintenancePreview{}, err
	}
	existing := make(map[int64]bool, len(issued))
	for _, id := range issued {
		existing[id] = true
	}
	missing := make([]models.MaintenanceFlat, 0, len(flats))
	for _, flat := range flats {
		if !existing[flat.ID] {
			missing = append(missing, flat)
		}
	}
	return previewFlats(id, m, v, missing)
}
func (s *Service) Preview(ctx context.Context, id, user int64, month string) (models.MaintenancePreview, error) {
	ctx, cancel := context.WithTimeout(ctx, FinancialTimeout)
	defer cancel()
	var result models.MaintenancePreview
	err := s.store.Locked(ctx, id, func(ctx context.Context) error {
		if err := s.admin(ctx, id, user); err != nil {
			return err
		}
		if err := s.operational.EnsureSocietyOperational(ctx, id); err != nil {
			return err
		}
		v, err := s.store.Settings(ctx, id)
		if err != nil {
			return err
		}
		if !v.Enabled {
			return models.NewAppError("MAINTENANCE_DISABLED", "Maintenance billing is disabled", 409, nil)
		}
		if _, err := time.Parse("2006-01", month); err != nil {
			return invalid("billing_month must be YYYY-MM")
		}
		previous, err := s.store.Run(ctx, id, month)
		if err != nil {
			return err
		}
		if previous != nil {
			result, err = s.previewMissing(ctx, id, month)
			return err
		}
		result, err = s.preview(ctx, id, month, v)
		return err
	})
	return result, err
}
func (s *Service) Generate(ctx context.Context, id, user int64, month string) (models.MaintenanceRunResult, error) {
	ctx, cancel := context.WithTimeout(ctx, FinancialTimeout)
	defer cancel()
	if err := s.admin(ctx, id, user); err != nil {
		return models.MaintenanceRunResult{}, err
	}
	return s.generate(ctx, id, user, month, false)
}
func (s *Service) generate(ctx context.Context, id, user int64, month string, scheduled bool) (models.MaintenanceRunResult, error) {
	var result models.MaintenanceRunResult
	err := s.store.Locked(ctx, id, func(ctx context.Context) error {
		if user > 0 {
			if err := s.admin(ctx, id, user); err != nil {
				return err
			}
		}
		if err := s.operational.EnsureSocietyOperational(ctx, id); err != nil {
			return err
		}
		v, err := s.store.Settings(ctx, id)
		if err != nil {
			return err
		}
		if !v.Enabled {
			if scheduled {
				return nil
			}
			return models.NewAppError("MAINTENANCE_DISABLED", "Maintenance billing is disabled", 409, nil)
		}
		if scheduled {
			loc, err := time.LoadLocation(v.Timezone)
			if err != nil {
				return err
			}
			now := s.now().In(loc)
			if now.Day() < v.BillingDay {
				return nil
			}
			month = now.Format("2006-01")
		}
		if _, err := time.Parse("2006-01", month); err != nil {
			return invalid("billing_month must be YYYY-MM")
		}
		previous, err := s.store.Run(ctx, id, month)
		if err != nil {
			return err
		}
		if previous != nil {
			result = *previous
			if scheduled {
				return nil
			}
			p, err := s.previewMissing(ctx, id, month)
			if err != nil {
				return err
			}
			if len(p.Issues) > 0 {
				return &ValidationIssues{Issues: p.Issues}
			}
			result, err = s.store.IssueMissing(ctx, id, user, *previous, p.Bills)
			return err
		}
		p, err := s.preview(ctx, id, month, v)
		if err != nil {
			return err
		}
		if len(p.Issues) > 0 {
			return &ValidationIssues{Issues: p.Issues}
		}
		v.EligibleStatuses = append([]string(nil), s.statuses...)
		result, err = s.store.Issue(ctx, id, user, month, v, p.Bills)
		return err
	})
	return result, err
}
func (s *Service) List(ctx context.Context, f models.MaintenanceBillFilter) (models.MaintenanceBillList, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	result := models.MaintenanceBillList{Items: []models.MaintenanceBill{}}
	if !f.Resident {
		if err := s.admin(ctx, f.SocietyID, f.UserID); err != nil {
			return result, err
		}
	}
	if f.Status != "" && f.Status != "unpaid" && f.Status != "overdue" && f.Status != "paid" {
		return result, invalid("status must be unpaid, overdue or paid")
	}
	if err := validStatus(f.DisplayStatus, "unpaid", "overdue", "paid", "pending_review", "rejected", "outstanding"); err != nil {
		return result, err
	}
	if f.Month != "" {
		if _, err := time.Parse("2006-01", f.Month); err != nil {
			return result, invalid("billing_month must be YYYY-MM")
		}
	}
	if f.Limit < 1 || f.Limit > 100 {
		f.Limit = 25
	}
	limit := f.Limit
	if f.Page > 0 {
		if f.Page > 100000 || f.BeforeID != 0 {
			return result, invalid("Invalid page or mixed pagination")
		}
		counter, ok := s.store.(interface {
			CountBills(context.Context, models.MaintenanceBillFilter) (int64, error)
		})
		if !ok {
			return result, errors.New("billing repository does not support page counts")
		}
		total, err := counter.CountBills(ctx, f)
		if err != nil {
			return result, err
		}
		result.TotalCount = &total
		result.Page = f.Page
		result.TotalPages = int32((total + int64(limit) - 1) / int64(limit))
		f.Offset = (f.Page - 1) * limit
	}
	f.Limit++
	bills, err := s.store.Bills(ctx, f)
	if err != nil {
		return result, err
	}
	if len(bills) > int(limit) {
		result.HasMore = true
		bills = bills[:limit]
		id := bills[len(bills)-1].ID
		result.NextCursor = &id
	}
	for i := range bills {
		if err := bills[i].PopulatePresentation(s.now()); err != nil {
			return result, err
		}
	}
	result.Items = bills
	return result, nil
}
func (s *Service) Get(ctx context.Context, f models.MaintenanceBillFilter) (models.MaintenanceBill, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	f.Limit = 1
	list, err := s.List(ctx, f)
	if err != nil {
		return models.MaintenanceBill{}, err
	}
	if len(list.Items) == 0 {
		return models.MaintenanceBill{}, models.NewAppError("MAINTENANCE_BILL_NOT_FOUND", "Bill not found", 404, nil)
	}
	return list.Items[0], nil
}

// RunScheduled checks the current month on every trigger, so a missed billing-day
// trigger is recovered without creating retroactive debt for previous months.
func (s *Service) RunScheduled(ctx context.Context) error {
	if store, ok := s.store.(billingPersistence); ok {
		if count, err := store.MaintenancePendingClaimsCount(ctx); err == nil {
			pendingClaims.Set(float64(count))
		}
	}
	ids, err := s.store.Societies(ctx)
	if err != nil {
		return err
	}
	var failures []error
	for _, id := range ids {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		billCtx, cancel := context.WithTimeout(ctx, FinancialTimeout)
		_, err = s.generate(billCtx, id, 0, "", true)
		observeFinancial("billing/scheduled", err)
		cancel()
		if err != nil {
			failures = append(failures, fmt.Errorf("society %d: %w", id, err))
		}
	}
	return errors.Join(failures...)
}

// Deliver remains as a compatibility entry point while callers migrate to the
// independent worker. It never claims Maintenance's domain staging table.
func (s *Service) Deliver(ctx context.Context) error {
	if shared, ok := s.push.(interface{ DeliverOutbox(context.Context) error }); ok {
		return shared.DeliverOutbox(ctx)
	}
	return errors.New("shared notification materializer unavailable")
}

func (s *Service) BillingRun(ctx context.Context, society, user int64, month string) (models.MaintenanceRunStatus, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	result := models.MaintenanceRunStatus{Status: "not_generated"}
	if err := s.admin(ctx, society, user); err != nil {
		return result, err
	}
	if _, err := time.Parse("2006-01", month); err != nil {
		return result, invalid("billing_month must be YYYY-MM")
	}
	run, err := s.store.Run(ctx, society, month)
	if err != nil {
		return result, err
	}
	if run != nil {
		result.Status = "completed"
		result.RunID = run.RunID
		result.BillCount = run.Existing
	}
	return result, nil
}
