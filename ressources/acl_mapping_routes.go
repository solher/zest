package ressources

import "github.com/solher/zest/usecases"

func setAclMappingRoutes(routeDir *usecases.RouteDirectory, controller *AclMappingCtrl) {
	key := usecases.NewDirectoryKey("aclMappings")

	routeDir.Add(key.For("Create"), &usecases.Route{Method: "POST", Path: "/aclMappings", Handler: controller.Create, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("Find"), &usecases.Route{Method: "GET", Path: "/aclMappings", Handler: controller.Find, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("FindByID"), &usecases.Route{Method: "GET", Path: "/aclMappings/:id", Handler: controller.FindByID, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("Upsert"), &usecases.Route{Method: "PUT", Path: "/aclMappings", Handler: controller.Upsert, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("UpdateByID"), &usecases.Route{Method: "PUT", Path: "/aclMappings/:id", Handler: controller.UpdateByID, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("DeleteAll"), &usecases.Route{Method: "DELETE", Path: "/aclMappings", Handler: controller.DeleteAll, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("DeleteByID"), &usecases.Route{Method: "DELETE", Path: "/aclMappings/:id", Handler: controller.DeleteByID, Visible: true, CheckPermissions: true})

	routeDir.Add(key.For("CreateRelated"), &usecases.Route{Method: "POST", Path: "/aclMappings/:pk/:related", Handler: controller.Related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("FindRelated"), &usecases.Route{Method: "GET", Path: "/aclMappings/:pk/:related", Handler: controller.Related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("FindByIDRelated"), &usecases.Route{Method: "GET", Path: "/aclMappings/:pk/:related/:fk", Handler: controller.RelatedOne, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("UpsertRelated"), &usecases.Route{Method: "PUT", Path: "/aclMappings/:pk/:related", Handler: controller.Related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("UpdateByIDRelated"), &usecases.Route{Method: "PUT", Path: "/aclMappings/:pk/:related/:fk", Handler: controller.RelatedOne, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("DeleteAllRelated"), &usecases.Route{Method: "DELETE", Path: "/aclMappings/:pk/:related", Handler: controller.Related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("DeleteByIDRelated"), &usecases.Route{Method: "DELETE", Path: "/aclMappings/:pk/:related/:fk", Handler: controller.RelatedOne, Visible: true, CheckPermissions: false})
}
