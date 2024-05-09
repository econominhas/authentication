package adapters

type SendVerificationCodeEmailInput struct {
	To   string
	Code string
}

type EmailAdapter interface {
	SendVerificationCodeEmail(i *SendVerificationCodeEmailInput) error
}
