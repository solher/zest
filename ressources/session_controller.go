// Generated by: main
// TypeWriter: controller
// Directive: +gen on Session

package ressources

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Solher/auth-scaffold/apierrors"
	"github.com/Solher/auth-scaffold/interfaces"
	"github.com/julienschmidt/httprouter"
)

type AbstractSessionInter interface {
	Create(sessions []Session) ([]Session, error)
	CreateOne(session *Session) (*Session, error)
	Find(filter *interfaces.Filter) ([]Session, error)
	FindByID(id int, filter *interfaces.Filter) (*Session, error)
	Upsert(sessions []Session) ([]Session, error)
	UpsertOne(session *Session) (*Session, error)
	DeleteAll(filter *interfaces.Filter) error
	DeleteByID(id int) error
	CurrentFromToken(authToken string) (*Session, error)
}

type SessionCtrl struct {
	interactor AbstractSessionInter
	render     interfaces.Render
}

func NewSessionCtrl(interactor AbstractSessionInter, render interfaces.Render, routesDir interfaces.RouteDirectory) *SessionCtrl {
	controller := &SessionCtrl{interactor: interactor, render: render}

	if routesDir != nil {
		addSessionRoutes(routesDir, controller)
	}

	return controller
}

func (c *SessionCtrl) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session := &Session{}
	var sessions []Session

	buffer, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(buffer, session)
	if err != nil {
		err := json.Unmarshal(buffer, &sessions)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
			return
		}
	}

	if sessions == nil {
		session.ScopeModel()
		session, err = c.interactor.CreateOne(session)
		if err != nil {
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
			return
		}

		c.render.JSON(w, http.StatusCreated, session)
	} else {
		sessions, err = c.interactor.Create(sessions)
		if err != nil {
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
			return
		}

		c.render.JSON(w, http.StatusCreated, sessions)
	}
}

func (c *SessionCtrl) Find(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	filter, err := interfaces.GetQueryFilter(r)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.FilterDecodingError, err)
		return
	}

	sessions, err := c.interactor.Find(filter)
	if err != nil {
		c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		return
	}

	c.render.JSON(w, http.StatusOK, sessions)
}

func (c *SessionCtrl) FindByID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.InvalidPathParams, err)
		return
	}

	filter, err := interfaces.GetQueryFilter(r)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.FilterDecodingError, err)
		return
	}

	session, err := c.interactor.FindByID(id, filter)
	if err != nil {
		c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
		return
	}

	c.render.JSON(w, http.StatusOK, session)
}

func (c *SessionCtrl) Upsert(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session := &Session{}
	var sessions []Session

	buffer, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(buffer, session)
	if err != nil {
		err := json.Unmarshal(buffer, &sessions)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
			return
		}
	}

	if sessions == nil {
		session, err = c.interactor.UpsertOne(session)
		if err != nil {
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
			return
		}

		c.render.JSON(w, http.StatusCreated, session)
	} else {
		sessions, err = c.interactor.Upsert(sessions)
		if err != nil {
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
			return
		}

		c.render.JSON(w, http.StatusCreated, sessions)
	}
}

func (c *SessionCtrl) DeleteAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	filter, err := interfaces.GetQueryFilter(r)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.FilterDecodingError, err)
		return
	}

	err = c.interactor.DeleteAll(filter)
	if err != nil {
		c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		return
	}

	c.render.JSON(w, http.StatusNoContent, nil)
}

func (c *SessionCtrl) DeleteByID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.InvalidPathParams, err)
		return
	}

	err = c.interactor.DeleteByID(id)
	if err != nil {
		c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
		return
	}

	c.render.JSON(w, http.StatusNoContent, nil)
}
