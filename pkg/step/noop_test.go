package step

import (
	"bytes"
	"testing"
)

func TestNoOpStep_Invoke(t *testing.T) {
	step := NoOpStep{}
	in := []byte("noop")
	out := bytes.Buffer{}
	err := step.Invoke(bytes.NewReader(in),&out)
	if err != nil {
		t.Fatal(err)
	}
	if out.String() != "noop" {
		t.Fatalf("want: noop, got: %s\n",out.String())
	}
}
