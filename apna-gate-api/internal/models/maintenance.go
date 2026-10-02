package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type MaintenanceSettings struct {
	FirstEnabledMonth string           `json:"first_enabled_month,omitempty"`
	Enabled           bool             `json:"enabled"`
	PricingModel      string           `json:"pricing_model"`
	Currency          string           `json:"currency"`
	BillingDay        int              `json:"billing_day"`
	DueDay            int              `json:"due_day"`
	Timezone          string           `json:"timezone"`
	FixedPaise        int64            `json:"fixed_paise"`
	AreaRatePaise     int64            `json:"area_rate_paise"`
	TypeRates         map[string]int64 `json:"type_rates"`
	EligibleStatuses  []string         `json:"eligible_flat_statuses,omitempty"`
}

func DefaultMaintenanceSettings() MaintenanceSettings {
	return MaintenanceSettings{PricingModel: "fixed", Currency: "INR", BillingDay: 1, DueDay: 10, Timezone: "Asia/Kolkata", TypeRates: map[string]int64{}}
}

func (s MaintenanceSettings) Validate() error {
	if s.Currency != "INR" {
		return errors.New("currency must be INR")
	}
	if s.BillingDay < 1 || s.BillingDay > 28 || s.DueDay < 1 || s.DueDay > 28 {
		return errors.New("billing_day and due_day must be between 1 and 28")
	}
	if s.Timezone == "" || s.Timezone == "Local" {
		return errors.New("an explicit IANA timezone is required")
	}
	if _, err := time.LoadLocation(s.Timezone); err != nil {
		return errors.New("invalid IANA timezone")
	}
	if s.FixedPaise < 0 || s.AreaRatePaise < 0 {
		return errors.New("rates cannot be negative")
	}
	for k, v := range s.TypeRates {
		if strings.TrimSpace(k) != k || k == "" || len(k) > 50 || v <= 0 {
			return errors.New("type_rates require nonempty trimmed types up to 50 bytes and positive paise amounts")
		}
	}
	switch s.PricingModel {
	case "fixed":
		if s.Enabled && s.FixedPaise <= 0 {
			return errors.New("fixed_paise must be positive")
		}
	case "per_sqft":
		if s.Enabled && s.AreaRatePaise <= 0 {
			return errors.New("area_rate_paise must be positive")
		}
	case "flat_type":
		if s.Enabled && len(s.TypeRates) == 0 {
			return errors.New("type_rates are required")
		}
	case "hybrid":
		if s.Enabled && (s.FixedPaise <= 0 || s.AreaRatePaise <= 0) {
			return errors.New("hybrid requires positive fixed and area rates")
		}
	default:
		return errors.New("pricing_model must be fixed, per_sqft, flat_type or hybrid")
	}
	return nil
}

// ParseFlatArea stores hundredths of a square foot without floating point.
func ParseFlatArea(value string) (int64, error) {
	parts := strings.Split(value, ".")
	if len(parts) > 2 || parts[0] == "" {
		return 0, errors.New("area_sqft must be a positive decimal with at most two fractional digits")
	}
	fraction := "00"
	if len(parts) == 2 {
		if len(parts[1]) < 1 || len(parts[1]) > 2 {
			return 0, errors.New("area_sqft accepts at most two fractional digits")
		}
		fraction = parts[1] + strings.Repeat("0", 2-len(parts[1]))
	}
	digits := parts[0] + fraction
	for _, c := range digits {
		if c < '0' || c > '9' {
			return 0, errors.New("invalid area_sqft")
		}
	}
	n, err := strconv.ParseInt(digits, 10, 64)
	if err != nil || n <= 0 {
		return 0, errors.New("area_sqft must be positive and within int64 range")
	}
	return n, nil
}
func FlatAreaString(n *int64) *string {
	if n == nil {
		return nil
	}
	v := fmt.Sprintf("%d.%02d", *n/100, *n%100)
	return &v
}
func FlatAreaValue(s *string) (*int64, error) {
	if s == nil {
		return nil, nil
	}
	n, err := ParseFlatArea(*s)
	return &n, err
}
func ValidateFlatBillingFields(kind, area *string) error {
	if kind != nil && (strings.TrimSpace(*kind) != *kind || *kind == "" || len(*kind) > 50) {
		return errors.New("flat_type must contain 1 to 50 bytes without surrounding spaces")
	}
	if area != nil {
		_, err := ParseFlatArea(*area)
		return err
	}
	return nil
}

