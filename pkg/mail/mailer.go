package mail

import (
	"fmt"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type emailTemplate struct {
	isHTML  bool
	subject string
	body    string
}

type Mailer struct {
	apiKey           string
	emailSender      string
	publicAddress    string
	emailTemplateMap map[string]*emailTemplate
}

func NewMailer(apiKey string, emailSender string, publicAddress string) *Mailer {
	mailer := &Mailer{
		apiKey:           apiKey,
		emailSender:      emailSender,
		publicAddress:    publicAddress,
		emailTemplateMap: make(map[string]*emailTemplate),
	}

	mailer.AddEmailTemplate("AccountVerification", "AMP Account verification", true, accountVerificationBody)
	mailer.AddEmailTemplate("AccountCreated", "AMP account creation confirmation", true, accountCreationBody)
	mailer.AddEmailTemplate("AccountRemoved", "AMP account removed", true, accountRemovedBody)

	mailer.AddEmailTemplate("organizationCreated", "AMP organization creation confirmation", true, organizationCreationBody)
	mailer.AddEmailTemplate("userAddedInOrganization", "AMP user added in organization confirmation", true, addUserToOrganizationBody)
	mailer.AddEmailTemplate("userRemoveFromOrganization", "AMP user removed from organization confirmation", true, removeUserFromOrganizationBody)
	mailer.AddEmailTemplate("organizationRemoved", "AMP organization removed confirmation", true, organizationRemoveBody)

	mailer.AddEmailTemplate("teamCreated", "AMP team creation confirmation", true, teamCreationBody)
	mailer.AddEmailTemplate("userAddedInTeam", "AMP user added in team confirmation", true, addUserToTeamBody)
	mailer.AddEmailTemplate("userRemoveFromTeam", "AMP user removed from team confirmation", true, removeUserFromTeamBody)
	mailer.AddEmailTemplate("teamRemoved", "AMP team removed confirmation", true, teamRemoveBody)

	mailer.AddEmailTemplate("AccountResetPassword", "AMP reset password", true, accountResetPasswordEmailBody)
	mailer.AddEmailTemplate("AccountPasswordConfirmation", "AMP reset password confirmation", true, accountPasswordConfirmationEmailBody)
	mailer.AddEmailTemplate("AccountNameReminder", "AMP account name reminder", true, accountNameReminderBody)

	return mailer
}

// SendAccountVerificationEmail send mail
func (m *Mailer) SendAccountVerificationEmail(to string, accountName string, token string) error {
	//config := conf.GetRegularConfig(false)
	variables := map[string]string{
		"accountName": accountName,
		"token":       token,
		"ampAddress":  m.publicAddress,
	}
	if err := m.SendTemplateEmail(to, "AccountVerification", variables); err != nil {
		return err
	}
	return nil
}

// SendAccountCreatedEmail send mail
func (m *Mailer) SendAccountCreatedEmail(to string, accountName string) error {
	//config := conf.GetRegularConfig(false)
	variables := map[string]string{
		"accountName": accountName,
	}
	if err := m.SendTemplateEmail(to, "AccountCreated", variables); err != nil {
		return err
	}
	return nil
}

// SendAccountRemovedEmail send mail
func (m *Mailer) SendAccountRemovedEmail(to string, accountName string) error {
	//config := conf.GetRegularConfig(false)
	variables := map[string]string{
		"accountName": accountName,
	}
	if err := m.SendTemplateEmail(to, "AccountRemoved", variables); err != nil {
		return err
	}
	return nil
}

//SendOrganizationCreatedEmail send mail
func (m *Mailer) SendOrganizationCreatedEmail(to string, organization string) error {
	variables := map[string]string{
		"organization": organization,
	}
	if err := m.SendTemplateEmail(to, "organizationCreated", variables); err != nil {
		return err
	}
	return nil
}

//SendUserAddedInOrganizationEmail send mail
func (m *Mailer) SendUserAddedInOrganizationEmail(to string, organization string, user string) error {
	variables := map[string]string{
		"organization": organization,
		"user":         user,
	}
	if err := m.SendTemplateEmail(to, "userAddedInOrganization", variables); err != nil {
		return err
	}
	return nil
}

//SendUserRemovedFromOrganizationEmail send mail
func (m *Mailer) SendUserRemovedFromOrganizationEmail(to string, organization string, user string) error {
	variables := map[string]string{
		"organization": organization,
		"user":         user,
	}
	if err := m.SendTemplateEmail(to, "userRemoveFromOrganization", variables); err != nil {
		return err
	}
	return nil
}

//SendOrganizationRemovedEmail send mail
func (m *Mailer) SendOrganizationRemovedEmail(to string, organization string) error {
	variables := map[string]string{
		"organization": organization,
	}
	if err := m.SendTemplateEmail(to, "organizationRemoved", variables); err != nil {
		return err
	}
	return nil
}

//SendTeamCreatedEmail send mail
func (m *Mailer) SendTeamCreatedEmail(to string, team string) error {
	variables := map[string]string{
		"team": team,
	}
	if err := m.SendTemplateEmail(to, "teamCreated", variables); err != nil {
		return err
	}
	return nil
}

//SendUserAddedInTeamEmail send mail
func (m *Mailer) SendUserAddedInTeamEmail(to string, team string, user string) error {
	variables := map[string]string{
		"team": team,
		"user": user,
	}
	if err := m.SendTemplateEmail(to, "userAddedInTeam", variables); err != nil {
		return err
	}
	return nil
}

//SendUserRemovedFromTeamEmail send mail
func (m *Mailer) SendUserRemovedFromTeamEmail(to string, team string, user string) error {
	variables := map[string]string{
		"team": team,
		"user": user,
	}
	if err := m.SendTemplateEmail(to, "userRemoveFromTeam", variables); err != nil {
		return err
	}
	return nil
}

//SendTeamRemovedEmail send mail
func (m *Mailer) SendTeamRemovedEmail(to string, team string) error {
	variables := map[string]string{
		"team": team,
	}
	if err := m.SendTemplateEmail(to, "teamRemoved", variables); err != nil {
		return err
	}
	return nil
}

// SendAccountResetPasswordEmail send a AccountResetPassword email template
func (m *Mailer) SendAccountResetPasswordEmail(to string, accountName string, token string) error {
	//config := conf.GetRegularConfig(false)
	variables := map[string]string{
		"accountName": accountName,
		"token":       token,
		"ampAddress":  m.publicAddress,
	}
	if err := m.SendTemplateEmail(to, "AccountResetPassword", variables); err != nil {
		return err
	}
	return nil
}

// SendAccountNameReminderEmail send a AccountNameReminder email template
func (m *Mailer) SendAccountNameReminderEmail(to string, accountName string) error {
	variables := map[string]string{
		"accountName": accountName,
	}
	if err := m.SendTemplateEmail(to, "AccountNameReminder", variables); err != nil {
		return err
	}
	return nil
}

// SendTemplateEmail send a tempalte email
func (m *Mailer) SendTemplateEmail(to string, templateEmailName string, variableMap map[string]string) error {
	email, err := m.getEmailTemplate(templateEmailName)
	if err != nil {
		return err
	}
	email.setVariables(variableMap)
	if err := m.SendMail(to, email.subject, email.isHTML, email.body); err != nil {
		return err
	}
	return nil
}

// SendMail send an email to "to" with subject and body, use configuration
func (m *Mailer) SendMail(to string, subject string, isHTML bool, body string) error {
	from := mail.NewEmail("amp", m.emailSender)
	target := mail.NewEmail(strings.Split(to, "@")[0], to)
	cType := "text/plain"
	if isHTML {
		cType = "text/html"
	}
	content := mail.NewContent(cType, body)
	mailer := mail.NewV3MailInit(from, subject, target, content)

	request := sendgrid.GetRequest(m.apiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(mailer)
	if _, err := sendgrid.API(request); err != nil {
		return err
	}
	return nil
}

//----------------------------------------------------------------------------------------------------
// emailTemplate functions

// AddEmailTemplate add template in the map
func (m *Mailer) AddEmailTemplate(templateName string, subject string, isHTML bool, body string) {
	template := &emailTemplate{
		subject: subject,
		isHTML:  isHTML,
		body:    body,
	}
	m.emailTemplateMap[templateName] = template
}

// get and return a copy of email template instance
func (m *Mailer) getEmailTemplate(templateName string) (*emailTemplate, error) {
	email, ok := m.emailTemplateMap[templateName]
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
