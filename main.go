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
	mustExit := handleOsArgs()
	if mustExit {
		return
	}

	app := negroni.New()
	router := httptreemux.New()
	router.RedirectBehavior = httptreemux.UseHandler
	render := infrastructure.NewRender()
	store := infrastructure.NewGormStore()
	sessionCache := infrastructure.NewLRUCacheStore(1024)
	roleCache := infrastructure.NewCacheStore()
	aclCache := infrastructure.NewCacheStore()

	initApp(app, router, render, store, sessionCache, roleCache, aclCache)
	defer closeApp(store)

	app.Run(Port)
}

func handleOsArgs() bool {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "resetDB":
			resetDatabase()
			return true
		case "migrateDB":
			migrateDatabase()
			return true
		}
	}

	return false
}

func initApp(app *negroni.Negroni, router *httptreemux.TreeMux, render *infrastructure.Render,
	store *infrastructure.GormStore, lruCacheStore *infrastructure.LRUCacheStore, roleCacheStore, aclCacheStore *infrastructure.CacheStore) {

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

	routes := usecases.NewRouteDirectory(accountInter, render)

	ressources.NewUserCtrl(userInter, render, routes)
	ressources.NewSessionCtrl(sessionInter, render, routes)
	ressources.NewAccountCtrl(accountInter, render, routes)
	ressources.NewRoleMappingCtrl(roleMappingInter, render, routes)
	ressources.NewRoleCtrl(roleInter, render, routes)
	ressources.NewAclMappingCtrl(aclMappingInter, render, routes)
	ressources.NewAclCtrl(aclInter, render, routes)

	routes.Register(router)

	app.Use(negroni.NewLogger())
	app.Use(negroni.NewRecovery())
	app.Use(middlewares.NewSessions(accountInter))

	app.UseHandler(router)
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
