package main

import (
	"log"

	"github.com/aparnasukesh/user-admin-super_admin-svc/config"
	"github.com/aparnasukesh/user-admin-super_admin-svc/internal/di"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	server, err := di.InitResources(cfg)

	if err := server(); err != nil {
		log.Fatal(err)
	}

}
