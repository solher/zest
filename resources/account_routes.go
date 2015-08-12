package resources

import "github.com/solher/zest/usecases"

func setAccountRoutes(routeDir *usecases.RouteDirectory, controller *AccountCtrl) {
	key := usecases.NewDirectoryKey("accounts")

	routeDir.Add(key.For("Signin"), &usecases.Route{Method: "POST", Path: "/accounts/signin", Handler: controller.Signin, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("Signout"), &usecases.Route{Method: "POST", Path: "/accounts/signout", Handler: controller.Signout, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("Signup"), &usecases.Route{Method: "POST", Path: "/accounts/signup", Handler: controller.Signup, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("Current"), &usecases.Route{Method: "GET", Path: "/accounts/me", Handler: controller.Current, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("DeleteCurrent"), &usecases.Route{Method: "DELETE", Path: "/accounts/me", Handler: controller.DeleteCurrent, Visible: true, CheckPermissions: true})

	routeDir.Add(key.For("Create"), &usecases.Route{Method: "POST", Path: "/accounts", Handler: controller.Create, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("Find"), &usecases.Route{Method: "GET", Path: "/accounts", Handler: controller.Find, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("FindByID"), &usecases.Route{Method: "GET", Path: "/accounts/:id", Handler: controller.FindByID, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("Upsert"), &usecases.Route{Method: "PUT", Path: "/accounts", Handler: controller.Upsert, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("UpdateByID"), &usecases.Route{Method: "PUT", Path: "/accounts/:id", Handler: controller.UpdateByID, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("DeleteAll"), &usecases.Route{Method: "DELETE", Path: "/accounts", Handler: controller.DeleteAll, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("DeleteByID"), &usecases.Route{Method: "DELETE", Path: "/accounts/:id", Handler: controller.DeleteByID, Visible: true, CheckPermissions: true})

	routeDir.Add(key.For("UpdatePasswordRelated"), &usecases.Route{Method: "POST", Path: "/accounts/:pk/users/:fk/updatePassword", Handler: controller.UpdatePasswordRelated, Visible: true, CheckPermissions: false})

	routeDir.Add(key.For("CreateRelated"), &usecases.Route{Method: "POST", Path: "/accounts/:pk/:related", Handler: controller.CreateRelated, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("FindRelated"), &usecases.Route{Method: "GET", Path: "/accounts/:pk/:related", Handler: controller.FindRelated, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("FindByIDRelated"), &usecases.Route{Method: "GET", Path: "/accounts/:pk/:related/:fk", Handler: controller.FindByIDRelated, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("UpsertRelated"), &usecases.Route{Method: "PUT", Path: "/accounts/:pk/:related", Handler: controller.UpsertRelated, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("UpdateByIDRelated"), &usecases.Route{Method: "PUT", Path: "/accounts/:pk/:related/:fk", Handler: controller.UpdateByIDRelated, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("DeleteAllRelated"), &usecases.Route{Method: "DELETE", Path: "/accounts/:pk/:related", Handler: controller.DeleteAllRelated, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("DeleteByIDRelated"), &usecases.Route{Method: "DELETE", Path: "/accounts/:pk/:related/:fk", Handler: controller.DeleteByIDRelated, Visible: true, CheckPermissions: false})
}
