package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	target := flag.String("target", "varInputter", "select target")
	flag.Parse()

	switch *target {
	case "tutorial1":
		runTutorial1()
	case "tutorial2":
		runTutorial2()
	case "textInput":
		runTextInput()
	case "textInputs":
		runTextInputs()
	case "table":
		runTable()
	case "varInputter":
		runVarInputter()
	default:
		fmt.Println("invalid target")
	}
}

func runTutorial1() {
	p := tea.NewProgram(initialModel1())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func runTutorial2() {
	if _, err := tea.NewProgram(model2{}).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}

func runTextInput() {
	p := tea.NewProgram(initialTextInput())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func runTextInputs() {
	if _, err := tea.NewProgram(initialTextInputs()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}

func runTable() {
	m := initTable()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}

func runVarInputter() {
	p := tea.NewProgram(initVarInputter())
	result, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	vi, ok := result.(VarInputter)
	if ok {
		fmt.Printf("result: %v", vi.GetValues())
	}
}
