package repositories

import (
	"errors"
	"time"

	"github.com/econominhas/authentication/internal/adapters"
	"github.com/econominhas/authentication/internal/models"
)

type RefreshTokenRepository struct {
	IdAdapter    adapters.IdAdapter
	TokenAdapter adapters.TokenAdapter
}

func (rep *RefreshTokenRepository) Create(i *models.CreateRefreshTokenInput) (*models.CreateRefreshTokenOutput, error) {
	refreshToken, err := rep.TokenAdapter.GenRefresh(&adapters.GenRefreshInput{
		AccountId: i.AccountId,
	})
	if err != nil {
		return nil, errors.New("fail to generate refresh token")
	}

	_, err = i.Db.Exec(
		"INSERT INTO auth.refresh_tokens (account_id, refresh_token) VALUES ($1, $2)",
		i.AccountId,
		refreshToken,
	)
	if err != nil {
		return nil, errors.New("fail to save refresh token")
	}

	return &models.CreateRefreshTokenOutput{
		RefreshToken: refreshToken,
	}, nil
}

func (rep *RefreshTokenRepository) Get(i *models.GetRefreshTokenInput) (*models.GetRefreshTokenOutput, error) {
	row := i.Db.QueryRow(
		"SELECT * FROM auth.refresh_tokens WHERE refresh_token = $1",
		i.RefreshToken,
	)

	var accountId string
	var createdAt time.Time

	err := row.Scan(
		&accountId,
		&createdAt,
	)

	if err != nil {
		return nil, errors.New("refresh token not found")
	}

	return &models.GetRefreshTokenOutput{
		AccountId: accountId,
		CreatedAt: createdAt,
	}, nil
}
