package usecases

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRoutes(t *testing.T) {
	// router := httprouter.New()
	//
	Convey("Testing routes...", t, func() {
		// 	Convey("Should be able to add routes to a directory and register it.", func() {
		// 		routes := NewRouteDirectory()
		// 		key := NewDirectoryKey(nil)
		//
		// 		routes[key.For("key1")] = Route{Method: "GET", Path: "/", Handler: nil}
		// 		So(routes[DirectoryKey{Handler: "key1"}].Method, ShouldEqual, "GET")
		//
		// 		routes[key.For("key2")] = Route{Method: "POST", Path: "/", Handler: nil}
		// 		So(routes[DirectoryKey{Handler: "key2"}].Method, ShouldEqual, "POST")
		//
		// 		routes.Register(router)
		// 	})
		//
		// 	Convey("Should be able to mock HTTP requests.", func() {
		// 		route := Route{Method: "GET", Path: "/", Handler: func(w http.ResponseWriter, r *http.Request, _ map[string]string) {}}
		// 		res := MockHTTPRequest(route, "", "", nil)
		// 		So(res, ShouldEqual, "")
		//
		// 		route.Handler = nil
		// 		So(func() { MockHTTPRequest(route, "", "", nil) }, ShouldPanic)
		// 	})
	})
}
