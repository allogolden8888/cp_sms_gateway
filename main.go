package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	cpsmpp "cp_sms_gateway/smpp"

	"github.com/fiorix/go-smpp/smpp/pdu"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
)

func main() {

	cpsmpp.ParseDLR("id:1 sub:001 dlvrd:001 submit date:2603201425 done date:2603201425 stat:DELIVRD err:000 Text:Hello world!")

	host := flag.String("host", "localhost", "gosmpp Server Address")
	port := flag.Int("port", 2775, "gosmpp Server Port")
	username := flag.String("username", "admin", "gosmpp system_id")
	password := flag.String("password", "admin", "gosmpp Password")
	from := flag.String("from", "Test", "Source address")
	to := flag.String("to", "998901331835", "Destination address")
	message := flag.String("message", "Hello world!", "Message text")
	bindType := flag.String("bind_type", "tx", "Bind type")
	encoding := flag.String("encoding", "gsm7", "data_coding")
	validity := flag.String("validity_period", "24h", "validity_period")
	register := flag.Int("registered_delivery", 1, "registered_delivery")
	priority := flag.Int("priority_flag", 1, "priority_flag")

	flag.Parse()

	var timeout time.Duration
	var err error
	if *validity != "" {
		timeout, err = cpsmpp.ParseValidity(*validity)
		if err != nil {
			fmt.Println("Invalid Validity:", err)
			os.Exit(1)
		}
	} else {
		timeout = time.Duration(24 * time.Hour)
	}

	if *encoding != "gsm7" && *encoding != "ucs2" && *encoding != "latin1" {
		fmt.Printf("Wrong data_coding selected: %s. Allowed encodings: gsm7, ucs2, latin1", *encoding)
		os.Exit(1)
	}

	done := make(chan string, 10)

	client, status := cpsmpp.Connect(*bindType, fmt.Sprintf("%s:%d", *host, *port), *username, *password, func(p pdu.Body) {
		if p.Header().ID == 0x00000005 {
			text := p.Fields()[pdufield.ShortMessage].String()
			dlr, err := cpsmpp.ParseDLR(text)
			if err != nil {
				fmt.Println("DLR parse error:", err)
			} else {
				fmt.Printf("DLR received:\n  message_id: %s\n  status: %s\n  error: %s\n  done: %s\n",
					dlr.MessageID, dlr.Status, dlr.ErrorCode, dlr.DoneDate)
			}
			done <- dlr.MessageID
		}

	})

	if status != nil {
		fmt.Println(status.Error())
		os.Exit(1)
	}

	defer client.Close()
	sm, err, partsCount := cpsmpp.SendMessage(*from, *to, *message, *encoding, *validity, client, *register, *priority)
	if err != nil {
		fmt.Println("Submit error:", err)
		os.Exit(1)
	}

	fmt.Printf(`Sending message via %s:%d;
				Source addr: %s;
				Destination addr: %s;
				System ID: %s;
				Password: %s;
				Message Text: %s
				`, *host, *port, *from, *to, *username, *password, *message)

	for i, part := range sm {
		fmt.Printf("Part %d: message_id=%s\n", i+1, part.RespID())
	}
	fmt.Println("Waiting for DLR, timeout:", timeout)
	for i := 0; i < partsCount; i++ {
		select {
		case id := <-done:
			fmt.Printf("DLR для части message_id=%s\n", id)
		case <-time.After(timeout):
			fmt.Println("DLR Timeout")
			return
		}
	}

}
