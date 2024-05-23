package google

import (
	"net/http"

	"github.com/econominhas/authentication/internal/models"
)

type GoogleAdapter struct {
	logger models.Logger

	httpClient *http.Client
}

func NewGoogle(logger models.Logger) *GoogleAdapter {
	return &GoogleAdapter{
		logger:     logger,
		httpClient: &http.Client{},
	}
}
