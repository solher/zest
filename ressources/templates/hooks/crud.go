package interfaces

import (
	"github.com/Solher/zest/ressources/templates"
	"github.com/clipperhouse/typewriter"
)

func init() {
	imports := []typewriter.ImportSpec{}

	err := typewriter.Register(templates.NewWrite("hooks", slice, imports))
	if err != nil {
		panic(err)
	}
}

var slice = typewriter.TemplateSlice{
	hooks,
}

var hooks = &typewriter.Template{
	Name: "Hooks",
	Text: `
	func (i *{{.Type}}Inter) BeforeCreate({{.Name}} *domain.{{.Type}}) error {
		{{.Name}}.ID = 0
		{{.Name}}.CreatedAt = time.Time{}
		{{.Name}}.UpdatedAt = time.Time{}

		err := {{.Name}}.ScopeModel()
		if err != nil {
			return err
		}

		return nil
	}

	func (i *{{.Type}}Inter) AfterCreate({{.Name}} *domain.{{.Type}}) error {
		return nil
	}

	func (i *{{.Type}}Inter) BeforeUpdate({{.Name}} *domain.{{.Type}}) error {
		{{.Name}}.ID = 0
		{{.Name}}.CreatedAt = time.Time{}
		{{.Name}}.UpdatedAt = time.Time{}

		err := {{.Name}}.ScopeModel()
		if err != nil {
			return err
		}

		return nil
	}

	func (i *{{.Type}}Inter) AfterUpdate({{.Name}} *domain.{{.Type}}) error {
		return nil
	}

	func (i *{{.Type}}Inter) BeforeDelete({{.Name}} *domain.{{.Type}}) error {
		return nil
	}

	func (i *{{.Type}}Inter) AfterDelete({{.Name}} *domain.{{.Type}}) error {
		return nil
	}
`}
