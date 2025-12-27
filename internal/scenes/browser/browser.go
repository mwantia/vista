package browser

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mwantia/vista/internal/resources/theme"
	"github.com/mwantia/vista/internal/vfs"
)

// Init initializes the browser scene
func (m Model) Init() tea.Cmd {
	// Load root directory on initialization
	return m.vfs.LoadDirectory(m.path)
}

// Update handles messages for the browser scene
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case vfs.DirectoryLoadedMsg:
		m.entries = msg.Entries
		// Sort entries: folders first, then files
		SortEntries(m.entries, m.sortDirection)
		m.loading = false
		m.err = nil
		m.cursor = 0
		m.offset = 0
		m.path = msg.Path
		m.previewGeneration++                                      // Invalidate old previews
		m.breadcrumbSegments = CalculateBreadcrumbSegments(m.path) // Calculate breadcrumb segments on path change
		return m, m.loadPreview()

	case vfs.DirectoryLoadErrorMsg:
		m.err = msg.Error
		m.loading = false
		return m, nil

	case previewLoadedMsg:
		// Only update if generation matches (prevents race conditions)
		if msg.generation == m.previewGeneration {
			m.previewContent = msg.content
			m.previewError = msg.err
		}
		return m, nil

	case tea.MouseMsg:
		return m.handleMouseEvent(msg)

	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	return m, nil
}

// handleKeyPress handles keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Don't handle keys if loading
	if m.loading {
		return m, nil
	}

	switch {
	case key.Matches(msg, m.keys.Up):
		if m.cursor > 0 {
			m.cursor--
			m.adjustOffset()
			m.previewGeneration++ // Invalidate old preview
			return m, m.loadPreview()
		}
		return m, nil

	case key.Matches(msg, m.keys.Down):
		if m.cursor < len(m.entries)-1 {
			m.cursor++
			m.adjustOffset()
			m.previewGeneration++ // Invalidate old preview
			return m, m.loadPreview()
		}
		return m, nil

	case key.Matches(msg, m.keys.Enter):
		if m.cursor < len(m.entries) {
			entry := m.entries[m.cursor]

			if entry.Mode.IsDir() || entry.Mode.IsMount() {
				// Navigate into directory or mount point (keep current view visible while loading)
				return m, m.vfs.LoadDirectory(entry.Path)
			}
			// For files, just emit a selection message for now
			return m, func() tea.Msg {
				return FileSelectedMsg{Entry: entry}
			}
		}
		return m, nil

	case key.Matches(msg, m.keys.Back):
		// Go up one directory (keep current view visible while loading)
		if m.path != "/" {
			parentPath := GetParentPath(m.path)
			return m, m.vfs.LoadDirectory(parentPath)
		}
		return m, nil

	case key.Matches(msg, m.keys.TogglePreview):
		m.showPreview = !m.showPreview
		if m.showPreview {
			m.previewGeneration++ // Invalidate old preview
			return m, m.loadPreview()
		}
		return m, nil

	case key.Matches(msg, m.keys.ToggleSortDirection):
		// Toggle sort direction
		if m.sortDirection == SortAscending {
			m.sortDirection = SortDescending
		} else {
			m.sortDirection = SortAscending
		}
		// Re-sort current entries
		SortEntries(m.entries, m.sortDirection)
		return m, nil
	}

	return m, nil
}

