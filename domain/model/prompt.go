package model

import "fmt"

const maxVarPromptWidth = 18

// Prompt は入力のプロンプト
// 変数の値を入力する時に利用される
type Prompt struct {
	VarName string
}

func NewPrompt(varName string) *Prompt {
	return &Prompt{
		VarName: varName,
	}
}

func (p *Prompt) Generate(maxNameLen int) string {
	pad := maxNameLen + 2 // plus double quote
	if pad > maxVarPromptWidth {
		pad = maxVarPromptWidth
	}
	name := fmt.Sprintf(`"%s"`, p.VarName)

	return fmt.Sprintf(`Enter %-*s: `, pad, name)
}
