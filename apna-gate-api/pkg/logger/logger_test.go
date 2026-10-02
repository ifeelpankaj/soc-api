package logger

import (
	"errors"
	"testing"
)

func TestIsIgnorableSyncError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil", want: false},
		{name: "windows stdout invalid handle", err: errors.New("sync /dev/stdout: The handle is invalid."), want: true},
		{name: "stderr invalid argument", err: errors.New("sync /dev/stderr: invalid argument"), want: true},
		{name: "real file sync error", err: errors.New("sync logs/info.production.log: access denied"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isIgnorableSyncError(tt.err); got != tt.want {
				t.Fatalf("isIgnorableSyncError() = %v, want %v", got, tt.want)
			}
		})
	}
}
