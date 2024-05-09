package adapters

type SendVerificationCodeSmsInput struct {
	To   string
	Code string
}

type SmsAdapter interface {
	SendVerificationCodeSms(i *SendVerificationCodeSmsInput) error
}
