package domain

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("not found")

type Listing struct {
	ID          string
	Title       string
	Description string
	Price       int64
	City        string
	CreatedAt   time.Time
}