type MaintenanceMonthRequest struct {
	BillingMonth               string `json:"billing_month"`
	CatchUp                    bool   `json:"catch_up"`
	ReviewToken                string `json:"review_token,omitempty"`
	CurrentPricingAcknowledged bool   `json:"current_pricing_acknowledged"`
}
type MaintenanceItem struct {
	Description string `json:"description"`
	AmountPaise int64  `json:"amount_paise"`
}
type MaintenanceFlat struct {
	ID             int64          `json:"id"`
	FlatNumber     string         `json:"flat_number"`
	Block          *string        `json:"block,omitempty"`
	FlatType       *string        `json:"flat_type,omitempty"`
	AreaHundredths *int64         `json:"area_sqft_hundredths,omitempty"`
	BilledParty    map[string]any `json:"billed_party,omitempty"`
}
type MaintenanceBill struct {
	DisplayStatus          string                  `json:"display_status"`
	DueMessage             string                  `json:"due_message"`
	PaidOn                 string                  `json:"paid_on,omitempty"`
	Outstanding            *MaintenanceOutstanding `json:"-"`
	IsCatchUp              bool                    `json:"is_catch_up"`
	IssuedAt               time.Time               `json:"issued_at"`
	Issuer                 MaintenanceIssuer       `json:"issuer"`
	PaidAmountPaise        int64                   `json:"paid_amount_paise"`
	OutstandingAmountPaise int64                   `json:"outstanding_amount_paise"`
	PaymentClaimStatus     string                  `json:"payment_claim_status"`
	ID                     int64                   `json:"id"`
	RunID                  int64                   `json:"run_id"`
	SocietyID              int64                   `json:"society_id"`
	FlatID                 int64                   `json:"flat_id"`
	BillNumber             string                  `json:"bill_number"`
	BillingMonth           string                  `json:"billing_month"`
	DueDate                string                  `json:"due_date"`
	Timezone               string                  `json:"timezone"`
	Currency               string                  `json:"currency"`
	TotalPaise             int64                   `json:"total_paise"`
	Status                 string                  `json:"status"`
	Flat                   MaintenanceFlat         `json:"flat"`
	Items                  []MaintenanceItem       `json:"items"`
	BilledParty            map[string]any          `json:"billed_party,omitempty"`
	CreatedAt              time.Time               `json:"created_at"`
}
type MaintenanceRunResult struct {
	RunID    int64 `json:"run_id"`
	Created  int   `json:"created"`
	Existing int   `json:"existing"`
}
type MaintenanceFlatIssue struct {
	FlatID int64  `json:"flat_id"`
	Reason string `json:"reason"`
}
type MaintenancePreview struct {
	ReviewToken string                 `json:"review_token,omitempty"`
	ExpiresAt   *time.Time             `json:"expires_at,omitempty"`
	IsCatchUp   bool                   `json:"is_catch_up"`
	Bills       []MaintenanceBill      `json:"bills"`
	Issues      []MaintenanceFlatIssue `json:"issues"`
	TotalPaise  int64                  `json:"total_paise"`
}
type MaintenanceBillList struct {
	Items      []MaintenanceBill `json:"items"`
	NextCursor *int64            `json:"next_cursor,omitempty"`
	TotalCount *int64            `json:"total_count,omitempty"`
	Page       int32             `json:"page,omitempty"`
	TotalPages int32             `json:"total_pages"`
	HasMore    bool              `json:"has_more"`
}
type MaintenanceBillFilter struct {
	SocietyID, UserID, ID, FlatID, BeforeID int64
	Month, Status, BillNumber               string
	Block, FlatNumber, Search               string
	Limit                                   int32
	Resident                                bool
	DisplayStatus                           string
	Page, Offset                            int32
}

type MaintenanceOutstandingFlatSummary struct {
	Flat                  MaintenanceFlat `json:"flat"`
	TotalOutstandingPaise int64           `json:"total_outstanding_paise"`
	UnpaidBillCount       int64           `json:"unpaid_bill_count"`
	OldestUnpaidMonth     string          `json:"oldest_unpaid_month,omitempty"`
}
type MaintenanceOutstandingFlatList struct {
	Items      []MaintenanceOutstandingFlatSummary `json:"items"`
	NextCursor *int64                              `json:"next_cursor,omitempty"`
}
type MaintenanceOutstandingFlatFilter struct {
	SocietyID, UserID, BeforeID int64
	Block, FlatNumber, Search   string
	Limit                       int32
}

type MaintenanceSettingsResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    MaintenanceSettings `json:"data"`
}

type MaintenancePreviewResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    MaintenancePreview `json:"data"`
}

type MaintenanceRunResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Data    MaintenanceRunResult `json:"data"`
}

type MaintenanceBillResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    MaintenanceBill `json:"data"`
}

type MaintenanceBillListResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    MaintenanceBillList `json:"data"`
}

type MaintenanceRunStatus struct {
	Status    string `json:"status"`
	RunID     int64  `json:"run_id,omitempty"`
	BillCount int    `json:"bill_count"`
}

type MaintenanceOutstandingBill struct {
	ID               int64  `json:"id"`
	BillNumber       string `json:"bill_number"`
	BillingMonth     string `json:"billing_month"`
	DueDate          string `json:"due_date"`
	OutstandingPaise int64  `json:"outstanding_paise"`
	Status           string `json:"status"`
}
type MaintenanceOutstanding struct {
	OverduePaise             int64                        `json:"overdue_paise"`
	TotalPaidPaise           int64                        `json:"total_paid_paise"`
	CurrentBill              *MaintenanceBill             `json:"current_bill,omitempty"`
	FlatID                   int64                        `json:"flat_id"`
	Currency                 string                       `json:"currency"`
	CurrentMonth             string                       `json:"current_month"`
	CalculatedAt             time.Time                    `json:"calculated_at"`
	CurrentMonthPaise        int64                        `json:"current_month_paise"`
	PreviousOutstandingPaise int64                        `json:"previous_outstanding_paise"`
	TotalOutstandingPaise    int64                        `json:"total_outstanding_paise"`
	UnpaidBills              []MaintenanceOutstandingBill `json:"unpaid_bills"`
}
