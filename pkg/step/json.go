package step

import (
	"encoding/json"
	"github.com/Masterminds/sprig"
	"io"
	"io/ioutil"
	"text/template"
)

type JsonStep struct {
	Template string `json:"template"`
}

func UnmarshalJsonStep(reader io.Reader) (JsonStep,error){
	var step JsonStep
	err := json.NewDecoder(reader).Decode(&step)
	return step,err
}

func (j JsonStep) Invoke(reader io.Reader, writer io.Writer) error {
	data,err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	in := map[string]interface{}{}

	err = json.Unmarshal(data,&in)
	if err != nil {
		return err
	}

	tmpl,err := template.New("").Funcs(sprig.TxtFuncMap()).Parse(j.Template)
	if err != nil {
		return err
	}

	return tmpl.Execute(writer,in)
}
