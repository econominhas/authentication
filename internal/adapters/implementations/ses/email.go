package ses

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/econominhas/authentication/internal/adapters"
)

type sendEmailInput struct {
	senderEmail   string
	receiverEmail string
	charset       string
	subject       string
	body          string
}

func (adp *SesAdapter) sendEmail(i *sendEmailInput) (*ses.SendEmailOutput, error) {
	ctx, cf := context.WithTimeout(context.Background(), time.Second*5)
	defer cf()

	return adp.ses.SendEmail(ctx, &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{i.receiverEmail},
		},
		Source: &i.senderEmail,
		Message: &types.Message{
			Subject: &types.Content{
				Charset: &i.charset,
				Data:    &i.subject,
			},
			Body: &types.Body{
				Html: &types.Content{
					Charset: &i.charset,
					Data:    &i.body,
				},
			},
		},
	})
}

func (adp *SesAdapter) SendVerificationCodeEmail(i *adapters.SendVerificationCodeEmailInput) error {
	_, err := adp.sendEmail(&sendEmailInput{
		senderEmail:   "foo@bar.com",
		receiverEmail: i.To,
		charset:       "utf-8",
		subject:       "Bem vindo ao Econominhas!",
		body:          "" + i.Code,
	})

	return err
}
