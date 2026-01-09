package vfs

import (
	"strings"

	"github.com/mwantia/vfs/data"
)

// EntryFromMetadata converts data.Metadata to our Entry model
func EntryFromMetadata(path string, meta data.Metadata) Entry {
	entry := Entry{
		Name:        strings.TrimPrefix(meta.Key, path+"/"),
		Path:        meta.Key,
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
