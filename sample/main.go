package main

import (
	Mailer "github.com/artziel/go-mailer"
)

type viewData struct {
	Fullname string
	Username string
	Password string
}

const welcomeEmailTplFile = "templates/email/welcome.ghtml"

func main() {

	tpl, err := Mailer.ReadTemplateFile(welcomeEmailTplFile)
	if err != nil {
		panic(err)
	}

	cnf := Mailer.MailerConfig{
		Host:     "",
		Port:     25,
		Username: "",
		Password: "",
		From:     "",
		TLS:      true,
	}

	if srv, err := Mailer.Service(cnf); err != nil {
		panic(err)
	} else {

		if srv.SendMail(
			Mailer.Mail{
				To:       []string{"sample@sample.com"},
				Bcc:      []string{"unknown.a@sample.com", "unknown.b@sample.com"},
				Subject:  "Esta es una prueba de correo con Golang",
				Body:     tpl,
				ViewData: viewData{},
				Attachments: []Mailer.File{
					Mailer.InlineAttachment("assets/img/email/alus-logo.png"),
				},
			},
		); err != nil {
			panic(err)
		}
	}

}
