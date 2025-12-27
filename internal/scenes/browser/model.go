package browser

import (
	"github.com/mwantia/vista/internal/resources/keys"
	"github.com/mwantia/vista/internal/resources/theme"
	"github.com/mwantia/vista/internal/vfs"
)

// SortDirection represents the sort direction
type SortDirection int

const (
	SortAscending SortDirection = iota
	SortDescending
)

// Model represents the browser scene state
type Model struct {
	// Dependencies
	vfs   *vfs.Manager
	theme *theme.Theme
	keys  keys.KeyMap

	// State
	path    string
	entries []vfs.Entry
	cursor  int
	offset  int
	width   int
	height  int
	loading bool
	err     error

	// Preview state
	showPreview       bool
	previewContent    string
	previewError      error
	previewGeneration int

	// Mouse state
	lastClickTime int64
	lastClickX    int
	lastClickY    int

	// Breadcrumb state
	breadcrumbSegments []BreadcrumbSegment

	// Sort state
	sortDirection SortDirection
}

// New creates a new browser scene
func New(vfsManager *vfs.Manager, theme *theme.Theme, keys keys.KeyMap) Model {
	return Model{
		vfs:         vfsManager,
		theme:       theme,
		keys:        keys,
		path:        "/",
		showPreview: true,
	}
}
