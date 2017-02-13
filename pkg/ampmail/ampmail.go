package ampmail

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"log"
	"net/smtp"
	"strings"

	conf "github.com/appcelerator/amp/pkg/config"
)

var emailTemplateMap map[string]*emailTemplate
var config *conf.Configuration

type emailTemplate struct {
	isHtml  bool
	subject string
	body    string
}

type pCipher struct {
	key    []byte
	nonce  []byte
	block  cipher.Block
	buffer []byte
}

func init() {
	config = conf.GetRegularConfig(false)
	emailTemplateMap = make(map[string]*emailTemplate)
	AddEmailTemplate("AccountVerification", "AMP account activation", true, accountVerificationBody)
	AddEmailTemplate("AccountResetPassword", "AMP account password reset", true, accountResetPasswordEmailBody)
	AddEmailTemplate("AccountPasswordConfirmation", "AMP account password confirmation", true, accountPasswordConfirmationEmailBody)
}

// SendAccountVerificationEmail send a AccountVerification email template
func SendAccountVerificationEmail(to string, accountName string, token string) error {
	//config := conf.GetRegularConfig(false)
	variables := map[string]string{
		"accountName": accountName,
		"token":       token,
		"ampAddress":  config.AmpAddress,
	}
	if err := SendTemplateEmail(to, "AccountVerification", variables); err != nil {
		return err
	}
	return nil
}

// SendAccountResetPasswordEmail send a AccountResetPassword email template
func SendAccountResetPasswordEmail(to string, accountName string, token string) error {
	//config := conf.GetRegularConfig(false)
	variables := map[string]string{
		"accountName": accountName,
		"token":       token,
		"ampAddress":  config.AmpAddress,
	}
	if err := SendTemplateEmail(to, "AccountResetPassword", variables); err != nil {
		return err
	}
	return nil
}

// SendAccountResetPasswordEmail send a AccountResetPassword email template
func SendAccountPasswordConfirmationEmail(to string, accountName string) error {
	variables := map[string]string{
		"accountName": accountName,
	}
	if err := SendTemplateEmail(to, "AccountPasswordConfirmation", variables); err != nil {
		return err
	}
	return nil
}

// SendTemplateEmail send a tempalte email
func SendTemplateEmail(to string, templateEmailName string, variableMap map[string]string) error {
	email, err := getEmailTemplate(templateEmailName)
	if err != nil {
		return err
	}
	email.setVariables(variableMap)
	if err := SendMail(to, email.subject, email.isHtml, email.body); err != nil {
		return err
	}
	return nil
}

// SendMail send an eamail to "to" with subject and body, use configuration
func SendMail(to string, subject string, isHtml bool, body string) error {
	//config := conf.GetRegularConfig(false)
	from := config.EmailSender
	servername := config.EmailServerAddress
	serverport := config.EmailServerPort
	pass := config.EmailPwd
	if pass == "" {
		p, err := newPCipher(config)
		if err != nil {
			return err
		}
		datah, _ := hex.DecodeString("94aacd567dff5c9bc02860")
		data, errd := p.decrypt(datah)
		if errd != nil {
			return errd
		}
		pass = string(data)
	}
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n", from, to, subject, body)
	if isHtml {
		msg = fmt.Sprintf("From: %s\r\nTo: %s\r\nMIME-Version: 1.0\r\nContent-type: text/html\r\nSubject: %s\r\n\r\n%s\r\n", from, to, subject, body)
	}
	auth := smtp.PlainAuth("", from, pass, servername)
	if err := smtp.SendMail(fmt.Sprintf("%s:%s", servername, serverport), auth, from, []string{to}, []byte(msg)); err != nil {
		return err
	}
	log.Printf("email subject=%s has been sent to: %s\n", subject, to)
	return nil
}

//UpdateAmpMailConfig update email config
func UpdateAmpMailConfig(serverAddress string, port string, sender string, pwd string) {
	config.EmailServerAddress=serverAddress
	config.EmailServerPort=port
	config.EmailSender=sender
	config.EmailPwd=pwd
}

//DisplayEncryptedPwd encrypt a pwd to write it in configuration file
func DisplayEncryptedPwd(pwd string) {
	//config := conf.GetRegularConfig(false)
	p, err := newPCipher(config)
	if err != nil {
		log.Println(err)
	}
	data, errd := p.encrypt([]byte(pwd))
	if errd != nil {
		log.Println(errd)
	}
	fmt.Printf("Encrypted pwd: %x\n", string(data))
}

//----------------------------------------------------------------------------------------------------
// emailTemplate functions

func AddEmailTemplate(templateName string, subject string, isHtml bool, body string) {
	template := &emailTemplate{
		subject: subject,
		isHtml:  isHtml,
		body:    body,
	}
	emailTemplateMap[templateName] = template
}

// get and return a copy of email template instance
func getEmailTemplate(templateName string) (*emailTemplate, error) {
	email, ok := emailTemplateMap[templateName]
	if !ok {
		return nil, fmt.Errorf("The templateEmail %s doesn't exist", templateName)
	}
	return &emailTemplate{
		isHtml:  email.isHtml,
		subject: email.subject,
		body:    email.body,
	}, nil
}

func (t *emailTemplate) setVariables(variableMap map[string]string) {
	for name, value := range variableMap {
		t.subject = strings.Replace(t.subject, fmt.Sprintf("{%s}", name), value, -1)
		t.body = strings.Replace(t.body, fmt.Sprintf("{%s}", name), value, -1)
	}
}

//----------------------------------------------------------------------------------------------------
// pCipher function
func newPCipher(config *conf.Configuration) (*pCipher, error) {
	p := &pCipher{}
	addr := config.EmailServerAddress
	for len(addr) < 32 {
		addr = fmt.Sprintf("%s%s", addr, addr)
	}
	addr = addr[0:32]
	block, err := aes.NewCipher([]byte(addr))
	if err != nil {
		return nil, err
	}
	p.block = block
	p.buffer = make([]byte, aes.BlockSize)
	return p, nil
}

func (p *pCipher) encrypt(data []byte) ([]byte, error) {
	stream := cipher.NewCFBEncrypter(p.block, p.buffer)
	stream.XORKeyStream(data, data)
	return data, nil
}

func (p *pCipher) decrypt(data []byte) ([]byte, error) {
	stream := cipher.NewCFBDecrypter(p.block, p.buffer)
	stream.XORKeyStream(data, data)
	return data, nil
}
