package service

import (
	"context"
	"log/slog"

	"github.com/Tahsin005/olx-api/internal/domain"
	"github.com/Tahsin005/olx-api/internal/repository"
)

type ListingService struct {
	repo   repository.ListingRepository
	logger *slog.Logger
}

func NewListingService(repo repository.ListingRepository, logger *slog.Logger) *ListingService {
	return &ListingService{
		repo:   repo,
		logger: logger,
	}
}

func (s *ListingService) List(ctx context.Context) ([]domain.Listing, error) {
	listings, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	return listings, nil
}

func (s *ListingService) Create(ctx context.Context, l domain.Listing) (domain.Listing, error) {
	created, err := s.repo.Create(ctx, l)
	if err != nil {
		return domain.Listing{}, err
	}

	s.logger.Info("listing created", "listing_id", created.ID, "created_at", created.CreatedAt)

	return created, nil
}

func (s *ListingService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.logger.Info("listing deleted", "listing_id", id)

	return nil
}
