package interfaces

import (
	"errors"

	"github.com/dimfeld/httptreemux"
)

type (
	Route struct {
		Method, Path string
		Handler      *httptreemux.HandlerFunc
		Visible      bool
	}

	DirectoryKey struct {
		Ressources, Method string
	}

	RouteDirectory struct {
		accountInter AbstractAccountInter
		routes       map[DirectoryKey]Route
		render       AbstractRender
	}
)

func NewDirectoryKey(ressources string) *DirectoryKey {
	return &DirectoryKey{Ressources: ressources}
}

func (k *DirectoryKey) For(method string) *DirectoryKey {
	k.Method = method
	return k
}

func NewRouteDirectory(accountInter AbstractAccountInter, render AbstractRender) *RouteDirectory {
	return &RouteDirectory{accountInter: accountInter, routes: make(map[DirectoryKey]Route), render: render}
}

func (routeDir *RouteDirectory) Add(key *DirectoryKey, route *Route, checkPermissions bool) {
	if checkPermissions {
		permissionGate := NewPermissionGate(routeDir.accountInter, route.Handler, routeDir, routeDir.render)
		gatedHandler := httptreemux.HandlerFunc(permissionGate.Handler)
		route.Handler = &gatedHandler
	}

	routeDir.routes[*key] = *route
}

func (routeDir *RouteDirectory) Get(key *DirectoryKey) *Route {
	route := routeDir.routes[*key]
	return &route
}

func (routeDir *RouteDirectory) GetKey(handler *httptreemux.HandlerFunc) (*DirectoryKey, error) {
	for key, route := range routeDir.routes {
		if route.Handler == handler {
			return &key, nil
		}
	}

	return nil, errors.New("Handler not found.")
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
			router.Handle(route.Method, route.Path, *route.Handler)
		}
	}
}
