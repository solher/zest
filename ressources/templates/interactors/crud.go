package interactors

import (
	"github.com/solher/zest/ressources/templates"
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
	createOne,
	find,
	findByID,
	upsert,
	upsertOne,
	updateByID,
	deleteAll,
	deleteByID,
}

var interactor = &typewriter.Template{
	Name: "Interactor",
	Text: `
	func init() {
		usecases.DependencyDirectory.Register(New{{.Type}}Inter)
	}

	type Abstract{{.Type}}Repo interface {
		Create({{.Name}}s []domain.{{.Type}}) ([]domain.{{.Type}}, error)
		CreateOne({{.Name}} *domain.{{.Type}}) (*domain.{{.Type}}, error)
		Find(context usecases.QueryContext) ([]domain.{{.Type}}, error)
		FindByID(id int, context usecases.QueryContext) (*domain.{{.Type}}, error)
		Update({{.Name}}s []domain.{{.Type}}, context usecases.QueryContext) ([]domain.{{.Type}}, error)
		UpdateByID(id int, {{.Name}} *domain.{{.Type}},	context usecases.QueryContext) (*domain.{{.Type}}, error)
		DeleteAll(context usecases.QueryContext) error
		DeleteByID(id int, context usecases.QueryContext) error
		Raw(query string, values ...interface{}) (*sql.Rows, error)
	}

	type {{.Type}}Inter struct {
		repo Abstract{{.Type}}Repo
	}

	func New{{.Type}}Inter(repo Abstract{{.Type}}Repo) *{{.Type}}Inter {
		return &{{.Type}}Inter{repo: repo}
	}

	func (i *{{.Type}}Inter) BeforeSave({{.Name}} *domain.{{.Type}}) error {
		{{.Name}}.ID = 0
		{{.Name}}.CreatedAt = time.Time{}
		{{.Name}}.UpdatedAt = time.Time{}

		err := {{.Name}}.ScopeModel()
		if err != nil {
			return err
		}

		return nil
	}
`}

var create = &typewriter.Template{
	Name: "Create",
	Text: `
	func (i *{{.Type}}Inter) Create({{.Name}}s []domain.{{.Type}}) ([]domain.{{.Type}}, error) {
		var err error

		for k := range {{.Name}}s {
			err := i.BeforeSave(&{{.Name}}s[k])
			if err != nil {
				return nil, err
			}
		}

		{{.Name}}s, err = i.repo.Create({{.Name}}s)
		if err != nil {
			return nil, err
		}

		return {{.Name}}s, nil
	}
`}

var createOne = &typewriter.Template{
	Name: "CreateOne",
	Text: `
	func (i *{{.Type}}Inter) CreateOne({{.Name}} *domain.{{.Type}}) (*domain.{{.Type}}, error) {
		err := i.BeforeSave({{.Name}})
		if err != nil {
			return nil, err
		}

		{{.Name}}, err = i.repo.CreateOne({{.Name}})
		if err != nil {
			return nil, err
		}

		return {{.Name}}, nil
	}
`}

var find = &typewriter.Template{
	Name: "Find",
	Text: `
	func (i *{{.Type}}Inter) Find(context usecases.QueryContext) ([]domain.{{.Type}}, error) {
		{{.Name}}s, err := i.repo.Find(context)
		if err != nil {
			return nil, err
		}

		return {{.Name}}s, nil
	}
`}

var findByID = &typewriter.Template{
	Name: "FindByID",
	Text: `
	func (i *{{.Type}}Inter) FindByID(id int, context usecases.QueryContext) (*domain.{{.Type}}, error) {
		{{.Name}}, err := i.repo.FindByID(id, context)
		if err != nil {
			return nil, err
		}

		return {{.Name}}, nil
	}
`}

var upsert = &typewriter.Template{
	Name: "Upsert",
	Text: `
	func (i *{{.Type}}Inter) Upsert({{.Name}}s []domain.{{.Type}}, context usecases.QueryContext) ([]domain.{{.Type}}, error) {
		{{.Name}}sToUpdate := []domain.{{.Type}}{}
		{{.Name}}sToCreate := []domain.{{.Type}}{}

		for k := range {{.Name}}s {
			err := i.BeforeSave(&{{.Name}}s[k])
			if err != nil {
				return nil, err
			}

			if {{.Name}}s[k].ID != 0 {
				{{.Name}}sToUpdate = append({{.Name}}sToUpdate, {{.Name}}s[k])
			} else {
				{{.Name}}sToCreate = append({{.Name}}sToCreate, {{.Name}}s[k])
			}
		}

		{{.Name}}sToUpdate, err := i.repo.Update({{.Name}}sToUpdate, context)
		if err != nil {
			return nil, err
		}

		{{.Name}}sToCreate, err = i.repo.Create({{.Name}}sToCreate)
		if err != nil {
			return nil, err
		}

		return append({{.Name}}sToUpdate, {{.Name}}sToCreate...), nil
	}
`}

var upsertOne = &typewriter.Template{
	Name: "UpsertOne",
	Text: `
	func (i *{{.Type}}Inter) UpsertOne({{.Name}} *domain.{{.Type}}, context usecases.QueryContext) (*domain.{{.Type}}, error) {
		err := i.BeforeSave({{.Name}})
		if err != nil {
			return nil, err
		}

		if {{.Name}}.ID != 0 {
			{{.Name}}, err = i.repo.UpdateByID({{.Name}}.ID, {{.Name}}, context)
		} else {
			{{.Name}}, err = i.repo.CreateOne({{.Name}})
		}

		if err != nil {
			return nil, err
		}

		return {{.Name}}, nil
	}
`}

var updateByID = &typewriter.Template{
	Name: "UpdateByID",
	Text: `
	func (i *{{.Type}}Inter) UpdateByID(id int, {{.Name}} *domain.{{.Type}},
		context usecases.QueryContext) (*domain.{{.Type}}, error) {

		err := i.BeforeSave({{.Name}})
		if err != nil {
			return nil, err
		}

		{{.Name}}, err = i.repo.UpdateByID(id, {{.Name}}, context)
		if err != nil {
			return nil, err
		}

		return {{.Name}}, nil
	}
`}

var deleteAll = &typewriter.Template{
	Name: "DeleteAll",
	Text: `
	func (i *{{.Type}}Inter) DeleteAll(context usecases.QueryContext) error {
		err := i.repo.DeleteAll(context)
		if err != nil {
			return err
		}

		return nil
	}
`}

var deleteByID = &typewriter.Template{
	Name: "DeleteByID",
	Text: `
	func (i *{{.Type}}Inter) DeleteByID(id int, context usecases.QueryContext) error {
		err := i.repo.DeleteByID(id, context)
		if err != nil {
			return err
		}

		return nil
	}
`}
