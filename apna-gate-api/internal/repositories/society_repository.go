package repository

import (
	"context"
	"encoding/json"
	"errors"
	"go-server/internal/db"
	"go-server/internal/models"
	"go-server/internal/repositories/contracts"
	"go-server/pkg/database"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type SocietyRepository = contracts.SocietyRepository

type societyRepository struct {
	db *database.Database
}

func NewSocietyRepository(db *database.Database) SocietyRepository {
	return &societyRepository{db: db}
}

func (r *societyRepository) Create(ctx context.Context, society *models.Society) error {
	return NewTransactionManager(r.db).WithTransaction(ctx, func(ctx context.Context) error {
		metadata, err := json.Marshal(society.Metadata)
		if err != nil {
			return err
		}
		if string(metadata) == "null" {
			metadata = []byte("{}")
		}

		row, err := GetQueries(ctx, r.db).CreateSociety(ctx, db.CreateSocietyParams{
			Name:         society.Name,
			SocietyCode:  society.SocietyCode,
			Email:        society.Email,
			PhoneNumber:  society.PhoneNumber,
			AddressLine1: society.AddressLine1,
			AddressLine2: society.AddressLine2,
			Landmark:     society.Landmark,
			City:         society.City,
			State:        society.State,
			Pincode:      society.Pincode,
			TotalFlats:   society.TotalFlats,
			TotalBlocks:  society.TotalBlocks,
			CreatedBy:    society.CreatedBy,
			Country:      society.Country,
			Metadata:     metadata,
		})
		if err != nil {
			return err
		}
		*society = *societyFromDB(row)
		return GetQueries(ctx, r.db).HubProvisionChannels(ctx, society.ID)
	})
}

func (r *societyRepository) Get(ctx context.Context, filter models.GetSocietyFilter) (*models.Society, error) {
	row, err := GetQueries(ctx, r.db).GetSociety(ctx, db.GetSocietyParams{
		ID:             filter.ID,
		Code:           filter.Code,
		CreatedBy:      filter.CreatedBy,
		Status:         dbSocietyStatusPtr(filter.Status),
		IncludeDeleted: filter.IncludeDeleted,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return societyFromDB(row), nil
}

func (r *societyRepository) List(ctx context.Context, filter models.ListSocietiesFilter) ([]*models.Society, error) {
	params := listSocietiesParams(filter)
	rows, err := GetQueries(ctx, r.db).ListSocieties(ctx, params)
	if err != nil {
		return nil, err
	}
	result := make([]*models.Society, 0, len(rows))
	for _, row := range rows {
		result = append(result, societyFromDB(row))
	}
	return result, nil
}

func (r *societyRepository) Count(ctx context.Context, filter models.ListSocietiesFilter) (int64, error) {
	params := listSocietiesParams(filter)
	return GetQueries(ctx, r.db).CountSocieties(ctx, db.CountSocietiesParams{
		ID: params.ID, Status: params.Status, CreatedBy: params.CreatedBy, ApprovedBy: params.ApprovedBy,
		RejectedBy: params.RejectedBy, SuspendedBy: params.SuspendedBy, CreatedFrom: params.CreatedFrom,
		CreatedTo: params.CreatedTo, Code: params.Code, Name: params.Name, City: params.City,
		State: params.State, Country: params.Country, Pincode: params.Pincode, Search: params.Search,
		SearchMode: params.SearchMode,
	})
}

func (r *societyRepository) Update(ctx context.Context, societyID int64, req contracts.UpdateSocietyInput) (*models.Society, error) {
	var metadata []byte
	if req.Metadata != nil {
		var err error
		metadata, err = json.Marshal(req.Metadata)
		if err != nil {
			return nil, err
		}
	}

	row, err := GetQueries(ctx, r.db).UpdateSociety(ctx, db.UpdateSocietyParams{
		Name: req.Name, Email: req.Email, PhoneNumber: req.PhoneNumber, AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2, Landmark: req.Landmark, City: req.City, State: req.State,
		Pincode: req.Pincode, Country: req.Country, TotalFlats: req.TotalFlats, TotalBlocks: req.TotalBlocks,
		Metadata: metadata, ID: societyID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return societyFromDB(row), nil
}

func (r *societyRepository) Approve(ctx context.Context, societyID int64, approvedBy int64) (*models.Society, error) {
	row, err := GetQueries(ctx, r.db).ApproveSociety(ctx, db.ApproveSocietyParams{ID: societyID, ApprovedBy: &approvedBy})
	return societyFromDBNoRows(row, err)
}

func (r *societyRepository) Reject(ctx context.Context, societyID int64, rejectedBy int64, reason string) (*models.Society, error) {
	row, err := GetQueries(ctx, r.db).RejectSociety(ctx, db.RejectSocietyParams{ID: societyID, RejectedBy: &rejectedBy, RejectionReason: nullableString(reason)})
	return societyFromDBNoRows(row, err)
}

func (r *societyRepository) Suspend(ctx context.Context, societyID int64, suspendedBy int64, reason string) (*models.Society, error) {
	row, err := GetQueries(ctx, r.db).SuspendSociety(ctx, db.SuspendSocietyParams{ID: societyID, SuspendedBy: &suspendedBy, SuspensionReason: nullableString(reason)})
	return societyFromDBNoRows(row, err)
}

func (r *societyRepository) Reactivate(ctx context.Context, societyID int64, reactivatedBy int64) (*models.Society, error) {
	row, err := GetQueries(ctx, r.db).ReactivateSociety(ctx, db.ReactivateSocietyParams{ID: societyID, ApprovedBy: &reactivatedBy})
	return societyFromDBNoRows(row, err)
}

func (r *societyRepository) Restore(ctx context.Context, societyID int64) (*models.Society, error) {
	row, err := GetQueries(ctx, r.db).RestoreSociety(ctx, societyID)
	return societyFromDBNoRows(row, err)
}

func (r *societyRepository) SoftDelete(ctx context.Context, societyID int64) error {
	return GetQueries(ctx, r.db).SoftDeleteSociety(ctx, societyID)
}

func (r *societyRepository) CountPendingByCreator(ctx context.Context, createdBy int64) (int64, error) {
	return GetQueries(ctx, r.db).CountPendingSocietiesByCreator(ctx, createdBy)
}

func societyFromDBNoRows(row db.Society, err error) (*models.Society, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return societyFromDB(row), nil
}

func societyFromDB(row db.Society) *models.Society {
	metadata := map[string]any{}
	if len(row.Metadata) > 0 {
		_ = json.Unmarshal(row.Metadata, &metadata)
	}
	return &models.Society{
		ID: row.ID, Name: row.Name, SocietyCode: row.SocietyCode, Email: row.Email,
		PhoneNumber: row.PhoneNumber, AddressLine1: row.AddressLine1, AddressLine2: row.AddressLine2,
		Landmark: row.Landmark, City: row.City, State: row.State, Pincode: row.Pincode, Country: row.Country,
		TotalFlats: row.TotalFlats, TotalBlocks: row.TotalBlocks, Status: models.SocietyStatus(row.Status),
		CreatedBy: row.CreatedBy, ApprovedBy: row.ApprovedBy, ApprovedAt: pgTimestamptzToTimePtr(row.ApprovedAt),
		RejectedBy: row.RejectedBy, RejectedAt: pgTimestamptzToTimePtr(row.RejectedAt), RejectionReason: row.RejectionReason,
		SuspendedBy: row.SuspendedBy, SuspendedAt: pgTimestamptzToTimePtr(row.SuspendedAt), SuspensionReason: row.SuspensionReason,
		Metadata: metadata, CreatedAt: pgTimestamptzToTime(row.CreatedAt), UpdatedAt: pgTimestamptzToTime(row.UpdatedAt),
		DeletedAt: pgTimestamptzToTimePtr(row.DeletedAt),
	}
}

func listSocietiesParams(filter models.ListSocietiesFilter) db.ListSocietiesParams {
	return db.ListSocietiesParams{
		ID: filter.ID, Status: dbSocietyStatusPtr(filter.Status), CreatedBy: filter.CreatedBy, ApprovedBy: filter.ApprovedBy,
		RejectedBy: filter.RejectedBy, SuspendedBy: filter.SuspendedBy, CreatedFrom: timePtrToPgTimestamptz(filter.CreatedFrom),
		CreatedTo: timePtrToPgTimestamptz(filter.CreatedTo), Code: filter.Code, Name: filter.Name, City: filter.City,
		State: filter.State, Country: filter.Country, Pincode: filter.Pincode, Search: filter.Search,
		SearchMode: normalizeSearchMode(filter.SearchMode, "created_by", "approved_by", "rejected_by", "suspended_by"),
		SortBy:     normalizeSocietySort(filter.SortBy), SortOrder: normalizeSortOrder(filter.SortOrder),
		Limit: normalizeLimit(filter.Limit), Offset: normalizeOffset(filter.Offset),
	}
}

func dbSocietyStatusPtr(status *string) *db.SocietyStatus {
	if status == nil || *status == "" {
		return nil
	}
	value := db.SocietyStatus(*status)
	return &value
}

func timePtrToPgTimestamptz(value *time.Time) pgtype.Timestamptz {
	if value == nil {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: *value, Valid: true}
}

func normalizeLimit(limit int32) int32 {
	if limit <= 0 {
		return 20
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func normalizeOffset(offset int32) int32 {
	if offset < 0 {
		return 0
	}
	return offset
}

func normalizeSortOrder(order string) string {
	if order == "asc" {
		return "asc"
	}
	return "desc"
}

func normalizeSocietySort(sortBy string) string {
	switch sortBy {
	case "name", "city", "status", "updated_at", "created_at":
		return sortBy
	default:
		return "created_at"
	}
}

func nullableString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
