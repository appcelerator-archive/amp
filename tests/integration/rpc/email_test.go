package tests

import (
	"bytes"
	"log"
	"net"
	"net/mail"
	"testing"
	"time"
	"fmt"

	"github.com/appcelerator/amp/pkg/ampmail"
	"github.com/mhale/smtpd"
)

type emailMessage struct {
	from string
	to []string
	subject string
	data  []byte
}

var emailReceived *emailMessage
var from="amp@axway.com"
var to="user@axway.com"
var accountName="myAccount"
var token="1234567890"


func TestAmpMail(t *testing.T) {
	ampmail.UpdateAmpMailConfig("localhost", "2525", from, "")
	emailReceived=nil
	ampmail.SendAccountVerificationEmail(to, accountName, token)
	waitTestForEmail(t, "AMP account activation", "AMP account activation")
	emailReceived=nil
	ampmail.SendAccountResetPasswordEmail(to, accountName, token)
	waitTestForEmail(t, "AccountResetPassword", "AMP account password reset")
	emailReceived=nil
	ampmail.SendAccountPasswordConfirmationEmail(to, accountName)
	waitTestForEmail(t, "AccountPasswordConfirmation", "AMP account password confirmation")

}

func initMailServer() {
	go func() {
		fmt.Printf("server mail started\n")
		smtpd.ListenAndServe("127.0.0.1:2525", mailHandler, "MailServerTest", "")
	}()
}

func mailHandler(origin net.Addr, from string, to []string, data []byte) {
	msg, _ := mail.ReadMessage(bytes.NewReader(data))
	subject := msg.Header.Get("Subject")
	emailReceived = &emailMessage{
		from : from,
		to: to,
		subject: subject,
		data: data,
	}
	log.Printf("Received mail from %s for %s with subject %s", from, to[0], subject)
}

func waitTestForEmail(t *testing.T, template string, subject string) {
	t0:=time.Now()
	for emailReceived==nil {
		time.Sleep(3*time.Second)
		if time.Now().Sub(t0).Seconds()>3 {
			t.Fatalf("Email %s not received", template)
		}
	}
	if emailReceived.subject!=subject {
		t.Fatalf("Email %s bad subject, should be %s, got %s", template, subject, emailReceived.subject)
	}
	if emailReceived.from!=from {
		t.Fatalf("Email %s bad sender, should be %s, got %s", template, from, emailReceived.from)
	}
	if len(emailReceived.to)!=1 {
		t.Fatalf("Email %s bad receiver number, should be 1, got %d", template, len(emailReceived.to))
	}
	if emailReceived.to[0]!=to {
		t.Fatalf("Email %s bad receiver, should be %s, got %s", template, to, emailReceived.to[0])
	}
}
