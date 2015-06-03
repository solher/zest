package ressources

import (
	"github.com/Solher/zest/usecases"
	"github.com/dimfeld/httptreemux"
)

func setAccountRoutes(routeDir *usecases.RouteDirectory, controller *AccountCtrl) {
	key := usecases.NewDirectoryKey("accounts")
	signin := httptreemux.HandlerFunc(controller.Signin)
	signout := httptreemux.HandlerFunc(controller.Signout)
	signup := httptreemux.HandlerFunc(controller.Signup)
	current := httptreemux.HandlerFunc(controller.Current)
	related := httptreemux.HandlerFunc(controller.Related)
	relatedOne := httptreemux.HandlerFunc(controller.RelatedOne)

	routeDir.Add(key.For("Signin"), &usecases.Route{Method: "POST", Path: "/accounts/signin", Handler: &signin, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("Signout"), &usecases.Route{Method: "POST", Path: "/accounts/signout", Handler: &signout, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("Signup"), &usecases.Route{Method: "POST", Path: "/accounts/signup", Handler: &signup, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("Current"), &usecases.Route{Method: "GET", Path: "/accounts/me", Handler: &current, Visible: true, CheckPermissions: true})

	routeDir.Add(key.For("CreateRelated"), &usecases.Route{Method: "POST", Path: "/accounts/me/:related", Handler: &related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("FindRelated"), &usecases.Route{Method: "GET", Path: "/accounts/me/:related", Handler: &related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("FindByIDRelated"), &usecases.Route{Method: "GET", Path: "/accounts/me/:related/:fk", Handler: &relatedOne, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("UpsertRelated"), &usecases.Route{Method: "PUT", Path: "/accounts/me/:related", Handler: &related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("DeleteAllRelated"), &usecases.Route{Method: "DELETE", Path: "/accounts/me/:related", Handler: &related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("DeleteByIDRelated"), &usecases.Route{Method: "DELETE", Path: "/accounts/me/:related/:fk", Handler: &relatedOne, Visible: true, CheckPermissions: false})
}
