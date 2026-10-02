package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-server/internal/db"
	"go-server/internal/models"
	"go-server/internal/repositories/contracts"
	"go-server/pkg/database"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type MaintenanceRepository struct {
	*financialRepository
	database *database.Database
	tx       TransactionManager
}

func NewMaintenanceRepository(d *database.Database, tx TransactionManager) *MaintenanceRepository {
	return &MaintenanceRepository{newFinancialRepository(d), d, tx}
}
func (r *MaintenanceRepository) Locked(ctx context.Context, id int64, fn func(context.Context) error) error {
	return r.tx.WithTransaction(ctx, func(ctx context.Context) error {
		if err := GetQueries(ctx, r.database).MaintenanceLock(ctx, id); err != nil {
			return err
		}
		return fn(ctx)
	})
}
func (r *MaintenanceRepository) Admin(ctx context.Context, society, user int64) (bool, error) {
	return GetQueries(ctx, r.database).MaintenanceAdmin(ctx, db.MaintenanceAdminParams{SocietyID: society, UserID: user})
}
func (r *MaintenanceRepository) Settings(ctx context.Context, id int64) (models.MaintenanceSettings, error) {
	b, err := GetQueries(ctx, r.database).GetMaintenanceSettings(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.DefaultMaintenanceSettings(), nil
	}
	if err != nil {
		return models.MaintenanceSettings{}, err
	}
	var s models.MaintenanceSettings
	err = json.Unmarshal(b, &s)
	if err == nil {
		err = r.firstMonth(ctx, id, &s)
	}
	return s, err
}
func (r *MaintenanceRepository) SaveSettings(ctx context.Context, id, user int64, s models.MaintenanceSettings) error {
	q := GetQueries(ctx, r.database)
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	if err = q.SaveMaintenanceSettings(ctx, db.SaveMaintenanceSettingsParams{SocietyID: id, Enabled: s.Enabled, Config: b, UpdatedBy: &user}); err != nil {
		return err
	}
	if err = r.initializeFirstMonth(ctx, id, s); err != nil {
		return err
	}
	if err = q.DeleteMaintenanceTypeRates(ctx, id); err != nil {
		return err
	}
	for k, v := range s.TypeRates {
		if err = q.SaveMaintenanceTypeRate(ctx, db.SaveMaintenanceTypeRateParams{SocietyID: id, FlatType: k, AmountPaise: v}); err != nil {
			return err
		}
	}
	return nil
}
func (r *MaintenanceRepository) Flats(ctx context.Context, id int64, statuses []string) ([]models.MaintenanceFlat, error) {
	rows, err := GetQueries(ctx, r.database).ListMaintenanceFlats(ctx, db.ListMaintenanceFlatsParams{SocietyID: id, Statuses: statuses})
	if err != nil {
		return nil, err
	}
	result := make([]models.MaintenanceFlat, 0, len(rows))
	for _, v := range rows {
		var party map[string]any
		if err = json.Unmarshal(v.BilledParty, &party); err != nil {
			return nil, err
		}
		result = append(result, models.MaintenanceFlat{ID: v.ID, FlatNumber: v.FlatNumber, Block: v.Block, FlatType: v.FlatType, AreaHundredths: v.AreaSqftHundredths, BilledParty: party})
	}
	return result, nil
}
func maintenanceDate(s string) pgtype.Date {
	t, err := time.Parse("2006-01-02", s)
	return pgtype.Date{Time: t, Valid: err == nil}
}
func (r *MaintenanceRepository) Run(ctx context.Context, id int64, month string) (*models.MaintenanceRunResult, error) {
	row, err := GetQueries(ctx, r.database).GetMaintenanceRun(ctx, db.GetMaintenanceRunParams{SocietyID: id, BillingMonth: maintenanceDate(month + "-01")})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &models.MaintenanceRunResult{RunID: row.ID, Existing: int(row.BillCount)}, nil
}

func (r *MaintenanceRepository) RunTerms(ctx context.Context, id int64, month string) (models.MaintenanceSettings, []int64, error) {
	q := GetQueries(ctx, r.database)
	row, err := q.GetMaintenanceRun(ctx, db.GetMaintenanceRunParams{SocietyID: id, BillingMonth: maintenanceDate(month + "-01")})
	var snapshot struct {
		Settings models.MaintenanceSettings `json:"settings"`
	}
	if err != nil {
		return snapshot.Settings, nil, err
	}
	if err := json.Unmarshal(row.Snapshot, &snapshot); err != nil {
		return snapshot.Settings, nil, err
	}
	if err := snapshot.Settings.Validate(); err != nil {
		return snapshot.Settings, nil, err
	}
	if len(snapshot.Settings.EligibleStatuses) == 0 {
		return snapshot.Settings, nil, errors.New("issued month has no saved eligibility policy")
	}
	ids, err := q.ListMaintenanceIssuedFlats(ctx, db.ListMaintenanceIssuedFlatsParams{SocietyID: id, BillingMonth: maintenanceDate(month + "-01")})
	return snapshot.Settings, ids, err
}

