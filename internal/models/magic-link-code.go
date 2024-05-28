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
	AccountId     string    `db:"account_id" json:"accountId,omitempty"`
	Code          string    `db:"code" json:"code,omitempty"`
	IsFirstAccess bool      `db:"is_first_access" json:"isFirstAccess,omitempty"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt,omitempty"`
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
