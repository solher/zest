package main

import (
	"os"
	"runtime"

	"github.com/Solher/auth-scaffold/infrastructure"
	"github.com/Solher/auth-scaffold/interfaces"
	"github.com/Solher/auth-scaffold/middlewares"
	"github.com/Solher/auth-scaffold/ressources/sessions"
	"github.com/Solher/auth-scaffold/ressources/users"
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

	router := httprouter.New()
	render := infrastructure.NewRender()
	store := infrastructure.NewGormStore()

	negroni := negroni.New(negroni.NewRecovery(), negroni.NewLogger(), middlewares.NewSessions(store))
	negroni.UseHandler(router)

	initApp(router, render, store)
	defer closeApp(store)

	negroni.Run(":3001")
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

func initApp(router *httprouter.Router, render *infrastructure.Render, store *infrastructure.GormStore) {
	err := store.Connect("sqlite3", "database.db")
	if err != nil {
		panic("Could not connect to database.")
	}
	routes := interfaces.NewRouteDirectory()

	usersRepository := users.NewRepository(store)
	usersInteractor := users.NewInteractor(usersRepository)
	users.NewController(usersInteractor, render, routes)

	sessionsRepository := sessions.NewRepository(store)
	sessionsInteractor := sessions.NewInteractor(sessionsRepository)
	sessions.NewController(sessionsInteractor, render, routes)

	routes.Register(router)
}

func closeApp(store *infrastructure.GormStore) {
	store.Close()
}
