package app

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mwantia/vista/internal/resources/keys"
	"github.com/mwantia/vista/internal/resources/theme"
	"github.com/mwantia/vista/internal/scenes/browser"
	"github.com/mwantia/vista/internal/vfs"
)

// Model is the root Bubble Tea model that manages scenes
type Model struct {
	// Dependencies
	vfs   *vfs.Manager
	theme *theme.Theme
	keys  keys.KeyMap

	// State
	currentScene SceneType
	browserScene browser.Model

	// Dimensions
	width  int
	height int
}

// New creates a new root application model
func New(vfsManager *vfs.Manager) Model {
	theme := theme.DefaultTheme()
	keys := keys.DefaultKeyMap()
	scene := browser.New(vfsManager, theme, keys)

	return Model{
		vfs:          vfsManager,
		theme:        theme,
		keys:         keys,
		currentScene: SceneBrowser,
		browserScene: scene,
	}
}

// Init initializes the application
func (m Model) Init() tea.Cmd {
	return m.browserScene.Init()
}

// Update handles messages for the application
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Propagate to current scene
		switch m.currentScene {
		case SceneBrowser:
			updatedScene, cmd := m.browserScene.Update(msg)
			m.browserScene = updatedScene.(browser.Model)
			return m, cmd
		}
		return m, nil

	case tea.KeyMsg:
		// Handle global shortcuts first
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.BrowserScene):
			m.currentScene = SceneBrowser
			return m, m.browserScene.Init()

		// F2 and F3 scenes not implemented yet, so ignore
		case key.Matches(msg, m.keys.MountsScene):
			// TODO: Switch to mounts scene when implemented
			return m, nil

		case key.Matches(msg, m.keys.SettingsScene):
			// TODO: Switch to settings scene when implemented
			return m, nil
		}

	case SceneChangeMsg:
		m.currentScene = msg.Scene
		switch msg.Scene {
		case SceneBrowser:
			return m, m.browserScene.Init()
		}
		return m, nil

	case ErrorMsg:
		// Global error handling
		// For now, just pass through
		// TODO: Could show a modal or status message
		return m, nil
	}

	// Delegate to current scene
	switch m.currentScene {
	case SceneBrowser:
		updatedScene, cmd := m.browserScene.Update(msg)
		m.browserScene = updatedScene.(browser.Model)
		return m, cmd
	}

	return m, nil
}

// View renders the current scene
func (m Model) View() string {
	if m.width == 0 {
		return "Initializing Vista..."
	}

	// Render current scene
	switch m.currentScene {
	case SceneBrowser:
		return m.browserScene.View()
	}

	return "No scene active"
}
