package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Tahsin005/olx-api/internal/domain"
)

type ListingRepository interface {
	List(ctx context.Context) ([]domain.Listing, error)
	Create(ctx context.Context, l domain.Listing) (domain.Listing, error)
	Delete(ctx context.Context, id string) error
}

type listingRepo struct {
	db *sql.DB
}

func NewListingRepository(db *sql.DB) ListingRepository {
	return &listingRepo{db: db}
}

func (r *listingRepo) List(ctx context.Context) ([]domain.Listing, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`select id, title, description, price, city, created_at 
		from listings
		order by created_at desc 
		limit 100`,
	)
	if err != nil {
		return nil, fmt.Errorf("query listings: %w", err)
	}
	defer rows.Close()

	var listings []domain.Listing

	for rows.Next() {
		var l domain.Listing
		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan listing: %w", err)
		}

		listings = append(listings, l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate listings: %w", err)
	}

	return listings, nil
}

func (r *listingRepo) Create(ctx context.Context, l domain.Listing) (domain.Listing, error) {
	row := r.db.QueryRowContext(ctx,
		`insert into listings (title, description, price, city) values ($1, $2, $3, $4) returning id, created_at`,
		l.Title, l.Description, l.Price, l.City,
	)

	if err := row.Scan(&l.ID, &l.CreatedAt); err != nil {
		return domain.Listing{}, fmt.Errorf("insert listing: %w", err)
	}

	return l, nil
}

func (r *listingRepo) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(
		ctx,
		`delete from listings where id = $1`, id,
	)

	if err != nil {
		return fmt.Errorf("delete listing: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}
