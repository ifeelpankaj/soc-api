package repository

import (
	"context"
	"encoding/json"
	"errors"
	"go-server/internal/repositories/contracts"

	"go-server/internal/db"
	"go-server/internal/models"
	"go-server/pkg/database"

	"github.com/jackc/pgx/v5"
)

type FlatRepository = contracts.FlatRepository

type flatRepository struct {
	db *database.Database
}

func NewFlatRepository(db *database.Database) FlatRepository {
	return &flatRepository{db: db}
}

func (r *flatRepository) Create(ctx context.Context, flat *models.Flat) error {
	area, err := models.FlatAreaValue(flat.AreaSqft)
	if err != nil {
		return err
	}
	metadata, err := jsonMap(flat.Metadata)
	if err != nil {
		return err
	}
	row, err := GetQueries(ctx, r.db).CreateFlat(ctx, db.CreateFlatParams{
		FlatType: flat.FlatType, AreaSqftHundredths: area, SocietyID: flat.SocietyID, Block: flat.Block, Floor: flat.Floor,
		FlatNumber: flat.FlatNumber, CreatedBy: flat.CreatedBy, Metadata: metadata,
	})
	if err != nil {
		return err
	}
	*flat = *flatFromDB(row)
	return nil
}

func (r *flatRepository) Get(ctx context.Context, filter *models.FlatFilter) (*models.Flat, error) {
	row, err := GetQueries(ctx, r.db).GetFlat(ctx, flatGetParams(filter))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return flatFromGetRow(row), nil
}

func (r *flatRepository) List(ctx context.Context, filter *models.FlatFilter) ([]*models.Flat, error) {
	params := flatListParams(filter)
	rows, err := GetQueries(ctx, r.db).ListFlats(ctx, params)
	if err != nil {
		return nil, err
	}
	items := make([]*models.Flat, 0, len(rows))
	for _, row := range rows {
		items = append(items, flatFromListRow(row))
	}
	return items, nil
}

func (r *flatRepository) Count(ctx context.Context, filter *models.FlatFilter) (int64, error) {
	params := flatCountParams(filter)
	return GetQueries(ctx, r.db).CountFlats(ctx, params)
}

func (r *flatRepository) Stats(ctx context.Context, societyID int64) (*models.FlatStatsResponse, error) {
	row, err := GetQueries(ctx, r.db).GetFlatStats(ctx, societyID)
	if err != nil {
		return nil, err
	}
	return &models.FlatStatsResponse{
		SocietyID: row.SocietyID, TotalFlats: row.TotalFlats, VacantFlats: row.VacantFlats,
		OccupiedFlats: row.OccupiedFlats, BlockedFlats: row.BlockedFlats,
		ActiveFlats: row.ActiveFlats, InactiveFlats: row.InactiveFlats,
	}, nil
}

func (r *flatRepository) Update(ctx context.Context, filter *models.FlatFilter, req *contracts.UpdateFlatInput) (*models.Flat, error) {
	area, err := models.FlatAreaValue(req.AreaSqft)
	if err != nil {
		return nil, err
	}
	var metadata []byte
	if req.Metadata != nil {
		var err error
		metadata, err = jsonMap(req.Metadata)
		if err != nil {
			return nil, err
		}
	}
	row, err := GetQueries(ctx, r.db).UpdateFlat(ctx, db.UpdateFlatParams{
		FlatType: req.FlatType, AreaSqftHundredths: area, Block: req.Block, Floor: req.Floor, FlatNumber: req.FlatNumber,
		IsActive: req.IsActive, Metadata: metadata,
		ID: idFromFlatFilter(filter), SocietyID: societyIDFromFlatFilter(filter),
	})
	return flatFromDBNoRows(row, err)
}

func (r *flatRepository) Deactivate(ctx context.Context, filter *models.FlatFilter) error {
	return GetQueries(ctx, r.db).DeactivateFlat(ctx, db.DeactivateFlatParams{
		ID: idFromFlatFilter(filter), SocietyID: societyIDFromFlatFilter(filter),
	})
}

