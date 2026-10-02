package notificationsvc

import (
	"context"
	"errors"
	"fmt"
	"go-server/internal/models"
	"go-server/internal/repositories/contracts"
	"go-server/pkg/logger"
	"strconv"
	"sync"
	"time"

	"firebase.google.com/go/v4/messaging"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

var pipelineHeartbeat = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "apna_gate_notification_worker_heartbeat_seconds", Help: "Last successful loop iteration as Unix time."}, []string{"loop"})
var pipelineBacklog = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "apna_gate_notification_pipeline_backlog", Help: "Pipeline backlog by state."}, []string{"state"})
var providerFailures = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "apna_gate_push_provider_failures_total", Help: "Push provider failures."}, []string{"provider", "code"})
var pushLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "apna_gate_push_delivery_latency_seconds", Help: "Provider call duration.", Buckets: prometheus.DefBuckets}, []string{"provider"})
var loopLastSuccess sync.Map

func PipelineReady() bool {
	for _, name := range []string{"materializer", "hub_fanout", "push_delivery"} {
		value, ok := loopLastSuccess.Load(name)
		if !ok || time.Since(value.(time.Time)) > 60*time.Second {
			return false
		}
	}
	return true
}

func init() {
	prometheus.MustRegister(pipelineHeartbeat, pipelineBacklog, providerFailures, pushLatency)
}

// RunPipeline supervises three independent loops. A transient failure cannot
// terminate a loop; each successful iteration updates its own heartbeat.
func (s *notificationService) RunPipeline(ctx context.Context) error {
	store, ok := s.notifications.(contracts.PipelineRepository)
	if !ok {
		return errors.New("notification pipeline repository is unavailable")
	}
	var wg sync.WaitGroup
	for name, run := range map[string]func(context.Context) error{
		"materializer":  func(c context.Context) error { _, e := s.deliverOutboxBatch(c); return e },
		"hub_fanout":    func(c context.Context) error { return s.runFanoutBatch(c, store) },
		"push_delivery": func(c context.Context) error { return s.runPushBatch(c, store) },
	} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pipelineLoop(ctx, name, run)
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		pipelineLoop(ctx, "metrics", func(c context.Context) error {
			stats, err := store.PipelineStats(c)
			if err != nil {
				return err
			}
			for k, v := range stats {
				pipelineBacklog.WithLabelValues(k).Set(v)
			}
			return nil
		})
	}()
	<-ctx.Done()
	wg.Wait()
	return ctx.Err()
}

