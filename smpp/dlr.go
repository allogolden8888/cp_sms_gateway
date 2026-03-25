package smpp

import (
	"fmt"
	"regexp"
)

type DLR struct {
	MessageID  string
	Submitted  string
	Delivered  string
	SubmitDate string
	DoneDate   string
	Status     string
	ErrorCode  string
	Text       string
}

func ParseDLR(text string) (*DLR, error) {
	re := regexp.MustCompile(`id:(\w+) sub:(\w+) dlvrd:(\w+) submit date:(\w+) done date:(\w+) stat:(\w+) err:(\w+) Text:(.+)`)
	result := re.FindStringSubmatch(text)

	if result == nil {
		return nil, fmt.Errorf("failed to parse DLR: %s", text)
	}

	return &DLR{
		MessageID:  result[1],
		Submitted:  result[2],
		Delivered:  result[3],
		SubmitDate: result[4],
		DoneDate:   result[5],
		Status:     result[6],
		ErrorCode:  result[7],
		Text:       result[8],
	}, nil

}
