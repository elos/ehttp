package templates

import (
	"net/http"
	"text/template"

	"github.com/elos/transfer"
	"github.com/julienschmidt/httprouter"
)

type (
	TemplateSet map[Name][]string

	TemplateMap map[Name]*template.Template

	Context struct {
		rootDir string
		tset    *TemplateSet
		tmap    *TemplateMap
	}
)

func NewContext(rootDir string, tset *TemplateSet) *Context {
	return &Context{
		rootDir: rootDir,
		tset:    tset,
		tmap:    new(TemplateMap),
	}
}

/*
	Show is for rendering templates that require
	no specific data
*/
func (c *Context) Show(name Name) httprouter.Handle {
	return c.Template(name, nil)
}

/*
	Template is a httprouter.Handle curried function to inject
	the template name and data

	You can only really use this if the data is constant.
*/
func (c *Context) Template(name Name, data interface{}) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		conn := transfer.NewHTTPConnection(w, r, nil)
		CatchError(conn, c.Render(conn, name, data))
	}
}

/*
	renderTemplate is the internally used implementation of rendering a named
	template with the supplied data
*/
func (c *Context) Render(connection *transfer.HTTPConnection, name Name, data interface{}) error {
	err := c.ParseHTMLTemplates()
	if err != nil {
		return err
	}

	t, ok := (*c.tmap)[name]

	if !ok {
		return NewNotFoundError(name)
	}

	if err := t.Execute(connection.ResponseWriter(), data); err != nil {
		return NewRenderError(err)
	}

	return nil
}

func (c *Context) ParseHTMLTemplates() error {
	for name, set := range *c.tset {
		t, err := template.ParseFiles(JoinDir(c.rootDir, set)...)
		if err != nil {
			return err
		}
		(*c.tmap)[name] = t
	}
	return nil
}
