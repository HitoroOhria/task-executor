package io

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

	headerMessage = "Input variable values."
	footerMessage = "(enter to finish)"
)

var ellipsisLen = utf8.RuneCountInString(ellipsis)

// VariableInputTable は変数入力テーブル
type VariableInputTable struct {
	textInputs []textinput.Model
	cursor     int // number of row index

	Variables []Variable
}

type Variable struct {
	Name         string
	Required     bool
	DefaultValue string
	InputValue   string
}

func NewVariable(name string, required bool, defaultValue string) Variable {
	return Variable{
		Name:         name,
		Required:     required,
		DefaultValue: defaultValue,
		InputValue:   "",
	}
}

func newTextInputModel(isFocus bool) textinput.Model {
	ti := textinput.New()
	ti.Prompt = ""
	ti.Width = inputValueWidth

	if isFocus {
		ti.Focus()
	}

	return ti
}

func initVariableInputTable(vars []Variable) VariableInputTable {
	tis := make([]textinput.Model, len(vars))
	for i, _ := range vars {
		tis[i] = newTextInputModel(i == 0) // default focus is first variable
	}

	return VariableInputTable{
		textInputs: tis,
		Variables:  vars,
		cursor:     0,
	}
}

// RunVariableInputTable は変数入力テーブルを実行し、値を受け取る
func RunVariableInputTable(vars []Variable) ([]Variable, error) {
	p := tea.NewProgram(initVariableInputTable(vars))
	result, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("tea.NewProgram.Run: %w", err)
	}

	vit, ok := result.(VariableInputTable)
	if !ok {
		return nil, fmt.Errorf("result is not VariableInputTable")
	}

	return vit.Variables, nil
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

func (m VariableInputTable) Separator() string {
	width := m.NameColumnLen() + 2 +
		m.RequiredColumnLen() + 2 +
		m.DefaultValueColumnLen() + 2 +
		inputValueWidth

	return strings.Repeat("─", width)
}

func (m VariableInputTable) NameColumnLen() int {
	var maxLen int
	for _, v := range m.Variables {
		if len(v.Name) > maxLen {
			maxLen = len(v.Name)
		}
	}

	if maxLen > maxNameWidth {
		maxLen = maxNameWidth
	}

	return max(maxLen, len(nameHeader))
}

func (m VariableInputTable) RequiredColumnLen() int {
	return len(requiredHeader)
}

func (m VariableInputTable) DefaultValueColumnLen() int {
	var maxLen int
	for _, v := range m.Variables {
		if len(v.DefaultValue) > maxLen {
			maxLen = len(v.DefaultValue)
		}
	}

	if maxLen > maxDefaultValueWidth {
		maxLen = maxDefaultValueWidth
	}

	return max(maxLen, len(defaultValueHeader)) + 2 // add "[]"
}

func (m VariableInputTable) Init() tea.Cmd {
	return textinput.Blink
}

func (m VariableInputTable) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.cursor < len(m.Variables)-1 {
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
		m.Variables[i].InputValue = ti.Value()
	}

	return m, cmd
}

func (m VariableInputTable) View() string {
	s := fmt.Sprintf("%s\n\n", headerMessage)

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

	for i, v := range m.Variables {
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

	s += fmt.Sprintf("\n%s\n", footerMessage)

	return s
}
