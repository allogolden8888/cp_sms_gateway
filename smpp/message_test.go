package smpp

import (
	"testing"
	"time"
)

func TestParseValidity(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		want         time.Duration
		wantPositive bool
		wantErr      bool
	}{
		{
			name:  "relative 1h correct",
			input: "1h",
			want:  time.Hour,
		},
		{
			name:  "relative 30m correct",
			input: "30m",
			want:  30 * time.Minute,
		},
		{
			name:         "absolute 2099-01-01T00:00:00 correct",
			input:        "2099-01-01T00:00:00",
			wantPositive: true,
		},
		{
			name:    "incorrect abc",
			input:   "abc",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseValidity(tt.input)
			if !tt.wantErr {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if tt.wantPositive {
					if got > 0 {
						return
					} else {
						t.Fatalf("%v test failed. expected positive or %v", tt.name, tt.want)
					}
				}
				if got != tt.want {
					t.Errorf("%v test failed. expected %v, got %v", tt.name, tt.want, got)
				}

			} else {
				if err == nil {
					t.Fatalf("%v test failed. expected error, got nil", tt.name)
				}
			}
		})
	}
}
