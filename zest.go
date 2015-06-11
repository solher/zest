package zest

import (
	"os"

	"github.com/Solher/zest/infrastructure"
	"github.com/Solher/zest/middlewares"
	"github.com/Solher/zest/ressources"
	"github.com/Solher/zest/usecases"
	"github.com/codegangsta/negroni"
	"github.com/dimfeld/httptreemux"
	"github.com/rs/cors"
)

type Zest struct {
	Port, Environment, DatabaseURL string

	injector *infrastructure.Injector
	app      *negroni.Negroni

	InitApp        func(z *Zest) error
	CloseApp       func(z *Zest) error
	ResetDatabase  func(z *Zest) error
	SeedDatabase   func(z *Zest) error
	UpdateDatabase func(z *Zest) error
}

func Classic() *Zest {
	zest := &Zest{app: negroni.New()}

	env := os.Getenv("GOENV")
	if env == "development" || env == "production" {
		zest.Environment = env
	} else {
		zest.Environment = "development"
	}

	port := os.Getenv("PORT")
	if port != "" {
		zest.Port = ":" + os.Getenv("PORT")
	} else {
		zest.Port = ":3000"
	}

	zest.DatabaseURL = os.Getenv("DATABASE_URL")

	return zest
}

func (z *Zest) Run() {
	err := z.InitApp(z)
	if err != nil {
		panic(err)
	}

	defer z.CloseApp(z)

	mustExit, err := z.handleOsArgs()
	if err != nil {
		panic(err)
	}

	if mustExit {
		return
	}

	z.app.Run(z.Port)
}

func (z *Zest) handleOsArgs() (bool, error) {
	var err error

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "resetDB":
			err = z.ResetDatabase(z)
			if err != nil {
				return true, err
			}
			return true, nil
		case "updateDB":
			err = z.UpdateDatabase(z)
			if err != nil {
				return true, err
			}
			return true, nil
		}
	}

	return true, nil
}

func initApp(z *Zest) error {
	deps := usecases.DependencyDirectory.Get()

	deps = append(
		deps,
		httptreemux.New(),
		infrastructure.NewGormStore(),
		infrastructure.NewRender(),
		infrastructure.NewLRUCacheStore(1024),
		infrastructure.NewCacheStore(),
		infrastructure.NewCacheStore(),

		usecases.NewRouteDirectory,
		usecases.NewPermissionCacheInter,
		usecases.NewSessionCacheInter,
	)

	z.injector.RegisterMultiple(deps)

	err := z.injector.Populate()
	if err != nil {
		return err
	}

	type dependencies struct {
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

	d := &dependencies{}
	err = z.injector.Get(d)
	if err != nil {
		return err
	}

	if z.DatabaseURL != "" {
		err = d.Store.Connect("postgres", z.DatabaseURL)
	} else {
		err = d.Store.Connect("sqlite3", "database.db")
	}
	if err != nil {
		return err
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

	return nil
}

func closeApp(z *Zest) {
	type dependencies struct {
		Store *infrastructure.GormStore
	}

	d := &dependencies{}
	err := z.injector.Get(d)
	if err != nil {
		panic(err)
	}

	d.Store.Close()
}
