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

// TODO: заменить *bytes.Reader на io.Reader — здесь нужен только binary.Read, который принимает io.Reader.
// Убрать проверку r.Len() < 16: binary.Read сам вернёт ошибку если байт не хватит.
// Убрать импорт "bytes", добавить "io".
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
