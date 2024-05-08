package facebook

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type getAppTokenApiOutput struct {
	AccessToken string `json:"access_token"`
}

type tokenDebugApiOutput struct {
	Data struct {
		AppId       int    `json:"app_id"`
		Type        string `json:"type"`
		Application string `json:"application"`
		ExpiresAt   int    `json:"expires_at"`
		IsValid     bool   `json:"is_valid"`
		IssuedAt    int    `json:"issued_at"`
		Metadata    struct {
			Sso string `json:"sso"`
		} `json:"metadata"`
		Scopes []string `json:"scopes"`
		UserId string   `json:"user_id"`
	} `json:"data"`
}

type getTokenDataOutput struct {
	AppId       int
	Type        string
	Application string
	ExpiresAt   int
	IsValid     bool
	IssuedAt    int
	Scopes      []string
	UserId      string
	AppToken    string
}

func (adp *FacebookAdapter) getTokenData(accessToken string) (*getTokenDataOutput, error) {
	appTokenRequest, err := http.NewRequest(
		http.MethodGet,
		"https://graph.facebook.com/oauth/access_token"+
			"?client_id="+os.Getenv("FACEBOOK_CLIENT_ID")+
			"&client_secret="+os.Getenv("FACEBOOK_CLIENT_SECRET")+
			"&grant_type=client_credentials",
		nil,
	)
	appTokenRequest.Header.Add("Accept", "application/json")
	if err != nil {
		return nil, errors.New("fail to build request")
	}

	appTokenResp, err := adp.httpClient.Do(appTokenRequest)
	if err != nil {
		return nil, errors.New("fail to make request")
	}
	defer appTokenResp.Body.Close()

	appToken := getAppTokenApiOutput{}
	err = json.NewDecoder(appTokenResp.Body).Decode(&appToken)
	if err != nil {
		return nil, errors.New("fail to decode request body")
	}

	// Token debug

	tokenDebugRequest, err := http.NewRequest(
		http.MethodGet,
		"https://graph.facebook.com/debug_token"+
			"?input_token="+accessToken+
			"&access_token="+appToken.AccessToken,
		nil,
	)
	tokenDebugRequest.Header.Add("Accept", "application/json")
	if err != nil {
		return nil, errors.New("fail to build request")
	}

	tokenDebugRes, err := adp.httpClient.Do(tokenDebugRequest)
	if err != nil {
		return nil, errors.New("fail to make request")
	}
	defer tokenDebugRes.Body.Close()

	tokenDebug := tokenDebugApiOutput{}
	err = json.NewDecoder(tokenDebugRes.Body).Decode(&tokenDebug)
	if err != nil {
		return nil, errors.New("fail to decode request body")
	}

	// Return

	return &getTokenDataOutput{
		AppId:       tokenDebug.Data.AppId,
		Type:        tokenDebug.Data.Type,
		Application: tokenDebug.Data.Application,
		ExpiresAt:   tokenDebug.Data.ExpiresAt,
		IsValid:     tokenDebug.Data.IsValid,
		IssuedAt:    tokenDebug.Data.IssuedAt,
		Scopes:      tokenDebug.Data.Scopes,
		UserId:      tokenDebug.Data.UserId,
		AppToken:    appToken.AccessToken,
	}, nil
}
