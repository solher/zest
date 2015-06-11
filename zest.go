package zest

import (
	"os"

	_ "github.com/clipperhouse/typewriter"
	"github.com/codegangsta/negroni"
	"github.com/dimfeld/httptreemux"
	"github.com/rs/cors"
	"github.com/solher/zest/infrastructure"
	"github.com/solher/zest/middlewares"
	"github.com/solher/zest/ressources"
	"github.com/solher/zest/usecases"
)

type Zest struct {
	Port, Environment, DatabaseURL string

	app      *negroni.Negroni
	injector *infrastructure.Injector

	HandleOsArgs   func(z *Zest) (bool, error)
	BuildApp       func(z *Zest) error
	InitApp        func(z *Zest) error
	CloseApp       func(z *Zest) error
	ReinitDatabase func(z *Zest) error
	SeedDatabase   func(z *Zest) error
	UpdateDatabase func(z *Zest) error
}

func New() *Zest {
	return &Zest{
		Environment: "development",
		Port:        ":3000",
		app:         negroni.New(),
		injector:    infrastructure.NewInjector(),
	}
}

func Classic() *Zest {
	zest := &Zest{
		Environment:    "development",
		Port:           ":3000",
		app:            negroni.New(),
		injector:       infrastructure.NewInjector(),
		HandleOsArgs:   handleOsArgs,
		BuildApp:       buildApp,
		InitApp:        initApp,
		CloseApp:       closeApp,
		ReinitDatabase: reinitDatabase,
		SeedDatabase:   seedDatabase,
		UpdateDatabase: updateDatabase,
	}

	env := os.Getenv("GOENV")
	if env == "development" || env == "production" {
		zest.Environment = env
	}

	port := os.Getenv("PORT")
	if port != "" {
		zest.Port = ":" + os.Getenv("PORT")
	}

	zest.DatabaseURL = os.Getenv("DATABASE_URL")

	return zest
}

func (z *Zest) Run() {
	err := z.BuildApp(z)
	if err != nil {
		panic(err)
	}

	mustExit, err := z.HandleOsArgs(z)
	if err != nil {
		panic(err)
	}

	if mustExit {
		return
	}

	err = z.InitApp(z)
	if err != nil {
		panic(err)
	}
	defer z.CloseApp(z)

	z.app.Run(z.Port)
}

func handleOsArgs(z *Zest) (bool, error) {
	var err error

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "resetDB":
			err = z.ReinitDatabase(z)
			if err != nil {
				return true, err
			}
			err = z.SeedDatabase(z)
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

	return false, nil
}

func buildApp(z *Zest) error {
	deps := usecases.DependencyDirectory.Get()

	store := infrastructure.NewGormStore()
	accountRepo := ressources.NewAccountRepo(store)
	aclRepo := ressources.NewAclRepo(store)
	sessionRepo := ressources.NewSessionRepo(store)

	deps = append(
		deps,

		store,
		accountRepo,
		aclRepo,
		sessionRepo,
		httptreemux.New(),
		infrastructure.NewRender(),
		usecases.NewPermissionCacheInter(accountRepo, aclRepo, infrastructure.NewCacheStore(), infrastructure.NewCacheStore()),
		usecases.NewSessionCacheInter(sessionRepo, infrastructure.NewLRUCacheStore(1024)),
		usecases.NewRouteDirectory,
	)

	z.injector.RegisterMultiple(deps)

	err := z.injector.Populate()
	if err != nil {
		return err
	}

	return nil
}

func initApp(z *Zest) error {
	type dependencies struct {
		Router *httptreemux.TreeMux

		Store *infrastructure.GormStore

		SessionCacheInter    *usecases.SessionCacheInter
		PermissionCacheInter *usecases.PermissionCacheInter
		AccountInter         *ressources.AccountInter

		RouteDir *usecases.RouteDirectory
		Render   *infrastructure.Render
	}

	d := &dependencies{}
	err := z.injector.Get(d)
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

	z.app.Use(middlewares.NewLogger())
	z.app.Use(middlewares.NewRecovery(d.Render))
	z.app.Use(cors.Default())
	z.app.Use(middlewares.NewSessions(d.AccountInter))

	d.Router.RedirectBehavior = httptreemux.UseHandler

	z.app.UseHandler(d.Router)

	return nil
}

func closeApp(z *Zest) error {
	type dependencies struct {
		Store *infrastructure.GormStore
	}

	d := &dependencies{}
	err := z.injector.Get(d)
	if err != nil {
		return err
	}

	d.Store.Close()

	return nil
}
