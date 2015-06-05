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

	userRepo := ressources.NewUserRepo(store)
	sessionRepo := ressources.NewSessionRepo(store)
	accountRepo := ressources.NewAccountRepo(store)
	roleMappingRepo := ressources.NewRoleMappingRepo(store)
	roleRepo := ressources.NewRoleRepo(store)
	aclMappingRepo := ressources.NewAclMappingRepo(store)
	aclRepo := ressources.NewAclRepo(store)

	sessionCacheInter := usecases.NewSessionCacheInter(sessionRepo, lruCacheStore)
	permissionCacheInter := usecases.NewPermissionCacheInter(accountRepo, aclRepo, roleCacheStore, aclCacheStore)
	sessionCacheInter.Refresh()
	permissionCacheInter.Refresh()

	userInter := ressources.NewUserInter(userRepo)
	sessionInter := ressources.NewSessionInter(sessionRepo)
	accountInter := ressources.NewAccountInter(accountRepo, userInter, sessionInter, sessionCacheInter, permissionCacheInter)
	roleMappingInter := ressources.NewRoleMappingInter(roleMappingRepo)
	roleInter := ressources.NewRoleInter(roleRepo)
	aclMappingInter := ressources.NewAclMappingInter(aclMappingRepo)
	aclInter := ressources.NewAclInter(aclRepo)

	routeDir := usecases.NewRouteDirectory(accountInter, render)

	ressources.NewUserCtrl(userInter, render, routeDir)
	ressources.NewSessionCtrl(sessionInter, render, routeDir)
	ressources.NewAccountCtrl(accountInter, render, routeDir)
	ressources.NewRoleMappingCtrl(roleMappingInter, render, routeDir)
	ressources.NewRoleCtrl(roleInter, render, routeDir)
	ressources.NewAclMappingCtrl(aclMappingInter, render, routeDir)
	ressources.NewAclCtrl(aclInter, render, routeDir)

	routeDir.Register(router)
	aclInter.RefreshFromRoutes(routeDir.Routes())

	app.Use(negroni.NewLogger())
	app.Use(middlewares.NewRecovery(render))
	app.Use(cors.Default())
	app.Use(middlewares.NewSessions(accountInter))

	app.UseHandler(router)

	return routeDir.Routes()
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
