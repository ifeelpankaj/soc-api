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

type FlatClaimRepository = contracts.FlatClaimRepository

type flatClaimRepository struct {
	db *database.Database
}

func NewFlatClaimRepository(db *database.Database) FlatClaimRepository {
	return &flatClaimRepository{db: db}
}

func (r *flatClaimRepository) Submit(ctx context.Context, claim *models.FlatClaim) error {
	metadata, err := jsonMap(claim.Metadata)
	if err != nil {
		return err
	}
	row, err := GetQueries(ctx, r.db).SubmitFlatClaim(ctx, db.SubmitFlatClaimParams{
		SocietyID: claim.SocietyID, FlatID: claim.FlatID, UserID: claim.UserID,
		RequestedRole: db.FlatResidentRole(claim.RequestedRole), RequestedPrimary: claim.RequestedPrimary,
		Note: claim.Note, Metadata: metadata,
	})
	if err != nil {
		return err
	}
	*claim = *flatClaimFromDB(row)
	return nil
}

func (r *flatClaimRepository) Get(ctx context.Context, filter *models.FlatClaimFilter) (*models.FlatClaim, error) {
	row, err := GetQueries(ctx, r.db).GetFlatClaim(ctx, claimGetParams(filter))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return flatClaimFromGetRow(row), nil
}

func (r *flatClaimRepository) List(ctx context.Context, filter *models.FlatClaimFilter) ([]*models.FlatClaim, error) {
	rows, err := GetQueries(ctx, r.db).ListFlatClaims(ctx, claimListParams(filter))
	if err != nil {
		return nil, err
	}
	items := make([]*models.FlatClaim, 0, len(rows))
	for _, row := range rows {
		items = append(items, flatClaimFromListRow(row))
	}
	return items, nil
}

func (r *flatClaimRepository) Stats(ctx context.Context, societyID int64) (*models.FlatClaimStatsResponse, error) {
	row, err := GetQueries(ctx, r.db).GetFlatClaimStats(ctx, societyID)
	if err != nil {
		return nil, err
	}
	return &models.FlatClaimStatsResponse{
		TotalClaims:     row.TotalClaims,
		PendingClaims:   row.PendingClaims,
		ApprovedClaims:  row.ApprovedClaims,
		RejectedClaims:  row.RejectedClaims,
		CancelledClaims: row.CancelledClaims,
	}, nil
}

func (r *flatClaimRepository) Approve(ctx context.Context, societyID int64, claimID int64, reviewedBy int64) (*models.FlatClaim, error) {
	row, err := GetQueries(ctx, r.db).ApproveFlatClaim(ctx, db.ApproveFlatClaimParams{ID: claimID, SocietyID: societyID, ReviewedBy: &reviewedBy})
	return flatClaimFromDBNoRows(row, err)
}

func (r *flatClaimRepository) Reject(ctx context.Context, societyID int64, claimID int64, reviewedBy int64, reason string) (*models.FlatClaim, error) {
	row, err := GetQueries(ctx, r.db).RejectFlatClaim(ctx, db.RejectFlatClaimParams{ID: claimID, SocietyID: societyID, ReviewedBy: &reviewedBy, RejectionReason: nullableString(reason)})
	return flatClaimFromDBNoRows(row, err)
}

func (r *flatClaimRepository) Cancel(ctx context.Context, claimID int64, userID int64) (*models.FlatClaim, error) {
	row, err := GetQueries(ctx, r.db).CancelMyFlatClaim(ctx, db.CancelMyFlatClaimParams{ID: claimID, UserID: userID})
	return flatClaimFromDBNoRows(row, err)
}

func flatClaimFromDBNoRows(row db.FlatClaimRequest, err error) (*models.FlatClaim, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return flatClaimFromDB(row), nil
}

func flatClaimFromDB(row db.FlatClaimRequest) *models.FlatClaim {
	return &models.FlatClaim{
		ID: row.ID, SocietyID: row.SocietyID, FlatID: row.FlatID, UserID: row.UserID,
		RequestedRole: models.FlatResidentRole(row.RequestedRole), RequestedPrimary: row.RequestedPrimary,
		Status: models.FlatClaimStatus(row.Status), Note: row.Note, RejectionReason: row.RejectionReason,
		ReviewedBy: row.ReviewedBy, ReviewedAt: pgTimestamptzToTimePtr(row.ReviewedAt),
		CancelledAt: pgTimestamptzToTimePtr(row.CancelledAt), Metadata: metadataFromJSON(row.Metadata),
		CreatedAt: pgTimestamptzToTime(row.CreatedAt), UpdatedAt: pgTimestamptzToTime(row.UpdatedAt),
	}
}

func flatClaimFromGetRow(row db.GetFlatClaimRow) *models.FlatClaim {
	return flatClaimFromParts(
		row.ID, row.SocietyID, row.FlatID, row.UserID, row.RequestedRole, row.RequestedPrimary,
		row.Status, row.Note, row.RejectionReason, row.ReviewedBy, row.ReviewedAt, row.CancelledAt,
		row.Metadata, row.CreatedAt, row.UpdatedAt, row.UserName, row.UserEmail, row.UserPhone,
		row.FlatNumber, row.Block, row.Floor, row.FlatStatus, row.SocietyName, row.SocietyCode,
		row.ReviewerName, row.ReviewerEmail, row.ReviewerPhone,
	)
}

