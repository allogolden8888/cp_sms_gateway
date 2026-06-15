package smpp

import (
	"testing"
)

func TestParseDLR(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *DLR
		wantErr bool
	}{
		{
			name:  "Base - fully correct",
			input: "id:abc123 sub:001 dlvrd:001 submit date:2601011200 done date:2601011201 stat:DELIVRD err:000 Text:Test",
			want: &DLR{MessageID: "abc123", Submitted: "001", Delivered: "001",
				SubmitDate: "2601011200", DoneDate: "2601011201", Status: "DELIVRD", ErrorCode: "000", Text: "Test"},
		},
		{
			name:    "Empty string",
			input:   "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Incorrect without 'Text'",
			input:   "id:abc123 sub:001 dlvrd:001 submit date:2601011200 done date:2601011201 stat:DELIVRD err:000",
			want:    nil,
			wantErr: true,
		},
		{
			name:  "Correct and text with multiple words",
			input: "id:abc123 sub:001 dlvrd:001 submit date:2601011200 done date:2601011201 stat:DELIVRD err:000 Text:Test Test",
			want: &DLR{MessageID: "abc123", Submitted: "001", Delivered: "001",
				SubmitDate: "2601011200", DoneDate: "2601011201", Status: "DELIVRD", ErrorCode: "000", Text: "Test Test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDLR(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("%v test: got %v. expected %v", tt.name, got, tt.want)
			} else {
				if *got != *tt.want {
					t.Fatalf("%v test failed! %v expected, got %v", tt.name, tt.want, got)
				}
			}
		})
	}
}
