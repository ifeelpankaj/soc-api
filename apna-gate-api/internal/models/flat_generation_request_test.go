package models

import (
	"strings"
	"testing"
)

func int32Pointer(value int32) *int32    { return &value }
func stringPointer(value string) *string { return &value }

func TestGenerateFlatsRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     GenerateFlatsRequest
		wantErr string
	}{
		{
			name: "continuous permits zero and defaults floor",
			req: GenerateFlatsRequest{
				Block: " A ", NumberingMode: FlatNumberingModeContinuous, FlatsPerFloor: 3,
				SequenceStart: stringPointer("000"), TotalFlats: int32Pointer(7),
			},
		},
		{
			name: "floor based permits zero unit and floor",
			req: GenerateFlatsRequest{
				Block: "A", NumberingMode: FlatNumberingModeFloorBased, StartFloor: int32Pointer(0), FlatsPerFloor: 3,
				UnitStart: stringPointer("00"), NumberOfFloors: int32Pointer(2),
			},
		},
		{
			name: "rejects negative floor",
			req: GenerateFlatsRequest{
				Block: "A", NumberingMode: FlatNumberingModeContinuous, StartFloor: int32Pointer(-1), FlatsPerFloor: 3,
				SequenceStart: stringPointer("001"), TotalFlats: int32Pointer(3),
			},
			wantErr: "start_floor",
		},
		{
			name: "rejects non decimal start",
			req: GenerateFlatsRequest{
				Block: "A", NumberingMode: FlatNumberingModeContinuous, FlatsPerFloor: 3,
				SequenceStart: stringPointer("A01"), TotalFlats: int32Pointer(3),
			},
			wantErr: "sequence_start",
		},
		{
			name: "rejects fields from other mode",
			req: GenerateFlatsRequest{
				Block: "A", NumberingMode: FlatNumberingModeContinuous, FlatsPerFloor: 3,
				SequenceStart: stringPointer("001"), TotalFlats: int32Pointer(3), UnitStart: stringPointer("01"),
			},
			wantErr: "not allowed",
		},
		{
			name: "rejects count above cap before generation",
			req: GenerateFlatsRequest{
				Block: "A", NumberingMode: FlatNumberingModeFloorBased, FlatsPerFloor: 101,
				UnitStart: stringPointer("000"), NumberOfFloors: int32Pointer(100),
			},
			wantErr: "10000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.req.Sanitize()
			err := tt.req.Validate()
			if tt.wantErr == "" && err != nil {
				t.Fatalf("Validate() error = %v", err)
			}
			if tt.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tt.wantErr)) {
				t.Fatalf("Validate() error = %v, want containing %q", err, tt.wantErr)
			}
			if tt.name == "continuous permits zero and defaults floor" {
				if tt.req.Block != "A" || tt.req.StartFloorValue() != 1 {
					t.Fatalf("sanitized block/start floor = %q/%d", tt.req.Block, tt.req.StartFloorValue())
				}
			}
		})
	}
}
