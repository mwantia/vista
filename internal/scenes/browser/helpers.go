package browser

import (
	"sort"
	"strings"

	"github.com/mwantia/vista/internal/vfs"
)

// truncate truncates a string to a maximum width
func Truncate(s string, width int) string {
	if len(s) <= width {
		return s
	}
	if width < 3 {
		return s[:width]
	}
	return s[:width-3] + "..."
}

// getParentPath returns the parent directory path
func GetParentPath(path string) string {
	if path == "/" {
		return "/"
	}

	// Remove trailing slash if present
	path = strings.TrimSuffix(path, "/")

	// Find last slash
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash == 0 {
		return "/"
	}

	return path[:lastSlash]
}

// SortEntries sorts entries with folders first, then files, alphabetically by name
func SortEntries(entries []vfs.Entry, direction SortDirection) {
	sort.SliceStable(entries, func(i, j int) bool {
		entryI := entries[i]
		entryJ := entries[j]

		// Check if entries are directories/mounts
		isDirI := entryI.Mode.IsDir() || entryI.Mode.IsMount()
		isDirJ := entryJ.Mode.IsDir() || entryJ.Mode.IsMount()

		// Folders first, then files
		if isDirI != isDirJ {
			return isDirI // Folders come first
		}

		// Within same type (both folders or both files), sort alphabetically
		nameI := strings.ToLower(entryI.Name)
		nameJ := strings.ToLower(entryJ.Name)

		if direction == SortAscending {
			return nameI < nameJ
		}
		return nameI > nameJ
	})
}

// TruncateContentToHeight truncates content to fit within maxLines, accounting for line wrapping
func TruncateContentToHeight(content string, maxLines int, maxWidth int) string {
	if maxLines <= 0 || maxWidth <= 0 {
		return ""
	}

	lines := strings.Split(content, "\n")
	var result []string
	visualLineCount := 0
	truncated := false

	for i, line := range lines {
		// Calculate how many visual lines this line will take when wrapped
		// Account for ANSI codes which don't contribute to visual width
		visualLength := estimateVisualLength(line)
		wrappedLines := 1
		if visualLength > maxWidth {
			wrappedLines = (visualLength + maxWidth - 1) / maxWidth // Ceiling division
		}

		// Check if adding this line would exceed maxLines
		if visualLineCount+wrappedLines > maxLines {
			// Mark as truncated if there are more lines to show
			if i < len(lines)-1 {
				truncated = true
			}
			break
		}

		result = append(result, line)
		visualLineCount += wrappedLines
	}

	// Add truncation indicator if content was truncated
	if truncated {
		result = append(result, "...")
	}

	return strings.Join(result, "\n")
}

// estimateVisualLength estimates the visual length of a string, ignoring ANSI escape codes
func estimateVisualLength(s string) int {
	// Simple heuristic: count characters that aren't part of ANSI escape sequences
	length := 0
	inEscape := false

	for _, r := range s {
		if r == '\x1b' { // ESC character starts ANSI sequence
			inEscape = true
			continue
		}
		if inEscape {
			// ANSI sequences end with a letter (simplified)
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				inEscape = false
			}
			continue
		}
		length++
	}

	return length
}
