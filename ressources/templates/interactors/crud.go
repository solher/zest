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
	type {{.Type}}Repository interface {
		Create({{.Name}}s []{{.Type}}) ([]{{.Type}}, error)
		Find(filter *interfaces.Filter) ([]{{.Type}}, error)
		FindByID(id int) (*{{.Type}}, error)
		Upsert({{.Name}}s []{{.Type}}) ([]{{.Type}}, error)
		DeleteAll(filter *interfaces.Filter) error
		DeleteByID(id int) error
	}

	type Interactor struct {
		repo {{.Type}}Repository
	}

	func NewInteractor(repo {{.Type}}Repository) *Interactor {
		return &Interactor{repo: repo}
	}
`}

var create = &typewriter.Template{
	Name: "Create",
	Text: `
	func (i *Interactor) Create({{.Name}}s []{{.Type}}) ([]{{.Type}}, error) {
		{{.Name}}s, err := i.repo.Create({{.Name}}s)
		return {{.Name}}s, err
	}
`}

var find = &typewriter.Template{
	Name: "Find",
	Text: `
	func (i *Interactor) Find(filter *interfaces.Filter) ([]{{.Type}}, error) {
		{{.Name}}s, err := i.repo.Find(filter)
		return {{.Name}}s, err
	}
`}

var findByID = &typewriter.Template{
	Name: "FindByID",
	Text: `
	func (i *Interactor) FindByID(id int) (*{{.Type}}, error) {
		{{.Name}}, err := i.repo.FindByID(id)
		return {{.Name}}, err
	}
`}

var upsert = &typewriter.Template{
	Name: "Upsert",
	Text: `
	func (i *Interactor) Upsert({{.Name}}s []{{.Type}}) ([]{{.Type}}, error) {
		{{.Name}}s, err := i.repo.Upsert({{.Name}}s)
		return {{.Name}}s, err
	}
`}

var deleteAll = &typewriter.Template{
	Name: "DeleteAll",
	Text: `
	func (i *Interactor) DeleteAll(filter *interfaces.Filter) error {
		err := i.repo.DeleteAll(filter)
		return err
	}
`}

var deleteByID = &typewriter.Template{
	Name: "DeleteByID",
	Text: `
	func (i *Interactor) DeleteByID(id int) error {
		err := i.repo.DeleteByID(id)
		return err
	}
`}
