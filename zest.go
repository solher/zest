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
	"strconv"
	"time"

	"gopkg.in/tylerb/graceful.v1"

	"github.com/codegangsta/cli"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	"github.com/rs/cors"
	"github.com/solher/syringe"
)

// Injector provide a quick access to an instanciated injector.
var Injector = syringe.New()

// SeqFunc is the prototype of the functions present in the launch/exit sequences.
type SeqFunc func(z *Zest) error

// Zest is an aggregation of well known and efficient package, also providing
// a simple launch/exit process to the user.
//
// The launch sequence is divided into three steps:
// - The register sequence is run, allowing the user to register dependencies
// into the injector.
// - The injection is run.
// - The init sequence is run, allowing the user to properly initialize the
// freshly built app.
// Launch and exit sequences are run following the order of the arrays, at each
// start/stop of the app, thanks to Cli and the tylerb/graceful module.
type Zest struct {
	cli     *cli.App
	Context *cli.Context

	Port        int
	ExitTimeout time.Duration

	Server   *negroni.Negroni
	Injector *syringe.Syringe

	RegisterSequence []SeqFunc
	InitSequence     []SeqFunc
	ExitSequence     []SeqFunc
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
func (z *Zest) Run() error {
	return z.cli.Run(os.Args)
}

// Classic returns a new instance of Zest, with some default register and init steps:
// "classicRegister" which registers the default dependencies (Render, Bone) in the injector.
// "classicInit" which initialize the Bone router and the default middlewares in Negroni.
func Classic() *Zest {
	z := New()

	z.cli.Flags = []cli.Flag{
		cli.IntFlag{
			Name:   "port,p",
			Value:  3000,
			Usage:  "listening port",
			EnvVar: "ZEST_PORT",
		},
		cli.DurationFlag{
			Name:   "exitTimeout,t",
			Value:  10 * time.Second,
			Usage:  "graceful shutdown timeout (0 for infinite)",
			EnvVar: "ZEST_TIMEOUT",
		},
	}

	z.RegisterSequence = append(z.RegisterSequence, classicRegister)
	z.InitSequence = append(z.InitSequence, classicInit)

	return z
}

// New returns a new instance of Zest.
func New() *Zest {
	z := &Zest{
		cli:         cli.NewApp(),
		Server:      negroni.New(),
		Injector:    Injector,
		Port:        3000,
		ExitTimeout: 10 * time.Second,
	}

	z.cli.Usage = "A Zest powered service."
	z.cli.Before = z.init
	z.cli.After = z.exit
	z.cli.Action = z.run

	return z
}

func (z *Zest) init(c *cli.Context) error {
	z.Context = c

	for _, f := range z.RegisterSequence {
		if err := f(z); err != nil {
			return err
		}
	}

	if err := z.Injector.Inject(); err != nil {
		return err
	}

	for _, f := range z.InitSequence {
		if err := f(z); err != nil {
			return err
		}
	}

	return nil
}

func (z *Zest) run(c *cli.Context) {
	fmt.Printf("\n[Zest] Listening on %d\n", z.Port)

	graceful.Run(":"+strconv.Itoa(z.Port), z.ExitTimeout, z.Server)
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
	z.Injector.Register(NewRender(), bone.New())

	return nil
}

func classicInit(z *Zest) error {
	d := &struct{ Router *bone.Mux }{}

	if err := z.Injector.Get(d); err != nil {
		return err
	}

	z.Server.Use(NewRecovery())
	z.Server.Use(NewLogger())
	z.Server.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}))

	z.Server.UseHandler(d.Router)

	z.Port = z.Context.GlobalInt("port")
	z.ExitTimeout = z.Context.GlobalDuration("exitTimeout")

	return nil
}
