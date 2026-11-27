package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Tahsin005/olx-api/internal/httpx"
	"github.com/Tahsin005/olx-api/internal/middlewares"
)

type ListingHandler struct {
	db *sql.DB
	logger *slog.Logger
}

func NewListingHandler(db *sql.DB, logger *slog.Logger) *ListingHandler {
	return &ListingHandler{
		db: db, 
		logger: logger,
	}
}

func (lh *ListingHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middlewares.RequestIDFromContext(ctx)

	rows, err := lh.db.QueryContext(
		ctx,
		`select id, title, description, price, city, created_at 
		from listings
		order by created_at desc 
		limit 100`,
	)
	if err != nil {
		lh.logger.Error("failed to query listings", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}
	defer rows.Close()

	listings := []ListListingsResponse{}

	for rows.Next() {
		var l ListListingsResponse
		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.CreatedAt); err != nil {
			lh.logger.Error("failed to scan listing row", "request_id", requestId, "err", err)
			httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
			return
		}

		listings = append(listings, l)
	}

	if err := rows.Err(); err != nil {
		lh.logger.Error("failed to iterate listing rows", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(listings)
}

func (lh *ListingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middlewares.RequestIDFromContext(ctx)
	id := r.PathValue("id")

	result, err := lh.db.ExecContext(
		ctx,
		`delete from listings where id = $1`, id,
	)

	if err != nil {
		lh.logger.Error("failed to delete listing", "listing_id", id, "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		httpx.Error(w, http.StatusNotFound, "listing not found", httpx.CodeNotFound)
		return
	}

	lh.logger.Info("listing deleted", "request_id", requestId, "listing_id", id)

	w.WriteHeader(http.StatusNoContent)
}

func (lh *ListingHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middlewares.RequestIDFromContext(ctx)

	var req CreateListingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lh.logger.Error("failed to decode request body", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusBadRequest, "invalid body", httpx.CodeMalformedJSON)
		return
	}

	if err := req.Validate(); err != nil {
		var verr *ValidationError
		errors.As(err, &verr)
		lh.logger.Error("validation failed", "request_id", requestId, "err", err)
		httpx.ValidationError(w, http.StatusUnprocessableEntity, err.Error(), httpx.CodeValidationFailed, verr.Field)
		return
	}

	row := lh.db.QueryRowContext(ctx,
		`insert into listings (title, description, price, city) values ($1, $2, $3, $4) returning id, created_at`,
		req.Title, req.Description, req.Price, req.City,
	)

	var out CreateListingResponse
	if err := row.Scan(&out.ID, &out.CreatedAt); err != nil {
		lh.logger.Error("failed to insert listing", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	lh.logger.Info("listing created", "request_id", requestId, "listing_id", out.ID, "created_at", out.CreatedAt)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(out)
}
