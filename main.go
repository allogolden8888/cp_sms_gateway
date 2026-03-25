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
	configPath := flag.String("config", "config.json", "path to config")

	flag.Parse()

	config, err := cpsmpp.LoadConfig(*configPath)

	var timeout time.Duration

	if config.Validity != "" {
		timeout, err = cpsmpp.ParseValidity(config.Validity)
		if err != nil {
			fmt.Println("Invalid Validity:", err)
			os.Exit(1)
		}
	} else {
		timeout = time.Duration(24 * time.Hour)
	}

	if config.Encoding != "gsm7" && config.Encoding != "ucs2" && config.Encoding != "latin1" {
		fmt.Printf("Wrong data_coding selected: %s. Allowed encodings: gsm7, ucs2, latin1", config.Encoding)
		os.Exit(1)
	}

	done := make(chan string, 10)

	respIDs := make(map[string]bool)

	client, status := cpsmpp.Connect(config.BindType, fmt.Sprintf("%s:%d", config.Host, config.Port), config.Username, config.Password, func(p pdu.Body) {
		if p.Header().ID == 0x00000005 {
			text := p.Fields()[pdufield.ShortMessage].String()
			dlr, err := cpsmpp.ParseDLR(text)
			if err != nil {
				fmt.Println("DLR parse error:", err)
			} else {
				fmt.Printf("DLR received:\n  message_id: %s\n  status: %s\n  error: %s\n  done: %s\n",
					dlr.MessageID, dlr.Status, dlr.ErrorCode, dlr.DoneDate)
			}
			_, ok := respIDs[dlr.MessageID]
			if ok {
				fmt.Println("Expected DLR ID")
				done <- dlr.MessageID
			}

		}

	})

	if status != nil {
		fmt.Println(status.Error())
		os.Exit(1)
	}

	sm, err := cpsmpp.SendMessage(config.From, config.To, config.Message, config.Encoding, config.Validity, client, config.Register, config.Priority)
	if err != nil {
		fmt.Println("Submit error:", err)
		os.Exit(1)
	}

	for _, part := range sm {
		respIDs[part.RespID()] = true
	}

	fmt.Printf(`Sending message via %s:%d;
				Source addr: %s;
				Destination addr: %s;
				System ID: %s;
				Password: %s;
				Message Text: %s
				`, config.Host, config.Port, config.From, config.To, config.Username, config.Password, config.Message)

	for i, part := range sm {
		fmt.Printf("Part %d: message_id=%s\n", i+1, part.RespID())
	}
	fmt.Println("Waiting for DLR, timeout:", timeout)
	for i := 0; i < len(sm); i++ {
		select {
		case id := <-done:
			fmt.Printf("DLR для части message_id=%s\n", id)
		case <-time.After(timeout):
			fmt.Println("DLR Timeout")
			return
		}
	}

}
