package facebook

import (
	"net/http"

	"github.com/econominhas/authentication/internal/models"
)

type FacebookAdapter struct {
	logger models.Logger

	httpClient *http.Client
}

func NewFacebook(logger models.Logger) *FacebookAdapter {
	return &FacebookAdapter{
		logger:     logger,
		httpClient: &http.Client{},
	}
}
