package theme

import "github.com/charmbracelet/lipgloss"

// Theme represents a complete color scheme and style set for the TUI
type Theme struct {
	Name   string
	Colors ColorPalette
	Styles StyleSet
}

// ColorPalette defines all colors used in the theme
// Uses TerminalColor interface to support both solid colors and adaptive colors
type ColorPalette struct {
	Primary       lipgloss.TerminalColor
	Secondary     lipgloss.TerminalColor
	Background    lipgloss.TerminalColor
	Foreground    lipgloss.TerminalColor
	Border        lipgloss.TerminalColor
	Highlight     lipgloss.TerminalColor
	HighlightText lipgloss.TerminalColor
	Dim           lipgloss.TerminalColor
	Error         lipgloss.TerminalColor
	Success       lipgloss.TerminalColor
	Warning       lipgloss.TerminalColor
}

// StyleSet contains pre-built Lipgloss styles
type StyleSet struct {
	Title        lipgloss.Style
	StatusBar    lipgloss.Style
	SelectedItem lipgloss.Style
	NormalItem   lipgloss.Style
	Directory    lipgloss.Style
	File         lipgloss.Style
	Border       lipgloss.Style
	Error        lipgloss.Style
	Loading      lipgloss.Style
	Help         lipgloss.Style
	Breadcrumb   lipgloss.Style
	PreviewPane  lipgloss.Style
}
