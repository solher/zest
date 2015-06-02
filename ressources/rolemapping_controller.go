// Generated by: main
// TypeWriter: controller
// Directive: +gen on RoleMapping

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

type AbstractRoleMappingInter interface {
	Create(rolemappings []domain.RoleMapping) ([]domain.RoleMapping, error)
	CreateOne(rolemapping *domain.RoleMapping) (*domain.RoleMapping, error)
	Find(filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.RoleMapping, error)
	FindByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.RoleMapping, error)
	Upsert(rolemappings []domain.RoleMapping, filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.RoleMapping, error)
	UpsertOne(rolemapping *domain.RoleMapping, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.RoleMapping, error)
	UpdateByID(id int, rolemapping *domain.RoleMapping, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.RoleMapping, error)
	DeleteAll(filter *usecases.Filter, ownerRelations []domain.Relation) error
	DeleteByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) error
}

type RoleMappingCtrl struct {
	interactor AbstractRoleMappingInter
	render     interfaces.AbstractRender
	routeDir   *interfaces.RouteDirectory
}

func NewRoleMappingCtrl(interactor AbstractRoleMappingInter, render interfaces.AbstractRender, routeDir *interfaces.RouteDirectory) *RoleMappingCtrl {
	controller := &RoleMappingCtrl{interactor: interactor, render: render, routeDir: routeDir}

	if routeDir != nil {
		setRoleMappingAccess(routeDir, controller)
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
		rolemapping.SetRelatedID(lastRessource.IDKey, lastRessource.ID)
		rolemapping, err = c.interactor.CreateOne(rolemapping)
	} else {
		for i := range rolemappings {
			(&rolemappings[i]).SetRelatedID(lastRessource.IDKey, lastRessource.ID)
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
		rolemapping.BeforeRender()
		c.render.JSON(w, http.StatusCreated, rolemapping)
	} else {
		for i := range rolemappings {
			(&rolemappings[i]).BeforeRender()
		}
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
	filter = interfaces.FilterIfOwnerRelations(r, filter)
	relations := interfaces.GetOwnerRelations(r)

	rolemappings, err := c.interactor.Find(filter, relations)
	if err != nil {
		c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		return
	}

	for i := range rolemappings {
		(&rolemappings[i]).BeforeRender()
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

	filter = interfaces.FilterIfOwnerRelations(r, filter)
	relations := interfaces.GetOwnerRelations(r)

	rolemapping, err := c.interactor.FindByID(id, filter, relations)
	if err != nil {
		c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
		return
	}

	rolemapping.BeforeRender()
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
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	ownerRelations := interfaces.GetOwnerRelations(r)

	if rolemappings == nil {
		rolemapping.SetRelatedID(lastRessource.IDKey, lastRessource.ID)
		rolemapping, err = c.interactor.UpsertOne(rolemapping, filter, ownerRelations)
	} else {
		for i := range rolemappings {
			(&rolemappings[i]).SetRelatedID(lastRessource.IDKey, lastRessource.ID)
		}
		rolemappings, err = c.interactor.Upsert(rolemappings, filter, ownerRelations)
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
		rolemapping.BeforeRender()
		c.render.JSON(w, http.StatusCreated, rolemapping)
	} else {
		for i := range rolemappings {
			(&rolemappings[i]).BeforeRender()
		}
		c.render.JSON(w, http.StatusCreated, rolemappings)
	}
}

func (c *RoleMappingCtrl) UpdateByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.InvalidPathParams, err)
		return
	}

	rolemapping := &domain.RoleMapping{}

	err = json.NewDecoder(r.Body).Decode(rolemapping)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
		return
	}

	lastRessource := interfaces.GetLastRessource(r)
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	ownerRelations := interfaces.GetOwnerRelations(r)

	rolemapping.SetRelatedID(lastRessource.IDKey, lastRessource.ID)
	rolemapping, err = c.interactor.UpdateByID(id, rolemapping, filter, ownerRelations)

	if err != nil {
		switch err.(type) {
		case *internalerrors.ViolatedConstraint:
			c.render.JSONError(w, 422, apierrors.ViolatedConstraint, err)
		default:
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		}
		return
	}

	rolemapping.BeforeRender()
	c.render.JSON(w, http.StatusCreated, rolemapping)
}

func (c *RoleMappingCtrl) DeleteAll(w http.ResponseWriter, r *http.Request, _ map[string]string) {
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

func (c *RoleMappingCtrl) DeleteByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
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
