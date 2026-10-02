package models

import "time"

// MaintenanceIssuer is frozen by PostgreSQL at issuance or migration rollout.
type MaintenanceIssuer struct {
	Name          string    `json:"name"`
	SocietyCode   string    `json:"society_code"`
	AddressLine1  string    `json:"address_line1"`
	AddressLine2  string    `json:"address_line2"`
	Landmark      string    `json:"landmark"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	Pincode       string    `json:"pincode"`
	Country       string    `json:"country"`
	Email         string    `json:"email"`
	PhoneNumber   string    `json:"phone_number"`
	CapturedAt    time.Time `json:"captured_at"`
	CaptureSource string    `json:"capture_source"`
}

type MaintenancePDF struct {
	Bytes    []byte
	Filename string
}
