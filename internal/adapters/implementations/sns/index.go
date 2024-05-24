package sns

import (
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/econominhas/authentication/internal/models"
)

type SnsAdapter struct {
	logger models.Logger

	sns *sns.Client
}

func NewSns(logger models.Logger) *SnsAdapter {
	snsClient := sns.New(sns.Options{})

	return &SnsAdapter{
		logger: logger,
		sns:    snsClient,
	}
}
