// Generated by: main
// TypeWriter: controller
// Directive: +gen on RoleMapping

package ressources

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Solher/auth-scaffold/apierrors"
	"github.com/Solher/auth-scaffold/domain"
	"github.com/Solher/auth-scaffold/interfaces"
	"github.com/Solher/auth-scaffold/internalerrors"
	"github.com/dimfeld/httptreemux"
	"github.com/gorilla/context"
)

type AbstractRoleMappingInter interface {
	Create(rolemappings []domain.RoleMapping) ([]domain.RoleMapping, error)
	CreateOne(rolemapping *domain.RoleMapping) (*domain.RoleMapping, error)
	Find(filter *interfaces.Filter) ([]domain.RoleMapping, error)
	FindByID(id int, filter *interfaces.Filter) (*domain.RoleMapping, error)
	Upsert(rolemappings []domain.RoleMapping) ([]domain.RoleMapping, error)
	UpsertOne(rolemapping *domain.RoleMapping) (*domain.RoleMapping, error)
	DeleteAll(filter *interfaces.Filter) error
	DeleteByID(id int) error
}

type RoleMappingCtrl struct {
	interactor AbstractRoleMappingInter
	render     interfaces.AbstractRender
	routeDir   *interfaces.RouteDirectory
}

func NewRoleMappingCtrl(interactor AbstractRoleMappingInter, render interfaces.AbstractRender, routeDir *interfaces.RouteDirectory) *RoleMappingCtrl {
	controller := &RoleMappingCtrl{interactor: interactor, render: render, routeDir: routeDir}

	if routeDir != nil {
		setRoleMappingAccessOptions(routeDir, controller)
	}

	return controller
}

func (c *RoleMappingCtrl) Create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	rolemapping := &domain.RoleMapping{}
	var rolemappings []domain.RoleMapping

	buffer, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(buffer, rolemapping)
	if err != nil {
		err := json.Unmarshal(buffer, &rolemappings)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
			return
		}
	}

	lastRessource := interfaces.GetLastRessource(r)

	if rolemappings == nil {
		rolemapping.ScopeModel(lastRessource.ID)
		rolemapping, err = c.interactor.CreateOne(rolemapping)
	} else {
		for i := range rolemappings {
			(&rolemappings[i]).ScopeModel(lastRessource.ID)
		}
		rolemappings, err = c.interactor.Create(rolemappings)
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

	if rolemappings == nil {
		c.render.JSON(w, http.StatusCreated, rolemapping)
	} else {
		c.render.JSON(w, http.StatusCreated, rolemappings)
	}
}

func (c *RoleMappingCtrl) Find(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	filter, err := interfaces.GetQueryFilter(r)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.FilterDecodingError, err)
		return
	}

	filter = interfaces.FilterIfLastRessource(r, filter)

	rolemappings, err := c.interactor.Find(filter)
	if err != nil {
		c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		return
	}

	c.render.JSON(w, http.StatusOK, rolemappings)
}

func (c *RoleMappingCtrl) FindByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	rolemapping, err := c.interactor.FindByID(id, filter)
	if err != nil {
		c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
		return
	}

	c.render.JSON(w, http.StatusOK, rolemapping)
}

func (c *RoleMappingCtrl) Upsert(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	rolemapping := &domain.RoleMapping{}
	var rolemappings []domain.RoleMapping

	buffer, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(buffer, rolemapping)
	if err != nil {
		err := json.Unmarshal(buffer, &rolemappings)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
			return
		}
	}

	lastRessource := interfaces.GetLastRessource(r)

	if rolemappings == nil {
		rolemapping.ScopeModel(lastRessource.ID)
		rolemapping, err = c.interactor.UpsertOne(rolemapping)
	} else {
		for i := range rolemappings {
			(&rolemappings[i]).ScopeModel(lastRessource.ID)
		}
		rolemappings, err = c.interactor.Upsert(rolemappings)
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

	if rolemappings == nil {
		c.render.JSON(w, http.StatusCreated, rolemapping)
	} else {
		c.render.JSON(w, http.StatusCreated, rolemappings)
	}
}

func (c *RoleMappingCtrl) DeleteAll(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	filter, err := interfaces.GetQueryFilter(r)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.FilterDecodingError, err)
		return
	}

	filter = interfaces.FilterIfLastRessource(r, filter)

	err = c.interactor.DeleteAll(filter)
	if err != nil {
		c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		return
	}

	c.render.JSON(w, http.StatusNoContent, nil)
}

func (c *RoleMappingCtrl) DeleteByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
	id, err := strconv.Atoi(params["id"])
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

func (c *RoleMappingCtrl) Related(w http.ResponseWriter, r *http.Request, params map[string]string) {
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
		handler = c.routeDir.Get(key.For("Create")).Handler
	case "GET":
		handler = c.routeDir.Get(key.For("Find")).Handler
	case "PUT":
		handler = c.routeDir.Get(key.For("Upsert")).Handler
	case "DELETE":
		handler = c.routeDir.Get(key.For("DeleteAll")).Handler
	}

	if handler == nil {
		c.render.JSON(w, http.StatusNotFound, nil)
		return
	}

	context.Set(r, "lastRessource", &interfaces.Ressource{Name: related, IDKey: "rolemappingID", ID: pk})

	(*handler)(w, r, params)
}

func (c *RoleMappingCtrl) RelatedOne(w http.ResponseWriter, r *http.Request, params map[string]string) {
	params["id"] = params["fk"]

	related := params["related"]
	key := interfaces.NewDirectoryKey(related)

	var handler httptreemux.HandlerFunc

	switch r.Method {
	case "GET":
		handler = *c.routeDir.Get(key.For("FindByID")).Handler
	case "DELETE":
		handler = *c.routeDir.Get(key.For("DeleteByID")).Handler
	}

	if handler == nil {
		c.render.JSON(w, http.StatusNotFound, nil)
		return
	}

	handler(w, r, params)
}
