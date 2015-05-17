package interfaces

import (
	"github.com/Solher/auth-scaffold/usecases"
	"github.com/julienschmidt/httprouter"
)

type (
	Route struct {
		Method  string
		Path    string
		Handler *httprouter.Handle
	}

	DirectoryKey struct {
		Controller interface{}
		Handler    string
	}

	RouteDirectory map[DirectoryKey]Route
)

func NewDirectoryKey(controller interface{}) *DirectoryKey {
	return &DirectoryKey{Controller: controller}
}

func (k *DirectoryKey) For(handler string) DirectoryKey {
	k.Handler = handler
	return *k
}

func NewRouteDirectory() RouteDirectory {
	return RouteDirectory{}
}

func (routes RouteDirectory) Register(router *httprouter.Router, permissions usecases.PermissionDirectory, render AbstractRender) {
	var keys []DirectoryKey

	for k := range routes {
		keys = append(keys, k)
	}

	for _, k := range keys {
		route := routes[k]
		permissionGate := NewPermissionGate(route.Handler, routes, permissions, render)
		router.Handle(route.Method, route.Path, permissionGate.Handler)
	}
}
