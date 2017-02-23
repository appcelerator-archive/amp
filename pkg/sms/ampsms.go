package sms

import (
	"bytes"
	"fmt"
	"github.com/appcelerator/amp/pkg/config"
	"net/http"
	"net/url"
	"os"
)

//SendSms send a sms
func SendSms(to string, message string) error {
	config := amp.GetConfig()
	url := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", config.SmsAccountID)
	body := encodeBody(config.SmsSender, to, message)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	req.SetBasicAuth(config.SmsAccountID, config.SmsKey)
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
