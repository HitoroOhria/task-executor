package console

type Deps struct {
	Runner  Runner
	Printer Printer
}

func NewDeps(runner Runner, printer Printer) *Deps {
	return &Deps{
		Runner:  runner,
		Printer: printer,
	}
}
