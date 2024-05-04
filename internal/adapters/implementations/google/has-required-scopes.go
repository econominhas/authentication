package google

import "github.com/econominhas/authentication/internal/utils"

func (adp *GoogleAdapter) HasRequiredScopes(scopes []string) bool {
	requiredScopes := []string{
		"https://www.googleapis.com/auth/userinfo.profile",
		"openid",
		"https://www.googleapis.com/auth/userinfo.email",
	}

	return utils.AllInArray(requiredScopes, scopes)
}
