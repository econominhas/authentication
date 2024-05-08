package email

import (
	"github.com/aws/aws-sdk-go-v2/service/ses"
)

type SesAdapter struct {
	Ses *ses.Client
}

func NewSesAdapter() *SesAdapter {
	sesClient := ses.New(ses.Options{})

	return &SesAdapter{
		Ses: sesClient,
	}
}
