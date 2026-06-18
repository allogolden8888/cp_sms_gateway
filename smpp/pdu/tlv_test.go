package pdu

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestParseTLVLine(t *testing.T) {
	test_case := bytes.NewReader(makeTLVBytes())
	want := TLV{
		Tag:    0x001E,
		Length: 7,
		Value:  []byte("Hellop\x00"),
	}

	got, err := parseTLVLine(test_case)
	if err != nil {
		t.Errorf("Expected %v, got err: %v", want, err)
	}

	if want.Tag != got.Tag {
		t.Errorf("Expected %v, got %v", want.Tag, got.Tag)
	}

	if want.Length != got.Length {
		t.Errorf("Expected %v, got %v", want.Length, got.Length)
	}

	if !bytes.Equal(want.Value, got.Value) {
		t.Errorf("Expected %v, got %v", want.Value, got.Value)
	}
}

func TestParseTLVs(t *testing.T) {
	var got []TLV
	var want []TLV

	test := bytes.NewReader(makeTLVBytes())
	test2 := bytes.NewReader(makeTLVBytes())

	buf, err := parseTLVLine(test)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	want = append(want, *buf)

	got, err = parseTLVs(test2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	for i := range len(want) {
		comparing1 := want[i]
		comparing2 := got[i]

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

func makeTLVBytes() []byte {
	var buf bytes.Buffer

	binary.Write(&buf, binary.BigEndian, uint16(0x001E))
	binary.Write(&buf, binary.BigEndian, uint16(7))
	buf.WriteString("Hellop\x00") // short_message (5 байт)

	return buf.Bytes()
}
