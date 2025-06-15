package adapter

import (
	"fmt"

	"github.com/HitoroOhria/task-executor/domain/console"
	"github.com/HitoroOhria/task-executor/io"
)

type VariableInputterImpl struct {
}

func NewVariableInputter() console.VariableInputter {
	return &VariableInputterImpl{}
}

func (v *VariableInputterImpl) Input(vars []*console.Variable) ([]*console.Variable, error) {
	inputs := make([]io.Variable, 0, len(vars))
	for _, v := range vars {
		inputs = append(inputs, io.NewVariable(v.Name, v.Required, v.DefaultValue))
	}

	result, err := io.RunVariableInputTable(inputs)
	if err != nil {
		return nil, fmt.Errorf("io.RunVariableInputTable: %w", err)
	}

	res := make([]*console.Variable, 0, len(vars))
	for _, v := range result {
		res = append(res, &console.Variable{
			Name:         v.Name,
			Required:     v.Required,
			DefaultValue: v.DefaultValue,
			InputValue:   v.InputValue,
		})
	}

	return res, nil
}
