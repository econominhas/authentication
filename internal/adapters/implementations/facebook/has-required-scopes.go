package facebook

import "github.com/econominhas/authentication/internal/utils"

func (adp *FacebookAdapter) HasRequiredScopes(scopes []string) bool {
	requiredScopes := []string{
		"public_profile",
		"email",
	}

	return utils.AllInArray(requiredScopes, scopes)
}
