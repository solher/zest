// @SubApi Account resource [/accounts]
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
	usecases.DependencyDirectory.Register(NewAccountCtrl)
}

type Credentials struct {
	Email      string `json:"email,omitempty"`
	Password   string `json:"password,omitempty"`
	RememberMe bool   `json:"rememberMe,omitempty"`
}

type AbstractAccountInter interface {
	Create(accounts []domain.Account) ([]domain.Account, error)
	CreateOne(account *domain.Account) (*domain.Account, error)
	Find(context usecases.QueryContext) ([]domain.Account, error)
	FindByID(id int, context usecases.QueryContext) (*domain.Account, error)
	Upsert(accounts []domain.Account, context usecases.QueryContext) ([]domain.Account, error)
	UpsertOne(account *domain.Account, context usecases.QueryContext) (*domain.Account, error)
	UpdateByID(id int, account *domain.Account, context usecases.QueryContext) (*domain.Account, error)
	DeleteAll(context usecases.QueryContext) error
	DeleteByID(id int, context usecases.QueryContext) error
}

type AbstractAccountGuestInter interface {
	Signin(ip, userAgent string, credentials *Credentials) (*domain.Session, error)
	Signout(currentSession *domain.Session) error
	Signup(user *domain.User) (*domain.Account, error)
	Current(currentSession *domain.Session) (*domain.Account, error)
	CurrentSessionFromToken(authToken string) (*domain.Session, error)
}

type AccountCtrl struct {
	interactor AbstractAccountInter
	guestInter AbstractAccountGuestInter
	render     interfaces.AbstractRender
	routeDir   *usecases.RouteDirectory
}

func NewAccountCtrl(interactor AbstractAccountInter, guestInter AbstractAccountGuestInter, render interfaces.AbstractRender, routeDir *usecases.RouteDirectory) *AccountCtrl {
	controller := &AccountCtrl{interactor: interactor, guestInter: guestInter, render: render, routeDir: routeDir}

	if routeDir != nil {
		setAccountRoutes(routeDir, controller)
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

	session, err := c.guestInter.Signin(r.RemoteAddr, r.UserAgent(), &credentials)
	if err != nil {
		c.render.JSONError(w, http.StatusUnauthorized, apierrors.InvalidCredentials, err)
		return
	}

	cookie := http.Cookie{Name: "authToken", Value: session.AuthToken, Expires: session.ValidTo, Path: "/"}
	http.SetCookie(w, &cookie)

	session.BeforeRender()
	c.render.JSON(w, http.StatusCreated, session)
}

func (c *AccountCtrl) Signout(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	sessionCtx := context.Get(r, "currentSession")

	if sessionCtx == nil {
		c.render.JSONError(w, http.StatusUnauthorized, apierrors.SessionNotFound, nil)
		return
	}
	session := sessionCtx.(domain.Session)

	err := c.guestInter.Signout(&session)

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

	account, err := c.guestInter.Signup(&user)
	if err != nil {
		switch err.(type) {
		case *internalerrors.ViolatedConstraint:
			c.render.JSONError(w, 422, apierrors.AlreadyExistingEmail, err)
		default:
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		}
		return
	}

	account.BeforeRender()
	c.render.JSON(w, http.StatusCreated, account)
}

func (c *AccountCtrl) Current(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	sessionCtx := context.Get(r, "currentSession")

	if sessionCtx == nil {
		c.render.JSONError(w, http.StatusUnauthorized, apierrors.SessionNotFound, nil)
		return
	}
	session := sessionCtx.(domain.Session)

	account, err := c.guestInter.Current(&session)
	if err != nil {
		c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		return
	}

	account.BeforeRender()
	c.render.JSON(w, http.StatusOK, account)
}

// @Title Create
// @Description Create one or multiple Account instances
// @Accept  json
// @Param   Account body domain.Account true "Account instance(s) data"
// @Success 201 {object} domain.Account "Request was successful"
// @Router /accounts [post]
func (c *AccountCtrl) Create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	account := &domain.Account{}
	var accounts []domain.Account

	buffer, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(buffer, account)
	if err != nil {
		err := json.Unmarshal(buffer, &accounts)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
			return
		}
	}

	lastResource := interfaces.GetLastResource(r)

	if accounts == nil {
		account.SetRelatedID(lastResource.IDKey, lastResource.ID)
		account, err = c.interactor.CreateOne(account)
	} else {
		for i := range accounts {
			(&accounts[i]).SetRelatedID(lastResource.IDKey, lastResource.ID)
		}
		accounts, err = c.interactor.Create(accounts)
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

	if accounts == nil {
		account.BeforeRender()
		c.render.JSON(w, http.StatusCreated, account)
	} else {
		for i := range accounts {
			(&accounts[i]).BeforeRender()
		}
		c.render.JSON(w, http.StatusCreated, accounts)
	}
}

