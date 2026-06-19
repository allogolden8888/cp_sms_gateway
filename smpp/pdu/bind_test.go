package pdu

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestParseBind(t *testing.T) {
	tests := []struct {
		name              string
		data              []byte
		expectedCommandID uint32
		wantErr           bool
		wantHeader        *PDUHeader
		wantBody          *BindBody
	}{
		{
			name:              "valid transmitter",
			data:              makeBindBytes(CommandBindTransmitter),
			expectedCommandID: CommandBindTransmitter,
			wantHeader: &PDUHeader{
				Length:         32,
				CommandID:      CommandBindTransmitter,
				CommandStatus:  0,
				SequenceNumber: 1,
			},
			wantBody: &BindBody{
				SystemID: "20100",
				Password: "test",
			},
		},
		{
			name:              "wrong command ID",
			data:              makeBindBytes(CommandBindTransmitter),
			expectedCommandID: CommandBindReceiver,
			wantErr:           true,
		},
		{
			name:              "too short for header",
			data:              []byte{0x00, 0x01},
			expectedCommandID: CommandBindTransmitter,
			wantErr:           true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body, header, err := parseBind(bytes.NewReader(tc.data), tc.expectedCommandID)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.wantHeader.Length != header.Length {
				t.Errorf("header.Length: want %d, got %d", tc.wantHeader.Length, header.Length)
			}
			if tc.wantHeader.CommandID != header.CommandID {
				t.Errorf("header.CommandID: want %d, got %d", tc.wantHeader.CommandID, header.CommandID)
			}
			if tc.wantHeader.CommandStatus != header.CommandStatus {
				t.Errorf("header.CommandStatus: want %d, got %d", tc.wantHeader.CommandStatus, header.CommandStatus)
			}
			if tc.wantHeader.SequenceNumber != header.SequenceNumber {
				t.Errorf("header.SequenceNumber: want %d, got %d", tc.wantHeader.SequenceNumber, header.SequenceNumber)
			}

			if tc.wantBody.SystemID != body.SystemID {
				t.Errorf("body.SystemID: want %s, got %s", tc.wantBody.SystemID, body.SystemID)
			}
			if tc.wantBody.Password != body.Password {
				t.Errorf("body.Password: want %s, got %s", tc.wantBody.Password, body.Password)
			}
		})
	}
}

func TestParseBindBody(t *testing.T) {
	test := bytes.NewReader(makeBindBodyBytes())

	want := &BindBody{
		SystemID: "20100",
		Password: "test",
	}

	got, err := parseBindBody(test)
	if err != nil {
		t.Fatalf("parse error: %v", err.Error())
	}

	if want.SystemID != got.SystemID {
		t.Errorf("Expected %v, got %v.", want.SystemID, got.SystemID)
	}
	if want.Password != got.Password {
		t.Errorf("Expected %v, got %v.", want.Password, got.Password)
	}

}

// TODO: переписать TestParseBindTransmitter, TestParseBindReceiver, TestParseBindTransceiver
// в один table-driven тест. Проблема: функции возвращают разные типы, поэтому в parseFn
// возвращай только нужные поля:
//
//   tests := []struct {
//       name      string
//       commandID uint32
//       parseFn   func(*bytes.Reader) (systemID, password string, length, commandID uint32, err error)
//   }{
//       {
//           name:      "transmitter",
//           commandID: CommandBindTransmitter,
//           parseFn: func(r *bytes.Reader) (string, string, uint32, uint32, error) {
//               got, err := ParseBindTransmitter(r)
//               if err != nil { return "", "", 0, 0, err }
//               return got.SystemID, got.Password, got.Length, got.CommandID, nil
//           },
//       },
//       // аналогично для Receiver и Transceiver
//   }
func TestParseBindTransmitter(t *testing.T) {
	test := bytes.NewReader(makeBindBytes(CommandBindTransmitter))

	want := &BindTransmitter{
		PDUHeader: PDUHeader{
			Length:         32,
			CommandID:      CommandBindTransmitter,
			CommandStatus:  0,
			SequenceNumber: 1,
		},
		BindBody: BindBody{
			SystemID: "20100",
			Password: "test",
		},
	}

	got, err := ParseBindTransmitter(test)
	if err != nil {
		t.Errorf("%v", got)
		t.Fatalf("parse error: %v", err.Error())
	}

	if want.SystemID != got.SystemID {
		t.Errorf("Expected %v, got %v.", want.SystemID, got.SystemID)
	}
	if want.Password != got.Password {
		t.Errorf("Expected %v, got %v.", want.Password, got.Password)
	}
	if want.Length != got.Length {
		t.Errorf("Expected %v, got %v.", want.Length, got.Length)
	}
	if want.CommandID != got.CommandID {
		t.Errorf("Expected %v, got %v.", want.CommandID, got.CommandID)
	}
	if got.CommandID != CommandBindTransmitter {
		t.Fatalf("Expected command ID %d, got: %d", CommandBindTransmitter, got.CommandID)
	}
}

