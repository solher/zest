package ressources

import (
	"encoding/json"
	"net/http"

	"github.com/Solher/auth-scaffold/apierrors"
	"github.com/Solher/auth-scaffold/domain"
	"github.com/Solher/auth-scaffold/interfaces"
	"github.com/Solher/auth-scaffold/internalerrors"
	"github.com/dimfeld/httptreemux"
	"github.com/gorilla/context"
)

type Credentials struct {
	Email      string `json:"email,omitempty"`
	Password   string `json:"password,omitempty"`
	RememberMe bool   `json:"rememberMe,omitempty"`
}

type AbstractAccountInter interface {
	Signin(ip, userAgent string, credentials *Credentials) (*domain.Session, error)
	Signout(currentSession *domain.Session) error
	Signup(user *domain.User) (*domain.Account, error)
	Current(currentSession *domain.Session) (*domain.Account, error)
	CurrentSessionFromToken(authToken string) (*domain.Session, error)
}

type AccountCtrl struct {
	interactor AbstractAccountInter
	render     interfaces.AbstractRender
	routeDir   *interfaces.RouteDirectory
}

func NewAccountCtrl(interactor AbstractAccountInter, render interfaces.AbstractRender, routeDir *interfaces.RouteDirectory) *AccountCtrl {
	controller := &AccountCtrl{interactor: interactor, render: render, routeDir: routeDir}

	if routeDir != nil {
		setAccountAccess(routeDir, controller)
	}

	return controller
}

func (c *AccountCtrl) Signin(w http.ResponseWriter, r *http.Request, _ map[string]string) {
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

	session, err := c.interactor.Signin(r.RemoteAddr, r.UserAgent(), &credentials)
	if err != nil {
		c.render.JSONError(w, http.StatusUnauthorized, apierrors.InvalidCredentials, err)
		return
	}

	cookie := http.Cookie{Name: "authToken", Value: session.AuthToken, Expires: session.ValidTo, Path: "/"}
	http.SetCookie(w, &cookie)

	c.render.JSON(w, http.StatusCreated, session)
}

func (c *AccountCtrl) Signout(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	sessionCtx := context.Get(r, "currentSession")

	if sessionCtx == nil {
		c.render.JSONError(w, http.StatusUnauthorized, apierrors.SessionNotFound, nil)
		return
	}
	session := sessionCtx.(domain.Session)

	err := c.interactor.Signout(&session)

	if err != nil {
		c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *AccountCtrl) Signup(w http.ResponseWriter, r *http.Request, _ map[string]string) {
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

	user := domain.User{
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

func (c *AccountCtrl) Current(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	sessionCtx := context.Get(r, "currentSession")

	if sessionCtx == nil {
		c.render.JSONError(w, http.StatusUnauthorized, apierrors.SessionNotFound, nil)
		return
	}
	session := sessionCtx.(domain.Session)

	account, err := c.interactor.Current(&session)
	if err != nil {
		c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		return
	}

	c.render.JSON(w, http.StatusOK, account)
}

func (c *AccountCtrl) Related(w http.ResponseWriter, r *http.Request, params map[string]string) {
	sessionCtx := context.Get(r, "currentSession")

	if sessionCtx == nil {
		c.render.JSONError(w, http.StatusUnauthorized, apierrors.SessionNotFound, nil)
		return
	}
	session := sessionCtx.(domain.Session)

	related := params["related"]
	key := interfaces.NewDirectoryKey(related)

	var handler *httptreemux.HandlerFunc
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

	context.Set(r, "lastRessource", &interfaces.Ressource{Name: related, IDKey: "accountID", ID: session.AccountID})

	(*handler)(w, r, params)
}

func (c *AccountCtrl) RelatedOne(w http.ResponseWriter, r *http.Request, params map[string]string) {
	params["id"] = params["fk"]

	related := params["related"]
	key := interfaces.NewDirectoryKey(related)

	var handler httptreemux.HandlerFunc

	switch r.Method {
	case "GET":
		handler = *c.routeDir.Get(key.For("FindByID")).EffectiveHandler
	case "DELETE":
		handler = *c.routeDir.Get(key.For("DeleteByID")).EffectiveHandler
	}

	handler(w, r, params)
}
