package repositories

import (
	"github.com/Solher/zest/ressources/templates"
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
	update,
	updateByID,
	deleteAll,
	deleteByID,
	raw,
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
				}

				return nil, internalerrors.DatabaseError
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
			}

			return nil, internalerrors.DatabaseError
		}

		return {{.Name}}, nil
	}
`}

var find = &typewriter.Template{
	Name: "Find",
	Text: `
	func (r *{{.Type}}Repo) Find(filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.{{.Type}}, error) {
		query, err := r.store.BuildQuery(filter, ownerRelations)
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
	func (r *{{.Type}}Repo) FindByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.{{.Type}}, error) {
		query, err := r.store.BuildQuery(filter, ownerRelations)
		if err != nil {
			return nil, internalerrors.DatabaseError
		}

		{{.Name}} := domain.{{.Type}}{}

		err = query.Where("{{.Name}}s.id = ?", id).First(&{{.Name}}).Error
		if err != nil {
			return nil, internalerrors.DatabaseError
		}

		return &{{.Name}}, nil
	}
`}

var update = &typewriter.Template{
	Name: "Update",
	Text: `
	func (r *{{.Type}}Repo) Update({{.Name}}s []domain.{{.Type}}, filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.{{.Type}}, error) {
		db := r.store.GetDB()
		transaction := db.Begin()

		query, err := r.store.BuildQuery(filter, ownerRelations)
		if err != nil {
			return nil, internalerrors.DatabaseError
		}

		for i, {{.Name}} := range {{.Name}}s {
			queryCopy := *query
			oldUser := domain.{{.Type}}{}

			err := queryCopy.Where("{{.Name}}s.id = ?", {{.Name}}.ID).First(&oldUser).Updates({{.Name}}s[i]).Error
			if err != nil {
				transaction.Rollback()

				if strings.Contains(err.Error(), "constraint") {
					return nil, internalerrors.NewViolatedConstraint(err.Error())
				}

				return nil, internalerrors.DatabaseError
			}
		}

		transaction.Commit()
		return {{.Name}}s, nil
	}
`}
var updateByID = &typewriter.Template{
	Name: "UpdateByID",
	Text: `
	func (r *{{.Type}}Repo) UpdateByID(id int, {{.Name}} *domain.{{.Type}},
		filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.{{.Type}}, error) {

		query, err := r.store.BuildQuery(filter, ownerRelations)
		if err != nil {
			return nil, internalerrors.DatabaseError
		}

		oldUser := domain.{{.Type}}{}

		err = query.Where("{{.Name}}s.id = ?", id).First(&oldUser).Updates({{.Name}}).Error
		if err != nil {
			if strings.Contains(err.Error(), "constraint") {
				return nil, internalerrors.NewViolatedConstraint(err.Error())
			}

			return nil, internalerrors.DatabaseError
		}

		return {{.Name}}, nil
	}
`}

var deleteAll = &typewriter.Template{
	Name: "DeleteAll",
	Text: `
	func (r *{{.Type}}Repo) DeleteAll(filter *usecases.Filter, ownerRelations []domain.Relation) error {
		query, err := r.store.BuildQuery(filter, ownerRelations)
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
	func (r *{{.Type}}Repo) DeleteByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) error {
		query, err := r.store.BuildQuery(filter, ownerRelations)
		if err != nil {
			return internalerrors.DatabaseError
		}

		err = query.Delete(&domain.{{.Type}}{GormModel: domain.GormModel{ID: id}}).Error
		if err != nil {
			return internalerrors.DatabaseError
		}

		return nil
	}
`}

var raw = &typewriter.Template{
	Name: "Raw",
	Text: `
func (r *{{.Type}}Repo) Raw(query string, values ...interface{}) (*sql.Rows, error) {
	db := r.store.GetDB()

	rows, err := db.Raw(query, values...).Rows()
	if err != nil {
		if strings.Contains(err.Error(), "constraint") {
			return nil, internalerrors.NewViolatedConstraint(err.Error())
		}

		return nil, internalerrors.DatabaseError
	}

	return rows, nil
}
`}
