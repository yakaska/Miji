package link

import (
	"Miji/internal/core/domain"
	coreerrors "Miji/internal/core/errors"
	"context"
	"errors"
	"fmt"
)

var ErrSlugExists = errors.New("slug already exists")

type Repository interface {
	Create(ctx context.Context, link domain.Link) (domain.Link, error)
	ExistsBySlug(ctx context.Context, slug string) (bool, error)
}

type SlugGenerator interface {
	Generate() (string, error)
}
type Service struct {
	repo          Repository
	slugGenerator SlugGenerator
}

func NewService(repo Repository, slugGenerator SlugGenerator) *Service {
	return &Service{
		repo:          repo,
		slugGenerator: slugGenerator,
	}
}

const maxSlugGenerationAttempts = 5

func (s *Service) Create(ctx context.Context, link domain.Link) (domain.Link, error) {
	if link.Slug == domain.UninitializedSlug {
		slug, err := s.generateUniqueSlug(ctx)
		if err != nil {
			return domain.Link{}, err
		}
		link.Slug = slug
	} else {
		if err := s.ensureSlugAvailable(ctx, link.Slug); err != nil {
			return domain.Link{}, err
		}
	}

	return s.repo.Create(ctx, link)
}

func (s *Service) generateUniqueSlug(ctx context.Context) (string, error) {
	for range maxSlugGenerationAttempts {
		slug, err := s.slugGenerator.Generate()
		if err != nil {
			return "", fmt.Errorf("generate slug: %w", err)
		}

		exists, err := s.repo.ExistsBySlug(ctx, slug)
		if err != nil {
			return "", fmt.Errorf("check slug uniqueness: %w", err)
		}
		if !exists {
			return slug, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique slug after %d attempts", maxSlugGenerationAttempts)
}

func (s *Service) ensureSlugAvailable(ctx context.Context, slug string) error {
	exists, err := s.repo.ExistsBySlug(ctx, slug)
	if err != nil {
		return fmt.Errorf("check slug uniqueness: %w", err)
	}
	if exists {
		return fmt.Errorf("%w: %w", coreerrors.ErrConflict, ErrSlugExists)
	}
	return nil
}
