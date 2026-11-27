package handlers

import (
	"fmt"
	"strings"
	"time"
)

type CreateListingRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	City        string `json:"city"`
}

type CreateListingResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type ListListingsResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"created_at"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func (req CreateListingRequest) Validate() error {
	if strings.TrimSpace(req.Title) == "" {
		return &ValidationError{Field: "title", Message: "title is required"}
	}
	if req.Description == "" {
		return &ValidationError{Field: "description", Message: "description is required"}
	}
	if req.Price <= 0 {
		return &ValidationError{Field: "price", Message: "price must be greater than zero"}
	}
	if req.City == "" {
		return &ValidationError{Field: "city", Message: "city is required"}
	}
	return nil
}