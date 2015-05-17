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
	func set{{.Type}}AccessOptions(routesDir interfaces.RouteDirectory, permissionDir usecases.PermissionDirectory, controller *{{.Type}}Ctrl) {
		key := interfaces.NewDirectoryKey(controller)
		create := httprouter.Handle(controller.Create)
		find := httprouter.Handle(controller.Find)
		findByID := httprouter.Handle(controller.FindByID)
		upsert := httprouter.Handle(controller.Upsert)
		deleteAll := httprouter.Handle(controller.DeleteAll)
		deleteByID := httprouter.Handle(controller.DeleteByID)

  	routesDir[key.For("Create")] = interfaces.Route{Method: "POST", Path: "/{{.Name}}s", Handler: &create}
  	routesDir[key.For("Find")] = interfaces.Route{Method: "GET", Path: "/{{.Name}}s", Handler: &find}
  	routesDir[key.For("FindByID")] = interfaces.Route{Method: "GET", Path: "/{{.Name}}s/:id", Handler: &findByID}
  	routesDir[key.For("Upsert")] = interfaces.Route{Method: "PUT", Path: "/{{.Name}}s", Handler: &upsert}
  	routesDir[key.For("DeleteAll")] = interfaces.Route{Method: "DELETE", Path: "/{{.Name}}s", Handler: &deleteAll}
  	routesDir[key.For("DeleteByID")] = interfaces.Route{Method: "DELETE", Path: "/{{.Name}}s/:id", Handler: &deleteByID}

		permissions := permissionDir["admin"]
		permissions.Add(&create).Add(&find).Add(&findByID).Add(&upsert).Add(&deleteAll).Add(&deleteByID)
		permissions = permissionDir["authenticated"]
		permissions.Add(&create).Add(&find).Add(&findByID).Add(&upsert).Add(&deleteAll).Add(&deleteByID)
		permissions = permissionDir["guest"]
	}
`}
