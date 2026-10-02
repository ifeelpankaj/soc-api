package repository

import (
	"context"
	"errors"
	"go-server/internal/repositories/contracts"
	"time"

	"go-server/internal/db"
	"go-server/internal/models"
	"go-server/pkg/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type SubscriptionRepository = contracts.SubscriptionRepository

type subscriptionRepository struct {
	db *database.Database
}

func NewSubscriptionRepository(db *database.Database) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) CreatePending(ctx context.Context, societyID, planID, createdBy int64) (*models.SocietySubscription, error) {
	row, err := GetQueries(ctx, r.db).CreatePendingSubscription(ctx, db.CreatePendingSubscriptionParams{
		SocietyID: societyID, PlanID: planID, CreatedBy: createdBy, Metadata: []byte("{}"),
	})
	return subscriptionFromDBNoRows(row, err)
}

func (r *subscriptionRepository) CreateTrial(ctx context.Context, societyID, planID, createdBy int64, req *contracts.CreateTrialSubscriptionInput) (*models.SocietySubscription, error) {
	metadata, err := jsonMap(req.Metadata)
	if err != nil {
		return nil, err
	}
	row, err := GetQueries(ctx, r.db).CreateTrialSubscription(ctx, db.CreateTrialSubscriptionParams{
		SocietyID: societyID, PlanID: planID, CreatedBy: createdBy,
		StartsAt: timePtrToPgTimestamptz(&req.StartsAt), EndsAt: timePtrToPgTimestamptz(req.EndsAt),
		TrialEndsAt: timePtrToPgTimestamptz(&req.TrialEndsAt), Metadata: metadata,
	})
	return subscriptionFromDBNoRows(row, err)
}

func (r *subscriptionRepository) Activate(ctx context.Context, subscriptionID, activatedBy int64, req *contracts.ActivateSubscriptionInput) (*models.SocietySubscription, error) {
	metadata, err := jsonMap(req.Metadata)
	if err != nil {
		return nil, err
	}
	row, err := GetQueries(ctx, r.db).ActivateSubscription(ctx, db.ActivateSubscriptionParams{
		ID: subscriptionID, StartsAt: timePtrToPgTimestamptz(&req.StartsAt),
		EndsAt: timePtrToPgTimestamptz(&req.EndsAt), ActivatedBy: &activatedBy, Metadata: metadata,
	})
	return subscriptionFromDBNoRows(row, err)
}

func (r *subscriptionRepository) Renew(ctx context.Context, subscriptionID, renewedBy int64, req *contracts.RenewSubscriptionInput) (*models.SocietySubscription, error) {
	metadata, err := jsonMap(req.Metadata)
	if err != nil {
		return nil, err
	}
	row, err := GetQueries(ctx, r.db).RenewSubscription(ctx, db.RenewSubscriptionParams{
		ID: subscriptionID, StartsAt: timePtrToPgTimestamptz(&req.StartsAt),
		EndsAt: timePtrToPgTimestamptz(&req.EndsAt), ActivatedBy: &renewedBy, Metadata: metadata,
	})
	return subscriptionFromDBNoRows(row, err)
}

func (r *subscriptionRepository) Cancel(ctx context.Context, subscriptionID, cancelledBy int64, req *contracts.CancelSubscriptionInput) (*models.SocietySubscription, error) {
	metadata, err := jsonMap(req.Metadata)
	if err != nil {
		return nil, err
	}
	row, err := GetQueries(ctx, r.db).CancelSubscription(ctx, db.CancelSubscriptionParams{
		ID: subscriptionID, CancelledBy: &cancelledBy, CancellationReason: nullableString(req.Reason), Metadata: metadata,
	})
	return subscriptionFromDBNoRows(row, err)
}

func (r *subscriptionRepository) Expire(ctx context.Context, subscriptionID int64) (*models.SocietySubscription, error) {
	row, err := GetQueries(ctx, r.db).ExpireSubscription(ctx, subscriptionID)
	return subscriptionFromDBNoRows(row, err)
}

func (r *subscriptionRepository) ExpireDue(ctx context.Context) (int64, error) {
	return GetQueries(ctx, r.db).ExpireDueSubscriptions(ctx)
}

func (r *subscriptionRepository) ChangePlan(ctx context.Context, subscriptionID, newPlanID int64) (*models.SocietySubscription, error) {
	row, err := GetQueries(ctx, r.db).ChangeSubscriptionPlan(ctx, db.ChangeSubscriptionPlanParams{ID: subscriptionID, ID_2: newPlanID})
	return subscriptionFromDBNoRows(row, err)
}

