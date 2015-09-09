// @SubApi User resource [/users]
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
	"github.com/solher/zest/utils"
)

func init() {
	usecases.DependencyDirectory.Register(NewUserCtrl)
}

type AbstractUserInter interface {
	Create(users []domain.User) ([]domain.User, error)
	CreateOne(user *domain.User) (*domain.User, error)
	Find(context usecases.QueryContext) ([]domain.User, error)
	FindByID(id int, context usecases.QueryContext) (*domain.User, error)
	Upsert(users []domain.User, context usecases.QueryContext) ([]domain.User, error)
	UpsertOne(user *domain.User, context usecases.QueryContext) (*domain.User, error)
	UpdateByID(id int, user *domain.User, context usecases.QueryContext) (*domain.User, error)
	DeleteAll(context usecases.QueryContext) error
	DeleteByID(id int, context usecases.QueryContext) error
}

type AbstractGuestUserInter interface {
	UpdateByID(id int, user *domain.User, context usecases.QueryContext) (*domain.User, error)
	UpdatePassword(id int, context usecases.QueryContext, oldPassword, newPassword string) (*domain.User, error)
}

type UserCtrl struct {
	interactor AbstractUserInter
	guestInter AbstractGuestUserInter
	render     interfaces.AbstractRender
	routeDir   *usecases.RouteDirectory
}

func NewUserCtrl(interactor AbstractUserInter, guestInter AbstractGuestUserInter, render interfaces.AbstractRender, routeDir *usecases.RouteDirectory) *UserCtrl {
	controller := &UserCtrl{interactor: interactor, guestInter: guestInter, render: render, routeDir: routeDir}

	if routeDir != nil {
		setUserRoutes(routeDir, controller)
	}

	return controller
}

type PasswordForm struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// @Title UpdatePassword
// @Description Update the user password
// @Accept  json
// @Param   id path int true "User id"
// @Param   PasswordForm body PasswordForm true "The old and the new password"
// @Success 200 {object} domain.User "Request was successful"
// @Router /users/{id}/updatePassword [post]
func (c *UserCtrl) UpdatePassword(w http.ResponseWriter, r *http.Request, params map[string]string) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.InvalidPathParams, err)
		return
	}

	form := &PasswordForm{}

	err = json.NewDecoder(r.Body).Decode(form)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
		return
	}

	filter := interfaces.FilterIfOwnerRelations(r, nil)
	relations := interfaces.GetOwnerRelations(r)

	user, err := c.guestInter.UpdatePassword(id, usecases.QueryContext{Filter: filter, OwnerRelations: relations}, form.OldPassword, form.NewPassword)

	if err != nil {
		switch err {
		case internalerrors.NotFound:
			c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
		case internalerrors.InvalidCredentials:
			c.render.JSONError(w, http.StatusUnauthorized, apierrors.InvalidCredentials, err)
		default:
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		}
		return
	}

	user.BeforeRender()
	c.render.JSON(w, http.StatusOK, user)
}

// @Title Create
// @Description Create one or multiple User instances
// @Accept  json
// @Param   User body domain.User true "User instance(s) data"
// @Success 201 {object} domain.User "Request was successful"
// @Router /users [post]
func (c *UserCtrl) Create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	user := &domain.User{}
	var users []domain.User

	buffer, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(buffer, user)
	if err != nil {
		err := json.Unmarshal(buffer, &users)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
			return
		}
	}

	lastResource := interfaces.GetLastResource(r)

	if users == nil {
		user.SetRelatedID(lastResource.IDKey, lastResource.ID)
		user, err = c.interactor.CreateOne(user)
	} else {
		for i := range users {
			(&users[i]).SetRelatedID(lastResource.IDKey, lastResource.ID)
		}
		users, err = c.interactor.Create(users)
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

	if users == nil {
		user.BeforeRender()
		c.render.JSON(w, http.StatusCreated, user)
	} else {
		for i := range users {
			(&users[i]).BeforeRender()
		}
		c.render.JSON(w, http.StatusCreated, users)
	}
}

