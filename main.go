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
	router.RedirectBehavior = httptreemux.UseHandler
	render := infrastructure.NewRender()
	store := infrastructure.NewGormStore()
	sessionCache := infrastructure.NewLRUCacheStore(1024)
	roleCache := infrastructure.NewCacheStore()
	aclCache := infrastructure.NewCacheStore()

	routes := initApp(app, router, render, store, sessionCache, roleCache, aclCache)
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

func initApp(app *negroni.Negroni, router *httptreemux.TreeMux, render *infrastructure.Render,
	store *infrastructure.GormStore, lruCacheStore *infrastructure.LRUCacheStore,
	roleCacheStore, aclCacheStore *infrastructure.CacheStore) map[usecases.DirectoryKey]usecases.Route {

	err := connectDB(store)
	if err != nil {
		panic("Could not connect to database.")
	}

	dependencies := []interface{}{
		store,
		render,
		lruCacheStore,
		roleCacheStore,
		aclCacheStore,

		usecases.NewRouteDirectory,
		usecases.NewSessionCacheInter,
		usecases.NewPermissionCacheInter,

		ressources.NewUserRepo,
		ressources.NewUserInter,
		ressources.NewUserCtrl,

		ressources.NewSessionRepo,
		ressources.NewSessionInter,
		ressources.NewSessionCtrl,

		ressources.NewAccountRepo,
		ressources.NewAccountInter,
		ressources.NewAccountCtrl,

		ressources.NewRoleMappingRepo,
		ressources.NewRoleMappingInter,
		ressources.NewRoleMappingCtrl,

		ressources.NewRoleRepo,
		ressources.NewRoleInter,
		ressources.NewRoleCtrl,

		ressources.NewAclMappingRepo,
		ressources.NewAclMappingInter,
		ressources.NewAclMappingCtrl,

		ressources.NewAclRepo,
		ressources.NewAclInter,
		ressources.NewAclCtrl,
	}

	injector := infrastructure.NewInjector()
	injector.RegisterMultiple(dependencies)
	err = injector.Populate()

	if err != nil {
		panic(err)
	} else {
		type dependencies struct {
			SessionCacheInter    *usecases.SessionCacheInter
			PermissionCacheInter *usecases.PermissionCacheInter
			AclInter             *ressources.AclInter
			AccountInter         *ressources.AccountInter

			RouteDir *usecases.RouteDirectory
			Render   *infrastructure.Render
		}

		deps := &dependencies{}
		err = injector.Get(deps)
		if err != nil {
			panic(err)
		}

		deps.SessionCacheInter.Refresh()
		deps.PermissionCacheInter.Refresh()

		deps.RouteDir.Register(router)
		deps.AclInter.RefreshFromRoutes(deps.RouteDir.Routes())

		app.Use(negroni.NewLogger())
		app.Use(middlewares.NewRecovery(deps.Render))
		app.Use(cors.Default())
		app.Use(middlewares.NewSessions(deps.AccountInter))

		app.UseHandler(router)

		return deps.RouteDir.Routes()
	}

	return nil
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
