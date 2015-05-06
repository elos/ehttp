package templates

import (
	"fmt"
	"io"
	"text/template"
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
		fmap          template.FuncMap
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
		fmap:          e.fmap,
		globalContext: c,
		everyload:     e.everyload,
	}
}

func (e *Engine) WithEveryLoad() *Engine {
	return &Engine{
		rootDir:       e.rootDir,
		tset:          e.tset,
		tmap:          e.tmap,
		fmap:          e.fmap,
		globalContext: e.globalContext,
		everyload:     true,
	}
}

func (e *Engine) WithFuncMap(fm template.FuncMap) *Engine {
	return &Engine{
		rootDir:       e.rootDir,
		tset:          e.tset,
		tmap:          e.tmap,
		fmap:          fm,
		globalContext: e.globalContext,
		everyload:     e.everyload,
	}
}

func (e *Engine) Execute(w io.Writer, name Name, data interface{}) error {
	// everyload is like a development mode, always re-compile templates
	if e.everyload {
		err := e.Parse()
		if err != nil {
			return err
		}
	}

	t, ok := (*e.tmap)[name]

	if !ok {
		return NewNotFoundError(name)
	}

	if e.globalContext != nil {
		data = e.globalContext.WithData(data)
	}

	if err := t.Execute(w, data); err != nil {
		return NewRenderError(err)
	}

	return nil
}

// Must be executed at least once to load templates, if the template set changes post-hoc
// you must Parse() again to see the changes
// The requirement engine specifies on templates is that each set define a "ROOT" element
func (e *Engine) Parse() error {
	for name, set := range *e.tset {
		t := template.New("")

		if e.fmap != nil {
			t.Funcs(e.fmap)
		}

		if _, err := t.ParseFiles(JoinDir(e.rootDir, set)...); err != nil {
			return err
		}

		t = t.Lookup("ROOT")
		if t == nil {
			return fmt.Errorf("ROOT template not found in %v", set)
		}

		(*e.tmap)[name] = t
	}
	return nil
}
