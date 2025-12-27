package vfs

import (
	"context"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"path/filepath"
	"strings"

	"github.com/eliukblau/pixterm/pkg/ansimage"
	"github.com/mwantia/vfs"
	"github.com/mwantia/vfs/data"
	"github.com/mwantia/vista/internal/preview"
	"golang.org/x/image/draw"
)

// PreviewType represents how a file should be previewed
type PreviewType int

const (
	PreviewMetadata    PreviewType = iota // File details only
	PreviewText                           // UTF-8 text files
	PreviewImage                          // Images rendered as ANSI art
	PreviewBinary                         // Hex dump
	PreviewUnsupported                    // Unsupported types
)

// FileTypeInfo contains information about how to preview a file
type FileTypeInfo struct {
	Type        PreviewType
	Description string
}

// DetectFileType determines the appropriate preview type for a file
// Uses extension-based detection only (no content validation)
func DetectFileType(filename string) FileTypeInfo {
	ext := strings.ToLower(filepath.Ext(filename))

	// Image files
	imageExts := map[string]bool{
		".png": true, ".jpg": true, ".jpeg": true,
		".gif": true, ".bmp": true, ".webp": true,
		".svg": true, ".ico": true,
	}
	if imageExts[ext] {
		return FileTypeInfo{
			Type:        PreviewBinary,
			Description: "Image file",
		}
	}

	// Text files
	textExts := map[string]bool{
		".txt": true, ".md": true, ".go": true, ".js": true,
		".ts": true, ".tsx": true, ".jsx": true,
		".py": true, ".java": true, ".c": true, ".cpp": true,
		".h": true, ".hpp": true, ".rs": true, ".sh": true,
		".bash": true, ".zsh": true, ".fish": true,
		".json": true, ".xml": true, ".yaml": true, ".yml": true,
		".toml": true, ".ini": true, ".cfg": true, ".conf": true,
		".html": true, ".css": true, ".scss": true, ".sass": true,
		".sql": true, ".log": true, ".csv": true, ".tsv": true,
		".gitignore": true, ".dockerfile": true, ".env": true,
		".rb": true, ".php": true, ".vue": true,
	}
	if textExts[ext] {
		return FileTypeInfo{
			Type:        PreviewText,
			Description: "Text file",
		}
	}

	// Binary files that shouldn't be previewed as text
	binaryExts := map[string]bool{
		".zip": true, ".gz": true, ".tar": true, ".bz2": true,
		".7z": true, ".rar": true, ".xz": true, ".tgz": true,
		".exe": true, ".dll": true, ".so": true, ".dylib": true,
		".bin": true, ".dat": true, ".db": true, ".sqlite": true,
		".pdf": true, ".doc": true, ".docx": true, ".xls": true,
		".xlsx": true, ".ppt": true, ".pptx": true,
		".mp3": true, ".mp4": true, ".avi": true, ".mkv": true,
		".wav": true, ".flac": true, ".ogg": true, ".mov": true,
		".wmv": true, ".flv": true,
		".a": true, ".o": true,
	}
	if binaryExts[ext] {
		return FileTypeInfo{
			Type:        PreviewBinary,
			Description: "Binary file",
		}
	}

	// Default to text for unknown extensions
	return FileTypeInfo{
		Type:        PreviewText,
		Description: "Unknown text file",
	}
}

// GenerateTextPreview creates a text preview of a file with syntax highlighting
func GenerateTextPreview(vfs vfs.VirtualFileSystem, ctx context.Context, path string, filename string, maxBytes int) (string, error) {
	file, err := vfs.OpenFile(ctx, path, data.AccessModeRead)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read up to maxBytes
	buf := make([]byte, maxBytes)
	n, err := io.ReadFull(file, buf)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return "", err
	}
	buf = buf[:n]

	// Apply syntax highlighting
	highlighter := preview.NewHighlighter()
	highlighted := highlighter.Highlight(buf, filename)

	return highlighted, nil
}

