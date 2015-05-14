package ressources

import (
	"encoding/json"
	"net/http"

	"github.com/Solher/auth-scaffold/apierrors"
	"github.com/Solher/auth-scaffold/interfaces"
	"github.com/Solher/auth-scaffold/internalerrors"
	"github.com/gorilla/context"

	"github.com/julienschmidt/httprouter"
)

type Credentials struct {
	Email      string `json:"email,omitempty"`
	Password   string `json:"password,omitempty"`
	RememberMe bool   `json:"rememberMe,omitempty"`
}

type AbstractAccountInter interface {
	Signin(ip, userAgent string, credentials *Credentials) (*Session, error)
	Signout(currentSession *Session) error
	Signup(user *User) (*Account, error)
	Current(currentSession *Session) (*Account, error)
}

type AccountCtrl struct {
	interactor AbstractAccountInter
	render     interfaces.Render
}

func NewAccountCtrl(interactor AbstractAccountInter, render interfaces.Render, routesDir interfaces.RouteDirectory) *AccountCtrl {
	controller := &AccountCtrl{interactor: interactor, render: render}

	if routesDir != nil {
		addAccountRoutes(routesDir, controller)
	}

	return controller
}

func (c *AccountCtrl) Signin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var credentials Credentials

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
		return
	}

	if credentials.Password == "" {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.BlankParam("password"), err)
		return
	}

	if credentials.Email == "" {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.BlankParam("email"), err)
		return
	}

	session, err := c.interactor.Signin(r.Host, r.UserAgent(), &credentials)
	if err != nil {
		c.render.JSONError(w, http.StatusUnauthorized, apierrors.InvalidCredentials, err)
		return
	}

	cookie := http.Cookie{Name: "authToken", Value: session.AuthToken, Expires: session.ValidTo}
	http.SetCookie(w, &cookie)

	c.render.JSON(w, http.StatusCreated, session)
}

func (c *AccountCtrl) Signout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sessionCtx := context.Get(r, "currentSession")

	if sessionCtx == nil {
		c.render.JSONError(w, http.StatusUnauthorized, apierrors.SessionNotFound, nil)
		return
	}
	session := sessionCtx.(Session)

	err := c.interactor.Signout(&session)

	if err != nil {
		c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *AccountCtrl) Signup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type Params struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Password  string `json:"password"`
		Email     string `json:"email"`
	}
	var params Params

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
		return
	}

	if params.Password == "" {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.BlankParam("password"), err)
		return
	}

	if params.Email == "" {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.BlankParam("email"), err)
		return
	}

	user := User{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Password:  params.Password,
		Email:     params.Email,
	}

	account, err := c.interactor.Signup(&user)
	if err != nil {
		switch err.(type) {
		case *internalerrors.ViolatedConstraint:
			c.render.JSONError(w, 422, apierrors.AlreadyExistingEmail, err)
		default:
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		}
		return
	}

	c.render.JSON(w, http.StatusCreated, account)
}

func (c *AccountCtrl) Current(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sessionCtx := context.Get(r, "currentSession")

	if sessionCtx == nil {
		c.render.JSONError(w, http.StatusUnauthorized, apierrors.SessionNotFound, nil)
		return
	}
	session := sessionCtx.(Session)

	account, err := c.interactor.Current(&session)
	if err != nil {
		c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		return
	}

	c.render.JSON(w, http.StatusOK, account)
}
