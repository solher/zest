package interfaces

import (
	"github.com/Solher/auth-scaffold/ressources/templates"
	"github.com/clipperhouse/typewriter"
)

func init() {
	imports := []typewriter.ImportSpec{}

	err := typewriter.Register(templates.NewWrite("routes", slice, imports))
	if err != nil {
		panic(err)
	}
}

var slice = typewriter.TemplateSlice{
	routes,
}

var routes = &typewriter.Template{
	Name: "Routes",
	Text: `
  func add{{.Type}}Routes(routesDir interfaces.RouteDirectory, controller *{{.Type}}Ctrl) {
  	key := interfaces.NewDirectoryKey(controller)

  	routesDir[key.For("Create")] = interfaces.Route{Method: "POST", Path: "/{{.Name}}s", Handler: controller.Create}
  	routesDir[key.For("Find")] = interfaces.Route{Method: "GET", Path: "/{{.Name}}s", Handler: controller.Find}
  	routesDir[key.For("FindByID")] = interfaces.Route{Method: "GET", Path: "/{{.Name}}s/:id", Handler: controller.FindByID}
  	routesDir[key.For("Upsert")] = interfaces.Route{Method: "PUT", Path: "/{{.Name}}s", Handler: controller.Upsert}
  	routesDir[key.For("DeleteAll")] = interfaces.Route{Method: "DELETE", Path: "/{{.Name}}s", Handler: controller.DeleteAll}
  	routesDir[key.For("DeleteByID")] = interfaces.Route{Method: "DELETE", Path: "/{{.Name}}s/:id", Handler: controller.DeleteByID}
  }
`}
