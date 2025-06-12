package model

import (
	"fmt"
	"regexp"
)

type VarValue string

func NewVarValue(value any) VarValue {
	v := ""
	switch value.(type) {
	case string:
		v = value.(string)
	case *string:
		v = *value.(*string)
	default:
		fmt.Printf("[🚨WARNING!!] unknown var value type: type = %T. value = %+v\n", value, value)
	}

	return VarValue(v)
}

func (v VarValue) IsSelfValue(name string) bool {
	return string(v) == v.selfValue(name)
}

func (v VarValue) selfValue(name string) string {
	return fmt.Sprintf("{{.%s}}", name)
}

func (v VarValue) IsSelfValueWithDefault(name string) bool {
	re := v.selfWithDefaultRegex(name)
	return re.MatchString(string(v))
}

func (v VarValue) selfWithDefaultRegex(name string) *regexp.Regexp {
	re := fmt.Sprintf(`%s ?| default ".+"`, v.selfValue(name))
	return regexp.MustCompile(re)
}
