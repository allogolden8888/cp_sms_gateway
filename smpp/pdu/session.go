package pdu

import (
	"bytes"
	"fmt"
)

type GenerickNACK struct {
	PDUHeader
}

type EnquireLink struct {
	PDUHeader
}

type EnquireLinkResp struct {
	PDUHeader
}

const (
	CommandGenericNACK     uint32 = 0x80000000
	CommandEnquireLink     uint32 = 0x00000015
	CommandEnquireLinkResp uint32 = 0x80000015
)

func ParseGenerickNACK(r *bytes.Reader) (*GenerickNACK, error) {
	header, err := ParsePDUHeader(r)
	if err != nil {
		return nil, err
	}

	if header.CommandID != CommandGenericNACK {
		return nil, fmt.Errorf("failed to parse GenerickNACK, wrong CommandID: %d", header.CommandID)
	}

	result := GenerickNACK{
		PDUHeader: *header,
	}

	return &result, nil
}

func ParseEnquireLink(r *bytes.Reader) (*EnquireLink, error) {
	header, err := ParsePDUHeader(r)
	if err != nil {
		return nil, err
	}

	if header.CommandID != CommandEnquireLink {
		return nil, fmt.Errorf("failed to parse EnquireLink, wrong CommandID: %d", header.CommandID)
	}

	result := EnquireLink{
		PDUHeader: *header,
	}

	return &result, nil
}

func ParseEnquireLinkResp(r *bytes.Reader) (*EnquireLinkResp, error) {
	header, err := ParsePDUHeader(r)
	if err != nil {
		return nil, err
	}

	if header.CommandID != CommandEnquireLinkResp {
		return nil, fmt.Errorf("failed to parse EnquireLinkResp, wrong CommandID: %d", header.CommandID)
	}

	result := EnquireLinkResp{
		PDUHeader: *header,
	}

	return &result, nil
}
