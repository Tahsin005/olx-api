package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxKey int
const (
	requestIDKey ctxKey = iota
)
const (
	requestId = "X-Request-ID"
)

func RequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	id, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return ""
	}

	return id
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(requestId)
		if id == "" {
			id = uuid.NewString()
		}

		ctx := context.WithValue(r.Context(), requestIDKey, id)
		r = r.WithContext(ctx)

		w.Header().Add(requestId, id)

		next.ServeHTTP(w, r)
	})
}