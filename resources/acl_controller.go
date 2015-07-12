// @SubApi Acl resource [/acls]
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
	usecases.DependencyDirectory.Register(NewAclCtrl)
}

type AbstractAclInter interface {
	Create(acls []domain.Acl) ([]domain.Acl, error)
	CreateOne(acl *domain.Acl) (*domain.Acl, error)
	Find(context usecases.QueryContext) ([]domain.Acl, error)
	FindByID(id int, context usecases.QueryContext) (*domain.Acl, error)
	Upsert(acls []domain.Acl, context usecases.QueryContext) ([]domain.Acl, error)
	UpsertOne(acl *domain.Acl, context usecases.QueryContext) (*domain.Acl, error)
	UpdateByID(id int, acl *domain.Acl, context usecases.QueryContext) (*domain.Acl, error)
	DeleteAll(context usecases.QueryContext) error
	DeleteByID(id int, context usecases.QueryContext) error
}

type AclCtrl struct {
	interactor AbstractAclInter
	render     interfaces.AbstractRender
	routeDir   *usecases.RouteDirectory
}

func NewAclCtrl(interactor AbstractAclInter, render interfaces.AbstractRender, routeDir *usecases.RouteDirectory) *AclCtrl {
	controller := &AclCtrl{interactor: interactor, render: render, routeDir: routeDir}

	if routeDir != nil {
		setAclRoutes(routeDir, controller)
	}

	return controller
}

// @Title Create
// @Description Create one or multiple Acl instances
// @Accept  json
// @Param   Acl body domain.Acl true "Acl instance(s) data"
// @Success 201 {object} domain.Acl "Request was successful"
// @Router /acls [post]
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

	lastResource := interfaces.GetLastResource(r)

	if acls == nil {
		acl.SetRelatedID(lastResource.IDKey, lastResource.ID)
		acl, err = c.interactor.CreateOne(acl)
	} else {
		for i := range acls {
			(&acls[i]).SetRelatedID(lastResource.IDKey, lastResource.ID)
		}
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
		acl.BeforeRender()
		c.render.JSON(w, http.StatusCreated, acl)
	} else {
		for i := range acls {
			(&acls[i]).BeforeRender()
		}
		c.render.JSON(w, http.StatusCreated, acls)
	}
}

// @Title Find
// @Description Find all Acl instances matched by filter
// @Accept  json
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 200 {object} domain.Acl "Request was successful"
// @Router /acls [get]
func (c *AclCtrl) Find(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	filter, err := interfaces.GetQueryFilter(r)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.FilterDecodingError, err)
		return
	}

	filter = interfaces.FilterIfLastResource(r, filter)
	filter = interfaces.FilterIfOwnerRelations(r, filter)
	relations := interfaces.GetOwnerRelations(r)

	acls, err := c.interactor.Find(usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	if err != nil {
		c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		return
	}

	for i := range acls {
		(&acls[i]).BeforeRender()
	}
	c.render.JSON(w, http.StatusOK, acls)
}

// @Title FindByID
// @Description Find a Acl instance
// @Accept  json
// @Param   id path int true "Acl id"
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 200 {object} domain.Acl "Request was successful"
// @Router /acls/{id} [get]
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

	acl, err := c.interactor.FindByID(id, usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	if err != nil {
		switch err {
		case internalerrors.NotFound:
			c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
		default:
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		}
		return
	}

	acl.BeforeRender()
	c.render.JSON(w, http.StatusOK, acl)
}

// @Title Upsert
// @Description Upsert one or multiple Acl instances
// @Accept  json
// @Param   Acl body domain.Acl true "Acl instance(s) data"
// @Success 201 {object} domain.Acl "Request was successful"
// @Router /acls [put]
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

	lastResource := interfaces.GetLastResource(r)
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	relations := interfaces.GetOwnerRelations(r)

	if acls == nil {
		acl.SetRelatedID(lastResource.IDKey, lastResource.ID)
		acl, err = c.interactor.UpsertOne(acl, usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	} else {
		for i := range acls {
			(&acls[i]).SetRelatedID(lastResource.IDKey, lastResource.ID)
		}
		acls, err = c.interactor.Upsert(acls, usecases.QueryContext{Filter: filter, OwnerRelations: relations})
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

	if acls == nil {
		acl.BeforeRender()
		c.render.JSON(w, http.StatusCreated, acl)
	} else {
		for i := range acls {
			(&acls[i]).BeforeRender()
		}
		c.render.JSON(w, http.StatusCreated, acls)
	}
}

// @Title UpdateByID
// @Description Update attributes of a Acl instance
// @Accept  json
// @Param   id path int true "Acl id"
// @Param   Acl body domain.Acl true "Acl instance data"
// @Success 201 {object} domain.Acl
// @Router /acls/{id} [put]
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

	lastResource := interfaces.GetLastResource(r)
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	relations := interfaces.GetOwnerRelations(r)

	acl.SetRelatedID(lastResource.IDKey, lastResource.ID)
	acl, err = c.interactor.UpdateByID(id, acl, usecases.QueryContext{Filter: filter, OwnerRelations: relations})

	if err != nil {
		switch err {
		case internalerrors.NotFound:
			c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
		default:
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		}
		return
	}

	acl.BeforeRender()
	c.render.JSON(w, http.StatusCreated, acl)
}

