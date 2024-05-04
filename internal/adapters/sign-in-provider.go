package adapters

import "time"

type ExchangeCodeInput struct {
	Code      string
	OriginUrl string
}

type ExchangeCodeOutput struct {
	Scopes       []string
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

type GetAuthenticatedUserDataOutput struct {
	Id              string
	Name            string
	Email           string
	IsEmailVerified bool
}

type SignInProviderAdapter interface {
	ExchangeCode(i *ExchangeCodeInput) (*ExchangeCodeOutput, error)

	GetUserData(accessToken string) (*GetAuthenticatedUserDataOutput, error)

	HasRequiredScopes(scopes []string) bool
}
