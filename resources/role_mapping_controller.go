// @SubApi RoleMapping resource [/roleMappings]
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
	usecases.DependencyDirectory.Register(NewRoleMappingCtrl)
}

type AbstractRoleMappingInter interface {
	Create(roleMappings []domain.RoleMapping) ([]domain.RoleMapping, error)
	CreateOne(roleMapping *domain.RoleMapping) (*domain.RoleMapping, error)
	Find(context usecases.QueryContext) ([]domain.RoleMapping, error)
	FindByID(id int, context usecases.QueryContext) (*domain.RoleMapping, error)
	Upsert(roleMappings []domain.RoleMapping, context usecases.QueryContext) ([]domain.RoleMapping, error)
	UpsertOne(roleMapping *domain.RoleMapping, context usecases.QueryContext) (*domain.RoleMapping, error)
	UpdateByID(id int, roleMapping *domain.RoleMapping, context usecases.QueryContext) (*domain.RoleMapping, error)
	DeleteAll(context usecases.QueryContext) error
	DeleteByID(id int, context usecases.QueryContext) error
}

type RoleMappingCtrl struct {
	interactor AbstractRoleMappingInter
	render     interfaces.AbstractRender
	routeDir   *usecases.RouteDirectory
}

func NewRoleMappingCtrl(interactor AbstractRoleMappingInter, render interfaces.AbstractRender, routeDir *usecases.RouteDirectory) *RoleMappingCtrl {
	controller := &RoleMappingCtrl{interactor: interactor, render: render, routeDir: routeDir}

	if routeDir != nil {
		setRoleMappingRoutes(routeDir, controller)
	}

	return controller
}

// @Title Create
// @Description Create one or multiple RoleMapping instances
// @Accept  json
// @Param   RoleMapping body domain.RoleMapping true "RoleMapping instance(s) data"
// @Success 201 {object} domain.RoleMapping "Request was successful"
// @Router /roleMappings [post]
func (c *RoleMappingCtrl) Create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	roleMapping := &domain.RoleMapping{}
	var roleMappings []domain.RoleMapping

	buffer, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(buffer, roleMapping)
	if err != nil {
		err := json.Unmarshal(buffer, &roleMappings)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
			return
		}
	}

	lastResource := interfaces.GetLastResource(r)

	if roleMappings == nil {
		roleMapping.SetRelatedID(lastResource.IDKey, lastResource.ID)
		roleMapping, err = c.interactor.CreateOne(roleMapping)
	} else {
		for i := range roleMappings {
			(&roleMappings[i]).SetRelatedID(lastResource.IDKey, lastResource.ID)
		}
		roleMappings, err = c.interactor.Create(roleMappings)
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

	if roleMappings == nil {
		roleMapping.BeforeRender()
		c.render.JSON(w, http.StatusCreated, roleMapping)
	} else {
		for i := range roleMappings {
			(&roleMappings[i]).BeforeRender()
		}
		c.render.JSON(w, http.StatusCreated, roleMappings)
	}
}

// @Title Find
// @Description Find all RoleMapping instances matched by filter
// @Accept  json
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 200 {object} domain.RoleMapping "Request was successful"
// @Router /roleMappings [get]
func (c *RoleMappingCtrl) Find(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	filter, err := interfaces.GetQueryFilter(r)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.FilterDecodingError, err)
		return
	}

	filter = interfaces.FilterIfLastResource(r, filter)
	filter = interfaces.FilterIfOwnerRelations(r, filter)
	relations := interfaces.GetOwnerRelations(r)

	roleMappings, err := c.interactor.Find(usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	if err != nil {
		c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		return
	}

	for i := range roleMappings {
		(&roleMappings[i]).BeforeRender()
	}
	c.render.JSON(w, http.StatusOK, roleMappings)
}

