// @SubApi Session resource [/sessions]
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
	usecases.DependencyDirectory.Register(NewSessionCtrl)
}

type AbstractSessionInter interface {
	Create(sessions []domain.Session) ([]domain.Session, error)
	CreateOne(session *domain.Session) (*domain.Session, error)
	Find(context usecases.QueryContext) ([]domain.Session, error)
	FindByID(id int, context usecases.QueryContext) (*domain.Session, error)
	Upsert(sessions []domain.Session, context usecases.QueryContext) ([]domain.Session, error)
	UpsertOne(session *domain.Session, context usecases.QueryContext) (*domain.Session, error)
	UpdateByID(id int, session *domain.Session, context usecases.QueryContext) (*domain.Session, error)
	DeleteAll(context usecases.QueryContext) error
	DeleteByID(id int, context usecases.QueryContext) error
}

type SessionCtrl struct {
	interactor AbstractSessionInter
	render     interfaces.AbstractRender
	routeDir   *usecases.RouteDirectory
}

func NewSessionCtrl(interactor AbstractSessionInter, render interfaces.AbstractRender, routeDir *usecases.RouteDirectory) *SessionCtrl {
	controller := &SessionCtrl{interactor: interactor, render: render, routeDir: routeDir}

	if routeDir != nil {
		setSessionRoutes(routeDir, controller)
	}

	return controller
}

// @Title Create
// @Description Create one or multiple Session instances
// @Accept  json
// @Param   Session body domain.Session true "Session instance(s) data"
// @Success 201 {object} domain.Session "Request was successful"
// @Router /sessions [post]
func (c *SessionCtrl) Create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	session := &domain.Session{}
	var sessions []domain.Session

	buffer, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(buffer, session)
	if err != nil {
		err := json.Unmarshal(buffer, &sessions)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
			return
		}
	}

	lastResource := interfaces.GetLastResource(r)

	if sessions == nil {
		session.SetRelatedID(lastResource.IDKey, lastResource.ID)
		session, err = c.interactor.CreateOne(session)
	} else {
		for i := range sessions {
			(&sessions[i]).SetRelatedID(lastResource.IDKey, lastResource.ID)
		}
		sessions, err = c.interactor.Create(sessions)
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

	if sessions == nil {
		session.BeforeRender()
		c.render.JSON(w, http.StatusCreated, session)
	} else {
		for i := range sessions {
			(&sessions[i]).BeforeRender()
		}
		c.render.JSON(w, http.StatusCreated, sessions)
	}
}

// @Title Find
// @Description Find all Session instances matched by filter
// @Accept  json
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 200 {object} domain.Session "Request was successful"
// @Router /sessions [get]
func (c *SessionCtrl) Find(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	filter, err := interfaces.GetQueryFilter(r)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.FilterDecodingError, err)
		return
	}

	filter = interfaces.FilterIfLastResource(r, filter)
	filter = interfaces.FilterIfOwnerRelations(r, filter)
	relations := interfaces.GetOwnerRelations(r)

	sessions, err := c.interactor.Find(usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	if err != nil {
		c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		return
	}

	for i := range sessions {
		(&sessions[i]).BeforeRender()
	}
	c.render.JSON(w, http.StatusOK, sessions)
}

// @Title FindByID
// @Description Find a Session instance
// @Accept  json
// @Param   id path int true "Session id"
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 200 {object} domain.Session "Request was successful"
// @Router /sessions/{id} [get]
func (c *SessionCtrl) FindByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	session, err := c.interactor.FindByID(id, usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	if err != nil {
		switch err {
		case internalerrors.NotFound:
			c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
		default:
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		}
		return
	}

	session.BeforeRender()
	c.render.JSON(w, http.StatusOK, session)
}

// @Title Upsert
// @Description Upsert one or multiple Session instances
// @Accept  json
// @Param   Session body domain.Session true "Session instance(s) data"
// @Success 201 {object} domain.Session "Request was successful"
// @Router /sessions [put]
func (c *SessionCtrl) Upsert(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	session := &domain.Session{}
	var sessions []domain.Session

	buffer, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(buffer, session)
	if err != nil {
		err := json.Unmarshal(buffer, &sessions)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
			return
		}
	}

	lastResource := interfaces.GetLastResource(r)
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	relations := interfaces.GetOwnerRelations(r)

	if sessions == nil {
		session.SetRelatedID(lastResource.IDKey, lastResource.ID)
		session, err = c.interactor.UpsertOne(session, usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	} else {
		for i := range sessions {
			(&sessions[i]).SetRelatedID(lastResource.IDKey, lastResource.ID)
		}
		sessions, err = c.interactor.Upsert(sessions, usecases.QueryContext{Filter: filter, OwnerRelations: relations})
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

	if sessions == nil {
		session.BeforeRender()
		c.render.JSON(w, http.StatusCreated, session)
	} else {
		for i := range sessions {
			(&sessions[i]).BeforeRender()
		}
		c.render.JSON(w, http.StatusCreated, sessions)
	}
}

// @Title UpdateByID
// @Description Update attributes of a Session instance
// @Accept  json
// @Param   id path int true "Session id"
// @Param   Session body domain.Session true "Session instance data"
// @Success 201 {object} domain.Session
// @Router /sessions/{id} [put]
func (c *SessionCtrl) UpdateByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.InvalidPathParams, err)
		return
	}

	session := &domain.Session{}

	err = json.NewDecoder(r.Body).Decode(session)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
		return
	}

	lastResource := interfaces.GetLastResource(r)
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	relations := interfaces.GetOwnerRelations(r)

	session.SetRelatedID(lastResource.IDKey, lastResource.ID)
	session, err = c.interactor.UpdateByID(id, session, usecases.QueryContext{Filter: filter, OwnerRelations: relations})

	if err != nil {
		switch err {
		case internalerrors.NotFound:
			c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
		default:
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		}
		return
	}

	session.BeforeRender()
	c.render.JSON(w, http.StatusCreated, session)
}

