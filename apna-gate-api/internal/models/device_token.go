package models

import "time"

type DevicePlatform string

const (
	DevicePlatformIOS     DevicePlatform = "ios"
	DevicePlatformAndroid DevicePlatform = "android"
	DevicePlatformWeb     DevicePlatform = "web"
)

func (p DevicePlatform) IsValid() bool {
	switch p {
	case DevicePlatformIOS, DevicePlatformAndroid, DevicePlatformWeb:
		return true
	default:
		return false
	}
}

type DeviceToken struct {
	ID         int64          `json:"id"`
	UserID     int64          `json:"user_id"`
	Token      string         `json:"token"`
	Platform   DevicePlatform `json:"platform"`
	DeviceID   *string        `json:"device_id,omitempty"`
	LastSeenAt time.Time      `json:"last_seen_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type RegisterDeviceTokenRequest struct {
	SessionVersion *int64         `json:"-" swaggerignore:"true"`
	Token          string         `json:"token" binding:"required"`
	Platform       DevicePlatform `json:"platform" binding:"required"`
	DeviceID       *string        `json:"device_id,omitempty"`
}

type UnregisterDeviceTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

type NotificationPayload struct {
	Title      string            `json:"title"`
	Body       string            `json:"body"`
	Data       map[string]string `json:"data,omitempty"`
	ImageURL   string            `json:"image_url,omitempty"`
	CategoryID string            `json:"-"`
}
