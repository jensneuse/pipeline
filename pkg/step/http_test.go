package step

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHttpStep_Invoke(t *testing.T) {

	outBody := `{"foo":"bar"}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(outBody))
	}))

	defer srv.Close()

	in := HttpStepInput{
		Method: http.MethodGet,
		URL:    srv.URL,
	}

	inBytes,err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}

	step := HttpStep{
		DefaultTimeout: time.Second * time.Duration(10),
		DefaultMethod:  http.MethodGet,
	}

	out := bytes.Buffer{}

	err = step.Invoke(bytes.NewReader(inBytes),&out)
	if err != nil {
		t.Fatal(err)
	}

	var output HttpStepOutput
	err = json.Unmarshal(out.Bytes(),&output)
	if err != nil {
		t.Fatal(err)
	}

	if output.StatusCode != http.StatusOK {
		t.Fatal("want 200 OK")
	}
	if output.Body != outBody {
		t.Fatalf("want: %s, got: %s\n",outBody,output.Body)
	}
}
