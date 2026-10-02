package repository

import (
	"context"
	"encoding/json"
	"errors"
	"go-server/internal/db"
	"go-server/internal/models"
	"go-server/internal/repositories/contracts"
	"go-server/pkg/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type SocietyMemberRepository = contracts.SocietyMemberRepository

type societyMemberRepository struct {
	db *database.Database
}

func NewSocietyMemberRepository(db *database.Database) SocietyMemberRepository {
	return &societyMemberRepository{db: db}
}

func (r *societyMemberRepository) Add(ctx context.Context, member *models.SocietyMember) error {
	metadata, err := json.Marshal(member.Metadata)
	if err != nil {
		return err
	}
	if string(metadata) == "null" {
		metadata = []byte("{}")
	}
	row, err := GetQueries(ctx, r.db).AddSocietyMember(ctx, db.AddSocietyMemberParams{
		SocietyID: member.SocietyID, UserID: member.UserID, Role: db.SocietyMemberRole(member.Role),
		InvitedBy: member.InvitedBy, Metadata: metadata,
	})
	if err != nil {
		return err
	}
	*member = *societyMemberFromDB(row)
	return nil
}

func (r *societyMemberRepository) Get(ctx context.Context, filter models.GetSocietyMemberFilter) (*models.SocietyMember, error) {
	row, err := GetQueries(ctx, r.db).GetSocietyMember(ctx, db.GetSocietyMemberParams{
		ID: filter.ID, SocietyID: filter.SocietyID, UserID: filter.UserID,
		Role: dbSocietyMemberRolePtr(filter.Role), Status: dbSocietyMemberStatusPtr(filter.Status),
		Email: filter.Email, Phone: filter.Phone,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return societyMemberFromGetRow(row), nil
}

func (r *societyMemberRepository) List(ctx context.Context, filter models.ListSocietyMembersFilter) ([]*models.SocietyMember, error) {
	params := listMembersParams(filter)
	rows, err := GetQueries(ctx, r.db).ListSocietyMembers(ctx, params)
	if err != nil {
		return nil, err
	}
	result := make([]*models.SocietyMember, 0, len(rows))
	for _, row := range rows {
		result = append(result, societyMemberFromListRow(row))
	}
	return result, nil
}

func (r *societyMemberRepository) ListByUser(ctx context.Context, userID int64) ([]*models.SocietyMember, error) {
	rows, err := GetQueries(ctx, r.db).ListSocietyMembershipsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]*models.SocietyMember, 0, len(rows))
	for _, row := range rows {
		result = append(result, societyMemberFromListMyRow(row))
	}
	return result, nil
}

func (r *societyMemberRepository) ListMySocietiesByUser(ctx context.Context, userID int64) ([]*models.MySocietyResponse, error) {
	rows, err := GetQueries(ctx, r.db).ListMySocietiesByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]*models.MySocietyResponse, 0, len(rows))
	for _, row := range rows {
		member := societyMemberFromListMySocietiesRow(row)
		society := societyFromListMySocietiesRow(row)
		result = append(result, &models.MySocietyResponse{
			Society: society.ToResponse(),
			Member:  member.ToResponse(),
		})
	}
	return result, nil
}

func (r *societyMemberRepository) Count(ctx context.Context, filter models.ListSocietyMembersFilter) (int64, error) {
	params := listMembersParams(filter)
	return GetQueries(ctx, r.db).CountSocietyMembers(ctx, db.CountSocietyMembersParams{
		SocietyID: params.SocietyID, Role: params.Role, Status: params.Status, UserID: params.UserID,
		InvitedBy: params.InvitedBy, RemovedBy: params.RemovedBy, JoinedFrom: params.JoinedFrom,
		JoinedTo: params.JoinedTo, Search: params.Search, SearchMode: params.SearchMode,
	})
}

func (r *societyMemberRepository) ChangeRole(ctx context.Context, societyID int64, userID int64, role models.SocietyMemberRole) (*models.SocietyMember, error) {
	row, err := GetQueries(ctx, r.db).ChangeSocietyMemberRole(ctx, db.ChangeSocietyMemberRoleParams{
		SocietyID: societyID, UserID: userID, Role: db.SocietyMemberRole(role),
	})
	return societyMemberFromDBNoRows(row, err)
}

func (r *societyMemberRepository) Suspend(ctx context.Context, societyID int64, userID int64) (*models.SocietyMember, error) {
	row, err := GetQueries(ctx, r.db).SuspendSocietyMember(ctx, db.SuspendSocietyMemberParams{SocietyID: societyID, UserID: userID})
	return societyMemberFromDBNoRows(row, err)
}

