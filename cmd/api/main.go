package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Tahsin005/olx-api/internal/config"
	"github.com/Tahsin005/olx-api/internal/db"
	"github.com/Tahsin005/olx-api/internal/handlers"
	"github.com/Tahsin005/olx-api/internal/middlewares"
	"github.com/Tahsin005/olx-api/internal/repository"
	"github.com/Tahsin005/olx-api/internal/service"
)

func main() {
	cfg := config.MustLoad()
	db, err := db.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("main.db.connect: %v", err)
	}

	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level: slog.LevelInfo,
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	fmt.Println("database connected")
	fmt.Println("starting olx server [" + cfg.Env + "]")

	repo := repository.NewListingRepository(db)
	svc := service.NewListingService(repo, logger)
	lh := handlers.NewListingHandler(svc, logger)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handlers.Health)
	mux.HandleFunc("GET /listings", lh.List)
	mux.HandleFunc("DELETE /listings/{id}", lh.Delete)
	mux.HandleFunc("POST /listings", lh.Create)

	handler := middlewares.RequestID(mux)
	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	log.Printf("server is listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
