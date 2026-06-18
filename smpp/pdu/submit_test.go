package pdu

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestParseSubmitBody(t *testing.T) {
	test := bytes.NewReader(makeSubmitSMBytes())
	want := &SubmitSMBody{
		ServiceType:        "",
		SourceAddrTON:      1,
		SourceAddrNPI:      1,
		SourceAddr:         "sender",
		DestAddrTON:        1,
		DestAddrNPI:        1,
		DestinationAddr:    "79001234567",
		RegisteredDelivery: 1,
		SMLength:           5,
		ShortMessage:       []byte("Hello"),
	}

	got, err := parseSubmitSMBody(test)
	if err != nil {
		t.Errorf("parse error: %v", err.Error())
	}
	if want.ServiceType != got.ServiceType {
		t.Errorf("Expected %v, got %v.", want.ServiceType, got.ServiceType)
	}
	if want.SourceAddrNPI != got.SourceAddrNPI {
		t.Errorf("Expected %v, got %v.", want.SourceAddrNPI, got.SourceAddrNPI)
	}
	if want.SourceAddr != got.SourceAddr {
		t.Errorf("Expected %v, got %v.", want.SourceAddr, got.SourceAddr)
	}
	if want.SourceAddrTON != got.SourceAddrTON {
		t.Errorf("Expected %v, got %v.", want.SourceAddrTON, got.SourceAddrTON)
	}
	if want.DestAddrTON != got.DestAddrTON {
		t.Errorf("Expected %v, got %v.", want.DestAddrTON, got.DestAddrTON)
	}
	if want.DestAddrNPI != got.DestAddrNPI {
		t.Errorf("Expected %v, got %v.", want.DestAddrNPI, got.DestAddrNPI)
	}
	if want.DestinationAddr != got.DestinationAddr {
		t.Errorf("Expected %v, got %v.", want.DestinationAddr, got.DestinationAddr)
	}
	if want.RegisteredDelivery != got.RegisteredDelivery {
		t.Errorf("Expected %v, got %v.", want.RegisteredDelivery, got.RegisteredDelivery)
	}
	if want.SMLength != got.SMLength {
		t.Errorf("Expected %v, got %v.", want.SMLength, got.SMLength)
	}
	if !bytes.Equal(got.ShortMessage, want.ShortMessage) {
		t.Errorf("Expected %v, got %v.", want.ShortMessage, got.ShortMessage)
	}

}

func TestParseSubmitSM(t *testing.T) {
	test := bytes.NewReader(makeFullSubmitSMBytes())
	want := &SubmitSM{
		PDUHeader: PDUHeader{
			Length:         65,
			CommandID:      5,
			CommandStatus:  0,
			SequenceNumber: 1,
		},
		SubmitSMBody: SubmitSMBody{
			ServiceType:        "",
			SourceAddrTON:      1,
			SourceAddrNPI:      1,
			SourceAddr:         "sender",
			DestAddrTON:        1,
			DestAddrNPI:        1,
			DestinationAddr:    "79001234567",
			RegisteredDelivery: 1,
			SMLength:           5,
			ShortMessage:       []byte("Hello"),
		},
		TLVs: []TLV{
			{
				Tag:    0x001E,
				Length: 7,
				Value:  []byte("Hellop\x00"),
			},
		},
	}

	got, err := parseSubmitSM(test)
	if err != nil {
		t.Fatalf("parse error: %v", err.Error())
	}

	if want.Length != got.Length {
		t.Errorf("Expected %v, got %v.", want.Length, got.Length)
	}
	if want.CommandID != got.CommandID {
		t.Errorf("Expected %v, got %v.", want.CommandID, got.CommandID)
	}
	if want.CommandStatus != got.CommandStatus {
		t.Errorf("Expected %v, got %v.", want.CommandStatus, got.CommandStatus)
	}
	if want.SequenceNumber != got.SequenceNumber {
		t.Errorf("Expected %v, got %v.", want.SequenceNumber, got.SequenceNumber)
	}

	if want.ServiceType != got.ServiceType {
		t.Errorf("Expected %v, got %v.", want.ServiceType, got.ServiceType)
	}
	if want.SourceAddrNPI != got.SourceAddrNPI {
		t.Errorf("Expected %v, got %v.", want.SourceAddrNPI, got.SourceAddrNPI)
	}
	if want.SourceAddr != got.SourceAddr {
		t.Errorf("Expected %v, got %v.", want.SourceAddr, got.SourceAddr)
	}
	if want.SourceAddrTON != got.SourceAddrTON {
		t.Errorf("Expected %v, got %v.", want.SourceAddrTON, got.SourceAddrTON)
	}
	if want.DestAddrTON != got.DestAddrTON {
		t.Errorf("Expected %v, got %v.", want.DestAddrTON, got.DestAddrTON)
	}
	if want.DestAddrNPI != got.DestAddrNPI {
		t.Errorf("Expected %v, got %v.", want.DestAddrNPI, got.DestAddrNPI)
	}
	if want.DestinationAddr != got.DestinationAddr {
		t.Errorf("Expected %v, got %v.", want.DestinationAddr, got.DestinationAddr)
	}
	if want.RegisteredDelivery != got.RegisteredDelivery {
		t.Errorf("Expected %v, got %v.", want.RegisteredDelivery, got.RegisteredDelivery)
	}
	if want.SMLength != got.SMLength {
		t.Errorf("Expected %v, got %v.", want.SMLength, got.SMLength)
	}
	if !bytes.Equal(got.ShortMessage, want.ShortMessage) {
		t.Errorf("Expected %v, got %v.", want.ShortMessage, got.ShortMessage)
	}

	if len(got.TLVs) != len(want.TLVs) {
		t.Fatalf("want TLVs length doesnt match got TLVs length")
	}
	for i := range len(want.TLVs) {
		comparing1 := want.TLVs[i]
		comparing2 := got.TLVs[i]

		if comparing1.Tag != comparing2.Tag {
			t.Errorf("Expected %v, got %v", comparing1.Tag, comparing2.Tag)
		}

		if comparing1.Length != comparing2.Length {
			t.Errorf("Expected %v, got %v", comparing1.Length, comparing2.Length)
		}

		if !bytes.Equal(comparing1.Value, comparing2.Value) {
			t.Errorf("Expected %v, got %v", comparing1.Value, comparing2.Value)
		}

	}

}

