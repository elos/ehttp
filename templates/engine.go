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

	Context interface {
		WithData(interface{}) Context
	}

	Engine struct {
		rootDir       string
		tset          *TemplateSet
		tmap          *TemplateMap
		globalContext Context
		everyload     bool
	}
)

func NewEngine(rootDir string, tset *TemplateSet) *Engine {
	tm := make(TemplateMap)

	return &Engine{
		rootDir: rootDir,
		tset:    tset,
		tmap:    &tm,
	}
}

func (e *Engine) WithContext(c Context) *Engine {
	return &Engine{
		rootDir:       e.rootDir,
		tset:          e.tset,
		tmap:          e.tmap,
		globalContext: c,
	}
}

func (e *Engine) WithEveryLoad() *Engine {
	return &Engine{
		rootDir:       e.rootDir,
		tset:          e.tset,
		tmap:          e.tmap,
		globalContext: e.globalContext,
		everyload:     true,
	}
}

// Show is a shortcut for rendering templates that
// have no specific data. Perhaps an index page.
func (e *Engine) Show(name Name) httprouter.Handle {
	return e.Handle(name, e.globalContext)
}

// Handle is a httprouter.Handle that is partially curried to
// inject the template nate and data
// Note: you can only really use this if the data is constant.
func (e *Engine) Handle(name Name, data interface{}) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		conn := transfer.NewHTTPConnection(w, r, nil)
		CatchError(conn, e.Render(conn, name, data))
	}
}

// Render is used by most other engine functions and actually performs the lookup of the template
// name, the wrapping of the data in the engine's globalContext
func (e *Engine) Render(connection *transfer.HTTPConnection, name Name, data interface{}) error {
	if e.everyload {
		err := e.ParseHTMLTemplates()
		if err != nil {
			return err
		}
	}

	t, ok := (*e.tmap)[name]

	if !ok {
		return NewNotFoundError(name)
	}

	if err := t.Execute(connection.ResponseWriter(), e.globalContext.WithData(data)); err != nil {
		return NewRenderError(err)
	}

	return nil
}

// Must be executed at least once to load templates, if the template set changes post-hoc
// you must recall PaseHTMLTemplates() to see the changes
func (c *Engine) ParseHTMLTemplates() error {
	for name, set := range *c.tset {
		t, err := template.ParseFiles(JoinDir(c.rootDir, set)...)
		if err != nil {
			return err
		}
		(*c.tmap)[name] = t
	}
	return nil
}