func flatClaimFromListRow(row db.ListFlatClaimsRow) *models.FlatClaim {
	return flatClaimFromParts(
		row.ID, row.SocietyID, row.FlatID, row.UserID, row.RequestedRole, row.RequestedPrimary,
		row.Status, row.Note, row.RejectionReason, row.ReviewedBy, row.ReviewedAt, row.CancelledAt,
		row.Metadata, row.CreatedAt, row.UpdatedAt, row.UserName, row.UserEmail, row.UserPhone,
		row.FlatNumber, row.Block, row.Floor, row.FlatStatus, row.SocietyName, row.SocietyCode,
		row.ReviewerName, row.ReviewerEmail, row.ReviewerPhone,
	)
}

func flatClaimFromParts(id, societyID, flatID, userID int64, role db.FlatResidentRole, requestedPrimary bool, status db.FlatClaimStatus, note, rejectionReason *string, reviewedBy *int64, reviewedAt, cancelledAt pgtype.Timestamptz, rawMetadata []byte, createdAt, updatedAt pgtype.Timestamptz, userName string, userEmail, userPhone *string, flatNumber string, block, floor *string, flatStatus db.FlatStatus, societyName, societyCode string, reviewerName, reviewerEmail, reviewerPhone *string) *models.FlatClaim {
	userNameCopy, flatNumberCopy, societyNameCopy, societyCodeCopy := userName, flatNumber, societyName, societyCode
	modelFlatStatus := models.FlatStatus(flatStatus)
	return &models.FlatClaim{
		ID: id, SocietyID: societyID, FlatID: flatID, UserID: userID,
		RequestedRole: models.FlatResidentRole(role), RequestedPrimary: requestedPrimary,
		Status: models.FlatClaimStatus(status), Note: note, RejectionReason: rejectionReason,
		ReviewedBy: reviewedBy, ReviewedAt: pgTimestamptzToTimePtr(reviewedAt), CancelledAt: pgTimestamptzToTimePtr(cancelledAt),
		Metadata: metadataFromJSON(rawMetadata), CreatedAt: pgTimestamptzToTime(createdAt), UpdatedAt: pgTimestamptzToTime(updatedAt),
		UserName: &userNameCopy, UserEmail: userEmail, UserPhone: userPhone, FlatNumber: &flatNumberCopy,
		Block: block, Floor: floor, FlatStatus: &modelFlatStatus, SocietyName: &societyNameCopy, SocietyCode: &societyCodeCopy,
		ReviewerName: reviewerName, ReviewerEmail: reviewerEmail, ReviewerPhone: reviewerPhone,
	}
}

func claimGetParams(filter *models.FlatClaimFilter) db.GetFlatClaimParams {
	return db.GetFlatClaimParams{
		ID: claimID(filter), SocietyID: claimSocietyID(filter), FlatID: claimFlatID(filter),
		UserID: claimUserID(filter), Status: dbClaimStatusPtr(claimStatus(filter)),
	}
}

func claimListParams(filter *models.FlatClaimFilter) db.ListFlatClaimsParams {
	return db.ListFlatClaimsParams{
		ID: claimID(filter), SocietyID: claimSocietyID(filter), FlatID: claimFlatID(filter),
		UserID: claimUserID(filter), Status: dbClaimStatusPtr(claimStatus(filter)),
		Search:     claimSearch(filter),
		SearchMode: normalizeSearchMode(claimSearchMode(filter), "claimant", "flat", "reviewer"),
		Limit:      normalizeLimit(claimLimit(filter)), Offset: normalizeOffset(claimOffset(filter)),
	}
}

func claimID(filter *models.FlatClaimFilter) *int64 {
	if filter == nil {
		return nil
	}
	return filter.ID
}

func claimSocietyID(filter *models.FlatClaimFilter) *int64 {
	if filter == nil {
		return nil
	}
	return filter.SocietyID
}

func claimFlatID(filter *models.FlatClaimFilter) *int64 {
	if filter == nil {
		return nil
	}
	return filter.FlatID
}

func claimUserID(filter *models.FlatClaimFilter) *int64 {
	if filter == nil {
		return nil
	}
	return filter.UserID
}

func claimStatus(filter *models.FlatClaimFilter) *string {
	if filter == nil {
		return nil
	}
	return filter.Status
}

func claimSearch(filter *models.FlatClaimFilter) string {
	if filter == nil {
		return ""
	}
	return filter.Search
}

func claimSearchMode(filter *models.FlatClaimFilter) string {
	if filter == nil {
		return ""
	}
	return filter.SearchMode
}

func claimLimit(filter *models.FlatClaimFilter) int32 {
	if filter == nil {
		return 0
	}
	return filter.Limit
}

func claimOffset(filter *models.FlatClaimFilter) int32 {
	if filter == nil {
		return 0
	}
	return filter.Offset
}

func dbClaimStatusPtr(status *string) *db.FlatClaimStatus {
	if status == nil || *status == "" {
		return nil
	}
	value := db.FlatClaimStatus(*status)
	return &value
}
