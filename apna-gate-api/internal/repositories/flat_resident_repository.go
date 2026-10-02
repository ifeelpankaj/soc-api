package repository

import (
	"context"
	"errors"
	"go-server/internal/repositories/contracts"

	"go-server/internal/db"
	"go-server/internal/models"
	"go-server/pkg/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type FlatResidentRepository = contracts.FlatResidentRepository

type flatResidentRepository struct {
	db *database.Database
}

func NewFlatResidentRepository(db *database.Database) FlatResidentRepository {
	return &flatResidentRepository{db: db}
}

func (r *flatResidentRepository) Add(ctx context.Context, resident *models.FlatResident) error {
	metadata, err := jsonMap(resident.Metadata)
	if err != nil {
		return err
	}
	row, err := GetQueries(ctx, r.db).AddFlatResident(ctx, db.AddFlatResidentParams{
		SocietyID: resident.SocietyID, FlatID: resident.FlatID, UserID: resident.UserID,
		Role: db.FlatResidentRole(resident.Role), IsPrimary: resident.IsPrimary,
		CreatedBy: resident.CreatedBy, Metadata: metadata,
	})
	if err != nil {
		return err
	}
	*resident = *flatResidentFromDB(row)
	return nil
}

func (r *flatResidentRepository) Get(ctx context.Context, filter *models.FlatResidentFilter) (*models.FlatResident, error) {
	row, err := GetQueries(ctx, r.db).GetFlatResident(ctx, residentGetParams(filter))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return flatResidentFromGetRow(row), nil
}

func (r *flatResidentRepository) List(ctx context.Context, filter *models.FlatResidentFilter) ([]*models.FlatResident, error) {
	rows, err := GetQueries(ctx, r.db).ListFlatResidents(ctx, residentListParams(filter))
	if err != nil {
		return nil, err
	}
	items := make([]*models.FlatResident, 0, len(rows))
	for _, row := range rows {
		items = append(items, flatResidentFromListRow(row))
	}
	return items, nil
}

func (r *flatResidentRepository) Remove(ctx context.Context, filter *models.FlatResidentFilter) error {
	return GetQueries(ctx, r.db).RemoveFlatResident(ctx, db.RemoveFlatResidentParams{
		ID: residentID(filter), SocietyID: residentSocietyID(filter), FlatID: residentFlatID(filter), UserID: residentUserID(filter),
	})
}

func (r *flatResidentRepository) MoveOut(ctx context.Context, filter *models.FlatResidentFilter) (*models.FlatResident, error) {
	row, err := GetQueries(ctx, r.db).MoveOutFlatResident(ctx, db.MoveOutFlatResidentParams{
		ID: residentID(filter), SocietyID: residentSocietyID(filter), FlatID: residentFlatID(filter), UserID: residentUserID(filter),
	})
	return flatResidentFromDBNoRows(row, err)
}

func (r *flatResidentRepository) ClearPrimary(ctx context.Context, societyID int64, flatID int64) error {
	return GetQueries(ctx, r.db).ClearPrimaryFlatResident(ctx, db.ClearPrimaryFlatResidentParams{SocietyID: societyID, FlatID: flatID})
}

func (r *flatResidentRepository) SetPrimary(ctx context.Context, societyID int64, flatID int64, residentID int64) (*models.FlatResident, error) {
	row, err := GetQueries(ctx, r.db).SetPrimaryFlatResident(ctx, db.SetPrimaryFlatResidentParams{SocietyID: societyID, FlatID: flatID, ID: residentID})
	return flatResidentFromDBNoRows(row, err)
}

func (r *flatResidentRepository) UpdateRole(ctx context.Context, filter *models.FlatResidentFilter, role models.FlatResidentRole) (*models.FlatResident, error) {
	row, err := GetQueries(ctx, r.db).UpdateFlatResidentRole(ctx, db.UpdateFlatResidentRoleParams{
		Role: db.FlatResidentRole(role), ID: residentID(filter), SocietyID: residentSocietyID(filter),
		FlatID: residentFlatID(filter), UserID: residentUserID(filter),
	})
	return flatResidentFromDBNoRows(row, err)
}

func (r *flatResidentRepository) CountActive(ctx context.Context, societyID int64, flatID int64) (int64, error) {
	return GetQueries(ctx, r.db).CountActiveFlatResidents(ctx, db.CountActiveFlatResidentsParams{SocietyID: societyID, FlatID: flatID})
}

func (r *flatResidentRepository) CountPrimary(ctx context.Context, societyID int64, flatID int64) (int64, error) {
	return GetQueries(ctx, r.db).CountActivePrimaryFlatResidents(ctx, db.CountActivePrimaryFlatResidentsParams{SocietyID: societyID, FlatID: flatID})
}

func flatResidentFromDBNoRows(row db.FlatResident, err error) (*models.FlatResident, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return flatResidentFromDB(row), nil
}

func flatResidentFromDB(row db.FlatResident) *models.FlatResident {
	return &models.FlatResident{
		ID: row.ID, SocietyID: row.SocietyID, FlatID: row.FlatID, UserID: row.UserID,
		Role: models.FlatResidentRole(row.Role), Status: models.FlatResidentStatus(row.Status),
		IsPrimary: row.IsPrimary, MovedInAt: pgTimestamptzToTime(row.MovedInAt),
		MovedOutAt: pgTimestamptzToTimePtr(row.MovedOutAt), Metadata: metadataFromJSON(row.Metadata),
		CreatedBy: row.CreatedBy, CreatedAt: pgTimestamptzToTime(row.CreatedAt), UpdatedAt: pgTimestamptzToTime(row.UpdatedAt),
	}
}

func flatResidentFromGetRow(row db.GetFlatResidentRow) *models.FlatResident {
	userName, flatNumber, societyName, societyCode := row.UserName, row.FlatNumber, row.SocietyName, row.SocietyCode
	resident := flatResidentFromParts(row.ID, row.SocietyID, row.FlatID, row.UserID, row.Role, row.Status, row.IsPrimary, row.MovedInAt, row.MovedOutAt, row.Metadata, row.CreatedBy, row.CreatedAt, row.UpdatedAt)
	resident.UserName = &userName
	resident.UserEmail = row.UserEmail
	resident.UserPhone = row.UserPhone
	resident.FlatNumber = &flatNumber
	resident.Block = row.Block
	resident.Floor = row.Floor
	resident.SocietyName = &societyName
	resident.SocietyCode = &societyCode
	return resident
}

func flatResidentFromListRow(row db.ListFlatResidentsRow) *models.FlatResident {
	userName, flatNumber, societyName, societyCode := row.UserName, row.FlatNumber, row.SocietyName, row.SocietyCode
	resident := flatResidentFromParts(row.ID, row.SocietyID, row.FlatID, row.UserID, row.Role, row.Status, row.IsPrimary, row.MovedInAt, row.MovedOutAt, row.Metadata, row.CreatedBy, row.CreatedAt, row.UpdatedAt)
	resident.UserName = &userName
	resident.UserEmail = row.UserEmail
	resident.UserPhone = row.UserPhone
	resident.FlatNumber = &flatNumber
	resident.Block = row.Block
	resident.Floor = row.Floor
	resident.SocietyName = &societyName
	resident.SocietyCode = &societyCode
	return resident
}

func flatResidentFromParts(id, societyID, flatID, userID int64, role db.FlatResidentRole, status db.FlatResidentStatus, isPrimary bool, movedInAt, movedOutAt pgtype.Timestamptz, rawMetadata []byte, createdBy *int64, createdAt, updatedAt pgtype.Timestamptz) *models.FlatResident {
	return &models.FlatResident{
		ID: id, SocietyID: societyID, FlatID: flatID, UserID: userID,
		Role: models.FlatResidentRole(role), Status: models.FlatResidentStatus(status),
		IsPrimary: isPrimary, MovedInAt: pgTimestamptzToTime(movedInAt), MovedOutAt: pgTimestamptzToTimePtr(movedOutAt),
		Metadata: metadataFromJSON(rawMetadata), CreatedBy: createdBy,
		CreatedAt: pgTimestamptzToTime(createdAt), UpdatedAt: pgTimestamptzToTime(updatedAt),
	}
}

func residentGetParams(filter *models.FlatResidentFilter) db.GetFlatResidentParams {
	return db.GetFlatResidentParams{
		ID: residentID(filter), SocietyID: residentSocietyID(filter), FlatID: residentFlatID(filter),
		UserID: residentUserID(filter), Role: dbResidentRolePtr(residentRole(filter)),
		Status: dbResidentStatusPtr(residentStatus(filter)), IsPrimary: residentIsPrimary(filter),
	}
}

func residentListParams(filter *models.FlatResidentFilter) db.ListFlatResidentsParams {
	return db.ListFlatResidentsParams{
		ID: residentID(filter), SocietyID: residentSocietyID(filter), FlatID: residentFlatID(filter),
		UserID: residentUserID(filter), Role: dbResidentRolePtr(residentRole(filter)),
		Status: dbResidentStatusPtr(residentStatus(filter)), IsPrimary: residentIsPrimary(filter),
		Search:     residentSearch(filter),
		SearchMode: normalizeSearchMode(residentSearchMode(filter), "resident", "society", "flat"),
		Limit:      normalizeLimit(residentLimit(filter)), Offset: normalizeOffset(residentOffset(filter)),
	}
}

func residentID(filter *models.FlatResidentFilter) *int64 {
	if filter == nil {
		return nil
	}
	return filter.ID
}

func residentSocietyID(filter *models.FlatResidentFilter) *int64 {
	if filter == nil {
		return nil
	}
	return filter.SocietyID
}

func residentFlatID(filter *models.FlatResidentFilter) *int64 {
	if filter == nil {
		return nil
	}
	return filter.FlatID
}

func residentUserID(filter *models.FlatResidentFilter) *int64 {
	if filter == nil {
		return nil
	}
	return filter.UserID
}

func residentRole(filter *models.FlatResidentFilter) *string {
	if filter == nil {
		return nil
	}
	return filter.Role
}

func residentStatus(filter *models.FlatResidentFilter) *string {
	if filter == nil {
		return nil
	}
	return filter.Status
}

func residentIsPrimary(filter *models.FlatResidentFilter) *bool {
	if filter == nil {
		return nil
	}
	return filter.IsPrimary
}

func residentSearch(filter *models.FlatResidentFilter) string {
	if filter == nil {
		return ""
	}
	return filter.Search
}

func residentSearchMode(filter *models.FlatResidentFilter) string {
	if filter == nil {
		return ""
	}
	return filter.SearchMode
}

func residentLimit(filter *models.FlatResidentFilter) int32 {
	if filter == nil {
		return 0
	}
	return filter.Limit
}

func residentOffset(filter *models.FlatResidentFilter) int32 {
	if filter == nil {
		return 0
	}
	return filter.Offset
}

func dbResidentRolePtr(role *string) *db.FlatResidentRole {
	if role == nil || *role == "" {
		return nil
	}
	value := db.FlatResidentRole(*role)
	return &value
}

func dbResidentStatusPtr(status *string) *db.FlatResidentStatus {
	if status == nil || *status == "" {
		return nil
	}
	value := db.FlatResidentStatus(*status)
	return &value
}
