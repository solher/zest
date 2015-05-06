package templates

import (
	"io"
	"strings"

	"github.com/clipperhouse/typewriter"
)

type (
	Write struct {
		name      string
		templates typewriter.TemplateSlice
	}

	model struct {
		Type typewriter.Type
		Name string
	}
)

func NewWrite(name string, templates typewriter.TemplateSlice) *Write {
	return &Write{name: name, templates: templates}
}

func (sw *Write) Name() string {
	return sw.name
}

func (sw *Write) Imports(t typewriter.Type) (result []typewriter.ImportSpec) {
	// none
	return result
}

func (sw *Write) Write(w io.Writer, t typewriter.Type) error {
	tag, found := t.FindTag(sw)

	if !found {
		// nothing to be done
		return nil
	}

	tmpl, err := sw.templates.ByTag(t, tag)
	if err != nil {
		return err
	}

	m := model{
		Type: t,
		Name: strings.ToLower(t.String()),
	}

	if err := tmpl.Execute(w, m); err != nil {
		return err
	}

	values := tag.Values

	for _, value := range values {
		tmpl, err := sw.templates.ByTagValue(t, value)

		if err != nil {
			return err
		}

		if err := tmpl.Execute(w, m); err != nil {
			return err
		}
	}

	return nil
}
