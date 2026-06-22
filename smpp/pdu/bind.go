package pdu

import (
	"bytes"
	"fmt"
)

type BindBody struct {
	SystemID         string
	Password         string
	SystemType       string
	InterfaceVersion uint8
	AddrTON          uint8
	AddrNPI          uint8
	AddressRange     string
}

type BindTransmitter struct {
	PDUHeader
	BindBody
}

type BindReceiver struct {
	PDUHeader
	BindBody
}

type BindTransceiver struct {
	PDUHeader
	BindBody
}

type BindResp struct {
	PDUHeader
	SystemID string
	TLVs     []TLV
}

type UnbindResp struct {
	PDUHeader
}

type Unbind struct {
	PDUHeader
}

const (
	CommandBindTransmitter     uint32 = 0x00000002
	CommandBindReceiver        uint32 = 0x00000001
	CommandBindTransceiver     uint32 = 0x00000009
	CommandBindTransmitterResp uint32 = 0x80000002
	CommandBindReceiverResp    uint32 = 0x80000001
	CommandBindTransceiverResp uint32 = 0x80000009
	CommandUnbind              uint32 = 0x00000006
	CommandUnbindResp          uint32 = 0x80000006
)

func ParseUnbind(r *bytes.Reader) (*Unbind, error) {
	header, err := ParsePDUHeader(r)
	if err != nil {
		return nil, err
	}

	if header.CommandID != CommandUnbind {
		return nil, fmt.Errorf("failed to parse Unbind, wrong CommandID: %d", header.CommandID)
	}

	result := Unbind{
		PDUHeader: *header,
	}

	return &result, nil
}

func ParseUnbindResp(r *bytes.Reader) (*UnbindResp, error) {
	header, err := ParsePDUHeader(r)
	if err != nil {
		return nil, err
	}

	if header.CommandID != CommandUnbindResp {
		return nil, fmt.Errorf("failed to parse UnbindResp, wrong CommandID: %d", header.CommandID)
	}

	result := UnbindResp{
		PDUHeader: *header,
	}

	return &result, nil
}


func ParseBindResp(r *bytes.Reader, commandID uint32) (*BindResp, error) {
	var systemID string
	var tlvs []TLV
	header, err := ParsePDUHeader(r)
	if err != nil {
		return nil, err
	}

	if header.CommandID != commandID {
		return nil, fmt.Errorf("failed to parse BindResp, wrong CommandID: %d", header.CommandID)
	}

	systemID, err = readCString(r)
	if err != nil {
		return nil, err
	}

	tlvs, err = parseTLVs(r)
	if err != nil {
		return nil, err
	}

	result := BindResp{
		PDUHeader: *header,
		SystemID:  systemID,
		TLVs:      tlvs,
	}

	return &result, nil
}

// TODO: заменить *bytes.Reader на io.ByteReader — функция вызывает readCString и r.ReadByte(),
// оба требуют только io.ByteReader. После изменений pdu.go и header.go это станет возможным.
// parseBind и публичные функции можно оставить на *bytes.Reader — он удовлетворяет обоим интерфейсам.
func parseBindBody(r *bytes.Reader) (*BindBody, error) {
	var result BindBody
	var err error

	result.SystemID, err = readCString(r)
	if err != nil {
		return nil, err
	}

	result.Password, err = readCString(r)
	if err != nil {
		return nil, err
	}

	result.SystemType, err = readCString(r)
	if err != nil {
		return nil, err
	}

	result.InterfaceVersion, err = r.ReadByte()
	if err != nil {
		return nil, err
	}

	result.AddrTON, err = r.ReadByte()
	if err != nil {
		return nil, err
	}

	result.AddrNPI, err = r.ReadByte()
	if err != nil {
		return nil, err
	}

	result.AddressRange, err = readCString(r)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func ParseBindTransmitter(r *bytes.Reader) (*BindTransmitter, error) {
	body, header, err := parseBind(r, CommandBindTransmitter)
	if err != nil {
		return nil, err
	}

	if header.CommandID != CommandBindTransmitter {
		return nil, fmt.Errorf("failed to parse Transmitter, wrong CommandID: %d", header.CommandID)
	}

	return &BindTransmitter{PDUHeader: *header, BindBody: *body}, nil
}

func ParseBindReceiver(r *bytes.Reader) (*BindReceiver, error) {
	body, header, err := parseBind(r, CommandBindReceiver)
	if err != nil {
		return nil, err
	}

	if header.CommandID != CommandBindReceiver {
		return nil, fmt.Errorf("failed to parse Receiver, wrong CommandID: %d", header.CommandID)
	}

	return &BindReceiver{PDUHeader: *header, BindBody: *body}, nil
}

func ParseBindTransceiver(r *bytes.Reader) (*BindTransceiver, error) {
	body, header, err := parseBind(r, CommandBindTransceiver)
	if err != nil {
		return nil, err
	}

	if header.CommandID != CommandBindTransceiver {
		return nil, fmt.Errorf("failed to parse Transceiver, wrong CommandID: %d", header.CommandID)
	}

	return &BindTransceiver{PDUHeader: *header, BindBody: *body}, nil
}

func parseBind(r *bytes.Reader, expectedCommandID uint32) (*BindBody, *PDUHeader, error) {
	var err error
	var header *PDUHeader
	var body *BindBody

	header, err = ParsePDUHeader(r)
	if err != nil {
		return nil, nil, err
	}

	if header.CommandID != expectedCommandID {
		return nil, nil, fmt.Errorf("failed to parse %v, wrong CommandID: %d", expectedCommandID, header.CommandID)
	}

	body, err = parseBindBody(r)
	if err != nil {
		return nil, nil, err
	}

	return body, header, nil
}
