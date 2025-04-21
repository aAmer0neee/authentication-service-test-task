package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aAmer0neee/authentication-service-test-task/internal/api"
	"github.com/aAmer0neee/authentication-service-test-task/internal/auth"
	"github.com/aAmer0neee/authentication-service-test-task/internal/config"
	"github.com/aAmer0neee/authentication-service-test-task/internal/logger"
	"github.com/aAmer0neee/authentication-service-test-task/internal/notify"
	"github.com/aAmer0neee/authentication-service-test-task/internal/repository"
	"github.com/aAmer0neee/authentication-service-test-task/internal/token"
)

func main() {
	cfg := config.LoadConfig()
	repo, err := repository.New(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	service := auth.New(repo,
		token.New(cfg.AuthSecret),
		notify.New(&cfg),
		logger.New(cfg.Server.Env))

	srv := api.ConfigureApi(cfg, service)
	go func() {
		srv.ListenAndServe()
	}()

	fmt.Println("SERVER STARTING AT -> ", cfg.Server.Host+":"+cfg.Server.Port)
	shutdown(srv)
}

func shutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Graceful Shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	srv.Shutdown(ctx)
}
