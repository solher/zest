package controllers

import (
	"github.com/Solher/zest/ressources/templates"
	"github.com/clipperhouse/typewriter"
)

func init() {
	imports := []typewriter.ImportSpec{}

	err := typewriter.Register(templates.NewWrite("controller", slice, imports))
	if err != nil {
		panic(err)
	}
}

var slice = typewriter.TemplateSlice{
	controller,
	create,
	find,
	findByID,
	upsert,
	updateByID,
	deleteAll,
	deleteByID,
	related,
	relatedOne,
}

var controller = &typewriter.Template{
	Name: "Controller",
	Text: `
	type Abstract{{.Type}}Inter interface {
		Create({{.Name}}s []domain.{{.Type}}) ([]domain.{{.Type}}, error)
		CreateOne({{.Name}} *domain.{{.Type}}) (*domain.{{.Type}}, error)
		Find(context usecases.QueryContext) ([]domain.{{.Type}}, error)
		FindByID(id int, context usecases.QueryContext) (*domain.{{.Type}}, error)
		Upsert({{.Name}}s []domain.{{.Type}}, context usecases.QueryContext) ([]domain.{{.Type}}, error)
		UpsertOne({{.Name}} *domain.{{.Type}}, context usecases.QueryContext) (*domain.{{.Type}}, error)
		UpdateByID(id int, {{.Name}} *domain.{{.Type}},	context usecases.QueryContext) (*domain.{{.Type}}, error)
		DeleteAll(context usecases.QueryContext) error
		DeleteByID(id int, context usecases.QueryContext) error
	}

	type {{.Type}}Ctrl struct {
		interactor Abstract{{.Type}}Inter
		render     interfaces.AbstractRender
		routeDir   *usecases.RouteDirectory
	}

	func New{{.Type}}Ctrl(interactor Abstract{{.Type}}Inter, render interfaces.AbstractRender, routeDir *usecases.RouteDirectory) *{{.Type}}Ctrl {
		controller := &{{.Type}}Ctrl{interactor: interactor, render: render, routeDir: routeDir}

		if routeDir != nil {
			set{{.Type}}Routes(routeDir, controller)
		}

		return controller
	}
`}

var create = &typewriter.Template{
	Name: "Create",
	Text: `
	func (c *{{.Type}}Ctrl) Create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		{{.Name}} := &domain.{{.Type}}{}
		var {{.Name}}s []domain.{{.Type}}

		buffer, _ := ioutil.ReadAll(r.Body)

		err := json.Unmarshal(buffer, {{.Name}})
		if err != nil {
			err := json.Unmarshal(buffer, &{{.Name}}s)
			if err != nil {
				c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
				return
			}
		}

		lastRessource := interfaces.GetLastRessource(r)

		if {{.Name}}s == nil {
			{{.Name}}.SetRelatedID(lastRessource.IDKey, lastRessource.ID)
			{{.Name}}, err = c.interactor.CreateOne({{.Name}})
		} else {
			for i := range {{.Name}}s {
				(&{{.Name}}s[i]).SetRelatedID(lastRessource.IDKey, lastRessource.ID)
			}
			{{.Name}}s, err = c.interactor.Create({{.Name}}s)
		}

		if err != nil {
			switch err.(type) {
			case *internalerrors.ViolatedConstraint:
				c.render.JSONError(w, 422, apierrors.ViolatedConstraint, err)
			default:
				c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
			}
			return
		}

		if {{.Name}}s == nil {
			{{.Name}}.BeforeRender()
			c.render.JSON(w, http.StatusCreated, {{.Name}})
		} else {
			for i := range {{.Name}}s {
				(&{{.Name}}s[i]).BeforeRender()
			}
			c.render.JSON(w, http.StatusCreated, {{.Name}}s)
		}
	}
`}

var find = &typewriter.Template{
	Name: "Find",
	Text: `
	func (c *{{.Type}}Ctrl) Find(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		filter, err := interfaces.GetQueryFilter(r)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.FilterDecodingError, err)
			return
		}

		filter = interfaces.FilterIfLastRessource(r, filter)
		filter = interfaces.FilterIfOwnerRelations(r, filter)
		relations := interfaces.GetOwnerRelations(r)

		{{.Name}}s, err := c.interactor.Find(usecases.QueryContext{Filter:filter, OwnerRelations: relations})
		if err != nil {
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
			return
		}

		for i := range {{.Name}}s {
			(&{{.Name}}s[i]).BeforeRender()
		}
		c.render.JSON(w, http.StatusOK, {{.Name}}s)
	}
`}

var findByID = &typewriter.Template{
	Name: "FindByID",
	Text: `
	func (c *{{.Type}}Ctrl) FindByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.InvalidPathParams, err)
			return
		}

		filter, err := interfaces.GetQueryFilter(r)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.FilterDecodingError, err)
			return
		}

		filter = interfaces.FilterIfOwnerRelations(r, filter)
		relations := interfaces.GetOwnerRelations(r)

		{{.Name}}, err := c.interactor.FindByID(id, usecases.QueryContext{Filter:filter, OwnerRelations: relations})
		if err != nil {
			c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
			return
		}

		{{.Name}}.BeforeRender()
		c.render.JSON(w, http.StatusOK, {{.Name}})
	}
`}

