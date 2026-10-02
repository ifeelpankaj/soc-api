package repository

import (
	"context"
	"errors"
	"go-server/internal/repositories/contracts"

	"go-server/internal/db"
	"go-server/internal/models"
	"go-server/pkg/database"

	"github.com/jackc/pgx/v5"
)

type PlanRepository = contracts.PlanRepository

type planRepository struct {
	db *database.Database
}

func NewPlanRepository(db *database.Database) PlanRepository {
	return &planRepository{db: db}
}

func (r *planRepository) Create(ctx context.Context, plan *models.Plan) error {
	features, err := jsonMap(plan.Features)
	if err != nil {
		return err
	}
	row, err := GetQueries(ctx, r.db).CreatePlan(ctx, db.CreatePlanParams{
		Name: plan.Name, Code: plan.Code, Description: plan.Description,
		PriceAmountPaise: plan.PriceAmountPaise, Currency: plan.Currency,
		BillingCycle: db.BillingCycle(plan.BillingCycle), MaxFlats: plan.MaxFlats,
		MaxAdmins: plan.MaxAdmins, MaxStaff: plan.MaxStaff, MaxResidents: plan.MaxResidents, Features: features,
	})
	if err != nil {
		return err
	}
	*plan = *planFromDB(row)
	return nil
}

func (r *planRepository) Get(ctx context.Context, filter *models.PlanFilter) (*models.Plan, error) {
	row, err := GetQueries(ctx, r.db).GetPlan(ctx, planGetParams(filter))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return planFromDB(row), nil
}

func (r *planRepository) List(ctx context.Context, filter *models.PlanFilter) ([]*models.Plan, error) {
	rows, err := GetQueries(ctx, r.db).ListPlans(ctx, planListParams(filter))
	if err != nil {
		return nil, err
	}
	items := make([]*models.Plan, 0, len(rows))
	for _, row := range rows {
		items = append(items, planFromDB(row))
	}
	return items, nil
}

func (r *planRepository) Count(ctx context.Context, filter *models.PlanFilter) (int64, error) {
	return GetQueries(ctx, r.db).CountPlans(ctx, planCountParams(filter))
}

func (r *planRepository) Update(ctx context.Context, planID int64, req *contracts.UpdatePlanInput) (*models.Plan, error) {
	var features []byte
	if req.Features != nil {
		var err error
		features, err = jsonMap(req.Features)
		if err != nil {
			return nil, err
		}
	}
	row, err := GetQueries(ctx, r.db).UpdatePlan(ctx, db.UpdatePlanParams{
		ID: planID, Name: req.Name, Code: req.Code, Description: req.Description,
		PriceAmountPaise: req.PriceAmountPaise, Currency: req.Currency,
		BillingCycle: dbBillingCyclePtrFromModel(req.BillingCycle), MaxFlats: req.MaxFlats,
		MaxAdmins: req.MaxAdmins, MaxStaff: req.MaxStaff, MaxResidents: req.MaxResidents, Features: features,
	})
	return planFromDBNoRows(row, err)
}

func (r *planRepository) Activate(ctx context.Context, planID int64) (*models.Plan, error) {
	row, err := GetQueries(ctx, r.db).ActivatePlan(ctx, planID)
	return planFromDBNoRows(row, err)
}

func (r *planRepository) Deactivate(ctx context.Context, planID int64) (*models.Plan, error) {
	row, err := GetQueries(ctx, r.db).DeactivatePlan(ctx, planID)
	return planFromDBNoRows(row, err)
}

func planFromDBNoRows(row db.Plan, err error) (*models.Plan, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return planFromDB(row), nil
}

func planFromDB(row db.Plan) *models.Plan {
	return &models.Plan{
		ID: row.ID, Name: row.Name, Code: row.Code, Description: row.Description,
		PriceAmountPaise: row.PriceAmountPaise, Currency: row.Currency,
		BillingCycle: models.BillingCycle(row.BillingCycle), MaxFlats: row.MaxFlats,
		MaxAdmins: row.MaxAdmins, MaxStaff: row.MaxStaff, MaxResidents: row.MaxResidents, Features: metadataFromJSON(row.Features),
		IsActive: row.IsActive, CreatedAt: pgTimestamptzToTime(row.CreatedAt),
		UpdatedAt: pgTimestamptzToTime(row.UpdatedAt),
	}
}

func planGetParams(filter *models.PlanFilter) db.GetPlanParams {
	return db.GetPlanParams{
		ID: planID(filter), Code: planCode(filter), Name: planName(filter),
		BillingCycle: dbBillingCyclePtr(planBillingCycle(filter)), IsActive: planIsActive(filter),
		MinPricePaise: planMinPrice(filter), MaxPricePaise: planMaxPrice(filter), Search: planSearch(filter),
	}
}

func planCountParams(filter *models.PlanFilter) db.CountPlansParams {
	return db.CountPlansParams{
		ID: planID(filter), Code: planCode(filter), Name: planName(filter),
		BillingCycle: dbBillingCyclePtr(planBillingCycle(filter)), IsActive: planIsActive(filter),
		MinPricePaise: planMinPrice(filter), MaxPricePaise: planMaxPrice(filter), Search: planSearch(filter),
	}
}

func planListParams(filter *models.PlanFilter) db.ListPlansParams {
	return db.ListPlansParams{
		ID: planID(filter), Code: planCode(filter), Name: planName(filter),
		BillingCycle: dbBillingCyclePtr(planBillingCycle(filter)), IsActive: planIsActive(filter),
		MinPricePaise: planMinPrice(filter), MaxPricePaise: planMaxPrice(filter), Search: planSearch(filter),
		Limit: normalizeLimit(planLimit(filter)), Offset: normalizeOffset(planOffset(filter)),
	}
}

func planID(filter *models.PlanFilter) *int64 {
	if filter == nil {
		return nil
	}
	return filter.ID
}
func planCode(filter *models.PlanFilter) *string {
	if filter == nil {
		return nil
	}
	return filter.Code
}
func planName(filter *models.PlanFilter) *string {
	if filter == nil {
		return nil
	}
	return filter.Name
}
func planBillingCycle(filter *models.PlanFilter) *string {
	if filter == nil {
		return nil
	}
	return filter.BillingCycle
}
func planIsActive(filter *models.PlanFilter) *bool {
	if filter == nil {
		return nil
	}
	return filter.IsActive
}
func planSearch(filter *models.PlanFilter) *string {
	if filter == nil {
		return nil
	}
	return filter.Search
}
func planMinPrice(filter *models.PlanFilter) *int64 {
	if filter == nil {
		return nil
	}
	return filter.MinPricePaise
}
func planMaxPrice(filter *models.PlanFilter) *int64 {
	if filter == nil {
		return nil
	}
	return filter.MaxPricePaise
}
func planLimit(filter *models.PlanFilter) int32 {
	if filter == nil {
		return 0
	}
	return filter.Limit
}
func planOffset(filter *models.PlanFilter) int32 {
	if filter == nil {
		return 0
	}
	return filter.Offset
}

func dbBillingCyclePtr(value *string) *db.BillingCycle {
	if value == nil || *value == "" {
		return nil
	}
	cycle := db.BillingCycle(*value)
	return &cycle
}

func dbBillingCyclePtrFromModel(value *models.BillingCycle) *db.BillingCycle {
	if value == nil {
		return nil
	}
	cycle := db.BillingCycle(*value)
	return &cycle
}
