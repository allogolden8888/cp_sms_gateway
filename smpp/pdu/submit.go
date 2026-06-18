package pdu

import (
	"bytes"
	"encoding/binary"
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

func parseSubmitSMBody(r *bytes.Reader) (*SubmitSMBody, error) {
	var body SubmitSMBody
	var err error

	body.ServiceType, err = readCString(r)
	if err != nil {
		return nil, err
	}

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
