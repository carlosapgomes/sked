package mocks

import "carlosapgomes.com/sked/internal/mailer"

// MailerMockSvc mocks a mailer
type MailerMockSvc struct {
	Rec *EmailRecorder
}

//EmailRecorder stores email data
type EmailRecorder struct {
	Email string
	Name  string
	Msg   string
}

// NewMailerMock returns a pointer to a MailerMockSvc
func NewMailerMock(rec *EmailRecorder) *MailerMockSvc {
	return &MailerMockSvc{
		rec,
	}
}

// Send email
func (m MailerMockSvc) Send(toName, toEmail, subject, htmlContent string) (*mailer.Response, error) {
	if m.Rec != nil {
		m.Rec.Msg = htmlContent
		m.Rec.Email = toEmail
		m.Rec.Name = toName
	}
	return nil, nil
}