func (r *subscriptionRepository) Get(ctx context.Context, filter *models.SubscriptionFilter) (*models.SocietySubscription, error) {
	row, err := GetQueries(ctx, r.db).GetSubscription(ctx, subscriptionGetParams(filter))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return subscriptionFromGetRow(row), nil
}

func (r *subscriptionRepository) List(ctx context.Context, filter *models.SubscriptionFilter) ([]*models.SocietySubscription, error) {
	rows, err := GetQueries(ctx, r.db).ListSubscriptions(ctx, subscriptionListParams(filter))
	if err != nil {
		return nil, err
	}
	items := make([]*models.SocietySubscription, 0, len(rows))
	for _, row := range rows {
		items = append(items, subscriptionFromListRow(row))
	}
	return items, nil
}

func (r *subscriptionRepository) Stats(ctx context.Context, filter *models.SubscriptionFilter) (*models.SubscriptionStatsResponse, error) {
	row, err := GetQueries(ctx, r.db).GetSubscriptionStats(ctx, db.GetSubscriptionStatsParams{
		SocietyID: subSocietyID(filter), PlanID: subPlanID(filter), Status: dbSubscriptionStatusPtr(subStatus(filter)),
		PlanCode: subPlanCode(filter), BillingCycle: dbBillingCyclePtr(subBillingCycle(filter)),
	})
	if err != nil {
		return nil, err
	}
	return &models.SubscriptionStatsResponse{
		TotalSubscriptions: row.TotalSubscriptions, PendingSubscriptions: row.PendingSubscriptions,
		TrialSubscriptions: row.TrialSubscriptions, ActiveSubscriptions: row.ActiveSubscriptions,
		ExpiredSubscriptions: row.ExpiredSubscriptions, CancelledSubscriptions: row.CancelledSubscriptions,
	}, nil
}

func (r *subscriptionRepository) CountActiveFlats(ctx context.Context, societyID int64) (int64, error) {
	return GetQueries(ctx, r.db).CountActiveFlatsForQuota(ctx, societyID)
}

func (r *subscriptionRepository) CountActiveAdmins(ctx context.Context, societyID int64) (int64, error) {
	return GetQueries(ctx, r.db).CountActiveAdminsForQuota(ctx, societyID)
}

func (r *subscriptionRepository) CountActiveStaff(ctx context.Context, societyID int64) (int64, error) {
	return GetQueries(ctx, r.db).CountActiveStaffForQuota(ctx, societyID)
}

func (r *subscriptionRepository) CountActiveResidents(ctx context.Context, societyID int64) (int64, error) {
	return GetQueries(ctx, r.db).CountActiveResidentsForQuota(ctx, societyID)
}

func (r *subscriptionRepository) GetActiveForUpdate(ctx context.Context, societyID int64) (*models.SocietySubscription, error) {
	row, err := GetQueries(ctx, r.db).GetActiveSubscriptionForUpdate(ctx, societyID)
	return subscriptionFromDBNoRows(row, err)
}

func subscriptionFromDBNoRows(row db.SocietySubscription, err error) (*models.SocietySubscription, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return subscriptionFromDB(row), nil
}

func subscriptionFromDB(row db.SocietySubscription) *models.SocietySubscription {
	return &models.SocietySubscription{
		ID: row.ID, SocietyID: row.SocietyID, PlanID: row.PlanID, Status: models.SubscriptionStatus(row.Status),
		StartsAt: pgTimestamptzToTimePtr(row.StartsAt), EndsAt: pgTimestamptzToTimePtr(row.EndsAt),
		TrialEndsAt: pgTimestamptzToTimePtr(row.TrialEndsAt), PlanName: row.PlanName, PlanCode: row.PlanCode,
		PriceAmountPaise: row.PriceAmountPaise, Currency: row.Currency, BillingCycle: models.BillingCycle(row.BillingCycle),
		MaxFlats: row.MaxFlats, MaxAdmins: row.MaxAdmins, MaxStaff: row.MaxStaff, MaxResidents: row.MaxResidents, Features: metadataFromJSON(row.Features),
		ActivatedAt: pgTimestamptzToTimePtr(row.ActivatedAt), ActivatedBy: row.ActivatedBy,
		ExpiredAt: pgTimestamptzToTimePtr(row.ExpiredAt), CancelledAt: pgTimestamptzToTimePtr(row.CancelledAt),
		CancelledBy: row.CancelledBy, CancellationReason: row.CancellationReason, Metadata: metadataFromJSON(row.Metadata),
		CreatedBy: row.CreatedBy, CreatedAt: pgTimestamptzToTime(row.CreatedAt), UpdatedAt: pgTimestamptzToTime(row.UpdatedAt),
	}
}

