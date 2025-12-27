package browser

import "github.com/mwantia/vista/internal/vfs"

// FileSelectedMsg is emitted when a user selects a file
type FileSelectedMsg struct {
	Entry vfs.Entry
}

// NavigateMsg requests navigation to a specific path
type NavigateMsg struct {
	Path string
}

// previewLoadedMsg is emitted when a preview has been generated
type previewLoadedMsg struct {
	content    string
	err        error
	generation int // Match against model.previewGeneration to prevent race conditions
}
