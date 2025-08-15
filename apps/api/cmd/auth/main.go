package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/kobayashiyabako16g/passkey-auth-example/internal/config"
	"github.com/kobayashiyabako16g/passkey-auth-example/internal/domain/repository"
	"github.com/kobayashiyabako16g/passkey-auth-example/internal/ui/handler"
	"github.com/kobayashiyabako16g/passkey-auth-example/internal/ui/middleware"
	"github.com/kobayashiyabako16g/passkey-auth-example/internal/ui/router"
	"github.com/kobayashiyabako16g/passkey-auth-example/internal/usecase"
	"github.com/kobayashiyabako16g/passkey-auth-example/pkg/db"
	"github.com/kobayashiyabako16g/passkey-auth-example/pkg/kvstore"
	"github.com/kobayashiyabako16g/passkey-auth-example/pkg/logger"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	kvClient, err := kvstore.NewValKeyClient(cfg.ValKeyConfig)
	if err != nil {
		panic(err)
	}
	dbClient, err := db.NewClient("postgres", "postgres://postgres:postgres@db:5432/app")
	if err != nil {
		panic(err)
	}

	// webautn
	wconfig := &webauthn.Config{
		RPDisplayName: "Passkey Demo",                    // Display Name for your site
		RPID:          "localhost",                       // Generally the domain name for your site
		RPOrigins:     []string{"http://localhost:5173"}, // Vite dev server origin
	}
	webAuthn, err := webauthn.New(wconfig)
	if err != nil {
		panic(err)
	}

	// Repository
	sessionRepository := repository.NewSession(kvClient)
	userRepository := repository.NewUser(dbClient)

	// Usecase
	authUsecase := usecase.NewAuth(sessionRepository, userRepository, webAuthn)

	mux := http.NewServeMux()
	auth := handler.NewAuth(authUsecase)
	rt := router.NewRouter(auth)
	rt.HandleRequest(mux)

	server := middleware.CORSMiddleware(mux, cfg.AllowOrigin)
	server = middleware.LogMiddleware(server)
	port := fmt.Sprintf(":%s", cfg.Port)
	logger.Info(ctx, fmt.Sprintf("Starting server on port %s", port))
	if err := http.ListenAndServe(port, server); err != nil {
		panic(err)
	}
}
