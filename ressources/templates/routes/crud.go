package interfaces

import (
	"github.com/Solher/zest/ressources/templates"
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
	func set{{.Type}}Routes(routeDir *usecases.RouteDirectory, controller *{{.Type}}Ctrl) {
		key := usecases.NewDirectoryKey("{{.Name}}s")

		routeDir.Add(key.For("Create"), &usecases.Route{Method: "POST", Path: "/{{.Name}}s", Handler: controller.Create, Visible: true, CheckPermissions: true})
  	routeDir.Add(key.For("Find"), &usecases.Route{Method: "GET", Path: "/{{.Name}}s", Handler: controller.Find, Visible: true, CheckPermissions: true})
  	routeDir.Add(key.For("FindByID"), &usecases.Route{Method: "GET", Path: "/{{.Name}}s/:id", Handler: controller.FindByID, Visible: true, CheckPermissions: true})
  	routeDir.Add(key.For("Upsert"), &usecases.Route{Method: "PUT", Path: "/{{.Name}}s", Handler: controller.Upsert, Visible: true, CheckPermissions: true})
		routeDir.Add(key.For("UpdateByID"), &usecases.Route{Method: "PUT", Path: "/{{.Name}}s/:id", Handler: controller.UpdateByID, Visible: true, CheckPermissions: true})
  	routeDir.Add(key.For("DeleteAll"), &usecases.Route{Method: "DELETE", Path: "/{{.Name}}s", Handler: controller.DeleteAll, Visible: true, CheckPermissions: true})
  	routeDir.Add(key.For("DeleteByID"), &usecases.Route{Method: "DELETE", Path: "/{{.Name}}s/:id", Handler: controller.DeleteByID, Visible: true, CheckPermissions: true})

		routeDir.Add(key.For("CreateRelated"), &usecases.Route{Method: "POST", Path: "/{{.Name}}s/:pk/:related", Handler: controller.Related, Visible: true, CheckPermissions: false})
		routeDir.Add(key.For("FindRelated"), &usecases.Route{Method: "GET", Path: "/{{.Name}}s/:pk/:related", Handler: controller.Related, Visible: true, CheckPermissions: false})
		routeDir.Add(key.For("FindByIDRelated"), &usecases.Route{Method: "GET", Path: "/{{.Name}}s/:pk/:related/:fk", Handler: controller.RelatedOne, Visible: true, CheckPermissions: false})
		routeDir.Add(key.For("UpsertRelated"), &usecases.Route{Method: "PUT", Path: "/{{.Name}}s/:pk/:related", Handler: controller.Related, Visible: true, CheckPermissions: false})
		routeDir.Add(key.For("UpsertRelated"), &usecases.Route{Method: "PUT", Path: "/{{.Name}}s/:pk/:related/:fk", Handler: controller.RelatedOne, Visible: true, CheckPermissions: false})
		routeDir.Add(key.For("DeleteAllRelated"), &usecases.Route{Method: "DELETE", Path: "/{{.Name}}s/:pk/:related", Handler: controller.Related, Visible: true, CheckPermissions: false})
		routeDir.Add(key.For("DeleteByIDRelated"), &usecases.Route{Method: "DELETE", Path: "/{{.Name}}s/:pk/:related/:fk", Handler: controller.RelatedOne, Visible: true, CheckPermissions: false})
	}
`}
