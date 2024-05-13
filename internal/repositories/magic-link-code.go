package repositories

import (
	"errors"
	"time"

	"github.com/econominhas/authentication/internal/adapters"
	"github.com/econominhas/authentication/internal/models"
)

type MagicLinkCodeRepository struct {
	SecretAdapter adapters.SecretAdapter
}

func (rep *MagicLinkCodeRepository) Upsert(i *models.UpsertMagicLinkRefreshTokenInput) (*models.MagicLinkCodeEntity, error) {
	secret, err := rep.SecretAdapter.GenSecret(16)

	if err != nil {
		return nil, errors.New("fail to generate secret")
	}

	createdAt := time.Now()

	_, err = i.Db.Exec(
		`
		INSERT INTO auth.magic_link_codes (account_id, code, is_first_access)
		VALUES ($1, $2, $3)
		ON CONFLICT (account_id)
		DO UPDATE SET
			code = $2,
			is_first_access = $3,
			created_at = $4
    `,
		i.AccountId,
		secret,
		i.IsFirstAccess,
		createdAt,
	)

	if err != nil {
		return nil, errors.New("fail to insert or create magic link code")
	}

	return &models.MagicLinkCodeEntity{
		AccountId:     i.AccountId,
		Code:          i.AccountId,
		IsFirstAccess: i.IsFirstAccess,
		CreatedAt:     createdAt,
	}, nil
}

func (rep *MagicLinkCodeRepository) Get(i *models.GetMagicLinkRefreshTokenInput) (*models.MagicLinkCodeEntity, error) {
	var data models.MagicLinkCodeEntity

	if err := i.Db.QueryRow(
		`
		SELECT *
		FROM auth.magic_link_codes mlc
		WHERE
			mlc.account_id = $1,
			mlc.code = $2
		LIMIT 1
	`,
		i.AccountId,
		i.Code,
	).Scan(&data); err != nil {
		return nil, errors.New("fail to get magic link code")
	}

	return &data, nil
}
