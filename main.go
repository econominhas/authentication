package main

import (
	"net/http"
	"os"

	"github.com/econominhas/authentication/internal/adapters/implementations/facebook"
	"github.com/econominhas/authentication/internal/adapters/implementations/google"
	"github.com/econominhas/authentication/internal/services"
)

func main() {
	router := http.NewServeMux()

	server := http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: router,
	}

	googleAdapter := google.NewGoogle()
	facebookAdapter := facebook.NewFacebook()

	accountService := &services.AccountService{
		GoogleAdapter:   googleAdapter,
		FacebookAdapter: facebookAdapter,
	}

	server.ListenAndServe()
}
