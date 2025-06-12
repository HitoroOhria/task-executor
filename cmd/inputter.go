package main

import "fmt"

type Inputter struct {
}

func NewInputter() *Inputter {
	return &Inputter{}
}

func (i *Inputter) Prompt(required bool, maxNameLen int, varName string) string {
	necessity := "optional"
	if required {
		necessity = "required"
	}

	pad := maxNameLen + 2 // plus double quote
	name := fmt.Sprintf(`"%s"`, varName)

	return fmt.Sprintf(`Enter %-*s (%s): `, pad, name, necessity)
}

func (i *Inputter) Input(prompt string) string {
	fmt.Print(prompt)
	return readInput()
}
