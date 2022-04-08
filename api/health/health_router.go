package health

import "github.com/sanctumlabs/curtz/server/router"

// healthRouter is a router to talk to when we are checking health of the api
type healthRouter struct {
	routes []router.Route
}

// NewHealthRouter initializes a new router
func NewHealthRouter() router.Router {
	r := &healthRouter{}
	r.initRoutes()
	return r
}

// Routes returns the available routes to the controller
func (h *healthRouter) Routes() []router.Route {
	return h.routes
}

func (h *healthRouter) initRoutes() {
	h.routes = []router.Route{
		router.NewGetRoute("/health", h.health),
	}
}
