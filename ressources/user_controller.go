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

type UserCtrl struct {
	interactor AbstractUserInter
	render     interfaces.AbstractRender
	routeDir   *usecases.RouteDirectory
}

func NewUserCtrl(interactor AbstractUserInter, render interfaces.AbstractRender, routeDir *usecases.RouteDirectory) *UserCtrl {
	controller := &UserCtrl{interactor: interactor, render: render, routeDir: routeDir}

	if routeDir != nil {
		setUserRoutes(routeDir, controller)
	}

	return controller
}

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

	lastRessource := interfaces.GetLastRessource(r)

	if users == nil {
		user.SetRelatedID(lastRessource.IDKey, lastRessource.ID)
		user, err = c.interactor.CreateOne(user)
	} else {
		for i := range users {
			(&users[i]).SetRelatedID(lastRessource.IDKey, lastRessource.ID)
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

func (c *UserCtrl) Find(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	filter, err := interfaces.GetQueryFilter(r)
	if err != nil {
		c.render.JSONError(w, http.StatusBadRequest, apierrors.FilterDecodingError, err)
		return
	}

	filter = interfaces.FilterIfLastRessource(r, filter)
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

	lastRessource := interfaces.GetLastRessource(r)
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	relations := interfaces.GetOwnerRelations(r)

	if users == nil {
		user.SetRelatedID(lastRessource.IDKey, lastRessource.ID)
		user, err = c.interactor.UpsertOne(user, usecases.QueryContext{Filter: filter, OwnerRelations: relations})
	} else {
		for i := range users {
			(&users[i]).SetRelatedID(lastRessource.IDKey, lastRessource.ID)
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

	lastRessource := interfaces.GetLastRessource(r)
	filter := interfaces.FilterIfOwnerRelations(r, nil)
	relations := interfaces.GetOwnerRelations(r)

	user.SetRelatedID(lastRessource.IDKey, lastRessource.ID)
	user, err = c.interactor.UpdateByID(id, user, usecases.QueryContext{Filter: filter, OwnerRelations: relations})

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
	c.render.JSON(w, http.StatusCreated, user)
}

func (c *UserCtrl) DeleteAll(w http.ResponseWriter, r *http.Request, _ map[string]string) {
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

func (c *UserCtrl) Related(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	context.Set(r, "lastRessource", &interfaces.Ressource{Name: related, IDKey: "userID", ID: pk})

	handler(w, r, params)
}

func (c *UserCtrl) RelatedOne(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	context.Set(r, "lastRessource", &interfaces.Ressource{Name: related, IDKey: "userID", ID: pk})

	handler(w, r, params)
}
