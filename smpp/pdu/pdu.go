package pdu

import "bytes"

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
