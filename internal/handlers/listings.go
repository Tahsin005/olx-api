package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Tahsin005/olx-api/internal/domain"
	"github.com/Tahsin005/olx-api/internal/httpx"
	"github.com/Tahsin005/olx-api/internal/middlewares"
	"github.com/Tahsin005/olx-api/internal/service"
)

type ListingHandler struct {
	svc    *service.ListingService
	logger *slog.Logger
}

func NewListingHandler(svc *service.ListingService, logger *slog.Logger) *ListingHandler {
	return &ListingHandler{
		svc:    svc,
		logger: logger,
	}
}

func (lh *ListingHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middlewares.RequestIDFromContext(ctx)

	listings, err := lh.svc.List(ctx)
	if err != nil {
		lh.logger.Error("failed to list listings", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	out := make([]ListListingsResponse, len(listings))
	for i, l := range listings {
		out[i] = ListListingsResponse{
			ID:          l.ID,
			Title:       l.Title,
			Description: l.Description,
			Price:       l.Price,
			City:        l.City,
			CreatedAt:   l.CreatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(out)
}

func (lh *ListingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middlewares.RequestIDFromContext(ctx)
	id := r.PathValue("id")

	if err := lh.svc.Delete(ctx, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			httpx.Error(w, http.StatusNotFound, "listing not found", httpx.CodeNotFound)
			return
		}

		lh.logger.Error("failed to delete listing", "listing_id", id, "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

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

	listing, err := lh.svc.Create(ctx, domain.Listing{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		City:        req.City,
	})
	if err != nil {
		lh.logger.Error("failed to create listing", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	out := CreateListingResponse{
		ID:        listing.ID,
		CreatedAt: listing.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(out)
}
