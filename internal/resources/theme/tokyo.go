package theme

import "github.com/charmbracelet/lipgloss"

// TokyoNightTheme returns the Tokyo Night color scheme (dark)
func TokyoNightTheme() *Theme {
	colors := ColorPalette{
		Primary:       lipgloss.Color("#7AA2F7"), // Blue
		Secondary:     lipgloss.Color("#BB9AF7"), // Purple
		Background:    lipgloss.Color("#1A1B26"), // Dark background
		Foreground:    lipgloss.Color("#C0CAF5"), // Light text
		Border:        lipgloss.Color("#414868"), // Gray border
		Highlight:     lipgloss.Color("#7AA2F7"), // Blue highlight
		HighlightText: lipgloss.Color("#1A1B26"), // Dark text on highlight
		Dim:           lipgloss.Color("#565F89"), // Dimmed text
		Error:         lipgloss.Color("#F7768E"), // Red
		Success:       lipgloss.Color("#9ECE6A"), // Green
		Warning:       lipgloss.Color("#E0AF68"), // Yellow
	}

	return &Theme{
		Name:   "Tokyo Night",
		Colors: colors,
		Styles: buildStyles(colors),
	}
}

// TokyoNightLightTheme returns the Tokyo Night Light color scheme (light)
func TokyoSunriseTheme() *Theme {
	colors := ColorPalette{
		Primary:       lipgloss.Color("#2E7DE9"), // Darker blue for light bg
		Secondary:     lipgloss.Color("#9854F1"), // Darker purple for light bg
		Background:    lipgloss.Color("#D5D6DB"), // Light background
		Foreground:    lipgloss.Color("#343B58"), // Dark text
		Border:        lipgloss.Color("#A8ADB7"), // Light gray border
		Highlight:     lipgloss.Color("#2E7DE9"), // Blue highlight
		HighlightText: lipgloss.Color("#FFFFFF"), // White text on highlight
		Dim:           lipgloss.Color("#6172B0"), // Dimmed text (still readable)
		Error:         lipgloss.Color("#F52A65"), // Bright red
		Success:       lipgloss.Color("#587539"), // Dark green
		Warning:       lipgloss.Color("#8C6C3E"), // Dark yellow/brown
	}

	return &Theme{
		Name:   "Tokyo Sunrise",
		Colors: colors,
		Styles: buildStyles(colors),
	}
}

// TokyoAutoTheme returns an adaptive theme that adjusts to terminal background
func TokyoAutoTheme() *Theme {
	colors := ColorPalette{
		Primary: lipgloss.AdaptiveColor{
			Light: "#2E7DE9", // Dark blue for light terminals
			Dark:  "#7AA2F7", // Light blue for dark terminals
		},
		Secondary: lipgloss.AdaptiveColor{
			Light: "#9854F1", // Dark purple for light terminals
			Dark:  "#BB9AF7", // Light purple for dark terminals
		},
		Background: lipgloss.AdaptiveColor{
			Light: "#D5D6DB", // Light background
			Dark:  "#1A1B26", // Dark background
		},
		Foreground: lipgloss.AdaptiveColor{
			Light: "#343B58", // Dark text
			Dark:  "#C0CAF5", // Light text
		},
		Border: lipgloss.AdaptiveColor{
			Light: "#A8ADB7", // Light gray border
			Dark:  "#414868", // Dark gray border
		},
		Highlight: lipgloss.AdaptiveColor{
			Light: "#2E7DE9", // Dark blue highlight
			Dark:  "#7AA2F7", // Light blue highlight
		},
		HighlightText: lipgloss.AdaptiveColor{
			Light: "#FFFFFF", // White text on highlight
			Dark:  "#1A1B26", // Dark text on highlight
		},
		Dim: lipgloss.AdaptiveColor{
			Light: "#6172B0", // Dimmed but readable
			Dark:  "#565F89", // Dimmed text
		},
		Error: lipgloss.AdaptiveColor{
			Light: "#F52A65", // Bright red
			Dark:  "#F7768E", // Softer red
		},
		Success: lipgloss.AdaptiveColor{
			Light: "#587539", // Dark green
			Dark:  "#9ECE6A", // Light green
		},
		Warning: lipgloss.AdaptiveColor{
			Light: "#8C6C3E", // Dark yellow/brown
			Dark:  "#E0AF68", // Light yellow
		},
	}

	return &Theme{
		Name:   "Tokyo",
		Colors: colors,
		Styles: buildStyles(colors),
	}
}
