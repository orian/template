// package template wraps a html/template standard package. It allows

package template

import (
	htmlTmpl "html/template"
	"io"
)

var (
	debug = false
)

type Template interface {
	Name() string
	Execute(io.Writer, interface{}) error
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

// delays parsing until execution
type reloadTemplate struct {
	t *htmlTmpl.Template
	ops []op
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
	t.ops = append(t.ops, &opParseFiles{filenames})
	return t, nil
}

func (t *reloadTemplate) ParseGlob(pattern string) (Template, error) {
	t.ops = append(t.ops, &opParseGlob{pattern})
	return t, nil
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

type instantTemplate struct {
	*htmlTmpl.Template
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

func (r *instantTemplate) Execute(wr io.Writer, data interface{}) error{
	return r.Template.Execute(wr, data)
}

func Must(t Template, err error) Template {
	if err!= nil {
		panic(err)
	}
	return t
}

func New(name string) Template {
	if debug {
		return &reloadTemplate{t: htmlTmpl.New(name)}
	}
	return &instantTemplate{htmlTmpl.New(name)}
}

func ParseFiles(filenames ...string) (Template, error) {
	if debug {
		return (&reloadTemplate{}).ParseFiles(filenames...)
	}
	t,err := htmlTmpl.ParseFiles(filenames...)
	return &instantTemplate{t},err
}

func ParseGlob(pattern string) (Template, error) {
	if debug {
		return (&reloadTemplate{}).ParseGlob(pattern)
	}
	t,err := htmlTmpl.ParseGlob(pattern)
	return &instantTemplate{t},err
}