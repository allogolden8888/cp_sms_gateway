package pdu

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

func ParsePDUHeader(r *bytes.Reader) (*PDUHeader, error) {

	var h PDUHeader

	if r.Len() < 16 {
		return nil, fmt.Errorf("PDU too short: got %d bytes, need 16", r.Len())
	}

	err := binary.Read(r, binary.BigEndian, &h)
	if err != nil {
		return nil, err
	}

	return &h, nil
}
