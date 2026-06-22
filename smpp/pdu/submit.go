package pdu

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type SubmitSMBody struct {
	ServiceType          string
	SourceAddrTON        uint8
	SourceAddrNPI        uint8
	SourceAddr           string
	DestAddrTON          uint8
	DestAddrNPI          uint8
	DestinationAddr      string
	ESMClass             uint8
	ProtocolID           uint8
	PriorityFlag         uint8
	ScheduleDeliveryTime string
	ValidityPeriod       string
	RegisteredDelivery   uint8
	ReplaceIfPresentFlag uint8
	DataCoding           uint8
	SMDefaultMsgID       uint8
	SMLength             uint8
	ShortMessage         []byte
}

type SubmitSM struct {
	PDUHeader
	SubmitSMBody
	TLVs []TLV
}

type SubmitSMResp struct {
	PDUHeader
	MessageID string
}

const (
	CommandSubmitSM     uint32 = 0x00000004
	CommandSubmitSMResp uint32 = 0x80000004
)

func ParseSubmitSMResp(r *bytes.Reader) (*SubmitSMResp, error) {
	var messageID string
	header, err := ParsePDUHeader(r)
	if err != nil {
		return nil, err
	}

	if header.CommandID != CommandSubmitSMResp {
		return nil, fmt.Errorf("failed to parse SubmitSMResp, wrong CommandID: %d", header.CommandID)
	}

	messageID, err = readCString(r)
	if err != nil {
		return nil, err
	}

	result := SubmitSMResp{
		PDUHeader: *header,
		MessageID: messageID,
	}

	return &result, nil
}

// TODO: принимать io.Reader вместо *bytes.Reader — Seek и Len здесь не нужны,
// привязка к конкретному типу без причины усложняет тестирование и переиспользование.
func parseSubmitSMBody(r *bytes.Reader) (*SubmitSMBody, error) {
	var body SubmitSMBody
	var err error

	body.ServiceType, err = readCString(r)
	if err != nil {
		return nil, err
	}

	// TODO: для uint8 binary.Read избыточен — он использует рефлексию чтобы прочитать 1 байт,
	// а BigEndian для однобайтового значения вообще ничего не значит.
	// Заменить на r.ReadByte() — короче, быстрее, выразительнее.
	err = binary.Read(r, binary.BigEndian, &body.SourceAddrTON)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &body.SourceAddrNPI)
	if err != nil {
		return nil, err
	}

	body.SourceAddr, err = readCString(r)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &body.DestAddrTON)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &body.DestAddrNPI)
	if err != nil {
		return nil, err
	}

	body.DestinationAddr, err = readCString(r)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &body.ESMClass)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &body.ProtocolID)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &body.PriorityFlag)
	if err != nil {
		return nil, err
	}

	body.ScheduleDeliveryTime, err = readCString(r)
	if err != nil {
		return nil, err
	}

	body.ValidityPeriod, err = readCString(r)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &body.RegisteredDelivery)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &body.ReplaceIfPresentFlag)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &body.DataCoding)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &body.SMDefaultMsgID)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &body.SMLength)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, body.SMLength)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	body.ShortMessage = buf

	return &body, nil
}

func parseSubmitSM(r *bytes.Reader) (*SubmitSM, error) {
	var result SubmitSM
	var header *PDUHeader
	var body *SubmitSMBody
	var tlvs []TLV
	var err error

	header, err = ParsePDUHeader(r)
	if err != nil {
		return nil, err
	}

	if header.CommandID != CommandSubmitSM {
		return nil, fmt.Errorf("failed to parse SubmitSM, wrong CommandID: %d", header.CommandID)
	}

	result.PDUHeader = *header

	body, err = parseSubmitSMBody(r)
	if err != nil {
		return nil, err
	}

	result.SubmitSMBody = *body

	tlvs, err = parseTLVs(r)
	if err != nil {
		return nil, err
	}

	result.TLVs = tlvs

	return &result, nil

}
