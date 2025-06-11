package main

import "fmt"

// Vars は変数名と値のセット
type Vars []*Var

type Var struct {
	Required bool
	Name     string
	Value    string
}

func (vs *Vars) FindByName(name string) *Var {
	for _, v := range *vs {
		if v.Name == name {
			return v
		}
	}

	return nil
}

func (vs *Vars) SetNameAsOptional(name string) {
	*vs = append(*vs, &Var{
		Required: false,
		Name:     name,
		Value:    "",
	})
}

func (vs *Vars) SetNameAsRequired(name string) {
	*vs = append(*vs, &Var{
		Required: true,
		Name:     name,
		Value:    "",
	})
}

func (vs *Vars) SetValue(name, value string) error {
	v := vs.FindByName(name)
	if v == nil {
		return fmt.Errorf("variable %s is not found. vs = %+v", name, *vs)
	}

	if v.Required {
		if value == "" {
			return fmt.Errorf("variable %s is required", name)
		}
	}

	v.Value = value
	return nil
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
