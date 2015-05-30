package interfaces

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/Solher/auth-scaffold/apierrors"
	"github.com/Solher/auth-scaffold/domain"
	"github.com/gorilla/context"
)

func GetQueryFilter(r *http.Request) (*Filter, error) {
	param := r.URL.Query().Get("filter")
	if param == "" {
		return nil, nil
	}

	filterReader := strings.NewReader(param)

	filter := &Filter{}
	err := json.NewDecoder(filterReader).Decode(filter)
	if err != nil {
		return nil, err
	}

	return filter, nil
}

func MockHTTPRequest(route Route, body, filter string, params map[string]string) string {
	if route.Method == "" || route.Path == "" || route.Handler == nil {
		panic("Non existing or incomplete route when mocking HTTP request.")
	}

	w := httptest.NewRecorder()

	path := route.Path
	if filter != "" {
		path = path + "?filter=" + filter
	}

	req, _ := http.NewRequest(route.Method, path, bytes.NewBufferString(body))
	(*route.Handler)(w, req, params)

	return w.Body.String()
}

func GetErrorCode(res string) string {
	apiError := &apierrors.APIError{}
	_ = json.Unmarshal([]byte(res), apiError)

	return apiError.ErrorCode
}

func GetOwnerRelations(r *http.Request) []domain.Relation {
	ownerRelationsCtx := context.Get(r, "ownerRelations")
	var ownerRelations []domain.Relation
	if ownerRelationsCtx != nil {
		ownerRelations = ownerRelationsCtx.([]domain.Relation)
	}

	return ownerRelations
}

func FilterIfOwnerRelations(r *http.Request, filter *Filter) *Filter {
	ownerRelationsCtx := context.Get(r, "ownerRelations")
	if ownerRelationsCtx != nil {
		currentSession := context.Get(r, "currentSession").(domain.Session)

		if filter == nil {
			filter = &Filter{
				where: map[string]interface{}{"accountId": currentSession.AccountID},
			}
		} else {
			if filter.where == nil {
				filter.where = map[string]interface{}{"accountId": currentSession.AccountID}
			} else {
				filter.where["accountId"] = currentSession.AccountID
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

func FilterIfLastRessource(r *http.Request, filter *Filter) *Filter {
	lastRessourceCtx := context.Get(r, "lastRessource")

	if lastRessourceCtx != nil {
		lastRessource := lastRessourceCtx.(*Ressource)

		if filter == nil {
			filter = &Filter{
				where: map[string]interface{}{lastRessource.IDKey: lastRessource.ID},
			}
		} else {
			if filter.where == nil {
				filter.where = map[string]interface{}{lastRessource.IDKey: lastRessource.ID}
			} else {
				filter.where[lastRessource.IDKey] = lastRessource.ID
			}
		}
	}

	return filter
}
