package console

type Deps struct {
	Runner           Runner
	Printer          Printer
	VariableInputter VariableInputter
}

func NewDeps(runner Runner, printer Printer, variableInputter VariableInputter) *Deps {
	return &Deps{
		Runner:           runner,
		Printer:          printer,
		VariableInputter: variableInputter,
	}
}
