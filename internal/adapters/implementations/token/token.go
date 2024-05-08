package ulid

import (
	"errors"
	"math/rand"
	"os"
	"time"
	"unsafe"

	"aidanwoods.dev/go-paseto"
	"github.com/econominhas/authentication/internal/adapters"
)

const (
	tokenLength   = 64
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
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

func (adp *TokenAdapter) GenRefresh(i *adapters.GenRefreshInput) (string, error) {
	var src = rand.NewSource(time.Now().UnixNano())

	b := make([]byte, tokenLength)

	for i, cache, remain := tokenLength-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b)), nil
}

func NewTokenAdapter() *TokenAdapter {
	return &TokenAdapter{}
}
