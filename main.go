package main

import (
	"flag"
	"fmt"
	"os"

	cpsmpp "cp_sms_gateway/smpp"
)

func main() {

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

	if *encoding != "gsm7" && *encoding != "ucs2" && *encoding != "latin1" {
		fmt.Printf("Wrong data_coding selected: %s. Allowed encodings: gsm7, ucs2, latin1", *encoding)
		os.Exit(1)
	}

	client, status := cpsmpp.Connect(*bindType, fmt.Sprintf("%s:%d", *host, *port), *username, *password)

	if status != nil {
		fmt.Println(status.Error())
	} else {
		defer client.Close()
		sm, err := cpsmpp.SendMessage(*from, *to, *message, *encoding, *validity, client, *register, *priority)
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
	}

}
