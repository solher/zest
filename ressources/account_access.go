package ressources

import (
	"github.com/Solher/auth-scaffold/interfaces"
	"github.com/Solher/auth-scaffold/usecases"
	"github.com/dimfeld/httptreemux"
)

func setAccountAccessOptions(routeDir *interfaces.RouteDirectory, permissionDir usecases.PermissionDirectory, controller *AccountCtrl) {
	key := interfaces.NewDirectoryKey("accounts")
	signin := httptreemux.HandlerFunc(controller.Signin)
	signout := httptreemux.HandlerFunc(controller.Signout)
	signup := httptreemux.HandlerFunc(controller.Signup)
	current := httptreemux.HandlerFunc(controller.Current)
	related := httptreemux.HandlerFunc(controller.Related)
	relatedOne := httptreemux.HandlerFunc(controller.RelatedOne)

	routeDir.Add(key.For("Signin"), &interfaces.Route{Method: "POST", Path: "/accounts/signin", Handler: &signin, Visible: true}, true)
	routeDir.Add(key.For("Signout"), &interfaces.Route{Method: "POST", Path: "/accounts/signout", Handler: &signout, Visible: true}, true)
	routeDir.Add(key.For("Signup"), &interfaces.Route{Method: "POST", Path: "/accounts/signup", Handler: &signup, Visible: true}, true)
	routeDir.Add(key.For("Me"), &interfaces.Route{Method: "GET", Path: "/accounts/me", Handler: &current, Visible: true}, true)

	routeDir.Add(key.For("CreateRelated"), &interfaces.Route{Method: "POST", Path: "/accounts/me/:related", Handler: &related, Visible: true}, false)
	routeDir.Add(key.For("FindRelated"), &interfaces.Route{Method: "GET", Path: "/accounts/me/:related", Handler: &related, Visible: true}, false)
	routeDir.Add(key.For("FindByIDRelated"), &interfaces.Route{Method: "GET", Path: "/accounts/me/:related/:fk", Handler: &relatedOne, Visible: true}, false)
	routeDir.Add(key.For("UpsertRelated"), &interfaces.Route{Method: "PUT", Path: "/accounts/me/:related", Handler: &related, Visible: true}, false)
	routeDir.Add(key.For("DeleteAllRelated"), &interfaces.Route{Method: "DELETE", Path: "/accounts/me/:related", Handler: &related, Visible: true}, false)
	routeDir.Add(key.For("DeleteByIDRelated"), &interfaces.Route{Method: "DELETE", Path: "/accounts/me/:related/:fk", Handler: &relatedOne, Visible: true}, false)

	permissions := permissionDir["admin"]
	permissions.Add(&signin).Add(&signout).Add(&signup).Add(&current).Add(&related).Add(&relatedOne)
	permissions = permissionDir["authenticated"]
	permissions.Add(&signin).Add(&signout).Add(&signup).Add(&current).Add(&related).Add(&relatedOne)
	permissions = permissionDir["guest"]
	permissions.Add(&signin).Add(&signout).Add(&signup).Add(&current)
}
