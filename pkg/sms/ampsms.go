package sms

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
)

const (
	DefaultSender = "amp"
)

type Sms struct {
	accountID string
	sender    string
	apiKey    string
}

func NewSms(accountID string, apiKey string, sender string) *Sms {
	return &Sms{
		accountID: accountID,
		apiKey:    apiKey,
		sender:    sender,
	}
}

//SendSms send a sms
func (s *Sms) SendSms(to string, message string) error {
	url := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", s.accountID)
	body := encodeBody(s.sender, to, message)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	req.SetBasicAuth(s.accountID, s.apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func encodeBody(from string, to string, message string) string {
	v := url.Values{}
	v.Set("From", from)
	v.Add("To", to)
	v.Add("Body", message)
	return v.Encode()
}
