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

	lastRessource := interfaces.GetLastRessource(r)

	if sessions == nil {
		session.SetRelatedID(lastRessource.IDKey, lastRessource.ID)
		session, err = c.interactor.CreateOne(session)
	} else {
		for i := range sessions {
			(&sessions[i]).SetRelatedID(lastRessource.IDKey, lastRessource.ID)
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

func (c *SessionCtrl) Find(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	filter, err := interfaces.GetQueryFilter(r)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.FilterDecodingError, err)
		return
	}

	filter = interfaces.FilterIfLastRessource(r, filter)
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

	lastRessource := interfaces.GetLastRessource(r)
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	relations := interfaces.GetOwnerRelations(r)

	if sessions == nil {
		session.SetRelatedID(lastRessource.IDKey, lastRessource.ID)
		session, err = c.interactor.UpsertOne(session, usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	} else {
		for i := range sessions {
			(&sessions[i]).SetRelatedID(lastRessource.IDKey, lastRessource.ID)
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

	lastRessource := interfaces.GetLastRessource(r)
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	relations := interfaces.GetOwnerRelations(r)

	session.SetRelatedID(lastRessource.IDKey, lastRessource.ID)
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

func (c *SessionCtrl) DeleteAll(w http.ResponseWriter, r *http.Request, _ map[string]string) {
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

func (c *SessionCtrl) Related(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	context.Set(r, "lastRessource", &interfaces.Ressource{Name: related, IDKey: "sessionID", ID: pk})

	handler(w, r, params)
}

func (c *SessionCtrl) RelatedOne(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	context.Set(r, "lastRessource", &interfaces.Ressource{Name: related, IDKey: "sessionID", ID: pk})

	handler(w, r, params)
}
