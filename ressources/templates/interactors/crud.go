package interactors

import (
	"github.com/Solher/zest/ressources/templates"
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
	type Abstract{{.Type}}Repo interface {
		Create({{.Name}}s []domain.{{.Type}}) ([]domain.{{.Type}}, error)
		CreateOne({{.Name}} *domain.{{.Type}}) (*domain.{{.Type}}, error)
		Find(filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.{{.Type}}, error)
		FindByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.{{.Type}}, error)
		Update({{.Name}}s []domain.{{.Type}}, filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.{{.Type}}, error)
		UpdateByID(id int, {{.Name}} *domain.{{.Type}},	filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.{{.Type}}, error)
		DeleteAll(filter *usecases.Filter, ownerRelations []domain.Relation) error
		DeleteByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) error
		Raw(query string, values ...interface{}) (*sql.Rows, error)
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
	func (i *{{.Type}}Inter) Create({{.Name}}s []domain.{{.Type}}) ([]domain.{{.Type}}, error) {
		var err error

		for i := range {{.Name}}s {
			err = (&{{.Name}}s[i]).BeforeCreate()
			if err != nil {
				return nil, err
			}
		}

		{{.Name}}s, err = i.repo.Create({{.Name}}s)
		if err != nil {
			return nil, err
		}

		for i := range {{.Name}}s {
			err = (&{{.Name}}s[i]).AfterCreate()
			if err != nil {
				return nil, err
			}
		}

		return {{.Name}}s, nil
	}
`}

var createOne = &typewriter.Template{
	Name: "CreateOne",
	Text: `
	func (i *{{.Type}}Inter) CreateOne({{.Name}} *domain.{{.Type}}) (*domain.{{.Type}}, error) {
		err := {{.Name}}.BeforeCreate()
		if err != nil {
			return nil, err
		}

		{{.Name}}, err = i.repo.CreateOne({{.Name}})
		if err != nil {
			return nil, err
		}

		err = {{.Name}}.AfterCreate()
		if err != nil {
			return nil, err
		}

		return {{.Name}}, nil
	}
`}

var find = &typewriter.Template{
	Name: "Find",
	Text: `
	func (i *{{.Type}}Inter) Find(filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.{{.Type}}, error) {
		{{.Name}}s, err := i.repo.Find(filter, ownerRelations)
		if err != nil {
			return nil, err
		}

		return {{.Name}}s, nil
	}
`}

var findByID = &typewriter.Template{
	Name: "FindByID",
	Text: `
	func (i *{{.Type}}Inter) FindByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.{{.Type}}, error) {
		{{.Name}}, err := i.repo.FindByID(id, filter, ownerRelations)
		if err != nil {
			return nil, err
		}

		return {{.Name}}, nil
	}
`}

var upsert = &typewriter.Template{
	Name: "Upsert",
	Text: `
	func (i *{{.Type}}Inter) Upsert({{.Name}}s []domain.{{.Type}}, filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.{{.Type}}, error) {
		{{.Name}}sToUpdate := []domain.{{.Type}}{}
		{{.Name}}sToCreate := []domain.{{.Type}}{}

		for i := range {{.Name}}s {
			var err error

			if {{.Name}}s[i].ID != 0 {
				err = (&{{.Name}}s[i]).BeforeUpdate()
				{{.Name}}sToUpdate = append({{.Name}}sToUpdate, {{.Name}}s[i])
			} else {
				err = (&{{.Name}}s[i]).BeforeCreate()
				{{.Name}}sToCreate = append({{.Name}}sToCreate, {{.Name}}s[i])
			}

			if err != nil {
				return nil, err
			}
		}

		{{.Name}}sToUpdate, err := i.repo.Update({{.Name}}sToUpdate, filter, ownerRelations)
		if err != nil {
			return nil, err
		}

		{{.Name}}sToCreate, err = i.repo.Create({{.Name}}sToCreate)
		if err != nil {
			return nil, err
		}

		for i := range {{.Name}}sToUpdate {
			err = (&{{.Name}}s[i]).AfterUpdate()
			if err != nil {
				return nil, err
			}
		}

		for i := range {{.Name}}sToCreate {
			err = (&{{.Name}}s[i]).AfterCreate()
			if err != nil {
				return nil, err
			}
		}

		return append({{.Name}}sToUpdate, {{.Name}}sToCreate...), nil
	}
`}

var upsertOne = &typewriter.Template{
	Name: "UpsertOne",
	Text: `
	func (i *{{.Type}}Inter) UpsertOne({{.Name}} *domain.{{.Type}}, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.{{.Type}}, error) {
		if {{.Name}}.ID != 0 {
			err := {{.Name}}.BeforeUpdate()
			if err != nil {
				return nil, err
			}

			{{.Name}}, err = i.repo.UpdateByID({{.Name}}.ID, {{.Name}}, filter, ownerRelations)
			if err != nil {
				return nil, err
			}

			err = {{.Name}}.AfterUpdate()
			if err != nil {
				return nil, err
			}
		} else {
			err := {{.Name}}.BeforeCreate()
			if err != nil {
				return nil, err
			}

			{{.Name}}, err = i.repo.CreateOne({{.Name}})
			if err != nil {
				return nil, err
			}

			err = {{.Name}}.AfterCreate()
			if err != nil {
				return nil, err
			}
		}

		return {{.Name}}, nil
	}
`}

var updateByID = &typewriter.Template{
	Name: "UpdateByID",
	Text: `
	func (i *{{.Type}}Inter) UpdateByID(id int, {{.Name}} *domain.{{.Type}},
		filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.{{.Type}}, error) {

		err := {{.Name}}.BeforeUpdate()
		if err != nil {
			return nil, err
		}

		{{.Name}}, err = i.repo.UpdateByID(id, {{.Name}}, filter, ownerRelations)
		if err != nil {
			return nil, err
		}

		err = {{.Name}}.AfterUpdate()
		if err != nil {
			return nil, err
		}

		return {{.Name}}, nil
	}
`}

var deleteAll = &typewriter.Template{
	Name: "DeleteAll",
	Text: `
	func (i *{{.Type}}Inter) DeleteAll(filter *usecases.Filter, ownerRelations []domain.Relation) error {
		{{.Name}}s, err := i.repo.Find(filter, ownerRelations)
		if err != nil {
			return err
		}

		for i := range {{.Name}}s {
			err = (&{{.Name}}s[i]).BeforeDelete()
			if err != nil {
				return err
			}
		}

		err = i.repo.DeleteAll(filter, ownerRelations)
		if err != nil {
			return err
		}

		for i := range {{.Name}}s {
			err = (&{{.Name}}s[i]).AfterDelete()
			if err != nil {
				return err
			}
		}

		return nil
	}
`}

var deleteByID = &typewriter.Template{
	Name: "DeleteByID",
	Text: `
	func (i *{{.Type}}Inter) DeleteByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) error {
		{{.Name}}, err := i.repo.FindByID(id, filter, ownerRelations)
		if err != nil {
			return err
		}

		err = {{.Name}}.BeforeDelete()
		if err != nil {
			return err
		}

		err = i.repo.DeleteByID(id, filter, ownerRelations)
		if err != nil {
			return err
		}

		err = {{.Name}}.AfterDelete()
		if err != nil {
			return err
		}

		return nil
	}
`}
