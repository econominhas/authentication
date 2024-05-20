package paseto

import (
	"errors"
	"os"
	"time"

	goPaseto "aidanwoods.dev/go-paseto"
	"github.com/econominhas/authentication/internal/adapters"
)

type PasetoAdapter struct{}

func (adp *PasetoAdapter) GenAccess(i *adapters.GenAccessInput) (*adapters.GenAccessOutput, error) {
	secretKey, err := goPaseto.NewV4AsymmetricSecretKeyFromHex(
		os.Getenv("PASETO_PRIVATE_KEY"),
	)
	if err != nil {
		return nil, errors.New("fail to get paseto private key")
	}

	expiresAt := time.Now().Add(15 * time.Minute)

	token := goPaseto.NewToken()
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

func NewPaseto() *PasetoAdapter {
	return &PasetoAdapter{}
}
