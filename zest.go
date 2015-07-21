package zest

import (
	"fmt"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/dimfeld/httptreemux"
	"github.com/rs/cors"
	"github.com/solher/zest/infrastructure"
	"github.com/solher/zest/middlewares"
	"github.com/solher/zest/resources"
	"github.com/solher/zest/usecases"
)

type Zest struct {
	Port, Environment, DatabaseURL string

	App      *negroni.Negroni
	Injector *infrastructure.Injector
	Swagger  *infrastructure.Swagger

	HandleOsArgs     func(z *Zest) (bool, error)
	Build            func(z *Zest) error
	AfterBuild       func(z *Zest) error
	Init             func(z *Zest) error
	AfterInit        func(z *Zest) error
	Close            func(z *Zest) error
	ReinitDatabase   func(z *Zest) error
	SeedDatabase     func(z *Zest) error
	UserSeedDatabase func(z *Zest) error
	UpdateDatabase   func(z *Zest) error
}

func New() *Zest {
	return &Zest{
		Environment: "development",
		Port:        "3000",
		App:         negroni.New(),
		Swagger:     infrastructure.NewSwagger(),
		Injector:    infrastructure.NewInjector(),
	}
}

func Classic() *Zest {
	zest := &Zest{
		Environment:      "development",
		Port:             "3000",
		App:              negroni.New(),
		Swagger:          infrastructure.NewSwagger(),
		Injector:         infrastructure.NewInjector(),
		HandleOsArgs:     handleOsArgs,
		Build:            buildApp,
		AfterBuild:       func(z *Zest) error { return nil },
		Init:             initApp,
		AfterInit:        func(z *Zest) error { return nil },
		Close:            closeApp,
		ReinitDatabase:   reinitDatabase,
		SeedDatabase:     seedDatabase,
		UserSeedDatabase: func(z *Zest) error { return nil },
		UpdateDatabase:   updateDatabase,
	}

	env := os.Getenv("GOENV")
	if env == "development" || env == "production" {
		zest.Environment = env
	}

	port := os.Getenv("PORT")
	if port != "" {
		zest.Port = os.Getenv("PORT")
	}

	zest.DatabaseURL = os.Getenv("DATABASE_URL")

	return zest
}

func (z *Zest) Run() {
	err := z.Build(z)
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

	err = z.AfterBuild(z)
	if err != nil {
		panic(err)
	}

	err = z.Init(z)
	if err != nil {
		panic(err)
	}
	defer z.Close(z)

	err = z.AfterInit(z)
	if err != nil {
		panic(err)
	}

	z.App.Run(":" + z.Port)
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
		case "generateDoc":
			err = z.Swagger.Generate()
			if err != nil {
				return true, err
			}
			return true, nil
		default:
			fmt.Println("Unknown command. Available: resetDB, updateDB, generateDoc")
			return true, nil
		}
	}

	return false, nil
}

func buildApp(z *Zest) error {
	deps := usecases.DependencyDirectory.Get()

	store := infrastructure.NewGormStore()
	accountRepo := resources.NewAccountRepo(store)
	aclRepo := resources.NewAclRepo(store)
	sessionRepo := resources.NewSessionRepo(store)

	deps = append(
		deps,

		store,
		accountRepo,
		aclRepo,
		sessionRepo,
		httptreemux.New(),
		infrastructure.NewRender(),
		infrastructure.NewCacheStore(),
		usecases.NewPermissionCacheInter(accountRepo, aclRepo, infrastructure.NewCacheStore(), infrastructure.NewCacheStore()),
		usecases.NewSessionCacheInter(sessionRepo, infrastructure.NewLRUCacheStore(1024)),
		usecases.NewRouteDirectory,
		usecases.NewPermissionInter,
	)

	z.Injector.RegisterMultiple(deps)

	err := z.Injector.Populate()
	if err != nil {
		return err
	}

	type dependencies struct {
		RouteDir *usecases.RouteDirectory
		Render   *infrastructure.Render
	}

	d := &dependencies{}
	err = z.Injector.Get(d)
	if err != nil {
		return err
	}

	z.Swagger.AddRoutes(d.RouteDir)

	return nil
}

func initApp(z *Zest) error {
	type dependencies struct {
		Router *httptreemux.TreeMux

		Store *infrastructure.GormStore

		SessionCacheInter    *usecases.SessionCacheInter
		PermissionCacheInter *usecases.PermissionCacheInter
		AccountGuestInter    *resources.AccountGuestInter

		RouteDir *usecases.RouteDirectory
		Render   *infrastructure.Render
	}

	d := &dependencies{}
	err := z.Injector.Get(d)
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

	z.App.Use(middlewares.NewLogger())
	z.App.Use(middlewares.NewRecovery(d.Render))
	z.App.Use(cors.Default())
	z.App.Use(middlewares.NewSessions(d.AccountGuestInter))

	d.Router.RedirectBehavior = httptreemux.UseHandler

	z.App.UseHandler(d.Router)

	return nil
}

func closeApp(z *Zest) error {
	type dependencies struct {
		Store *infrastructure.GormStore
	}

	d := &dependencies{}
	err := z.Injector.Get(d)
	if err != nil {
		return err
	}

	d.Store.Close()

	return nil
}
