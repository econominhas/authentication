package google

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

func (adp *GoogleAdapter) ExchangeCode(i *adapters.ExchangeCodeInput) (*adapters.ExchangeCodeOutput, error) {
	// ALERT: The order of the properties is important, don't change it!
	body := url.Values{}
	body.Set("code", i.Code)
	body.Set("client_id", os.Getenv("GOOGLE_CLIENT_ID"))
	body.Set("client_secret", os.Getenv("GOOGLE_CLIENT_SECRET"))
	if i.OriginUrl != "" {
		body.Set("redirect_uri", i.OriginUrl)
	}
	body.Set("grant_type", "authorization_code")
	// ALERT: The order of the properties is important, don't change it!

	req, err := http.NewRequest(
		http.MethodPost,
		"https://oauth2.googleapis.com/token",
		strings.NewReader(body.Encode()),
	)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil, errors.New("fail to build request")
	}

	codeRes, err := adp.httpClient.Do(req)
	if err != nil {
		return nil, errors.New("fail to make request")
	}
	defer codeRes.Body.Close()

	exchangeCode := exchangeTokenApiOutput{}
	err = json.NewDecoder(codeRes.Body).Decode(&exchangeCode)
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

	output := adapters.ExchangeCodeOutput{
		AccessToken:  exchangeCode.AccessToken,
		RefreshToken: exchangeCode.RefreshToken,
		Scopes:       strings.Split(exchangeCode.Scope, " "),
		ExpiresAt:    expDate,
	}

	return &output, nil
}
