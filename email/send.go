package email

type EmailPayload struct {
	Subject string
	Content string
}

type Emailer interface {
	Send()
}
