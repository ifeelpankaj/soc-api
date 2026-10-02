package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type UpdateUserRequest struct {
	ID int64 `json:"-"`

	FirstName   *string    `json:"first_name,omitempty" validate:"omitempty,min=2,max=100,alphanumeric_space"`
	LastName    *string    `json:"last_name,omitempty" validate:"omitempty,max=100,alphanumeric_space"`
	PhoneNumber *string    `json:"phone_number,omitempty" validate:"omitempty,phone_intl,max=20"`
	AvatarURL   *string    `json:"avatar_url,omitempty" validate:"omitempty,url,max=500"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Gender      *string    `json:"gender,omitempty" validate:"omitempty,oneof=male female other prefer_not_to_say"`
	Timezone    *string    `json:"timezone,omitempty" validate:"omitempty,max=100"`
	Language    *string    `json:"language,omitempty" validate:"omitempty,max=20"`
}

type updateUserRequestJSON struct {
	FirstName   *string `json:"first_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	DateOfBirth *string `json:"date_of_birth,omitempty"`
	Gender      *string `json:"gender,omitempty"`
	Timezone    *string `json:"timezone,omitempty"`
	Language    *string `json:"language,omitempty"`
}

func (r *UpdateUserRequest) UnmarshalJSON(data []byte) error {
	var raw updateUserRequestJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	r.FirstName = raw.FirstName
	r.LastName = raw.LastName
	r.PhoneNumber = raw.PhoneNumber
	r.AvatarURL = raw.AvatarURL
	r.Gender = raw.Gender
	r.Timezone = raw.Timezone
	r.Language = raw.Language

	if raw.DateOfBirth == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*raw.DateOfBirth)
	if trimmed == "" {
		return nil
	}

	parsed, err := ParseFlexibleDate(trimmed)
	if err != nil {
		return err
	}
	r.DateOfBirth = &parsed
	return nil
}

func ParseFlexibleDate(raw string) (time.Time, error) {
	layouts := []string{
		"2006-01-02",
		"2006/1/2",
		"2006/01/02",
		"2/1/2006",
		"02/01/2006",
		time.RFC3339,
	}

	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, raw); err == nil {
			return time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC), nil
		}
	}

	return time.Time{}, fmt.Errorf("date_of_birth must use YYYY-MM-DD format")
}

func (r *UpdateUserRequest) Sanitize() {
	if r.FirstName != nil {
		v := strings.TrimSpace(*r.FirstName)
		r.FirstName = &v
	}

	if r.LastName != nil {
		v := strings.TrimSpace(*r.LastName)
		r.LastName = &v
	}

	if r.PhoneNumber != nil {
		v := strings.TrimSpace(*r.PhoneNumber)
		r.PhoneNumber = &v
	}

	if r.AvatarURL != nil {
		v := strings.TrimSpace(*r.AvatarURL)
		r.AvatarURL = &v
	}

	if r.Gender != nil {
		v := strings.ToLower(strings.TrimSpace(*r.Gender))
		r.Gender = &v
	}

	if r.Timezone != nil {
		v := strings.TrimSpace(*r.Timezone)
		r.Timezone = &v
	}

	if r.Language != nil {
		v := strings.ToLower(strings.TrimSpace(*r.Language))
		r.Language = &v
	}
}

func (r *UpdateUserRequest) Validate() error {
	if r.AvatarURL != nil && strings.TrimSpace(*r.AvatarURL) != "" {
		return ErrAvatarUploadRequired
	}
	r.AvatarURL = nil
	if r.FirstName == nil &&
		r.LastName == nil &&
		r.PhoneNumber == nil &&
		r.AvatarURL == nil &&
		r.DateOfBirth == nil &&
		r.Gender == nil &&
		r.Timezone == nil &&
		r.Language == nil {
		return errors.New("at least one field must be provided for update")
	}

	if r.FirstName != nil && strings.TrimSpace(*r.FirstName) == "" {
		return errors.New("first_name cannot be empty")
	}

	return nil
}
