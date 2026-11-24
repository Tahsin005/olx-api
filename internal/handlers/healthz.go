package handlers

import (
	"log/slog"
	"net/http"

	"github.com/Tahsin005/olx-api/internal/middlewares"
)

func Health(w http.ResponseWriter, r *http.Request) {
	requestId := middlewares.RequestIDFromContext(r.Context())
	slog.Info("health check", "request_id", requestId)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	w.Write([]byte(`{"status": "ok"}`))
}