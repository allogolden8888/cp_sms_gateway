package pdu

import (
	"bytes"
	"fmt"
)

// TODO: заменить *bytes.Reader на io.ByteReader — функция использует только r.ReadByte(),
// этот метод входит в интерфейс io.ByteReader. Импорт "bytes" заменить на "io".
func readCString(r *bytes.Reader) (string, error) {
	var result []byte

	for {
		x, err := r.ReadByte()
		if err != nil {
			return "", err
		}
		if x == 0x00 {
			break
		}

		result = append(result, x)

	}

	return string(result), nil
}

func ParsePDU(b []byte) (interface{}, error) {
	buf := bytes.NewReader(b)

	header, err := ParsePDUHeader(buf)
	if err != nil {
		return nil, err
	}

	bufCopy := bytes.NewReader(b)

	switch header.CommandID {
	case CommandGenericNACK:
		result, err := ParseGenerickNACK(bufCopy)
		if err != nil {
			return nil, err
		}
		return result, nil
	case CommandBindReceiver:
		result, err := ParseBindReceiver(bufCopy)
		if err != nil {
			return nil, err
		}
		return result, nil
	case CommandBindReceiverResp:
		result, err := ParseBindResp(bufCopy, CommandBindReceiverResp)
		if err != nil {
			return nil, err
		}
		return result, nil
	case CommandBindTransmitter:
		result, err := ParseBindTransmitter(bufCopy)
		if err != nil {
			return nil, err
		}
		return result, nil
	case CommandBindTransmitterResp:
		result, err := ParseBindResp(bufCopy, CommandBindTransmitterResp)
		if err != nil {
			return nil, err
		}
		return result, nil
	case CommandBindTransceiver:
		result, err := ParseBindTransceiver(bufCopy)
		if err != nil {
			return nil, err
		}
		return result, nil
	case CommandBindTransceiverResp:
		result, err := ParseBindResp(bufCopy, CommandBindTransceiverResp)
		if err != nil {
			return nil, err
		}
		return result, nil
	case CommandSubmitSM:
		result, err := parseSubmitSM(bufCopy)
		if err != nil {
			return nil, err
		}
		return result, nil
	case CommandSubmitSMResp:
		result, err := ParseSubmitSMResp(bufCopy)
		if err != nil {
			return nil, err
		}
		return result, nil
	case CommandDeliverSM:
		result, err := ParseDeliverSM(bufCopy)
		if err != nil {
			return nil, err
		}
		return result, nil
	case CommandDeliverSMResp:
		result, err := ParseDeliverSMResp(bufCopy)
		if err != nil {
			return nil, err
		}
		return result, nil
	case CommandUnbind:
		result, err := ParseUnbind(bufCopy)
		if err != nil {
			return nil, err
		}
		return result, nil
	case CommandUnbindResp:
		result, err := ParseUnbindResp(bufCopy)
		if err != nil {
			return nil, err
		}
		return result, nil
	case CommandEnquireLink:
		result, err := ParseEnquireLink(bufCopy)
		if err != nil {
			return nil, err
		}
		return result, nil
	case CommandEnquireLinkResp:
		result, err := ParseEnquireLinkResp(bufCopy)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	return nil, fmt.Errorf("Unsupported PDU with CommandID: %v", header.CommandID)

}
