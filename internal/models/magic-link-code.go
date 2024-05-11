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

type UpsertMagicLinkRefreshTokenInput struct {
	Db sql.Tx

	AccountId     string
	IsFirstAccess bool
}

type GetInput struct {
	Db sql.Tx

	AccountId string
	Code      string
}

type MagicLinkCodeRepository interface {
	Upsert(i *UpsertMagicLinkRefreshTokenInput) (*MagicLinkCodeEntity, error)
	Get(i *GetInput) (*MagicLinkCodeEntity, error)
}
