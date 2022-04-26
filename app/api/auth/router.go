package auth

import (
	"github.com/sanctumlabs/curtz/app/internal/core/contracts"
	"github.com/sanctumlabs/curtz/app/server/router"
)

type authRouter struct {
	svc    contracts.UserService
	routes []router.Route
}

func NewRouter(svc contracts.UserService) router.Router {
	r := &authRouter{
		svc: svc,
	}
	r.initRoutes()
	return r
}

func (hdl authRouter) Routes() []router.Route {
	return hdl.routes
}

func (hdl *authRouter) initRoutes() {
	hdl.routes = []router.Route{
		router.NewPostRoute("/", hdl.signUp),
		router.NewPostRoute("/", hdl.login),
	}
}
