# Zest

Zest is a lightweight framework based on the [Cli](https://github.com/codegangsta/cli) package allowing clean and easy command line interfaces, the [Negroni](https://github.com/codegangsta/negroni) middleware handler, and the  [Syringe](https://github.com/solher/syringe) injector.

Zest encourages the use of small, well chosen individual dependencies instead of high productivity, full-stack frameworks.

## Installation

To install Zest:

    go get github.com/solher/zest

## Init/exit sequences
## zest.Classic()


## Example

**main.go**

```go
package main

import (
	"github.com/dimfeld/httptreemux"
	"github.com/solher/zest"
)

func main() {
	app := zest.Classic()

	cli := app.Cli()
	cli.Name = "Test"
	cli.Usage = "Test app"
	app.SetCli(cli)

	app.InitSequence = append(app.InitSequence, SetRoutes)

	app.Run()
}

func SetRoutes(z *zest.Zest) error {
	type deps struct {
		Router *httptreemux.TreeMux
		Ctrl   *Controller
	}

	d := &deps{}

	if err := z.Injector.Get(d); err != nil {
		return err
	}

	d.Router.GET("/", d.Ctrl.Handler)

	return nil
}
```

**controller.go**

```go
package main

import (
	"net/http"

	"github.com/solher/zest"
)

func init() {
	zest.Injector.Register(NewController)
}

type Controller struct {
	m *Model
	r *zest.Render
}

func NewController(m *Model, r *zest.Render) *Controller {
	return &Controller{m: m, r: r}
}

func (c *Controller) Handler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	result, err := c.m.Action()
	if err != nil {
		apiErr := &zest.APIError{Description: "An error occured. Please retry later.", ErrorCode: "UNKNOWN_ERROR"}
		c.r.JSONError(w, http.StatusInternalServerError, apiErr, err)
		return
	}

	c.r.JSON(w, http.StatusOK, result)
}
```

**model.go**

```go
package main

import "github.com/solher/zest"

func init() {
	zest.Injector.Register(NewModel)
}

type Model struct {
	s *Store
}

func NewModel(s *Store) *Model {
	return &Model{s: s}
}

func (m *Model) Action() (string, error) {
	result, err := m.s.DBAction()

	return result, err
}
```

**store.go**

```go
package main

import "github.com/solher/zest"

func init() {
	zest.Injector.Register(NewStore)
}

type Store struct{}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) DBAction() (string, error) {
	return "foobar", nil
}
```

## Why Zest ?

Because having a good cli interface, a simple init/exit process and having your app built automatically should be the basis of your development.

## About

Thanks to the [Code Gangsta](http://codegangsta.io/) for his amazing work on [Negroni](https://github.com/codegangsta/negroni) and [Cli](https://github.com/codegangsta/cli).

## License

MIT