func TestParseBindReceiver(t *testing.T) {
	test := bytes.NewReader(makeBindBytes(CommandBindReceiver))

	want := &BindReceiver{
		PDUHeader: PDUHeader{
			Length:         32,
			CommandID:      CommandBindReceiver,
			CommandStatus:  0,
			SequenceNumber: 1,
		},
		BindBody: BindBody{
			SystemID: "20100",
			Password: "test",
		},
	}

	got, err := ParseBindReceiver(test)
	if err != nil {
		t.Errorf("%v", got)
		t.Fatalf("parse error: %v", err.Error())
	}

	if want.SystemID != got.SystemID {
		t.Errorf("Expected %v, got %v.", want.SystemID, got.SystemID)
	}
	if want.Password != got.Password {
		t.Errorf("Expected %v, got %v.", want.Password, got.Password)
	}
	if want.Length != got.Length {
		t.Errorf("Expected %v, got %v.", want.Length, got.Length)
	}
	if want.CommandID != got.CommandID {
		t.Errorf("Expected %v, got %v.", want.CommandID, got.CommandID)
	}
	if got.CommandID != CommandBindReceiver {
		t.Fatalf("Expected command ID %d, got: %d", CommandBindReceiver, got.CommandID)
	}
}

func TestParseBindTransceiver(t *testing.T) {
	test := bytes.NewReader(makeBindBytes(CommandBindTransceiver))

	want := &BindTransceiver{
		PDUHeader: PDUHeader{
			Length:         32,
			CommandID:      CommandBindTransceiver,
			CommandStatus:  0,
			SequenceNumber: 1,
		},
		BindBody: BindBody{
			SystemID: "20100",
			Password: "test",
		},
	}

	got, err := ParseBindTransceiver(test)
	if err != nil {
		t.Errorf("%v", got)
		t.Fatalf("parse error: %v", err.Error())
	}

	if want.SystemID != got.SystemID {
		t.Errorf("Expected %v, got %v.", want.SystemID, got.SystemID)
	}
	if want.Password != got.Password {
		t.Errorf("Expected %v, got %v.", want.Password, got.Password)
	}
	if want.Length != got.Length {
		t.Errorf("Expected %v, got %v.", want.Length, got.Length)
	}
	if want.CommandID != got.CommandID {
		t.Errorf("Expected %v, got %v.", want.CommandID, got.CommandID)
	}
	if got.CommandID != CommandBindTransceiver {
		t.Fatalf("Expected command ID %d, got: %d", CommandBindTransceiver, got.CommandID)
	}
}

func makeBindBodyBytes() []byte {
	var buf bytes.Buffer

	buf.WriteString("20100")
	buf.WriteByte(0x00) // system_id нуль-терминатор
	buf.WriteString("test")
	buf.WriteByte(0x00) // password нуль-терминатор
	buf.WriteByte(0x00) // system_type: ""
	buf.WriteByte(0x00) // interface_version
	buf.WriteByte(0x00) // addr_ton
	buf.WriteByte(0x00) // addr_npi
	buf.WriteByte(0x00) // address_range: ""

	return buf.Bytes()
}

func makeBindBytes(commandID uint32) []byte {
	var buf bytes.Buffer

	must(binary.Write(&buf, binary.BigEndian, PDUHeader{
		Length:         0,
		CommandID:      commandID,
		CommandStatus:  0,
		SequenceNumber: 1,
	}))

	buf.WriteString("20100")
	buf.WriteByte(0x00)
	buf.WriteString("test")
	buf.WriteByte(0x00)
	buf.WriteByte(0x00)
	buf.WriteByte(0x00)
	buf.WriteByte(0x00)
	buf.WriteByte(0x00)
	buf.WriteByte(0x00)

	b := buf.Bytes()
	binary.BigEndian.PutUint32(b[0:4], uint32(len(b)))

	return b
}
