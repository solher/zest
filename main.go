package main

import (
	"os"
	"runtime"

	"github.com/Solher/auth-scaffold/infrastructure"
	"github.com/Solher/auth-scaffold/interfaces"
	"github.com/Solher/auth-scaffold/middlewares"
	"github.com/Solher/auth-scaffold/ressources"
	"github.com/Solher/auth-scaffold/usecases"
	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	mustExit := handleOsArgs()
	if mustExit {
		return
	}

	app := negroni.New()
	router := httprouter.New()
	render := infrastructure.NewRender()
	store := infrastructure.NewGormStore()

	initApp(app, router, render, store)
	defer closeApp(store)

	app.Run(":3001")
}

func handleOsArgs() bool {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "reinitDB":
			reinitDatabase()
			return true
		case "migrateDB":
			migrateDatabase()
			return true
		}
	}

	return false
}

func initApp(app *negroni.Negroni, router *httprouter.Router, render *infrastructure.Render, store *infrastructure.GormStore) {
	err := store.Connect("sqlite3", "database.db")
	if err != nil {
		panic("Could not connect to database.")
	}
	routes := interfaces.NewRouteDirectory()
	permissions := usecases.NewPermissionDirectory()

	userRepository := ressources.NewUserRepo(store)
	userInteractor := ressources.NewUserInter(userRepository)
	ressources.NewUserCtrl(userInteractor, render, routes, permissions)

	sessionRepository := ressources.NewSessionRepo(store)
	sessionInteractor := ressources.NewSessionInter(sessionRepository)
	ressources.NewSessionCtrl(sessionInteractor, render, routes, permissions)

	accountRepository := ressources.NewAccountRepo(store)
	accountInteractor := ressources.NewAccountInter(accountRepository, userRepository, sessionRepository)
	ressources.NewAccountCtrl(accountInteractor, render, routes, permissions)

	routes.Register(router, permissions, render)

	app.Use(negroni.NewLogger())
	app.Use(negroni.NewRecovery())
	app.Use(middlewares.NewSessions(accountRepository, userRepository, sessionRepository))

	app.UseHandler(router)
}

func closeApp(store *infrastructure.GormStore) {
	store.Close()
}