func (r *societyMemberRepository) Reactivate(ctx context.Context, societyID int64, userID int64) (*models.SocietyMember, error) {
	row, err := GetQueries(ctx, r.db).ReactivateSocietyMember(ctx, db.ReactivateSocietyMemberParams{SocietyID: societyID, UserID: userID})
	return societyMemberFromDBNoRows(row, err)
}

func (r *societyMemberRepository) Remove(ctx context.Context, societyID int64, userID int64, removedBy int64, reason string) error {
	return GetQueries(ctx, r.db).RemoveSocietyMember(ctx, db.RemoveSocietyMemberParams{
		SocietyID: societyID, UserID: userID, RemovedBy: &removedBy, RemoveReason: nullableString(reason),
	})
}

func (r *societyMemberRepository) CountActiveOwners(ctx context.Context, societyID int64) (int64, error) {
	return GetQueries(ctx, r.db).CountActiveOwners(ctx, societyID)
}

func (r *societyMemberRepository) DemoteActiveOwners(ctx context.Context, societyID int64, exceptUserID int64) error {
	return GetQueries(ctx, r.db).DemoteActiveOwners(ctx, db.DemoteActiveOwnersParams{SocietyID: societyID, UserID: exceptUserID})
}

func (r *societyMemberRepository) PromoteToOwner(ctx context.Context, societyID int64, userID int64) (*models.SocietyMember, error) {
	row, err := GetQueries(ctx, r.db).PromoteMemberToOwner(ctx, db.PromoteMemberToOwnerParams{SocietyID: societyID, UserID: userID})
	return societyMemberFromDBNoRows(row, err)
}

func (r *societyMemberRepository) UpsertResident(ctx context.Context, societyID int64, userID int64, invitedBy int64) (*models.SocietyMember, error) {
	row, err := GetQueries(ctx, r.db).UpsertResidentSocietyMember(ctx, db.UpsertResidentSocietyMemberParams{
		SocietyID: societyID, UserID: userID, InvitedBy: &invitedBy,
	})
	return societyMemberFromDBNoRows(row, err)
}

func societyMemberFromDBNoRows(row db.SocietyMember, err error) (*models.SocietyMember, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return societyMemberFromDB(row), nil
}

func societyMemberFromDB(row db.SocietyMember) *models.SocietyMember {
	metadata := map[string]any{}
	if len(row.Metadata) > 0 {
		_ = json.Unmarshal(row.Metadata, &metadata)
	}
	return &models.SocietyMember{
		ID: row.ID, SocietyID: row.SocietyID, UserID: row.UserID, Role: models.SocietyMemberRole(row.Role),
		Status: models.SocietyMemberStatus(row.Status), InvitedBy: row.InvitedBy, JoinedAt: pgTimestamptzToTime(row.JoinedAt),
		RemovedBy: row.RemovedBy, RemovedAt: pgTimestamptzToTimePtr(row.RemovedAt), RemoveReason: row.RemoveReason,
		Metadata: metadata, CreatedAt: pgTimestamptzToTime(row.CreatedAt), UpdatedAt: pgTimestamptzToTime(row.UpdatedAt),
	}
}

func societyMemberFromGetRow(row db.GetSocietyMemberRow) *models.SocietyMember {
	member := societyMemberFromParts(row.ID, row.SocietyID, row.UserID, string(row.Role), string(row.Status), row.InvitedBy, row.JoinedAt, row.RemovedBy, row.RemovedAt, row.RemoveReason, row.Metadata, row.CreatedAt, row.UpdatedAt)
	member.UserFullName = &row.UserFullName
	member.UserEmail = row.UserEmail
	member.UserPhone = row.UserPhone
	return member
}

func societyMemberFromListRow(row db.ListSocietyMembersRow) *models.SocietyMember {
	member := societyMemberFromParts(row.ID, row.SocietyID, row.UserID, string(row.Role), string(row.Status), row.InvitedBy, row.JoinedAt, row.RemovedBy, row.RemovedAt, row.RemoveReason, row.Metadata, row.CreatedAt, row.UpdatedAt)
	member.UserFullName = &row.UserFullName
	member.UserEmail = row.UserEmail
	member.UserPhone = row.UserPhone
	return member
}

func societyMemberFromListMyRow(row db.ListSocietyMembershipsByUserRow) *models.SocietyMember {
	member := societyMemberFromParts(row.ID, row.SocietyID, row.UserID, string(row.Role), string(row.Status), row.InvitedBy, row.JoinedAt, row.RemovedBy, row.RemovedAt, row.RemoveReason, row.Metadata, row.CreatedAt, row.UpdatedAt)
	member.UserFullName = &row.UserFullName
	member.UserEmail = row.UserEmail
	member.UserPhone = row.UserPhone
	return member
}

