package console

type Deps struct {
	Command Command
	Printer Printer
}

func NewDeps(cmd Command, printer Printer) *Deps {
	return &Deps{
		Command: cmd,
		Printer: printer,
	}
}
