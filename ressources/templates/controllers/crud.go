package controllers

import (
	"github.com/Solher/auth-scaffold/ressources/templates"
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
		Find(filter *interfaces.Filter, ownerRelations []domain.Relation) ([]domain.{{.Type}}, error)
		FindByID(id int, filter *interfaces.Filter, ownerRelations []domain.Relation) (*domain.{{.Type}}, error)
		Upsert({{.Name}}s []domain.{{.Type}}, filter *interfaces.Filter, ownerRelations []domain.Relation) ([]domain.{{.Type}}, error)
		UpsertOne({{.Name}} *domain.{{.Type}}, filter *interfaces.Filter, ownerRelations []domain.Relation) (*domain.{{.Type}}, error)
		DeleteAll(filter *interfaces.Filter, ownerRelations []domain.Relation) error
		DeleteByID(id int, filter *interfaces.Filter, ownerRelations []domain.Relation) error
	}

	type {{.Type}}Ctrl struct {
		interactor Abstract{{.Type}}Inter
		render     interfaces.AbstractRender
		routeDir   *interfaces.RouteDirectory
	}

	func New{{.Type}}Ctrl(interactor Abstract{{.Type}}Inter, render interfaces.AbstractRender, routeDir *interfaces.RouteDirectory) *{{.Type}}Ctrl {
		controller := &{{.Type}}Ctrl{interactor: interactor, render: render, routeDir: routeDir}

		if routeDir != nil {
			set{{.Type}}Access(routeDir, controller)
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
			{{.Name}}.ScopeModel(lastRessource.ID)
			{{.Name}}, err = c.interactor.CreateOne({{.Name}})
		} else {
			for i := range {{.Name}}s {
				(&{{.Name}}s[i]).ScopeModel(lastRessource.ID)
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
			c.render.JSON(w, http.StatusCreated, {{.Name}})
		} else {
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

		{{.Name}}s, err := c.interactor.Find(filter, relations)
		if err != nil {
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
			return
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

		{{.Name}}, err := c.interactor.FindByID(id, filter, relations)
		if err != nil {
			c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
			return
		}

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
		ownerRelations := interfaces.GetOwnerRelations(r)

		if {{.Name}}s == nil {
			{{.Name}}.ScopeModel(lastRessource.ID)
			{{.Name}}, err = c.interactor.UpsertOne({{.Name}}, filter, ownerRelations)
		} else {
			for i := range {{.Name}}s {
				(&{{.Name}}s[i]).ScopeModel(lastRessource.ID)
			}
			{{.Name}}s, err = c.interactor.Upsert({{.Name}}s, filter, ownerRelations)
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
			c.render.JSON(w, http.StatusCreated, {{.Name}})
		} else {
			c.render.JSON(w, http.StatusCreated, {{.Name}}s)
		}
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

		err = c.interactor.DeleteAll(filter, relations)
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
		ownerRelations := interfaces.GetOwnerRelations(r)

		err = c.interactor.DeleteByID(id, filter, ownerRelations)
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
		key := interfaces.NewDirectoryKey(related)

		var handler *httptreemux.HandlerFunc
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

		(*handler)(w, r, params)
	}
`}

var relatedOne = &typewriter.Template{
	Name: "RelatedOne",
	Text: `
	func (c *{{.Type}}Ctrl) RelatedOne(w http.ResponseWriter, r *http.Request, params map[string]string) {
		params["id"] = params["fk"]

		related := params["related"]
		key := interfaces.NewDirectoryKey(related)

		var handler httptreemux.HandlerFunc

		switch r.Method {
		case "GET":
			handler = *c.routeDir.Get(key.For("FindByID")).EffectiveHandler
		case "DELETE":
			handler = *c.routeDir.Get(key.For("DeleteByID")).EffectiveHandler
		}

		if handler == nil {
			c.render.JSON(w, http.StatusNotFound, nil)
			return
		}

		handler(w, r, params)
	}
`}
