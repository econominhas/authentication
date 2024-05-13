package ulid

import (
	"errors"
	"os"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/econominhas/authentication/internal/adapters"
)

type TokenAdapter struct{}

func (adp *TokenAdapter) GenAccess(i *adapters.GenAccessInput) (*adapters.GenAccessOutput, error) {
	secretKey, err := paseto.NewV4AsymmetricSecretKeyFromHex(
		os.Getenv("PASETO_PRIVATE_KEY"),
	)
	if err != nil {
		return nil, errors.New("fail to get paseto private key")
	}

	expiresAt := time.Now().Add(15 * time.Minute)

	token := paseto.NewToken()
	token.SetSubject(i.AccountId)
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(expiresAt)

	accessToken := token.V4Sign(secretKey, nil)

	return &adapters.GenAccessOutput{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
	}, nil
}

func NewTokenAdapter() *TokenAdapter {
	return &TokenAdapter{}
}
