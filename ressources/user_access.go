// Generated by: main
// TypeWriter: access
// Directive: +gen on User

package ressources

import (
	"github.com/Solher/auth-scaffold/interfaces"
	"github.com/dimfeld/httptreemux"
)

func setUserAccessOptions(routeDir *interfaces.RouteDirectory, controller *UserCtrl) {
	key := interfaces.NewDirectoryKey("users")
	create := httptreemux.HandlerFunc(controller.Create)
	find := httptreemux.HandlerFunc(controller.Find)
	findByID := httptreemux.HandlerFunc(controller.FindByID)
	upsert := httptreemux.HandlerFunc(controller.Upsert)
	deleteAll := httptreemux.HandlerFunc(controller.DeleteAll)
	deleteByID := httptreemux.HandlerFunc(controller.DeleteByID)
	related := httptreemux.HandlerFunc(controller.Related)
	relatedOne := httptreemux.HandlerFunc(controller.RelatedOne)

	routeDir.Add(key.For("Create"), &interfaces.Route{Method: "POST", Path: "/users", Handler: &create, Visible: true}, true)
	routeDir.Add(key.For("Find"), &interfaces.Route{Method: "GET", Path: "/users", Handler: &find, Visible: true}, true)
	routeDir.Add(key.For("FindByID"), &interfaces.Route{Method: "GET", Path: "/users/:id", Handler: &findByID, Visible: true}, true)
	routeDir.Add(key.For("Upsert"), &interfaces.Route{Method: "PUT", Path: "/users", Handler: &upsert, Visible: true}, true)
	routeDir.Add(key.For("DeleteAll"), &interfaces.Route{Method: "DELETE", Path: "/users", Handler: &deleteAll, Visible: true}, true)
	routeDir.Add(key.For("DeleteByID"), &interfaces.Route{Method: "DELETE", Path: "/users/:id", Handler: &deleteByID, Visible: true}, true)

	routeDir.Add(key.For("CreateRelated"), &interfaces.Route{Method: "POST", Path: "/users/:pk/:related", Handler: &related, Visible: true}, false)
	routeDir.Add(key.For("FindRelated"), &interfaces.Route{Method: "GET", Path: "/users/:pk/:related", Handler: &related, Visible: true}, false)
	routeDir.Add(key.For("FindByIDRelated"), &interfaces.Route{Method: "GET", Path: "/users/:pk/:related/:fk", Handler: &relatedOne, Visible: true}, false)
	routeDir.Add(key.For("UpsertRelated"), &interfaces.Route{Method: "PUT", Path: "/users/:pk/:related", Handler: &related, Visible: true}, false)
	routeDir.Add(key.For("DeleteAllRelated"), &interfaces.Route{Method: "DELETE", Path: "/users/:pk/:related", Handler: &related, Visible: true}, false)
	routeDir.Add(key.For("DeleteByIDRelated"), &interfaces.Route{Method: "DELETE", Path: "/users/:pk/:related/:fk", Handler: &relatedOne, Visible: true}, false)
}
