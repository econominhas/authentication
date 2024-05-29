package services

import (
	"database/sql"
	"net/http"
	"sync"

	"github.com/econominhas/authentication/internal/adapters"
	"github.com/econominhas/authentication/internal/models"
	"github.com/econominhas/authentication/internal/utils"
)

type AccountService struct {
	Logger models.Logger

	GoogleAdapter   adapters.SignInProviderAdapter
	FacebookAdapter adapters.SignInProviderAdapter
	DiscordAdapter  adapters.SignInProviderAdapter
	TokenAdapter    adapters.TokenAdapter
	EmailAdapter    adapters.EmailAdapter
	SmsAdapter      adapters.SmsAdapter

	Db *sql.DB

	AccountRepository       models.AccountRepository
	RefreshTokenRepository  models.RefreshTokenRepository
	MagicLinkCodeRepository models.MagicLinkCodeRepository
}

type genAuthOutputInput struct {
	db *sql.Tx

	accountId string
	// Tells if it's the user's first access
	isFirstAccess bool
	// Tells if should generate and return refresh token
	refresh bool
	// Tells if the account is complete or partial
	isComplete bool
}

type createFromExternalProviderInput struct {
	db *sql.Tx

	providerService *adapters.SignInProviderAdapter
	providerType    models.ProviderType
	code            string
	originUrl       string
}

func (serv *AccountService) genAuthOutput(i *genAuthOutputInput) (*models.AuthOutput, *utils.HttpError) {
	var wg sync.WaitGroup
	var refreshToken *models.CreateRefreshTokenOutput
	var accessToken *adapters.GenAccessOutput
	var err error

	if i.refresh {
		wg.Add(1)
		defer wg.Done()
		go func() {
			refreshToken, err = serv.RefreshTokenRepository.Create(&models.CreateRefreshTokenInput{
				Db: i.db,

				AccountId: i.accountId,
			})
		}()
	}

	wg.Add(1)
	defer wg.Done()
	go func() {
		accessToken, err = serv.TokenAdapter.GenAccess(&adapters.GenAccessInput{
			AccountId:  i.accountId,
			IsComplete: i.isComplete,
		})
	}()

	wg.Wait()

	if err != nil {
		i.db.Rollback()
		return nil, &utils.HttpError{
			Message:    "fail to generate auth output",
			StatusCode: http.StatusInternalServerError,
		}
	}

	i.db.Commit()

	return &models.AuthOutput{
		AccessToken:  accessToken.AccessToken,
		ExpiresAt:    accessToken.ExpiresAt,
		RefreshToken: refreshToken.RefreshToken,
	}, nil
}

