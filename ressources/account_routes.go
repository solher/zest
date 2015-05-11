package ressources

import "github.com/Solher/auth-scaffold/interfaces"

func addAccountRoutes(routesDir interfaces.RouteDirectory, controller *AccountCtrl) {
	key := interfaces.NewDirectoryKey(controller)

	routesDir[key.For("Signin")] = interfaces.Route{Method: "POST", Path: "/accounts/signin", Handler: controller.Signin}
	routesDir[key.For("Signout")] = interfaces.Route{Method: "POST", Path: "/accounts/signout", Handler: controller.Signout}
	routesDir[key.For("Signup")] = interfaces.Route{Method: "POST", Path: "/accounts/signup", Handler: controller.Signup}
	routesDir[key.For("Current")] = interfaces.Route{Method: "GET", Path: "/accounts/current", Handler: controller.Current}
}
