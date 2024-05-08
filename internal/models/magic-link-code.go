package models

import "time"

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

type UpsertInput struct {
	AccountId     string
	IsFirstAccess bool
}

type GetInput struct {
	AccountId string
	Code      string
}

type MagicLinkCodeRepository interface {
	Upsert(i *UpsertInput) (*MagicLinkCodeEntity, error)
	Get(i *GetInput) (*MagicLinkCodeEntity, error)
}