// @Title FindByID
// @Description Find a RoleMapping instance
// @Accept  json
// @Param   id path int true "RoleMapping id"
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 200 {object} domain.RoleMapping "Request was successful"
// @Router /roleMappings/{id} [get]
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

	roleMapping, err := c.interactor.FindByID(id, usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	if err != nil {
		switch err {
		case internalerrors.NotFound:
			c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
		default:
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		}
		return
	}

	roleMapping.BeforeRender()
	c.render.JSON(w, http.StatusOK, roleMapping)
}

// @Title Upsert
// @Description Upsert one or multiple RoleMapping instances
// @Accept  json
// @Param   RoleMapping body domain.RoleMapping true "RoleMapping instance(s) data"
// @Success 201 {object} domain.RoleMapping "Request was successful"
// @Router /roleMappings [put]
func (c *RoleMappingCtrl) Upsert(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	roleMapping := &domain.RoleMapping{}
	var roleMappings []domain.RoleMapping

	buffer, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(buffer, roleMapping)
	if err != nil {
		err := json.Unmarshal(buffer, &roleMappings)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
			return
		}
	}

	lastResource := interfaces.GetLastResource(r)
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	relations := interfaces.GetOwnerRelations(r)

	if roleMappings == nil {
		roleMapping.SetRelatedID(lastResource.IDKey, lastResource.ID)
		roleMapping, err = c.interactor.UpsertOne(roleMapping, usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	} else {
		for i := range roleMappings {
			(&roleMappings[i]).SetRelatedID(lastResource.IDKey, lastResource.ID)
		}
		roleMappings, err = c.interactor.Upsert(roleMappings, usecases.QueryContext{Filter: filter, OwnerRelations: relations})
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

	if roleMappings == nil {
		roleMapping.BeforeRender()
		c.render.JSON(w, http.StatusCreated, roleMapping)
	} else {
		for i := range roleMappings {
			(&roleMappings[i]).BeforeRender()
		}
		c.render.JSON(w, http.StatusCreated, roleMappings)
	}
}

// @Title UpdateByID
// @Description Update attributes of a RoleMapping instance
// @Accept  json
// @Param   id path int true "RoleMapping id"
// @Param   RoleMapping body domain.RoleMapping true "RoleMapping instance data"
// @Success 201 {object} domain.RoleMapping
// @Router /roleMappings/{id} [put]
func (c *RoleMappingCtrl) UpdateByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.InvalidPathParams, err)
		return
	}

	roleMapping := &domain.RoleMapping{}

	err = json.NewDecoder(r.Body).Decode(roleMapping)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
		return
	}

	lastResource := interfaces.GetLastResource(r)
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	relations := interfaces.GetOwnerRelations(r)

	roleMapping.SetRelatedID(lastResource.IDKey, lastResource.ID)
	roleMapping, err = c.interactor.UpdateByID(id, roleMapping, usecases.QueryContext{Filter: filter, OwnerRelations: relations})

	if err != nil {
		switch err {
		case internalerrors.NotFound:
			c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
		default:
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		}
		return
	}

	roleMapping.BeforeRender()
	c.render.JSON(w, http.StatusCreated, roleMapping)
}

// @Title DeleteAll
// @Description Delete all RoleMapping instances matched by filter
// @Accept  json
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 204 {object} error "Request was successful"
// @Router /roleMappings [delete]
func (c *RoleMappingCtrl) DeleteAll(w http.ResponseWriter, r *http.Request, _ map[string]string) {
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

// @Title DeleteByID
// @Description Delete a RoleMapping instance
// @Accept  json
// @Param   id path int true "RoleMapping id"
// @Success 204 {object} error "Request was successful"
// @Router /roleMappings/{id} [delete]
func (c *RoleMappingCtrl) DeleteByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

// @Title CreateRelated
// @Description Create one or multiple RoleMapping instances of a related resource
// @Accept  json
// @Param   pk path int true "RoleMapping id"
// @Param   relatedResource path string true "Related resource name"
// @Param   RoleMapping body domain.RoleMapping true "RoleMapping instance(s) data"
// @Success 201 {object} domain.RoleMapping "Request was successful"
// @Router /roleMappings/{pk}/{relatedResource} [post]
func (c *RoleMappingCtrl) CreateRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.related(w, r, params)
}

