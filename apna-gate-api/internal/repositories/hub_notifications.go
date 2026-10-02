package repository

import (
	"context"
	"go-server/internal/db"
	"go-server/internal/repositories/contracts"
)

func (r *HubRepository) HubEnqueue(ctx context.Context, arg contracts.HubEnqueueInput) error {
	return persistenceError(GetQueries(ctx, r.database).HubEnqueue(ctx, db.HubEnqueueParams{UserID: arg.UserID, SocietyID: arg.SocietyID, EventKey: arg.EventKey, Payload: arg.Payload, PushEnabled: arg.PushEnabled}))
}

// HubQueueAnnouncementFanout is called inside the post transaction.
func (r *HubRepository) HubQueueAnnouncementFanout(ctx context.Context, societyID, postID, actorID int64, important bool) error {
	_, err := GetExecutor(ctx, r.database).Exec(ctx, `INSERT INTO hub_notification_fanout_jobs(society_id,post_id,actor_user_id,important)
 VALUES($1,$2,$3,$4) ON CONFLICT(society_id,post_id) DO NOTHING`, societyID, postID, actorID, important)
	return persistenceError(err)
}
