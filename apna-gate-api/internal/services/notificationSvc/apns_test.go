package notificationsvc

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestAPNSSendUsesNativeTokenAndClassifiesFailure(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	status, reason := http.StatusOK, ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/3/device/native-token" || r.Header.Get("apns-topic") != "org.apnagate.app" || !strings.HasPrefix(r.Header.Get("authorization"), "bearer ") {
			t.Errorf("invalid APNs request: %s", r.URL.Path)
		}
		var body map[string]any
		if e := json.NewDecoder(r.Body).Decode(&body); e != nil {
			t.Error(e)
		}
		aps, _ := body["aps"].(map[string]any)
		if aps["category"] != "visitor_decision" || body["notification_id"] != "inbox-id" {
			t.Errorf("invalid APNs payload: %#v", body)
		}
		w.Header().Set("apns-id", "provider-id")
		w.WriteHeader(status)
		if reason != "" {
			_, _ = w.Write([]byte(`{"reason":"` + reason + `"}`))
		}
	}))
	defer server.Close()
	c := &apnsClient{key: jwt.SigningMethodES256, privateKey: key, keyID: "key", teamID: "team", topic: "org.apnagate.app", endpoint: server.URL + "/3/device/", http: &http.Client{Timeout: time.Second}}
	firstToken, err := c.authorizationToken()
	if err != nil {
		t.Fatal(err)
	}
	secondToken, err := c.authorizationToken()
	if err != nil || firstToken != secondToken {
		t.Fatal("APNs provider token was not reused")
	}
	data := map[string]string{"notification_id": "inbox-id", "category_id": "visitor_decision"}
	id, code, permanent, err := c.Send(context.Background(), "native-token", "Visitor", "Waiting", data)
	if err != nil || id != "provider-id" || code != "" || permanent {
		t.Fatalf("success result: %q %q %v %v", id, code, permanent, err)
	}
	status, reason = http.StatusGone, "Unregistered"
	_, code, permanent, err = c.Send(context.Background(), "native-token", "Visitor", "Waiting", data)
	if err == nil || code != "Unregistered" || !permanent {
		t.Fatalf("invalid token result: %q %v %v", code, permanent, err)
	}
	status, reason = http.StatusForbidden, "ExpiredProviderToken"
	_, code, permanent, err = c.Send(context.Background(), "native-token", "Visitor", "Waiting", data)
	if err == nil || code != "ExpiredProviderToken" || permanent {
		t.Fatalf("configuration failure result: %q %v %v", code, permanent, err)
	}
	if c.authToken != "" {
		t.Fatal("expired APNs provider token was kept in cache")
	}
}
