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
	func set{{.Type}}AccessOptions(routeDir *interfaces.RouteDirectory, controller *{{.Type}}Ctrl) {
		key := interfaces.NewDirectoryKey("{{.Name}}s")
		create := httptreemux.HandlerFunc(controller.Create)
		find := httptreemux.HandlerFunc(controller.Find)
		findByID := httptreemux.HandlerFunc(controller.FindByID)
		upsert := httptreemux.HandlerFunc(controller.Upsert)
		deleteAll := httptreemux.HandlerFunc(controller.DeleteAll)
		deleteByID := httptreemux.HandlerFunc(controller.DeleteByID)

		routeDir.Add(key.For("Create"), &interfaces.Route{Method: "POST", Path: "/{{.Name}}s", Handler: &create, Visible: true}, true)
  	routeDir.Add(key.For("Find"), &interfaces.Route{Method: "GET", Path: "/{{.Name}}s", Handler: &find, Visible: true}, true)
  	routeDir.Add(key.For("FindByID"), &interfaces.Route{Method: "GET", Path: "/{{.Name}}s/:id", Handler: &findByID, Visible: true}, true)
  	routeDir.Add(key.For("Upsert"), &interfaces.Route{Method: "PUT", Path: "/{{.Name}}s", Handler: &upsert, Visible: true}, true)
  	routeDir.Add(key.For("DeleteAll"), &interfaces.Route{Method: "DELETE", Path: "/{{.Name}}s", Handler: &deleteAll, Visible: true}, true)
  	routeDir.Add(key.For("DeleteByID"), &interfaces.Route{Method: "DELETE", Path: "/{{.Name}}s/:id", Handler: &deleteByID, Visible: true}, true)
	}
`}
