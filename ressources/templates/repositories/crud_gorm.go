package repositories

import (
	"github.com/Solher/auth-scaffold/ressources/templates"
	"github.com/clipperhouse/typewriter"
)

func init() {
	imports := []typewriter.ImportSpec{}

	err := typewriter.Register(templates.NewWrite("repository", slice, imports))
	if err != nil {
		panic(err)
	}
}

var slice = typewriter.TemplateSlice{
	repository,
	create,
	createOne,
	find,
	findByID,
	upsert,
	upsertOne,
	deleteAll,
	deleteByID,
}

var repository = &typewriter.Template{
	Name: "Repository",
	Text: `
  type {{.Type}}Repo struct {
  	store interfaces.AbstractGormStore
  }

  func New{{.Type}}Repo(store interfaces.AbstractGormStore) *{{.Type}}Repo {
  	return &{{.Type}}Repo{store: store}
  }
`}

var create = &typewriter.Template{
	Name: "Create",
	Text: `
  func (r *{{.Type}}Repo) Create({{.Name}}s []{{.Type}}) ([]{{.Type}}, error) {
  	db := r.store.GetDB()
  	transaction := db.Begin()

  	for i, {{.Name}} := range {{.Name}}s {
  		err := db.Create(&{{.Name}}).Error
  		if err != nil {
  			transaction.Rollback()
  			return nil, err
  		}

      {{.Name}}s[i] = {{.Name}}
  	}

  	transaction.Commit()
  	return {{.Name}}s, nil
  }
`}

var createOne = &typewriter.Template{
	Name: "CreateOne",
	Text: `
	func (r *{{.Type}}Repo) CreateOne({{.Name}} *{{.Type}}) (*{{.Type}}, error) {
		db := r.store.GetDB()

		err := db.Create({{.Name}}).Error
		if err != nil {
			return nil, err
		}

		return {{.Name}}, nil
	}
`}

var find = &typewriter.Template{
	Name: "Find",
	Text: `
	func (r *{{.Type}}Repo) Find(filter *interfaces.Filter) ([]{{.Type}}, error) {
		query, err := r.store.BuildQuery(filter)
		if err != nil {
			return nil, err
		}

		{{.Name}}s := []{{.Type}}{}

		err = query.Find(&{{.Name}}s).Error
		if err != nil {
			return nil, err
		}

		return {{.Name}}s, nil
	}
`}

var findByID = &typewriter.Template{
	Name: "FindByID",
	Text: `
	func (r *{{.Type}}Repo) FindByID(id int, filter *interfaces.Filter) (*{{.Type}}, error) {
		query, err := r.store.BuildQuery(filter)
		if err != nil {
			return nil, err
		}

		{{.Name}} := {{.Type}}{}

		err = query.First(&{{.Name}}, id).Error
		if err != nil {
			return nil, err
		}

		return &{{.Name}}, nil
	}
`}

var upsert = &typewriter.Template{
	Name: "Upsert",
	Text: `
	func (r *{{.Type}}Repo) Upsert({{.Name}}s []{{.Type}}) ([]{{.Type}}, error) {
		db := r.store.GetDB()
		transaction := db.Begin()

		for i, {{.Name}} := range {{.Name}}s {
			if {{.Name}}.ID != 0 {
				oldUser := {{.Type}}{}

				err := db.First(&oldUser, {{.Name}}.ID).Updates({{.Name}}).Error
				if err != nil {
					transaction.Rollback()
					return nil, err
				}
			} else {
				err := db.Create(&{{.Name}}).Error
				if err != nil {
					transaction.Rollback()
					return nil, err
				}
			}

			{{.Name}}s[i] = {{.Name}}
		}

		transaction.Commit()
		return {{.Name}}s, nil
	}
`}

var upsertOne = &typewriter.Template{
	Name: "UpsertOne",
	Text: `
	func (r *{{.Type}}Repo) UpsertOne({{.Name}} *{{.Type}}) (*{{.Type}}, error) {
		db := r.store.GetDB()

		if {{.Name}}.ID != 0 {
			oldUser := {{.Type}}{}

			err := db.First(&oldUser, {{.Name}}.ID).Updates({{.Name}}).Error
			if err != nil {
				return nil, err
			}
		} else {
			err := db.Create(&{{.Name}}).Error
			if err != nil {
				return nil, err
			}
		}

		return {{.Name}}, nil
	}
`}

var deleteAll = &typewriter.Template{
	Name: "DeleteAll",
	Text: `
	func (r *{{.Type}}Repo) DeleteAll(filter *interfaces.Filter) error {
		query, err := r.store.BuildQuery(filter)
		if err != nil {
			return err
		}

		err = query.Delete({{.Type}}{}).Error
		if err != nil {
			return err
		}

		return nil
	}
`}

var deleteByID = &typewriter.Template{
	Name: "DeleteByID",
	Text: `
	func (r *{{.Type}}Repo) DeleteByID(id int) error {
		db := r.store.GetDB()

		err := db.Delete(&{{.Type}}{ID: id}).Error
		if err != nil {
			return err
		}

		return nil
	}
`}
