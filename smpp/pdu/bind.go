package pdu

import (
	"bytes"
	"encoding/binary"
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

const (
	CommandBindTransmitter     uint32 = 0x00000002
	CommandBindReceiver        uint32 = 0x00000001
	CommandBindTransceiver     uint32 = 0x00000009
	CommandGenerickNACK        uint32 = 0x80000000
	CommandBindTransmitterResp uint32 = 0x80000002
	CommandBindReceiverResp    uint32 = 0x80000001
	CommandBindTransceiverResp uint32 = 0x80000009
	CommandUnbind              uint32 = 0x00000006
	CommandUnbindResp          uint32 = 0x80000006
)

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

	err = binary.Read(r, binary.BigEndian, &result.InterfaceVersion)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &result.AddrTON)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &result.AddrNPI)
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
	var result BindTransmitter
	var err error
	var header *PDUHeader
	var body *BindBody

	header, err = ParsePDUHeader(r)
	if err != nil {
		return nil, err
	}

	if header.CommandID != CommandBindTransmitter {
		return nil, fmt.Errorf("failed to parse Transmitter, wrong CommandID: %d", header.CommandID)
	}

	result.PDUHeader = *header

	body, err = parseBindBody(r)
	if err != nil {
		return nil, err
	}

	result.BindBody = *body

	return &result, nil
}

func ParseBindReceiver(r *bytes.Reader) (*BindReceiver, error) {
	var result BindReceiver
	var err error
	var header *PDUHeader
	var body *BindBody

	header, err = ParsePDUHeader(r)
	if err != nil {
		return nil, err
	}

	if header.CommandID != CommandBindReceiver {
		return nil, fmt.Errorf("failed to parse Receiver, wrong CommandID: %d", header.CommandID)
	}

	result.PDUHeader = *header

	body, err = parseBindBody(r)
	if err != nil {
		return nil, err
	}

	result.BindBody = *body

	return &result, nil
}

func ParseBindTransceiver(r *bytes.Reader) (*BindTransceiver, error) {
	var result BindTransceiver
	var err error
	var header *PDUHeader
	var body *BindBody

	header, err = ParsePDUHeader(r)
	if err != nil {
		return nil, err
	}

	if header.CommandID != CommandBindTransceiver {
		return nil, fmt.Errorf("failed to parse Transceiver, wrong CommandID: %d", header.CommandID)
	}

	result.PDUHeader = *header

	body, err = parseBindBody(r)
	if err != nil {
		return nil, err
	}

	result.BindBody = *body

	return &result, nil
}
