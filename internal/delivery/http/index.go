package http

import (
	"net/http"
	"os"

	"github.com/econominhas/authentication/internal/models"
)

type HttpDelivery struct {
	server *http.Server
	router *http.ServeMux

	accountService *models.AccountService
}

type NewHttpDeliveryInput struct {
	AccountService models.AccountService
}

func (dlv *HttpDelivery) Listen() {
	dlv.auth()

	dlv.server.ListenAndServe()
}

func NewHttpDelivery(i *NewHttpDeliveryInput) *HttpDelivery {
	router := http.NewServeMux()

	server := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: router,
	}

	return &HttpDelivery{
		server: server,
		router: router,

		accountService: &i.AccountService,
	}
}