// @Title Find
// @Description Find all Account instances matched by filter
// @Accept  json
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 200 {object} domain.Account "Request was successful"
// @Router /accounts [get]
func (c *AccountCtrl) Find(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	filter, err := interfaces.GetQueryFilter(r)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.FilterDecodingError, err)
		return
	}

	filter = interfaces.FilterIfLastResource(r, filter)
	filter = interfaces.FilterIfOwnerRelations(r, filter)
	relations := interfaces.GetOwnerRelations(r)

	accounts, err := c.interactor.Find(usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	if err != nil {
		c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		return
	}

	for i := range accounts {
		(&accounts[i]).BeforeRender()
	}
	c.render.JSON(w, http.StatusOK, accounts)
}

// @Title FindByID
// @Description Find a Account instance
// @Accept  json
// @Param   id path int true "Account id"
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 200 {object} domain.Account "Request was successful"
// @Router /accounts/{id} [get]
func (c *AccountCtrl) FindByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var (
		id  int
		err error
	)

	if params["id"] == "me" {
		sessionCtx := context.Get(r, "currentSession")
		if sessionCtx == nil {
			c.render.JSONError(w, http.StatusUnauthorized, apierrors.SessionNotFound, nil)
			return
		}

		id = sessionCtx.(domain.Session).AccountID
	} else {
		id, err = strconv.Atoi(params["id"])
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.InvalidPathParams, err)
			return
		}
	}

	filter, err := interfaces.GetQueryFilter(r)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.FilterDecodingError, err)
		return
	}

	filter = interfaces.FilterIfOwnerRelations(r, filter)
	relations := interfaces.GetOwnerRelations(r)

	account, err := c.interactor.FindByID(id, usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	if err != nil {
		switch err {
		case internalerrors.NotFound:
			c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
		default:
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		}
		return
	}

	account.BeforeRender()
	c.render.JSON(w, http.StatusOK, account)
}

// @Title Upsert
// @Description Upsert one or multiple Account instances
// @Accept  json
// @Param   Account body domain.Account true "Account instance(s) data"
// @Success 201 {object} domain.Account "Request was successful"
// @Router /accounts [put]
func (c *AccountCtrl) Upsert(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	account := &domain.Account{}
	var accounts []domain.Account

	buffer, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(buffer, account)
	if err != nil {
		err := json.Unmarshal(buffer, &accounts)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
			return
		}
	}

	lastResource := interfaces.GetLastResource(r)
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	relations := interfaces.GetOwnerRelations(r)

	if accounts == nil {
		account.SetRelatedID(lastResource.IDKey, lastResource.ID)
		account, err = c.interactor.UpsertOne(account, usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	} else {
		for i := range accounts {
			(&accounts[i]).SetRelatedID(lastResource.IDKey, lastResource.ID)
		}
		accounts, err = c.interactor.Upsert(accounts, usecases.QueryContext{Filter: filter, OwnerRelations: relations})
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

	if accounts == nil {
		account.BeforeRender()
		c.render.JSON(w, http.StatusCreated, account)
	} else {
		for i := range accounts {
			(&accounts[i]).BeforeRender()
		}
		c.render.JSON(w, http.StatusCreated, accounts)
	}
}

// @Title UpdateByID
// @Description Update attributes of a Account instance
// @Accept  json
// @Param   id path int true "Account id"
// @Param   Account body domain.Account true "Account instance data"
// @Success 201 {object} domain.Account
// @Router /accounts/{id} [put]
func (c *AccountCtrl) UpdateByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var (
		id  int
		err error
	)

	if params["id"] == "me" {
		sessionCtx := context.Get(r, "currentSession")
		if sessionCtx == nil {
			c.render.JSONError(w, http.StatusUnauthorized, apierrors.SessionNotFound, nil)
			return
		}

		id = sessionCtx.(domain.Session).AccountID
	} else {
		id, err = strconv.Atoi(params["id"])
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.InvalidPathParams, err)
			return
		}
	}

	account := &domain.Account{}

	err = json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
		return
	}

	lastResource := interfaces.GetLastResource(r)
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	relations := interfaces.GetOwnerRelations(r)

	account.SetRelatedID(lastResource.IDKey, lastResource.ID)
	account, err = c.interactor.UpdateByID(id, account, usecases.QueryContext{Filter: filter, OwnerRelations: relations})

	if err != nil {
		switch err {
		case internalerrors.NotFound:
			c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
		default:
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		}
		return
	}

	account.BeforeRender()
	c.render.JSON(w, http.StatusCreated, account)
}

