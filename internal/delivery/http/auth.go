package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/econominhas/authentication/internal/delivery"
	"github.com/econominhas/authentication/internal/models"
	"github.com/econominhas/authentication/internal/utils"
)

type AuthController struct {
	prefix string

	router    *http.ServeMux
	validator delivery.Validator

	accountService models.AccountService
}

func (c *AuthController) CreateFromEmailProvider() {
	if os.Getenv("FF_CREATE_FROM_EMAIL") != "1" {
		return
	}

	route := fmt.Sprintf("POST %s/email", c.prefix)

	c.router.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		body := &models.CreateAccountFromEmailInput{}
		err := json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = c.validator.Validate(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = c.accountService.CreateFromEmailProvider(body)
		if err != nil {
			http.Error(w, err.Error(), err.(*utils.HttpError).HttpStatusCode())
			return
		}
	})
}

func (c *AuthController) CreateFromPhoneProvider() {
	if os.Getenv("FF_CREATE_FROM_PHONE") != "1" {
		return
	}

	route := fmt.Sprintf("POST %s/phone", c.prefix)

	c.router.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		body := &models.CreateAccountFromPhoneInput{}
		err := json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = c.validator.Validate(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = c.accountService.CreateFromPhoneProvider(body)
		if err != nil {
			http.Error(w, err.Error(), err.(*utils.HttpError).HttpStatusCode())
			return
		}
	})
}

func (c *AuthController) CreateFromFacebookProvider() {
	if os.Getenv("FF_CREATE_FROM_FACEBOOK") != "1" {
		return
	}

	route := fmt.Sprintf("POST %s/facebook", c.prefix)

	c.router.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		body := &models.CreateAccountFromExternalProviderInput{}
		err := json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = c.validator.Validate(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := c.accountService.CreateFromFacebookProvider(body)
		if err != nil {
			http.Error(w, err.Error(), err.(*utils.HttpError).HttpStatusCode())
			return
		}

		json.NewEncoder(w).Encode(result)
	})
}

func (c *AuthController) CreateFromGoogleProvider() {
	if os.Getenv("FF_CREATE_FROM_GOOGLE") != "1" {
		return
	}

	route := fmt.Sprintf("POST %s/google", c.prefix)

	c.router.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		body := &models.CreateAccountFromExternalProviderInput{}
		err := json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = c.validator.Validate(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := c.accountService.CreateFromGoogleProvider(body)
		if err != nil {
			http.Error(w, err.Error(), err.(*utils.HttpError).HttpStatusCode())
			return
		}

		json.NewEncoder(w).Encode(result)
	})
}

func (c *AuthController) CreateFromDiscordProvider() {
	if os.Getenv("FF_CREATE_FROM_DISCORD") != "1" {
		return
	}

	route := fmt.Sprintf("POST %s/discord", c.prefix)

	c.router.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		body := &models.CreateAccountFromExternalProviderInput{}
		err := json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = c.validator.Validate(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := c.accountService.CreateFromDiscordProvider(body)
		if err != nil {
			http.Error(w, err.Error(), err.(*utils.HttpError).HttpStatusCode())
			return
		}

		json.NewEncoder(w).Encode(result)
	})
}

func (c *AuthController) PartialCreateFromDiscordId() {
	if os.Getenv("FF_PARTIAL_CREATE_FROM_DISCORD") != "1" {
		return
	}

	route := fmt.Sprintf("POST %s/partial/discord", c.prefix)

	c.router.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		body := &models.PartialCreateFromDiscordIdInput{}
		err := json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = c.validator.Validate(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := c.accountService.PartialCreateFromDiscordId(body)
		if err != nil {
			http.Error(w, err.Error(), err.(*utils.HttpError).HttpStatusCode())
			return
		}

		json.NewEncoder(w).Encode(result)
	})
}

func (c *AuthController) ExchangeCode() {
	// If none of the passwordless methods are allowed,
	// then this feature should not be allowed either
	if os.Getenv("FF_CREATE_FROM_EMAIL") != "1" &&
		os.Getenv("FF_CREATE_FROM_PHONE") != "1" {
		return
	}

	route := fmt.Sprintf("POST %s/code", c.prefix)

	c.router.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		body := &models.ExchangeAccountCodeInput{}
		err := json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = c.validator.Validate(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := c.accountService.ExchangeCode(body)
		if err != nil {
			http.Error(w, err.Error(), err.(*utils.HttpError).HttpStatusCode())
			return
		}

		json.NewEncoder(w).Encode(result)
	})
}

func (c *AuthController) RefreshToken() {
	route := fmt.Sprintf("POST %s/refresh", c.prefix)

	c.router.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		body := &models.RefreshAccountTokenInput{}
		err := json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = c.validator.Validate(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := c.accountService.RefreshToken(body)
		if err != nil {
			http.Error(w, err.Error(), err.(*utils.HttpError).HttpStatusCode())
			return
		}

		json.NewEncoder(w).Encode(result)
	})
}

func (dlv *HttpDelivery) AuthController() {
	controller := &AuthController{
		prefix:         "/auth",
		router:         dlv.router,
		validator:      dlv.validator,
		accountService: dlv.accountService,
	}

	controller.CreateFromEmailProvider()
	controller.CreateFromPhoneProvider()
	controller.CreateFromFacebookProvider()
	controller.CreateFromGoogleProvider()
	controller.CreateFromDiscordProvider()
	controller.PartialCreateFromDiscordId()
	controller.ExchangeCode()
	controller.RefreshToken()
}