func subscriptionFromGetRow(row db.GetSubscriptionRow) *models.SocietySubscription {
	sub := subscriptionFromParts(row.ID, row.SocietyID, row.PlanID, row.Status, row.StartsAt, row.EndsAt, row.TrialEndsAt, row.PlanName, row.PlanCode, row.PriceAmountPaise, row.Currency, row.BillingCycle, row.MaxFlats, row.MaxAdmins, row.MaxStaff, row.MaxResidents, row.Features, row.ActivatedAt, row.ActivatedBy, row.ExpiredAt, row.CancelledAt, row.CancelledBy, row.CancellationReason, row.Metadata, row.CreatedBy, row.CreatedAt, row.UpdatedAt)
	sub.SocietyName = &row.SocietyName
	sub.SocietyCode = &row.SocietyCode
	sub.CurrentPlanName = row.CurrentPlanName
	sub.CurrentPlanCode = row.CurrentPlanCode
	return sub
}

func subscriptionFromListRow(row db.ListSubscriptionsRow) *models.SocietySubscription {
	sub := subscriptionFromParts(row.ID, row.SocietyID, row.PlanID, row.Status, row.StartsAt, row.EndsAt, row.TrialEndsAt, row.PlanName, row.PlanCode, row.PriceAmountPaise, row.Currency, row.BillingCycle, row.MaxFlats, row.MaxAdmins, row.MaxStaff, row.MaxResidents, row.Features, row.ActivatedAt, row.ActivatedBy, row.ExpiredAt, row.CancelledAt, row.CancelledBy, row.CancellationReason, row.Metadata, row.CreatedBy, row.CreatedAt, row.UpdatedAt)
	sub.SocietyName = &row.SocietyName
	sub.SocietyCode = &row.SocietyCode
	sub.CurrentPlanName = row.CurrentPlanName
	sub.CurrentPlanCode = row.CurrentPlanCode
	return sub
}

func subscriptionFromParts(id, societyID int64, planID *int64, status db.SubscriptionStatus, startsAt, endsAt, trialEndsAt pgtype.Timestamptz, planName, planCode string, price int64, currency string, billing db.BillingCycle, maxFlats, maxAdmins, maxStaff, maxResidents int32, features []byte, activatedAt pgtype.Timestamptz, activatedBy *int64, expiredAt, cancelledAt pgtype.Timestamptz, cancelledBy *int64, reason *string, metadata []byte, createdBy *int64, createdAt, updatedAt pgtype.Timestamptz) *models.SocietySubscription {
	return &models.SocietySubscription{
		ID: id, SocietyID: societyID, PlanID: planID, Status: models.SubscriptionStatus(status),
		StartsAt: pgTimestamptzToTimePtr(startsAt), EndsAt: pgTimestamptzToTimePtr(endsAt), TrialEndsAt: pgTimestamptzToTimePtr(trialEndsAt),
		PlanName: planName, PlanCode: planCode, PriceAmountPaise: price, Currency: currency,
		BillingCycle: models.BillingCycle(billing), MaxFlats: maxFlats, MaxAdmins: maxAdmins, MaxStaff: maxStaff, MaxResidents: maxResidents,
		Features: metadataFromJSON(features), ActivatedAt: pgTimestamptzToTimePtr(activatedAt), ActivatedBy: activatedBy,
		ExpiredAt: pgTimestamptzToTimePtr(expiredAt), CancelledAt: pgTimestamptzToTimePtr(cancelledAt),
		CancelledBy: cancelledBy, CancellationReason: reason, Metadata: metadataFromJSON(metadata), CreatedBy: createdBy,
		CreatedAt: pgTimestamptzToTime(createdAt), UpdatedAt: pgTimestamptzToTime(updatedAt),
	}
}

