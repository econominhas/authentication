package models

import "database/sql"

// ----------------------------
//
//      Repository
//
// ----------------------------

type CreateRefreshTokenInput struct {
	Db sql.Tx

	AccountId string
}

type CreateRefreshTokenOutput struct {
	RefreshToken string
}

type RefreshTokenRepository interface {
	Create(i *CreateRefreshTokenInput) (*CreateRefreshTokenOutput, error)
}
