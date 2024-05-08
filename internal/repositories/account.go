package repositories

import (
	"errors"
	"time"

	"github.com/econominhas/authentication/internal/adapters"
	"github.com/econominhas/authentication/internal/models"
)

type AccountRepository struct {
	IdAdapter adapters.IdAdapter
}

func (rep *AccountRepository) Create(i *models.CreateAccountInput) (*models.CreateAccountOutput, error) {
	if i.Email == "" && i.Phone == "" {
		return nil, errors.New("email or phone is required")
	}

	accountId, err := rep.IdAdapter.GenId()
	if err != nil {
		return nil, errors.New("fail to generate id")
	}

	_, err = i.Db.Exec(
		"INSERT INTO auth.accounts (id) VALUES ($1)",
		accountId,
	)
	if err != nil {
		return nil, errors.New("fail to create account")
	}

	if i.Email != "" {
		_, err = i.Db.Exec(
			"INSERT INTO auth.contacts (account_id, contact, `type`, verified_at) VALUES ($1, $2, $3, $4, $5)",
			accountId,
			i.Email,
			"EMAIL",
			time.Now(),
		)
		if err != nil {
			return nil, errors.New("fail to create email contact")
		}
	}

	if i.Phone != "" {
		_, err = i.Db.Exec(
			"INSERT INTO auth.contacts (account_id, contact, `type`, verified_at) VALUES ($1, $2, $3, $4, $5)",
			accountId,
			i.Phone,
			"PHONE",
			time.Now(),
		)
		if err != nil {
			return nil, errors.New("fail to create phone contact")
		}
	}

	if len(i.SignInProviders) > 0 {
		for _, signInProvider := range i.SignInProviders {
			_, err = i.Db.Exec(
				"INSERT INTO auth.sign_in_providers (account_id, provider, provider_id, access_token, refresh_token, expires_at) VALUES ($1, $2, $3, $4, $5, $6)",
				accountId,
				signInProvider.Type,
				signInProvider.Id,
				signInProvider.AccessToken,
				&signInProvider.RefreshToken,
				signInProvider.ExpiresAt,
			)
			if err != nil {
				return nil, errors.New("fail to create sign in provider")
			}
		}
	}

	return &models.CreateAccountOutput{
		Id: accountId,
	}, nil
}

func (rep *AccountRepository) GetManyByProvider(i *models.GetManyAccountsByProviderInput) ([]models.GetManyAccountsByProviderOutput, error) {
	rows, err := i.Db.Query(
		`
		SELECT
			COALESCE(sp.account_id, c.account_id) as "AccountId"
			sp.provider_id as "ProviderId"
			sp.provider as "ProviderType"
			c.contact as "Email"
		FROM auth.sign_in_providers sp
			INNER JOIN auth.contacts c ON c.account_id = sp.account_id
		WHERE
			(sp.provider = $1 AND sp.provider_id = $2)
			OR
			(c.type = $3 AND c.contact = $4)
		`,
		i.ProviderType,
		i.ProviderId,
		"EMAIL",
		i.Email,
	)
	if err != nil {
		return nil, errors.New("fail to find related accounts")
	}

	var output []models.GetManyAccountsByProviderOutput
	for rows.Next() {
		var row models.GetManyAccountsByProviderOutput

		if err := rows.Scan(
			&row.AccountId,
			&row.ProviderId,
			&row.ProviderType,
			&row.Email,
		); err != nil {
			return nil, errors.New("fail to find related accounts")
		}

		output = append(output, row)
	}

	return output, nil
}
