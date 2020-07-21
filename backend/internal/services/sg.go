package services

import (
	"fmt"

	"carlosapgomes.com/gobackend/internal/mailer"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type sgService struct {
	key         string
	fromName    string
	fromAddress string
}

// NewMailerService returns a mailer
func NewMailerService(key, name, address string) mailer.Service {
	return &sgService{
		key,
		name,
		address,
	}
}

func (sg *sgService) Send(toName, toEmail, subject, htmlContent string) (*mailer.Response, error) {
	from := mail.NewEmail(sg.fromName, sg.fromAddress)
	to := mail.NewEmail(toName, toEmail)
	plainTextContent := "plain text"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(sg.key)
	var res mailer.Response
	if client != nil {
		r, err := client.Send(message)
		res.Code = r.StatusCode
		res.Msg = r.Body
		res.Msg += "\n"
		for k, v := range r.Headers {
			res.Msg += fmt.Sprintf("%s=\"%s\"\n", k, v)
		}
		if err != nil {
			return &res, err
		}
	}
	return &res, nil
}
