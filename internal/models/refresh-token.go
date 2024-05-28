package models

import (
	"database/sql"
	"time"
)

// ----------------------------
//
//      Entities
//
// ----------------------------

type RefreshTokenEntity struct {
	AccountId    string    `db:"account_id"`
	RefreshToken string    `db:"refresh_token"`
	CreatedAt    time.Time `db:"created_at"`
}

// ----------------------------
//
//      Repository
//
// ----------------------------

type CreateRefreshTokenInput struct {
	Db *sql.Tx

	AccountId string
}

type CreateRefreshTokenOutput struct {
	RefreshToken string `json:"refreshToken"`
}

type GetRefreshTokenInput struct {
	Db *sql.Tx

	AccountId    string
	RefreshToken string
}

type RefreshTokenRepository interface {
	Create(i *CreateRefreshTokenInput) (*CreateRefreshTokenOutput, error)

	Get(i *GetRefreshTokenInput) (bool, error)
}
