package handlers

import "time"

type CreateListingRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Price       int64  `json:"price" validate:"required,gt=0"`
	City        string `json:"city" validate:"required"`
}

type CreateListingResponse struct {
	ID string `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}