package smpp

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestParsePDUHeader(t *testing.T) {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, PDUHeader{
		Length:         16,
		CommandID:      5,
		CommandStatus:  0,
		SequenceNumber: 1,
	})
	data := buf.Bytes() // []byte длиной 16

	tests := []struct {
		name    string
		input   []byte
		want    *PDUHeader
		wantErr bool
	}{
		{
			name:  "valid header",
			input: data,
			want: &PDUHeader{
				Length:         16,
				CommandID:      5,
				CommandStatus:  0,
				SequenceNumber: 1,
			},
		},
		{
			name:    "invalid short header",
			input:   []byte{0x00, 0x01},
			wantErr: true,
		},
		{
			name:    "invalid empty header",
			input:   []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePDUHeader(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("test %v failed, expected error, got nil", tt.name)
				}
				return
			}
			if err != nil {
				t.Fatalf("test %v failed, error %v", tt.name, err)
			}
			if *tt.want != *got {
				t.Fatalf("test %v failed, expected %v, got %v", tt.name, tt.want, got)
			}
		})
	}
}