// GenerateImagePreview creates an ANSI art preview of an image
func GenerateImagePreview(vfs vfs.VirtualFileSystem, ctx context.Context, path string, width, height int) (string, error) {
	// First check file size to prevent loading huge images
	stat, err := vfs.StatMetadata(ctx, path)
	if err != nil {
		return "", fmt.Errorf("failed to stat image: %w", err)
	}

	// Skip images larger than 5MB - too slow to render
	const maxImageSize = 5 * 1024 * 1024
	if stat.Size > maxImageSize {
		return fmt.Sprintf("[Image too large to preview: %.1f MB]\n\nUse a dedicated image viewer for files > 5MB",
			float64(stat.Size)/(1024*1024)), nil
	}

	file, err := vfs.OpenFile(ctx, path, data.AccessModeRead)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	// Decode image
	img, format, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Calculate scaled dimensions
	maxHeight := height - 10 // Leave room for header
	maxWidth := width - 4    // Leave room for padding

	bounds := img.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	scale := minFloat64(float64(maxWidth)/float64(imgWidth), float64(maxHeight)/float64(imgHeight))
	if scale > 1.0 {
		scale = 1.0 // don't upscale
	}
	newW := int(float64(imgWidth) * scale)
	newH := int(float64(imgHeight) * scale)

	// Scale image
	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)

	// Create ANSI image
	ansImg, err := ansimage.NewFromImage(dst, color.Transparent, ansimage.NoDithering)
	if err != nil {
		return "", fmt.Errorf("failed to create ANSI image: %w", err)
	}

	rendered := ansImg.Render()
	header := fmt.Sprintf("Image: %s format, %dx%d pixels\n\n", format, imgWidth, imgHeight)

	return header + rendered, nil
}

// GenerateBinaryPreview creates a hex dump preview of a binary file
func GenerateBinaryPreview(vfs vfs.VirtualFileSystem, ctx context.Context, path string, maxBytes int) (string, error) {
	file, err := vfs.OpenFile(ctx, path, data.AccessModeRead)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Get file info
	stat, err := vfs.StatMetadata(ctx, path)
	if err != nil {
		return "", err
	}

	// Read up to maxBytes
	buf := make([]byte, maxBytes)
	n, err := io.ReadFull(file, buf)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return "", err
	}
	buf = buf[:n]

	var preview strings.Builder
	preview.WriteString(fmt.Sprintf("Binary file: %s\n", filepath.Base(path)))
	preview.WriteString(fmt.Sprintf("Size: %d bytes\n\n", stat.Size))
	preview.WriteString("Hex dump (first 512 bytes):\n")
	preview.WriteString(strings.Repeat("-", 60))
	preview.WriteString("\n")

	// Create hex dump
	dumper := hex.Dumper(&preview)
	dumper.Write(buf)
	dumper.Close()

	if stat.Size > int64(n) {
		preview.WriteString("\n... (truncated)")
	}

	return preview.String(), nil
}

// GenerateMetadataPreview creates a preview showing file metadata
func GenerateMetadataPreview(entry *Entry) string {
	var preview strings.Builder

	preview.WriteString(fmt.Sprintf("Name: %s\n", entry.Name))
	preview.WriteString(fmt.Sprintf("Size: %s\n", formatSize(entry.Size)))
	preview.WriteString(fmt.Sprintf("Modified: %s\n", entry.ModifyTime.Format("2006-01-02 15:04:05")))
	preview.WriteString(fmt.Sprintf("Mode: %s\n", entry.Mode))

	if entry.Mode.IsDir() {
		preview.WriteString("\nType: Directory\n")
	} else {
		ext := strings.ToLower(filepath.Ext(entry.Name))
		if ext != "" {
			preview.WriteString(fmt.Sprintf("\nExtension: %s\n", ext))
		}
		fileType := DetectFileType(entry.Name)
		preview.WriteString(fmt.Sprintf("File Type: %s\n", fileType.Description))
	}

	return preview.String()
}

// GeneratePreview generates an appropriate preview for any file
func GeneratePreview(vfs vfs.VirtualFileSystem, ctx context.Context, entry *Entry, width, height int) (string, error) {
	// For directories, show metadata only
	if entry.Mode.IsDir() {
		return GenerateMetadataPreview(entry), nil
	}

	fileInfo := DetectFileType(entry.Name)

	switch fileInfo.Type {
	case PreviewText:
		// Pass entry.Name for syntax highlighting (works with or without extension)
		content, err := GenerateTextPreview(vfs, ctx, entry.Path, entry.Name, 10240) // 10KB
		if err != nil {
			return "", err
		}
		return content, nil

	case PreviewImage:
		content, err := GenerateImagePreview(vfs, ctx, entry.Path, width, height)
		if err != nil {
			// If image rendering fails, fall back to binary preview
			return GenerateBinaryPreview(vfs, ctx, entry.Path, 512)
		}
		return content, nil

	case PreviewBinary:
		return GenerateBinaryPreview(vfs, ctx, entry.Path, 512) // 512 bytes hex dump

	case PreviewMetadata:
		return GenerateMetadataPreview(entry), nil

	case PreviewUnsupported:
		return fmt.Sprintf("[Cannot preview %s files]", fileInfo.Description), nil

	default:
		return "[Unknown file type]", nil
	}
}

// Helper functions
func minFloat64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
