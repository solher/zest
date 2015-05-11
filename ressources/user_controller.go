// Generated by: main
// TypeWriter: controller
// Directive: +gen on User

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

type AbstractUserInter interface {
	Create(users []User) ([]User, error)
	CreateOne(user *User) (*User, error)
	Find(filter *interfaces.Filter) ([]User, error)
	FindByID(id int, filter *interfaces.Filter) (*User, error)
	Upsert(users []User) ([]User, error)
	UpsertOne(user *User) (*User, error)
	DeleteAll(filter *interfaces.Filter) error
	DeleteByID(id int) error
}

type UserCtrl struct {
	interactor AbstractUserInter
	render     interfaces.Render
}

func NewUserCtrl(interactor AbstractUserInter, render interfaces.Render, routesDir interfaces.RouteDirectory) *UserCtrl {
	controller := &UserCtrl{interactor: interactor, render: render}

	if routesDir != nil {
		addUserRoutes(routesDir, controller)
	}

	return controller
}

func (c *UserCtrl) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := &User{}
	var users []User

	buffer, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(buffer, user)
	if err != nil {
		err := json.Unmarshal(buffer, &users)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
			return
		}
	}

	if users == nil {
		user, err = c.interactor.CreateOne(user)
		if err != nil {
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
			return
		}

		c.render.JSON(w, http.StatusCreated, user)
	} else {
		users, err = c.interactor.Create(users)
		if err != nil {
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
			return
		}

		c.render.JSON(w, http.StatusCreated, users)
	}
}

func (c *UserCtrl) Find(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	filter, err := interfaces.GetQueryFilter(r)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.FilterDecodingError, err)
		return
	}

	users, err := c.interactor.Find(filter)
	if err != nil {
		c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		return
	}

	c.render.JSON(w, http.StatusOK, users)
}

func (c *UserCtrl) FindByID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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

	user, err := c.interactor.FindByID(id, filter)
	if err != nil {
		c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
		return
	}

	c.render.JSON(w, http.StatusOK, user)
}

func (c *UserCtrl) Upsert(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := &User{}
	var users []User

	buffer, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(buffer, user)
	if err != nil {
		err := json.Unmarshal(buffer, &users)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
			return
		}
	}

	if users == nil {
		user, err = c.interactor.UpsertOne(user)
		if err != nil {
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
			return
		}

		c.render.JSON(w, http.StatusCreated, user)
	} else {
		users, err = c.interactor.Upsert(users)
		if err != nil {
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
			return
		}

		c.render.JSON(w, http.StatusCreated, users)
	}
}

func (c *UserCtrl) DeleteAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

func (c *UserCtrl) DeleteByID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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
