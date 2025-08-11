package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kobayashiyabako16g/passkey-auth-example/internal/config"
	"github.com/kobayashiyabako16g/passkey-auth-example/internal/domain/repository"
	"github.com/kobayashiyabako16g/passkey-auth-example/internal/handler"
	"github.com/kobayashiyabako16g/passkey-auth-example/internal/handler/middleware"
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
	// Repository
	sessionRepository := repository.NewSession(kvClient)
	userRepository := repository.NewUser(dbClient)

	mux := http.NewServeMux()
	auth := handler.NewAuth(sessionRepository, userRepository)
	router := handler.NewRouter(auth)
	router.HandleRequest(mux)

	server := middleware.CORSMiddleware(mux)
	server = middleware.LogMiddleware(server)
	port := fmt.Sprintf(":%s", cfg.Port)
	logger.Info(ctx, fmt.Sprintf("Starting server on port %s", port))
	if err := http.ListenAndServe(port, server); err != nil {
		panic(err)
	}
}