func (r *MaintenanceRepository) IssueMissing(ctx context.Context, id, user int64, previous models.MaintenanceRunResult, bills []models.MaintenanceBill) (models.MaintenanceRunResult, error) {
	q := GetQueries(ctx, r.database)
	ids, err := insertMaintenanceBills(ctx, q, previous.RunID, id, bills)
	if err != nil {
		return models.MaintenanceRunResult{}, err
	}
	result := models.MaintenanceRunResult{RunID: previous.RunID, Created: len(bills), Existing: previous.Existing}
	details, err := json.Marshal(struct {
		Result  models.MaintenanceRunResult `json:"result"`
		BillIDs []int64                     `json:"bill_ids"`
	}{result, ids})
	if err != nil {
		return result, err
	}
	err = q.AuditMaintenanceReconciliation(ctx, db.AuditMaintenanceReconciliationParams{SocietyID: id, ActorID: &user, EntityID: previous.RunID, Details: details})
	return result, err
}
func (r *MaintenanceRepository) Issue(ctx context.Context, id, user int64, month string, s models.MaintenanceSettings, bills []models.MaintenanceBill) (models.MaintenanceRunResult, error) {
	var result models.MaintenanceRunResult
	ids := make([]int64, 0, len(bills))
	for _, b := range bills {
		ids = append(ids, b.FlatID)
	}
	snap, err := json.Marshal(struct {
		Settings models.MaintenanceSettings `json:"settings"`
		FlatIDs  []int64                    `json:"flat_ids"`
	}{s, ids})
	if err != nil {
		return result, err
	}
	var actor *int64
	if user > 0 {
		actor = &user
	}
	q := GetQueries(ctx, r.database)
	run, err := q.CreateMaintenanceRun(ctx, db.CreateMaintenanceRunParams{SocietyID: id, BillingMonth: maintenanceDate(month + "-01"), Snapshot: snap, CreatedBy: actor})
	if err != nil {
		return result, err
	}
	if _, err := insertMaintenanceBills(ctx, q, run, id, bills); err != nil {
		return result, err
	}
	return models.MaintenanceRunResult{RunID: run, Created: len(bills)}, nil
}

func insertMaintenanceBills(ctx context.Context, q *db.Queries, run, id int64, bills []models.MaintenanceBill) ([]int64, error) {
	ids := make([]int64, 0, len(bills))
	society, err := q.GetSociety(ctx, db.GetSocietyParams{ID: &id})
	if err != nil {
		return nil, err
	}
	for _, b := range bills {
		zone, err := time.LoadLocation(b.Timezone)
		if err != nil {
			return nil, err
		}
		b.BillNumber = maintenanceBillNumber(society.SocietyCode, b.Flat, time.Now().In(zone))
		party, err := json.Marshal(b.BilledParty)
		if err != nil {
			return nil, err
		}
		b.BilledParty = nil
		b.Flat.BilledParty = nil
		snapshot, err := json.Marshal(b)
		if err != nil {
			return nil, err
		}
		billID, err := q.CreateMaintenanceBill(ctx, db.CreateMaintenanceBillParams{RunID: run, SocietyID: id, FlatID: b.FlatID, BillingMonth: maintenanceDate(b.BillingMonth + "-01"), BillNumber: b.BillNumber, DueDate: maintenanceDate(b.DueDate), Timezone: b.Timezone, TotalPaise: b.TotalPaise, Snapshot: snapshot, BilledParty: party})
		if err != nil {
			return nil, err
		}
		for i, item := range b.Items {
			if err = q.CreateMaintenanceBillItem(ctx, db.CreateMaintenanceBillItemParams{BillID: billID, Position: int32(i), Description: item.Description, AmountPaise: item.AmountPaise}); err != nil {
				return nil, err
			}
		}
		if err = q.EnqueueMaintenanceNotifications(ctx, billID); err != nil {
			return nil, err
		}
		ids = append(ids, billID)
	}
	return ids, nil
}
func maintenanceBillNumber(code string, flat models.MaintenanceFlat, issuedAt time.Time) string {
	number := flat.FlatNumber
	if flat.Block != nil && *flat.Block != "" {
		number = *flat.Block + "-" + number
	}
	return fmt.Sprintf("AG-%s-%s-%s", code, issuedAt.Format("20060102150405"), number)
}
func maintenanceBill(row db.MaintenanceBill, resident bool) (models.MaintenanceBill, error) {
	var b models.MaintenanceBill
	if err := json.Unmarshal(row.Snapshot, &b); err != nil {
		return b, err
	}
	if err := json.Unmarshal(row.IssuerSnapshot, &b.Issuer); err != nil {
		return b, err
	}
	b.ID = row.ID
	b.RunID = row.RunID
	b.SocietyID = row.SocietyID
	b.FlatID = row.FlatID
	b.BillNumber = row.BillNumber
	b.BillingMonth = row.BillingMonth.Time.Format("2006-01")
	b.DueDate = row.DueDate.Time.Format("2006-01-02")
	b.Timezone = row.Timezone
	b.TotalPaise = row.TotalPaise
	b.CreatedAt = row.CreatedAt.Time
	b.IssuedAt = row.CreatedAt.Time
	b.BilledParty = nil
	b.Flat.BilledParty = nil
	if !resident && len(row.BilledParty) > 0 {
		if err := json.Unmarshal(row.BilledParty, &b.BilledParty); err != nil {
			return b, err
		}
	}
	return b, nil
}
func MaintenanceFlatOptionalStrings(block, flatNumber string) (*string, *string) {
	var blockPtr, flatNumberPtr *string
	if block != "" {
		blockPtr = &block
	}
	if flatNumber != "" {
		flatNumberPtr = &flatNumber
	}
	return blockPtr, flatNumberPtr
}

