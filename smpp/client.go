package smpp

import (
	"fmt"
	"os"

	gosmpp "github.com/fiorix/go-smpp/smpp"
)

type Client interface {
	Bind() <-chan gosmpp.ConnStatus
	Close() error
	Submit(sm *gosmpp.ShortMessage) (*gosmpp.ShortMessage, error)
	SubmitLongMsg(sm *gosmpp.ShortMessage) ([]gosmpp.ShortMessage, error)
}

func Connect(bindType, addr, user, passwd string, handler gosmpp.HandlerFunc) (Client, error) {
	if bindType != "tx" && bindType != "trx" && bindType != "rx" {
		return nil, fmt.Errorf("wrong bind type: %s", bindType)
	} else {
		var client Client
		switch bindType {
		case "tx":
			client = &gosmpp.Transmitter{
				Addr:   addr,
				User:   user,
				Passwd: passwd,
			}
		case "trx":
			client = &gosmpp.Transceiver{
				Addr:    addr,
				User:    user,
				Passwd:  passwd,
				Handler: handler,
			}

		case "rx":
			fmt.Println("RX connection type is under development.")
			os.Exit(1)
		}
		conn := client.Bind()

		status := <-conn

		return client, status.Error()

	}

}