func (serv *AccountService) createFromExternal(i *createFromExternalProviderInput) (*models.AuthOutput, *utils.HttpError) {
	exchangeCode, err := (*i.providerService).ExchangeCode(&adapters.ExchangeCodeInput{
		Code:      i.code,
		OriginUrl: i.originUrl,
	})
	if err != nil {
		i.db.Rollback()
		return nil, &utils.HttpError{
			Message:    "fail to exchange code",
			StatusCode: http.StatusInternalServerError,
		}
	}

	hasRequiredScopes := (*i.providerService).HasRequiredScopes(exchangeCode.Scopes)
	if !hasRequiredScopes {
		i.db.Rollback()
		return nil, &utils.HttpError{
			Message:    "missing scopes",
			StatusCode: http.StatusBadRequest,
		}
	}

	providerData, err := (*i.providerService).GetUserData(exchangeCode.AccessToken)
	if err != nil {
		i.db.Rollback()
		return nil, &utils.HttpError{
			Message:    "fail to get external user data",
			StatusCode: http.StatusBadGateway,
		}
	}

	if !providerData.IsEmailVerified {
		i.db.Rollback()
		return nil, &utils.HttpError{
			Message:    "unverified email",
			StatusCode: http.StatusBadRequest,
		}
	}

	relatedAccounts, err := serv.AccountRepository.GetManyByProvider(&models.GetManyAccountsByProviderInput{
		Db: i.db,

		ProviderId:   providerData.Id,
		ProviderType: i.providerType,
		Email:        providerData.Email,
	})
	if err != nil {
		i.db.Rollback()
		return nil, &utils.HttpError{
			Message:    "fail to get related accounts",
			StatusCode: http.StatusInternalServerError,
		}
	}

	var accountId string
	var isFirstAccess bool

	if len(relatedAccounts) > 0 {
		var sameEmail *models.GetManyAccountsByProviderOutput = nil
		var sameProvider *models.GetManyAccountsByProviderOutput = nil
		for _, v := range relatedAccounts {
			if v.Email == providerData.Email {
				sameEmail = &v
			}
			if v.ProviderId == providerData.Id && v.ProviderType == i.providerType {
				sameProvider = &v
			}
			if sameEmail != nil && sameProvider != nil {
				break
			}
		}

		/*
		 * Has an account with the same email, and it
		 * isn't linked with another provider with the same type
		 */
		if sameEmail != nil && sameProvider == nil && sameEmail.ProviderType != i.providerType {
			accountId = sameEmail.AccountId
		}

		/*
		 * Account with same provider id (it can have a different email,
		 * in case that the user updated it in provider or on our platform)
		 * More descriptive IF:
		 * if ((sameProviderId && !sameEmail) || (sameProviderId && sameEmail)) {
		 */
		if sameProvider != nil {
			accountId = sameProvider.AccountId
		}

		if accountId == "" {
			i.db.Rollback()
			return nil, &utils.HttpError{
				Message:    "fail to relate account",
				StatusCode: http.StatusInternalServerError,
			}
		}

		/*
		 * Updates the account because it can be partially created
		 * or the user can add a new email, different sign in provider,
		 * so we need to add the extra missing information to make
		 * it complete
		 */
		err := serv.AccountRepository.Update(&models.UpdateAccountInput{
			Db: i.db,

			Email: providerData.Email,
			SignInProviders: []models.CreateAccountSignInProvider{
				{
					Id:           providerData.Id,
					Type:         i.providerType,
					AccessToken:  exchangeCode.AccessToken,
					RefreshToken: &exchangeCode.RefreshToken,
					ExpiresAt:    exchangeCode.ExpiresAt,
				},
			},
		})
		if err != nil {
			i.db.Rollback()
			return nil, &utils.HttpError{
				Message:    "fail to update account",
				StatusCode: http.StatusInternalServerError,
			}
		}
	} else {
		result, err := serv.AccountRepository.Create(&models.CreateAccountInput{
			Db: i.db,

			Email: providerData.Email,
			SignInProviders: []models.CreateAccountSignInProvider{
				{
					Id:           providerData.Id,
					Type:         i.providerType,
					AccessToken:  exchangeCode.AccessToken,
					RefreshToken: &exchangeCode.RefreshToken,
					ExpiresAt:    exchangeCode.ExpiresAt,
				},
			},
		})
		if err != nil {
			i.db.Rollback()
			return nil, &utils.HttpError{
				Message:    "fail to create account",
				StatusCode: http.StatusInternalServerError,
			}
		}

		accountId = result.Id
		isFirstAccess = true
	}

	return serv.genAuthOutput(&genAuthOutputInput{
		accountId:     accountId,
		isFirstAccess: isFirstAccess,
		refresh:       true,
		isComplete:    true,
	})
}

func (serv *AccountService) CreateFromFacebookProvider(i *models.CreateAccountFromExternalProviderInput) (*models.AuthOutput, *utils.HttpError) {
	tx, err := serv.Db.Begin()
	if err != nil {
		return nil, &utils.HttpError{
			Message:    "fail to create transaction",
			StatusCode: http.StatusInternalServerError,
		}
	}

	return serv.createFromExternal(&createFromExternalProviderInput{
		db: tx,

		providerService: &serv.FacebookAdapter,
		providerType:    models.ProviderTypeFacebookEnum,
		code:            i.Code,
		originUrl:       i.OriginUrl,
	})
}

func (serv *AccountService) CreateFromGoogleProvider(i *models.CreateAccountFromExternalProviderInput) (*models.AuthOutput, *utils.HttpError) {
	tx, err := serv.Db.Begin()
	if err != nil {
		return nil, &utils.HttpError{
			Message:    "fail to create transaction",
			StatusCode: http.StatusInternalServerError,
		}
	}

	return serv.createFromExternal(&createFromExternalProviderInput{
		db: tx,

		providerService: &serv.GoogleAdapter,
		providerType:    models.ProviderTypeGoogleEnum,
		code:            i.Code,
		originUrl:       i.OriginUrl,
	})
}

func (serv *AccountService) CreateFromDiscordProvider(i *models.CreateAccountFromExternalProviderInput) (*models.AuthOutput, *utils.HttpError) {
	tx, err := serv.Db.Begin()
	if err != nil {
		return nil, &utils.HttpError{
			Message:    "fail to create transaction",
			StatusCode: http.StatusInternalServerError,
		}
	}

	return serv.createFromExternal(&createFromExternalProviderInput{
		db: tx,

		providerService: &serv.DiscordAdapter,
		providerType:    models.ProviderTypeDiscordEnum,
		code:            i.Code,
		originUrl:       i.OriginUrl,
	})
}

