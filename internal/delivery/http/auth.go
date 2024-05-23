package http

import (
	"encoding/json"
	"net/http"

	"github.com/econominhas/authentication/internal/models"
)

func (dlv *HttpDelivery) auth() {
	const prefix = "auth"

	dlv.router.HandleFunc("POST /"+prefix+"/google", func(w http.ResponseWriter, r *http.Request) {
		body := &models.CreateAccountFromExternalProviderInput{}
		err := json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := (*dlv.accountService).CreateFromGoogleProvider(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(result)
	})
}