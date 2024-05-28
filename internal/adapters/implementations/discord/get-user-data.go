package discord

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/econominhas/authentication/internal/adapters"
)

type getUserDataApiOutput struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Banner        string `json:"banner"`
	Bot           bool   `json:"bot"`
	System        bool   `json:"system"`
	MfaEnabled    bool   `json:"mfa_enabled"`
	AccentColor   int    `json:"accent_color"`
	Locale        string `json:"locale"`
	Verified      bool   `json:"verified"`
	Email         string `json:"email"`
	Flags         int    `json:"flags"`
	PremiumType   int    `json:"premium_type"`
	PublicFlags   int    `json:"public_flags"`
}

func (adp *DiscordAdapter) GetUserData(accessToken string) (*adapters.GetAuthenticatedUserDataOutput, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://discord.com/api/v10/users/@me",
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
		Id:              userData.Id,
		Name:            userData.Username,
		Email:           userData.Email,
		IsEmailVerified: userData.Verified,
	}, nil
}
