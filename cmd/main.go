package main

import (
	"log"

	"github.com/aAmer0neee/authentication-service-test-task/internal/api"
	"github.com/aAmer0neee/authentication-service-test-task/internal/auth"
	"github.com/aAmer0neee/authentication-service-test-task/internal/config"
	"github.com/aAmer0neee/authentication-service-test-task/internal/repository"
)

func main() {
	cfg := config.LoadConfig()
	repo, err := repository.New(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	service := auth.New(repo, cfg)
	srv := api.ConfigureApi(cfg, service)

	srv.ListenAndServe()

}
