package secret

import (
	"github.com/econominhas/authentication/internal/models"
	"github.com/econominhas/authentication/internal/utils"
)

type SecretAdapter struct {
	logger models.Logger
}

func (adp *SecretAdapter) GenSecret(length int) (string, error) {
	return utils.GenRandomString(length), nil
}

func NewSecret(logger models.Logger) *SecretAdapter {
	return &SecretAdapter{
		logger: logger,
	}
}
