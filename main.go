package main

import (
	"os"
	"runtime"

	"github.com/Solher/zest/infrastructure"
	"github.com/Solher/zest/interfaces"
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

	userRepository := ressources.NewUserRepo(store)
	sessionRepository := ressources.NewSessionRepo(store)
	accountRepository := ressources.NewAccountRepo(store)
	roleMappingRepository := ressources.NewRoleMappingRepo(store)
	roleRepository := ressources.NewRoleRepo(store)
	aclMappingRepository := ressources.NewAclMappingRepo(store)
	aclRepository := ressources.NewAclRepo(store)

	sessionCacheInter := usecases.NewSessionCacheInter(sessionRepository, lruCacheStore)
	permissionCacheInter := usecases.NewPermissionCacheInter(accountRepository, aclRepository, roleCacheStore, aclCacheStore)
	sessionCacheInter.Refresh()
	permissionCacheInter.Refresh()

	userInteractor := ressources.NewUserInter(userRepository)
	sessionInteractor := ressources.NewSessionInter(sessionRepository)
	accountInteractor := ressources.NewAccountInter(accountRepository, userInteractor, sessionInteractor, sessionCacheInter, permissionCacheInter)
	roleMappingInteractor := ressources.NewRoleMappingInter(roleMappingRepository)
	roleInteractor := ressources.NewRoleInter(roleRepository)
	aclMappingInteractor := ressources.NewAclMappingInter(aclMappingRepository)
	aclInteractor := ressources.NewAclInter(aclRepository)

	routes := interfaces.NewRouteDirectory(accountInteractor, render)

	ressources.NewUserCtrl(userInteractor, render, routes)
	ressources.NewSessionCtrl(sessionInteractor, render, routes)
	ressources.NewAccountCtrl(accountInteractor, render, routes)
	ressources.NewRoleMappingCtrl(roleMappingInteractor, render, routes)
	ressources.NewRoleCtrl(roleInteractor, render, routes)
	ressources.NewAclMappingCtrl(aclMappingInteractor, render, routes)
	ressources.NewAclCtrl(aclInteractor, render, routes)

	routes.Register(router)

	app.Use(negroni.NewLogger())
	app.Use(negroni.NewRecovery())
	app.Use(cors.Default())
	app.Use(middlewares.NewSessions(accountInteractor))

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
