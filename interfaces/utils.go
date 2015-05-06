package interfaces

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/Solher/auth-scaffold/apierrors"
	"github.com/julienschmidt/httprouter"
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

func MockHTTPRequest(route Route, body, filter string, params httprouter.Params) string {
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
