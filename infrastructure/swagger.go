package infrastructure

import (
	"errors"
	"net/http"
	"strings"
	"text/template"

	"github.com/solher/swagger/generator"

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
		ApiPackage:      "github.com/solher/zest," + s.ExternalAPIPackage,
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
	routeDir.Add(dirKey.For("Favicon"), &usecases.Route{Method: "GET", Path: "/favicon.ico", Handler: s.FaviconHandler, Visible: true, CheckPermissions: false})
}

func (s *Swagger) FaviconHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	http.ServeFile(w, r, s.SwaggerDir+"images/favicon.ico")
}

func (s *Swagger) UIHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	shortPath := getShortPath(r.URL.Path)

	if shortPath == "" && !strings.HasSuffix(r.URL.Path, "/") {
		w.Header().Set("Location", shortPath+"explorer/")
		w.WriteHeader(http.StatusMovedPermanently)
		return
	}

	http.ServeFile(w, r, s.SwaggerDir+shortPath)
}

func (s *Swagger) ResourcesHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	r.URL.Path = getShortPath(r.URL.Path)

	if splittedPath := strings.Split(r.URL.Path, "/"); !strings.Contains(splittedPath[len(splittedPath)-1], ".") {
		if !strings.HasSuffix(r.URL.Path, "/") {
			r.URL.Path = r.URL.Path + "/"
		}

		r.URL.Path = r.URL.Path + "index.json"

		t, e := template.ParseFiles(s.SwaggerDir + r.URL.Path)
		if e != nil {
			panic(e)
		} else {
			t.Execute(w, "http://"+r.Host)
		}
	} else {
		http.ServeFile(w, r, s.SwaggerDir+r.URL.Path)
	}
}

func getShortPath(url string) string {
	splittedPath := strings.Split(url, "/")
	shortPath := ""

	for i, path := range splittedPath {
		if path == "explorer" {
			shortPath = strings.Join(splittedPath[i+1:], "/")
			break
		}
	}

	return shortPath
}
