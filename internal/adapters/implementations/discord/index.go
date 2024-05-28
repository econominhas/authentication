package discord

import (
	"net/http"

	"github.com/econominhas/authentication/internal/models"
)

type DiscordAdapter struct {
	logger models.Logger

	httpClient *http.Client
}

func NewDiscord(logger models.Logger) *DiscordAdapter {
	return &DiscordAdapter{
		logger:     logger,
		httpClient: &http.Client{},
	}
}
