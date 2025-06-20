package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func printTutorial() {
	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		PaddingTop(2).
		PaddingLeft(4).
		Width(22)

	fmt.Println(style.Render("Hello, kitty"))
}
