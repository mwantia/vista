package theme

import "strings"

// Icon constants for different file types
const (
	TextFileIcon    string = "📑"
	ImageFileIcon   string = "🖼️"
	VideoFileIcon   string = "🎞️"
	ArchiveFileIcon string = "📦"
	CodeFileIcon    string = "📇"
	MountFileIcon   string = "🗃️"
	DeviceFileIcon  string = "💽"
	FolderFileIcon  string = "📂"
	DefaultFileIcon string = "📄"
	VirtualFileIcon string = "🗒️"
)

// IsTextFile checks if a file is a text file based on extension
func IsTextFile(name string) bool {
	exts := []string{".txt", ".md", ".log", ".conf", ".cfg", ".ini", ".yaml", ".yml", ".json", ".xml"}
	return hasExtension(name, exts)
}

// IsImageFile checks if a file is an image based on extension
func IsImageFile(name string) bool {
	exts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".webp", ".ico"}
	return hasExtension(name, exts)
}

// IsVideoFile checks if a file is a video based on extension
func IsVideoFile(name string) bool {
	exts := []string{".mp4", ".avi", ".mkv", ".mov", ".wmv", ".flv", ".webm"}
	return hasExtension(name, exts)
}

// IsArchiveFile checks if a file is an archive based on extension
func IsArchiveFile(name string) bool {
	exts := []string{".zip", ".tar", ".gz", ".bz2", ".7z", ".rar", ".tgz", ".tar.gz"}
	return hasExtension(name, exts)
}

// IsCodeFile checks if a file is a source code file based on extension
func IsCodeFile(name string) bool {
	exts := []string{".go", ".js", ".ts", ".py", ".java", ".c", ".cpp", ".h", ".hpp", ".rs", ".rb", ".php", ".sh", ".bash", ".sql", ".css", ".html", ".vue", ".jsx", ".tsx"}
	return hasExtension(name, exts)
}

// hasExtension checks if a filename has any of the given extensions (case-insensitive)
func hasExtension(name string, exts []string) bool {
	nameLower := strings.ToLower(name)
	for _, ext := range exts {
		if strings.HasSuffix(nameLower, ext) {
			return true
		}
	}
	return false
}
