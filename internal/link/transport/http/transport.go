package link_transport_http

type LinkHTTPHandler struct {
	linkService LinkService
}

type LinkService interface {
}

func NewLinkHTTPHandler(
	linkService LinkService,
) *LinkHTTPHandler {
	return &LinkHTTPHandler{
		linkService: linkService,
	}
}
