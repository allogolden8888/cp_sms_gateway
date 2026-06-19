package pdu

import "bytes"

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
