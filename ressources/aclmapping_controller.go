// Generated by: main
// TypeWriter: controller
// Directive: +gen on AclMapping

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

type AbstractAclMappingInter interface {
	Create(aclmappings []domain.AclMapping) ([]domain.AclMapping, error)
	CreateOne(aclmapping *domain.AclMapping) (*domain.AclMapping, error)
	Find(filter *interfaces.Filter) ([]domain.AclMapping, error)
	FindByID(id int, filter *interfaces.Filter) (*domain.AclMapping, error)
	Upsert(aclmappings []domain.AclMapping) ([]domain.AclMapping, error)
	UpsertOne(aclmapping *domain.AclMapping) (*domain.AclMapping, error)
	DeleteAll(filter *interfaces.Filter) error
	DeleteByID(id int) error
}

type AclMappingCtrl struct {
	interactor AbstractAclMappingInter
	render     interfaces.AbstractRender
	routeDir   *interfaces.RouteDirectory
}

func NewAclMappingCtrl(interactor AbstractAclMappingInter, render interfaces.AbstractRender, routeDir *interfaces.RouteDirectory) *AclMappingCtrl {
	controller := &AclMappingCtrl{interactor: interactor, render: render, routeDir: routeDir}

	if routeDir != nil {
		setAclMappingAccess(routeDir, controller)
	}

	return controller
}

func (c *AclMappingCtrl) Create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	aclmapping := &domain.AclMapping{}
	var aclmappings []domain.AclMapping

	buffer, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(buffer, aclmapping)
	if err != nil {
		err := json.Unmarshal(buffer, &aclmappings)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
			return
		}
	}

	lastRessource := interfaces.GetLastRessource(r)

	if aclmappings == nil {
		aclmapping.ScopeModel(lastRessource.ID)
		aclmapping, err = c.interactor.CreateOne(aclmapping)
	} else {
		for i := range aclmappings {
			(&aclmappings[i]).ScopeModel(lastRessource.ID)
		}
		aclmappings, err = c.interactor.Create(aclmappings)
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

	if aclmappings == nil {
		c.render.JSON(w, http.StatusCreated, aclmapping)
	} else {
		c.render.JSON(w, http.StatusCreated, aclmappings)
	}
}

func (c *AclMappingCtrl) Find(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	filter, err := interfaces.GetQueryFilter(r)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.FilterDecodingError, err)
		return
	}

	filter = interfaces.FilterIfLastRessource(r, filter)

	aclmappings, err := c.interactor.Find(filter)
	if err != nil {
		c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		return
	}

	c.render.JSON(w, http.StatusOK, aclmappings)
}

func (c *AclMappingCtrl) FindByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	aclmapping, err := c.interactor.FindByID(id, filter)
	if err != nil {
		c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
		return
	}

	c.render.JSON(w, http.StatusOK, aclmapping)
}

func (c *AclMappingCtrl) Upsert(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	aclmapping := &domain.AclMapping{}
	var aclmappings []domain.AclMapping

	buffer, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(buffer, aclmapping)
	if err != nil {
		err := json.Unmarshal(buffer, &aclmappings)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
			return
		}
	}

	lastRessource := interfaces.GetLastRessource(r)

	if aclmappings == nil {
		aclmapping.ScopeModel(lastRessource.ID)
		aclmapping, err = c.interactor.UpsertOne(aclmapping)
	} else {
		for i := range aclmappings {
			(&aclmappings[i]).ScopeModel(lastRessource.ID)
		}
		aclmappings, err = c.interactor.Upsert(aclmappings)
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

	if aclmappings == nil {
		c.render.JSON(w, http.StatusCreated, aclmapping)
	} else {
		c.render.JSON(w, http.StatusCreated, aclmappings)
	}
}

func (c *AclMappingCtrl) DeleteAll(w http.ResponseWriter, r *http.Request, _ map[string]string) {
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

func (c *AclMappingCtrl) DeleteByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

func (c *AclMappingCtrl) Related(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	context.Set(r, "lastRessource", &interfaces.Ressource{Name: related, IDKey: "aclmappingID", ID: pk})

	(*handler)(w, r, params)
}

func (c *AclMappingCtrl) RelatedOne(w http.ResponseWriter, r *http.Request, params map[string]string) {
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
