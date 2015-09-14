// Package zest contains a lightweight framework based on the codegangsta/cli package allowing
// clean and easy command line interfaces, the codegangsta/negroni middleware
// handler, and the solher/syringe injector.
//
// Zest encourages the use of small, well chosen individual dependencies
// instead of high productivity, full-stack frameworks.
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

// Injector provide a quick access to an instanciated injector.
var Injector = syringe.New()

// SeqFunc is the prototype of the functions present in the exit/init sequences.
type SeqFunc func(z *Zest) error

// Zest is an aggregation of well known and efficient package, also providing
// a simple init/exit process to the user.
//
// Init and exit sequences are run following the order of the array, at each
// start/stop of the app, thanks to Cli and the tylerb/graceful module.
type Zest struct {
	cli     *cli.App
	Context *cli.Context

	Server   *negroni.Negroni
	Injector *syringe.Syringe

	InitSequence []SeqFunc
	ExitSequence []SeqFunc
}

// Cli returns a copy of the embedded Cli app.
func (z *Zest) Cli() cli.App {
	return *z.cli
}

// SetCli sets a copy of the embedded Cli app.
func (z *Zest) SetCli(cli cli.App) {
	*z.cli = cli
}

// Run starts the Cli app.
func (z *Zest) Run() {
	z.cli.Run(os.Args)
}

// Classic returns a new instance of Zest, with some default init steps, in this order :
// "classicRegister" which registers the default dependencies (Render, Httptreemux) in the injector.
// "classicBuild" which triggers the dependency injection of the app.
// "classicInit" which initialize the Httptreemux router and the default middlewares in Negroni.
func Classic() *Zest {
	z := New()

	z.InitSequence = append([]SeqFunc{classicRegister}, z.InitSequence...)
	z.InitSequence = append(z.InitSequence, classicInit)

	return z
}

// New returns a new instance of Zest.
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

	fmt.Println("\n[Zest] Listening on " + port)

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
	d := &struct{ Router *httptreemux.TreeMux }{}

	if err := z.Injector.Get(d); err != nil {
		return err
	}

	z.Server.Use(NewRecovery())
	z.Server.Use(NewLogger())
	z.Server.Use(cors.Default())

	z.Server.UseHandler(d.Router)

	return nil
}
