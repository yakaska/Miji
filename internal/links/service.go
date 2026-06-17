package links

import (
	"Miji/internal/core/domain"
	"context"
)

type Service struct {
	repo Repository
}

type Repository interface {
	Create(ctx context.Context, link domain.Link) (domain.Link, error)
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, link domain.Link) (domain.Link, error) {

	return s.repo.Create(ctx, link)
}