// @Title Find
// @Description Find all User instances matched by filter
// @Accept  json
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 200 {object} domain.User "Request was successful"
// @Router /users [get]
func (c *UserCtrl) Find(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	filter, err := interfaces.GetQueryFilter(r)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.FilterDecodingError, err)
		return
	}

	filter = interfaces.FilterIfLastResource(r, filter)
	filter = interfaces.FilterIfOwnerRelations(r, filter)
	relations := interfaces.GetOwnerRelations(r)

	users, err := c.interactor.Find(usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	if err != nil {
		c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		return
	}

	for i := range users {
		(&users[i]).BeforeRender()
	}
	c.render.JSON(w, http.StatusOK, users)
}

// @Title FindByID
// @Description Find a User instance
// @Accept  json
// @Param   id path int true "User id"
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 200 {object} domain.User "Request was successful"
// @Router /users/{id} [get]
func (c *UserCtrl) FindByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	user, err := c.interactor.FindByID(id, usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	if err != nil {
		switch err {
		case internalerrors.NotFound:
			c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
		default:
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		}
		return
	}

	user.BeforeRender()
	c.render.JSON(w, http.StatusOK, user)
}

// @Title Upsert
// @Description Upsert one or multiple User instances
// @Accept  json
// @Param   User body domain.User true "User instance(s) data"
// @Success 201 {object} domain.User "Request was successful"
// @Router /users [put]
func (c *UserCtrl) Upsert(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	user := &domain.User{}
	var users []domain.User

	buffer, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(buffer, user)
	if err != nil {
		err := json.Unmarshal(buffer, &users)
		if err != nil {
			c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
			return
		}
	}

	lastResource := interfaces.GetLastResource(r)
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	relations := interfaces.GetOwnerRelations(r)

	if users == nil {
		user.SetRelatedID(lastResource.IDKey, lastResource.ID)
		user, err = c.interactor.UpsertOne(user, usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	} else {
		for i := range users {
			(&users[i]).SetRelatedID(lastResource.IDKey, lastResource.ID)
		}
		users, err = c.interactor.Upsert(users, usecases.QueryContext{Filter: filter, OwnerRelations: relations})
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

	if users == nil {
		user.BeforeRender()
		c.render.JSON(w, http.StatusCreated, user)
	} else {
		for i := range users {
			(&users[i]).BeforeRender()
		}
		c.render.JSON(w, http.StatusCreated, users)
	}
}

// @Title UpdateByID
// @Description Update attributes of a User instance
// @Accept  json
// @Param   id path int true "User id"
// @Param   User body domain.User true "User instance data"
// @Success 200 {object} domain.User
// @Router /users/{id} [put]
func (c *UserCtrl) UpdateByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.InvalidPathParams, err)
		return
	}

	user := &domain.User{}

	err = json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.BodyDecodingError, err)
		return
	}

	lastResource := interfaces.GetLastResource(r)
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	relations := interfaces.GetOwnerRelations(r)

	user.SetRelatedID(lastResource.IDKey, lastResource.ID)

	if roles := context.Get(r, "roles"); roles != nil && utils.ContainsStr(roles.([]string), "Admin") {
		user, err = c.interactor.UpdateByID(id, user, usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	} else {
		user, err = c.guestInter.UpdateByID(id, user, usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	}

	if err != nil {
		switch err {
		case internalerrors.NotFound:
			c.render.JSONError(w, http.StatusUnauthorized, apierrors.Unauthorized, err)
		default:
			c.render.JSONError(w, http.StatusInternalServerError, apierrors.InternalServerError, err)
		}
		return
	}

	user.BeforeRender()
	c.render.JSON(w, http.StatusOK, user)
}

// @Title DeleteAll
// @Description Delete all User instances matched by filter
// @Accept  json
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 204 {object} error "Request was successful"
// @Router /users [delete]
func (c *UserCtrl) DeleteAll(w http.ResponseWriter, r *http.Request, _ map[string]string) {
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
// @Description Delete a User instance
// @Accept  json
// @Param   id path int true "User id"
// @Success 204 {object} error "Request was successful"
// @Router /users/{id} [delete]
func (c *UserCtrl) DeleteByID(w http.ResponseWriter, r *http.Request, params map[string]string) {
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
// @Description Create one or multiple User instances of a related resource
// @Accept  json
// @Param   pk path int true "User id"
// @Param   relatedResource path string true "Related resource name"
// @Param   User body domain.User true "User instance(s) data"
// @Success 201 {object} domain.User "Request was successful"
// @Router /users/{pk}/{relatedResource} [post]
func (c *UserCtrl) CreateRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.related(w, r, params)
}

// @Title FindRelated
// @Description Find all User instances  of a related resource matched by filter
// @Accept  json
// @Param   pk path int true "User id"
// @Param   relatedResource path string true "Related resource name"
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 200 {object} domain.User "Request was successful"
// @Router /users/{pk}/{relatedResource} [get]
func (c *UserCtrl) FindRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.related(w, r, params)
}

// @Title UpsertRelated
// @Description Upsert one or multiple User instances of a related resource
// @Accept  json
// @Param   pk path int true "User id"
// @Param   relatedResource path string true "Related resource name"
// @Param   User body domain.User true "User instance(s) data"
// @Success 201 {object} domain.User "Request was successful"
// @Router /users/{pk}/{relatedResource} [put]
func (c *UserCtrl) UpsertRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.related(w, r, params)
}

