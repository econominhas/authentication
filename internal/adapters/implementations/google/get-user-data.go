package google

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/econominhas/authentication/internal/adapters"
)

type getUserDataApiOutput struct {
	Sub           string `json:"sub"`
	GivenName     string `json:"given_name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

func (adp *GoogleAdapter) GetUserData(accessToken string) (*adapters.GetAuthenticatedUserDataOutput, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://openidconnect.googleapis.com/v1/userinfo",
		nil,
	)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	if err != nil {
		return nil, errors.New("fail to build request")
	}

	userDataRes, err := adp.httpClient.Do(req)
	if err != nil {
		return nil, errors.New("fail to make request")
	}
	defer userDataRes.Body.Close()

	userData := getUserDataApiOutput{}
	err = json.NewDecoder(userDataRes.Body).Decode(&userData)
	if err != nil {
		return nil, errors.New("fail to decode request body")
	}

	return &adapters.GetAuthenticatedUserDataOutput{
		Id:              userData.Sub,
		Name:            userData.GivenName,
		Email:           userData.Email,
		IsEmailVerified: userData.EmailVerified,
	}, nil
}
