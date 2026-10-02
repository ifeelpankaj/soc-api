package notificationsvc

import (
	"context"
	"fmt"
	"os"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"go-server/internal/config"
	"go-server/internal/models"
	"go-server/pkg/logger"

	"go.uber.org/zap"
	"google.golang.org/api/option"
)

type fcmClient interface {
	SendEachForMulticast(ctx context.Context, message *messaging.MulticastMessage) (*messaging.BatchResponse, error)
	Close() error
}

type firebaseFCMClient struct {
	client *messaging.Client
}

func newFirebaseFCMClient(ctx context.Context, credentialsPath string) (*firebaseFCMClient, error) {
	// SECURITY-REVIEW: loads Firebase Admin credentials from a filesystem path configured via env var.
	credentials, err := os.ReadFile(credentialsPath) //nolint:gosec // Path comes from trusted service configuration for Firebase credentials.
	if err != nil {
		return nil, fmt.Errorf("read firebase credentials: %w", err)
	}

	app, err := firebase.NewApp(ctx, nil, option.WithAuthCredentialsJSON(option.ServiceAccount, credentials))
	if err != nil {
		return nil, fmt.Errorf("initialize firebase app: %w", err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("initialize firebase messaging client: %w", err)
	}

	return &firebaseFCMClient{client: client}, nil
}

func (c *firebaseFCMClient) SendEachForMulticast(ctx context.Context, message *messaging.MulticastMessage) (*messaging.BatchResponse, error) {
	return c.client.SendEachForMulticast(ctx, message)
}

func (c *firebaseFCMClient) Close() error {
	return nil
}

type noopFCMClient struct{}

func (noopFCMClient) SendEachForMulticast(context.Context, *messaging.MulticastMessage) (*messaging.BatchResponse, error) {
	return &messaging.BatchResponse{}, nil
}

func (noopFCMClient) Close() error {
	return nil
}

func buildMulticastMessage(tokens []string, payload models.NotificationPayload) *messaging.MulticastMessage {
	data := make(map[string]string, len(payload.Data)+1)
	for key, value := range payload.Data {
		data[key] = value
	}
	categoryID := strings.TrimSpace(payload.CategoryID)
	if categoryID != "" {
		// Expo Notifications SDK 57 reads categoryId from direct FCM data on Android.
		data["categoryId"] = categoryID
	}
	message := &messaging.MulticastMessage{
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title: payload.Title,
			Body:  payload.Body,
		},
		Data: data,
	}
	if categoryID != "" {
		message.APNS = &messaging.APNSConfig{Payload: &messaging.APNSPayload{
			Aps: &messaging.Aps{Category: categoryID},
		}}
	}

	if strings.TrimSpace(payload.ImageURL) != "" {
		message.Notification.ImageURL = payload.ImageURL
	}

	return message
}

func sendMulticast(ctx context.Context, client fcmClient, tokens []string, payload models.NotificationPayload, onInvalidToken func(token string)) error {
	if len(tokens) == 0 {
		return nil
	}

	response, err := client.SendEachForMulticast(ctx, buildMulticastMessage(tokens, payload))
	if err != nil {
		return err
	}

	for index, item := range response.Responses {
		if item.Success || item.Error == nil {
			continue
		}

		if messaging.IsUnregistered(item.Error) || messaging.IsInvalidArgument(item.Error) {
			onInvalidToken(tokens[index])
			continue
		}

		logger.Warn("failed to deliver push notification",
			zap.Int("token_index", index),
			zap.Error(item.Error),
		)
	}

	return nil
}

func validateFCMConfig(cfg *config.Config) error {
	if !cfg.FCMEnabled {
		return nil
	}
	if strings.TrimSpace(cfg.FCMCredentialsPath) == "" {
		return fmt.Errorf("FCM_CREDENTIALS_PATH is required when FCM is enabled")
	}
	return nil
}
