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

type CreateRefreshTokenInput struct {
	Db sql.Tx

	AccountId string
}

type CreateRefreshTokenOutput struct {
	RefreshToken string
}

type GetRefreshTokenInput struct {
	Db sql.Tx

	RefreshToken string
}

type GetRefreshTokenOutput struct {
	AccountId string
	CreatedAt time.Time
}

type RefreshTokenRepository interface {
	Create(i *CreateRefreshTokenInput) (*CreateRefreshTokenOutput, error)
	Get(i *GetRefreshTokenInput) (*GetRefreshTokenOutput, error)
}
