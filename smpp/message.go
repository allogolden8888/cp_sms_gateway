package smpp

import (
	"fmt"
	"time"

	gosmpp "github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"
)

func parseValidity(value string) (time.Duration, error) {
	d, err := time.ParseDuration(value)
	if err == nil {
		return d, nil
	}
	t, err := time.Parse("2006-01-02T15:04:05", value)
	if err != nil {
		return 0, fmt.Errorf("invalid validity: %s", value)
	}
	return time.Until(t), nil
}

func SendMessage(src, dst, text, encoding, validity string, client Client, register, priority int) ([]*gosmpp.ShortMessage, error) {
	shortMessage := gosmpp.ShortMessage{
		Src:            src,
		Dst:            dst,
		ServiceType:    "",
		SourceAddrTON:  0,
		SourceAddrNPI:  0,
		DestAddrTON:    1,
		DestAddrNPI:    1,
		ESMClass:       0,
		ProtocolID:     0,
		SMDefaultMsgID: 0,
		Register:       pdufield.DeliverySetting(register),
		PriorityFlag:   uint8(priority),
	}

	if validity != "" {
		d, err := parseValidity(validity)
		if err != nil {
			return nil, err
		}
		shortMessage.Validity = d
	}

	var partMaxLen int

	encodeCheck := validateEncoding(text, encoding)
	if encodeCheck != nil {
		return nil, encodeCheck
	}

	switch encoding {
	case "gsm7":
		shortMessage.Text = pdutext.GSM7(text)
		partMaxLen = 160
	case "ucs2":
		shortMessage.Text = pdutext.UCS2(text)
		partMaxLen = 70
	case "latin1":
		shortMessage.Text = pdutext.Latin1(text)
		partMaxLen = 140
	}

	var textLen int
	if encoding == "ucs2" {
		textLen = len([]rune(text))
	} else {
		textLen = len(text)
	}

	if textLen > partMaxLen {
		parts, err := client.SubmitLongMsg(&shortMessage)
		result := make([]*gosmpp.ShortMessage, len(parts))
		for i := range parts {
			result[i] = &parts[i]
		}
		return result, err
	} else {
		sm, err := client.Submit(&shortMessage)
		return []*gosmpp.ShortMessage{sm}, err
	}

}
