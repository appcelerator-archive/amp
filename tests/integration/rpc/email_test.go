package tests

import (
	"bytes"
	"log"
	"net"
	netmail "net/mail"

	"github.com/mhale/smtpd"
)

type emailMessage struct {
	from    string
	to      []string
	subject string
	data    []byte
}

var emailReceived *emailMessage

func initMailServer() {
	go func() {
		smtpd.ListenAndServe("127.0.0.1:2525", mailHandler, "MailServerTest", "")
	}()
}

func mailHandler(origin net.Addr, from string, to []string, data []byte) {
	msg, _ := netmail.ReadMessage(bytes.NewReader(data))
	subject := msg.Header.Get("Subject")
	emailReceived = &emailMessage{
		from:    from,
		to:      to,
		subject: subject,
		data:    data,
	}
	log.Printf("Received mail from %s for %s with subject %s", from, to[0], subject)
}