// @Title DeleteAll
// @Description Delete all Acl instances matched by filter
// @Accept  json
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 204 {object} error "Request was successful"
// @Router /acls [delete]
func (c *AclCtrl) DeleteAll(w http.ResponseWriter, r *http.Request, _ map[string]string) {
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
// @Description Delete a Acl instance
// @Accept  json
// @Param   id path int true "Acl id"
// @Success 204 {object} error "Request was successful"
// @Router /acls/{id} [delete]
func (c *AclCtrl) DeleteByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
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
// @Description Create one or multiple Acl instances of a related resource
// @Accept  json
// @Param   pk path int true "Acl id"
// @Param   relatedResource path string true "Related resource name"
// @Param   Acl body domain.Acl true "Acl instance(s) data"
// @Success 201 {object} domain.Acl "Request was successful"
// @Router /acls/{pk}/{relatedResource} [post]
func (c *AclCtrl) CreateRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.related(w, r, params)
}

// @Title FindRelated
// @Description Find all Acl instances  of a related resource matched by filter
// @Accept  json
// @Param   pk path int true "Acl id"
// @Param   relatedResource path string true "Related resource name"
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 200 {object} domain.Acl "Request was successful"
// @Router /acls/{pk}/{relatedResource} [get]
func (c *AclCtrl) FindRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.related(w, r, params)
}

// @Title UpsertRelated
// @Description Upsert one or multiple Acl instances of a related resource
// @Accept  json
// @Param   pk path int true "Acl id"
// @Param   relatedResource path string true "Related resource name"
// @Param   Acl body domain.Acl true "Acl instance(s) data"
// @Success 201 {object} domain.Acl "Request was successful"
// @Router /acls/{pk}/{relatedResource} [put]
func (c *AclCtrl) UpsertRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.related(w, r, params)
}

// @Title DeleteAllRelated
// @Description Delete all Acl instances of a related resource matched by filter
// @Accept  json
// @Param   pk path int true "Acl id"
// @Param   relatedResource path string true "Related resource name"
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 204 {object} error "Request was successful"
// @Router /acls/{pk}/{relatedResource} [delete]
func (c *AclCtrl) DeleteAllRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.related(w, r, params)
}

func (c *AclCtrl) related(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	context.Set(r, "lastResource", &interfaces.Resource{Name: related, IDKey: "aclID", ID: pk})

	handler(w, r, params)
}

// @Title FindByIDRelated
// @Description Find a Acl instance of a related resource
// @Accept  json
// @Param   pk path int true "Acl id"
// @Param   relatedResource path string true "Related resource name"
// @Param   fk path int true "Related resource id"
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 200 {object} domain.Acl "Request was successful"
// @Router /acls/{pk}/{relatedResource}/{fk} [get]
func (c *AclCtrl) FindByIDRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.relatedOne(w, r, params)
}

// @Title UpdateByIDRelated
// @Description Update attributes of a Acl instance of a related resource
// @Accept  json
// @Param   pk path int true "Acl id"
// @Param   relatedResource path string true "Related resource name"
// @Param   fk path int true "Related resource id"
// @Param   Acl body domain.Acl true "Acl instance data"
// @Success 201 {object} domain.Acl
// @Router /acls/{pk}/{relatedResource}/{fk} [put]
func (c *AclCtrl) UpdateByIDRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.relatedOne(w, r, params)
}

// @Title DeleteByIDRelated
// @Description Delete a Acl instance of a related resource
// @Accept  json
// @Param   pk path int true "Acl id"
// @Param   relatedResource path string true "Related resource name"
// @Param   fk path int true "Related resource id"
// @Success 204 {object} error "Request was successful"
// @Router /acls/{pk}/{relatedResource}/{fk} [delete]
func (c *AclCtrl) DeleteByIDRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.relatedOne(w, r, params)
}

func (c *AclCtrl) relatedOne(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	context.Set(r, "lastResource", &interfaces.Resource{Name: related, IDKey: "aclID", ID: pk})

	handler(w, r, params)
}
