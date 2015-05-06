package interactors

import (
	"github.com/Solher/auth-scaffold/ressources/templates"
	"github.com/clipperhouse/typewriter"
)

func init() {
	imports := []typewriter.ImportSpec{
		typewriter.ImportSpec{
			Name: ".",
			Path: "github.com/smartystreets/goconvey/convey",
		},
	}

	err := typewriter.Register(templates.NewWrite("interactor_test", testSlice, imports))
	if err != nil {
		panic(err)
	}
}

var testSlice = typewriter.TemplateSlice{
	interactorTest,
}

var interactorTest = &typewriter.Template{
	Name: "Interactor_test",
	Text: `
	type stubRepository struct{}

	func (r *stubRepository) Create({{.Name}}s []{{.Type}}) ([]{{.Type}}, error) {
		return {{.Name}}s, nil
	}

	func (r *stubRepository) Find(filter *interfaces.Filter) ([]{{.Type}}, error) {
		return []{{.Type}}{}, nil
	}

	func (r *stubRepository) FindByID(id int) (*{{.Type}}, error) {
		return &{{.Type}}{}, nil
	}

	func (r *stubRepository) Upsert({{.Name}}s []{{.Type}}) ([]{{.Type}}, error) {
		return {{.Name}}s, nil
	}

	func (r *stubRepository) DeleteAll(filter *interfaces.Filter) error {
		return nil
	}

	func (r *stubRepository) DeleteByID(id int) error {
		return nil
	}

	func TestInteractor(t *testing.T) {
		repo := &stubRepository{}
		interactor := NewInteractor(repo)

		Convey("Testing {{.Name}}s interactor...", t, func() {
			Convey("Should be able to create {{.Name}}s.", func() {
				{{.Name}}s := []{{.Type}}{}
				_, err := interactor.Create({{.Name}}s)

				So(err, ShouldBeNil)
			})

			Convey("Should be able to find {{.Name}}s.", func() {
				_, err := interactor.Find(nil)

				So(err, ShouldBeNil)
			})

			Convey("Should be able to find a {{.Name}} by id.", func() {
				_, err := interactor.FindByID(1)

				So(err, ShouldBeNil)
			})

			Convey("Should be able to upsert {{.Name}}s.", func() {
				{{.Name}}s := []{{.Type}}{}
				_, err := interactor.Upsert({{.Name}}s)

				So(err, ShouldBeNil)
			})

			Convey("Should be able to delete {{.Name}}s.", func() {
				err := interactor.DeleteAll(nil)
				So(err, ShouldBeNil)
			})

			Convey("Should be able to delete a {{.Name}} by id.", func() {
				err := interactor.DeleteByID(1)
				So(err, ShouldBeNil)
			})
		})
	}
`}
