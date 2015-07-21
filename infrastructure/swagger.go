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
	ExternalAPIPackage, SwaggerDir string
}

func NewSwagger() *Swagger {
	return &Swagger{SwaggerDir: "./swagger/"}
}

func (s *Swagger) Init(externalAPIPackage string) {
	s.ExternalAPIPackage = externalAPIPackage
}

func (s *Swagger) Generate() error {
	if s.ExternalAPIPackage == "" {
		return errors.New("You must specify an api package for the swagger doc to be generated")
	}

	params := generator.Params{
		ApiPackage:      "github.com/solher/zest/resources," + s.ExternalAPIPackage + "/resources",
		MainApiFile:     s.ExternalAPIPackage + "/main.go",
		OutputFormat:    "swagger",
		ControllerClass: "(Ctrl)$",
		OutputSpec:      "swagger/resources",
	}

	err := generator.Run(params)
	if err != nil {
		return err
	}

	return nil
}

func (s *Swagger) AddRoutes(routeDir *usecases.RouteDirectory) {
	dirKey := usecases.NewDirectoryKey("swagger")

	routeDir.Add(dirKey.For("UI"), &usecases.Route{Method: "GET", Path: "/explorer", Handler: s.UIHandler, Visible: true, CheckPermissions: false})
	routeDir.Add(dirKey.For("Resources"), &usecases.Route{Method: "GET", Path: "/explorer/*path", Handler: s.ResourcesHandler, Visible: true, CheckPermissions: false})
}

func (s *Swagger) UIHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	relativePath, _ := getRelativePath(r.URL.Path)

	if relativePath == "" && !strings.HasSuffix(r.URL.Path, "/") {
		w.Header().Set("Location", relativePath+"explorer/")
		w.WriteHeader(http.StatusMovedPermanently)
		return
	}

	http.ServeFile(w, r, s.SwaggerDir+relativePath)
}

func (s *Swagger) ResourcesHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var pathPrefix string
	r.URL.Path, pathPrefix = getRelativePath(r.URL.Path)

	if splittedPath := strings.Split(r.URL.Path, "/"); !strings.Contains(splittedPath[len(splittedPath)-1], ".") {
		if !strings.HasSuffix(r.URL.Path, "/") {
			r.URL.Path = r.URL.Path + "/"
		}

		r.URL.Path = r.URL.Path + "index.json"

		t, e := template.ParseFiles(s.SwaggerDir + r.URL.Path)
		if e != nil {
			panic(e)
		} else {
			t.Execute(w, "http://"+r.Host+pathPrefix)
		}
	} else {
		http.ServeFile(w, r, s.SwaggerDir+r.URL.Path)
	}
}

func getRelativePath(url string) (string, string) {
	splittedPath := strings.Split(url, "/")
	var relativePath, prefix string

	for i, path := range splittedPath {
		if path == "explorer" {
			relativePath = strings.Join(splittedPath[i+1:], "/")
			prefix = strings.Join(splittedPath[:i], "/")
			break
		}
	}

	return relativePath, prefix
}
