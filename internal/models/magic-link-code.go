package models

import (
	"database/sql"
	"time"
)

type MagicLinkCodeEntity struct {
	AccountId     string
	Code          string
	IsFirstAccess bool
	CreatedAt     time.Time
}

// ----------------------------
//
//      Repository
//
// ----------------------------

type UpsertRefreshTokenInput struct {
	Db sql.Tx

	AccountId     string
	IsFirstAccess bool
}

type GetRefreshTokenInput struct {
	Db sql.Tx

	AccountId string
	Code      string
}

type MagicLinkCodeRepository interface {
	Upsert(i *UpsertRefreshTokenInput) (*MagicLinkCodeEntity, error)
	Get(i *GetRefreshTokenInput) (*MagicLinkCodeEntity, error)
}
