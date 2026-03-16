package smpp

import (
	"fmt"
	"strings"
)

const gsm7Chars = "@£$¥èéùìòÇØøÅåΔ_ΦΓΛΩΠΨΣΘΞ\x1bÆæßÉ !\"#¤%&'()*+,-./0123456789:;<=>?¡ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÑÜ§¿abcdefghijklmnopqrstuvwxyzäöñüà\r\n"

func validateEncoding(text, encoding string) error {
	for _, r := range text {
		switch encoding {
		case "latin1":
			if r > 0xFF {
				return fmt.Errorf("character '%c' (U+%04X) not supported in latin1", r, r)
			}
		case "gsm7":

			if !strings.ContainsRune(gsm7Chars, r) {
				return fmt.Errorf("character '%c' (U+%04X) not supported in gsm7", r, r)
			}

		case "ucs2":
			// нет ограничений
		}
	}
	return nil
}
