package link_transport_http

import "net/http"

type CreateLinkRequest struct {
	Slug        *string
	OriginalURL string
}

func (h *LinkHTTPHandler) CreateLink(w http.ResponseWriter, r *http.Request) {

}
