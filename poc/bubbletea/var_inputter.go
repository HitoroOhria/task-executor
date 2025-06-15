package main

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	ellipsis = "…"

	nameHeader         = "Variable"
	requiredHeader     = "Req."
	defaultValueHeader = "Def."
	inputValueHeader   = "Value"

	maxNameWidth         = 18
	maxDefaultValueWidth = 18
	inputValueWidth      = 20
)

var ellipsisLen = utf8.RuneCountInString(ellipsis)

type VarInputter struct {
	textInputs []textinput.Model
	cursor     int // number of row index

	Vars []Var
}

type Var struct {
	Name         string
	Required     bool
	DefaultValue string
	InputValue   string
}

func initVarInputter() VarInputter {
	vars := []Var{
		{"OPTIONAL", false, "", ""},
		{"REQUIRED", true, "", ""},
		{"HAS_DEFAULT", false, "foo", ""},
		{"LOOOOOOOOOOOOOOOOOOOOOOOOOOOOONG", false, "Lorem Ipsum is simply dummy text of the printing and typesetting industry", ""},
	}

	tis := make([]textinput.Model, len(vars))
	for i, _ := range vars {
		ti := textinput.New()
		ti.Prompt = ""
		if i == 0 {
			ti.Focus()
		}
		ti.Width = inputValueWidth

		tis[i] = ti
	}

	return VarInputter{
		textInputs: tis,
		Vars:       vars,
		cursor:     0,
	}
}

func truncate(s string, max int) string {
	if utf8.RuneCountInString(s) <= max {
		return s
	}

	runes := []rune(s)
	// maxが小さすぎると「...」を入れる余地がない
	if max <= ellipsisLen {
		return string(runes[:max])
	}

	return string(runes[:max-ellipsisLen]) + ellipsis
}

func (m VarInputter) Separator() string {
	width := m.NameColumnLen() + 2 +
		m.RequiredColumnLen() + 2 +
		m.DefaultValueColumnLen() + 2 +
		inputValueWidth

	return strings.Repeat("─", width)
}

func (m VarInputter) NameColumnLen() int {
	var maxLen int
	for _, v := range m.Vars {
		if len(v.Name) > maxLen {
			maxLen = len(v.Name)
		}
	}
	if maxLen > maxNameWidth {
		maxLen = maxNameWidth
	}

	return max(maxLen, len(nameHeader))
}

func (m VarInputter) RequiredColumnLen() int {
	return len(requiredHeader)
}

func (m VarInputter) DefaultValueColumnLen() int {
	var maxLen int
	for _, v := range m.Vars {
		if len(v.DefaultValue) > maxLen {
			maxLen = len(v.DefaultValue)
		}
	}

	if maxLen > maxDefaultValueWidth {
		maxLen = maxDefaultValueWidth
	}

	return max(maxLen, len(defaultValueHeader)) + 2 // add "[]"
}

func (m VarInputter) Init() tea.Cmd {
	return textinput.Blink
}

func (m VarInputter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyUp:
			if m.cursor > 0 {
				m.textInputs[m.cursor].Blur()
				m.cursor--
				m.textInputs[m.cursor].Focus()
			}
		case tea.KeyDown:
			if m.cursor < len(m.Vars)-1 {
				m.textInputs[m.cursor].Blur()
				m.cursor++
				m.textInputs[m.cursor].Focus()
			}
		default:
			for i, ti := range m.textInputs {
				if !ti.Focused() {
					continue
				}

				m.textInputs[i], cmd = ti.Update(msg)
			}
		}
	}

	for i, ti := range m.textInputs {
		m.Vars[i].InputValue = ti.Value()
	}

	return m, cmd
}

func (m VarInputter) View() string {
	// The header
	s := "Input variable values.\n\n"

	s += fmt.Sprintf(
		"%-*s  %s  %-*s  %-*s\n",
		m.NameColumnLen(),
		nameHeader,
		requiredHeader,
		m.DefaultValueColumnLen(),
		defaultValueHeader,
		inputValueWidth,
		inputValueHeader,
	)

	s += m.Separator() + "\n"

	for i, v := range m.Vars {
		name := truncate(v.Name, maxNameWidth)

		required := ""
		if v.Required {
			required = " ✓ "
		}

		defaultValue := v.DefaultValue
		if defaultValue != "" {
			truncated := truncate(defaultValue, maxDefaultValueWidth)
			defaultValue = fmt.Sprintf("[%s]", truncated)
		}

		s += fmt.Sprintf(
			"%-*s  %-*s  %-*s  %s\n",
			m.NameColumnLen(),
			name,
			m.RequiredColumnLen(),
			required,
			m.DefaultValueColumnLen(),
			defaultValue,
			m.textInputs[i].View(),
		)
	}

	s += "\n(enter to finish)\n\n"

	return s
}

func (m VarInputter) GetValues() []string {
	values := make([]string, len(m.Vars))
	for i, v := range m.Vars {
		if v.InputValue == "" {
			values[i] = v.DefaultValue
			continue
		}

		values[i] = v.InputValue
	}

	return values
}
