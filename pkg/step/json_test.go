package step

import (
	"bytes"
	"strings"
	"testing"
)

func TestJsonStep_Invoke(t *testing.T) {

	input := `{"foo":"bar"}`
	config := `{"template":"{\"baz\":\"{{ .foo }}\"}"}`

	step,err := UnmarshalJsonStep(strings.NewReader(config))
	if err != nil {
		t.Fatal(err)
	}

	buf := bytes.Buffer{}
	err = step.Invoke(strings.NewReader(input),&buf)
	if err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	want := `{"baz":"bar"}`
	if want != got {
		t.Fatalf("want: %s\ngot: %s\n",want,got)
	}
}

func TestJsonStep_Invoke_NewJSON(t *testing.T) {

	input := `{"foo":"bar"}`
	tmpl := "{\"baz\":\"{{ .foo }}\"}"

	step,err := NewJSON(tmpl)
	if err != nil {
		t.Fatal(err)
	}

	buf := bytes.Buffer{}
	err = step.Invoke(strings.NewReader(input),&buf)
	if err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	want := `{"baz":"bar"}`
	if want != got {
		t.Fatalf("want: %s\ngot: %s\n",want,got)
	}
}