func (r *flatRepository) Block(ctx context.Context, filter *models.FlatFilter) (*models.Flat, error) {
	row, err := GetQueries(ctx, r.db).BlockFlat(ctx, db.BlockFlatParams{ID: idFromFlatFilter(filter), SocietyID: societyIDFromFlatFilter(filter)})
	return flatFromDBNoRows(row, err)
}

func (r *flatRepository) Unblock(ctx context.Context, filter *models.FlatFilter) (*models.Flat, error) {
	row, err := GetQueries(ctx, r.db).UnblockFlat(ctx, db.UnblockFlatParams{ID: idFromFlatFilter(filter), SocietyID: societyIDFromFlatFilter(filter)})
	return flatFromDBNoRows(row, err)
}

func (r *flatRepository) MarkOccupied(ctx context.Context, societyID int64, flatID int64) (*models.Flat, error) {
	row, err := GetQueries(ctx, r.db).MarkFlatOccupied(ctx, db.MarkFlatOccupiedParams{ID: flatID, SocietyID: societyID})
	return flatFromDBNoRows(row, err)
}

func (r *flatRepository) MarkVacant(ctx context.Context, societyID int64, flatID int64) (*models.Flat, error) {
	row, err := GetQueries(ctx, r.db).MarkFlatVacant(ctx, db.MarkFlatVacantParams{ID: flatID, SocietyID: societyID})
	return flatFromDBNoRows(row, err)
}

func flatFromDBNoRows(row db.Flat, err error) (*models.Flat, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return flatFromDB(row), nil
}

func flatFromDB(row db.Flat) *models.Flat {
	return &models.Flat{
		FlatType: row.FlatType, AreaSqft: models.FlatAreaString(row.AreaSqftHundredths), ID: row.ID, SocietyID: row.SocietyID, Block: row.Block, Floor: row.Floor,
		FlatNumber: row.FlatNumber, Status: models.FlatStatus(row.Status), IsActive: row.IsActive,
		Metadata: metadataFromJSON(row.Metadata), CreatedBy: row.CreatedBy,
		CreatedAt: pgTimestamptzToTime(row.CreatedAt), UpdatedAt: pgTimestamptzToTime(row.UpdatedAt),
	}
}

func flatFromGetRow(row db.GetFlatRow) *models.Flat {
	societyName, societyCode := row.SocietyName, row.SocietyCode
	return &models.Flat{
		FlatType: row.FlatType, AreaSqft: models.FlatAreaString(row.AreaSqftHundredths), ID: row.ID, SocietyID: row.SocietyID, Block: row.Block, Floor: row.Floor,
		FlatNumber: row.FlatNumber, Status: models.FlatStatus(row.Status), IsActive: row.IsActive,
		Metadata: metadataFromJSON(row.Metadata), CreatedBy: row.CreatedBy,
		CreatedAt: pgTimestamptzToTime(row.CreatedAt), UpdatedAt: pgTimestamptzToTime(row.UpdatedAt),
		SocietyName: &societyName, SocietyCode: &societyCode,
	}
}

func flatFromListRow(row db.ListFlatsRow) *models.Flat {
	societyName, societyCode := row.SocietyName, row.SocietyCode
	var primaryResidentName *string
	if row.PrimaryResidentName != "" {
		name := row.PrimaryResidentName
		primaryResidentName = &name
	}
	return &models.Flat{
		FlatType: row.FlatType, AreaSqft: models.FlatAreaString(row.AreaSqftHundredths), ID: row.ID, SocietyID: row.SocietyID, Block: row.Block, Floor: row.Floor,
		FlatNumber: row.FlatNumber, Status: models.FlatStatus(row.Status), IsActive: row.IsActive,
		Metadata: metadataFromJSON(row.Metadata), CreatedBy: row.CreatedBy,
		CreatedAt: pgTimestamptzToTime(row.CreatedAt), UpdatedAt: pgTimestamptzToTime(row.UpdatedAt),
		SocietyName: &societyName, SocietyCode: &societyCode, PrimaryResidentName: primaryResidentName,
	}
}

