package secret

import "github.com/econominhas/authentication/internal/utils"

func GenSecret(length int) string {
	return utils.GenRandomString(length)
}
