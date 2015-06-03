// Generated by: main
// TypeWriter: routes
// Directive: +gen on Role

package ressources

import (
	"github.com/Solher/zest/usecases"
	"github.com/dimfeld/httptreemux"
)

func setRoleRoutes(routeDir *usecases.RouteDirectory, controller *RoleCtrl) {
	key := usecases.NewDirectoryKey("roles")
	create := httptreemux.HandlerFunc(controller.Create)
	find := httptreemux.HandlerFunc(controller.Find)
	findByID := httptreemux.HandlerFunc(controller.FindByID)
	upsert := httptreemux.HandlerFunc(controller.Upsert)
	updateByID := httptreemux.HandlerFunc(controller.UpdateByID)
	deleteAll := httptreemux.HandlerFunc(controller.DeleteAll)
	deleteByID := httptreemux.HandlerFunc(controller.DeleteByID)
	related := httptreemux.HandlerFunc(controller.Related)
	relatedOne := httptreemux.HandlerFunc(controller.RelatedOne)

	routeDir.Add(key.For("Create"), &usecases.Route{Method: "POST", Path: "/roles", Handler: &create, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("Find"), &usecases.Route{Method: "GET", Path: "/roles", Handler: &find, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("FindByID"), &usecases.Route{Method: "GET", Path: "/roles/:id", Handler: &findByID, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("Upsert"), &usecases.Route{Method: "PUT", Path: "/roles", Handler: &upsert, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("UpdateByID"), &usecases.Route{Method: "PUT", Path: "/roles/:id", Handler: &updateByID, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("DeleteAll"), &usecases.Route{Method: "DELETE", Path: "/roles", Handler: &deleteAll, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("DeleteByID"), &usecases.Route{Method: "DELETE", Path: "/roles/:id", Handler: &deleteByID, Visible: true, CheckPermissions: true})

	routeDir.Add(key.For("CreateRelated"), &usecases.Route{Method: "POST", Path: "/roles/:pk/:related", Handler: &related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("FindRelated"), &usecases.Route{Method: "GET", Path: "/roles/:pk/:related", Handler: &related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("FindByIDRelated"), &usecases.Route{Method: "GET", Path: "/roles/:pk/:related/:fk", Handler: &relatedOne, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("UpsertRelated"), &usecases.Route{Method: "PUT", Path: "/roles/:pk/:related", Handler: &related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("UpsertRelated"), &usecases.Route{Method: "PUT", Path: "/roles/:pk/:related/:fk", Handler: &relatedOne, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("DeleteAllRelated"), &usecases.Route{Method: "DELETE", Path: "/roles/:pk/:related", Handler: &related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("DeleteByIDRelated"), &usecases.Route{Method: "DELETE", Path: "/roles/:pk/:related/:fk", Handler: &relatedOne, Visible: true, CheckPermissions: false})
}
