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

		for k := range {{.Name}}s {
			err = i.BeforeCreate(&{{.Name}}s[k])
			if err != nil {
				return nil, err
			}
		}

		{{.Name}}s, err = i.repo.Create({{.Name}}s)
		if err != nil {
			return nil, err
		}

		for k := range {{.Name}}s {
			err = i.AfterCreate(&{{.Name}}s[k])
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
		err := i.BeforeCreate({{.Name}})
		if err != nil {
			return nil, err
		}

		{{.Name}}, err = i.repo.CreateOne({{.Name}})
		if err != nil {
			return nil, err
		}

		err = i.AfterCreate({{.Name}})
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

		for k := range {{.Name}}s {
			var err error

			if {{.Name}}s[k].ID != 0 {
				err = i.BeforeUpdate(&{{.Name}}s[k])
				{{.Name}}sToUpdate = append({{.Name}}sToUpdate, {{.Name}}s[k])
			} else {
				err = i.BeforeCreate(&{{.Name}}s[k])
				{{.Name}}sToCreate = append({{.Name}}sToCreate, {{.Name}}s[k])
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

		for k := range {{.Name}}sToUpdate {
			err = i.AfterUpdate(&{{.Name}}s[k])
			if err != nil {
				return nil, err
			}
		}

		for k := range {{.Name}}sToCreate {
			err = i.AfterCreate(&{{.Name}}s[k])
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
			err := i.BeforeUpdate({{.Name}})
			if err != nil {
				return nil, err
			}

			{{.Name}}, err = i.repo.UpdateByID({{.Name}}.ID, {{.Name}}, filter, ownerRelations)
			if err != nil {
				return nil, err
			}

			err = i.AfterUpdate({{.Name}})
			if err != nil {
				return nil, err
			}
		} else {
			err := i.BeforeCreate({{.Name}})
			if err != nil {
				return nil, err
			}

			{{.Name}}, err = i.repo.CreateOne({{.Name}})
			if err != nil {
				return nil, err
			}

			err = i.AfterCreate({{.Name}})
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

		err := i.BeforeUpdate({{.Name}})
		if err != nil {
			return nil, err
		}

		{{.Name}}, err = i.repo.UpdateByID(id, {{.Name}}, filter, ownerRelations)
		if err != nil {
			return nil, err
		}

		err = i.AfterUpdate({{.Name}})
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

		for k := range {{.Name}}s {
			err = i.BeforeDelete(&{{.Name}}s[k])
			if err != nil {
				return err
			}
		}

		err = i.repo.DeleteAll(filter, ownerRelations)
		if err != nil {
			return err
		}

		for k := range {{.Name}}s {
			err = i.AfterDelete(&{{.Name}}s[k])
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

		err = i.BeforeDelete({{.Name}})
		if err != nil {
			return err
		}

		err = i.repo.DeleteByID(id, filter, ownerRelations)
		if err != nil {
			return err
		}

		err = i.AfterDelete({{.Name}})
		if err != nil {
			return err
		}

		return nil
	}
`}
