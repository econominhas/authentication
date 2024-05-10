package secret

import "github.com/econominhas/authentication/internal/utils"

type SecretAdapter struct{}

func (adp *SecretAdapter) GenSecret(length int) (string, error) {
	return utils.GenRandomString(length), nil
}
