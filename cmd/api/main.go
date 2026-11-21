package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Tahsin005/olx-api/internal/config"
	"github.com/Tahsin005/olx-api/internal/handlers"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println("starting olx server [" + cfg.Env + "]")
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handlers.Health)

	srv := http.Server{
		Addr: ":" + cfg.Port,
		Handler: mux,
		ReadTimeout: time.Second  * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout: time.Second * 60,
	}

	log.Printf("server is listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}