package url

import (
	"fmt"

	"github.com/sanctumlabs/curtz/app/internal/core/contracts"
	"github.com/sanctumlabs/curtz/app/server/router"
)

// urlRouter is a router for the url API.
type urlRouter struct {
	svc     contracts.UrlService
	routes  []router.Route
	baseUri string
}

// NewUrlRouter initializes a new router
func NewUrlRouter(baseUri string, s contracts.UrlService) router.Router {
	r := &urlRouter{
		svc:     s,
		baseUri: baseUri,
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
		router.NewGetRoute(fmt.Sprintf("%s/:identifier", route.baseUri), route.getUrl),
		router.NewPostRoute(fmt.Sprintf("%s", route.baseUri), route.createShortUrl),
	}
}
