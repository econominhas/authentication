package ses

import (
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/econominhas/authentication/internal/models"
)

type SesAdapter struct {
	logger models.Logger

	ses *ses.Client
}

func NewSes(logger models.Logger) *SesAdapter {
	sesClient := ses.New(ses.Options{})

	return &SesAdapter{
		logger: logger,
		ses:    sesClient,
	}
}
