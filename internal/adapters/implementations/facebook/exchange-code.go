package facebook

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

type exchangeCodeApiOutput struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func (adp *FacebookAdapter) ExchangeCode(i *adapters.ExchangeCodeInput) (*adapters.ExchangeCodeOutput, error) {
	// ALERT: The order of the properties is important, don't change it!
	body := url.Values{}
	body.Set("code", i.Code)
	body.Set("client_id", os.Getenv("FACEBOOK_CLIENT_ID"))
	body.Set("client_secret", os.Getenv("FACEBOOK_CLIENT_SECRET"))
	if i.OriginUrl != "" {
		body.Set("redirect_uri", i.OriginUrl)
	}
	body.Set("grant_type", "authorization_code")
	// ALERT: The order of the properties is important, don't change it!

	request, err := http.NewRequest(
		http.MethodPost,
		"https://graph.facebook.com/v19.0/oauth/access_token",
		strings.NewReader(body.Encode()),
	)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil, errors.New("fail to build request")
	}

	res, err := adp.httpClient.Do(request)
	if err != nil {
		return nil, errors.New("fail to make request")
	}
	defer res.Body.Close()

	exchangeCode := exchangeCodeApiOutput{}
	err = json.NewDecoder(res.Body).Decode(&exchangeCode)
	if err != nil {
		return nil, errors.New("fail to decode request body")
	}

	expDate := time.
		Now().
		Add(
			time.Duration(
				exchangeCode.ExpiresIn,
			),
		)

	tokenDebug, err := adp.getTokenData(exchangeCode.AccessToken)
	if err != nil {
		return nil, errors.New("fail to get app token")
	}

	return &adapters.ExchangeCodeOutput{
		AccessToken: exchangeCode.AccessToken,
		Scopes:      tokenDebug.Scopes,
		ExpiresAt:   expDate,
	}, nil
}
