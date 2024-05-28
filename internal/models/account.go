package models

import (
	"database/sql"
	"time"

	"github.com/econominhas/authentication/internal/utils"
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
	Db *sql.Tx

	Email           string
	Phone           CreateAccountPhone
	SignInProviders []CreateAccountSignInProvider
}

type CreateAccountOutput struct {
	Id string
}

type GetAccountByEmailInput struct {
	Db *sql.Tx

	Email string
}

type GetAccountByEmailOutput struct {
	AccountId string
	Email     string
}

type GetAccountByPhoneInput struct {
	Db *sql.Tx

	CountryCode string
	Number      string
}

type GetAccountByPhoneOutput struct {
	AccountId   string
	CountryCode string
	Number      string
}

type GetManyAccountsByProviderInput struct {
	Db *sql.Tx

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

	GetByPhone(i *GetAccountByPhoneInput) (*GetAccountByPhoneOutput, error)
}

// ----------------------------
//
//      Service
//
// ----------------------------

type AuthOutput struct {
	RefreshToken string    `json:"refreshToken,omitempty"`
	AccessToken  string    `json:"accessToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

type CreateAccountFromEmailInput struct {
	Email string `validate:"required,email"`
}

type CreateAccountFromPhoneInput struct {
	Phone CreateAccountPhone `validate:"required,numeric"`
}

type CreateAccountFromExternalProviderInput struct {
	Code      string `validate:"required,alphanum"`
	OriginUrl string `validate:"url"`
}

type ExchangeAccountCodeInput struct {
	AccountId string `validate:"required,ulid"`
	Code      string `validate:"required,alphanum"`
}

type RefreshAccountTokenInput struct {
	AccountId    string `validate:"required,ulid"`
	RefreshToken string `validate:"required,alphanum"`
}

type RefreshAccountTokenOutput struct {
	AccessToken string    `json:"accessToken"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

type AccountService interface {
	CreateFromEmailProvider(i *CreateAccountFromEmailInput) *utils.HttpError

	CreateFromPhoneProvider(i *CreateAccountFromPhoneInput) *utils.HttpError

	CreateFromGoogleProvider(i *CreateAccountFromExternalProviderInput) (*AuthOutput, *utils.HttpError)

	CreateFromFacebookProvider(i *CreateAccountFromExternalProviderInput) (*AuthOutput, *utils.HttpError)

	ExchangeCode(i *ExchangeAccountCodeInput) (*AuthOutput, *utils.HttpError)

	RefreshToken(i *RefreshAccountTokenInput) (*RefreshAccountTokenOutput, *utils.HttpError)
}
