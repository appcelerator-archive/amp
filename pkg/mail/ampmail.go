package mail

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"github.com/appcelerator/amp/pkg/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var emailTemplateMap map[string]*emailTemplate
var config *amp.Config

type emailTemplate struct {
	isHTML  bool
	subject string
	body    string
}

func init() {
	config = amp.GetConfig()
	emailTemplateMap = make(map[string]*emailTemplate)
	AddEmailTemplate("AccountVerification", "AMP Account verification", true, accountVerificationBody)
	AddEmailTemplate("AccountCreated", "AMP account creation confirmation", true, accountCreationBody)

	AddEmailTemplate("organizationCreated", "AMP organization creation confirmation", true, organizationCreationBody)
	AddEmailTemplate("userAddedInOrganization", "AMP user added in organization confirmation", true, addUserToOrganizationBody)
	AddEmailTemplate("userRemoveFromOrganization", "AMP user removed from organization confirmation", true, removeUserFromOrganizationBody)
	AddEmailTemplate("organizationRemoved", "AMP organization removed confirmation", true, organizationRemoveBody)

	AddEmailTemplate("teamCreated", "AMP team creation confirmation", true, teamCreationBody)
	AddEmailTemplate("userAddedInTeam", "AMP user added in team confirmation", true, addUserToTeamBody)
	AddEmailTemplate("userRemoveFromTeam", "AMP user removed from team confirmation", true, removeUserFromTeamBody)
	AddEmailTemplate("teamRemoved", "AMP team removed confirmation", true, teamRemoveBody)

	AddEmailTemplate("AccountResetPassword", "AMP reset password", true, accountResetPasswordEmailBody)
	AddEmailTemplate("AccountPasswordConfirmation", "AMP reset password confirmation", true, accountPasswordConfirmationEmailBody)
	AddEmailTemplate("AccountNameReminder", "AMP account name reminder", true, accountNameReminderBody)
}

// SendAccountVerificationEmail send a AccountVerification email template
func SendAccountVerificationEmail(to string, accountName string, token string) error {
	//config := conf.GetRegularConfig(false)
	variables := map[string]string{
		"accountName": accountName,
		"token":       token,
		"ampAddress":  config.ServerAddress,
	}
	if err := SendTemplateEmail(to, "AccountVerification", variables); err != nil {
		return err
	}
	return nil
}

// SendAccountVerificationEmail send a AccountVerification email template
func SendAccountCreatedEmail(to string, accountName string) error {
	//config := conf.GetRegularConfig(false)
	variables := map[string]string{
		"accountName": accountName,
	}
	if err := SendTemplateEmail(to, "AccountVerification", variables); err != nil {
		return err
	}
	return nil
}

//SendOrganizationCreatedEmail send mail
func SendOrganizationCreatedEmail(to string, organization string) error {
	variables := map[string]string{
		"organization": organization,
	}
	if err := SendTemplateEmail(to, "organizationCreated", variables); err != nil {
		return err
	}
	return nil
}

//SendUserAddedInOrganizationEmail send mail
func SendUserAddedInOrganizationEmail(to string, organization string, user string) error {
	variables := map[string]string{
		"organization": organization,
		"user":         user,
	}
	if err := SendTemplateEmail(to, "userAddedInOrganization", variables); err != nil {
		return err
	}
	return nil
}

//SendUserRemovedFromOrganizationEmail send mail
func SendUserRemovedFromOrganizationEmail(to string, organization string, user string) error {
	variables := map[string]string{
		"organization": organization,
		"user":         user,
	}
	if err := SendTemplateEmail(to, "userRemoveFromOrganization", variables); err != nil {
		return err
	}
	return nil
}

//SendOrganizationRemovedEmail send mail
func SendOrganizationRemovedEmail(to string, organization string) error {
	variables := map[string]string{
		"organization": organization,
	}
	if err := SendTemplateEmail(to, "organizationRemoveBody", variables); err != nil {
		return err
	}
	return nil
}

//SendTeamCreatedEmail send mail
func SendTeamCreatedEmail(to string, team string) error {
	variables := map[string]string{
		"team": team,
	}
	if err := SendTemplateEmail(to, "teamCreated", variables); err != nil {
		return err
	}
	return nil
}

//SendUserAddedInTeamEmail send mail
func SendUserAddedInTeamEmail(to string, team string, user string) error {
	variables := map[string]string{
		"team": team,
		"user": user,
	}
	if err := SendTemplateEmail(to, "userAddedInTeam", variables); err != nil {
		return err
	}
	return nil
}

//SendUserRemovedFromTeamEmail send mail
func SendUserRemovedFromTeamEmail(to string, team string, user string) error {
	variables := map[string]string{
		"team": team,
		"user": user,
	}
	if err := SendTemplateEmail(to, "userRemoveFromTeam", variables); err != nil {
		return err
	}
	return nil
}

//SendTeamRemovedEmail send mail
func SendTeamRemovedEmail(to string, team string) error {
	variables := map[string]string{
		"team": team,
	}
	if err := SendTemplateEmail(to, "teamRemoveBody", variables); err != nil {
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
		"ampAddress":  config.ServerAddress,
	}
	if err := SendTemplateEmail(to, "AccountResetPassword", variables); err != nil {
		return err
	}
	return nil
}

// SendAccountPasswordConfirmationEmail send a AccountResetPassword email template
func SendAccountPasswordConfirmationEmail(to string, accountName string) error {
	variables := map[string]string{
		"accountName": accountName,
	}
	if err := SendTemplateEmail(to, "AccountPasswordConfirmation", variables); err != nil {
		return err
	}
	return nil
}

// SendAccountNameReminderEmail send a AccountNameReminder email template
func SendAccountNameReminderEmail(to string, accountName string) error {
	variables := map[string]string{
		"accountName": accountName,
	}
	if err := SendTemplateEmail(to, "AccountNameReminder", variables); err != nil {
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
	if err := SendMail(to, email.subject, email.isHTML, email.body); err != nil {
		return err
	}
	return nil
}

// SendMail send an eamail to "to" with subject and body, use configuration
func SendMail(to string, subject string, isHTML bool, body string) error {
	if config.EmailServerAddress == "" {
		return sendMailUsingSendGrid(to, subject, body)
	}
	from := config.EmailSender
	servername := config.EmailServerAddress
	serverport := config.EmailServerPort
	pass := config.EmailPwd
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n", from, to, subject, body)
	if isHTML {
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
	config.EmailServerAddress = serverAddress
	config.EmailServerPort = port
	config.EmailSender = sender
	config.EmailPwd = pwd
}

func sendMailUsingSendGrid(to string, subject string, body string) error {
	apiKey := config.EmailKey
	from := mail.NewEmail("amp", config.EmailSender)
	target := mail.NewEmail(strings.Split(to, "@")[0], to)
	content := mail.NewContent("text/html", body)
	m := mail.NewV3MailInit(from, subject, target, content)

	request := sendgrid.GetRequest(apiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	if _, err := sendgrid.API(request); err != nil {
		return err
	}
	return nil
}

//----------------------------------------------------------------------------------------------------
// emailTemplate functions

// AddEmailTemplate add templete in the map
func AddEmailTemplate(templateName string, subject string, isHTML bool, body string) {
	template := &emailTemplate{
		subject: subject,
		isHTML:  isHTML,
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
		isHTML:  email.isHTML,
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
