package smpp

import "testing"

func TestEncoding(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		encoding string
		wantErr  bool
	}{
		{
			name:     "valid gsm7",
			input:    "hello",
			encoding: "gsm7",
		},
		{
			name:     "invalid with cyrillic gsm7",
			input:    "helloк",
			encoding: "gsm7",
			wantErr:  true,
		},
		{
			name:     "emptystring gsm7",
			input:    "",
			encoding: "gsm7",
		},
		{
			name:     "valid latin1",
			input:    "café",
			encoding: "latin1",
		},
		{
			name:     "invalid latin1",
			input:    "привет",
			encoding: "latin1",
			wantErr:  true,
		},
		{
			name:     "valid ucs2",
			input:    "hello",
			encoding: "ucs2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEncoding(tt.input, tt.encoding)
			if !tt.wantErr {
				if err != nil {
					t.Errorf("%v test failed, unexpected err: %v", tt.name, err)
				}
			} else {
				if err == nil {
					t.Errorf("%v test failed, expected error, got nil", tt.name)
				}
			}
		})
	}
}
