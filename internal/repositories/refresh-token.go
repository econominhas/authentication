package repositories

import (
	"errors"

	"github.com/econominhas/authentication/internal/adapters"
	"github.com/econominhas/authentication/internal/models"
)

type RefreshTokenRepository struct {
	IdAdapter     adapters.IdAdapter
	SecretAdapter adapters.SecretAdapter
	TokenAdapter  adapters.TokenAdapter
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

func (rep *RefreshTokenRepository) Get(i *models.GetRefreshTokenInput) (bool, error) {
	var data models.RefreshTokenEntity

	if err := i.Db.QueryRow(
		`
		SELECT *
		FROM auth.refresh_tokens rf
		WHERE
			rf.account_id = $1,
			rf.refresh_token = $2
		LIMIT 1
	`,
		i.AccountId,
		i.RefreshToken,
	).Scan(&data); err != nil {
		return false, errors.New("fail to get refresh token")
	}

	if data.AccountId == "" {
		return false, nil
	}

	return true, nil
}
