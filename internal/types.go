// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package internal

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Types

type UserId struct {
	Data struct {
		ID int `json:"id"`
	} `json:"data"`
}

type RawKilonovaResponse struct {
	Status string          `json:"status"`
	Data   json.RawMessage `json:"data"`
}

type KilonovaResponse struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

type Problem struct {
	Id            int     `json:"id"`
	Name          string  `json:"name"`
	Time          float64 `json:"time_limit"`
	MemoryLimit   int     `json:"memory_limit"`
	SourceSize    int     `json:"source_size"`
	SourceCredits string  `json:"source_credits"`
	MaxScore      int     `json:"max_score"`
}
type ProblemInfo struct {
	Data Problem `json:"data"`
}

type ProblemList struct {
	Status string `json:"status"`
	Data   []Problem
}

// TEXT MODEL

type TextModel struct {
	viewport viewport.Model
	height   int
	width    int
	text     string
}

func (m *TextModel) Init() tea.Cmd {
	return nil
}

func (m *TextModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.String(); key {
		case "q", "esc":
			return m, tea.Quit
		case "up", "o":
			m.viewport.LineUp(1)
		case "down", "k":
			m.viewport.LineDown(1)
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 5

		m.viewport.SetContent(m.text)
	}
	return m, cmd
}

func (m *TextModel) View() string {
	style := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1)
	footer := "\n(Use ↑/↓ to scroll, 'q' to quit)"
	return style.Render(m.viewport.View()) + footer
}

func NewTextModel(text string) *TextModel {
	vp := viewport.New(80, 25)
	vp.SetContent(text)
	return &TextModel{
		viewport: vp,
		text:     text,
	}
}

// TABLE MODEL

type Model struct {
	table  table.Model
	width  int
	height int
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.String(); key {
		case "q", "esc":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	tableView := m.table.View()

	tableLines := strings.Count(tableView, "\n") + 1

	spaceLines := m.height - tableLines - 2
	if spaceLines < 0 {
		spaceLines = 0
	}

	spacing := strings.Repeat("\n", spaceLines)

	footer := spacing + "\n(Use ↑/↓ to navigate, 'q' to quit)"
	return lipgloss.NewStyle().Margin(1, 2).Render(m.table.View()) + footer
}

func NewTable(table table.Model) *Model {
	return &Model{table: table}
}

// Search Table

var ChosenProblem = ""

var GlobalRows []table.Row

type TableSearch struct {
	table  table.Model
	height int
	width  int
}

func (TableModel TableSearch) Init() tea.Cmd {
	return nil
}

func (TableModel TableSearch) Update(Message tea.Msg) (tea.Model, tea.Cmd) {
	switch Message := Message.(type) {
	case tea.KeyMsg:
		switch Message.String() {
		case "q":
			return TableModel, tea.Quit
		case "esc":
			return TableModel, tea.Quit
		case "enter":
			return TableModel.HandleSelection()
		}
	case tea.WindowSizeMsg:
		TableModel.width = Message.Width
		TableModel.height = Message.Height
	}

	var Command tea.Cmd
	TableModel.table, Command = TableModel.table.Update(Message)
	return TableModel, Command
}

func (TableModel TableSearch) HandleSelection() (tea.Model, tea.Cmd) {
	SelectedIndex := TableModel.table.Cursor()
	SelectedProblem := GlobalRows[SelectedIndex]

	ChosenProblem = SelectedProblem[0]

	return TableModel, tea.Quit
}

func (TableModel TableSearch) View() string {
	tableView := TableModel.table.View()

	tableLines := strings.Count(tableView, "\n") + 1

	spaceLines := TableModel.height - tableLines - 2
	if spaceLines < 0 {
		spaceLines = 0
	}

	spacing := strings.Repeat("\n", spaceLines)

	footer := spacing + "\n(Use ↑/↓ to navigate, 'q' to quit, 'enter' to get the statement)"
	return lipgloss.NewStyle().Margin(1, 2).Render(TableModel.table.View()) + footer
}

func NewSearchTable(table table.Model) *TableSearch {
	return &TableSearch{table: table}
}

// Create Tables

func CreateTable(Columns []table.Column, Rows []table.Row) table.Model {
	t := table.New(
		table.WithColumns(Columns),
		table.WithRows(Rows),
		table.WithFocused(true),
	)
	t.SetStyles(table.DefaultStyles())
	return t
}

func RenderTable(columns []table.Column, rows []table.Row, TableType int) {
	t := CreateTable(columns, rows)
	program := tea.NewProgram(NewTable(t), tea.WithAltScreen()) // 1 - Normal Table
	if TableType == 2 {                                         // 2 - Search Table
		program = tea.NewProgram(NewSearchTable(t), tea.WithAltScreen())
	}
	if _, err := program.Run(); err != nil {
		LogError(fmt.Errorf("error running program: %w", err))
	}
}
