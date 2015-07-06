package resources

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/context"
	"github.com/solher/zest/apierrors"
	"github.com/solher/zest/domain"
	"github.com/solher/zest/interfaces"
	"github.com/solher/zest/internalerrors"
	"github.com/solher/zest/usecases"
)

func init() {
	usecases.DependencyDirectory.Register(NewRoleCtrl)
}

type AbstractRoleInter interface {
	Create(roles []domain.Role) ([]domain.Role, error)
	CreateOne(role *domain.Role) (*domain.Role, error)
	Find(context usecases.QueryContext) ([]domain.Role, error)
	FindByID(id int, context usecases.QueryContext) (*domain.Role, error)
	Upsert(roles []domain.Role, context usecases.QueryContext) ([]domain.Role, error)
	UpsertOne(role *domain.Role, context usecases.QueryContext) (*domain.Role, error)
	UpdateByID(id int, role *domain.Role, context usecases.QueryContext) (*domain.Role, error)
	DeleteAll(context usecases.QueryContext) error
	DeleteByID(id int, context usecases.QueryContext) error
}

type RoleCtrl struct {
	interactor AbstractRoleInter
	render     interfaces.AbstractRender
	routeDir   *usecases.RouteDirectory
}

func NewRoleCtrl(interactor AbstractRoleInter, render interfaces.AbstractRender, routeDir *usecases.RouteDirectory) *RoleCtrl {
	controller := &RoleCtrl{interactor: interactor, render: render, routeDir: routeDir}

	if routeDir != nil {
		setRoleRoutes(routeDir, controller)
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

	lastResource := interfaces.GetLastResource(r)

	if roles == nil {
		role.SetRelatedID(lastResource.IDKey, lastResource.ID)
		role, err = c.interactor.CreateOne(role)
	} else {
		for i := range roles {
			(&roles[i]).SetRelatedID(lastResource.IDKey, lastResource.ID)
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
		role.BeforeRender()
		c.render.JSON(w, http.StatusCreated, role)
	} else {
		for i := range roles {
			(&roles[i]).BeforeRender()
		}
		c.render.JSON(w, http.StatusCreated, roles)
	}
}

func (c *RoleCtrl) Find(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	filter, err := interfaces.GetQueryFilter(r)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.FilterDecodingError, err)
		return
	}

	filter = interfaces.FilterIfLastResource(r, filter)
	filter = interfaces.FilterIfOwnerRelations(r, filter)
	relations := interfaces.GetOwnerRelations(r)

	roles, err := c.interactor.Find(usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	if err != nil {
		c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		return
	}

	for i := range roles {
		(&roles[i]).BeforeRender()
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

	role, err := c.interactor.FindByID(id, usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	if err != nil {
		switch err {
		case internalerrors.NotFound:
			c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
		default:
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		}
		return
	}

	role.BeforeRender()
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

	lastResource := interfaces.GetLastResource(r)
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	relations := interfaces.GetOwnerRelations(r)

	if roles == nil {
		role.SetRelatedID(lastResource.IDKey, lastResource.ID)
		role, err = c.interactor.UpsertOne(role, usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	} else {
		for i := range roles {
			(&roles[i]).SetRelatedID(lastResource.IDKey, lastResource.ID)
		}
		roles, err = c.interactor.Upsert(roles, usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	}

	if err != nil {
		switch err.(type) {
		case *internalerrors.ViolatedConstraint:
			c.render.JSONError(w, 422, apierrors.ViolatedConstraint, err)
		}

		switch err {
		case internalerrors.NotFound:
			c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
		default:
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		}

		return
	}

	if roles == nil {
		role.BeforeRender()
		c.render.JSON(w, http.StatusCreated, role)
	} else {
		for i := range roles {
			(&roles[i]).BeforeRender()
		}
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

	lastResource := interfaces.GetLastResource(r)
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	relations := interfaces.GetOwnerRelations(r)

	role.SetRelatedID(lastResource.IDKey, lastResource.ID)
	role, err = c.interactor.UpdateByID(id, role, usecases.QueryContext{Filter: filter, OwnerRelations: relations})

	if err != nil {
		switch err {
		case internalerrors.NotFound:
			c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
		default:
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		}
		return
	}

	role.BeforeRender()
	c.render.JSON(w, http.StatusCreated, role)
}

func (c *RoleCtrl) DeleteAll(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	filter, err := interfaces.GetQueryFilter(r)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.FilterDecodingError, err)
		return
	}

	filter = interfaces.FilterIfLastResource(r, filter)
	filter = interfaces.FilterIfOwnerRelations(r, filter)
	relations := interfaces.GetOwnerRelations(r)

	err = c.interactor.DeleteAll(usecases.QueryContext{Filter: filter, OwnerRelations: relations})
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
	relations := interfaces.GetOwnerRelations(r)

	err = c.interactor.DeleteByID(id, usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	if err != nil {
		switch err {
		case internalerrors.NotFound:
			c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
		default:
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		}
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

	context.Set(r, "lastResource", &interfaces.Resource{Name: related, IDKey: "roleID", ID: pk})

	handler(w, r, params)
}

func (c *RoleCtrl) RelatedOne(w http.ResponseWriter, r *http.Request, params map[string]string) {
	pk, err := strconv.Atoi(params["pk"])
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.InvalidPathParams, err)
		return
	}

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

	context.Set(r, "lastResource", &interfaces.Resource{Name: related, IDKey: "roleID", ID: pk})

	handler(w, r, params)
}