// @Title FindRelated
// @Description Find all RoleMapping instances  of a related resource matched by filter
// @Accept  json
// @Param   pk path int true "RoleMapping id"
// @Param   relatedResource path string true "Related resource name"
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 200 {object} domain.RoleMapping "Request was successful"
// @Router /roleMappings/{pk}/{relatedResource} [get]
func (c *RoleMappingCtrl) FindRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.related(w, r, params)
}

// @Title UpsertRelated
// @Description Upsert one or multiple RoleMapping instances of a related resource
// @Accept  json
// @Param   pk path int true "RoleMapping id"
// @Param   relatedResource path string true "Related resource name"
// @Param   RoleMapping body domain.RoleMapping true "RoleMapping instance(s) data"
// @Success 201 {object} domain.RoleMapping "Request was successful"
// @Router /roleMappings/{pk}/{relatedResource} [put]
func (c *RoleMappingCtrl) UpsertRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.related(w, r, params)
}

// @Title DeleteAllRelated
// @Description Delete all RoleMapping instances of a related resource matched by filter
// @Accept  json
// @Param   pk path int true "RoleMapping id"
// @Param   relatedResource path string true "Related resource name"
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 204 {object} error "Request was successful"
// @Router /roleMappings/{pk}/{relatedResource} [delete]
func (c *RoleMappingCtrl) DeleteAllRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.related(w, r, params)
}

func (c *RoleMappingCtrl) related(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	context.Set(r, "lastResource", &interfaces.Resource{Name: related, IDKey: "roleMappingID", ID: pk})

	handler(w, r, params)
}

// @Title FindByIDRelated
// @Description Find a RoleMapping instance of a related resource
// @Accept  json
// @Param   pk path int true "RoleMapping id"
// @Param   relatedResource path string true "Related resource name"
// @Param   fk path int true "Related resource id"
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 200 {object} domain.RoleMapping "Request was successful"
// @Router /roleMappings/{pk}/{relatedResource}/{fk} [get]
func (c *RoleMappingCtrl) FindByIDRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.relatedOne(w, r, params)
}

// @Title UpdateByIDRelated
// @Description Update attributes of a RoleMapping instance of a related resource
// @Accept  json
// @Param   pk path int true "RoleMapping id"
// @Param   relatedResource path string true "Related resource name"
// @Param   fk path int true "Related resource id"
// @Param   RoleMapping body domain.RoleMapping true "RoleMapping instance data"
// @Success 201 {object} domain.RoleMapping
// @Router /roleMappings/{pk}/{relatedResource}/{fk} [put]
func (c *RoleMappingCtrl) UpdateByIDRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.relatedOne(w, r, params)
}

// @Title DeleteByIDRelated
// @Description Delete a RoleMapping instance of a related resource
// @Accept  json
// @Param   pk path int true "RoleMapping id"
// @Param   relatedResource path string true "Related resource name"
// @Param   fk path int true "Related resource id"
// @Success 204 {object} error "Request was successful"
// @Router /roleMappings/{pk}/{relatedResource}/{fk} [delete]
func (c *RoleMappingCtrl) DeleteByIDRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.relatedOne(w, r, params)
}

func (c *RoleMappingCtrl) relatedOne(w http.ResponseWriter, r *http.Request, params map[string]string) {
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
	case "PUT":
		handler = c.routeDir.Get(key.For("UpdateByID")).EffectiveHandler
	case "DELETE":
		handler = c.routeDir.Get(key.For("DeleteByID")).EffectiveHandler
	}

	if handler == nil {
		c.render.JSON(w, http.StatusNotFound, nil)
		return
	}

	context.Set(r, "lastResource", &interfaces.Resource{Name: related, IDKey: "roleMappingID", ID: pk})

	handler(w, r, params)
}
