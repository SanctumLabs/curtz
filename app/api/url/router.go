package url

import (
	"github.com/sanctumlabs/curtz/app/internal/core/contracts"
	"github.com/sanctumlabs/curtz/app/server/router"
)

// urlRouter is a router for the url API.
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
func (route *urlRouter) Routes() []router.Route {
	return route.routes
}

func (route *urlRouter) initRoutes() {
	route.routes = []router.Route{
		router.NewGetRoute("/api/v1/:identifier", route.getUrl),
		router.NewPostRoute("/api/v1/", route.createUrl),
	}
}
