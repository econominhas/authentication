package adapters

type SendVerificationCodeEmailInput struct {
	To string
}

type EmailAdapter interface {
	SendVerificationCodeEmail(i *SendVerificationCodeEmailInput) error
}
