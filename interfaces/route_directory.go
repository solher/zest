package interfaces

import (
	"github.com/Solher/auth-scaffold/usecases"
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
		routes      map[DirectoryKey]Route
		permissions usecases.PermissionDirectory
		render      AbstractRender
	}
)

func NewDirectoryKey(controller string) *DirectoryKey {
	return &DirectoryKey{Controller: controller}
}

func (k *DirectoryKey) For(handler string) *DirectoryKey {
	k.Handler = handler
	return k
}

func NewRouteDirectory(permissions usecases.PermissionDirectory, render AbstractRender) *RouteDirectory {
	return &RouteDirectory{routes: make(map[DirectoryKey]Route), permissions: permissions, render: render}
}

func (routeDir *RouteDirectory) Add(key *DirectoryKey, route *Route, checkPermissions bool) {
	if checkPermissions {
		permissionGate := NewPermissionGate(route.Handler, routeDir, routeDir.permissions, routeDir.render)
		gatedHandler := httptreemux.HandlerFunc(permissionGate.Handler)
		route.Handler = &gatedHandler
	}

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
			router.Handle(route.Method, route.Path, *route.Handler)
		}
	}
}
