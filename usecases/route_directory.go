package usecases

import (
	"net/http"

	"github.com/dimfeld/httptreemux"
)

type (
	HandlerFunc func(w http.ResponseWriter, r *http.Request, params map[string]string)

	Routes map[DirectoryKey]Route

	Route struct {
		Method, Path              string
		Handler, EffectiveHandler HandlerFunc
		Visible, CheckPermissions bool
	}

	DirectoryKey struct {
		Ressource, Method string
	}

	RouteDirectory struct {
		accountInter AbstractAccountInter
		routes       map[DirectoryKey]Route
		render       AbstractRender
	}
)

func NewDirectoryKey(ressources string) *DirectoryKey {
	return &DirectoryKey{Ressource: ressources}
}

func (k *DirectoryKey) For(method string) *DirectoryKey {
	k.Method = method
	return k
}

func NewRouteDirectory(accountInter AbstractAccountInter, render AbstractRender) *RouteDirectory {
	return &RouteDirectory{accountInter: accountInter, routes: make(map[DirectoryKey]Route), render: render}
}

func (routeDir *RouteDirectory) Routes() map[DirectoryKey]Route {
	return routeDir.routes
}

func (routeDir *RouteDirectory) Add(key *DirectoryKey, route *Route) {
	routeDir.routes[*key] = *route
}

func (routeDir *RouteDirectory) Get(key *DirectoryKey) *Route {
	route := routeDir.routes[*key]
	return &route
}

func (routeDir *RouteDirectory) Register(router *httptreemux.TreeMux) {
	var keys []DirectoryKey
	routes := routeDir.routes

	for k := range routes {
		keys = append(keys, k)
	}

	for _, k := range keys {
		route := routes[k]

		if route.Visible {
			handler := route.Handler

			if route.CheckPermissions {
				permissionGate := NewPermissionGate(k.Ressource, k.Method, routeDir.accountInter, routeDir.render, route.Handler)
				handler = permissionGate.Handler
			}

			route.EffectiveHandler = handler
			routes[k] = route

			router.Handle(route.Method, route.Path, httptreemux.HandlerFunc(handler))
		}
	}
}
