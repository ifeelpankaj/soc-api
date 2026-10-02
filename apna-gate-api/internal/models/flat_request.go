package models

import (
	"errors"
	"fmt"
	"strings"
)

const MaxGeneratedFlatsPerRequest int64 = 10_000

type FlatNumberingMode string

const (
	FlatNumberingModeContinuous FlatNumberingMode = "continuous"
	FlatNumberingModeFloorBased FlatNumberingMode = "floor_based"
)

func (m FlatNumberingMode) IsValid() bool {
	return m == FlatNumberingModeContinuous || m == FlatNumberingModeFloorBased
}

type CreateFlatRequest struct {
	FlatType   *string        `json:"flat_type,omitempty"`
	AreaSqft   *string        `json:"area_sqft,omitempty"`
	Block      *string        `json:"block,omitempty" validate:"omitempty,max=50"`
	Floor      *string        `json:"floor,omitempty" validate:"omitempty,max=50"`
	FlatNumber string         `json:"flat_number" validate:"required,max=50"`
	Metadata   map[string]any `json:"metadata,omitempty"`
}

func (r *CreateFlatRequest) Sanitize() {
	r.Block = trimPtr(r.Block)
	r.Floor = trimPtr(r.Floor)
	r.FlatNumber = strings.TrimSpace(r.FlatNumber)
}

func (r *CreateFlatRequest) Validate() error {
	if err := ValidateFlatBillingFields(r.FlatType, r.AreaSqft); err != nil {
		return err
	}
	if strings.TrimSpace(r.FlatNumber) == "" {
		return errors.New("flat_number is required")
	}
	return nil
}

type BulkCreateFlatsRequest struct {
	Flats []CreateFlatRequest `json:"flats" validate:"required,min=1,dive"`
}

func (r *BulkCreateFlatsRequest) Sanitize() {
	for i := range r.Flats {
		r.Flats[i].Sanitize()
	}
}

func (r *BulkCreateFlatsRequest) Validate() error {
	if len(r.Flats) == 0 {
		return errors.New("at least one flat is required")
	}
	for i := range r.Flats {
		if err := r.Flats[i].Validate(); err != nil {
			return err
		}
	}
	return nil
}

type GenerateFlatsRequest struct {
	Block          string            `json:"block" validate:"required,max=50"`
	NumberingMode  FlatNumberingMode `json:"numbering_mode" validate:"required"`
	StartFloor     *int32            `json:"start_floor,omitempty"`
	FlatsPerFloor  int32             `json:"flats_per_floor" validate:"required,gt=0"`
	SequenceStart  *string           `json:"sequence_start,omitempty"`
	TotalFlats     *int32            `json:"total_flats,omitempty"`
	UnitStart      *string           `json:"unit_start,omitempty"`
	NumberOfFloors *int32            `json:"number_of_floors,omitempty"`
}

func (r *GenerateFlatsRequest) Sanitize() {
	if r == nil {
		return
	}
	r.Block = strings.TrimSpace(r.Block)
	r.SequenceStart = trimPtr(r.SequenceStart)
	r.UnitStart = trimPtr(r.UnitStart)
}

func (r *GenerateFlatsRequest) StartFloorValue() int64 {
	if r == nil || r.StartFloor == nil {
		return 1
	}
	return int64(*r.StartFloor)
}

func (r *GenerateFlatsRequest) GeneratedCount() (int64, error) {
	if r == nil {
		return 0, errors.New("generate flats request is required")
	}
	switch r.NumberingMode {
	case FlatNumberingModeContinuous:
		if r.TotalFlats == nil {
			return 0, errors.New("total_flats is required for continuous numbering")
		}
		return int64(*r.TotalFlats), nil
	case FlatNumberingModeFloorBased:
		if r.NumberOfFloors == nil {
			return 0, errors.New("number_of_floors is required for floor_based numbering")
		}
		return int64(*r.NumberOfFloors) * int64(r.FlatsPerFloor), nil
	default:
		return 0, errors.New("invalid numbering_mode")
	}
}

func (r *GenerateFlatsRequest) Validate() error {
	if r == nil {
		return errors.New("generate flats request is required")
	}
	if r.Block == "" {
		return errors.New("block is required")
	}
	if len(r.Block) > 50 {
		return errors.New("block must be at most 50 characters")
	}
	if !r.NumberingMode.IsValid() {
		return errors.New("numbering_mode must be continuous or floor_based")
	}
	if r.StartFloorValue() < 0 {
		return errors.New("start_floor must be zero or greater")
	}
	if r.FlatsPerFloor <= 0 {
		return errors.New("flats_per_floor must be positive")
	}

	switch r.NumberingMode {
	case FlatNumberingModeContinuous:
		if r.SequenceStart == nil || !isDecimalString(*r.SequenceStart) {
			return errors.New("sequence_start must contain decimal digits")
		}
		if len(*r.SequenceStart) > 50 {
			return errors.New("sequence_start must be at most 50 digits")
		}
		if r.TotalFlats == nil || *r.TotalFlats <= 0 {
			return errors.New("total_flats must be positive")
		}
		if r.UnitStart != nil || r.NumberOfFloors != nil {
			return errors.New("unit_start and number_of_floors are not allowed for continuous numbering")
		}
	case FlatNumberingModeFloorBased:
		if r.UnitStart == nil || !isDecimalString(*r.UnitStart) {
			return errors.New("unit_start must contain decimal digits")
		}
		if len(*r.UnitStart) > 50 {
			return errors.New("unit_start must be at most 50 digits")
		}
		if r.NumberOfFloors == nil || *r.NumberOfFloors <= 0 {
			return errors.New("number_of_floors must be positive")
		}
		if r.SequenceStart != nil || r.TotalFlats != nil {
			return errors.New("sequence_start and total_flats are not allowed for floor_based numbering")
		}
	}

	count, err := r.GeneratedCount()
	if err != nil {
		return err
	}
	if count <= 0 {
		return errors.New("generated flat count must be positive")
	}
	if count > MaxGeneratedFlatsPerRequest {
		return fmt.Errorf("generated flat count must not exceed %d", MaxGeneratedFlatsPerRequest)
	}
	return nil
}

