package pdu

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type DeliverSMBody struct {
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

type DeliverSM struct {
	PDUHeader
	DeliverSMBody
	TLVs []TLV
}

type DeliverSMResp struct {
	PDUHeader
	MessageID string
}

const (
	CommandDeliverSM     uint32 = 0x00000005
	CommandDeliverSMResp uint32 = 0x80000005
)

func ParseDeliverSMResp(r *bytes.Reader) (*DeliverSMResp, error) {
	var messageID string
	header, err := ParsePDUHeader(r)
	if err != nil {
		return nil, err
	}

	if header.CommandID != CommandDeliverSMResp {
		return nil, fmt.Errorf("failed to parse DeliverSMResp, wrong CommandID: %d", header.CommandID)
	}

	messageID, err = readCString(r)
	if err != nil {
		return nil, err
	}

	result := DeliverSMResp{
		PDUHeader: *header,
		MessageID: messageID,
	}

	return &result, nil
}

func ParseDeliverSM(r *bytes.Reader) (*DeliverSM, error) {
	// Total copy of parseSubmitSM, didnt reuse cause of readability and optimization
	var result DeliverSM
	var header *PDUHeader
	var body *DeliverSMBody
	var tlvs []TLV
	var err error

	header, err = ParsePDUHeader(r)
	if err != nil {
		return nil, err
	}

	if header.CommandID != CommandDeliverSM {
		return nil, fmt.Errorf("failed to parse DeliverSM, wrong CommandID: %d", header.CommandID)
	}

	result.PDUHeader = *header

	body, err = parseDeliverSMBody(r)
	if err != nil {
		return nil, err
	}

	result.DeliverSMBody = *body

	tlvs, err = parseTLVs(r)
	if err != nil {
		return nil, err
	}

	result.TLVs = tlvs

	return &result, nil
}

func parseDeliverSMBody(r *bytes.Reader) (*DeliverSMBody, error) {
	// Total copy of parseSubmitSMBody, didnt reuse cause of readability and optimization
	var err error
	var body DeliverSMBody

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
