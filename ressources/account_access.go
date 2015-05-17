package ressources

import (
	"github.com/Solher/auth-scaffold/interfaces"
	"github.com/Solher/auth-scaffold/usecases"
	"github.com/julienschmidt/httprouter"
)

func setAccountAccessOptions(routesDir interfaces.RouteDirectory, permissionDir usecases.PermissionDirectory, controller *AccountCtrl) {
	key := interfaces.NewDirectoryKey(controller)
	signin := httprouter.Handle(controller.Signin)
	signout := httprouter.Handle(controller.Signout)
	signup := httprouter.Handle(controller.Signup)
	current := httprouter.Handle(controller.Current)

	routesDir[key.For("signin")] = interfaces.Route{Method: "POST", Path: "/accounts/signin", Handler: &signin}
	routesDir[key.For("signout")] = interfaces.Route{Method: "POST", Path: "/accounts/signout", Handler: &signout}
	routesDir[key.For("signup")] = interfaces.Route{Method: "POST", Path: "/accounts/signup", Handler: &signup}
	routesDir[key.For("current")] = interfaces.Route{Method: "GET", Path: "/accounts/current", Handler: &current}

	permissions := permissionDir["admin"]
	permissions.Add(&signin).Add(&signout).Add(&signup).Add(&current)
	permissions = permissionDir["authenticated"]
	permissions.Add(&signin).Add(&signout).Add(&signup).Add(&current)
	permissions = permissionDir["guest"]
	permissions.Add(&signin).Add(&signout).Add(&signup).Add(&current)
}
