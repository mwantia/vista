package app

// SceneType represents the type of scene
type SceneType int

const (
	SceneBrowser SceneType = iota
	SceneMounts
	SceneSettings
)

// SceneChangeMsg requests a scene change
type SceneChangeMsg struct {
	Scene SceneType
}

// ErrorMsg represents a global error
type ErrorMsg struct {
	Error   error
	Context string
}
