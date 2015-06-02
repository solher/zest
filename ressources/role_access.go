// Generated by: main
// TypeWriter: access
// Directive: +gen on Role

package ressources

import (
	"github.com/Solher/auth-scaffold/interfaces"
	"github.com/dimfeld/httptreemux"
)

func setRoleAccess(routeDir *interfaces.RouteDirectory, controller *RoleCtrl) {
	key := interfaces.NewDirectoryKey("roles")
	create := httptreemux.HandlerFunc(controller.Create)
	find := httptreemux.HandlerFunc(controller.Find)
	findByID := httptreemux.HandlerFunc(controller.FindByID)
	upsert := httptreemux.HandlerFunc(controller.Upsert)
	updateByID := httptreemux.HandlerFunc(controller.UpdateByID)
	deleteAll := httptreemux.HandlerFunc(controller.DeleteAll)
	deleteByID := httptreemux.HandlerFunc(controller.DeleteByID)
	related := httptreemux.HandlerFunc(controller.Related)
	relatedOne := httptreemux.HandlerFunc(controller.RelatedOne)

	routeDir.Add(key.For("Create"), &interfaces.Route{Method: "POST", Path: "/roles", Handler: &create, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("Find"), &interfaces.Route{Method: "GET", Path: "/roles", Handler: &find, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("FindByID"), &interfaces.Route{Method: "GET", Path: "/roles/:id", Handler: &findByID, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("Upsert"), &interfaces.Route{Method: "PUT", Path: "/roles", Handler: &upsert, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("UpdateByID"), &interfaces.Route{Method: "PUT", Path: "/roles/:id", Handler: &updateByID, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("DeleteAll"), &interfaces.Route{Method: "DELETE", Path: "/roles", Handler: &deleteAll, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("DeleteByID"), &interfaces.Route{Method: "DELETE", Path: "/roles/:id", Handler: &deleteByID, Visible: true, CheckPermissions: true})

	routeDir.Add(key.For("CreateRelated"), &interfaces.Route{Method: "POST", Path: "/roles/:pk/:related", Handler: &related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("FindRelated"), &interfaces.Route{Method: "GET", Path: "/roles/:pk/:related", Handler: &related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("FindByIDRelated"), &interfaces.Route{Method: "GET", Path: "/roles/:pk/:related/:fk", Handler: &relatedOne, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("UpsertRelated"), &interfaces.Route{Method: "PUT", Path: "/roles/:pk/:related", Handler: &related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("UpsertRelated"), &interfaces.Route{Method: "PUT", Path: "/roles/:pk/:related/:fk", Handler: &relatedOne, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("DeleteAllRelated"), &interfaces.Route{Method: "DELETE", Path: "/roles/:pk/:related", Handler: &related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("DeleteByIDRelated"), &interfaces.Route{Method: "DELETE", Path: "/roles/:pk/:related/:fk", Handler: &relatedOne, Visible: true, CheckPermissions: false})
}
