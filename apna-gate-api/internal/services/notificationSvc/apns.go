package notificationsvc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-server/internal/config"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type apnsClient struct {
	key                            *jwt.SigningMethodECDSA
	privateKey                     any
	keyID, teamID, topic, endpoint string
	http                           *http.Client
	mu                             sync.Mutex
	authToken                      string
	authCreated                    time.Time
}

func newAPNSClient(cfg *config.Config) (*apnsClient, error) {
	if cfg.APNSKeyPath == "" || cfg.APNSKeyID == "" || cfg.APNSTeamID == "" || cfg.APNSTopic == "" {
		return nil, errors.New("APNS_KEY_PATH, APNS_KEY_ID, APNS_TEAM_ID and APNS_TOPIC are required")
	}
	b, err := os.ReadFile(cfg.APNSKeyPath)
	if err != nil {
		return nil, err
	}
	key, err := jwt.ParseECPrivateKeyFromPEM(b)
	if err != nil {
		return nil, err
	}
	endpoint := "https://api.sandbox.push.apple.com/3/device/"
	if cfg.APNSProduction {
		endpoint = "https://api.push.apple.com/3/device/"
	}
	return &apnsClient{key: jwt.SigningMethodES256, privateKey: key, keyID: cfg.APNSKeyID, teamID: cfg.APNSTeamID, topic: cfg.APNSTopic, endpoint: endpoint, http: &http.Client{Timeout: 20 * time.Second}}, nil
}

func (c *apnsClient) Send(ctx context.Context, token, title, body string, data map[string]string) (string, string, bool, error) {
	if c == nil {
		return "", "apns_disabled", false, errors.New("APNs is not configured")
	}
	auth, err := c.authorizationToken()
	if err != nil {
		return "", "signing_failed", false, err
	}
	aps := map[string]any{"alert": map[string]string{"title": title, "body": body}, "sound": "default"}
	if category := data["category_id"]; category != "" {
		aps["category"] = category
	}
	payload := map[string]any{"aps": aps}
	for k, v := range data {
		if k != "aps" {
			payload[k] = v
		}
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return "", "invalid_payload", true, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint+token, bytes.NewReader(b))
	if err != nil {
		return "", "request_failed", false, err
	}
	req.Header.Set("authorization", "bearer "+auth)
	req.Header.Set("apns-topic", c.topic)
	req.Header.Set("apns-push-type", "alert")
	req.Header.Set("apns-priority", "10")
	req.Header.Set("content-type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return "", "network_error", false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return resp.Header.Get("apns-id"), "", false, nil
	}
	var result struct {
		Reason string `json:"reason"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&result)
	code := strings.TrimSpace(result.Reason)
	if code == "" {
		code = fmt.Sprintf("http_%d", resp.StatusCode)
	}
	if code == "ExpiredProviderToken" || code == "InvalidProviderToken" {
		c.mu.Lock()
		c.authToken = ""
		c.mu.Unlock()
	}
	permanent := code == "BadDeviceToken" || code == "Unregistered" || code == "DeviceTokenNotForTopic" || code == "BadTopic" || code == "PayloadTooLarge"
	return "", code, permanent, fmt.Errorf("APNs rejected delivery: %s", code)
}

// Apple requires a fresh provider token within an hour, but rejects frequent
// token changes on a persistent connection. Reuse one signed JWT for 45 minutes.
func (c *apnsClient) authorizationToken() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.authToken != "" && time.Since(c.authCreated) < 45*time.Minute {
		return c.authToken, nil
	}
	now := time.Now()
	j := jwt.NewWithClaims(c.key, jwt.MapClaims{"iss": c.teamID, "iat": now.Unix()})
	j.Header["kid"] = c.keyID
	signed, err := j.SignedString(c.privateKey)
	if err != nil {
		return "", err
	}
	c.authToken = signed
	c.authCreated = now
	return signed, nil
}