func pipelineLoop(ctx context.Context, name string, run func(context.Context) error) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		attempt, cancel := context.WithTimeout(ctx, 50*time.Second)
		err := run(attempt)
		cancel()
		if err == nil {
			now := time.Now()
			loopLastSuccess.Store(name, now)
			pipelineHeartbeat.WithLabelValues(name).Set(float64(now.Unix()))
		} else if ctx.Err() == nil {
			logger.Warn("notification worker iteration failed", zap.String("loop", name), zap.Error(err))
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (s *notificationService) runFanoutBatch(ctx context.Context, store contracts.PipelineRepository) error {
	for i := 0; i < 10; i++ {
		job, err := store.ClaimHubFanout(ctx, uuid.New())
		if err != nil {
			return err
		}
		if job == nil {
			return nil
		}
		if _, err = store.FanoutHubBatch(ctx, *job, 200); err != nil {
			return err
		}
	}
	return nil
}

func (s *notificationService) runPushBatch(ctx context.Context, store contracts.PipelineRepository) error {
	for i := 0; i < 100; i++ {
		d, err := store.ClaimPushDelivery(ctx, uuid.New())
		if err != nil {
			return err
		}
		if d == nil {
			break
		}
		messageID, code, permanent, err := s.sendPushDelivery(ctx, d)
		if err == nil {
			err = store.FinishPushDelivery(ctx, d.ID, d.LeaseToken, messageID)
			if err != nil {
				return err
			}
		} else {
			if providerFailure(code) {
				providerFailures.WithLabelValues(d.Provider, code).Inc()
			}
			cleanup, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
			if invalidPushToken(code) && d.Token != "" {
				if deleteErr := s.deviceTokens.DeleteByToken(cleanup, d.Token); deleteErr != nil {
					cancel()
					return fmt.Errorf("retire invalid push token: %w", deleteErr)
				}
			}
			finishErr := store.FailPushDelivery(cleanup, d.ID, d.LeaseToken, code, err.Error(), permanent)
			cancel()
			if finishErr != nil {
				return errors.Join(err, finishErr)
			}
			logger.Warn("push delivery failed", zap.Int64("delivery_id", d.ID), zap.String("provider", d.Provider), zap.String("code", code), zap.Error(err))
		}
	}
	return nil
}

func providerFailure(code string) bool {
	switch code {
	case "token_missing", "platform_mismatch", "hub_post_missing", "hub_check_failed", "hub_post_removed", "access_check_failed", "access_revoked", "web_eligibility_failed", "web_session_revoked", "fcm_disabled", "web_push_disabled", "apns_disabled":
		return false
	}
	return true
}

func invalidPushToken(code string) bool {
	switch code {
	case "fcm_invalid_token", "BadDeviceToken", "Unregistered", "DeviceTokenNotForTopic":
		return true
	}
	return false
}

func (s *notificationService) sendPushDelivery(ctx context.Context, d *contracts.PushDelivery) (string, string, bool, error) {
	if d.Token == "" {
		return "", "token_missing", true, errors.New("device token no longer exists")
	}
	if (d.Platform == models.DevicePlatformIOS) != (d.Provider == "apns") {
		return "", "platform_mismatch", true, errors.New("token provider no longer matches its platform")
	}
	if d.SocietyID != nil {
		if d.Audience == "hub_member" {
			postID, parseErr := strconv.ParseInt(fmt.Sprint(d.Data["post_id"]), 10, 64)
			if parseErr != nil || postID <= 0 {
				return "", "hub_post_missing", true, errors.New("hub post missing")
			}
			if pipeline, ok := s.notifications.(contracts.PipelineRepository); ok {
				available, err := pipeline.HubContentAvailable(ctx, *d.SocietyID, postID)
				if err != nil {
					return "", "hub_check_failed", false, err
				}
				if !available {
					return "", "hub_post_removed", true, errors.New("hub post removed")
				}
			}
		}
		if access, ok := s.notifications.(contracts.OutboxRepository); ok {
			audience := d.Audience
			allowed, err := access.NotificationOutboxAccess(ctx, contracts.NotificationOutboxAccessInput{SocietyID: *d.SocietyID, UserID: d.UserID, Audience: audience, FlatID: d.FlatID})
			if err != nil {
				return "", "access_check_failed", false, err
			}
			if !allowed {
				return "", "access_revoked", true, errors.New("recipient access revoked")
			}
		}
	}
	if d.Platform == models.DevicePlatformWeb {
		if web, ok := s.deviceTokens.(contracts.WebDeviceTokenRepository); ok {
			society, flat := int64(0), int64(0)
			if d.SocietyID != nil {
				society = *d.SocietyID
			}
			if d.FlatID != nil {
				flat = *d.FlatID
			}
			eligible, err := web.ListEligibleWeb(ctx, d.UserID, society, flat)
			if err != nil {
				return "", "web_eligibility_failed", false, err
			}
			found := false
			for _, t := range eligible {
				if t.ID == *d.TokenID {
					found = true
					break
				}
			}
			if !found {
				return "", "web_session_revoked", true, errors.New("web session revoked")
			}
		}
	}
	data := stringMap(d.Data)
	data["notification_id"] = d.NotificationID
	payload := models.NotificationPayload{Title: d.Title, Body: d.Body, Data: data, CategoryID: data["category_id"]}
	start := time.Now()
	defer pushLatency.WithLabelValues(d.Provider).Observe(time.Since(start).Seconds())
	callCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	if d.Provider == "apns" {
		return s.apnsClient.Send(callCtx, d.Token, d.Title, d.Body, data)
	}
	if !s.enabled {
		return "", "fcm_disabled", false, errors.New("FCM disabled")
	}
	var msg *messaging.MulticastMessage
	if d.Platform == models.DevicePlatformWeb {
		if !s.webEnabled {
			return "", "web_push_disabled", false, errors.New("web push disabled")
		}
		msg = buildWebMulticast([]string{d.Token}, d.UserID, payload)
	} else {
		msg = buildMulticastMessage([]string{d.Token}, payload)
	}
	response, err := s.fcmClient.SendEachForMulticast(callCtx, msg)
	if err != nil {
		return "", "provider_error", false, err
	}
	if response == nil || len(response.Responses) != 1 {
		if response == nil {
			return "", "invalid_provider_response", false, errors.New("FCM returned no response")
		}
		return "", "invalid_provider_response", false, fmt.Errorf("FCM returned %d results", len(response.Responses))
	}
	result := response.Responses[0]
	if result.Success {
		return result.MessageID, "", false, nil
	}
	if result.Error == nil {
		return "", "unknown_provider_error", false, errors.New("FCM rejected delivery")
	}
	permanent := messaging.IsUnregistered(result.Error) || messaging.IsInvalidArgument(result.Error)
	code := "fcm_transient"
	if messaging.IsUnregistered(result.Error) {
		code = "fcm_invalid_token"
	} else if messaging.IsInvalidArgument(result.Error) {
		code = "fcm_invalid_argument"
	}
	return "", code, permanent, result.Error
}
