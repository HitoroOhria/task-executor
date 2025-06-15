package model

import "fmt"

const maxVarPromptWidth = 20

// Deprecated: Prompt は入力のプロンプト
// 変数の値を入力する時に利用される
type Prompt struct {
	VarName string
}

func NewPrompt(varName string) *Prompt {
	return &Prompt{
		VarName: varName,
	}
}

func (p *Prompt) Generate(maxDisplayLen int, defaultValue string) string {
	pad := maxDisplayLen
	if pad > maxVarPromptWidth {
		pad = maxVarPromptWidth
	}

	varDisplay := generateVarDisplay(p.VarName, defaultValue)

	return fmt.Sprintf(`Enter %-*s: `, pad, varDisplay)
}

func generateVarDisplay(varName string, defaultValue string) string {
	if defaultValue == "" {
		return fmt.Sprintf(`"%s" `, varName)
	}

	return fmt.Sprintf(`"%s" [%s]`, varName, defaultValue)
}
