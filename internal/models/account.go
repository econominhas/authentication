package models

import "time"

// ----------------------------
//
//      Repository
//
// ----------------------------

type CreateAccountInput struct {
	Email  string
	Phone  string
	Google struct {
		Id           string
		AccessToken  string
		RefreshToken string
		ExpiresAt    time.Time
	}
	Facebook struct {
		Id          string
		AccessToken string
		ExpiresAt   time.Time
	}
}

type CreateAccountOutput struct {
	Id string
}

type GetManyAccountsByProviderInput struct {
	ProviderId   string
	ProviderType string
	Email        string
}

type GetManyAccountsByProviderOutput struct {
	AccountId    string
	ProviderId   string
	ProviderType string
	Email        string
}

type AccountRepository interface {
	Create(i *CreateAccountInput) (*CreateAccountOutput, error)

	GetManyByProvider(i *GetManyAccountsByProviderInput) (*GetManyAccountsByProviderOutput, error)
}

// ----------------------------
//
//      Service
//
// ----------------------------

type AuthOutput struct {
	RefreshToken string
	AccessToken  string
	ExpiresAt    time.Time
}

type CreateAccountFromEmailInput struct {
	Email string
}

type CreateAccountFromPhoneInput struct {
	Phone string
}

type CreateAccountFromExternalProviderInput struct {
	Code      string
	OriginUrl string
}

type ExchangeAccountCodeInput struct {
	Code      string
	OriginUrl string
}
type RefreshAccountTokenInput struct {
	RefreshToken string
}

type RefreshAccountTokenOutput struct {
	AccessToken string
	ExpiresAt   time.Time
}

type AccountService interface {
	CreateFromEmail(i *CreateAccountFromEmailInput) (*CreateAccountOutput, error)

	CreateFromPhone(i *CreateAccountFromPhoneInput) (*CreateAccountOutput, error)

	CreateFromGoogleProvider(i *CreateAccountFromExternalProviderInput) (*AuthOutput, error)

	CreateFromFacebookProvider(i *CreateAccountFromExternalProviderInput) (*AuthOutput, error)

	ExchangeCode(i *ExchangeAccountCodeInput) (*AuthOutput, error)

	RefreshToken(i *ExchangeAccountCodeInput) (*RefreshAccountTokenOutput, error)
}
