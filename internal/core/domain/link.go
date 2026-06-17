package domain

import "time"

type Link struct {
	ID          int
	OwnerID     int
	Slug        string
	OriginalURL string
	CreatedAt   time.Time
	ExpiresAt   *time.Time
	IsActive    bool
	VisitCount  int
}
