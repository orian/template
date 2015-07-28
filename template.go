// package template wraps a html/template standard package. It allows

package template

import (
	htmlTmpl "html/template"
	"io"
)

var (
	Debug = false
)

type Template interface {
	Name() string
	Execute(io.Writer, interface{}) error
	Funcs(funcMap htmlTmpl.FuncMap) Template
	Parse(string) (Template, error)
	ParseFiles(filenames ...string) (Template, error)
	ParseGlob(pattern string) (Template, error)
}

type op interface {
	Run(*htmlTmpl.Template) (*htmlTmpl.Template, error)
}

type opParse struct {
	t string
}

func (o *opParse) Run(t *htmlTmpl.Template) (*htmlTmpl.Template, error) {
	return t.Parse(o.t)
}

type opParseFiles struct {
	files []string
}

func (o *opParseFiles) Run(t *htmlTmpl.Template) (*htmlTmpl.Template, error) {
	if t == nil {
		return htmlTmpl.ParseFiles(o.files...)
	}
	return t.ParseFiles(o.files...)
}

type opParseGlob struct{
	pattern string
}

func (o *opParseGlob) Run(t *htmlTmpl.Template) (*htmlTmpl.Template, error) {
	if t == nil {
		return htmlTmpl.ParseGlob(o.pattern)
	}
	return t.ParseGlob(o.pattern)
}

type opFuncs struct {
	funcs htmlTmpl.FuncMap
}

func copyFuncMap(f htmlTmpl.FuncMap) htmlTmpl.FuncMap {
	fm := htmlTmpl.FuncMap{}
	for k,v := range f {
		fm[k] = v
	}
	return fm
}

func (o *opFuncs) Run(t *htmlTmpl.Template) (*htmlTmpl.Template, error) {
	return t.Funcs(o.funcs), nil
}

// delays parsing until execution
type reloadTemplate struct {
	t *htmlTmpl.Template
	ops []op
}

func (t *reloadTemplate) Execute(wr io.Writer, data interface{}) error{
	tmpl := t.t
	var err error
	for _,o := range t.ops {
		tmpl, err = o.Run(tmpl)
		if err != nil {
			return err
		}
	}
	return tmpl.Execute(wr, data)
}

func(r *reloadTemplate) Funcs(funcMap htmlTmpl.FuncMap) Template {
	r.ops = append(r.ops, &opFuncs{copyFuncMap(funcMap)})
	return r
}

func (r *reloadTemplate) Name() string{
	// if not created with new, than reload
	return r.t.Name()
}

func (t *reloadTemplate) Parse(src string) (Template, error) {
	t.ops = append(t.ops, &opParse{src})
	return t, nil
}

func (t *reloadTemplate) ParseFiles(filenames ...string) (Template, error) {
	c := make([]string, len(filenames))
	t.ops = append(t.ops, &opParseFiles{c})
	return t, nil
}

func (t *reloadTemplate) ParseGlob(pattern string) (Template, error) {
	t.ops = append(t.ops, &opParseGlob{pattern})
	return t, nil
}


type instantTemplate struct {
	*htmlTmpl.Template
}

func (r *instantTemplate) Execute(wr io.Writer, data interface{}) error{
	return r.Template.Execute(wr, data)
}

func(r *instantTemplate) Funcs(funcMap htmlTmpl.FuncMap) Template {
	// TODO creating new instantTemplate is a problem
	return &instantTemplate{r.Template.Funcs(funcMap)}
}

func (r *instantTemplate) Name() string {
	return r.Template.Name()
}

func (r *instantTemplate) Parse(src string) (Template, error) {
	t,err := r.Template.Parse(src)
	return &instantTemplate{t},err
}

func (r *instantTemplate) ParseFiles(filenames ...string) (Template, error) {
	t,err := r.Template.ParseFiles(filenames...)
	return &instantTemplate{t},err
}

func (r *instantTemplate) ParseGlob(pattern string) (Template, error) {
	t,err := r.Template.Parse(pattern)
	return &instantTemplate{t},err
}

func Must(t Template, err error) Template {
	if err!= nil {
		panic(err)
	}
	return t
}

func New(name string) Template {
	if Debug {
		return &reloadTemplate{t: htmlTmpl.New(name)}
	}
	return &instantTemplate{htmlTmpl.New(name)}
}

func ParseFiles(filenames ...string) (Template, error) {
	if Debug {
		return (&reloadTemplate{}).ParseFiles(filenames...)
	}
	t,err := htmlTmpl.ParseFiles(filenames...)
	return &instantTemplate{t},err
}

func ParseGlob(pattern string) (Template, error) {
	if Debug {
		return (&reloadTemplate{}).ParseGlob(pattern)
	}
	t,err := htmlTmpl.ParseGlob(pattern)
	return &instantTemplate{t},err
}