func (serv *AccountService) CreateFromEmailProvider(i *models.CreateAccountFromEmailInput) *utils.HttpError {
	tx, err := serv.Db.Begin()
	if err != nil {
		return &utils.HttpError{
			Message:    "fail to create transaction",
			StatusCode: http.StatusInternalServerError,
		}
	}

	var accountId string
	var isFirstAccess bool

	existentAccount, err := serv.AccountRepository.GetByEmail(&models.GetAccountByEmailInput{
		Db: tx,

		Email: i.Email,
	})
	if err != nil {
		tx.Rollback()
		return &utils.HttpError{
			Message:    "fail to get account",
			StatusCode: http.StatusInternalServerError,
		}
	}

	if existentAccount == nil {
		createdAccount, err := serv.AccountRepository.Create(&models.CreateAccountInput{
			Db: tx,

			Email: i.Email,
		})
		if err != nil {
			tx.Rollback()
			return &utils.HttpError{
				Message:    "fail to create account",
				StatusCode: http.StatusInternalServerError,
			}
		}

		accountId = createdAccount.Id
		isFirstAccess = true
	} else {
		accountId = existentAccount.AccountId
	}

	magicLinkCode, err := serv.MagicLinkCodeRepository.Upsert(&models.UpsertMagicLinkRefreshTokenInput{
		Db: tx,

		AccountId:     accountId,
		IsFirstAccess: isFirstAccess,
	})
	if err != nil {
		tx.Rollback()
		return &utils.HttpError{
			Message:    "fail to create account",
			StatusCode: http.StatusInternalServerError,
		}
	}

	err = serv.EmailAdapter.SendVerificationCodeEmail(&adapters.SendVerificationCodeEmailInput{
		To:   i.Email,
		Code: magicLinkCode.Code,
	})
	if err != nil {
		tx.Rollback()
		return &utils.HttpError{
			Message:    "fail to send sms",
			StatusCode: http.StatusInternalServerError,
		}
	}

	tx.Commit()

	return nil
}

func (serv *AccountService) CreateFromPhoneProvider(i *models.CreateAccountFromPhoneInput) *utils.HttpError {
	tx, err := serv.Db.Begin()
	if err != nil {
		return &utils.HttpError{
			Message:    "fail to create transaction",
			StatusCode: http.StatusInternalServerError,
		}
	}

	var accountId string
	var isFirstAccess bool

	existentAccount, err := serv.AccountRepository.GetByPhone(&models.GetAccountByPhoneInput{
		Db: tx,

		CountryCode: i.Phone.CountryCode,
		Number:      i.Phone.Number,
	})
	if err != nil {
		tx.Rollback()
		return &utils.HttpError{
			Message:    "fail to get account",
			StatusCode: http.StatusInternalServerError,
		}
	}

	if existentAccount == nil {
		createdAccount, err := serv.AccountRepository.Create(&models.CreateAccountInput{
			Db: tx,

			Phone: &i.Phone,
		})
		if err != nil {
			tx.Rollback()
			return &utils.HttpError{
				Message:    "fail to create account",
				StatusCode: http.StatusInternalServerError,
			}
		}

		accountId = createdAccount.Id
		isFirstAccess = true
	} else {
		accountId = existentAccount.AccountId
	}

	magicLinkCode, err := serv.MagicLinkCodeRepository.Upsert(&models.UpsertMagicLinkRefreshTokenInput{
		Db: tx,

		AccountId:     accountId,
		IsFirstAccess: isFirstAccess,
	})
	if err != nil {
		tx.Rollback()
		return &utils.HttpError{
			Message:    "fail to create account",
			StatusCode: http.StatusInternalServerError,
		}
	}

	err = serv.SmsAdapter.SendVerificationCodeSms(&adapters.SendVerificationCodeSmsInput{
		To:   i.Phone.CountryCode + i.Phone.Number,
		Code: magicLinkCode.Code,
	})
	if err != nil {
		tx.Rollback()
		return &utils.HttpError{
			Message:    "fail to send sms",
			StatusCode: http.StatusInternalServerError,
		}
	}

	tx.Commit()

	return nil
}

