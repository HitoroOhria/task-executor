package adapter

import (
	"fmt"
	"strings"

	"github.com/HitoroOhria/task-executer/domain/model"
	"github.com/HitoroOhria/task-executer/domain/value"
	"github.com/fatih/color"
)

const (
	requiredHeader = "--- required ---"
	optionalHeader = "--- optional ---"
	endLine        = "---   end   ---"
)

type PrinterImpl struct{}

func NewPrinter() model.Printer {
	return &PrinterImpl{}
}

func (p *PrinterImpl) RequiredHeader() {
	color.Red(requiredHeader)
}

func (p *PrinterImpl) OptionalHeader() {
	color.Yellow(optionalHeader)
}

func (p *PrinterImpl) EndLine() {
	color.White(endLine)
}

func (p *PrinterImpl) LineBreaks() {
	fmt.Println()
}

func (p *PrinterImpl) ExecutionTask(taskfile string, fullName value.FullTaskName, args ...string) {
	color.White("run: task -t %s %s %s\n", taskfile, fullName, strings.Join(args, " "))
}
