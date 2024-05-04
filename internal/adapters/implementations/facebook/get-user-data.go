package facebook

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/econominhas/authentication/internal/adapters"
)

type getUserDataApiOutput struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (adp *FacebookAdapter) GetUserData(accessToken string) (*adapters.GetAuthenticatedUserDataOutput, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://graph.facebook.com/v19.0/me/"+
			"?fields=id,name,email"+
			"&access_token="+accessToken,
		nil,
	)
	req.Header.Add("Accept", "application/json")
	if err != nil {
		return nil, errors.New("fail to build request")
	}

	res, err := adp.httpClient.Do(req)
	if err != nil {
		return nil, errors.New("fail to make request")
	}
	defer res.Body.Close()

	userData := getUserDataApiOutput{}

	err = json.NewDecoder(res.Body).Decode(&userData)
	if err != nil {
		return nil, errors.New("fail to decode request body")
	}

	output := adapters.GetAuthenticatedUserDataOutput{
		Id:              userData.Id,
		Name:            userData.Name,
		Email:           userData.Email,
		IsEmailVerified: true,
	}

	return &output, nil
}
