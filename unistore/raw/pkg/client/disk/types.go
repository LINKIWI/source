package disk

import (
	"io"
	"io/fs"
	"time"
)

// WriteSeekCloser mimics the interfaces in standard library package io for an interface that embeds
// io.Writer, io.Seeker, and io.Closer.
type WriteSeekCloser interface {
	io.Writer
	io.Seeker
	io.Closer
}

// Stat is an abstraction over a fs.FileInfo that also includes the file's relative path.
type Stat struct {
	Path         string
	BaseName     string
	Size         uint64
	Mode         fs.FileMode
	ModifiedTime time.Time
	Attributes   map[string]string
}
