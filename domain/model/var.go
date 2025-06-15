package model

import "fmt"

func makeArg(name, value string) string {
	return fmt.Sprintf(`%s=%s`, name, value)
}
