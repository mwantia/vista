package vfs

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mwantia/vfs"
	"github.com/mwantia/vfs/mount"
)

// Manager wraps the VFS library and provides tea.Cmd operations
type Manager struct {
	vfs vfs.VirtualFileSystem
	ctx context.Context
}

// NewManager creates a new VFS manager and mounts ephemeral storage at root
func NewManager(ctx context.Context, uri string) (*Manager, error) {
	vfs, err := vfs.NewVirtualFileSystem()
	if err != nil {
		return nil, fmt.Errorf("failed to create VFS: %w", err)
	}

	// Use ephemeral as default mount
	if uri == "" {
		uri = "ephemeral://"
	}

	steps, err := mount.IdentifyMountSteps(ctx, uri)
	if err != nil {
		return nil, fmt.Errorf("failed to identify mount for root '/': %w", err)
	}

	// Mount ephemeral storage at root
	if err := vfs.Mount(ctx, "/", steps...); err != nil {
		return nil, fmt.Errorf("failed to mount root: %w", err)
	}

	return &Manager{
		vfs: vfs,
		ctx: ctx,
	}, nil
}

// Shutdown cleanly shuts down the VFS
func (m *Manager) Shutdown() error {
	return m.vfs.Shutdown(m.ctx)
}

// LoadDirectory returns a tea.Cmd that loads directory contents
func (m *Manager) LoadDirectory(path string) tea.Cmd {
	return func() tea.Msg {
		entries, err := m.vfs.ReadDirectory(m.ctx, path)
		if err != nil {
			return DirectoryLoadErrorMsg{
				Path:  path,
				Error: err,
			}
		}

		// Convert data.Metadata to our Entry model
		vfsEntries := make([]Entry, 0, len(entries))
		// Add regular entries
		for _, meta := range entries {
			vfsEntries = append(vfsEntries, EntryFromMetadata(path, meta))
		}

		return DirectoryLoadedMsg{
			Path:    path,
			Entries: vfsEntries,
		}
	}
}

// DirectoryLoadedMsg is emitted when a directory loads successfully
type DirectoryLoadedMsg struct {
	Path    string
	Entries []Entry
}

// DirectoryLoadErrorMsg is emitted when a directory fails to load
type DirectoryLoadErrorMsg struct {
	Path  string
	Error error
}

// GeneratePreview generates a preview for the given entry
func (m *Manager) GeneratePreview(entry *Entry, width, height int) (string, error) {
	return GeneratePreview(m.vfs, m.ctx, entry, width, height)
}
