package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
)

// Theme defines the visual styles for the CLI
type Theme struct {
	Name             string `json:"name"`
	PrimaryColor     string `json:"primary_color"`
	SecondaryColor   string `json:"secondary_color"`
	ErrorColor       string `json:"error_color"`
	SuccessColor     string `json:"success_color"`
	PanelBorderColor string `json:"panel_border_color"`
	TextColor        string `json:"text_color"`
}

var CurrentTheme *Theme

// Default themes
var themes = map[string]*Theme{
	"default": {
		Name:             "default",
		PrimaryColor:     "12", // Light Blue
		SecondaryColor:   "10", // Green
		ErrorColor:       "9",  // Red
		SuccessColor:     "10", // Green
		PanelBorderColor: "8",  // Gray
		TextColor:        "15", // White
	},
	"dark-balkan": {
		Name:             "dark-balkan",
		PrimaryColor:     "5",  // Magenta
		SecondaryColor:   "14", // Light Pink
		ErrorColor:       "1",  // Red
		SuccessColor:     "2",  // Green
		PanelBorderColor: "7",  // Light Gray
		TextColor:        "15", // White
	},
}

func InitTheme() {
	home, err := os.UserHomeDir()
	if err != nil {
		CurrentTheme = themes["default"]
		return
	}

	configPath := filepath.Join(home, ".kncli", "theme.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		CurrentTheme = themes["default"]
		return
	}

	var userTheme Theme
	if err := json.Unmarshal(data, &userTheme); err != nil {
		CurrentTheme = themes["default"]
		return
	}
	CurrentTheme = &userTheme
}

func SetTheme(name string) error {
	theme, ok := themes[name]
	if !ok {
		return fmt.Errorf("theme %q not found", name)
	}

	CurrentTheme = theme
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(home, ".kncli", "theme.json")
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(theme, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// GetStyle returns a lipgloss style based on the current theme
func GetStyle(colorKey string) lipgloss.Style {
	if CurrentTheme == nil {
		CurrentTheme = themes["default"]
	}

	var color string
	switch colorKey {
	case "primary":
		color = CurrentTheme.PrimaryColor
	case "secondary":
		color = CurrentTheme.SecondaryColor
	case "error":
		color = CurrentTheme.ErrorColor
	case "success":
		color = CurrentTheme.SuccessColor
	case "border":
		color = CurrentTheme.PanelBorderColor
	default:
		color = CurrentTheme.TextColor
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color))
}
