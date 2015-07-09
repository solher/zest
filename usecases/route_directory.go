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
		Resource, Method string
	}

	RouteDirectory struct {
		accountInter AbstractAccountInter
		Routes       map[DirectoryKey]Route
		render       AbstractRender
	}
)

func WrapHandler(handler http.Handler) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		handler.ServeHTTP(w, r)
	}
}

func NewDirectoryKey(resources string) *DirectoryKey {
	return &DirectoryKey{Resource: resources}
}

func (k *DirectoryKey) For(method string) *DirectoryKey {
	k.Method = method
	return k
}

func NewRouteDirectory(accountInter AbstractAccountInter, render AbstractRender) *RouteDirectory {
	return &RouteDirectory{accountInter: accountInter, Routes: make(map[DirectoryKey]Route), render: render}
}

func (routeDir *RouteDirectory) Add(key *DirectoryKey, route *Route) {
	routeDir.Routes[*key] = *route
}

func (routeDir *RouteDirectory) Get(key *DirectoryKey) *Route {
	route := routeDir.Routes[*key]
	return &route
}

func (routeDir *RouteDirectory) Register(router *httptreemux.TreeMux) {
	var keys []DirectoryKey
	routes := routeDir.Routes

	for k := range routes {
		keys = append(keys, k)
	}

	for _, k := range keys {
		route := routes[k]

		if route.Visible {
			handler := route.Handler

			if route.CheckPermissions {
				permissionGate := NewPermissionGate(k.Resource, k.Method, routeDir.accountInter, routeDir.render, route.Handler)
				handler = permissionGate.Handler
			}

			route.EffectiveHandler = handler
			routes[k] = route

			router.Handle(route.Method, route.Path, httptreemux.HandlerFunc(handler))
		}
	}
}
