package repository

import (
	"context"
	"errors"
	"go-server/internal/db"
	"go-server/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *MaintenanceRepository) Outstanding(ctx context.Context, society, user, flat int64, now time.Time) (models.MaintenanceOutstanding, error) {
	result := models.MaintenanceOutstanding{FlatID: flat, Currency: "INR", CalculatedAt: now, UnpaidBills: []models.MaintenanceOutstandingBill{}}
	// Use the same society lock as writers, keeping access and balance reads consistent.
	err := r.Locked(ctx, society, func(ctx context.Context) error {
		q := GetQueries(ctx, r.database)
		allowed, err := q.MaintenanceFlatAccess(ctx, db.MaintenanceFlatAccessParams{ID: flat, SocietyID: society, UserID: user})
		if err != nil {
			return err
		}
		if !allowed {
			return models.NewAppError("MAINTENANCE_BILL_NOT_FOUND", "Flat not found", 404, nil)
		}
		settings, err := r.Settings(ctx, society)
		if err != nil {
			return err
		}
		loc, err := time.LoadLocation(settings.Timezone)
		if err != nil {
			return err
		}
		result.CurrentMonth = now.In(loc).Format("2006-01")
		// Aggregate all financial records on the server; list pagination and
		// resident search filters never change the selected-flat summary.
		totals, err := q.UPICollectionSummary(ctx, db.UPICollectionSummaryParams{SocietyID: society, FlatID: flat})
		if err != nil {
			return err
		}
		result.TotalPaidPaise = totals.CollectedPaise
		result.OverduePaise = totals.OverduePaise
		current, err := r.Bills(ctx, models.MaintenanceBillFilter{SocietyID: society, UserID: user, FlatID: flat, Month: result.CurrentMonth, Resident: true, Limit: 1})
		if err != nil {
			return err
		}
		if len(current) > 0 {
			if err := current[0].PopulatePresentation(now); err != nil {
				return err
			}
			result.CurrentBill = &current[0]
		}
		rows, err := q.MaintenanceOutstandingBills(ctx, db.MaintenanceOutstandingBillsParams{SocietyID: society, FlatID: flat})
		if err != nil {
			return err
		}
		for _, b := range rows {
			month := b.BillingMonth.Time.Format("2006-01")
			result.UnpaidBills = append(result.UnpaidBills, models.MaintenanceOutstandingBill{ID: b.ID, BillNumber: b.BillNumber, BillingMonth: month, DueDate: b.DueDate.Time.Format("2006-01-02"), OutstandingPaise: b.OutstandingAmountPaise, Status: b.PaymentStatus})
			if month == result.CurrentMonth {
				result.CurrentMonthPaise += b.OutstandingAmountPaise
			} else if month < result.CurrentMonth {
				result.PreviousOutstandingPaise += b.OutstandingAmountPaise
			}
			result.TotalOutstandingPaise += b.OutstandingAmountPaise
		}
		return nil
	})
	return result, err
}

func (r *MaintenanceRepository) firstMonth(ctx context.Context, id int64, s *models.MaintenanceSettings) error {
	d, err := GetQueries(ctx, r.database).GetMaintenanceFirstMonth(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}
	s.FirstEnabledMonth = ""
	if d.Valid {
		s.FirstEnabledMonth = d.Time.Format("2006-01")
	}
	return nil
}

func (r *MaintenanceRepository) OutstandingFlats(ctx context.Context, f models.MaintenanceOutstandingFlatFilter) (models.MaintenanceOutstandingFlatList, error) {
	result := models.MaintenanceOutstandingFlatList{Items: []models.MaintenanceOutstandingFlatSummary{}}
	if err := r.adminOnly(ctx, f.SocietyID, f.UserID); err != nil {
		return result, err
	}
	limit := f.Limit
	if limit < 1 || limit > 100 {
		limit = 25
	}
	var block, flatNumber *string
	if f.Block != "" {
		block = &f.Block
	}
	if f.FlatNumber != "" {
		flatNumber = &f.FlatNumber
	}
	rows, err := GetQueries(ctx, r.database).ListMaintenanceOutstandingFlats(ctx, db.ListMaintenanceOutstandingFlatsParams{
		SocietyID:  f.SocietyID,
		BeforeID:   f.BeforeID,
		Block:      block,
		FlatNumber: flatNumber,
		Search:     f.Search,
		Limit:      limit + 1,
	})
	if err != nil {
		return result, err
	}
	if len(rows) > int(limit) {
		rows = rows[:limit]
		id := rows[len(rows)-1].ID
		result.NextCursor = &id
	}
	for _, row := range rows {
		item := models.MaintenanceOutstandingFlatSummary{
			Flat: models.MaintenanceFlat{
				ID:             row.ID,
				FlatNumber:     row.FlatNumber,
				Block:          row.Block,
				FlatType:       row.FlatType,
				AreaHundredths: row.AreaSqftHundredths,
			},
			TotalOutstandingPaise: row.TotalOutstandingPaise,
			UnpaidBillCount:       row.UnpaidBillCount,
		}
		if d, ok := row.OldestUnpaidMonth.(pgtype.Date); ok && d.Valid {
			item.OldestUnpaidMonth = d.Time.Format("2006-01")
		} else if t, ok := row.OldestUnpaidMonth.(time.Time); ok {
			item.OldestUnpaidMonth = t.Format("2006-01")
		}
		result.Items = append(result.Items, item)
	}
	return result, nil
}

func (r *MaintenanceRepository) adminOnly(ctx context.Context, society, user int64) error {
	ok, err := GetQueries(ctx, r.database).MaintenanceAdmin(ctx, db.MaintenanceAdminParams{SocietyID: society, UserID: user})
	if err != nil {
		return err
	}
	if !ok {
		return models.NewAppError("FORBIDDEN", "Forbidden", 403, nil)
	}
	return nil
}

func (r *MaintenanceRepository) initializeFirstMonth(ctx context.Context, id int64, s models.MaintenanceSettings) error {
	if !s.Enabled {
		return nil
	}
	d, err := time.Parse("2006-01", s.FirstEnabledMonth)
	if err != nil {
		return err
	}
	return GetQueries(ctx, r.database).SetMaintenanceFirstMonth(ctx, db.SetMaintenanceFirstMonthParams{SocietyID: id, FirstEnabledMonth: pgtype.Date{Time: d, Valid: true}})
}
