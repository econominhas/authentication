package main

import (
	"database/sql"
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
	logger := newLogger()

	logger.Info("Start")

	// ----------------------------
	//
	// Env
	//
	// ----------------------------

	logger.Debug("Validating env vars")

	validateEnvs(logger)

	logger.Info("Env vars validated")

	// ----------------------------
	//
	// Databases
	//
	// ----------------------------

	logger.Debug("Trying to connect to the database")

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Error(
			"Fail to connect to the database",
			"message", err.Error(),
		)
		panic(1)
	}
	defer db.Close()

	logger.Info("Connected to database")

	// ----------------------------
	//
	// Adapters
	//
	// ----------------------------

	logger.Debug("Initializing adapters")

	googleAdapter := google.NewGoogle(logger)
	facebookAdapter := facebook.NewFacebook(logger)
	pasetoAdapter := paseto.NewPaseto(logger)
	secretAdapter := secret.NewSecret(logger)
	sesAdapter := ses.NewSes(logger)
	snsAdapter := sns.NewSns(logger)
	ulidAdapter := ulid.NewUlid(logger)

	logger.Info("Adapters initialized")

	// ----------------------------
	//
	// Repositories
	//
	// ----------------------------

	logger.Debug("Initializing repositories")

	accountRepository := &repositories.AccountRepository{
		Logger: logger,

		IdAdapter: ulidAdapter,
	}
	magicLinkCodeRepository := &repositories.MagicLinkCodeRepository{
		Logger: logger,

		SecretAdapter: secretAdapter,
	}
	refreshTokenRepository := &repositories.RefreshTokenRepository{
		Logger: logger,

		IdAdapter:     ulidAdapter,
		SecretAdapter: secretAdapter,
		TokenAdapter:  pasetoAdapter,
	}

	logger.Info("Repositories initialized")

	// ----------------------------
	//
	// Services
	//
	// ----------------------------

	logger.Debug("Initializing services")

	accountService := &services.AccountService{
		Logger: logger,

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

	logger.Info("Services initialized")

	// ----------------------------
	//
	// Routers
	//
	// ----------------------------

	logger.Debug("Initializing http server")
	logger.Info("Http server initialized on port: " + os.Getenv("PORT"))

	http.NewHttpDelivery(&http.NewHttpDeliveryInput{
		AccountService: accountService,
	}).Listen()
}
