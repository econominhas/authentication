package sns

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/econominhas/authentication/internal/adapters"
)

func (adp *SnsAdapter) SendVerificationCodeSms(i *adapters.SendVerificationCodeSmsInput) error {
	body := "" + i.Code

	_, err := adp.sns.Publish(context.TODO(), &sns.PublishInput{
		PhoneNumber: &i.To,
		Message:     &body,
	})

	return err
}
