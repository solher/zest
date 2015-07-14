package infrastructure

import (
	"errors"
	"net/http"
	"strings"
	"text/template"

	"github.com/yvasiyarov/swagger/generator"

	"github.com/solher/zest/usecases"
)

type Swagger struct {
	ResourceListingJSON, ExternalAPIPackage string
	APIDescriptionsJSON                     map[string]string
}

func NewSwagger() *Swagger {
	return &Swagger{}
}

func (s *Swagger) Init(apiDescriptionsJSON map[string]string, resourceListingJSON, externalAPIPackage string) {
	s.APIDescriptionsJSON = apiDescriptionsJSON
	s.ResourceListingJSON = resourceListingJSON
	s.ExternalAPIPackage = externalAPIPackage
}

func (s *Swagger) Generate() error {
	if s.ExternalAPIPackage == "" {
		return errors.New("You must specify an api package for the swagger doc to be generated")
	}

	params := generator.Params{
		ApiPackage:      "github.com/solher/zest," + s.ExternalAPIPackage,
		MainApiFile:     s.ExternalAPIPackage + "/main.go",
		OutputFormat:    "go",
		ControllerClass: "(Ctrl)$",
	}

	err := generator.Run(params)
	if err != nil {
		return err
	}

	return nil
}

func (s *Swagger) AddRoutes(routeDir *usecases.RouteDirectory) {
	if s.ResourceListingJSON == "" && s.APIDescriptionsJSON == nil {
		return
	}

	fileServer := http.FileServer(http.Dir("./swagger-ui"))

	handlerFunc := func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		splittedPath := strings.Split(r.URL.Path, "/")
		var shortPath string

		for i, path := range splittedPath {
			if path == "explorer" {
				shortPath = strings.Join(splittedPath[i+1:], "/")
				break
			}
		}

		if shortPath == "" && !strings.HasSuffix(r.URL.Path, "/") {
			w.Header().Set("Location", shortPath+"explorer/")
			w.WriteHeader(http.StatusMovedPermanently)
			return
		}

		r.URL.Path = shortPath

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
		t, e := template.New("desc").Parse(json)
		if e != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		t.Execute(w, "http://localhost:3005/api")
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
