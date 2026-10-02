package repository

import (
	"context"
	"go-server/internal/db"
	"go-server/internal/repositories/contracts"
)

func (r *HubRepository) HubReport(ctx context.Context, arg contracts.HubReportInput) (contracts.ChannelReport, error) {
	value, err := GetQueries(ctx, r.database).HubReport(ctx, db.HubReportParams{PostID: arg.PostID, ReporterID: arg.ReporterID, Reason: arg.Reason})
	return mapChannelReport(value), persistenceError(err)
}

func (r *HubRepository) HubReportByID(ctx context.Context, arg contracts.HubReportByIDInput) (contracts.ChannelReport, error) {
	value, err := GetQueries(ctx, r.database).HubReportByID(ctx, db.HubReportByIDParams{SocietyID: arg.SocietyID, ID: arg.ID})
	return mapChannelReport(value), persistenceError(err)
}

func (r *HubRepository) HubReports(ctx context.Context, arg contracts.HubReportsInput) ([]contracts.ChannelReport, error) {
	value, err := GetQueries(ctx, r.database).HubReports(ctx, db.HubReportsParams{SocietyID: arg.SocietyID, Status: arg.Status, AfterID: arg.AfterID, PageLimit: arg.PageLimit})
	if err != nil {
		return nil, persistenceError(err)
	}
	out := make([]contracts.ChannelReport, len(value))
	for i := range value {
		out[i] = mapChannelReport(value[i])
	}
	return out, nil
}

func (r *HubRepository) HubResolvePostReports(ctx context.Context, arg contracts.HubResolvePostReportsInput) error {
	return persistenceError(GetQueries(ctx, r.database).HubResolvePostReports(ctx, db.HubResolvePostReportsParams{PostID: arg.PostID, ReviewerID: arg.ReviewerID, ResolutionNote: arg.ResolutionNote}))
}

func (r *HubRepository) HubResolveReport(ctx context.Context, arg contracts.HubResolveReportInput) error {
	return persistenceError(GetQueries(ctx, r.database).HubResolveReport(ctx, db.HubResolveReportParams{ID: arg.ID, Status: arg.Status, ReviewerID: arg.ReviewerID, ResolutionNote: arg.ResolutionNote}))
}
