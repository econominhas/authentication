package discord

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/econominhas/authentication/internal/adapters"
)

type exchangeTokenApiOutput struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

func (adp *DiscordAdapter) ExchangeCode(i *adapters.ExchangeCodeInput) (*adapters.ExchangeCodeOutput, error) {
	adp.logger.Info("Start ExchangeCode")

	adp.logger.Debug("Building exchange code body")

	// ALERT: The order of the properties is important, don't change it!
	body := url.Values{}
	body.Set("code", i.Code)
	body.Set("client_id", os.Getenv("DISCORD_CLIENT_ID"))
	body.Set("client_secret", os.Getenv("DISCORD_CLIENT_SECRET"))
	if i.OriginUrl != "" {
		body.Set("redirect_uri", i.OriginUrl)
	}
	body.Set("grant_type", "authorization_code")
	// ALERT: The order of the properties is important, don't change it!

	adp.logger.Debug("Exchange code  built")

	adp.logger.Debug("Building request to exchange code")

	req, err := http.NewRequest(
		http.MethodPost,
		"https://discord.com/api/v10/oauth2/token",
		strings.NewReader(body.Encode()),
	)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		adp.logger.Error(
			"Fail to build request to exchange code",
			"message", err.Error(),
		)

		return nil, errors.New("fail to build request")
	}

	adp.logger.Debug("Doing request to exchange code")

	codeRes, err := adp.httpClient.Do(req)
	if err != nil {
		adp.logger.Error(
			"Fail to make request to exchange code",
			"message", err.Error(),
		)
		return nil, errors.New("fail to make request")
	}
	defer codeRes.Body.Close()

	adp.logger.Debug("Request to exchange code done")

	adp.logger.Debug("Try to decode response body")

	exchangeCode := exchangeTokenApiOutput{}
	err = json.NewDecoder(codeRes.Body).Decode(&exchangeCode)
	if err != nil {
		adp.logger.Error(
			"Fail to make decode response body",
			"message", err.Error(),
		)

		return nil, errors.New("fail to decode request body")
	}

	adp.logger.Debug("Response body decoded")

	expDate := time.
		Now().
		Add(
			time.Duration(
				exchangeCode.ExpiresIn,
			),
		)

	output := adapters.ExchangeCodeOutput{
		AccessToken:  exchangeCode.AccessToken,
		RefreshToken: exchangeCode.RefreshToken,
		Scopes:       strings.Split(exchangeCode.Scope, " "),
		ExpiresAt:    expDate,
	}

	adp.logger.Info("Successfully finish ExchangeCode")

	return &output, nil
}
