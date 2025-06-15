package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type VarInputter struct {
	Vars     []string
	Cursor   int
	Selected map[int]struct{}
}

func initVarInputter() VarInputter {
	return VarInputter{
		Vars:     []string{"OPTIONAL", "REQUIRED", "HAS_DEFAULT"},
		Cursor:   0,
		Selected: make(map[int]struct{}),
	}
}

func (m VarInputter) Init() tea.Cmd {
	return nil
}

func (m VarInputter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Vars)-1 {
				m.Cursor++
			}
		case "enter", " ":
			_, ok := m.Selected[m.Cursor]
			if ok {
				delete(m.Selected, m.Cursor)
			} else {
				m.Selected[m.Cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m VarInputter) View() string {
	// The header
	s := "What should we buy at the market?\n\n"

	for i, v := range m.Vars {
		cursor := " "
		if m.Cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.Selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, v)
	}

	s += "\nPress q to quit.\n"

	return s
}
