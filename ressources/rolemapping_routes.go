// Generated by: main
// TypeWriter: routes
// Directive: +gen on RoleMapping

package ressources

import "github.com/Solher/zest/usecases"

func setRoleMappingRoutes(routeDir *usecases.RouteDirectory, controller *RoleMappingCtrl) {
	key := usecases.NewDirectoryKey("rolemappings")

	routeDir.Add(key.For("Create"), &usecases.Route{Method: "POST", Path: "/rolemappings", Handler: controller.Create, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("Find"), &usecases.Route{Method: "GET", Path: "/rolemappings", Handler: controller.Find, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("FindByID"), &usecases.Route{Method: "GET", Path: "/rolemappings/:id", Handler: controller.FindByID, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("Upsert"), &usecases.Route{Method: "PUT", Path: "/rolemappings", Handler: controller.Upsert, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("UpdateByID"), &usecases.Route{Method: "PUT", Path: "/rolemappings/:id", Handler: controller.UpdateByID, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("DeleteAll"), &usecases.Route{Method: "DELETE", Path: "/rolemappings", Handler: controller.DeleteAll, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("DeleteByID"), &usecases.Route{Method: "DELETE", Path: "/rolemappings/:id", Handler: controller.DeleteByID, Visible: true, CheckPermissions: true})

	routeDir.Add(key.For("CreateRelated"), &usecases.Route{Method: "POST", Path: "/rolemappings/:pk/:related", Handler: controller.Related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("FindRelated"), &usecases.Route{Method: "GET", Path: "/rolemappings/:pk/:related", Handler: controller.Related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("FindByIDRelated"), &usecases.Route{Method: "GET", Path: "/rolemappings/:pk/:related/:fk", Handler: controller.RelatedOne, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("UpsertRelated"), &usecases.Route{Method: "PUT", Path: "/rolemappings/:pk/:related", Handler: controller.Related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("UpsertRelated"), &usecases.Route{Method: "PUT", Path: "/rolemappings/:pk/:related/:fk", Handler: controller.RelatedOne, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("DeleteAllRelated"), &usecases.Route{Method: "DELETE", Path: "/rolemappings/:pk/:related", Handler: controller.Related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("DeleteByIDRelated"), &usecases.Route{Method: "DELETE", Path: "/rolemappings/:pk/:related/:fk", Handler: controller.RelatedOne, Visible: true, CheckPermissions: false})
}
