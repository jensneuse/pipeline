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
	tmpl *template.Template
}

func UnmarshalJsonStep(reader io.Reader) (JsonStep,error){
	var step JsonStep
	step.tmpl = template.New("")
	step.tmpl.Funcs(sprig.TxtFuncMap())
	err := json.NewDecoder(reader).Decode(&step)
	if err != nil {
		return step,err
	}
	step.tmpl,err = step.tmpl.Parse(step.Template)
	return step,err
}

func (j JsonStep) Invoke(reader io.Reader, writer io.Writer) error {
	data,err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	var in interface{}

	err = json.Unmarshal(data,&in)
	if err != nil {
		return err
	}

	return j.tmpl.Execute(writer,in)
}
