package smpp

import (
	"bytes"
	"cp_sms_gateway/smpp/pdu"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
)

func StartServer(port string) {

	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Printf("Connection error occured: %v\n", err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Connection error occured: %v\n", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var header *pdu.PDUHeader

	for {
		headerBuf := make([]byte, 16)

		_, err := io.ReadFull(conn, headerBuf)
		if err != nil {
			break
		}

		header, err = pdu.ParsePDUHeader(bytes.NewReader(headerBuf))
		if err != nil {
			break
		}

		bodyBuf := make([]byte, header.Length-16)

		_, err = io.ReadFull(conn, bodyBuf)
		if err != nil {
			break
		}

		pdubytes := append(headerBuf, bodyBuf...)
		parsedPDU, err := pdu.ParsePDU(pdubytes)
		fmt.Printf("Received PDU: %v", reflect.TypeOf(parsedPDU))

		switch v := parsedPDU.(type) {
		case *pdu.BindReceiver:
			var b pdu.BindResp

			b.CommandID = pdu.CommandBindReceiverResp
			b.CommandStatus = 0
			b.SequenceNumber = v.SequenceNumber
			b.SystemID = v.SystemID

			conn.Write(pdu.SerializeBindResp(&b))
			fmt.Println("Sent Bind Resp:", b.SequenceNumber)
		case *pdu.BindTransmitter:
			var b pdu.BindResp

			b.CommandID = pdu.CommandBindTransmitterResp
			b.CommandStatus = 0
			b.SequenceNumber = v.SequenceNumber
			b.SystemID = v.SystemID

			conn.Write(pdu.SerializeBindResp(&b))
			fmt.Println("Sent Bind Resp:", b.SequenceNumber)
		case *pdu.BindTransceiver:
			var b pdu.BindResp

			b.CommandID = pdu.CommandBindTransceiverResp
			b.CommandStatus = 0
			b.SequenceNumber = v.SequenceNumber
			b.SystemID = v.SystemID

			conn.Write(pdu.SerializeBindResp(&b))
			fmt.Printf("\nSent Bind Resp: %d\n", b.SequenceNumber)
		case *pdu.SubmitSM:
			var s pdu.SubmitSMResp

			s.CommandID = pdu.CommandBindTransceiverResp
			s.CommandStatus = 0
			s.SequenceNumber = v.SequenceNumber

			s.MessageID = "1"

			conn.Write(pdu.SerializeSubmitSMResp(&s))
			fmt.Printf("\nSent SubmitSMResp: %v\n", s.MessageID)
		}
	}

}
