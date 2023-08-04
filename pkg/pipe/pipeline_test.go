package pipe

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/jensneuse/pipeline/pkg/step"
)

func TestPipeline_Run_Simple(t *testing.T) {
	outBody := `{"foo":"bar"}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(outBody))
	}))

	defer srv.Close()

	in := step.HttpStepInput{
		URL: srv.URL,
	}

	jsonFile, err := os.Open("./testdata/simple.json")
	if err != nil {
		t.Fatal(err)
	}
	defer jsonFile.Close()

	var pipeline Pipeline
	err = pipeline.FromConfig(jsonFile)
	if err != nil {
		t.Fatal(err)
	}

	input := bytes.Buffer{}
	output := bytes.Buffer{}

	err = json.NewEncoder(&input).Encode(in)
	if err != nil {
		t.Fatal(err)
	}

	err = pipeline.Run(&input, &output)
	if err != nil {
		t.Fatal(err)
	}

	got := output.String()
	want := `{"result":"bar"}`

	if want != got {
		t.Fatalf("want: %s, got: %s\n", got, want)
	}
}

func TestPipeline_Run_Complex(t *testing.T) {
	srvBody := `
		{
			"policies": {
				"1": {
					"id": 1,
					"name": "pol1"
				},
				"2": {
					"id": 2,
					"name": "pol2"
				}
			}
		}
	`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		got := string(body)
		want := `{"policies":{"1":{"id":1,"name":"pol1"},"2":{"id":2,"name":"pol2"}}}
`
		if want != got {
			t.Fatalf("want: %s\ngot: %s\n", want, got)
		}
		_, _ = w.Write([]byte(srvBody))
	}))

	defer srv.Close()

	in := step.HttpStepInput{
		URL:    srv.URL,
		Method: http.MethodPost,
		Body: map[string]interface{}{
			"policies": []map[string]interface{}{
				{
					"id":   1,
					"name": "pol1",
				},
				{
					"id":   2,
					"name": "pol2",
				},
			},
		},
	}

	jsonFile, err := os.Open("./testdata/complex.json")
	if err != nil {
		t.Fatal(err)
	}
	defer jsonFile.Close()

	var pipeline Pipeline
	err = pipeline.FromConfig(jsonFile)
	if err != nil {
		t.Fatal(err)
	}

	input := bytes.Buffer{}
	output := bytes.Buffer{}

	err = json.NewEncoder(&input).Encode(in)
	if err != nil {
		t.Fatal(err)
	}

	err = pipeline.Run(&input, &output)
	if err != nil {
		t.Fatal(err)
	}

	got := output.String()
	want := `{
  "policies": [
	{
	  "id": 1,
	  "name": "pol1"
	},
	{
	  "id": 2,
	  "name": "pol2"
	}
  ]
}`

	want = pretty(want)

	if want != got {
		t.Fatalf("want: %s, got: %s\n", want, got)
	}
}

func pretty(input string) string {
	data := map[string]interface{}{}
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		panic(err)
	}
	tmpl, err := template.New("").Funcs(sprig.TxtFuncMap()).Parse("{{ toPrettyJson . }}")
	if err != nil {
		panic(err)
	}
	out := bytes.Buffer{}
	err = tmpl.Execute(&out, data)
	if err != nil {
		panic(err)
	}
	return out.String()
}

func BenchmarkPipeline_Run_Complex(b *testing.B) {
	srvBody := `
		{
			"policies": {
				"1": {
					"id": 1,
					"name": "pol1"
				},
				"2": {
					"id": 2,
					"name": "pol2"
				}
			}
		}
	`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		got := string(body)
		want := `{"policies":{"1":{"id":1,"name":"pol1"},"2":{"id":2,"name":"pol2"}}}
`
		if want != got {
			b.Fatalf("want: %s\ngot: %s\n", want, got)
		}
		_, _ = w.Write([]byte(srvBody))
	}))

	defer srv.Close()

	in := step.HttpStepInput{
		URL:    srv.URL,
		Method: http.MethodPost,
		Body: map[string]interface{}{
			"policies": []map[string]interface{}{
				{
					"id":   1,
					"name": "pol1",
				},
				{
					"id":   2,
					"name": "pol2",
				},
			},
		},
	}

	jsonFile, err := os.Open("./testdata/complex.json")
	if err != nil {
		b.Fatal(err)
	}
	defer jsonFile.Close()

	var pipeline Pipeline
	err = pipeline.FromConfig(jsonFile)
	if err != nil {
		b.Fatal(err)
	}

	input := bytes.Buffer{}
	output := bytes.Buffer{}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		output.Reset()
		input.Reset()

		err = json.NewEncoder(&input).Encode(in)
		if err != nil {
			b.Fatal(err)
		}

		err = pipeline.Run(&input, &output)
		if err != nil {
			b.Fatal(err)
		}
	}
}
