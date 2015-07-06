package interfaces

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/context"
	"github.com/solher/zest/apierrors"
	"github.com/solher/zest/domain"
	"github.com/solher/zest/usecases"
)

func GetQueryFilter(r *http.Request) (*usecases.Filter, error) {
	param := r.URL.Query().Get("filter")
	if param == "" {
		return nil, nil
	}

	filterReader := strings.NewReader(param)

	filter := &usecases.Filter{}
	err := json.NewDecoder(filterReader).Decode(filter)
	if err != nil {
		return nil, err
	}

	return filter, nil
}

func MockHTTPRequest(route usecases.Route, body, filter string, params map[string]string) string {
	if route.Method == "" || route.Path == "" || route.Handler == nil {
		panic("Non existing or incomplete route when mocking HTTP request.")
	}

	w := httptest.NewRecorder()

	path := route.Path
	if filter != "" {
		path = path + "?filter=" + filter
	}

	req, _ := http.NewRequest(route.Method, path, bytes.NewBufferString(body))
	route.Handler(w, req, params)

	return w.Body.String()
}

func GetErrorCode(res string) string {
	apiError := &apierrors.APIError{}
	_ = json.Unmarshal([]byte(res), apiError)

	return apiError.ErrorCode
}

func GetOwnerRelations(r *http.Request) []domain.DBRelation {
	ownerRelationsCtx := context.Get(r, "ownerRelations")
	var ownerRelations []domain.DBRelation
	if ownerRelationsCtx != nil {
		ownerRelations = ownerRelationsCtx.([]domain.DBRelation)
	}

	return ownerRelations
}

func FilterIfOwnerRelations(r *http.Request, filter *usecases.Filter) *usecases.Filter {
	ownerRelationsCtx := context.Get(r, "ownerRelations")
	if ownerRelationsCtx != nil {
		currentSession := context.Get(r, "currentSession").(domain.Session)

		idKey := "accountId"

		if context.Get(r, "resource").(string) == "accounts" {
			idKey = "id"
		}

		if filter == nil {
			filter = &usecases.Filter{
				Where: map[string]interface{}{idKey: currentSession.AccountID},
			}
		} else {
			if filter.Where == nil {
				filter.Where = map[string]interface{}{idKey: currentSession.AccountID}
			} else {
				filter.Where[idKey] = currentSession.AccountID
			}
		}
	}

	return filter
}

func GetLastResource(r *http.Request) *Resource {
	lastResourceCtx := context.Get(r, "lastResource")
	var lastResource *Resource
	if lastResourceCtx != nil {
		lastResource = lastResourceCtx.(*Resource)
	} else {
		lastResource = &Resource{}
	}

	return lastResource
}

func FilterIfLastResource(r *http.Request, filter *usecases.Filter) *usecases.Filter {
	lastResourceCtx := context.Get(r, "lastResource")

	if lastResourceCtx != nil {
		lastResource := lastResourceCtx.(*Resource)

		if filter == nil {
			filter = &usecases.Filter{
				Where: map[string]interface{}{lastResource.IDKey: lastResource.ID},
			}
		} else {
			if filter.Where == nil {
				filter.Where = map[string]interface{}{lastResource.IDKey: lastResource.ID}
			} else {
				filter.Where[lastResource.IDKey] = lastResource.ID
			}
		}
	}

	return filter
}
