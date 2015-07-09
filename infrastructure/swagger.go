package infrastructure

import (
	"net/http"
	"strings"

	"github.com/solher/zest/usecases"
)

type Swagger struct {
	ResourceListingJSON string
	APIDescriptionsJSON map[string]string
}

func NewSwagger() *Swagger {
	return &Swagger{}
}

func (s *Swagger) AddRoutes(routeDir *usecases.RouteDirectory) {
	if s.ResourceListingJSON == "" && s.APIDescriptionsJSON == nil {
		return
	}

	fileServer := http.FileServer(http.Dir("./swagger-ui"))

	handlerFunc := func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		splittedPath := strings.Split(r.URL.Path, "/")

		for i, path := range splittedPath {
			if path == "explorer" {
				r.URL.Path = strings.Join(splittedPath[i+1:], "/")
				break
			}
		}

		fileServer.ServeHTTP(w, r)
	}

	dirKey := usecases.NewDirectoryKey("swagger")

	routeDir.Add(dirKey.For("APIResources"), &usecases.Route{Method: "GET", Path: "/explorer/resources", Handler: s.IndexHandler, Visible: true, CheckPermissions: false})
	routeDir.Add(dirKey.For("UI"), &usecases.Route{Method: "GET", Path: "/explorer", Handler: handlerFunc, Visible: true, CheckPermissions: false})
	routeDir.Add(dirKey.For("UIResources"), &usecases.Route{Method: "GET", Path: "/explorer/*path", Handler: handlerFunc, Visible: true, CheckPermissions: false})

	for apiKey := range s.APIDescriptionsJSON {
		routeDir.Add(dirKey.For(apiKey), &usecases.Route{Method: "GET", Path: "/explorer/resources/" + apiKey, Handler: s.APIDescriptionHandler, Visible: true, CheckPermissions: false})
	}
}

func (s *Swagger) APIDescriptionHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	arrayURI := strings.Split(r.RequestURI, "/")
	apiKey := arrayURI[len(arrayURI)-1]

	if json, ok := s.APIDescriptionsJSON[apiKey]; ok {
		w.Write([]byte(json))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *Swagger) IndexHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	isJsonRequest := false

	if acceptHeaders, ok := r.Header["Accept"]; ok {
		for _, acceptHeader := range acceptHeaders {
			if strings.Contains(acceptHeader, "json") {
				isJsonRequest = true
				break
			}
		}
	}

	if isJsonRequest {
		w.Write([]byte(s.ResourceListingJSON))
	}
}
