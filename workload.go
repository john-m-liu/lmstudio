package main

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
)

var (
	AddJobType         = "ADD"
	SubtractJobType    = "SUBTRACT"
	PrintStringJobType = "PRINT_STRING"
	ShellJobType       = "SHELL"
	ErrorJobType       = "ERROR"
)

type Workload struct {
	WorkloadType string `json:"workload_type"`
	Body         any    `json:"body"`
}

type AddNumbersBody struct {
	Int1 int `json:"int1" mapstructure:"int1"`
	Int2 int `json:"int2" mapstructure:"int2"`
}

type SubtractNumbersBody struct {
	Int1 int `json:"int1" mapstructure:"int1"`
	Int2 int `json:"int2" mapstructure:"int2"`
}

type PrintBody struct {
	ToPrint string `json:"to_print" mapstructure:"to_print"`
}

type ShellBody struct {
	ToExecute string `json:"to_execute" mapstructure:"to_execute"`
}

func (w *Workload) UnmarshalJSON(data []byte) error {
	v := map[string]any{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	w.WorkloadType = v["workload_type"].(string)
	switch w.WorkloadType {
	case AddJobType:
		body := AddNumbersBody{}
		err := mapstructure.Decode(v["body"], &body)
		if err != nil {
			return err
		}
		w.Body = body
	case SubtractJobType:
		body := SubtractNumbersBody{}
		err := mapstructure.Decode(v["body"], &body)
		if err != nil {
			return err
		}
		w.Body = body
	case PrintStringJobType:
		body := PrintBody{}
		err := mapstructure.Decode(v["body"], &body)
		if err != nil {
			return err
		}
		w.Body = body
	case ShellJobType:
		body := ShellBody{}
		err := mapstructure.Decode(v["body"], &body)
		if err != nil {
			return err
		}
		w.Body = body
	case ErrorJobType:
		// the error job type has no body
	}

	return nil
}
