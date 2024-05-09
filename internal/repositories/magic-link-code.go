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

func (rep *MagicLinkCodeRepository) Upsert(i *models.UpsertInput) (*models.MagicLinkCodeEntity, error) {
	_, err := i.Db.Exec(
		`
		INSERT INTO auth.magic_link_codes (account_id, code, is_first_access)
        VALUES ($1, $2, $3)
        ON CONFLICT (account_id) DO UPDATE SET code = $2, is_first_access = $3
    `,
		i.AccountId,
		rep.SecretAdapter.GenSecret(16),
		i.IsFirstAccess,
	)

	if err != nil {
		return nil, errors.New("fail to insert or create magic link code")
	}

	return &models.MagicLinkCodeEntity{
		AccountId:     i.AccountId,
		Code:          i.AccountId,
		IsFirstAccess: i.IsFirstAccess,
		CreatedAt:     time.Now(),
	}
}
