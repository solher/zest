package main

import (
	"os"
	"runtime"

	"github.com/Solher/zest/infrastructure"
	"github.com/Solher/zest/middlewares"
	"github.com/Solher/zest/ressources"
	"github.com/Solher/zest/usecases"
	"github.com/codegangsta/negroni"
	"github.com/dimfeld/httptreemux"
	"github.com/rs/cors"

	_ "github.com/clipperhouse/typewriter" // Forced to allow vendoring.
)

var Environment, Port, DatabaseURL string

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	env := os.Getenv("GOENV")
	if env == "development" || env == "production" {
		Environment = env
	} else {
		Environment = "development"
	}

	port := os.Getenv("PORT")
	if port != "" {
		Port = ":" + os.Getenv("PORT")
	} else {
		Port = ":3000"
	}

	DatabaseURL = os.Getenv("DATABASE_URL")
}

func main() {
	app := negroni.New()
	router := httptreemux.New()
	store := infrastructure.NewGormStore()

	dependencies := []interface{}{
		app,
		router,
		store,
		infrastructure.NewRender(),
		infrastructure.NewLRUCacheStore(1024),
		infrastructure.NewCacheStore(),
		infrastructure.NewCacheStore(),
	}

	routes := initApp(dependencies)
	defer closeApp(store)

	mustExit := handleOsArgs(routes)
	if mustExit {
		return
	}

	app.Run(Port)
}

func handleOsArgs(routes map[usecases.DirectoryKey]usecases.Route) bool {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "resetDB":
			resetDatabase(routes)
			return true
		case "updateDB":
			updateDatabase(routes)
			return true
		}
	}

	return false
}

func initApp(dependencies []interface{}) usecases.Routes {
	depDir := usecases.DependencyDirectory
	depDir.RegisterMultiple(dependencies)
	depDir.Register(usecases.NewRouteDirectory)
	depDir.Register(usecases.NewPermissionCacheInter)
	depDir.Register(usecases.NewSessionCacheInter)

	err := depDir.Populate()
	if err != nil {
		panic(err)
	}

	type deps struct {
		App    *negroni.Negroni
		Router *httptreemux.TreeMux

		Store *infrastructure.GormStore

		SessionCacheInter    *usecases.SessionCacheInter
		PermissionCacheInter *usecases.PermissionCacheInter
		AclInter             *ressources.AclInter
		AccountInter         *ressources.AccountInter

		RouteDir *usecases.RouteDirectory
		Render   *infrastructure.Render
	}

	d := &deps{}
	err = depDir.Get(d)
	if err != nil {
		panic(err)
	}

	err = connectDB(d.Store)
	if err != nil {
		panic(err)
	}

	d.SessionCacheInter.Refresh()
	d.PermissionCacheInter.Refresh()

	d.RouteDir.Register(d.Router)
	d.AclInter.RefreshFromRoutes(d.RouteDir.Routes())

	d.App.Use(negroni.NewLogger())
	d.App.Use(middlewares.NewRecovery(d.Render))
	d.App.Use(cors.Default())
	d.App.Use(middlewares.NewSessions(d.AccountInter))

	d.Router.RedirectBehavior = httptreemux.UseHandler
	d.App.UseHandler(d.Router)

	return d.RouteDir.Routes()
}

func closeApp(store *infrastructure.GormStore) {
	store.Close()
}

func connectDB(store *infrastructure.GormStore) error {
	var err error

	if DatabaseURL != "" {
		err = store.Connect("postgres", DatabaseURL)
	} else {
		err = store.Connect("sqlite3", "database.db")
	}

	return err
}
