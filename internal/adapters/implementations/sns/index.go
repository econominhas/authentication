package sns

import "github.com/aws/aws-sdk-go-v2/service/sns"

type SnsAdapter struct {
	Sns *sns.Client
}

func NewSesAdapter() *SnsAdapter {
	snsClient := sns.New(sns.Options{})

	return &SnsAdapter{
		Sns: snsClient,
	}
}
