package adapter

import (
	"fmt"
	"strings"

	"github.com/HitoroOhria/task-executor/domain/console"
	"github.com/HitoroOhria/task-executor/domain/value"
	"github.com/charmbracelet/lipgloss"
)

const (
	requiredHeader = "--- required ---"
	optionalHeader = "--- optional ---"
	endLine        = "---   end   ---"
)

type PrinterImpl struct{}

func NewPrinter() console.Printer {
	return &PrinterImpl{}
}

func (p *PrinterImpl) LineBreaks() {
	fmt.Println()
}

func (p *PrinterImpl) ExecutionTask(taskfile string, fullName value.FullTaskName, args ...string) {
	command := fmt.Sprintf("task -t %s %s", taskfile, fullName)
	if len(args) != 0 {
		command += ` \`
	}
	run := "[run]"
	lines := append([]string{run, command}, makeFormattedArgs(args)...)

	lineStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("231")) // 白

	var styledLines []string
	for _, line := range lines {
		styledLines = append(styledLines, lineStyle.Render(line))
	}

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(0, 3, 0, 1)

	content := strings.Join(styledLines, "\n")
	box := borderStyle.Render(content)

	fmt.Println(box)
}

func makeFormattedArgs(args []string) []string {
	formattedArgs := make([]string, 0, len(args))
	for i, arg := range args {
		formatted := fmt.Sprintf("    %s", arg)
		if i != len(args)-1 {
			formatted += ` \`
		}
		formattedArgs = append(formattedArgs, formatted)
	}

	return formattedArgs
}
