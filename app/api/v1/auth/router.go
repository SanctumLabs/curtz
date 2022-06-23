package auth

import (
	"fmt"

	"github.com/sanctumlabs/curtz/app/internal/core/contracts"
	"github.com/sanctumlabs/curtz/app/server/router"
)

type authRouter struct {
	svc             contracts.UserService
	notificationSvc contracts.NotificationService
	authSvc         contracts.AuthService
	routes          []router.Route
	baseUri         string
}

func NewRouter(baseUri string, svc contracts.UserService, notificationSvc contracts.NotificationService, authSvc contracts.AuthService) router.Router {
	r := &authRouter{
		svc:             svc,
		notificationSvc: notificationSvc,
		authSvc:         authSvc,
		baseUri:         baseUri,
	}
	r.initRoutes()
	return r
}

func (hdl authRouter) Routes() []router.Route {
	return hdl.routes
}

func (hdl *authRouter) initRoutes() {
	hdl.routes = []router.Route{
		router.NewPostRoute(fmt.Sprintf("%s/auth/register", hdl.baseUri), hdl.register),
		router.NewGetRoute(fmt.Sprintf("%s/auth/login", hdl.baseUri), hdl.login),
	}
}