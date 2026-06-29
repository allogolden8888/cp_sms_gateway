//go:build ignore

// Запуск: go run simulator.go -addr localhost:8080 -from 79001234567 -to 79007654321 -msg "Hello"
package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"flag"
)

func cstr(s string) []byte {
	return append([]byte(s), 0x00)
}

func pduHeader(commandID, commandStatus, seq uint32, bodyLen int) []byte {
	h := make([]byte, 16)
	binary.BigEndian.PutUint32(h[0:], uint32(16+bodyLen))
	binary.BigEndian.PutUint32(h[4:], commandID)
	binary.BigEndian.PutUint32(h[8:], commandStatus)
	binary.BigEndian.PutUint32(h[12:], seq)
	return h
}

func buildBindTransceiver(systemID, password string, seq uint32) []byte {
	var body []byte
	body = append(body, cstr(systemID)...)
	body = append(body, cstr(password)...)
	body = append(body, cstr("")...)  // system_type
	body = append(body, 0x34)         // interface_version
	body = append(body, 0x00)         // addr_ton
	body = append(body, 0x00)         // addr_npi
	body = append(body, cstr("")...) // address_range

	const cmdBindTransceiver = 0x00000009
	return append(pduHeader(cmdBindTransceiver, 0, seq, len(body)), body...)
}

func buildSubmitSM(from, to, message string, seq uint32) []byte {
	msg := []byte(message)

	var body []byte
	body = append(body, cstr("")...)    // service_type
	body = append(body, 0x01)           // source_addr_ton (international)
	body = append(body, 0x01)           // source_addr_npi (ISDN)
	body = append(body, cstr(from)...)
	body = append(body, 0x01)           // dest_addr_ton
	body = append(body, 0x01)           // dest_addr_npi
	body = append(body, cstr(to)...)
	body = append(body, 0x00)           // esm_class
	body = append(body, 0x00)           // protocol_id
	body = append(body, 0x00)           // priority_flag
	body = append(body, cstr("")...)    // schedule_delivery_time
	body = append(body, cstr("")...)    // validity_period
	body = append(body, 0x00)           // registered_delivery
	body = append(body, 0x00)           // replace_if_present_flag
	body = append(body, 0x00)           // data_coding
	body = append(body, 0x00)           // sm_default_msg_id
	body = append(body, byte(len(msg))) // sm_length
	body = append(body, msg...)

	const cmdSubmitSM = 0x00000004
	return append(pduHeader(cmdSubmitSM, 0, seq, len(body)), body...)
}

func readPDU(conn net.Conn) ([]byte, error) {
	hdr := make([]byte, 16)
	if _, err := readFull(conn, hdr); err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	length := binary.BigEndian.Uint32(hdr[0:4])
	if length < 16 {
		return nil, fmt.Errorf("invalid PDU length: %d", length)
	}
	rest := make([]byte, length-16)
	if _, err := readFull(conn, rest); err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	return append(hdr, rest...), nil
}

func readFull(conn net.Conn, buf []byte) (int, error) {
	total := 0
	for total < len(buf) {
		n, err := conn.Read(buf[total:])
		total += n
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

func parsePDUHeader(b []byte) (commandID, commandStatus, seq uint32) {
	commandID = binary.BigEndian.Uint32(b[4:8])
	commandStatus = binary.BigEndian.Uint32(b[8:12])
	seq = binary.BigEndian.Uint32(b[12:16])
	return
}

func commandName(id uint32) string {
	switch id {
	case 0x80000009:
		return "bind_transceiver_resp"
	case 0x80000004:
		return "submit_sm_resp"
	case 0x80000000:
		return "generic_nack"
	default:
		return fmt.Sprintf("0x%08X", id)
	}
}

func main() {
	addr := flag.String("addr", "localhost:8080", "SMPP server address")
	systemID := flag.String("id", "test", "system_id")
	password := flag.String("pw", "test", "password")
	from := flag.String("from", "79001234567", "source address")
	to := flag.String("to", "79007654321", "destination address")
	msg := flag.String("msg", "Hello SMPP", "short message text")
	flag.Parse()

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "connect:", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("connected to", *addr)

	// --- bind ---
	bindPDU := buildBindTransceiver(*systemID, *password, 1)
	if _, err := conn.Write(bindPDU); err != nil {
		fmt.Fprintln(os.Stderr, "send bind:", err)
		os.Exit(1)
	}
	fmt.Printf("sent bind_transceiver (system_id=%q)\n", *systemID)

	resp, err := readPDU(conn)
	if err != nil {
		fmt.Fprintln(os.Stderr, "read bind_resp:", err)
		os.Exit(1)
	}
	cmdID, status, _ := parsePDUHeader(resp)
	fmt.Printf("recv %s status=0x%08X\n", commandName(cmdID), status)
	if status != 0 {
		fmt.Fprintln(os.Stderr, "bind failed")
		os.Exit(1)
	}

	// --- submit_sm ---
	submitPDU := buildSubmitSM(*from, *to, *msg, 2)
	if _, err := conn.Write(submitPDU); err != nil {
		fmt.Fprintln(os.Stderr, "send submit_sm:", err)
		os.Exit(1)
	}
	fmt.Printf("sent submit_sm from=%s to=%s msg=%q\n", *from, *to, *msg)

	resp, err = readPDU(conn)
	if err != nil {
		fmt.Fprintln(os.Stderr, "read submit_sm_resp:", err)
		os.Exit(1)
	}
	cmdID, status, _ = parsePDUHeader(resp)
	fmt.Printf("recv %s status=0x%08X\n", commandName(cmdID), status)
	if status != 0 {
		fmt.Fprintln(os.Stderr, "submit failed")
		os.Exit(1)
	}

	// message_id из resp[16:] до нулевого байта
	msgID := ""
	for i := 16; i < len(resp); i++ {
		if resp[i] == 0 {
			break
		}
		msgID += string(resp[i])
	}
	fmt.Println("message_id:", msgID)
}