var upsert = &typewriter.Template{
	Name: "Upsert",
	Text: `
	func (c *{{.Type}}Ctrl) Upsert(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		{{.Name}} := &domain.{{.Type}}{}
		var {{.Name}}s []domain.{{.Type}}

		buffer, _ := ioutil.ReadAll(r.Body)

		err := json.Unmarshal(buffer, {{.Name}})
		if err != nil {
			err := json.Unmarshal(buffer, &{{.Name}}s)
			if err != nil {
				c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
				return
			}
		}

		lastRessource := interfaces.GetLastRessource(r)
		filter := interfaces.FilterIfOwnerRelations(r, nil)
		relations := interfaces.GetOwnerRelations(r)

		if {{.Name}}s == nil {
			{{.Name}}.SetRelatedID(lastRessource.IDKey, lastRessource.ID)
			{{.Name}}, err = c.interactor.UpsertOne({{.Name}}, usecases.QueryContext{Filter:filter, OwnerRelations: relations})
		} else {
			for i := range {{.Name}}s {
				(&{{.Name}}s[i]).SetRelatedID(lastRessource.IDKey, lastRessource.ID)
			}
			{{.Name}}s, err = c.interactor.Upsert({{.Name}}s, usecases.QueryContext{Filter:filter, OwnerRelations: relations})
		}

		if err != nil {
			switch err.(type) {
			case *internalerrors.ViolatedConstraint:
				c.render.JSONError(w, 422, apierrors.ViolatedConstraint, err)
			default:
				c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
			}
			return
		}

		if {{.Name}}s == nil {
			{{.Name}}.BeforeRender()
			c.render.JSON(w, http.StatusCreated, {{.Name}})
		} else {
			for i := range {{.Name}}s {
				(&{{.Name}}s[i]).BeforeRender()
			}
			c.render.JSON(w, http.StatusCreated, {{.Name}}s)
		}
	}
`}

var updateByID = &typewriter.Template{
	Name: "UpdateByID",
	Text: `
	func (c *{{.Type}}Ctrl) UpdateByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.InvalidPathParams, err)
			return
		}

		{{.Name}} := &domain.{{.Type}}{}

		err = json.NewDecoder(r.Body).Decode({{.Name}})
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
			return
		}

		lastRessource := interfaces.GetLastRessource(r)
		filter := interfaces.FilterIfOwnerRelations(r, nil)
		relations := interfaces.GetOwnerRelations(r)

		{{.Name}}.SetRelatedID(lastRessource.IDKey, lastRessource.ID)
		{{.Name}}, err = c.interactor.UpdateByID(id, {{.Name}}, usecases.QueryContext{Filter:filter, OwnerRelations: relations})

		if err != nil {
			switch err.(type) {
			case *internalerrors.ViolatedConstraint:
				c.render.JSONError(w, 422, apierrors.ViolatedConstraint, err)
			default:
				c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
			}
			return
		}

		{{.Name}}.BeforeRender()
		c.render.JSON(w, http.StatusCreated, {{.Name}})
	}
`}

var deleteAll = &typewriter.Template{
	Name: "DeleteAll",
	Text: `
	func (c *{{.Type}}Ctrl) DeleteAll(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		filter, err := interfaces.GetQueryFilter(r)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.FilterDecodingError, err)
			return
		}

		filter = interfaces.FilterIfLastRessource(r, filter)
		filter = interfaces.FilterIfOwnerRelations(r, filter)
		relations := interfaces.GetOwnerRelations(r)

		err = c.interactor.DeleteAll(usecases.QueryContext{Filter:filter, OwnerRelations: relations})
		if err != nil {
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
			return
		}

		c.render.JSON(w, http.StatusNoContent, nil)
	}
`}

var deleteByID = &typewriter.Template{
	Name: "DeleteByID",
	Text: `
	func (c *{{.Type}}Ctrl) DeleteByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.InvalidPathParams, err)
			return
		}

		filter := interfaces.FilterIfOwnerRelations(r, nil)
		relations := interfaces.GetOwnerRelations(r)

		err = c.interactor.DeleteByID(id, usecases.QueryContext{Filter:filter, OwnerRelations: relations})
		if err != nil {
			c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
			return
		}

		c.render.JSON(w, http.StatusNoContent, nil)
	}
`}

var related = &typewriter.Template{
	Name: "Related",
	Text: `
	func (c *{{.Type}}Ctrl) Related(w http.ResponseWriter, r *http.Request, params map[string]string) {
		pk, err := strconv.Atoi(params["pk"])
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.InvalidPathParams, err)
			return
		}

		related := params["related"]
		key := usecases.NewDirectoryKey(related)

		var handler usecases.HandlerFunc
		switch r.Method {
		case "POST":
			handler = c.routeDir.Get(key.For("Create")).EffectiveHandler
		case "GET":
			handler = c.routeDir.Get(key.For("Find")).EffectiveHandler
		case "PUT":
			handler = c.routeDir.Get(key.For("Upsert")).EffectiveHandler
		case "DELETE":
			handler = c.routeDir.Get(key.For("DeleteAll")).EffectiveHandler
		}

		if handler == nil {
			c.render.JSON(w, http.StatusNotFound, nil)
			return
		}

		context.Set(r, "lastRessource", &interfaces.Ressource{Name: related, IDKey: "{{.Name}}ID", ID: pk})

		handler(w, r, params)
	}
`}

var relatedOne = &typewriter.Template{
	Name: "RelatedOne",
	Text: `
	func (c *{{.Type}}Ctrl) RelatedOne(w http.ResponseWriter, r *http.Request, params map[string]string) {
		params["id"] = params["fk"]

		related := params["related"]
		key := usecases.NewDirectoryKey(related)

		var handler usecases.HandlerFunc

		switch r.Method {
		case "GET":
			handler = c.routeDir.Get(key.For("FindByID")).EffectiveHandler
		case "DELETE":
			handler = c.routeDir.Get(key.For("DeleteByID")).EffectiveHandler
		}

		if handler == nil {
			c.render.JSON(w, http.StatusNotFound, nil)
			return
		}

		handler(w, r, params)
	}
`}
