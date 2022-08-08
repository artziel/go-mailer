package Mailer

import (
	"bytes"
	"html/template"
	"io/ioutil"

	mail "github.com/xhit/go-simple-mail/v2"
)

type File struct {
	FilePath string
	Inline   bool
}

type MailerConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
	TLS      bool   `yaml:"tls"`
}

type MailService struct {
	config MailerConfig
	server *mail.SMTPClient
}

type Mail struct {
	From        string
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	Body        string
	ViewData    interface{}
	Attachments []File
}

func (m *Mail) getAttachments() []mail.File {
	attachments := []mail.File{}
	for _, a := range m.Attachments {
		attachments = append(attachments, mail.File{FilePath: a.FilePath, Inline: a.Inline})

	}
	return attachments
}

func (s *MailService) SendMail(m Mail) error {

	if m.Body == "" {
		return ErrEmptyBody
	}
	if m.Subject == "" {
		return ErrEmptySubject
	}
	if m.To == nil || len(m.To) == 0 {
		return ErrNoRecipients
	}
	if m.From == "" && s.config.From == "" {
		return ErrNoRemitent
	}

	// Create email
	email := mail.NewMSG()
	if m.From == "" {
		email.SetFrom(s.config.From)
	} else {
		email.SetFrom(m.From)
	}

	for _, to := range m.To {
		email.AddTo(to)
	}
	for _, cc := range m.Cc {
		email.AddCc(cc)
	}
	for _, bcc := range m.Bcc {
		email.AddBcc(bcc)
	}

	for _, a := range m.getAttachments() {
		email.Attach(&a)
	}

	email.SetSubject(m.Subject)

	buf := &bytes.Buffer{}
	if tmpl, err := template.New("").Parse(m.Body); err != nil {
		return err
	} else {
		if err := tmpl.Execute(buf, m.ViewData); err != nil {
			return err
		}
	}

	email.SetBody(mail.TextHTML, buf.String())

	if err := email.Send(s.server); err != nil {
		return err
	}

	return nil
}

func InlineAttachment(filePath string) File {
	return File{FilePath: filePath, Inline: true}
}

func Attachment(filePath string) mail.File {
	return mail.File{FilePath: filePath}
}

func Service(cnf MailerConfig) (MailService, error) {
	srv := MailService{}

	server := mail.NewSMTPClient()
	server.Host = cnf.Host
	server.Port = cnf.Port
	server.Username = cnf.Username
	server.Password = cnf.Password
	if cnf.TLS {
		server.Encryption = mail.EncryptionTLS
	}

	smtpClient, err := server.Connect()
	if err != nil {
		return srv, err
	}

	srv = MailService{config: cnf, server: smtpClient}

	return srv, nil
}

func ReadTemplateFile(file string) (string, error) {
	tpl, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(tpl), nil
}
