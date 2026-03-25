package smpp

import "testing"

func TestParseDLR(t *testing.T) {
	result, err := ParseDLR("id:abc123 sub:001 dlvrd:001 submit date:2601011200 done date:2601011201 stat:DELIVRD err:000 Text:Test")
	if err != nil {
		t.Errorf("ParseDLR error")
	}
	expected := &DLR{
		MessageID:  "abc123",
		Submitted:  "001",
		Delivered:  "001",
		SubmitDate: "2601011200",
		DoneDate:   "2601011201",
		Status:     "DELIVRD",
		ErrorCode:  "000",
		Text:       "Test",
	}

	if result.MessageID != expected.MessageID || result.Status != expected.Status || result.ErrorCode != expected.ErrorCode {
		t.Errorf("got %v, expected %v", result, expected)
	}
}
