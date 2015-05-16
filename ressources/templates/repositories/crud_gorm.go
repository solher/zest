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
  func (r *{{.Type}}Repo) Create({{.Name}}s []domain.{{.Type}}) ([]domain.{{.Type}}, error) {
  	db := r.store.GetDB()
  	transaction := db.Begin()

  	for i, {{.Name}} := range {{.Name}}s {
  		err := db.Create(&{{.Name}}).Error
  		if err != nil {
  			transaction.Rollback()

				if strings.Contains(err.Error(), "constraint") {
					return nil, internalerrors.NewViolatedConstraint(err.Error())
				} else {
					return nil, internalerrors.DatabaseError
				}
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
	func (r *{{.Type}}Repo) CreateOne({{.Name}} *domain.{{.Type}}) (*domain.{{.Type}}, error) {
		db := r.store.GetDB()

		err := db.Create({{.Name}}).Error
		if err != nil {
			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			} else {
				return nil, internalerrors.DatabaseError
			}
		}

		return {{.Name}}, nil
	}
`}

var find = &typewriter.Template{
	Name: "Find",
	Text: `
	func (r *{{.Type}}Repo) Find(filter *interfaces.Filter) ([]domain.{{.Type}}, error) {
		query, err := r.store.BuildQuery(filter)
		if err != nil {
			return nil, internalerrors.DatabaseError
		}

		{{.Name}}s := []domain.{{.Type}}{}

		err = query.Find(&{{.Name}}s).Error
		if err != nil {
			return nil, internalerrors.DatabaseError
		}

		return {{.Name}}s, nil
	}
`}

var findByID = &typewriter.Template{
	Name: "FindByID",
	Text: `
	func (r *{{.Type}}Repo) FindByID(id int, filter *interfaces.Filter) (*domain.{{.Type}}, error) {
		query, err := r.store.BuildQuery(filter)
		if err != nil {
			return nil, internalerrors.DatabaseError
		}

		{{.Name}} := domain.{{.Type}}{}

		err = query.First(&{{.Name}}, id).Error
		if err != nil {
			return nil, internalerrors.DatabaseError
		}

		return &{{.Name}}, nil
	}
`}

var upsert = &typewriter.Template{
	Name: "Upsert",
	Text: `
	func (r *{{.Type}}Repo) Upsert({{.Name}}s []domain.{{.Type}}) ([]domain.{{.Type}}, error) {
		db := r.store.GetDB()
		transaction := db.Begin()

		for i, {{.Name}} := range {{.Name}}s {
			if {{.Name}}.ID != 0 {
				oldUser := domain.{{.Type}}{}

				err := db.First(&oldUser, {{.Name}}.ID).Updates({{.Name}}).Error
				if err != nil {
					transaction.Rollback()

					if strings.Contains(err.Error(), "constraint") {
						return nil, internalerrors.NewViolatedConstraint(err.Error())
					} else {
						return nil, internalerrors.DatabaseError
					}
				}
			} else {
				err := db.Create(&{{.Name}}).Error
				if err != nil {
					transaction.Rollback()

					if strings.Contains(err.Error(), "constraint") {
						return nil, internalerrors.NewViolatedConstraint(err.Error())
					} else {
						return nil, internalerrors.DatabaseError
					}
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
	func (r *{{.Type}}Repo) UpsertOne({{.Name}} *domain.{{.Type}}) (*domain.{{.Type}}, error) {
		db := r.store.GetDB()

		if {{.Name}}.ID != 0 {
			oldUser := domain.{{.Type}}{}

			err := db.First(&oldUser, {{.Name}}.ID).Updates({{.Name}}).Error
			if err != nil {
				if strings.Contains(err.Error(), "constraint") {
					return nil, internalerrors.NewViolatedConstraint(err.Error())
				} else {
					return nil, internalerrors.DatabaseError
				}
			}
		} else {
			err := db.Create(&{{.Name}}).Error
			if err != nil {
				if strings.Contains(err.Error(), "constraint") {
					return nil, internalerrors.NewViolatedConstraint(err.Error())
				} else {
					return nil, internalerrors.DatabaseError
				}
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
			return internalerrors.DatabaseError
		}

		err = query.Delete(domain.{{.Type}}{}).Error
		if err != nil {
			return internalerrors.DatabaseError
		}

		return nil
	}
`}

var deleteByID = &typewriter.Template{
	Name: "DeleteByID",
	Text: `
	func (r *{{.Type}}Repo) DeleteByID(id int) error {
		db := r.store.GetDB()

		err := db.Delete(&domain.{{.Type}}{GormModel: domain.GormModel{ID: id}}).Error
		if err != nil {
			return internalerrors.DatabaseError
		}

		return nil
	}
`}
