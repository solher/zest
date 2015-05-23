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
		Controller, Handler string
	}

	RouteDirectory struct {
		routes map[DirectoryKey]Route
		render AbstractRender
	}
)

func NewDirectoryKey(controller string) *DirectoryKey {
	return &DirectoryKey{Controller: controller}
}

func (k *DirectoryKey) For(handler string) *DirectoryKey {
	k.Handler = handler
	return k
}

func NewRouteDirectory(render AbstractRender) *RouteDirectory {
	return &RouteDirectory{routes: make(map[DirectoryKey]Route), render: render}
}

func (routeDir *RouteDirectory) Add(key *DirectoryKey, route *Route, checkPermissions bool) {
	if checkPermissions {
		permissionGate := NewPermissionGate(route.Handler, routeDir, routeDir.render)
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
