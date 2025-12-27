package browser

import "strings"

// BreadcrumbSegment represents a clickable path segment in the breadcrumb
type BreadcrumbSegment struct {
	path   string // Full path to this segment
	label  string // Display label
	startX int    // Starting X position
	endX   int    // Ending X position
}

// GetBreadcrumbSegmentAtPosition finds which breadcrumb segment was clicked
func (m Model) GetBreadcrumbSegmentAtPosition(x int) *BreadcrumbSegment {
	for i := range m.breadcrumbSegments {
		segment := &m.breadcrumbSegments[i]
		if x >= segment.startX && x < segment.endX {
			return segment
		}
	}

	return nil
}

// calculateBreadcrumbSegments calculates clickable breadcrumb segments from path
// This is a pure function called in Update() when the path changes
func CalculateBreadcrumbSegments(path string) []BreadcrumbSegment {
	prefix := "Vista - Virtual File-System Manager - "

	// Split path into segments
	segments := []string{}
	paths := []string{}

	if path == "/" {
		segments = append(segments, "/")
		paths = append(paths, "/")
	} else {
		parts := strings.Split(strings.Trim(path, "/"), "/")
		currentPath := ""
		for _, part := range parts {
			if part == "" {
				continue
			}
			currentPath += "/" + part
			segments = append(segments, part)
			paths = append(paths, currentPath)
		}
	}

	// Build breadcrumb segments and track positions
	result := make([]BreadcrumbSegment, 0)
	currentX := len(prefix)
	for i, segment := range segments {
		startX := currentX
		currentX += len(segment)

		result = append(result, BreadcrumbSegment{
			path:   paths[i],
			label:  segment,
			startX: startX,
			endX:   currentX,
		})

		// Add separator spacing if not last segment
		if i < len(segments)-1 {
			currentX += 3 // " → " is 3 characters
		}
	}

	return result
}

// renderBreadcrumb renders the path breadcrumb with clickable segments
// This is read-only - segments are pre-calculated by calculateBreadcrumbSegments()
func (m Model) renderBreadcrumb() string {
	prefix := "Vista - Virtual File-System Manager - "
	var breadcrumb strings.Builder
	breadcrumb.WriteString(prefix)

	for i, segment := range m.breadcrumbSegments {
		breadcrumb.WriteString(segment.label)
		// Add separator if not last segment
		if i < len(m.breadcrumbSegments)-1 {
			breadcrumb.WriteString(" → ")
		}
	}

	return m.theme.Styles.Title.Render(breadcrumb.String())
}
