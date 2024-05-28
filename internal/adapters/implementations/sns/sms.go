package sns

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/econominhas/authentication/internal/adapters"
)

type sendSmsInput struct {
	receiverPhoneNumber string
	body                string
}

func (adp *SnsAdapter) sendSms(i *sendSmsInput) (*sns.PublishOutput, error) {
	ctx, cf := context.WithTimeout(context.Background(), time.Second*5)
	defer cf()

	return adp.sns.Publish(ctx, &sns.PublishInput{
		PhoneNumber: &i.receiverPhoneNumber,
		Message:     &i.body,
	})
}

func (adp *SnsAdapter) SendVerificationCodeSms(i *adapters.SendVerificationCodeSmsInput) error {
	_, err := adp.sendSms(&sendSmsInput{
		receiverPhoneNumber: i.To,
		body:                "" + i.Code,
	})

	return err
}
