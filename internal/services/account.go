package services

import (
	"errors"
	"sync"

	"github.com/econominhas/authentication/internal/adapters"
	"github.com/econominhas/authentication/internal/models"
)

type AccountService struct {
	GoogleAdapter   adapters.SignInProviderAdapter
	FacebookAdapter adapters.SignInProviderAdapter
	TokenAdapter    adapters.TokenAdapter

	AccountRepository      models.AccountRepository
	RefreshTokenRepository models.RefreshTokenRepository
}

type genAuthOutputInput struct {
	accountId     string
	isFirstAccess bool
	refresh       bool
}

type createFromExternalProviderInput struct {
	providerService adapters.SignInProviderAdapter
	providerType    string
	code            string
	originUrl       string
}

func (serv *AccountService) genAuthOutput(i *genAuthOutputInput) (*models.AuthOutput, error) {
	var wg sync.WaitGroup
	var refreshToken *models.CreateRefreshTokenOutput
	var accessToken *adapters.GenAccessOutput
	var err error

	if i.refresh {
		wg.Add(1)
		defer wg.Done()
		go func() {
			refreshToken, err = serv.RefreshTokenRepository.Create(&models.CreateRefreshTokenInput{
				AccountId: i.accountId,
			})
		}()
	}

	wg.Add(1)
	defer wg.Done()
	go func() {
		accessToken, err = serv.TokenAdapter.GenAccess(&adapters.GenAccessInput{
			AccountId: i.accountId,
		})
	}()

	wg.Wait()

	if err != nil {
		return nil, errors.New("fail to generate auth output")
	}

	return &models.AuthOutput{
		AccessToken:  accessToken.AccessToken,
		ExpiresAt:    accessToken.ExpiresAt,
		RefreshToken: refreshToken.RefreshToken,
	}, nil
}

func (serv *AccountService) createFromExternal(i *createFromExternalProviderInput) (*models.AuthOutput, error) {
	exchangeCode, err := i.providerService.ExchangeCode(&adapters.ExchangeCodeInput{
		Code:      i.code,
		OriginUrl: i.originUrl,
	})
	if err != nil {
		return nil, errors.New("fail to exchange code")
	}

	hasRequiredScopes := i.providerService.HasRequiredScopes(exchangeCode.Scopes)
	if !hasRequiredScopes {
		return nil, errors.New("missing scopes")
	}

	providerData, err := i.providerService.GetUserData(exchangeCode.AccessToken)
	if err != nil {
		return nil, errors.New("fail to get external user data")
	}

	if !providerData.IsEmailVerified {
		return nil, errors.New("unverified email")
	}

	relatedAccounts, err := serv.AccountRepository.GetManyByProvider(&models.GetManyAccountsByProviderInput{
		ProviderId:   providerData.Id,
		ProviderType: i.providerType,
		Email:        providerData.Email,
	})
	if err != nil {
		return nil, errors.New("fail to get related accounts")
	}

	var accountId string
	var isFirstAccess bool

	if len(relatedAccounts) > 0 {
		sameEmail := new(models.GetManyAccountsByProviderOutput)
		sameProvider := new(models.GetManyAccountsByProviderOutput)
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
			return nil, errors.New("fail to relate account")
		}
	} else {
		result, err := serv.AccountRepository.Create(&models.CreateAccountInput{
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
			return nil, errors.New("fail to create account")
		}

		accountId = result.Id
		isFirstAccess = true
	}

	return serv.genAuthOutput(&genAuthOutputInput{
		accountId:     accountId,
		isFirstAccess: isFirstAccess,
		refresh:       true,
	})
}

func (serv *AccountService) CreateFromGoogleProvider(i *models.CreateAccountFromExternalProviderInput) (*models.AuthOutput, error) {
	return serv.createFromExternal(&createFromExternalProviderInput{
		providerService: serv.GoogleAdapter,
		providerType:    "GOOGLE",
		code:            i.Code,
		originUrl:       i.OriginUrl,
	})
}

func (serv *AccountService) CreateFromFacebookProvider(i *models.CreateAccountFromExternalProviderInput) (*models.AuthOutput, error) {
	return serv.createFromExternal(&createFromExternalProviderInput{
		providerService: serv.FacebookAdapter,
		providerType:    "FACEBOOK",
		code:            i.Code,
		originUrl:       i.OriginUrl,
	})
}