func (r *MaintenanceRepository) Bills(ctx context.Context, f models.MaintenanceBillFilter) ([]models.MaintenanceBill, error) {
	month := pgtype.Date{}
	if f.Month != "" {
		month = maintenanceDate(f.Month + "-01")
	}
	block, flatNumber := MaintenanceFlatOptionalStrings(f.Block, f.FlatNumber)
	rows, err := GetQueries(ctx, r.database).ListMaintenanceBills(ctx, db.ListMaintenanceBillsParams{
		SocietyID: f.SocietyID, BillID: f.ID, FlatID: f.FlatID, BeforeID: f.BeforeID, Month: month,
		Status: f.Status, BillNumber: f.BillNumber, Block: block, FlatNumber: flatNumber, Search: f.Search,
		Resident: f.Resident, UserID: f.UserID, Limit: f.Limit,
		DisplayStatus: f.DisplayStatus, PageOffset: f.Offset,
	})
	if err != nil {
		return nil, err
	}
	result := make([]models.MaintenanceBill, 0, len(rows))
	for _, row := range rows {
		b, err := maintenanceBill(row.MaintenanceBill, f.Resident)
		if err != nil {
			return nil, err
		}
		b.PaidAmountPaise = row.PaidAmountPaise
		b.OutstandingAmountPaise = row.OutstandingAmountPaise
		b.PaymentClaimStatus = row.PaymentClaimStatus
		b.Status = row.PaymentStatus
		if row.PaidOn.Valid {
			b.PaidOn = row.PaidOn.Time.Format("2006-01-02")
		}
		result = append(result, b)
	}
	return result, nil
}
func (r *MaintenanceRepository) CountBills(ctx context.Context, f models.MaintenanceBillFilter) (int64, error) {
	month := pgtype.Date{}
	if f.Month != "" {
		month = maintenanceDate(f.Month + "-01")
	}
	block, flatNumber := MaintenanceFlatOptionalStrings(f.Block, f.FlatNumber)
	return GetQueries(ctx, r.database).CountMaintenanceBills(ctx, db.CountMaintenanceBillsParams{
		SocietyID: f.SocietyID, BillID: f.ID, FlatID: f.FlatID, Month: month,
		Status: f.Status, DisplayStatus: f.DisplayStatus, BillNumber: f.BillNumber,
		Block: block, FlatNumber: flatNumber, Search: f.Search, Resident: f.Resident, UserID: f.UserID,
	})
}
func (r *MaintenanceRepository) Societies(ctx context.Context) ([]int64, error) {
	return GetQueries(ctx, r.database).ListMaintenanceSocieties(ctx)
}

type MaintenanceDelivery = contracts.MaintenanceDelivery

func (r *MaintenanceRepository) ClaimDelivery(ctx context.Context) (*MaintenanceDelivery, error) {
	q := GetQueries(ctx, r.database)
	row, err := q.ClaimMaintenanceDelivery(ctx, pgtype.UUID{Bytes: uuid.New(), Valid: true})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	d := &MaintenanceDelivery{ID: row.ID, UserID: row.UserID, Token: uuid.UUID(row.LeaseToken.Bytes), EventType: row.EventType, EventKey: row.EventKey}
	if err := json.Unmarshal(row.EventData, &d.EventData); err != nil {
		return nil, err
	}
	b, err := q.GetMaintenanceDeliveryBill(ctx, db.GetMaintenanceDeliveryBillParams{BillID: row.BillID, UserID: row.UserID, Audience: row.Audience})
	if errors.Is(err, pgx.ErrNoRows) {
		return d, nil
	}
	if err != nil {
		return nil, err
	}
	bill, err := maintenanceBill(b, true)
	if err != nil {
		return nil, err
	}
	d.Bill = &bill
	return d, nil
}
func (r *MaintenanceRepository) FinishDelivery(ctx context.Context, d *MaintenanceDelivery, cause error) error {
	q := GetQueries(ctx, r.database)
	if cause == nil {
		return q.CompleteMaintenanceDelivery(ctx, db.CompleteMaintenanceDeliveryParams{ID: d.ID, LeaseToken: pgtype.UUID{Bytes: d.Token, Valid: d.Token != uuid.Nil}})
	}
	message := cause.Error()
	if len(message) > 1000 {
		message = message[:1000]
	}
	return q.RetryMaintenanceDelivery(ctx, db.RetryMaintenanceDeliveryParams{ID: d.ID, LeaseToken: pgtype.UUID{Bytes: d.Token, Valid: d.Token != uuid.Nil}, LastError: &message})
}
