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
		create := httptreemux.HandlerFunc(controller.Create)
		find := httptreemux.HandlerFunc(controller.Find)
		findByID := httptreemux.HandlerFunc(controller.FindByID)
		upsert := httptreemux.HandlerFunc(controller.Upsert)
		updateByID := httptreemux.HandlerFunc(controller.UpdateByID)
		deleteAll := httptreemux.HandlerFunc(controller.DeleteAll)
		deleteByID := httptreemux.HandlerFunc(controller.DeleteByID)
		related := httptreemux.HandlerFunc(controller.Related)
		relatedOne := httptreemux.HandlerFunc(controller.RelatedOne)

		routeDir.Add(key.For("Create"), &usecases.Route{Method: "POST", Path: "/{{.Name}}s", Handler: &create, Visible: true, CheckPermissions: true})
  	routeDir.Add(key.For("Find"), &usecases.Route{Method: "GET", Path: "/{{.Name}}s", Handler: &find, Visible: true, CheckPermissions: true})
  	routeDir.Add(key.For("FindByID"), &usecases.Route{Method: "GET", Path: "/{{.Name}}s/:id", Handler: &findByID, Visible: true, CheckPermissions: true})
  	routeDir.Add(key.For("Upsert"), &usecases.Route{Method: "PUT", Path: "/{{.Name}}s", Handler: &upsert, Visible: true, CheckPermissions: true})
		routeDir.Add(key.For("UpdateByID"), &usecases.Route{Method: "PUT", Path: "/{{.Name}}s/:id", Handler: &updateByID, Visible: true, CheckPermissions: true})
  	routeDir.Add(key.For("DeleteAll"), &usecases.Route{Method: "DELETE", Path: "/{{.Name}}s", Handler: &deleteAll, Visible: true, CheckPermissions: true})
  	routeDir.Add(key.For("DeleteByID"), &usecases.Route{Method: "DELETE", Path: "/{{.Name}}s/:id", Handler: &deleteByID, Visible: true, CheckPermissions: true})

		routeDir.Add(key.For("CreateRelated"), &usecases.Route{Method: "POST", Path: "/{{.Name}}s/:pk/:related", Handler: &related, Visible: true, CheckPermissions: false})
		routeDir.Add(key.For("FindRelated"), &usecases.Route{Method: "GET", Path: "/{{.Name}}s/:pk/:related", Handler: &related, Visible: true, CheckPermissions: false})
		routeDir.Add(key.For("FindByIDRelated"), &usecases.Route{Method: "GET", Path: "/{{.Name}}s/:pk/:related/:fk", Handler: &relatedOne, Visible: true, CheckPermissions: false})
		routeDir.Add(key.For("UpsertRelated"), &usecases.Route{Method: "PUT", Path: "/{{.Name}}s/:pk/:related", Handler: &related, Visible: true, CheckPermissions: false})
		routeDir.Add(key.For("UpsertRelated"), &usecases.Route{Method: "PUT", Path: "/{{.Name}}s/:pk/:related/:fk", Handler: &relatedOne, Visible: true, CheckPermissions: false})
		routeDir.Add(key.For("DeleteAllRelated"), &usecases.Route{Method: "DELETE", Path: "/{{.Name}}s/:pk/:related", Handler: &related, Visible: true, CheckPermissions: false})
		routeDir.Add(key.For("DeleteByIDRelated"), &usecases.Route{Method: "DELETE", Path: "/{{.Name}}s/:pk/:related/:fk", Handler: &relatedOne, Visible: true, CheckPermissions: false})
	}
`}