func makeSubmitSMBytes() []byte {
	var buf bytes.Buffer

	buf.WriteByte(0x00)                // service_type: "" (пустая C-строка)
	buf.WriteByte(0x01)                // source_addr_ton
	buf.WriteByte(0x01)                // source_addr_npi
	buf.WriteString("sender\x00")      // source_addr
	buf.WriteByte(0x01)                // dest_addr_ton
	buf.WriteByte(0x01)                // dest_addr_npi
	buf.WriteString("79001234567\x00") // destination_addr
	buf.WriteByte(0x00)                // esm_class
	buf.WriteByte(0x00)                // protocol_id
	buf.WriteByte(0x00)                // priority_flag
	buf.WriteByte(0x00)                // schedule_delivery_time: ""
	buf.WriteByte(0x00)                // validity_period: ""
	buf.WriteByte(0x01)                // registered_delivery
	buf.WriteByte(0x00)                // replace_if_present_flag
	buf.WriteByte(0x00)                // data_coding
	buf.WriteByte(0x00)                // sm_default_msg_id
	buf.WriteByte(0x05)                // sm_length = 5
	buf.WriteString("Hello")           // short_message (5 байт)

	return buf.Bytes()
}

func makeFullSubmitSMBytes() []byte {
	var buf bytes.Buffer
	// TODO: ошибки binary.Write игнорируются. bytes.Buffer технически не может вернуть ошибку,
	// но линтеры будут ругаться. Вынести в must-хелпер: mustWrite(binary.Write(...))
	binary.Write(&buf, binary.BigEndian, PDUHeader{
		Length:         65,
		CommandID:      5,
		CommandStatus:  0,
		SequenceNumber: 1,
	})
	buf.WriteByte(0x00)                // service_type: "" (пустая C-строка)
	buf.WriteByte(0x01)                // source_addr_ton
	buf.WriteByte(0x01)                // source_addr_npi
	buf.WriteString("sender\x00")      // source_addr
	buf.WriteByte(0x01)                // dest_addr_ton
	buf.WriteByte(0x01)                // dest_addr_npi
	buf.WriteString("79001234567\x00") // destination_addr
	buf.WriteByte(0x00)                // esm_class
	buf.WriteByte(0x00)                // protocol_id
	buf.WriteByte(0x00)                // priority_flag
	buf.WriteByte(0x00)                // schedule_delivery_time: ""
	buf.WriteByte(0x00)                // validity_period: ""
	buf.WriteByte(0x01)                // registered_delivery
	buf.WriteByte(0x00)                // replace_if_present_flag
	buf.WriteByte(0x00)                // data_coding
	buf.WriteByte(0x00)                // sm_default_msg_id
	buf.WriteByte(0x05)                // sm_length = 5
	buf.WriteString("Hello")           // short_message (5 байт)

	binary.Write(&buf, binary.BigEndian, uint16(0x001E))
	binary.Write(&buf, binary.BigEndian, uint16(7))
	buf.WriteString("Hellop\x00")

	return buf.Bytes()
}
