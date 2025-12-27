package vfs

import (
	"fmt"
	"time"

	"github.com/mwantia/vfs/data"
	"github.com/mwantia/vista/internal/resources/theme"
)

type Entry struct {
	Name        string
	Path        string
	Mode        data.FileMode
	Size        int64
	AccessTime  time.Time
	ModifyTime  time.Time
	CreateTime  time.Time
	ContentType data.ContentType
	ETag        string
}

func (e Entry) DisplayName() string {
	return e.Name
}

func (e *Entry) DisplaySize() string {
	if e.Mode.IsMount() {
		return "<MNT>"
	}

	if e.Mode.IsDir() {
		return "<DIR>"
	}

	if e.Mode.IsDevice() {
		return "<DVC>"
	}

	const unit = 1024
	if e.Size < unit {
		return fmt.Sprintf("%d B", e.Size)
	}

	div, exp := int64(unit), 0
	for n := e.Size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(e.Size)/float64(div), "KMGTPE"[exp])
}

func (e *Entry) DisplayMode() string {
	return e.Mode.String()
}

func (e *Entry) DisplayContentType() string {
	if e.Mode.IsRegular() {
		if e.ContentType != "" {
			return string(e.ContentType)
		}

		return data.ContentTypeApplicationStream
	}

	return ""
}

func (e *Entry) GetETag() string {
	return e.ETag
}

func (e *Entry) DisplayModifyTime() string {
	return e.ModifyTime.Format("2006-01-02 15:04:05")
}

// Icon returns the appropriate icon for the entry
func (e Entry) Icon() string {
	if e.Mode.IsMount() {
		return theme.MountFileIcon
	}

	if e.Mode.IsDir() {
		return theme.FolderFileIcon
	}

	if e.Mode.IsDevice() {
		return theme.DeviceFileIcon
	}

	switch {
	case theme.IsTextFile(e.Name):
		return theme.TextFileIcon
	case theme.IsImageFile(e.Name):
		return theme.ImageFileIcon
	case theme.IsVideoFile(e.Name):
		return theme.VideoFileIcon
	case theme.IsArchiveFile(e.Name):
		return theme.ArchiveFileIcon
	case theme.IsCodeFile(e.Name):
		return theme.CodeFileIcon
	default:
		return theme.DefaultFileIcon
	}
}
