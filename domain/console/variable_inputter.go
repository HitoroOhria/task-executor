package console

type VariableInputter interface {
	Input(vars []*Variable) ([]*Variable, error)
}

type Variable struct {
	Name         string
	Required     bool
	DefaultValue string
	InputValue   string
}

func NewVariable(name string, required bool, defaultValue string) *Variable {
	return &Variable{
		Name:         name,
		Required:     required,
		DefaultValue: defaultValue,
		InputValue:   "",
	}
}