// @Title DeleteAll
// @Description Delete all Session instances matched by filter
// @Accept  json
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 204 {object} error "Request was successful"
// @Router /sessions [delete]
func (c *SessionCtrl) DeleteAll(w http.ResponseWriter, r *http.Request, _ map[string]string) {
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
// @Description Delete a Session instance
// @Accept  json
// @Param   id path int true "Session id"
// @Success 204 {object} error "Request was successful"
// @Router /sessions/{id} [delete]
func (c *SessionCtrl) DeleteByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
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
// @Description Create one or multiple Session instances of a related resource
// @Accept  json
// @Param   pk path int true "Session id"
// @Param   relatedResource path string true "Related resource name"
// @Param   Session body domain.Session true "Session instance(s) data"
// @Success 201 {object} domain.Session "Request was successful"
// @Router /sessions/{pk}/{relatedResource} [post]
func (c *SessionCtrl) CreateRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.related(w, r, params)
}

// @Title FindRelated
// @Description Find all Session instances  of a related resource matched by filter
// @Accept  json
// @Param   pk path int true "Session id"
// @Param   relatedResource path string true "Related resource name"
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 200 {object} domain.Session "Request was successful"
// @Router /sessions/{pk}/{relatedResource} [get]
func (c *SessionCtrl) FindRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.related(w, r, params)
}

// @Title UpsertRelated
// @Description Upsert one or multiple Session instances of a related resource
// @Accept  json
// @Param   pk path int true "Session id"
// @Param   relatedResource path string true "Related resource name"
// @Param   Session body domain.Session true "Session instance(s) data"
// @Success 201 {object} domain.Session "Request was successful"
// @Router /sessions/{pk}/{relatedResource} [put]
func (c *SessionCtrl) UpsertRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.related(w, r, params)
}

// @Title DeleteAllRelated
// @Description Delete all Session instances of a related resource matched by filter
// @Accept  json
// @Param   pk path int true "Session id"
// @Param   relatedResource path string true "Related resource name"
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 204 {object} error "Request was successful"
// @Router /sessions/{pk}/{relatedResource} [delete]
func (c *SessionCtrl) DeleteAllRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.related(w, r, params)
}

func (c *SessionCtrl) related(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	context.Set(r, "lastResource", &interfaces.Resource{Name: related, IDKey: "sessionID", ID: pk})

	handler(w, r, params)
}

// @Title FindByIDRelated
// @Description Find a Session instance of a related resource
// @Accept  json
// @Param   pk path int true "Session id"
// @Param   relatedResource path string true "Related resource name"
// @Param   fk path int true "Related resource id"
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 200 {object} domain.Session "Request was successful"
// @Router /sessions/{pk}/{relatedResource}/{fk} [get]
func (c *SessionCtrl) FindByIDRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.relatedOne(w, r, params)
}

// @Title UpdateByIDRelated
// @Description Update attributes of a Session instance of a related resource
// @Accept  json
// @Param   pk path int true "Session id"
// @Param   relatedResource path string true "Related resource name"
// @Param   fk path int true "Related resource id"
// @Param   Session body domain.Session true "Session instance data"
// @Success 201 {object} domain.Session
// @Router /sessions/{pk}/{relatedResource}/{fk} [put]
func (c *SessionCtrl) UpdateByIDRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.relatedOne(w, r, params)
}

// @Title DeleteByIDRelated
// @Description Delete a Session instance of a related resource
// @Accept  json
// @Param   pk path int true "Session id"
// @Param   relatedResource path string true "Related resource name"
// @Param   fk path int true "Related resource id"
// @Success 204 {object} error "Request was successful"
// @Router /sessions/{pk}/{relatedResource}/{fk} [delete]
func (c *SessionCtrl) DeleteByIDRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.relatedOne(w, r, params)
}

func (c *SessionCtrl) relatedOne(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	context.Set(r, "lastResource", &interfaces.Resource{Name: related, IDKey: "sessionID", ID: pk})

	handler(w, r, params)
}
