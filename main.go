package main

import (
	"flag"
	"fmt"
	"log/slog"
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
	if err != nil {
		slog.Error("config loading failed", "error", err)
		os.Exit(1)
	}

	var timeout time.Duration

	if config.Validity != "" {
		timeout, err = cpsmpp.ParseValidity(config.Validity)
		if err != nil {
			slog.Error("config validity period loading failed", "error", err)
			os.Exit(1)
		}
	} else {
		timeout = time.Duration(24 * time.Hour)
	}

	if config.Encoding != "gsm7" && config.Encoding != "ucs2" && config.Encoding != "latin1" {
		slog.Error("config encoding loading failed", "selected encoding", config.Encoding)
		os.Exit(1)
	}

	dlrTracker := cpsmpp.NewDLRTracker()

	client, status := cpsmpp.Connect(config.BindType, fmt.Sprintf("%s:%d", config.Host, config.Port),
		config.Username, config.Password, func(p pdu.Body) {
			if p.Header().ID == 0x00000005 {
				text := p.Fields()[pdufield.ShortMessage].String()
				dlr, err := cpsmpp.ParseDLR(text)
				if err != nil {
					slog.Error("dlr parsing failed", "error", err)
				} else {
					slog.Debug("dlr received", "message_id", dlr.MessageID, "status", dlr.Status, "error", dlr.ErrorCode, "done_date", dlr.DoneDate)
					dlrTracker.Receive(dlr)
				}

			}

		})

	if status != nil {
		slog.Error("connection failed", "error", status)
		os.Exit(1)
	}

	sm, err := cpsmpp.SendMessage(config.From, config.To, config.Message, config.Encoding, config.Validity, client, config.Register, config.Priority)
	if err != nil {
		slog.Error("submit failed", "error", err)
		os.Exit(1)
	}

	for _, part := range sm {
		dlrTracker.Expect(part.RespID())
	}

	slog.Info("message sent", "host", config.Host, "port", config.Port, "source_addr", config.From, "destination",
		config.To, "system_id", config.Username, "password", config.Password, "message_text", config.Message)

	for i, part := range sm {
		slog.Debug("displaying submit_sm_resps", "part", i+1, "message_id", part.RespID())
	}
	slog.Info("dlr waiting ", "timeout", timeout.String())
	for i := 0; i < len(sm); i++ {
		select {
		case id := <-dlrTracker.Done():
			slog.Debug("dlr received", "message_id", id)
		case <-time.After(timeout):
			slog.Info("dlr timeouted")
			return
		}
	}

}
