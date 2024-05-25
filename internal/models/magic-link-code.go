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

type MagicLinkCodeEntity struct {
	AccountId     string    `db:"account_id"`
	Code          string    `db:"code"`
	IsFirstAccess bool      `db:"is_first_access"`
	CreatedAt     time.Time `db:"created_at"`
}

// ----------------------------
//
//      Repository
//
// ----------------------------

type UpsertMagicLinkRefreshTokenInput struct {
	Db *sql.Tx

	AccountId     string
	IsFirstAccess bool
}

type GetMagicLinkRefreshTokenInput struct {
	Db *sql.Tx

	AccountId string
	Code      string
}

type MagicLinkCodeRepository interface {
	Upsert(i *UpsertMagicLinkRefreshTokenInput) (*MagicLinkCodeEntity, error)

	Get(i *GetMagicLinkRefreshTokenInput) (*MagicLinkCodeEntity, error)
}
