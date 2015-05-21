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
}

var controller = &typewriter.Template{
	Name: "Controller",
	Text: `
	type Abstract{{.Type}}Inter interface {
		Create({{.Name}}s []domain.{{.Type}}) ([]domain.{{.Type}}, error)
		CreateOne({{.Name}} *domain.{{.Type}}) (*domain.{{.Type}}, error)
		Find(filter *interfaces.Filter) ([]domain.{{.Type}}, error)
		FindByID(id int, filter *interfaces.Filter) (*domain.{{.Type}}, error)
		Upsert({{.Name}}s []domain.{{.Type}}) ([]domain.{{.Type}}, error)
		UpsertOne({{.Name}} *domain.{{.Type}}) (*domain.{{.Type}}, error)
		DeleteAll(filter *interfaces.Filter) error
		DeleteByID(id int) error
	}

	type {{.Type}}Ctrl struct {
		interactor Abstract{{.Type}}Inter
		render     interfaces.AbstractRender
	}

	func New{{.Type}}Ctrl(interactor Abstract{{.Type}}Inter, render interfaces.AbstractRender,
		routeDir interfaces.RouteDirectory, permissionDir usecases.PermissionDirectory) *{{.Type}}Ctrl {
		controller := &{{.Type}}Ctrl{interactor: interactor, render: render}

		if routeDir != nil && permissionDir != nil {
			set{{.Type}}AccessOptions(routeDir, permissionDir, controller)
		}

		return controller
	}
`}

var create = &typewriter.Template{
	Name: "Create",
	Text: `
	func (c *{{.Type}}Ctrl) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

		if {{.Name}}s == nil {
			{{.Name}}.ScopeModel()
			{{.Name}}, err = c.interactor.CreateOne({{.Name}})
		} else {
			for i := range {{.Name}}s {
				(&{{.Name}}s[i]).ScopeModel()
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
	func (c *{{.Type}}Ctrl) Find(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		filter, err := interfaces.GetQueryFilter(r)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.FilterDecodingError, err)
			return
		}

		{{.Name}}s, err := c.interactor.Find(filter)
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
	func (c *{{.Type}}Ctrl) FindByID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id, err := strconv.Atoi(params.ByName("id"))
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.InvalidPathParams, err)
			return
		}

		filter, err := interfaces.GetQueryFilter(r)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.FilterDecodingError, err)
			return
		}

		{{.Name}}, err := c.interactor.FindByID(id, filter)
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
	func (c *{{.Type}}Ctrl) Upsert(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

		if {{.Name}}s == nil {
			{{.Name}}.ScopeModel()
			{{.Name}}, err = c.interactor.UpsertOne({{.Name}})
		} else {
			for i := range {{.Name}}s {
				(&{{.Name}}s[i]).ScopeModel()
			}
			{{.Name}}s, err = c.interactor.Upsert({{.Name}}s)
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
	func (c *{{.Type}}Ctrl) DeleteAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		filter, err := interfaces.GetQueryFilter(r)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.FilterDecodingError, err)
			return
		}

		err = c.interactor.DeleteAll(filter)
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
	func (c *{{.Type}}Ctrl) DeleteByID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id, err := strconv.Atoi(params.ByName("id"))
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.InvalidPathParams, err)
			return
		}

		err = c.interactor.DeleteByID(id)
		if err != nil {
			c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
			return
		}

		c.render.JSON(w, http.StatusNoContent, nil)
	}
`}
