package ses

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/econominhas/authentication/internal/adapters"
)

func (adp *SesAdapter) SendVerificationCodeEmail(i *adapters.SendVerificationCodeEmailInput) error {
	senderEmail := "foo@bar.com"
	charset := "utf-8"
	subject := "Bem vindo ao Econominhas!"
	body := "" + i.Code

	_, err := adp.Ses.SendEmail(context.TODO(), &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{i.To},
		},
		Source: &senderEmail,
		Message: &types.Message{
			Subject: &types.Content{
				Charset: &charset,
				Data:    &subject,
			},
			Body: &types.Body{
				Html: &types.Content{
					Charset: &charset,
					Data:    &body,
				},
			},
		},
	})

	return err
}
