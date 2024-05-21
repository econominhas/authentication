package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/econominhas/authentication/internal/adapters/implementations/facebook"
	"github.com/econominhas/authentication/internal/adapters/implementations/google"
	"github.com/econominhas/authentication/internal/adapters/implementations/paseto"
	"github.com/econominhas/authentication/internal/adapters/implementations/secret"
	"github.com/econominhas/authentication/internal/adapters/implementations/ses"
	"github.com/econominhas/authentication/internal/adapters/implementations/sns"
	"github.com/econominhas/authentication/internal/adapters/implementations/ulid"
	"github.com/econominhas/authentication/internal/delivery/http"
	"github.com/econominhas/authentication/internal/repositories"
	"github.com/econominhas/authentication/internal/services"
	_ "github.com/lib/pq"
)

func main() {
	log.Print("Start")

	// ----------------------------
	//
	// Env
	//
	// ----------------------------

	validateEnvs()

	// ----------------------------
	//
	// Databases
	//
	// ----------------------------

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
		panic(1)
	}
	defer db.Close()
	log.Print("Connected to database")

	// ----------------------------
	//
	// Adapters
	//
	// ----------------------------

	googleAdapter := google.NewGoogle()
	facebookAdapter := facebook.NewFacebook()
	pasetoAdapter := paseto.NewPaseto()
	secretAdapter := secret.NewSecret()
	sesAdapter := ses.NewSes()
	snsAdapter := sns.NewSns()
	ulidAdapter := ulid.NewUlid()

	log.Print("Adapters initialized")

	// ----------------------------
	//
	// Repositories
	//
	// ----------------------------

	accountRepository := &repositories.AccountRepository{
		IdAdapter: ulidAdapter,
	}
	magicLinkCodeRepository := &repositories.MagicLinkCodeRepository{
		SecretAdapter: secretAdapter,
	}
	refreshTokenRepository := &repositories.RefreshTokenRepository{
		IdAdapter:     ulidAdapter,
		SecretAdapter: secretAdapter,
		TokenAdapter:  pasetoAdapter,
	}

	log.Print("Repositories initialized")

	// ----------------------------
	//
	// Services
	//
	// ----------------------------

	accountService := &services.AccountService{
		GoogleAdapter:   googleAdapter,
		FacebookAdapter: facebookAdapter,
		TokenAdapter:    pasetoAdapter,
		EmailAdapter:    sesAdapter,
		SmsAdapter:      snsAdapter,

		Db: db,

		AccountRepository:       accountRepository,
		MagicLinkCodeRepository: magicLinkCodeRepository,
		RefreshTokenRepository:  refreshTokenRepository,
	}

	log.Print("Services initialized")

	// ----------------------------
	//
	// Routers
	//
	// ----------------------------

	http.NewHttpDelivery(&http.NewHttpDeliveryInput{
		AccountService: accountService,
	}).Listen()
}
