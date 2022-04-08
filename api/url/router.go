package url

import (
	"github.com/sanctumlabs/curtz/internal/core/contracts"
	"github.com/sanctumlabs/curtz/server/router"
)

// urlRouter is a router to talk
type urlRouter struct {
	svc    contracts.UrlService
	routes []router.Route
}

// NewUrlRouter initializes a new router
func NewUrlRouter(s contracts.UrlService) router.Router {
	r := &urlRouter{
		svc: s,
	}
	r.initRoutes()
	return r
}

// Routes returns the available routes to the controller
func (hdl *urlRouter) Routes() []router.Route {
	return hdl.routes
}

func (hdl *urlRouter) initRoutes() {
	hdl.routes = []router.Route{
		router.NewGetRoute("/:identifier", hdl.getUrl),
		router.NewPostRoute("/", hdl.createUrl),
	}
}
