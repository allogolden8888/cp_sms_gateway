package smpp

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type PDUHeader struct {
	Length         uint32
	CommandID      uint32
	CommandStatus  uint32
	SequenceNumber uint32
}

func ParsePDUHeader(data []byte) (*PDUHeader, error) {
	reader := bytes.NewReader(data)

	var h PDUHeader

	if len(data) < 16 {
		return nil, fmt.Errorf("PDU too short: got %d bytes, need 16", len(data))
	}

	err := binary.Read(reader, binary.BigEndian, &h)
	if err != nil {
		return nil, err
	}

	return &h, nil
}

func readCString(r *bytes.Reader) (string, error) {
	for i := range len(r) {

	}
}
