package pdu

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

	test1 := bytes.NewReader(data)
	test2 := bytes.NewReader([]byte{0x00, 0x01})
	test3 := bytes.NewReader([]byte{})

	tests := []struct {
		name    string
		input   *bytes.Reader
		want    *PDUHeader
		wantErr bool
	}{
		{
			name:  "valid header",
			input: test1,
			want: &PDUHeader{
				Length:         16,
				CommandID:      5,
				CommandStatus:  0,
				SequenceNumber: 1,
			},
		},
		{
			name:    "invalid short header",
			input:   test2,
			wantErr: true,
		},
		{
			name:    "invalid empty header",
			input:   test3,
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
