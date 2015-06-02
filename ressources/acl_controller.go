// Generated by: main
// TypeWriter: controller
// Directive: +gen on Acl

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
	"github.com/Solher/auth-scaffold/usecases"
	"github.com/dimfeld/httptreemux"
	"github.com/gorilla/context"
)

type AbstractAclInter interface {
	Create(acls []domain.Acl) ([]domain.Acl, error)
	CreateOne(acl *domain.Acl) (*domain.Acl, error)
	Find(filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.Acl, error)
	FindByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.Acl, error)
	Upsert(acls []domain.Acl, filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.Acl, error)
	UpsertOne(acl *domain.Acl, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.Acl, error)
	UpdateByID(id int, acl *domain.Acl, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.Acl, error)
	DeleteAll(filter *usecases.Filter, ownerRelations []domain.Relation) error
	DeleteByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) error
}

type AclCtrl struct {
	interactor AbstractAclInter
	render     interfaces.AbstractRender
	routeDir   *interfaces.RouteDirectory
}

func NewAclCtrl(interactor AbstractAclInter, render interfaces.AbstractRender, routeDir *interfaces.RouteDirectory) *AclCtrl {
	controller := &AclCtrl{interactor: interactor, render: render, routeDir: routeDir}

	if routeDir != nil {
		setAclAccess(routeDir, controller)
	}

	return controller
}

func (c *AclCtrl) Create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	acl := &domain.Acl{}
	var acls []domain.Acl

	buffer, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(buffer, acl)
	if err != nil {
		err := json.Unmarshal(buffer, &acls)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
			return
		}
	}

	// lastRessource := interfaces.GetLastRessource(r)

	if acls == nil {
		// acl.ScopeModel(lastRessource.ID)
		acl, err = c.interactor.CreateOne(acl)
	} else {
		// for i := range acls {
		// 	(&acls[i]).ScopeModel(lastRessource.ID)
		// }
		acls, err = c.interactor.Create(acls)
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

	if acls == nil {
		c.render.JSON(w, http.StatusCreated, acl)
	} else {
		c.render.JSON(w, http.StatusCreated, acls)
	}
}

func (c *AclCtrl) Find(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	filter, err := interfaces.GetQueryFilter(r)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.FilterDecodingError, err)
		return
	}

	filter = interfaces.FilterIfLastRessource(r, filter)
	filter = interfaces.FilterIfOwnerRelations(r, filter)
	relations := interfaces.GetOwnerRelations(r)

	acls, err := c.interactor.Find(filter, relations)
	if err != nil {
		c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		return
	}

	c.render.JSON(w, http.StatusOK, acls)
}

func (c *AclCtrl) FindByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	acl, err := c.interactor.FindByID(id, filter, relations)
	if err != nil {
		c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
		return
	}

	c.render.JSON(w, http.StatusOK, acl)
}

func (c *AclCtrl) Upsert(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	acl := &domain.Acl{}
	var acls []domain.Acl

	buffer, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(buffer, acl)
	if err != nil {
		err := json.Unmarshal(buffer, &acls)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
			return
		}
	}

	// lastRessource := interfaces.GetLastRessource(r)
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	ownerRelations := interfaces.GetOwnerRelations(r)

	if acls == nil {
		// acl.ScopeModel(lastRessource.ID)
		acl, err = c.interactor.UpsertOne(acl, filter, ownerRelations)
	} else {
		// for i := range acls {
		// 	(&acls[i]).ScopeModel(lastRessource.ID)
		// }
		acls, err = c.interactor.Upsert(acls, filter, ownerRelations)
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

	if acls == nil {
		c.render.JSON(w, http.StatusCreated, acl)
	} else {
		c.render.JSON(w, http.StatusCreated, acls)
	}
}

func (c *AclCtrl) UpdateByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.InvalidPathParams, err)
		return
	}

	acl := &domain.Acl{}

	err = json.NewDecoder(r.Body).Decode(acl)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
		return
	}

	// lastRessource := interfaces.GetLastRessource(r)
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	ownerRelations := interfaces.GetOwnerRelations(r)

	// acl.ScopeModel(lastRessource.ID)
	acl, err = c.interactor.UpdateByID(id, acl, filter, ownerRelations)

	if err != nil {
		switch err.(type) {
		case *internalerrors.ViolatedConstraint:
			c.render.JSONError(w, 422, apierrors.ViolatedConstraint, err)
		default:
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		}
		return
	}

	c.render.JSON(w, http.StatusCreated, acl)
}

func (c *AclCtrl) DeleteAll(w http.ResponseWriter, r *http.Request, _ map[string]string) {
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

func (c *AclCtrl) DeleteByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

func (c *AclCtrl) Related(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	context.Set(r, "lastRessource", &interfaces.Ressource{Name: related, IDKey: "aclID", ID: pk})

	(*handler)(w, r, params)
}

func (c *AclCtrl) RelatedOne(w http.ResponseWriter, r *http.Request, params map[string]string) {
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
