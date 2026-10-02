package notificationsvc

import (
	"context"
	"errors"
	"strings"

	"go-server/internal/config"
	"go-server/internal/models"
	repository "go-server/internal/repositories/contracts"
	"go-server/pkg/logger"

	"go.uber.org/zap"
)

type notificationService struct {
	deviceTokens  repository.DeviceTokenRepository
	notifications repository.NotificationRepository
	residents     repository.FlatResidentRepository
	members       repository.SocietyMemberRepository
	fcmClient     fcmClient
	apnsClient    *apnsClient
	enabled       bool
	webEnabled    bool
	operational   interface {
		EnsureSocietyOperational(context.Context, int64) error
	}
}

func NewNotificationService(ctx context.Context, repo repository.DeviceTokenRepository, cfg *config.Config, deps ...any) (NotificationService, error) {
	if err := validateFCMConfig(cfg); err != nil {
		return nil, err
	}

	var client fcmClient = noopFCMClient{}
	if cfg.FCMEnabled {
		fcmClient, err := newFirebaseFCMClient(ctx, cfg.FCMCredentialsPath)
		if err != nil {
			return nil, err
		}
		client = fcmClient
	}

	service := &notificationService{
		deviceTokens: repo,
		fcmClient:    client,
		enabled:      cfg.FCMEnabled,
		webEnabled:   cfg.FCMWebEnabled,
	}
	if cfg.APNSEnabled {
		apns, err := newAPNSClient(cfg)
		if err != nil {
			return nil, err
		}
		service.apnsClient = apns
	}
	for _, dep := range deps {
		switch value := dep.(type) {
		case repository.NotificationRepository:
			service.notifications = value
		case repository.FlatResidentRepository:
			service.residents = value
		case repository.SocietyMemberRepository:
			service.members = value
		case interface {
			EnsureSocietyOperational(context.Context, int64) error
		}:
			service.operational = value
		}
	}

	return service, nil
}

func (s *notificationService) RegisterDeviceToken(ctx context.Context, userID int64, req models.RegisterDeviceTokenRequest) (*models.DeviceToken, error) {
	token := strings.TrimSpace(req.Token)
	if token == "" || !req.Platform.IsValid() {
		return nil, ErrInvalidDeviceToken
	}
	if req.Platform == models.DevicePlatformAndroid && isExpoPushToken(token) {
		return nil, ErrInvalidDeviceToken
	}

	if req.Platform == models.DevicePlatformWeb {
		if !s.enabled || !s.webEnabled {
			return nil, ErrNotificationDisabled
		}
		repo, ok := s.deviceTokens.(repository.WebDeviceTokenRepository)
		if !ok || req.DeviceID == nil || strings.TrimSpace(*req.DeviceID) == "" || req.SessionVersion == nil || isExpoPushToken(token) {
			return nil, ErrInvalidDeviceToken
		}
		result, err := repo.UpsertWeb(ctx, userID, token, strings.TrimSpace(*req.DeviceID), *req.SessionVersion)
		if errors.Is(err, repository.ErrNotFound) {
			return nil, models.NewAppError("WEB_PUSH_INELIGIBLE", "Sign in with an active resident account to enable notifications", 403, nil)
		}
		return result, err
	}
	return s.deviceTokens.Upsert(ctx, userID, token, req.Platform, req.DeviceID)
}

func (s *notificationService) UnregisterDeviceToken(ctx context.Context, userID int64, token string) error {
	token = strings.TrimSpace(token)
	if token == "" {
		return ErrInvalidDeviceToken
	}

	if err := s.deviceTokens.Delete(ctx, userID, token); err != nil {
		return err
	}

	return nil
}

func (s *notificationService) SendToUser(ctx context.Context, userID int64, payload models.NotificationPayload) error {
	return s.SendToUsers(ctx, []int64{userID}, payload)
}

func (s *notificationService) SendToUsers(ctx context.Context, userIDs []int64, payload models.NotificationPayload) error {
	return s.sendPushToUsers(ctx, userIDs, payload)
}

