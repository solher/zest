package ressources

import (
	"github.com/Solher/zest/interfaces"
	"github.com/dimfeld/httptreemux"
)

func setAccountRoutes(routeDir *interfaces.RouteDirectory, controller *AccountCtrl) {
	key := interfaces.NewDirectoryKey("accounts")
	signin := httptreemux.HandlerFunc(controller.Signin)
	signout := httptreemux.HandlerFunc(controller.Signout)
	signup := httptreemux.HandlerFunc(controller.Signup)
	current := httptreemux.HandlerFunc(controller.Current)
	related := httptreemux.HandlerFunc(controller.Related)
	relatedOne := httptreemux.HandlerFunc(controller.RelatedOne)

	routeDir.Add(key.For("Signin"), &interfaces.Route{Method: "POST", Path: "/accounts/signin", Handler: &signin, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("Signout"), &interfaces.Route{Method: "POST", Path: "/accounts/signout", Handler: &signout, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("Signup"), &interfaces.Route{Method: "POST", Path: "/accounts/signup", Handler: &signup, Visible: true, CheckPermissions: true})
	routeDir.Add(key.For("Current"), &interfaces.Route{Method: "GET", Path: "/accounts/me", Handler: &current, Visible: true, CheckPermissions: true})

	routeDir.Add(key.For("CreateRelated"), &interfaces.Route{Method: "POST", Path: "/accounts/me/:related", Handler: &related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("FindRelated"), &interfaces.Route{Method: "GET", Path: "/accounts/me/:related", Handler: &related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("FindByIDRelated"), &interfaces.Route{Method: "GET", Path: "/accounts/me/:related/:fk", Handler: &relatedOne, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("UpsertRelated"), &interfaces.Route{Method: "PUT", Path: "/accounts/me/:related", Handler: &related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("DeleteAllRelated"), &interfaces.Route{Method: "DELETE", Path: "/accounts/me/:related", Handler: &related, Visible: true, CheckPermissions: false})
	routeDir.Add(key.For("DeleteByIDRelated"), &interfaces.Route{Method: "DELETE", Path: "/accounts/me/:related/:fk", Handler: &relatedOne, Visible: true, CheckPermissions: false})
}
