package links

import (
	"Miji/internal/core/db"
	"Miji/internal/core/domain"
	"context"
	"fmt"
	"time"
)

type PostgresRepository struct {
	pool db.Pool
}

func NewPostgresRepository(pool db.Pool) *PostgresRepository {
	return &PostgresRepository{
		pool: pool,
	}
}

func (r *PostgresRepository) Create(ctx context.Context, link domain.Link) (domain.Link, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	row := r.pool.QueryRow(ctx,
		`INSERT INTO miji.links (owner_id, slug, original_url, expires_at, is_active, visit_count)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, created_at`,
		link.OwnerID, link.Slug, link.OriginalURL, link.ExpiresAt, link.IsActive, link.VisitCount,
	)

	var id int
	var createdAt time.Time
	if err := row.Scan(&id, &createdAt); err != nil {
		return domain.Link{}, fmt.Errorf("create link: %w", err)
	}

	link.ID = id
	link.CreatedAt = createdAt

	return link, nil
}