func (serv *AccountService) PartialCreateFromDiscordId(i *models.PartialCreateFromDiscordIdInput) (*models.AuthOutput, *utils.HttpError) {
	tx, err := serv.Db.Begin()
	if err != nil {
		return nil, &utils.HttpError{
			Message:    "fail to create transaction",
			StatusCode: http.StatusInternalServerError,
		}
	}

	relatedAccounts, err := serv.AccountRepository.GetManyByProvider(&models.GetManyAccountsByProviderInput{
		Db: tx,

		ProviderId:   i.Id,
		ProviderType: models.ProviderTypeDiscordEnum,
	})
	if err != nil {
		tx.Rollback()
		return nil, &utils.HttpError{
			Message:    "fail to get related accounts",
			StatusCode: http.StatusInternalServerError,
		}
	}

	var accountId string
	var isFirstAccess bool

	if len(relatedAccounts) > 0 {
		var sameProvider *models.GetManyAccountsByProviderOutput = nil
		for _, v := range relatedAccounts {
			if v.ProviderId == i.Id && v.ProviderType == models.ProviderTypeDiscordEnum {
				sameProvider = &v
				break
			}
		}

		if sameProvider != nil {
			accountId = sameProvider.AccountId
		}

		if accountId == "" {
			tx.Rollback()
			return nil, &utils.HttpError{
				Message:    "fail to relate account",
				StatusCode: http.StatusInternalServerError,
			}
		}
	} else {
		result, err := serv.AccountRepository.Create(&models.CreateAccountInput{
			Db: tx,

			SignInProviders: []models.CreateAccountSignInProvider{
				{
					Id:   i.Id,
					Type: models.ProviderTypeDiscordEnum,
				},
			},
		})
		if err != nil {
			tx.Rollback()
			return nil, &utils.HttpError{
				Message:    "fail to create account",
				StatusCode: http.StatusInternalServerError,
			}
		}

		accountId = result.Id
		isFirstAccess = true
	}

	return serv.genAuthOutput(&genAuthOutputInput{
		accountId:     accountId,
		isFirstAccess: isFirstAccess,
		refresh:       false,
		isComplete:    false,
	})
}

func (serv *AccountService) ExchangeCode(i *models.ExchangeAccountCodeInput) (*models.AuthOutput, *utils.HttpError) {
	tx, err := serv.Db.Begin()
	if err != nil {
		return nil, &utils.HttpError{
			Message:    "fail to create transaction",
			StatusCode: http.StatusInternalServerError,
		}
	}

	magicLinkCode, err := serv.MagicLinkCodeRepository.Get(&models.GetMagicLinkRefreshTokenInput{
		Db: tx,

		AccountId: i.AccountId,
		Code:      i.Code,
	})
	if err != nil {
		tx.Rollback()
		return nil, &utils.HttpError{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	if magicLinkCode == nil {
		tx.Rollback()
		return nil, &utils.HttpError{
			Message:    "magic link code doesn't exist",
			StatusCode: http.StatusNotFound,
		}
	}

	tx.Commit()

	return serv.genAuthOutput(&genAuthOutputInput{
		db: tx,

		accountId:     i.AccountId,
		isFirstAccess: magicLinkCode.IsFirstAccess,
		refresh:       true,
		isComplete:    true,
	})
}

func (serv *AccountService) RefreshToken(i *models.RefreshAccountTokenInput) (*models.RefreshAccountTokenOutput, *utils.HttpError) {
	tx, err := serv.Db.Begin()
	if err != nil {
		return nil, &utils.HttpError{
			Message:    "fail to create transaction",
			StatusCode: http.StatusInternalServerError,
		}
	}

	refreshToken, err := serv.RefreshTokenRepository.Get(&models.GetRefreshTokenInput{
		Db: tx,

		AccountId:    i.AccountId,
		RefreshToken: i.RefreshToken,
	})
	if err != nil {
		tx.Rollback()
		return nil, &utils.HttpError{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	if !refreshToken {
		tx.Rollback()
		return nil, &utils.HttpError{
			Message:    "refresh token doesn't exist",
			StatusCode: http.StatusNotFound,
		}
	}

	accessToken, err := serv.TokenAdapter.GenAccess(&adapters.GenAccessInput{
		AccountId: i.AccountId,
	})
	if err != nil {
		tx.Rollback()
		return nil, &utils.HttpError{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	tx.Commit()

	return &models.RefreshAccountTokenOutput{
		AccessToken: accessToken.AccessToken,
		ExpiresAt:   accessToken.ExpiresAt,
	}, nil
}
