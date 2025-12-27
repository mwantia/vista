package theme

import "github.com/charmbracelet/lipgloss"

// buildStyles creates a complete StyleSet from a ColorPalette
func buildStyles(colors ColorPalette) StyleSet {
	return StyleSet{
		Title: lipgloss.NewStyle().
			Foreground(colors.Primary).
			Bold(true).
			Padding(0, 1),

		StatusBar: lipgloss.NewStyle().
			Foreground(colors.Foreground).
			Background(colors.Border).
			Padding(0, 1),

		SelectedItem: lipgloss.NewStyle().
			Foreground(colors.HighlightText).
			Background(colors.Highlight).
			Bold(true),

		NormalItem: lipgloss.NewStyle().
			Foreground(colors.Foreground),

		Directory: lipgloss.NewStyle().
			Foreground(colors.Primary).
			Bold(true),

		File: lipgloss.NewStyle().
			Foreground(colors.Foreground),

		Border: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colors.Border).
			Padding(1),

		Error: lipgloss.NewStyle().
			Foreground(colors.Error).
			Bold(true),

		Loading: lipgloss.NewStyle().
			Foreground(colors.Dim).
			Italic(true),

		Help: lipgloss.NewStyle().
			Foreground(colors.Dim),

		Breadcrumb: lipgloss.NewStyle().
			Foreground(colors.Secondary),

		PreviewPane: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colors.Border).
			Padding(1).
			Foreground(colors.Foreground),
	}
}
