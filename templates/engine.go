package templates

import (
	"fmt"
	"io"
	"text/template"
)

type (
	// TemplateSet is the type for a set of templates, associated
	// from name to a a slice of file names from which to load the template
	TemplateSet map[Name][]string

	// templateMap is the type for a map of template names
	// to their actual loaded templates (used internally to engine)
	templateMap map[Name]*template.Template

	// Context represents any structure that conserve as the global
	// context for an egnine. The WithData(interface{}) is used
	// so that the context can deal with the data passed into the
	// execution of the template in any way it likes
	Context interface {
		WithData(interface{}) Context
	}

	// An Engine is the primary feature of the templates package,
	// it is considered an immutable structure which can render
	// templates.
	Engine struct {
		// Root directory from which to look up relative file names
		// in all template sets
		rootDir string

		// The user provided TemplateSet definition for templates
		tset *TemplateSet

		// The currently generated templateMap
		tmap *templateMap

		// The map of functions supposed to be applied
		// to each of the functions
		fmap template.FuncMap

		// The context which wraps whatever interface{} the template
		// is executed with
		globalContext Context

		// A flag indicating whether to reload the templates every time
		// they are executed
		everyload bool
	}
)

// NewEngine instantiates a new templates.Engine, rooted at the
// provided rootDir and based on the templateSet definition tset
func NewEngine(rootDir string, tset *TemplateSet) *Engine {
	tm := make(templateMap)

	return &Engine{
		rootDir: rootDir,
		tset:    tset,
		tmap:    &tm,
	}
}

// WithContext constructs a new templates.Engine based on the
// current tempaltes.Engine, e, but with the context provided, c.
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

// WithEveryLoad constructs a new templates.Engine based on the
// current templates.Engine, e, but with everload set to true.
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

// WithFuncMap constructs a new templates.Engine based on the
// current templates.Engine, e, but with the function map fm
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

// Execute executes the template referenced by name, with the global context
// wrapped data, and writes the output to the io.Writer, w. If everyload is
// set to true it first re-parses ALL the templates. if the template
// is not defined, it returns a templates.NotFoundError
// If the globalContext is nil, it doesn't wrap the provided data interface in
// anything. If there is an error while rendering the template, it returns a
// templates.RenderError
// If everyload is true, the client can also recieve any error that the Parse()
// function of this engine may return.
func (e *Engine) Execute(w io.Writer, name Name, data interface{}) error {
	// everyload is like a development mode, always re-compile templates
	if e.everyload {
		err := e.Parse()
		if err != nil {
			return err
		}
	}

	// retrieve template by name
	t, ok := (*e.tmap)[name]

	if !ok { // the name is undefined
		return NewNotFoundError(name)
	}

	// if we have a global context, wrap our data in it
	if e.globalContext != nil {
		data = e.globalContext.WithData(data)
	}

	// execute, writing to w, with data
	if err := t.Execute(w, data); err != nil {
		// execution error
		return NewRenderError(err)
	}

	// success
	return nil
}

// Parse loads the templates defined on e.tset, into the internal template. It must
// be executed at least once to load templates. If the template set changes post-hoc
// which it can (because the client should still hold a pointer to it), the client must
// Parse() again to see the changed.
// The templates package imposed requirement on all templates loaded is that they have a "ROOT"
// element - this handles which template to load without needing to synchronize template.Names with
// {{define "TemplateName"}}s
func (e *Engine) Parse() error {
	// For each of the template sets
	for name, set := range *e.tset {
		t := template.New("")

		if e.fmap != nil {
			t.Funcs(e.fmap)
		}

		if _, err := t.ParseFiles(JoinDir(e.rootDir, set)...); err != nil {
			return err
		}

		// look for the special ROOT template
		t = t.Lookup("ROOT")
		if t == nil {
			return fmt.Errorf("ROOT template not found in %v", set)
		}

		(*e.tmap)[name] = t
	}
	return nil
}
