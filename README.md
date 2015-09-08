# Zest

Zest is a lightweight framework based on the [Cli](https://github.com/codegangsta/cli) package allowing clean and easy command line interfaces, the [Negroni](https://github.com/codegangsta/negroni) middleware handler, and the [Syringe](https://github.com/solher/syringe) injector.

Zest encourages the use of small, well chosen individual dependencies instead of high productivity, full-stack frameworks.

## Overview

Having a good cli interface, a simple init/exit process and your app injected automatically should be the basis of your development.

Zest makes all that simple by aggregating well known and efficient packages. The `Classic` version also provides some default tools useful for most applications :
- [Gin](https://github.com/gin-gonic/gin) inspired logging, [CORS](https://github.com/rs/cors) (allowing all origins by default) and Recovery middlewares
- Pre-injected custom JSON renderer and [Httptreemux](https://github.com/dimfeld/httptreemux) router

## Installation

To install Zest:

    go get github.com/solher/zest

## Init/exit sequences

Init and exit sequences are run following the order of the array, at each start/stop of the app, thanks to Cli and the [graceful shutdown](https://github.com/tylerb/graceful) module.

```go
type ZestFunc func(z *Zest) error

type Zest struct {
	cli     *cli.App
	Context *cli.Context

	Server   *negroni.Negroni
	Injector *syringe.Syringe

	InitSequence []ZestFunc
	ExitSequence []ZestFunc
}
```

Each function is called with the Zest app in argument, allowing the user to interact with the cli context, the Negroni server or any dependency registered in the injector.

In the `New` version of Zest, the init sequence only contains a single `Build` step, triggering the dependency injection of the app.

In the `Classic` version, the init sequence contains three steps, in this order :
- `Register` which registers the default dependencies (Render, Httptreemux) in the injector.
- `Build` which triggers the dependency injection of the app.
- `Init` which initialize the Httptreemux router and the default middlewares in Negroni.

In both versions, the exit sequence is empty.

## API errors

One of the few conventions established by Zest is the API error messages style. It allows a consistent format between the recovery middleware and the render methods, and a better expressiveness.

```go
type APIError struct {
	Status      int    `json:"status"`
	Description string `json:"description"`
	Raw         string `json:"raw"`
	ErrorCode   string `json:"errorCode"`
}
```

- The `Status` is a copy of the status returned in the HTTP headers.
- The `Description` is the message describing what kind of error occured.
- The `Raw` field is the raw error message which triggered the API error. Its purpose is to allow a more efficient debugging and should not be used as an "error id" by the API client.
- The `ErrorCode` is the "machine oriented" description of the API error.

## Render

The custom `Render` module is based on the [Render](https://github.com/unrolled/render) package, made more expressive thanks to the coupling with the Zest API error format.

```go
func (r *Render) JSONError(w http.ResponseWriter, status int, apiError *APIError, err error)

func (r *Render) JSON(w http.ResponseWriter, status int, object interface{})
```

When the `JSONError` method is called, the status and the error are automatically copied into the final `APIError` struct so you don't have to worry about that.

In situation, it looks like that :

```go
var UnknownAPIError = &zest.APIError{Description: "An error occured. Please retry later.", ErrorCode: "UNKNOWN_ERROR"}

func (c *Controller) Handler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	result, err := c.m.Action()
	if err != nil {
		c.r.JSONError(w, http.StatusInternalServerError, UnknownAPIError, err)
		return
	}

	c.r.JSON(w, http.StatusOK, result)
}
```

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
	cli.Usage = "This is a test application."
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

## About

Thanks to the [Code Gangsta](http://codegangsta.io/) for his amazing work on [Negroni](https://github.com/codegangsta/negroni) and [Cli](https://github.com/codegangsta/cli).

## License

MIT