// handleMouseEvent handles mouse input
func (m Model) handleMouseEvent(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Don't handle mouse if loading
	if m.loading {
		return m, nil
	}

	// Handle mouse wheel events
	if msg.Action == tea.MouseActionPress {
		switch msg.Button {
		case tea.MouseButtonWheelUp:
			// Scroll up (move cursor up)
			if m.cursor > 0 {
				m.cursor--
				m.adjustOffset()
				m.previewGeneration++ // Invalidate old preview
				return m, m.loadPreview()
			}
			return m, nil

		case tea.MouseButtonWheelDown:
			// Scroll down (move cursor down)
			if m.cursor < len(m.entries)-1 {
				m.cursor++
				m.adjustOffset()
				m.previewGeneration++ // Invalidate old preview
				return m, m.loadPreview()
			}
			return m, nil

		case tea.MouseButtonLeft:
			// Check if click is on breadcrumb (y == 0)
			if msg.Y == 0 {
				if segment := m.GetBreadcrumbSegmentAtPosition(msg.X); segment != nil {
					// Navigate to clicked breadcrumb segment
					return m, m.vfs.LoadDirectory(segment.path)
				}
				return m, nil
			}

			// Left click on file list - select or navigate
			if clickedEntry := m.getEntryAtPosition(msg.X, msg.Y); clickedEntry >= 0 {
				// Check for double-click (within 500ms)
				now := time.Now().UnixMilli()
				isDoubleClick := (now-m.lastClickTime) < 500 &&
					msg.X == m.lastClickX &&
					msg.Y == m.lastClickY

				m.lastClickTime = now
				m.lastClickX = msg.X
				m.lastClickY = msg.Y

				if isDoubleClick {
					// Double-click - navigate into directory/mount or open file
					m.cursor = clickedEntry
					if m.cursor < len(m.entries) {
						entry := m.entries[m.cursor]
						if entry.Mode.IsDir() || entry.Mode.IsMount() {
							return m, m.vfs.LoadDirectory(entry.Path)
						}
						// For files, emit selection message
						return m, func() tea.Msg {
							return FileSelectedMsg{Entry: entry}
						}
					}
				} else {
					// Single click - select entry
					m.cursor = clickedEntry
					m.adjustOffset()
					m.previewGeneration++ // Invalidate old preview
					return m, m.loadPreview()
				}
			}
			return m, nil

		case tea.MouseButtonRight:
			// Right click - navigate back to parent
			if m.path != "/" {
				parentPath := GetParentPath(m.path)
				return m, m.vfs.LoadDirectory(parentPath)
			}
			return m, nil
		}
	}

	return m, nil
}

// getEntryAtPosition calculates which entry was clicked based on mouse position
func (m Model) getEntryAtPosition(x, y int) int {
	// Check if we have entries
	if len(m.entries) == 0 {
		return -1
	}

	// Calculate file list width (left side if preview is shown)
	fileListWidth := m.width
	if m.showPreview {
		fileListWidth = m.width / 2
	}

	// Check if click is within the file list area (not in preview pane)
	if m.showPreview && x >= fileListWidth {
		return -1 // Click was in preview pane
	}

	// Y position calculation:
	// Line 0: Breadcrumb
	// Line 1: Border top
	// Line 2: Padding top (empty line due to Padding(1))
	// Line 3: Header row
	// Line 4+: File entries
	const headerOffset = 4 // Breadcrumb + border + padding + header

	if y < headerOffset {
		return -1 // Click was in header area
	}

	// Calculate which entry was clicked
	relativeY := y - headerOffset
	entryIndex := relativeY + m.offset

	// Validate entry index
	if entryIndex < 0 || entryIndex >= len(m.entries) {
		return -1
	}

	return entryIndex
}

// loadPreview generates a preview for the currently selected entry
func (m Model) loadPreview() tea.Cmd {
	// Don't load preview if disabled or no entry selected
	if !m.showPreview || len(m.entries) == 0 || m.cursor >= len(m.entries) {
		return nil
	}

	generation := m.previewGeneration
	entry := m.entries[m.cursor]

	return func() tea.Msg {
		// Calculate preview dimensions (right half of screen)
		width := m.width / 2
		height := m.height - 5 // Account for breadcrumb, status bar, helper bar

		content, err := m.vfs.GeneratePreview(&entry, width, height)
		return previewLoadedMsg{
			content:    content,
			err:        err,
			generation: generation,
		}
	}
}