func societyMemberFromListMySocietiesRow(row db.ListMySocietiesByUserRow) *models.SocietyMember {
	member := societyMemberFromParts(row.ID, row.SocietyID, row.UserID, string(row.Role), string(row.Status), row.InvitedBy, row.JoinedAt, row.RemovedBy, row.RemovedAt, row.RemoveReason, row.Metadata, row.CreatedAt, row.UpdatedAt)
	member.UserFullName = &row.UserFullName
	member.UserEmail = row.UserEmail
	member.UserPhone = row.UserPhone
	return member
}

func societyFromListMySocietiesRow(row db.ListMySocietiesByUserRow) *models.Society {
	metadata := map[string]any{}
	if len(row.SocietyMetadata) > 0 {
		_ = json.Unmarshal(row.SocietyMetadata, &metadata)
	}
	return &models.Society{
		ID: row.SocietyRowID, Name: row.SocietyName, SocietyCode: row.SocietyCode, Email: row.SocietyEmail,
		PhoneNumber: row.SocietyPhoneNumber, AddressLine1: row.SocietyAddressLine1, AddressLine2: row.SocietyAddressLine2,
		Landmark: row.SocietyLandmark, City: row.SocietyCity, State: row.SocietyState, Pincode: row.SocietyPincode,
		Country: row.SocietyCountry, TotalFlats: row.SocietyTotalFlats, TotalBlocks: row.SocietyTotalBlocks,
		Status: models.SocietyStatus(row.SocietyStatus), CreatedBy: row.SocietyCreatedBy, ApprovedBy: row.SocietyApprovedBy,
		ApprovedAt: pgTimestamptzToTimePtr(row.SocietyApprovedAt), RejectedBy: row.SocietyRejectedBy,
		RejectedAt: pgTimestamptzToTimePtr(row.SocietyRejectedAt), RejectionReason: row.SocietyRejectionReason,
		SuspendedBy: row.SocietySuspendedBy, SuspendedAt: pgTimestamptzToTimePtr(row.SocietySuspendedAt),
		SuspensionReason: row.SocietySuspensionReason, Metadata: metadata,
		CreatedAt: pgTimestamptzToTime(row.SocietyCreatedAt), UpdatedAt: pgTimestamptzToTime(row.SocietyUpdatedAt),
		DeletedAt: pgTimestamptzToTimePtr(row.SocietyDeletedAt),
	}
}

func societyMemberFromParts(id, societyID, userID int64, role, status string, invitedBy *int64, joinedAt pgtype.Timestamptz, removedBy *int64, removedAt pgtype.Timestamptz, removeReason *string, rawMetadata []byte, createdAt, updatedAt pgtype.Timestamptz) *models.SocietyMember {
	metadata := map[string]any{}
	if len(rawMetadata) > 0 {
		_ = json.Unmarshal(rawMetadata, &metadata)
	}
	return &models.SocietyMember{
		ID: id, SocietyID: societyID, UserID: userID, Role: models.SocietyMemberRole(role),
		Status: models.SocietyMemberStatus(status), InvitedBy: invitedBy, JoinedAt: pgTimestamptzToTime(joinedAt),
		RemovedBy: removedBy, RemovedAt: pgTimestamptzToTimePtr(removedAt), RemoveReason: removeReason,
		Metadata: metadata, CreatedAt: pgTimestamptzToTime(createdAt), UpdatedAt: pgTimestamptzToTime(updatedAt),
	}
}

func listMembersParams(filter models.ListSocietyMembersFilter) db.ListSocietyMembersParams {
	return db.ListSocietyMembersParams{
		SocietyID: filter.SocietyID, Role: dbSocietyMemberRolePtr(filter.Role), Status: dbSocietyMemberStatusPtr(filter.Status),
		UserID: filter.UserID, InvitedBy: filter.InvitedBy, RemovedBy: filter.RemovedBy,
		JoinedFrom: timePtrToPgTimestamptz(filter.JoinedFrom), JoinedTo: timePtrToPgTimestamptz(filter.JoinedTo),
		Search:     filter.Search,
		SearchMode: normalizeSearchMode(filter.SearchMode, "resident", "invited_by", "removed_by"),
		SortBy:     normalizeMemberSort(filter.SortBy), SortOrder: normalizeSortOrder(filter.SortOrder),
		Limit: normalizeLimit(filter.Limit), Offset: normalizeOffset(filter.Offset),
	}
}

func dbSocietyMemberRolePtr(role *string) *db.SocietyMemberRole {
	if role == nil || *role == "" {
		return nil
	}
	value := db.SocietyMemberRole(*role)
	return &value
}

func dbSocietyMemberStatusPtr(status *string) *db.SocietyMemberStatus {
	if status == nil || *status == "" {
		return nil
	}
	value := db.SocietyMemberStatus(*status)
	return &value
}

func normalizeMemberSort(sortBy string) string {
	switch sortBy {
	case "role", "status", "joined_at":
		return sortBy
	default:
		return "joined_at"
	}
}