func flatGetParams(filter *models.FlatFilter) db.GetFlatParams {
	return db.GetFlatParams{
		ID: idFromFlatFilter(filter), SocietyID: societyIDFromFlatFilter(filter),
		Block: stringPtrFromFlatFilter(filter, "block"), Floor: stringPtrFromFlatFilter(filter, "floor"),
		FlatNumber: stringPtrFromFlatFilter(filter, "flat_number"), Status: dbFlatStatusPtr(filterStatus(filter)),
		IsActive: isActiveFromFlatFilter(filter),
	}
}

func flatListParams(filter *models.FlatFilter) db.ListFlatsParams {
	return db.ListFlatsParams{
		ID: idFromFlatFilter(filter), SocietyID: societyIDFromFlatFilter(filter),
		Block: stringPtrFromFlatFilter(filter, "block"), Floor: stringPtrFromFlatFilter(filter, "floor"),
		FlatNumber: stringPtrFromFlatFilter(filter, "flat_number"), Status: dbFlatStatusPtr(filterStatus(filter)),
		IsActive: isActiveFromFlatFilter(filter), Search: searchFromFlatFilter(filter),
		Limit: normalizeLimit(limitFromFlatFilter(filter)), Offset: normalizeOffset(offsetFromFlatFilter(filter)),
	}
}

func flatCountParams(filter *models.FlatFilter) db.CountFlatsParams {
	return db.CountFlatsParams{
		ID: idFromFlatFilter(filter), SocietyID: societyIDFromFlatFilter(filter),
		Block: stringPtrFromFlatFilter(filter, "block"), Floor: stringPtrFromFlatFilter(filter, "floor"),
		FlatNumber: stringPtrFromFlatFilter(filter, "flat_number"), Status: dbFlatStatusPtr(filterStatus(filter)),
		IsActive: isActiveFromFlatFilter(filter), Search: searchFromFlatFilter(filter),
	}
}

func idFromFlatFilter(filter *models.FlatFilter) *int64 {
	if filter == nil {
		return nil
	}
	return filter.ID
}

func societyIDFromFlatFilter(filter *models.FlatFilter) *int64 {
	if filter == nil {
		return nil
	}
	return filter.SocietyID
}

func stringPtrFromFlatFilter(filter *models.FlatFilter, name string) *string {
	if filter == nil {
		return nil
	}
	switch name {
	case "block":
		return filter.Block
	case "floor":
		return filter.Floor
	case "flat_number":
		return filter.FlatNumber
	default:
		return nil
	}
}

func filterStatus(filter *models.FlatFilter) *string {
	if filter == nil {
		return nil
	}
	return filter.Status
}

func isActiveFromFlatFilter(filter *models.FlatFilter) *bool {
	if filter == nil {
		return nil
	}
	return filter.IsActive
}

func limitFromFlatFilter(filter *models.FlatFilter) int32 {
	if filter == nil {
		return 0
	}
	return filter.Limit
}

func offsetFromFlatFilter(filter *models.FlatFilter) int32 {
	if filter == nil {
		return 0
	}
	return filter.Offset
}

func searchFromFlatFilter(filter *models.FlatFilter) string {
	if filter == nil {
		return ""
	}
	return filter.Search
}

func dbFlatStatusPtr(status *string) *db.FlatStatus {
	if status == nil || *status == "" {
		return nil
	}
	value := db.FlatStatus(*status)
	return &value
}

func dbFlatStatusPtrFromModel(status *models.FlatStatus) *db.FlatStatus {
	if status == nil {
		return nil
	}
	value := db.FlatStatus(*status)
	return &value
}

func jsonMap(value map[string]any) ([]byte, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	if string(raw) == "null" {
		return []byte("{}"), nil
	}
	return raw, nil
}

func metadataFromJSON(raw []byte) map[string]any {
	metadata := map[string]any{}
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &metadata)
	}
	return metadata
}