func subscriptionGetParams(filter *models.SubscriptionFilter) db.GetSubscriptionParams {
	return db.GetSubscriptionParams{ID: subID(filter), SocietyID: subSocietyID(filter), PlanID: subPlanID(filter), Status: dbSubscriptionStatusPtr(subStatus(filter)), PlanCode: subPlanCode(filter), BillingCycle: dbBillingCyclePtr(subBillingCycle(filter)), IsActiveOnly: subIsActiveOnly(filter), StartsAfter: timePtrToPgTimestamptz(subStartsAfter(filter)), StartsBefore: timePtrToPgTimestamptz(subStartsBefore(filter)), EndsAfter: timePtrToPgTimestamptz(subEndsAfter(filter)), EndsBefore: timePtrToPgTimestamptz(subEndsBefore(filter)), ExpiringBefore: timePtrToPgTimestamptz(subExpiringBefore(filter)), ExpiredOnly: subExpiredOnly(filter), Search: subSearch(filter), SearchMode: subSearchMode(filter)}
}

func subscriptionListParams(filter *models.SubscriptionFilter) db.ListSubscriptionsParams {
	return db.ListSubscriptionsParams{ID: subID(filter), SocietyID: subSocietyID(filter), PlanID: subPlanID(filter), Status: dbSubscriptionStatusPtr(subStatus(filter)), PlanCode: subPlanCode(filter), BillingCycle: dbBillingCyclePtr(subBillingCycle(filter)), IsActiveOnly: subIsActiveOnly(filter), StartsAfter: timePtrToPgTimestamptz(subStartsAfter(filter)), StartsBefore: timePtrToPgTimestamptz(subStartsBefore(filter)), EndsAfter: timePtrToPgTimestamptz(subEndsAfter(filter)), EndsBefore: timePtrToPgTimestamptz(subEndsBefore(filter)), ExpiringBefore: timePtrToPgTimestamptz(subExpiringBefore(filter)), ExpiredOnly: subExpiredOnly(filter), Search: subSearch(filter), SearchMode: subSearchMode(filter), Limit: normalizeLimit(subLimit(filter)), Offset: normalizeOffset(subOffset(filter))}
}

func subID(f *models.SubscriptionFilter) *int64 {
	if f == nil {
		return nil
	}
	return f.ID
}
func subSocietyID(f *models.SubscriptionFilter) *int64 {
	if f == nil {
		return nil
	}
	return f.SocietyID
}
func subPlanID(f *models.SubscriptionFilter) *int64 {
	if f == nil {
		return nil
	}
	return f.PlanID
}
func subStatus(f *models.SubscriptionFilter) *string {
	if f == nil {
		return nil
	}
	return f.Status
}
func subPlanCode(f *models.SubscriptionFilter) *string {
	if f == nil {
		return nil
	}
	return f.PlanCode
}
func subBillingCycle(f *models.SubscriptionFilter) *string {
	if f == nil {
		return nil
	}
	return f.BillingCycle
}
func subIsActiveOnly(f *models.SubscriptionFilter) *bool {
	if f == nil {
		return nil
	}
	return f.IsActiveOnly
}
func subStartsAfter(f *models.SubscriptionFilter) *time.Time {
	if f == nil {
		return nil
	}
	return f.StartsAfter
}
func subStartsBefore(f *models.SubscriptionFilter) *time.Time {
	if f == nil {
		return nil
	}
	return f.StartsBefore
}
func subEndsAfter(f *models.SubscriptionFilter) *time.Time {
	if f == nil {
		return nil
	}
	return f.EndsAfter
}
func subEndsBefore(f *models.SubscriptionFilter) *time.Time {
	if f == nil {
		return nil
	}
	return f.EndsBefore
}
func subExpiringBefore(f *models.SubscriptionFilter) *time.Time {
	if f == nil {
		return nil
	}
	return f.ExpiringBefore
}
func subExpiredOnly(f *models.SubscriptionFilter) *bool {
	if f == nil {
		return nil
	}
	return f.ExpiredOnly
}
func subSearch(f *models.SubscriptionFilter) *string {
	if f == nil {
		return nil
	}
	return f.Search
}
func subSearchMode(f *models.SubscriptionFilter) *string {
	mode := ""
	if f != nil && f.SearchMode != nil {
		mode = *f.SearchMode
	}
	normalized := normalizeSearchMode(mode, "society", "plan", "action_user", "resident_member")
	return &normalized
}
func subLimit(f *models.SubscriptionFilter) int32 {
	if f == nil {
		return 0
	}
	return f.Limit
}
func subOffset(f *models.SubscriptionFilter) int32 {
	if f == nil {
		return 0
	}
	return f.Offset
}

func dbSubscriptionStatusPtr(status *string) *db.SubscriptionStatus {
	if status == nil || *status == "" {
		return nil
	}
	value := db.SubscriptionStatus(*status)
	return &value
}
