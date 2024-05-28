package discord

import "github.com/econominhas/authentication/internal/utils"

func (adp *DiscordAdapter) HasRequiredScopes(scopes []string) bool {
	requiredScopes := []string{
		"identify",
		"email",
	}

	return utils.AllInArray(requiredScopes, scopes)
}
