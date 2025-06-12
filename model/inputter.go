package model

import (
	"fmt"
)

const maxVarPromptWidth = 18

var inputter *Inputter

type Inputter struct {
	readInput func() string
}

func NewInputter(readInput func() string) *Inputter {
	return &Inputter{
		readInput: readInput,
	}
}

func SetInputter(readInput func() string) {
	inputter = NewInputter(readInput)
}

func (i *Inputter) Prompt(required bool, maxNameLen int, varName string) string {
	necessity := "optional"
	if required {
		necessity = "required"
	}

	pad := maxNameLen + 2 // plus double quote
	if pad > maxVarPromptWidth {
		pad = maxVarPromptWidth
	}
	name := fmt.Sprintf(`"%s"`, varName)

	return fmt.Sprintf(`Enter %-*s (%s): `, pad, name, necessity)
}

func (i *Inputter) Input(prompt string) string {
	fmt.Print(prompt)
	return i.readInput()
}
