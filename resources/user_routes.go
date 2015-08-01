package resources

import "github.com/solher/zest/usecases"

func setUserRoutes(routeDir *usecases.RouteDirectory, controller *UserCtrl) {
	key := usecases.NewDirectoryKey("users")

	routeDir.Add(key.For("UpdatePassword"), &usecases.Route{Method: "POST", Path: "/users/:id/updatePassword", Handler: controller.UpdatePassword, Visible: true, CheckPermissions: true})

	routeDir.Add(key.For("Create"), &usecases.Route{Method: "POST", Path: "/users", Handler: controller.Create, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("Find"), &usecases.Route{Method: "GET", Path: "/users", Handler: controller.Find, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("FindByID"), &usecases.Route{Method: "GET", Path: "/users/:id", Handler: controller.FindByID, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("Upsert"), &usecases.Route{Method: "PUT", Path: "/users", Handler: controller.Upsert, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("UpdateByID"), &usecases.Route{Method: "PUT", Path: "/users/:id", Handler: controller.UpdateByID, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("DeleteAll"), &usecases.Route{Method: "DELETE", Path: "/users", Handler: controller.DeleteAll, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("DeleteByID"), &usecases.Route{Method: "DELETE", Path: "/users/:id", Handler: controller.DeleteByID, Visible: true, CheckPermissions: true})

	routeDir.Add(key.For("CreateRelated"), &usecases.Route{Method: "POST", Path: "/users/:pk/:related", Handler: controller.CreateRelated, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("FindRelated"), &usecases.Route{Method: "GET", Path: "/users/:pk/:related", Handler: controller.FindRelated, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("FindByIDRelated"), &usecases.Route{Method: "GET", Path: "/users/:pk/:related/:fk", Handler: controller.FindByIDRelated, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("UpsertRelated"), &usecases.Route{Method: "PUT", Path: "/users/:pk/:related", Handler: controller.UpsertRelated, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("UpdateByIDRelated"), &usecases.Route{Method: "PUT", Path: "/users/:pk/:related/:fk", Handler: controller.UpdateByIDRelated, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("DeleteAllRelated"), &usecases.Route{Method: "DELETE", Path: "/users/:pk/:related", Handler: controller.DeleteAllRelated, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("DeleteByIDRelated"), &usecases.Route{Method: "DELETE", Path: "/users/:pk/:related/:fk", Handler: controller.DeleteByIDRelated, Visible: true, CheckPermissions: false})
}
