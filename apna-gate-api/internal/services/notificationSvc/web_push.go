package notificationsvc

import (
	"context"
	"errors"
	"firebase.google.com/go/v4/messaging"
	"go-server/internal/models"
	repository "go-server/internal/repositories/contracts"
	"strconv"
)

// Data-only: the bundled service worker owns background display and click routing.
// Never include auth tokens or caller-supplied destinations.
func buildWebMulticast(tokens []string, userID int64, payload models.NotificationPayload) *messaging.MulticastMessage {
	data := map[string]string{"title": payload.Title, "body": payload.Body, "user_id": strconv.FormatInt(userID, 10)}
	for _, key := range []string{"notification_id", "type", "society_id", "flat_id", "entry_id", "invite_id", "invite_type"} {
		if value := payload.Data[key]; value != "" {
			data[key] = value
		}
	}
	return &messaging.MulticastMessage{Tokens: tokens, Data: data, Webpush: &messaging.WebpushConfig{Headers: map[string]string{"TTL": "300", "Urgency": "high"}}}
}

func (s *notificationService) sendWebPush(ctx context.Context, userIDs []int64, payload models.NotificationPayload) error {
	if !s.enabled || !s.webEnabled {
		return nil
	}
	repo, ok := s.deviceTokens.(repository.WebDeviceTokenRepository)
	if !ok {
		return ErrNotificationDisabled
	}
	societyID, _ := strconv.ParseInt(payload.Data["society_id"], 10, 64)
	flatID, _ := strconv.ParseInt(payload.Data["flat_id"], 10, 64)
	var failures []error
	for _, userID := range uniqueUserIDs(userIDs) {
		rows, err := repo.ListEligibleWeb(ctx, userID, societyID, flatID)
		if err != nil {
			failures = append(failures, err)
			continue
		}
		tokens := make([]string, 0, len(rows))
		seen := map[string]bool{}
		for _, row := range rows {
			if row.Token != "" && !seen[row.Token] {
				tokens = append(tokens, row.Token)
				seen[row.Token] = true
			}
		}
		for start := 0; start < len(tokens); start += 500 {
			end := min(start+500, len(tokens))
			batch := tokens[start:end]
			response, err := s.fcmClient.SendEachForMulticast(ctx, buildWebMulticast(batch, userID, payload))
			if err != nil {
				failures = append(failures, err)
				continue
			}
			for i, result := range response.Responses {
				if result.Success || result.Error == nil {
					continue
				}
				if messaging.IsUnregistered(result.Error) {
					s.deleteInvalidToken(ctx, "web-fcm", batch[i], "unregistered browser subscription")
				} else {
					failures = append(failures, result.Error)
				}
			}
		}
	}
	if err := errors.Join(failures...); err != nil {
		return ErrNotificationSend.WithCause(err)
	}
	return nil
}
