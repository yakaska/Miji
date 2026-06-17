package server

import (
	"Miji/internal/links"
	"net/http"
)

func AddRoutes(mux *http.ServeMux, linksHandler *links.HTTPHandler) {
	mux.HandleFunc("POST /api/links", linksHandler.CreateLink)
}
