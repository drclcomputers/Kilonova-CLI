// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package utility

import (
	"encoding/json"
	"fmt"

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

type ProblemInfo struct {
	Data struct {
		Name          string  `json:"name"`
		Time          float64 `json:"time_limit"`
		MemoryLimit   int     `json:"memory_limit"`
		SourceSize    int     `json:"source_size"`
		SourceCredits string  `json:"source_credits"`
	} `json:"data"`
}

// TEXT MODEL

type TextModel struct {
	viewport viewport.Model
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
		case "up", "k":
			m.viewport.LineUp(1)
		case "down", "j":
			m.viewport.LineDown(1)
		}
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
	return &TextModel{viewport: vp}
}

// TABLE MODEL

type Model struct {
	table table.Model
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
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	footer := "\n(Use ↑/↓ to navigate, 'q' to quit, 'enter' to get the statement)"
	return lipgloss.NewStyle().Margin(1, 2).Render(m.table.View()) + footer
}

func NewTable(table table.Model) *Model {
	return &Model{table: table}
}

// Search Table

var ChosenProblem string = ""

var GlobalRows []table.Row

type TableSearch struct {
	table table.Model
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
	}

	var Command tea.Cmd
	TableModel.table, Command = TableModel.table.Update(Message)
	return TableModel, Command
}

func (TableModel TableSearch) HandleSelection() (tea.Model, tea.Cmd) {
	SelectedIndex := TableModel.table.Cursor()
	SelectedProblem := GlobalRows[SelectedIndex]

	ChosenProblem = string(SelectedProblem[0])

	return TableModel, tea.Quit
}

func (TableModel TableSearch) View() string {
	return lipgloss.NewStyle().Margin(1, 2).Render(TableModel.table.View()) + "\n(Use ↑/↓ to navigate, 'q' to quit, 'enter' to get the statement)"
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
		table.WithHeight(20),
	)
	t.SetStyles(table.DefaultStyles())
	return t
}

func RenderTable(columns []table.Column, rows []table.Row, TableType int) {
	t := CreateTable(columns, rows)
	program := tea.NewProgram(&Model{table: t}, tea.WithAltScreen()) // 1 - Normal Table
	if TableType == 2 {                                              // 2 - Search Table
		program = tea.NewProgram(NewSearchTable(t), tea.WithAltScreen())
	}
	if _, err := program.Run(); err != nil {
		LogError(fmt.Errorf("error running program: %w", err))
	}
}
