package notificationsvc

import (
	"context"
	"errors"
	"go-server/internal/models"
	"testing"
)

type webTokenFake struct {
	deviceTokenRepoFake
	rows    []*models.DeviceToken
	version int64
	scope   [3]int64
	fail    bool
}

func (r *webTokenFake) UpsertWeb(_ context.Context, userID int64, token, device string, version int64) (*models.DeviceToken, error) {
	r.version = version
	return &models.DeviceToken{UserID: userID, Token: token}, nil
}
func (r *webTokenFake) ListEligibleWeb(_ context.Context, userID, societyID, flatID int64) ([]*models.DeviceToken, error) {
	r.scope = [3]int64{userID, societyID, flatID}
	if r.fail {
		return nil, errors.New("web lookup failed")
	}
	return r.rows, nil
}
func TestWebPushUsesDataOnlyAndResidentScope(t *testing.T) {
	repo := &webTokenFake{rows: []*models.DeviceToken{{Token: "browser"}}, deviceTokenRepoFake: deviceTokenRepoFake{tokensByUser: map[int64][]string{7: {"android"}}}}
	client := &fcmClientFake{}
	s := &notificationService{deviceTokens: repo, fcmClient: client, enabled: true, webEnabled: true}
	err := s.SendToUser(context.Background(), 7, models.NotificationPayload{Title: "Visitor", Body: "At gate", Data: map[string]string{"notification_id": "id", "society_id": "8", "flat_id": "9", "auth_token": "must not be sent", "url": "https://evil.example"}})
	if err != nil {
		t.Fatal(err)
	}
	if repo.scope != [3]int64{7, 8, 9} || len(client.messages) != 2 {
		t.Fatalf("scope/sends incorrect: %v %d", repo.scope, len(client.messages))
	}
	native, web := client.messages[0], client.messages[1]
	if web.Notification != nil || web.Webpush.Notification != nil || web.Data["user_id"] != "7" || web.Data["notification_id"] != "id" || web.Data["auth_token"] != "" || web.Data["url"] != "" {
		t.Fatal("unsafe/duplicate web payload")
	}
	if native.Notification == nil || native.Notification.Title != "Visitor" {
		t.Fatal("native display changed")
	}
}
func TestWebFailureDoesNotPreventAndroidDelivery(t *testing.T) {
	repo := &webTokenFake{fail: true, deviceTokenRepoFake: deviceTokenRepoFake{tokensByUser: map[int64][]string{7: {"android"}}}}
	client := &fcmClientFake{}
	s := &notificationService{deviceTokens: repo, fcmClient: client, enabled: true, webEnabled: true}
	if err := s.SendToUser(context.Background(), 7, models.NotificationPayload{Title: "Update"}); err == nil {
		t.Fatal("expected web error")
	}
	if len(client.messages) != 1 || client.messages[0].Notification == nil {
		t.Fatal("Android delivery was skipped")
	}
}
func TestWebRegistrationRequiresSessionInstallationAndRollout(t *testing.T) {
	repo := &webTokenFake{}
	s := &notificationService{deviceTokens: repo, enabled: true, webEnabled: true}
	version := int64(4)
	device := "browser-install"
	request := models.RegisterDeviceTokenRequest{Token: "fcm", Platform: models.DevicePlatformWeb, DeviceID: &device, SessionVersion: &version}
	if _, err := s.RegisterDeviceToken(context.Background(), 7, request); err != nil || repo.version != 4 {
		t.Fatalf("version not bound: %v", err)
	}
	request.SessionVersion = nil
	if _, err := s.RegisterDeviceToken(context.Background(), 7, request); err == nil {
		t.Fatal("missing version accepted")
	}
	request.SessionVersion = &version
	request.DeviceID = nil
	if _, err := s.RegisterDeviceToken(context.Background(), 7, request); err == nil {
		t.Fatal("missing installation accepted")
	}
	s.webEnabled = false
	if _, err := s.RegisterDeviceToken(context.Background(), 7, request); !errors.Is(err, ErrNotificationDisabled) {
		t.Fatal("rollout flag ignored")
	}
}
