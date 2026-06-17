package links

import (
	"Miji/internal/core/domain"
	"Miji/internal/core/logger"
	"Miji/internal/core/transport"
	"context"
	"net/http"
	"net/url"
	"strings"
)

type CreateLinkRequest struct {
	Slug        string `json:"slug"`
	OriginalURL string `json:"originalURL"`
}

func (r CreateLinkRequest) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if strings.TrimSpace(r.OriginalURL) == "" {
		problems["originalURL"] = "originalURL is required"
	} else if _, err := url.ParseRequestURI(r.OriginalURL); err != nil {
		problems["originalURL"] = "originalURL must be a valid URL"
	}

	if r.Slug != "" && len(r.Slug) > 120 {
		problems["slug"] = "slug must be at most 120 characters"
	}

	return problems
}

func (h *HTTPHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	req, problems, err := transport.DecodeValid[CreateLinkRequest](r)
	if err != nil {
		if problems != nil {
			transport.ValidationError(w, r, log.Logger, problems)
			return
		}
		transport.Error(w, r, log.Logger, err, "failed to create link")
		return
	}

	link, err := h.service.Create(r.Context(), domainFromDTO(req))
	if err != nil {
		transport.Error(w, r, log.Logger, err, "failed to create link")
		return
	}

	_ = transport.Encode(w, r, http.StatusCreated, link)
}

func domainFromDTO(req CreateLinkRequest) domain.Link {
	return domain.Link{
		ID:          domain.UninitializedID,
		OwnerID:     0,
		Slug:        req.Slug,
		OriginalURL: req.OriginalURL,
		ExpiresAt:   nil,
	}
}
