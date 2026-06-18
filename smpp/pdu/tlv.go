package pdu

import (
	"bytes"
	"encoding/binary"
	"io"
)

type TLV struct {
	Tag    uint16
	Length uint16
	Value  []byte
}

func parseTLVLine(r *bytes.Reader) (*TLV, error) {
	var result TLV
	err := binary.Read(r, binary.BigEndian, &result.Tag)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &result.Length)
	if err != nil {
		return nil, err
	}

	result.Value = make([]byte, result.Length)
	_, err = io.ReadFull(r, result.Value)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func parseTLVs(r *bytes.Reader) ([]TLV, error) {
	var result []TLV
	if r.Len() == 0 {
		return result, nil
	}
	for r.Len() > 0 {
		tlv, err := parseTLVLine(r)
		if err != nil {
			return nil, err
		}
		result = append(result, *tlv)
	}
	return result, nil
}
