package repository

import (
	"context"
	"go-server/internal/db"
	"go-server/internal/models"
	"go-server/internal/repositories/contracts"
)

type WebDeviceTokenRepository = contracts.WebDeviceTokenRepository

func (r *deviceTokenRepository) UpsertWeb(ctx context.Context, userID int64, token, deviceID string, version int64) (*models.DeviceToken, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, persistenceError(err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	q := db.New(tx)
	// Serialize installation transfers across accounts and token rotations.
	if _, err = tx.Exec(ctx, "SELECT pg_advisory_xact_lock(174294, 24)"); err != nil {
		return nil, persistenceError(err)
	}
	if _, err = q.LockWebPushUser(ctx, db.LockWebPushUserParams{UserID: userID, SessionVersion: version}); err != nil {
		return nil, persistenceError(err)
	}
	if err = q.DeletePreviousWebInstallation(ctx, db.DeletePreviousWebInstallationParams{Token: token, DeviceID: &deviceID}); err != nil {
		return nil, persistenceError(err)
	}
	row, err := q.InsertWebDeviceToken(ctx, db.InsertWebDeviceTokenParams{UserID: userID, Token: token, DeviceID: &deviceID, SessionVersion: &version})
	if err != nil {
		return nil, persistenceError(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, persistenceError(err)
	}
	return deviceTokenFromDB(row), nil
}

func (r *deviceTokenRepository) ListEligibleWeb(ctx context.Context, userID, societyID, flatID int64) ([]*models.DeviceToken, error) {
	rows, err := GetQueries(ctx, r.db).ListEligibleWebTokens(ctx, db.ListEligibleWebTokensParams{UserID: userID, SocietyID: societyID, FlatID: flatID})
	if err != nil {
		return nil, persistenceError(err)
	}
	result := make([]*models.DeviceToken, 0, len(rows))
	for _, row := range rows {
		result = append(result, deviceTokenFromDB(row))
	}
	return result, nil
}
