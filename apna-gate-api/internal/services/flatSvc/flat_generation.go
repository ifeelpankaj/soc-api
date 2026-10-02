package flatsvc

import (
	"errors"
	"strconv"
	"strings"

	"go-server/internal/models"
)

type flatCreateInput struct {
	FlatType   *string
	AreaSqft   *string
	Block      *string
	Floor      *string
	FlatNumber string
	Metadata   map[string]any
}

func generateFlatInputs(req *models.GenerateFlatsRequest) ([]flatCreateInput, error) {
	if req == nil {
		return nil, ErrInvalidFlatRequest
	}
	req.Sanitize()
	if err := req.Validate(); err != nil {
		return nil, ErrInvalidFlatRequest.WithCause(err)
	}

	count, err := req.GeneratedCount()
	if err != nil {
		return nil, ErrInvalidFlatRequest.WithCause(err)
	}
	block := req.Block
	items := make([]flatCreateInput, 0, int(count))

	switch req.NumberingMode {
	case models.FlatNumberingModeContinuous:
		width := len(*req.SequenceStart)
		for index := int64(0); index < count; index++ {
			flatNumber := leftPadDecimal(addDecimal(*req.SequenceStart, index), width)
			if len(flatNumber) > 50 {
				return nil, ErrInvalidFlatRequest.WithCause(errors.New("generated flat_number must be at most 50 characters"))
			}
			floor := strconv.FormatInt(req.StartFloorValue()+index/int64(req.FlatsPerFloor), 10)
			items = append(items, flatCreateInput{Block: &block, Floor: &floor, FlatNumber: flatNumber})
		}
	case models.FlatNumberingModeFloorBased:
		suffixWidth := len(*req.UnitStart)
		lastUnit := addDecimal(*req.UnitStart, int64(req.FlatsPerFloor)-1)
		if len(canonicalDecimal(lastUnit)) > suffixWidth {
			return nil, ErrInvalidFlatRequest.WithCause(errors.New("unit range exceeds the fixed unit_start width"))
		}
		for floorOffset := int64(0); floorOffset < int64(*req.NumberOfFloors); floorOffset++ {
			floor := strconv.FormatInt(req.StartFloorValue()+floorOffset, 10)
			for unitOffset := int64(0); unitOffset < int64(req.FlatsPerFloor); unitOffset++ {
				suffix := leftPadDecimal(addDecimal(*req.UnitStart, unitOffset), suffixWidth)
				flatNumber := floor + suffix
				if len(flatNumber) > 50 {
					return nil, ErrInvalidFlatRequest.WithCause(errors.New("generated flat_number must be at most 50 characters"))
				}
				floorValue := floor
				items = append(items, flatCreateInput{Block: &block, Floor: &floorValue, FlatNumber: flatNumber})
			}
		}
	default:
		return nil, ErrInvalidFlatRequest.WithCause(errors.New("invalid numbering_mode"))
	}

	return items, nil
}

func addDecimal(value string, increment int64) string {
	digits := []byte(canonicalDecimal(value))
	carry := increment
	for index := len(digits) - 1; index >= 0 && carry > 0; index-- {
		sum := int64(digits[index]-'0') + carry
		digits[index] = "0123456789"[sum%10]
		carry = sum / 10
	}
	if carry == 0 {
		return string(digits)
	}
	return strconv.FormatInt(carry, 10) + string(digits)
}

func canonicalDecimal(value string) string {
	canonical := strings.TrimLeft(value, "0")
	if canonical == "" {
		return "0"
	}
	return canonical
}

func leftPadDecimal(value string, width int) string {
	if len(value) >= width {
		return value
	}
	return strings.Repeat("0", width-len(value)) + value
}
