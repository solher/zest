package interfaces

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/Solher/zest/apierrors"
	"github.com/Solher/zest/domain"
	"github.com/Solher/zest/usecases"
	"github.com/gorilla/context"
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

		if filter == nil {
			filter = &usecases.Filter{
				Where: map[string]interface{}{"accountId": currentSession.AccountID},
			}
		} else {
			if filter.Where == nil {
				filter.Where = map[string]interface{}{"accountId": currentSession.AccountID}
			} else {
				filter.Where["accountId"] = currentSession.AccountID
			}
		}
	}

	return filter
}

func GetLastRessource(r *http.Request) *Ressource {
	lastRessourceCtx := context.Get(r, "lastRessource")
	var lastRessource *Ressource
	if lastRessourceCtx != nil {
		lastRessource = lastRessourceCtx.(*Ressource)
	} else {
		lastRessource = &Ressource{}
	}

	return lastRessource
}

func FilterIfLastRessource(r *http.Request, filter *usecases.Filter) *usecases.Filter {
	lastRessourceCtx := context.Get(r, "lastRessource")

	if lastRessourceCtx != nil {
		lastRessource := lastRessourceCtx.(*Ressource)

		if filter == nil {
			filter = &usecases.Filter{
				Where: map[string]interface{}{lastRessource.IDKey: lastRessource.ID},
			}
		} else {
			if filter.Where == nil {
				filter.Where = map[string]interface{}{lastRessource.IDKey: lastRessource.ID}
			} else {
				filter.Where[lastRessource.IDKey] = lastRessource.ID
			}
		}
	}

	return filter
}
