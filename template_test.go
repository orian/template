package template

import ("testing"
	"bytes"
)

func TestDelayParsing(t *testing.T) {
	tmplSrc := `{{define "test"}}test{{end}}`
	tmpl,err := New("test").Parse(tmplSrc)
	if err != nil {
		t.Errorf("want nil, got: %s", err)
	}
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, nil); err != nil {
		t.Errorf("want nil, got: %s", err)
	}
	if s:=buf.String(); "test" != s {
		t.Errorf("want test, got %q", s)
	}
}

func TestDelayParsingErr(t *testing.T) {
	Debug = true
	tmplSrc := `{{define "test"}}test{end}}`
	tmpl,err := New("test").Parse(tmplSrc)
	if err != nil {
		t.Fatalf("want nil, got: %s", err)
	}
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, nil); err == nil {
		t.Error("want err, got: nil")
		t.Errorf("want nil, got %q", buf.String())
	}
}

func TestParsingErr(t *testing.T) {
	Debug = false
	tmplSrc := `{{define "test"}}test{end}}`
	_,err := New("test").Parse(tmplSrc)
	if err == nil {
		t.Fatalf("want err, got: nil", err)
	}
}