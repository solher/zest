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
	usecases.DependencyDirectory.Register(NewAclMappingCtrl)
}

type AbstractAclMappingInter interface {
	Create(aclMappings []domain.AclMapping) ([]domain.AclMapping, error)
	CreateOne(aclMapping *domain.AclMapping) (*domain.AclMapping, error)
	Find(context usecases.QueryContext) ([]domain.AclMapping, error)
	FindByID(id int, context usecases.QueryContext) (*domain.AclMapping, error)
	Upsert(aclMappings []domain.AclMapping, context usecases.QueryContext) ([]domain.AclMapping, error)
	UpsertOne(aclMapping *domain.AclMapping, context usecases.QueryContext) (*domain.AclMapping, error)
	UpdateByID(id int, aclMapping *domain.AclMapping, context usecases.QueryContext) (*domain.AclMapping, error)
	DeleteAll(context usecases.QueryContext) error
	DeleteByID(id int, context usecases.QueryContext) error
}

type AclMappingCtrl struct {
	interactor AbstractAclMappingInter
	render     interfaces.AbstractRender
	routeDir   *usecases.RouteDirectory
}

func NewAclMappingCtrl(interactor AbstractAclMappingInter, render interfaces.AbstractRender, routeDir *usecases.RouteDirectory) *AclMappingCtrl {
	controller := &AclMappingCtrl{interactor: interactor, render: render, routeDir: routeDir}

	if routeDir != nil {
		setAclMappingRoutes(routeDir, controller)
	}

	return controller
}

func (c *AclMappingCtrl) Create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	aclMapping := &domain.AclMapping{}
	var aclMappings []domain.AclMapping

	buffer, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(buffer, aclMapping)
	if err != nil {
		err := json.Unmarshal(buffer, &aclMappings)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
			return
		}
	}

	lastRessource := interfaces.GetLastRessource(r)

	if aclMappings == nil {
		aclMapping.SetRelatedID(lastRessource.IDKey, lastRessource.ID)
		aclMapping, err = c.interactor.CreateOne(aclMapping)
	} else {
		for i := range aclMappings {
			(&aclMappings[i]).SetRelatedID(lastRessource.IDKey, lastRessource.ID)
		}
		aclMappings, err = c.interactor.Create(aclMappings)
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

	if aclMappings == nil {
		aclMapping.BeforeRender()
		c.render.JSON(w, http.StatusCreated, aclMapping)
	} else {
		for i := range aclMappings {
			(&aclMappings[i]).BeforeRender()
		}
		c.render.JSON(w, http.StatusCreated, aclMappings)
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

	aclMappings, err := c.interactor.Find(usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	if err != nil {
		c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		return
	}

	for i := range aclMappings {
		(&aclMappings[i]).BeforeRender()
	}
	c.render.JSON(w, http.StatusOK, aclMappings)
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

	aclMapping, err := c.interactor.FindByID(id, usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	if err != nil {
		switch err {
		case internalerrors.NotFound:
			c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
		default:
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		}
		return
	}

	aclMapping.BeforeRender()
	c.render.JSON(w, http.StatusOK, aclMapping)
}

func (c *AclMappingCtrl) Upsert(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	aclMapping := &domain.AclMapping{}
	var aclMappings []domain.AclMapping

	buffer, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(buffer, aclMapping)
	if err != nil {
		err := json.Unmarshal(buffer, &aclMappings)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
			return
		}
	}

	lastRessource := interfaces.GetLastRessource(r)
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	relations := interfaces.GetOwnerRelations(r)

	if aclMappings == nil {
		aclMapping.SetRelatedID(lastRessource.IDKey, lastRessource.ID)
		aclMapping, err = c.interactor.UpsertOne(aclMapping, usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	} else {
		for i := range aclMappings {
			(&aclMappings[i]).SetRelatedID(lastRessource.IDKey, lastRessource.ID)
		}
		aclMappings, err = c.interactor.Upsert(aclMappings, usecases.QueryContext{Filter: filter, OwnerRelations: relations})
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

	if aclMappings == nil {
		aclMapping.BeforeRender()
		c.render.JSON(w, http.StatusCreated, aclMapping)
	} else {
		for i := range aclMappings {
			(&aclMappings[i]).BeforeRender()
		}
		c.render.JSON(w, http.StatusCreated, aclMappings)
	}
}

func (c *AclMappingCtrl) UpdateByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.InvalidPathParams, err)
		return
	}

	aclMapping := &domain.AclMapping{}

	err = json.NewDecoder(r.Body).Decode(aclMapping)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
		return
	}

	lastRessource := interfaces.GetLastRessource(r)
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	relations := interfaces.GetOwnerRelations(r)

	aclMapping.SetRelatedID(lastRessource.IDKey, lastRessource.ID)
	aclMapping, err = c.interactor.UpdateByID(id, aclMapping, usecases.QueryContext{Filter: filter, OwnerRelations: relations})

	if err != nil {
		switch err {
		case internalerrors.NotFound:
			c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
		default:
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		}
		return
	}

	aclMapping.BeforeRender()
	c.render.JSON(w, http.StatusCreated, aclMapping)
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

	err = c.interactor.DeleteAll(usecases.QueryContext{Filter: filter, OwnerRelations: relations})
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

func (c *AclMappingCtrl) Related(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	context.Set(r, "lastRessource", &interfaces.Ressource{Name: related, IDKey: "aclMappingID", ID: pk})

	handler(w, r, params)
}

func (c *AclMappingCtrl) RelatedOne(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	handler(w, r, params)
}
