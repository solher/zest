package ressources

import "github.com/solher/zest/usecases"

func setAccountRoutes(routeDir *usecases.RouteDirectory, controller *AccountCtrl) {
	key := usecases.NewDirectoryKey("accounts")

	routeDir.Add(key.For("Signin"), &usecases.Route{Method: "POST", Path: "/accounts/signin", Handler: controller.Signin, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("Signout"), &usecases.Route{Method: "POST", Path: "/accounts/signout", Handler: controller.Signout, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("Signup"), &usecases.Route{Method: "POST", Path: "/accounts/signup", Handler: controller.Signup, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("Current"), &usecases.Route{Method: "GET", Path: "/accounts/me", Handler: controller.Current, Visible: true, CheckPermissions: true})

	routeDir.Add(key.For("CreateRelated"), &usecases.Route{Method: "POST", Path: "/accounts/me/:related", Handler: controller.Related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("FindRelated"), &usecases.Route{Method: "GET", Path: "/accounts/me/:related", Handler: controller.Related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("FindByIDRelated"), &usecases.Route{Method: "GET", Path: "/accounts/me/:related/:fk", Handler: controller.RelatedOne, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("UpsertRelated"), &usecases.Route{Method: "PUT", Path: "/accounts/me/:related", Handler: controller.Related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("DeleteAllRelated"), &usecases.Route{Method: "DELETE", Path: "/accounts/me/:related", Handler: controller.Related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("DeleteByIDRelated"), &usecases.Route{Method: "DELETE", Path: "/accounts/me/:related/:fk", Handler: controller.RelatedOne, Visible: true, CheckPermissions: false})
}
