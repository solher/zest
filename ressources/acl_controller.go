package ressources

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

	lastRessource := interfaces.GetLastRessource(r)

	if acls == nil {
		acl.SetRelatedID(lastRessource.IDKey, lastRessource.ID)
		acl, err = c.interactor.CreateOne(acl)
	} else {
		for i := range acls {
			(&acls[i]).SetRelatedID(lastRessource.IDKey, lastRessource.ID)
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

func (c *AclCtrl) Find(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	filter, err := interfaces.GetQueryFilter(r)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.FilterDecodingError, err)
		return
	}

	filter = interfaces.FilterIfLastRessource(r, filter)
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

	lastRessource := interfaces.GetLastRessource(r)
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	relations := interfaces.GetOwnerRelations(r)

	if acls == nil {
		acl.SetRelatedID(lastRessource.IDKey, lastRessource.ID)
		acl, err = c.interactor.UpsertOne(acl, usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	} else {
		for i := range acls {
			(&acls[i]).SetRelatedID(lastRessource.IDKey, lastRessource.ID)
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

	lastRessource := interfaces.GetLastRessource(r)
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	relations := interfaces.GetOwnerRelations(r)

	acl.SetRelatedID(lastRessource.IDKey, lastRessource.ID)
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

func (c *AclCtrl) DeleteAll(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	filter, err := interfaces.GetQueryFilter(r)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.FilterDecodingError, err)
		return
	}

	filter = interfaces.FilterIfLastRessource(r, filter)
	filter = interfaces.FilterIfOwnerRelations(r, filter)
	relations := interfaces.GetOwnerRelations(r)

	err = c.interactor.DeleteAll(usecases.QueryContext{Filter: filter, OwnerRelations: relations})
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

func (c *AclCtrl) Related(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	context.Set(r, "lastRessource", &interfaces.Ressource{Name: related, IDKey: "aclID", ID: pk})

	handler(w, r, params)
}

func (c *AclCtrl) RelatedOne(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	context.Set(r, "lastRessource", &interfaces.Ressource{Name: related, IDKey: "aclID", ID: pk})

	handler(w, r, params)
}
