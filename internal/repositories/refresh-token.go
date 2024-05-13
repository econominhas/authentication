package repositories

import (
	"errors"

	"github.com/econominhas/authentication/internal/adapters"
	"github.com/econominhas/authentication/internal/models"
)

type RefreshTokenRepository struct {
	IdAdapter     adapters.IdAdapter
	TokenAdapter  adapters.TokenAdapter
	SecretAdapter adapters.SecretAdapter
}

func (rep *RefreshTokenRepository) Create(i *models.CreateRefreshTokenInput) (*models.CreateRefreshTokenOutput, error) {
	refreshToken, err := rep.SecretAdapter.GenSecret(64)
	if err != nil {
		return nil, errors.New("fail to generate secret")
	}

	_, err = i.Db.Exec(
		`
		INSERT INTO auth.refresh_tokens (account_id, refresh_token)
		VALUES ($1, $2)
		`,
		i.AccountId,
		refreshToken,
	)
	if err != nil {
		return nil, errors.New("fail to create refresh token")
	}

	return &models.CreateRefreshTokenOutput{
		RefreshToken: refreshToken,
	}, nil
}
