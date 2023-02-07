package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/aopal/go-cache/pkg/config"
	"github.com/aopal/go-cache/pkg/handler"
	"github.com/aopal/go-cache/pkg/middlewares"
)

func main() {
	configPath := os.Args[1]
	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatalf("Could not initialize server: %+v\n", err)
	}

	server, err := handler.New(cfg)
	if err != nil {
		log.Fatalf("Could not initialize server: %+v\n", err)
	}

	http.HandleFunc("/", middlewares.WithLogging(server.Serve))

	log.Printf("Listening on :%s...", cfg.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
