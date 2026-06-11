package domain

import "time"

type Link struct {
	ID          int
	OwnerID     int
	Slug        string
	OriginalURL string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	IsActive    bool
	VisitCount  int
}

func NewLink(
	ID int,
	OwnerID int,
	Slug string,
	OriginalURL string,
	CreatedAt time.Time,
	UpdatedAt *time.Time,
	IsActive bool,
	VisitCount int,
) Link {
	return Link{
		ID:          ID,
		OwnerID:     OwnerID,
		Slug:        Slug,
		OriginalURL: OriginalURL,
		CreatedAt:   CreatedAt,
		UpdatedAt:   UpdatedAt,
		IsActive:    IsActive,
		VisitCount:  VisitCount,
	}
}
