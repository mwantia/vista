package vfs

import (
	"fmt"
	"path/filepath"

	"github.com/mwantia/vfs/data"
)

// EntryFromMetadata converts data.Metadata to our Entry model
func EntryFromMetadata(path string, meta *data.Metadata) Entry {
	entry := Entry{
		Name:        meta.Key,
		Path:        filepath.Join(path, meta.Key),
		Mode:        meta.Mode,
		Size:        meta.Size,
		AccessTime:  meta.AccessTime,
		ModifyTime:  meta.ModifyTime,
		CreateTime:  meta.CreateTime,
		ContentType: meta.ContentType,
		ETag:        meta.ETag,
	}

	return entry
}

// DisplaySize returns a human-readable size string
func (e Entry) _DisplaySize() string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	size := float64(e.Size)
	switch {
	case e.Size >= GB:
		return fmt.Sprintf("%6.2f GB", size/GB)
	case e.Size >= MB:
		return fmt.Sprintf("%6.2f MB", size/MB)
	case e.Size >= KB:
		return fmt.Sprintf("%6.2f KB", size/KB)
	default:
		return fmt.Sprintf("%6d B", e.Size)
	}
}