func isDecimalString(value string) bool {
	if value == "" {
		return false
	}
	for _, character := range value {
		if character < '0' || character > '9' {
			return false
		}
	}
	return true
}

type UpdateFlatRequest struct {
	FlatType   *string        `json:"flat_type,omitempty"`
	AreaSqft   *string        `json:"area_sqft,omitempty"`
	Block      *string        `json:"block,omitempty" validate:"omitempty,max=50"`
	Floor      *string        `json:"floor,omitempty" validate:"omitempty,max=50"`
	FlatNumber *string        `json:"flat_number,omitempty" validate:"omitempty,max=50"`
	Status     *FlatStatus    `json:"status,omitempty"`
	IsActive   *bool          `json:"is_active,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
}

func (r *UpdateFlatRequest) Sanitize() {
	r.Block = trimPtr(r.Block)
	r.Floor = trimPtr(r.Floor)
	r.FlatNumber = trimPtr(r.FlatNumber)
}

func (r *UpdateFlatRequest) Validate() error {
	if err := ValidateFlatBillingFields(r.FlatType, r.AreaSqft); err != nil {
		return err
	}
	if r.Block == nil && r.Floor == nil && r.FlatNumber == nil && r.Status == nil && r.IsActive == nil && r.Metadata == nil && r.FlatType == nil && r.AreaSqft == nil {
		return errors.New("at least one field must be provided")
	}
	if r.Status != nil {
		return errors.New("flat status is managed by resident occupancy and block actions")
	}
	return nil
}

type SubmitFlatClaimRequest struct {
	SocietyID        int64            `json:"society_id" validate:"required,gt=0"`
	FlatID           int64            `json:"flat_id" validate:"required,gt=0"`
	RequestedRole    FlatResidentRole `json:"requested_role" validate:"required"`
	RequestedPrimary bool             `json:"requested_primary"`
	Note             *string          `json:"note,omitempty" validate:"omitempty,max=500"`
	Metadata         map[string]any   `json:"metadata,omitempty"`
}

func (r *SubmitFlatClaimRequest) Sanitize() {
	r.Note = trimPtr(r.Note)
}

func (r *SubmitFlatClaimRequest) Validate() error {
	if r.SocietyID <= 0 || r.FlatID <= 0 {
		return errors.New("society_id and flat_id are required")
	}
	if !r.RequestedRole.IsValid() {
		return errors.New("invalid requested role")
	}
	return nil
}

type RejectFlatClaimRequest struct {
	Reason string `json:"reason" validate:"required,max=500"`
}

func (r *RejectFlatClaimRequest) Sanitize() {
	r.Reason = strings.TrimSpace(r.Reason)
}

type AddFlatResidentRequest struct {
	Role      FlatResidentRole `json:"role" validate:"required"`
	IsPrimary bool             `json:"is_primary"`
	Metadata  map[string]any   `json:"metadata,omitempty"`
}

func (r *AddFlatResidentRequest) Validate() error {
	if !r.Role.IsValid() {
		return errors.New("invalid resident role")
	}
	return nil
}

type UpdateFlatResidentRoleRequest struct {
	Role FlatResidentRole `json:"role" validate:"required"`
}

func (r *UpdateFlatResidentRoleRequest) Validate() error {
	if !r.Role.IsValid() {
		return errors.New("invalid resident role")
	}
	return nil
}

type FlatFilter struct {
	ID        *int64
	SocietyID *int64

	Block      *string
	Floor      *string
	FlatNumber *string
	Status     *string
	IsActive   *bool

	Search     string
	SearchMode string

	Limit  int32
	Offset int32
}

type FlatResidentFilter struct {
	ID        *int64
	SocietyID *int64
	FlatID    *int64
	UserID    *int64

	Role      *string
	Status    *string
	IsPrimary *bool

	Search     string
	SearchMode string

	Limit  int32
	Offset int32
}

type FlatClaimFilter struct {
	ID        *int64
	SocietyID *int64
	FlatID    *int64
	UserID    *int64

	Status *string

	Search     string
	SearchMode string

	Limit  int32
	Offset int32
}