// View renders the browser scene
func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	sections := []string{
		// Breadcrumb
		m.renderBreadcrumb(),
		// Main content area
		m.renderContent(),
		// Status bar
		m.renderStatusBar(),
		// Helper bar
		m.renderHelperBar(),
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderContent renders the main file list or error/loading state
func (m Model) renderContent() string {
	// Reserve space for: breadcrumb (1), status bar (1), helper bar (1) = 3 lines total
	contentHeight := m.height - 3

	if m.err != nil {
		errMsg := m.theme.Styles.Error.Render(fmt.Sprintf("Error: %v", m.err))
		return lipgloss.Place(m.width, contentHeight, lipgloss.Center, lipgloss.Center, errMsg)
	}

	if len(m.entries) == 0 {
		// Show loading only on initial load (no entries, no error)
		if m.loading {
			loading := m.theme.Styles.Loading.Render("Loading directory...")
			return lipgloss.Place(m.width, contentHeight, lipgloss.Center, lipgloss.Center, loading)
		}
		empty := m.theme.Styles.NormalItem.Render("(empty directory)")
		return lipgloss.Place(m.width, contentHeight, lipgloss.Center, lipgloss.Center, empty)
	}

	// Show entries even if loading (keeps UI responsive during navigation)
	if !m.showPreview {
		// Full width file list
		return m.renderFileList(m.width, contentHeight)
	}

	// Split view: file list (left) + preview (right)
	leftWidth := m.width / 2
	rightWidth := m.width - leftWidth

	fileList := m.renderFileList(leftWidth, contentHeight)
	preview := m.renderPreview(rightWidth, contentHeight)

	return lipgloss.JoinHorizontal(lipgloss.Top, fileList, preview)
}

// renderFileList renders the list of files and directories
func (m Model) renderFileList(width, height int) string {
	var lines []string

	// Add table header
	header := m.renderTableHeader(width)
	lines = append(lines, header)

	// Calculate visible lines accounting for: borders (2), padding (2), header (1) = 5 lines
	visibleLines := max(height-5, 1)
	start := m.offset
	end := min(m.offset+visibleLines, len(m.entries))

	for i := start; i < end; i++ {
		entry := m.entries[i]
		line := m.renderEntry(entry, i == m.cursor, i%2 == 0, width)
		lines = append(lines, line)
	}

	// Fill remaining space with empty lines
	for len(lines) < visibleLines+1 { // +1 for header
		lines = append(lines, "")
	}

	content := strings.Join(lines, "\n")

	// Add border around the file list
	// Subtract 4 for borders (2) + padding (2) to prevent overflow
	return m.theme.Styles.Border.
		Width(width - 6).
		Height(height - 4).
		Render(content)
}

// renderPreview renders the preview pane
func (m Model) renderPreview(width, height int) string {
	// No entry selected
	if len(m.entries) == 0 || m.cursor >= len(m.entries) {
		return m.theme.Styles.PreviewPane.
			Width(width).
			Height(height - 2). // Subtract 4 for borders (2) + padding (2)
			Render("No file selected")
	}

	entry := m.entries[m.cursor]

	// Preview error
	if m.previewError != nil {
		errorContent := fmt.Sprintf("Error loading preview:\n\n%v", m.previewError)
		return m.theme.Styles.PreviewPane.
			Width(width).
			Height(height - 2). // Subtract 4 for borders (2) + padding (2)
			Render(errorContent)
	}

	// Loading preview
	if m.previewContent == "" {
		return m.theme.Styles.PreviewPane.
			Width(width).
			Height(height - 2). // Subtract 4 for borders (2) + padding (2)
			Render("Loading preview...")
	}

	// Build header with file metadata
	var header strings.Builder
	header.WriteString(fmt.Sprintf("%s %s\n\n",
		entry.Icon(),
		entry.DisplayName()))
	header.WriteString(fmt.Sprintf("Size: %s | Modified: %s | Type: %s\n",
		entry.DisplaySize(),
		entry.DisplayModifyTime(),
		entry.DisplayContentType()))
	header.WriteString(strings.Repeat("─", width-4))
	header.WriteString("\n\n")

	// Combine header and content, then truncate to fit height
	fullContent := header.String() + m.previewContent

	// Truncate content to fit within available height, accounting for line wrapping
	// Note: lipgloss Height() already accounts for border and padding internally
	// We just need to account for the initial height-2 adjustment
	maxLines := max(height-5, 5)
	// For width, account for padding (2) on left+right
	maxWidth := max(width-2, 40)
	truncatedContent := TruncateContentToHeight(fullContent, maxLines, maxWidth)

	return m.theme.Styles.PreviewPane.
		Width(width).
		Height(height - 2). // Subtract 4 for borders (2) + padding (2)
		Render(truncatedContent)
}

// renderTableHeader renders the column headers
func (m Model) renderTableHeader(width int) string {
	// Calculate column widths
	contentTypeWidth := 30
	modifyTimeWidth := 20
	sizeWidth := 10
	iconWidth := 2
	padding := 8 // 2 spaces between each of 4 columns
	nameWidth := max(width-contentTypeWidth-modifyTimeWidth-sizeWidth-iconWidth-padding-12, 20)

	header := fmt.Sprintf("%s %-*s  %-*s  %-*s  %*s",
		"  ", // Icon column placeholder
		nameWidth, "Name",
		contentTypeWidth, "Content Type",
		modifyTimeWidth, "Modified",
		sizeWidth, "Size",
	)

	return m.theme.Styles.NormalItem.
		Bold(true).
		Render(header)
}

// renderEntry renders a single file or directory entry
func (m Model) renderEntry(entry vfs.Entry, selected bool, evenRow bool, width int) string {
	// Calculate column widths (must match header)
	contentTypeWidth := 30
	modifyTimeWidth := 20
	sizeWidth := 10
	iconWidth := 2
	padding := 8 // 2 spaces between each of 4 columns
	nameWidth := max(width-contentTypeWidth-modifyTimeWidth-sizeWidth-iconWidth-padding-12, 20)

	// Get formatted values
	icon := entry.Icon()
	name := Truncate(entry.DisplayName(), nameWidth)
	contentType := Truncate(entry.DisplayContentType(), contentTypeWidth)
	modifyTime := entry.DisplayModifyTime()
	size := entry.DisplaySize()
	// Dynamically adjust background brightness by 5%
	adjustedBg := theme.AdjustBrightness(m.theme.Colors.Background, 5)

	var style lipgloss.Style
	if selected {
		style = m.theme.Styles.SelectedItem
	} else if entry.Mode.IsMount() {
		style = m.theme.Styles.Directory.Foreground(m.theme.Colors.Secondary)
		if evenRow {
			style = style.Background(adjustedBg)
		}
	} else if entry.Mode.IsDir() {
		style = m.theme.Styles.Directory
		if evenRow {
			style = style.Background(adjustedBg)
		}
	} else {
		style = m.theme.Styles.File
		if evenRow {
			style = style.Background(adjustedBg)
		}
	}

	line := fmt.Sprintf("%s %-*s  %-*s  %-*s  %*s",
		icon,
		nameWidth, name,
		contentTypeWidth, contentType,
		modifyTimeWidth, modifyTime,
		sizeWidth, size,
	)

	return style.Render(line)
}

// renderStatusBar renders the bottom status bar
func (m Model) renderStatusBar() string {
	var status string
	if len(m.entries) == 0 {
		status = "0 Items"
	} else {
		status = fmt.Sprintf("%d/%d Items", m.cursor+1, len(m.entries))
	}
	return m.theme.Styles.StatusBar.
		Width(m.width).
		Render(status)
}

// renderHelperBar renders the keyboard shortcuts helper bar
func (m Model) renderHelperBar() string {
	sortDir := "ASC"
	if m.sortDirection == SortDescending {
		sortDir = "DESC"
	}
	helper := fmt.Sprintf("↑/k up • ↓/j down • enter/l open • bksp/h back • tab preview • s sort:%s • q/ctrl+c quit • ? help", sortDir)
	return m.theme.Styles.Help.
		Width(m.width).
		Render(helper)
}

// adjustOffset adjusts the scroll offset to keep the cursor visible
func (m *Model) adjustOffset() {
	visibleLines := max(m.height-7, 1)

	if m.cursor < m.offset {
		m.offset = m.cursor
	} else if m.cursor >= m.offset+visibleLines {
		m.offset = m.cursor - visibleLines + 1
	}
}