// @Title DeleteAll
// @Description Delete all Account instances matched by filter
// @Accept  json
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 204 {object} error "Request was successful"
// @Router /accounts [delete]
func (c *AccountCtrl) DeleteAll(w http.ResponseWriter, r *http.Request, _ map[string]string) {
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
// @Description Delete a Account instance
// @Accept  json
// @Param   id path int true "Account id"
// @Success 204 {object} error "Request was successful"
// @Router /accounts/{id} [delete]
func (c *AccountCtrl) DeleteByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var (
		id  int
		err error
	)

	if params["id"] == "me" {
		sessionCtx := context.Get(r, "currentSession")
		if sessionCtx == nil {
			c.render.JSONError(w, http.StatusUnauthorized, apierrors.SessionNotFound, nil)
			return
		}

		id = sessionCtx.(domain.Session).AccountID
	} else {
		id, err = strconv.Atoi(params["id"])
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.InvalidPathParams, err)
			return
		}
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
// @Description Create one or multiple Account instances of a related resource
// @Accept  json
// @Param   Account body domain.Account true "Account instance(s) data"
// @Success 201 {object} domain.Account "Request was successful"
// @Router /accounts/{id}/{relatedResource} [post]
func CreateRelated() {}

// @Title FindRelated
// @Description Find all Account instances  of a related resource matched by filter
// @Accept  json
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 200 {object} domain.Account "Request was successful"
// @Router /accounts/{id}/{relatedResource} [get]
func FindRelated() {}

// @Title UpsertRelated
// @Description Upsert one or multiple Account instances of a related resource
// @Accept  json
// @Param   Account body domain.Account true "Account instance(s) data"
// @Success 201 {object} domain.Account "Request was successful"
// @Router /accounts/{id}/{relatedResource} [put]
func UpsertRelated() {}

// @Title DeleteAllRelated
// @Description Delete all Account instances of a related resource matched by filter
// @Accept  json
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 204 {object} error "Request was successful"
// @Router /accounts/{id}/{relatedResource} [delete]
func DeleteAllRelated() {}

func (c *AccountCtrl) Related(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var (
		pk  int
		err error
	)

	if params["pk"] == "me" {
		sessionCtx := context.Get(r, "currentSession")
		if sessionCtx == nil {
			c.render.JSONError(w, http.StatusUnauthorized, apierrors.SessionNotFound, nil)
			return
		}

		pk = sessionCtx.(domain.Session).AccountID
	} else {
		pk, err = strconv.Atoi(params["pk"])
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.InvalidPathParams, err)
			return
		}
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

	context.Set(r, "lastResource", &interfaces.Resource{Name: related, IDKey: "accountID", ID: pk})

	handler(w, r, params)
}

// @Title FindByIDRelated
// @Description Find a Account instance of a related resource
// @Accept  json
// @Param   id path int true "Account id"
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 200 {object} domain.Account "Request was successful"
// @Router /accounts/{pk}/{relatedResource}/{fk} [get]
func FindByIDRelated() {}

// @Title UpdateByIDRelated
// @Description Update attributes of a Account instance of a related resource
// @Accept  json
// @Param   id path int true "Account id"
// @Param   Account body domain.Account true "Account instance data"
// @Success 201 {object} domain.Account
// @Router /accounts/{pk}/{relatedResource}/{fk} [put]
func UpdateByIDRelated() {}

// @Title DeleteByIDRelated
// @Description Delete a Account instance of a related resource
// @Accept  json
// @Param   id path int true "Account id"
// @Success 204 {object} error "Request was successful"
// @Router /accounts/{pk}/{relatedResource}/{fk} [delete]
func DeleteByIDRelated() {}

func (c *AccountCtrl) RelatedOne(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var (
		pk  int
		err error
	)

	if params["pk"] == "me" {
		sessionCtx := context.Get(r, "currentSession")
		if sessionCtx == nil {
			c.render.JSONError(w, http.StatusUnauthorized, apierrors.SessionNotFound, nil)
			return
		}

		pk = sessionCtx.(domain.Session).AccountID
	} else {
		pk, err = strconv.Atoi(params["pk"])
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.InvalidPathParams, err)
			return
		}
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
	case "PUT":
		handler = c.routeDir.Get(key.For("UpdateByID")).EffectiveHandler
	}

	if handler == nil {
		c.render.JSON(w, http.StatusNotFound, nil)
		return
	}

	context.Set(r, "lastResource", &interfaces.Resource{Name: related, IDKey: "accountID", ID: pk})

	handler(w, r, params)
}
