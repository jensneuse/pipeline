package pipeline

import (
	"bytes"
	"encoding/json"
	"github.com/jensneuse/pipeline/pkg/pipeline/step"
	"io"
)

type Step interface {
	Invoke(reader io.Reader, writer io.Writer) error
}

type Config struct {
	Steps []StepConfig `json:"steps"`
}

type StepConfig struct {
	Kind   string          `json:"kind"`
	Config json.RawMessage `json:"config"`
}

type Pipeline struct {
	Steps []Step
}

func (p *Pipeline) FromConfig(reader io.Reader) error {
	var config Config
	err := json.NewDecoder(reader).Decode(&config)
	if err != nil {
		return err
	}
	for i := range config.Steps {
		var next Step
		switch config.Steps[i].Kind {
		case "JSON":
			next, err = step.UnmarshalJsonStep(bytes.NewReader(config.Steps[i].Config))
		case "HTTP":
			next, err = step.UnmarshalHttpStep(bytes.NewReader(config.Steps[i].Config))
		}
		if err != nil {
			return err
		}
		p.Steps = append(p.Steps, next)
	}
	return nil
}

func (p *Pipeline) Run (reader io.Reader,writer io.Writer) error {

	readBuf := bytes.Buffer{}
	writeBuf := bytes.Buffer{}

	_,err := readBuf.ReadFrom(reader)
	if err != nil {
		return err
	}

	for i := range p.Steps {
		err = p.Steps[i].Invoke(&readBuf,&writeBuf)
		if err != nil {
			return err
		}
		readBuf.Reset()
		_,err = writeBuf.WriteTo(&readBuf)
		if err != nil {
			return err
		}
		writeBuf.Reset()
		intermediate := readBuf.String()
		if intermediate == "" {

		}
	}

	_,err = readBuf.WriteTo(writer)
	return err
}