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
	if i.Email == "" && i.Phone.Number == "" {
		return nil, errors.New("email or phone is required")
	}
	if (i.Phone.CountryCode != "" && i.Phone.Number == "") ||
		(i.Phone.CountryCode == "" && i.Phone.Number != "") {
		return nil, errors.New("both phone number and country code are required")
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
			"INSERT INTO auth.email_addresses (account_id, email_address, verified_at) VALUES ($1, $2, $3)",
			accountId,
			i.Email,
			time.Now(),
		)
		if err != nil {
			return nil, errors.New("fail to create email contact")
		}
	}

	if i.Phone.Number != "" {
		_, err = i.Db.Exec(
			"INSERT INTO auth.phone_numbers (account_id, country_code, phone_number, verified_at) VALUES ($1, $2, $3, $4)",
			accountId,
			i.Phone.CountryCode,
			i.Phone.Number,
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
			COALESCE(sp.account_id, ea.account_id) as "AccountId"
			sp.provider_id as "ProviderId"
			sp.provider as "ProviderType"
			ea.email_address as "Email"
		FROM auth.sign_in_providers sp
			INNER JOIN auth.email_addresses ea ON ea.account_id = sp.account_id
		WHERE
			(sp.provider = $1 AND sp.provider_id = $2)
			OR
			ea.email_address = $3
		`,
		i.ProviderType,
		i.ProviderId,
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

func (rep *AccountRepository) GetByEmail(i *models.GetAccountByEmailInput) (*models.GetAccountByEmailOutput, error) {
	var data models.GetAccountByEmailOutput

	if err := i.Db.QueryRow(
		`
		SELECT
			ea.account_id as "AccountId"
			ea.email_address as "Email"
		FROM auth.email_addresses ea
		WHERE
			ea.email_address = $1
		LIMIT 1
		`,
		i.Email,
	).Scan(data); err != nil {
		return nil, errors.New("fail to get account by email")
	}

	return &data, nil
}

func (rep *AccountRepository) GetByPhone(i *models.GetAccountByPhoneInput) (*models.GetAccountByPhoneOutput, error) {
	var data models.GetAccountByPhoneOutput

	if err := i.Db.QueryRow(
		`
		SELECT
			pn.account_id as "AccountId"
			pn.country_code as "CountryCode"
			pn.phone_number as "Phone"
		FROM auth.phone_numbers pn
		WHERE
			pn.country_code = $1
			pn.phone_number = $2
		LIMIT 1
		`,
		i.CountryCode,
		i.Number,
	).Scan(data); err != nil {
		return nil, errors.New("fail to get account by phone")
	}

	return &data, nil
}
