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
	type {{.Type}}Interactor interface {
		Create({{.Name}}s []{{.Type}}) ([]{{.Type}}, error)
		Find(filter *interfaces.Filter) ([]{{.Type}}, error)
		FindByID(id int, filter *interfaces.Filter) (*{{.Type}}, error)
		Upsert({{.Name}}s []{{.Type}}) ([]{{.Type}}, error)
		DeleteAll(filter *interfaces.Filter) error
		DeleteByID(id int) error
	}

	type Controller struct {
		interactor {{.Type}}Interactor
		render     interfaces.Render
	}

	func NewController(interactor {{.Type}}Interactor, render interfaces.Render, routesDir interfaces.RouteDirectory) *Controller {
		controller := &Controller{interactor: interactor, render: render}

		if routesDir != nil {
			addRoutes(routesDir, controller)
		}

		return controller
	}
`}

var create = &typewriter.Template{
	Name: "Create",
	Text: `
	func (c *Controller) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var {{.Name}}s []{{.Type}}
		err := json.NewDecoder(r.Body).Decode(&{{.Name}}s)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
			return
		}

		{{.Name}}s, err = c.interactor.Create({{.Name}}s)
		if err != nil {
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
			return
		}

		c.render.JSON(w, http.StatusCreated, {{.Name}}s)
	}
`}

var find = &typewriter.Template{
	Name: "Find",
	Text: `
	func (c *Controller) Find(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	func (c *Controller) FindByID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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
	func (c *Controller) Upsert(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var {{.Name}}s []{{.Type}}
		err := json.NewDecoder(r.Body).Decode(&{{.Name}}s)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
			return
		}

		{{.Name}}s, err = c.interactor.Upsert({{.Name}}s)
		if err != nil {
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
			return
		}

		c.render.JSON(w, http.StatusCreated, {{.Name}}s)
	}
`}

var deleteAll = &typewriter.Template{
	Name: "DeleteAll",
	Text: `
	func (c *Controller) DeleteAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	func (c *Controller) DeleteByID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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
