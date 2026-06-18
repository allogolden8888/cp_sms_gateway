package pdu

import (
	"bytes"
	"testing"
)

func TestReadCString(t *testing.T) {
	test1 := bytes.NewReader([]byte("hello\x00"))
	test2 := bytes.NewReader([]byte("\x00"))
	test3 := bytes.NewReader([]byte("hello"))

	tests := []struct {
		name    string
		input   *bytes.Reader
		want    string
		wantErr bool
	}{
		{
			name:  "valid bytes",
			input: test1,
			want:  "hello",
		},
		{
			name:  "valid empty bytes",
			input: test2,
			want:  "",
		},
		{
			name:    "invalid bytes with no zero byte",
			input:   test3,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readCString(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("test %v failed, expected error, got nil", tt.name)
				}
				return
			}
			if err != nil {
				t.Fatalf("test %v failed, error %v", tt.name, err)
			}
			if tt.want != got {
				t.Fatalf("test %v failed, expected %v, got %v", tt.name, tt.want, got)
			}
		})
	}

}
