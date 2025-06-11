package main

import "fmt"

// Vars は変数名と値のセット
type Vars []*Var

func (vs *Vars) FindByName(name string) *Var {
	for _, v := range *vs {
		if v.Name == name {
			return v
		}
	}

	return nil
}

func (vs *Vars) Append(v *Var) {
	*vs = append(*vs, v)
}

func (vs *Vars) SetValue(name, value string) error {
	v := vs.FindByName(name)
	if v == nil {
		return fmt.Errorf("variable %s is not found. vs = %+v", name, *vs)
	}

	return v.SetValue(value)
}

func (vs *Vars) GetMaxNameLen() int {
	maxLen := 0
	for _, v := range *vs {
		if len(v.Name) > maxLen {
			maxLen = len(v.Name)
		}
	}

	return maxLen
}

func (vs *Vars) GetInputPrompt(name string) (string, error) {
	v := vs.FindByName(name)
	if v == nil {
		return "", fmt.Errorf("variable %s is not found. vs = %+v", name, *vs)
	}

	maxLen := vs.GetMaxNameLen()
	return v.GetInputPrompt(maxLen), nil
}

// CommandArgs はコマンドの引数を組み立てる
// e.g. { "NAME": "john", "age": "25" } => [NAME="john", age="25"]
func (vs *Vars) CommandArgs() []string {
	args := make([]string, 0, len(*vs))
	for _, v := range *vs {
		arg := fmt.Sprintf(`%s="%s"`, v.Name, v.Value)
		args = append(args, arg)
	}

	return args
}

type Var struct {
	Required bool
	Name     string
	Value    string
}

func CreateRequiredVar(name string) *Var {
	return &Var{
		Required: true,
		Name:     name,
		Value:    "",
	}
}

func CreateOptionalVar(name string) *Var {
	return &Var{
		Required: false,
		Name:     name,
		Value:    "",
	}
}

func (v *Var) SetValue(value string) error {
	if v.Required {
		if value == "" {
			return fmt.Errorf("variable %s is required", v.Name)
		}
	}

	v.Value = value
	return nil
}

func (v *Var) GetInputPrompt(padding int) string {
	necessity := "optional"
	if v.Required {
		necessity = "required"
	}

	pad := padding + 2 // plus double quote
	varName := fmt.Sprintf(`"%s"`, v.Name)

	return fmt.Sprintf(`Enter %-*s (%s): `, pad, varName, necessity)

}
