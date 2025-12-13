package startscreen

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	config "sweep/config"
	misc "sweep/shared/consts/misc"
	styles "sweep/tui/styles"

	cursor "github.com/charmbracelet/bubbles/cursor"
	textinput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	focusedButton = styles.Seven.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", styles.Zero.Render("Submit"))
)

type model struct {
	focusIndex int
	inputs     []textinput.Model
	messages   [][]string
	config     *config.Config
	isValid    bool
}

var _ tea.Model = model{}

func (m model) validateNumber(s string) error {
	_, error := strconv.ParseInt(s, 10, 64)
	if error == nil {
		return nil
	}
	return errors.New("should be an integer")

}

func CreateModel(config *config.Config) model {
	m := model{
		inputs:     make([]textinput.Model, 3),
		focusIndex: 0,
		messages:   make([][]string, 3),
		config:     config,
		isValid:    false,
	}

	var input textinput.Model
	for i := range m.inputs {
		input.Cursor.Style = styles.Seven
		input = textinput.New()
		input.Validate = m.validateNumber
		input.Cursor = cursor.New()

		mines := strconv.FormatUint(uint64(config.Defaults.Mines), 10)
		height := strconv.FormatUint(uint64(config.Defaults.Height), 10)
		width := strconv.FormatUint(uint64(config.Defaults.Width), 10)

		input.Width = 5
		input.CharLimit = 5
		input.Placeholder = "0"
		switch i {
		case 0:
			input.Focus()
			input.PromptStyle = styles.BrightText
			input.TextStyle = styles.BrightText

			input.Prompt = "field width "
			if width != "0" {
				input.SetValue(width)
			}
		case 1:
			input.PromptStyle = styles.DimText
			input.TextStyle = styles.DimText

			input.Prompt = "field height "
			if height != "0" {
				input.SetValue(height)
			}
		case 2:
			input.PromptStyle = styles.DimText
			input.TextStyle = styles.DimText

			input.Prompt = "amount of mines "
			if mines != "0" {
				input.SetValue(mines)
			}
		}
		m.inputs[i] = input
	}
	return m
}

func (m model) Init() tea.Cmd {
	m.validateInputs()
	return tea.Batch(tea.SetWindowTitle(misc.APP_NAME), tea.ClearScreen)
}

func (m *model) validateInputs() {
	m.messages = make([][]string, 3)
	m.isValid = true
	for ix, input := range m.inputs {
		if input.Value() == "" {
			m.messages[ix] = append(m.messages[ix], fmt.Sprintf("%vmust not be empty", input.Prompt))
			m.isValid = false
		}
		if input.Err != nil {
			m.messages[ix] = append(m.messages[ix], fmt.Sprintf("%v%v\n", input.Prompt, input.Err.Error()))
			m.isValid = false
		}
	}

	const (
		widthIx  int = 0
		heightIx int = 1
		minesIx  int = 2
	)

	width, _ := strconv.ParseUint(m.inputs[widthIx].Value(), 10, 16)
	height, _ := strconv.ParseUint(m.inputs[heightIx].Value(), 10, 16)
	mines, _ := strconv.ParseUint(m.inputs[minesIx].Value(), 10, 16)

	if width == 0 {
		m.messages[widthIx] = append(m.messages[widthIx], fmt.Sprintln("field width cannot be zero"))
		m.isValid = false
	}
	if height == 0 {
		m.messages[heightIx] = append(m.messages[heightIx], fmt.Sprintln("field height cannot be zero"))
		m.isValid = false
	}
	if mines == 0 {
		m.messages[minesIx] = append(m.messages[minesIx], fmt.Sprintln("amount of mines cannot be zero"))
		m.isValid = false
	}
	if mines >= width*height {
		m.messages[minesIx] = append(m.messages[minesIx], fmt.Sprintln("amount of mines should be less than field area (width * height)"))
		m.isValid = false
	}
}

func (m model) updateInputs(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
			"backspace", "delete",
			"right", "left":

			cmds := make([]tea.Cmd, len(m.inputs))

			for i := range m.inputs {
				m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
			}

			return tea.Batch(cmds...)
		}
	}
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			os.Exit(0)

		case "tab", "shift+tab", "enter", "up", "down", "k", "j":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs) {
				if !m.isValid {
					return m, nil
				}
				widthStr := m.inputs[0].Value()
				heightStr := m.inputs[1].Value()
				minesStr := m.inputs[2].Value()

				width, _ := strconv.ParseUint(widthStr, 10, 16)
				height, _ := strconv.ParseUint(heightStr, 10, 16)
				mines, _ := strconv.ParseUint(minesStr, 10, 16)

				m.config.Width = uint16(width)
				m.config.Height = uint16(height)
				m.config.Mines = uint16(mines)

				return m, tea.Quit
			}

			if s == "up" || s == "shift+tab" || s == "k" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = styles.BrightText
					m.inputs[i].TextStyle = styles.BrightText
					continue
				}
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = styles.DimText
				m.inputs[i].TextStyle = styles.DimText
			}

			return m, tea.Batch(cmds...)
		}
	}

	m.validateInputs()
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString(styles.HeaderStyle.Render(misc.APP_ASCII_LOGO))
	b.WriteRune('\n')

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		isNotLast := i != len(m.inputs)
		if isNotLast {
			b.WriteRune('\n')
		}
		if len(m.messages[i]) > 0 {
			b.WriteString(styles.BrightText.Faint(true).Render(m.messages[i][0]))
		}
		if isNotLast {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n%s\n", *button)

	return b.String()
}
