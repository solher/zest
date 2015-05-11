package interactors

import (
	"github.com/Solher/auth-scaffold/ressources/templates"
	"github.com/clipperhouse/typewriter"
)

func init() {
	imports := []typewriter.ImportSpec{}

	err := typewriter.Register(templates.NewWrite("interactor", slice, imports))
	if err != nil {
		panic(err)
	}
}

var slice = typewriter.TemplateSlice{
	interactor,
	create,
	find,
	findByID,
	upsert,
	deleteAll,
	deleteByID,
}

var interactor = &typewriter.Template{
	Name: "Interactor",
	Text: `
	type Abstract{{.Type}}Repo interface {
		Create({{.Name}}s []{{.Type}}) ([]{{.Type}}, error)
		Find(filter *interfaces.Filter) ([]{{.Type}}, error)
		FindByID(id int, filter *interfaces.Filter) (*{{.Type}}, error)
		Upsert({{.Name}}s []{{.Type}}) ([]{{.Type}}, error)
		DeleteAll(filter *interfaces.Filter) error
		DeleteByID(id int) error
	}

	type {{.Type}}Inter struct {
		repo Abstract{{.Type}}Repo
	}

	func New{{.Type}}Inter(repo Abstract{{.Type}}Repo) *{{.Type}}Inter {
		return &{{.Type}}Inter{repo: repo}
	}
`}

var create = &typewriter.Template{
	Name: "Create",
	Text: `
	func (i *{{.Type}}Inter) Create({{.Name}}s []{{.Type}}) ([]{{.Type}}, error) {
		{{.Name}}s, err := i.repo.Create({{.Name}}s)
		return {{.Name}}s, err
	}
`}

var find = &typewriter.Template{
	Name: "Find",
	Text: `
	func (i *{{.Type}}Inter) Find(filter *interfaces.Filter) ([]{{.Type}}, error) {
		{{.Name}}s, err := i.repo.Find(filter)
		return {{.Name}}s, err
	}
`}

var findByID = &typewriter.Template{
	Name: "FindByID",
	Text: `
	func (i *{{.Type}}Inter) FindByID(id int, filter *interfaces.Filter) (*{{.Type}}, error) {
		{{.Name}}, err := i.repo.FindByID(id, filter)
		return {{.Name}}, err
	}
`}

var upsert = &typewriter.Template{
	Name: "Upsert",
	Text: `
	func (i *{{.Type}}Inter) Upsert({{.Name}}s []{{.Type}}) ([]{{.Type}}, error) {
		{{.Name}}s, err := i.repo.Upsert({{.Name}}s)
		return {{.Name}}s, err
	}
`}

var deleteAll = &typewriter.Template{
	Name: "DeleteAll",
	Text: `
	func (i *{{.Type}}Inter) DeleteAll(filter *interfaces.Filter) error {
		err := i.repo.DeleteAll(filter)
		return err
	}
`}

var deleteByID = &typewriter.Template{
	Name: "DeleteByID",
	Text: `
	func (i *{{.Type}}Inter) DeleteByID(id int) error {
		err := i.repo.DeleteByID(id)
		return err
	}
`}
