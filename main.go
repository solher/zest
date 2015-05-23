package main

import (
	"os"
	"runtime"

	"github.com/Solher/auth-scaffold/infrastructure"
	"github.com/Solher/auth-scaffold/interfaces"
	"github.com/Solher/auth-scaffold/middlewares"
	"github.com/Solher/auth-scaffold/ressources"
	"github.com/codegangsta/negroni"
	"github.com/dimfeld/httptreemux"

	_ "github.com/clipperhouse/typewriter" // Forced to permit vendoring.
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
	render := infrastructure.NewRender()
	store := infrastructure.NewGormStore()

	initApp(app, router, render, store)
	defer closeApp(store)

	app.Run(Port)
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

func connectDB(store *infrastructure.GormStore) error {
	var err error

	if DatabaseURL != "" {
		err = store.Connect("postgres", DatabaseURL)
	} else {
		err = store.Connect("sqlite3", "database.db")
	}

	return err
}

func initApp(app *negroni.Negroni, router *httptreemux.TreeMux, render *infrastructure.Render, store *infrastructure.GormStore) {
	err := connectDB(store)
	if err != nil {
		panic("Could not connect to database.")
	}

	routes := interfaces.NewRouteDirectory(render)

	userRepository := ressources.NewUserRepo(store)
	userInteractor := ressources.NewUserInter(userRepository)
	ressources.NewUserCtrl(userInteractor, render, routes)

	sessionRepository := ressources.NewSessionRepo(store)
	sessionInteractor := ressources.NewSessionInter(sessionRepository)
	ressources.NewSessionCtrl(sessionInteractor, render, routes)

	accountRepository := ressources.NewAccountRepo(store)
	accountInteractor := ressources.NewAccountInter(accountRepository, userRepository, sessionRepository)
	ressources.NewAccountCtrl(accountInteractor, render, routes)

	routes.Register(router)

	app.Use(negroni.NewLogger())
	app.Use(negroni.NewRecovery())
	app.Use(middlewares.NewSessions(accountRepository, userRepository, sessionRepository))

	app.UseHandler(router)
}

func closeApp(store *infrastructure.GormStore) {
	store.Close()
}
