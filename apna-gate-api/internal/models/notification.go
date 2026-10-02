package models

import "time"

type Notification struct {
	ID        string         `json:"id"`
	UserID    int64          `json:"user_id"`
	SocietyID *int64         `json:"society_id,omitempty"`
	FlatID    *int64         `json:"flat_id,omitempty"`
	Type      string         `json:"type"`
	Domain    string         `json:"domain"`
	Title     string         `json:"title"`
	Body      string         `json:"body"`
	Data      map[string]any `json:"data"`
	EventKey  *string        `json:"event_key,omitempty"`
	ReadAt    *time.Time     `json:"read_at,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

type NotificationCreate struct {
	ID        string
	UserID    int64
	SocietyID *int64
	FlatID    *int64
	Type      string
	Title     string
	Body      string
	Data      map[string]any
	EventKey  *string
}

type NotificationListFilter struct {
	UserID          int64
	CursorCreatedAt *time.Time
	CursorID        *string
	Limit           int32
}

type NotificationListResult struct {
	Items      []*Notification `json:"items"`
	NextCursor *string         `json:"next_cursor,omitempty"`
}

type NotificationUnreadCountResponse struct {
	UnreadCount int64 `json:"unread_count"`
}

type NotificationReadResponse struct {
	ID     string     `json:"id"`
	ReadAt *time.Time `json:"read_at,omitempty"`
}

type NotificationDetailAPIResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    *Notification `json:"data"`
}
