package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func printBox() {
	lines := []string{"foo", "bar", "baz"}

	lineStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("231")) // 白

	var styledLines []string
	for _, line := range lines {
		styledLines = append(styledLines, lineStyle.Render(line))
	}

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(0, 1)

	content := strings.Join(styledLines, "\n")
	box := borderStyle.Render(content)

	fmt.Println(box)
}
