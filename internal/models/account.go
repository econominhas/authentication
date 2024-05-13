package models

import (
	"database/sql"
	"time"
)

// ----------------------------
//
//      Repository
//
// ----------------------------

type CreateAccountSignInProvider struct {
	Id           string
	Type         string
	AccessToken  string
	RefreshToken *string
	ExpiresAt    time.Time
}

type CreateAccountPhone struct {
	CountryCode string
	Number      string
}

type CreateAccountInput struct {
	Db sql.Tx

	Email           string
	Phone           CreateAccountPhone
	SignInProviders []CreateAccountSignInProvider
}

type CreateAccountOutput struct {
	Id string
}

type GetAccountByEmailInput struct {
	Db sql.Tx

	Email string
}

type GetAccountByEmailOutput struct {
	AccountId string
	Email     string
}

type GetManyAccountsByProviderInput struct {
	Db sql.Tx

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

	GetManyByProvider(i *GetManyAccountsByProviderInput) ([]GetManyAccountsByProviderOutput, error)

	GetByEmail(i *GetAccountByEmailInput) (*GetAccountByEmailOutput, error)
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
	CreateFromEmail(i *CreateAccountFromEmailInput) error

	CreateFromPhone(i *CreateAccountFromPhoneInput) error

	CreateFromGoogleProvider(i *CreateAccountFromExternalProviderInput) (*AuthOutput, error)

	CreateFromFacebookProvider(i *CreateAccountFromExternalProviderInput) (*AuthOutput, error)

	ExchangeCode(i *ExchangeAccountCodeInput) (*AuthOutput, error)

	RefreshToken(i *ExchangeAccountCodeInput) (*RefreshAccountTokenOutput, error)
}
