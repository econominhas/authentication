package adapters

import "time"

type GenAccessInput struct {
	AccountId string
}

type GenAccessOutput struct {
	AccessToken string
	ExpiresAt   time.Time
}

type GenRefreshInput struct {
	AccountId string
}

type TokenAdapter interface {
	GenAccess(i *GenAccessInput) (*GenAccessOutput, error)

	GenRefresh(i *GenRefreshInput) (string, error)
}
