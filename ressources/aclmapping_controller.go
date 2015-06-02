// Generated by: main
// TypeWriter: controller
// Directive: +gen on AclMapping

package ressources

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Solher/zest/apierrors"
	"github.com/Solher/zest/domain"
	"github.com/Solher/zest/interfaces"
	"github.com/Solher/zest/internalerrors"
	"github.com/Solher/zest/usecases"
	"github.com/dimfeld/httptreemux"
	"github.com/gorilla/context"
)

type AbstractAclMappingInter interface {
	Create(aclmappings []domain.AclMapping) ([]domain.AclMapping, error)
	CreateOne(aclmapping *domain.AclMapping) (*domain.AclMapping, error)
	Find(filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.AclMapping, error)
	FindByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.AclMapping, error)
	Upsert(aclmappings []domain.AclMapping, filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.AclMapping, error)
	UpsertOne(aclmapping *domain.AclMapping, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.AclMapping, error)
	UpdateByID(id int, aclmapping *domain.AclMapping, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.AclMapping, error)
	DeleteAll(filter *usecases.Filter, ownerRelations []domain.Relation) error
	DeleteByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) error
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

	// lastRessource := interfaces.GetLastRessource(r)

	if aclmappings == nil {
		// aclmapping.ScopeModel(lastRessource.ID)
		aclmapping, err = c.interactor.CreateOne(aclmapping)
	} else {
		// for i := range aclmappings {
		// 	(&aclmappings[i]).ScopeModel(lastRessource.ID)
		// }
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
	filter = interfaces.FilterIfOwnerRelations(r, filter)
	relations := interfaces.GetOwnerRelations(r)

	aclmappings, err := c.interactor.Find(filter, relations)
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

	filter = interfaces.FilterIfOwnerRelations(r, filter)
	relations := interfaces.GetOwnerRelations(r)

	aclmapping, err := c.interactor.FindByID(id, filter, relations)
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

	// lastRessource := interfaces.GetLastRessource(r)
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	ownerRelations := interfaces.GetOwnerRelations(r)

	if aclmappings == nil {
		// aclmapping.ScopeModel(lastRessource.ID)
		aclmapping, err = c.interactor.UpsertOne(aclmapping, filter, ownerRelations)
	} else {
		// for i := range aclmappings {
		// 	(&aclmappings[i]).ScopeModel(lastRessource.ID)
		// }
		aclmappings, err = c.interactor.Upsert(aclmappings, filter, ownerRelations)
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

func (c *AclMappingCtrl) UpdateByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.InvalidPathParams, err)
		return
	}

	aclmapping := &domain.AclMapping{}

	err = json.NewDecoder(r.Body).Decode(aclmapping)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
		return
	}

	// lastRessource := interfaces.GetLastRessource(r)
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	ownerRelations := interfaces.GetOwnerRelations(r)

	// aclmapping.ScopeModel(lastRessource.ID)
	aclmapping, err = c.interactor.UpdateByID(id, aclmapping, filter, ownerRelations)

	if err != nil {
		switch err.(type) {
		case *internalerrors.ViolatedConstraint:
			c.render.JSONError(w, 422, apierrors.ViolatedConstraint, err)
		default:
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		}
		return
	}

	c.render.JSON(w, http.StatusCreated, aclmapping)
}

func (c *AclMappingCtrl) DeleteAll(w http.ResponseWriter, r *http.Request, _ map[string]string) {
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

func (c *AclMappingCtrl) DeleteByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
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