func (s *notificationService) sendPushToUsers(ctx context.Context, userIDs []int64, payload models.NotificationPayload) (resultErr error) {
	if strings.TrimSpace(payload.Title) == "" && strings.TrimSpace(payload.Body) == "" {
		return ErrInvalidDeviceToken
	}

	// Native delivery proceeds first; browser failures cannot skip Android delivery.
	defer func() { resultErr = errors.Join(resultErr, s.sendWebPush(ctx, userIDs, payload)) }()
	fcmTokens := make([]string, 0)
	seen := make(map[string]struct{})

	for _, userID := range userIDs {
		rows, err := s.deviceTokens.ListByUserID(ctx, userID)
		if err != nil {
			return ErrNotificationSend.WithCause(err)
		}

		for _, row := range rows {
			token := strings.TrimSpace(row.Token)
			if token == "" {
				continue
			}
			if _, exists := seen[token]; exists {
				continue
			}
			seen[token] = struct{}{}
			if isExpoPushToken(token) {
				s.deleteInvalidToken(ctx, "expo", token, "legacy Expo token is no longer supported")
				continue
			}
			if row.Platform == models.DevicePlatformIOS {
				logger.Debug("skipping iOS push token until APNs delivery is implemented",
					zap.String("type", notificationPayloadType(payload)),
					zap.Int64("user_id", userID),
					zap.String("token", redactPushToken(token)),
				)
				continue
			}
			if row.Platform != models.DevicePlatformAndroid {
				logger.Debug("skipping push token with unsupported platform",
					zap.String("type", notificationPayloadType(payload)),
					zap.Int64("user_id", userID),
					zap.String("platform", string(row.Platform)),
				)
				continue
			}
			fcmTokens = append(fcmTokens, token)
		}
	}

	notificationType := notificationPayloadType(payload)
	logger.Debug("push notification tokens loaded",
		zap.String("type", notificationType),
		zap.Int("recipients", len(uniqueUserIDs(userIDs))),
		zap.Int64s("recipient_user_ids", uniqueUserIDs(userIDs)),
		zap.Int("tokens", len(fcmTokens)),
		zap.Int("fcm_tokens", len(fcmTokens)),
	)

	if len(fcmTokens) == 0 {
		return nil
	}

	if !s.enabled {
		if len(fcmTokens) > 0 {
			logger.Debug("direct FCM push skipped because FCM is disabled",
				zap.String("type", notificationType),
				zap.Int("fcm_tokens", len(fcmTokens)),
			)
		}
		return nil
	}

	err := sendMulticast(ctx, s.fcmClient, fcmTokens, payload, func(token string) {
		s.deleteInvalidToken(ctx, "fcm", token, "invalid registration token")
	})
	if err != nil {
		return ErrNotificationSend.WithCause(err)
	}

	return nil
}

func (s *notificationService) deleteInvalidToken(ctx context.Context, provider string, token string, reason string) {
	if strings.TrimSpace(token) == "" {
		return
	}
	logger.Warn("deleting invalid push token",
		zap.String("provider", provider),
		zap.String("token", redactPushToken(token)),
		zap.String("reason", reason),
	)
	if deleteErr := s.deviceTokens.DeleteByToken(ctx, token); deleteErr != nil {
		logger.Warn("failed to delete invalid device token", zap.Error(deleteErr))
	}
}

func isExpoPushToken(token string) bool {
	token = strings.TrimSpace(token)
	return strings.HasPrefix(token, "ExpoPushToken[") || strings.HasPrefix(token, "ExponentPushToken[")
}

func redactPushToken(token string) string {
	token = strings.TrimSpace(token)
	if len(token) <= 12 {
		return "***"
	}
	return token[:8] + "..." + token[len(token)-4:]
}

func notificationPayloadType(payload models.NotificationPayload) string {
	if payload.Data == nil {
		return ""
	}
	return payload.Data["type"]
}

func (s *notificationService) Close() error {
	if s.apnsClient != nil {
		s.apnsClient.http.CloseIdleConnections()
	}
	if s.fcmClient == nil {
		return nil
	}
	return s.fcmClient.Close()
}
