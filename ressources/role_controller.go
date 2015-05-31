// Generated by: main
// TypeWriter: controller
// Directive: +gen on Role

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

type AbstractRoleInter interface {
	Create(roles []domain.Role) ([]domain.Role, error)
	CreateOne(role *domain.Role) (*domain.Role, error)
	Find(filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.Role, error)
	FindByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.Role, error)
	Upsert(roles []domain.Role, filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.Role, error)
	UpsertOne(role *domain.Role, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.Role, error)
	UpdateByID(id int, role *domain.Role, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.Role, error)
	DeleteAll(filter *usecases.Filter, ownerRelations []domain.Relation) error
	DeleteByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) error
}

type RoleCtrl struct {
	interactor AbstractRoleInter
	render     interfaces.AbstractRender
	routeDir   *interfaces.RouteDirectory
}

func NewRoleCtrl(interactor AbstractRoleInter, render interfaces.AbstractRender, routeDir *interfaces.RouteDirectory) *RoleCtrl {
	controller := &RoleCtrl{interactor: interactor, render: render, routeDir: routeDir}

	if routeDir != nil {
		setRoleAccess(routeDir, controller)
	}

	return controller
}

func (c *RoleCtrl) Create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	role := &domain.Role{}
	var roles []domain.Role

	buffer, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(buffer, role)
	if err != nil {
		err := json.Unmarshal(buffer, &roles)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
			return
		}
	}

	lastRessource := interfaces.GetLastRessource(r)

	if roles == nil {
		role.ScopeModel(lastRessource.ID)
		role, err = c.interactor.CreateOne(role)
	} else {
		for i := range roles {
			(&roles[i]).ScopeModel(lastRessource.ID)
		}
		roles, err = c.interactor.Create(roles)
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

	if roles == nil {
		c.render.JSON(w, http.StatusCreated, role)
	} else {
		c.render.JSON(w, http.StatusCreated, roles)
	}
}

func (c *RoleCtrl) Find(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	filter, err := interfaces.GetQueryFilter(r)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.FilterDecodingError, err)
		return
	}

	filter = interfaces.FilterIfLastRessource(r, filter)
	filter = interfaces.FilterIfOwnerRelations(r, filter)
	relations := interfaces.GetOwnerRelations(r)

	roles, err := c.interactor.Find(filter, relations)
	if err != nil {
		c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		return
	}

	c.render.JSON(w, http.StatusOK, roles)
}

func (c *RoleCtrl) FindByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	role, err := c.interactor.FindByID(id, filter, relations)
	if err != nil {
		c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
		return
	}

	c.render.JSON(w, http.StatusOK, role)
}

func (c *RoleCtrl) Upsert(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	role := &domain.Role{}
	var roles []domain.Role

	buffer, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(buffer, role)
	if err != nil {
		err := json.Unmarshal(buffer, &roles)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
			return
		}
	}

	lastRessource := interfaces.GetLastRessource(r)
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	ownerRelations := interfaces.GetOwnerRelations(r)

	if roles == nil {
		role.ScopeModel(lastRessource.ID)
		role, err = c.interactor.UpsertOne(role, filter, ownerRelations)
	} else {
		for i := range roles {
			(&roles[i]).ScopeModel(lastRessource.ID)
		}
		roles, err = c.interactor.Upsert(roles, filter, ownerRelations)
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

	if roles == nil {
		c.render.JSON(w, http.StatusCreated, role)
	} else {
		c.render.JSON(w, http.StatusCreated, roles)
	}
}

func (c *RoleCtrl) UpdateByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.InvalidPathParams, err)
		return
	}

	role := &domain.Role{}

	err = json.NewDecoder(r.Body).Decode(role)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
		return
	}

	lastRessource := interfaces.GetLastRessource(r)
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	ownerRelations := interfaces.GetOwnerRelations(r)

	role.ScopeModel(lastRessource.ID)
	role, err = c.interactor.UpdateByID(id, role, filter, ownerRelations)

	if err != nil {
		switch err.(type) {
		case *internalerrors.ViolatedConstraint:
			c.render.JSONError(w, 422, apierrors.ViolatedConstraint, err)
		default:
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		}
		return
	}

	c.render.JSON(w, http.StatusCreated, role)
}

func (c *RoleCtrl) DeleteAll(w http.ResponseWriter, r *http.Request, _ map[string]string) {
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

func (c *RoleCtrl) DeleteByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

func (c *RoleCtrl) Related(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	context.Set(r, "lastRessource", &interfaces.Ressource{Name: related, IDKey: "roleID", ID: pk})

	(*handler)(w, r, params)
}

func (c *RoleCtrl) RelatedOne(w http.ResponseWriter, r *http.Request, params map[string]string) {
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