// @Title DeleteAllRelated
// @Description Delete all User instances of a related resource matched by filter
// @Accept  json
// @Param   pk path int true "User id"
// @Param   relatedResource path string true "Related resource name"
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 204 {object} error "Request was successful"
// @Router /users/{pk}/{relatedResource} [delete]
func (c *UserCtrl) DeleteAllRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.related(w, r, params)
}

func (c *UserCtrl) related(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	context.Set(r, "lastResource", &interfaces.Resource{Name: related, IDKey: "userID", ID: pk})

	handler(w, r, params)
}

// @Title FindByIDRelated
// @Description Find a User instance of a related resource
// @Accept  json
// @Param   pk path int true "User id"
// @Param   relatedResource path string true "Related resource name"
// @Param   fk path int true "Related resource id"
// @Param   filter query string false "JSON filter defining fields and includes"
// @Success 200 {object} domain.User "Request was successful"
// @Router /users/{pk}/{relatedResource}/{fk} [get]
func (c *UserCtrl) FindByIDRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.relatedOne(w, r, params)
}

// @Title UpdateByIDRelated
// @Description Update attributes of a User instance of a related resource
// @Accept  json
// @Param   pk path int true "User id"
// @Param   relatedResource path string true "Related resource name"
// @Param   fk path int true "Related resource id"
// @Param   User body domain.User true "User instance data"
// @Success 201 {object} domain.User
// @Router /users/{pk}/{relatedResource}/{fk} [put]
func (c *UserCtrl) UpdateByIDRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.relatedOne(w, r, params)
}

// @Title DeleteByIDRelated
// @Description Delete a User instance of a related resource
// @Accept  json
// @Param   pk path int true "User id"
// @Param   relatedResource path string true "Related resource name"
// @Param   fk path int true "Related resource id"
// @Success 204 {object} error "Request was successful"
// @Router /users/{pk}/{relatedResource}/{fk} [delete]
func (c *UserCtrl) DeleteByIDRelated(w http.ResponseWriter, r *http.Request, params map[string]string) {
	c.relatedOne(w, r, params)
}

func (c *UserCtrl) relatedOne(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	context.Set(r, "lastResource", &interfaces.Resource{Name: related, IDKey: "userID", ID: pk})

	handler(w, r, params)
}
