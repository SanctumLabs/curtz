package email

type EmailService struct{}

func NewEmailSvc() *EmailService {
	return &EmailService{}
}

func (e *EmailService) SendEmail(recipient, subject, body string) error {
	panic("implement me")
}
