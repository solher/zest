package zest

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/tylerb/graceful.v1"

	"github.com/codegangsta/cli"
	"github.com/codegangsta/negroni"
	"github.com/dimfeld/httptreemux"
	"github.com/rs/cors"
	"github.com/solher/syringe"
)

var Injector = syringe.New()

type ZestFunc func(z *Zest) error

type Zest struct {
	cli     *cli.App
	Context *cli.Context

	Server   *negroni.Negroni
	Injector *syringe.Syringe

	InitSequence []ZestFunc
	ExitSequence []ZestFunc
}

func (z *Zest) Cli() cli.App {
	return *z.cli
}

func (z *Zest) SetCli(cli cli.App) {
	*z.cli = cli
}

func (z *Zest) Run() {
	z.cli.Run(os.Args)
}

func Classic() *Zest {
	z := New()

	z.InitSequence = append([]ZestFunc{classicRegister}, z.InitSequence...)
	z.InitSequence = append(z.InitSequence, classicInit)

	return z
}

func New() *Zest {
	z := &Zest{
		cli:      cli.NewApp(),
		Server:   negroni.New(),
		Injector: Injector,
	}

	z.InitSequence = append(z.InitSequence, classicBuild)

	z.cli.Usage = "A Zest powered service."
	z.cli.Before = z.init
	z.cli.After = z.exit
	z.cli.Action = z.run
	z.cli.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port,p",
			Value: 3000,
			Usage: "listening port",
		},
		cli.IntFlag{
			Name:  "exitTimeout,t",
			Value: 10,
			Usage: "graceful shutdown timeout in seconds (0 for infinite)",
		},
	}

	return z
}

func (z *Zest) init(c *cli.Context) error {
	z.Context = c

	for _, f := range z.InitSequence {
		if err := f(z); err != nil {
			return err
		}
	}

	return nil
}

func (z *Zest) run(c *cli.Context) {
	z.Context = c

	port := fmt.Sprintf(":%d", z.Context.GlobalInt("port"))
	exitTimeout := time.Duration(z.Context.GlobalInt("exitTimeout")) * time.Second

	graceful.Run(port, exitTimeout, z.Server)
}

func (z *Zest) exit(c *cli.Context) error {
	z.Context = c

	for _, f := range z.ExitSequence {
		if err := f(z); err != nil {
			return err
		}
	}

	return nil
}

func classicRegister(z *Zest) error {
	z.Injector.Register(NewRender(), httptreemux.New())

	return nil
}

func classicBuild(z *Zest) error {
	return z.Injector.Inject()
}

func classicInit(z *Zest) error {
	router := &httptreemux.TreeMux{}

	if err := z.Injector.GetOne(router); err != nil {
		return err
	}

	z.Server.UseHandler(router)

	z.Server.Use(NewLogger())
	z.Server.Use(NewRecovery())
	z.Server.Use(cors.Default())

	return nil
}
