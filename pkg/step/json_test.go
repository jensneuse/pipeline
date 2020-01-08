package step

import (
	"bytes"
	"strings"
	"testing"
)

func TestJsonStep_Invoke(t *testing.T) {
	input := `{"foo":"bar"}`

	step := JsonStep{
		Template: `{"baz":"{{ .foo }}"}`,
	}

	buf := bytes.Buffer{}
	err := step.Invoke(strings.NewReader(input),&buf)
	if err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	want := `{"baz":"bar"}`
	if want != got {
		t.Fatalf("want: %s\ngot: %s\n",want,got)
	}
}
