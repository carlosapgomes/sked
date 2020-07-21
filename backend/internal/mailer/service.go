package mailer

// Service describes a mailer interface
type Service interface {
	Send(toName, toEmail, subject, htmlContent string) (*Response, error)
}

// Response stores response data from mailer functions
type Response struct {
	Code int
	Msg  string
}
