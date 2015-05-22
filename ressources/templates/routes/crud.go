package interfaces

import (
	"github.com/Solher/auth-scaffold/ressources/templates"
	"github.com/clipperhouse/typewriter"
)

func init() {
	imports := []typewriter.ImportSpec{}

	err := typewriter.Register(templates.NewWrite("access", slice, imports))
	if err != nil {
		panic(err)
	}
}

var slice = typewriter.TemplateSlice{
	routes,
}

var routes = &typewriter.Template{
	Name: "Access",
	Text: `
	func setAccessOptions(routeDir interfaces.RouteDirectory, permissionDir usecases.PermissionDirectory, controller *{{.Type}}Ctrl) {
		key := interfaces.NewDirectoryKey(controller)
		create := httptreemux.HandlerFunc(controller.Create)
		find := httptreemux.HandlerFunc(controller.Find)
		findByID := httptreemux.HandlerFunc(controller.FindByID)
		upsert := httptreemux.HandlerFunc(controller.Upsert)
		deleteAll := httptreemux.HandlerFunc(controller.DeleteAll)
		deleteByID := httptreemux.HandlerFunc(controller.DeleteByID)

  	routeDir[key.For("Create")] = interfaces.Route{Method: "POST", Path: "/{{.Name}}s", Handler: &create}
  	routeDir[key.For("Find")] = interfaces.Route{Method: "GET", Path: "/{{.Name}}s", Handler: &find}
  	routeDir[key.For("FindByID")] = interfaces.Route{Method: "GET", Path: "/{{.Name}}s/:id", Handler: &findByID}
  	routeDir[key.For("Upsert")] = interfaces.Route{Method: "PUT", Path: "/{{.Name}}s", Handler: &upsert}
  	routeDir[key.For("DeleteAll")] = interfaces.Route{Method: "DELETE", Path: "/{{.Name}}s", Handler: &deleteAll}
  	routeDir[key.For("DeleteByID")] = interfaces.Route{Method: "DELETE", Path: "/{{.Name}}s/:id", Handler: &deleteByID}

		permissions := permissionDir["admin"]
		permissions.GrantAll()
		permissions = permissionDir["authenticated"]
		permissions.GrantAll()
		permissions = permissionDir["guest"]
	}
`}
