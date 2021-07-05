package transport

import (
	"os"
	"syscall"
)

// File is a Transport for writing bytes to a file on disk.
type File struct {
	file *os.File
}

// NewFile creates a File transport for a target file path on disk.
// The file is opened in append mode, and will be created if it does not already exist.
func NewFile(path string) (Transport, error) {
	file, err := os.OpenFile(path, syscall.O_WRONLY|syscall.O_APPEND|syscall.O_CREAT, 0664)
	if err != nil {
		return nil, err
	}

	return &File{file}, nil
}

// Send writes a chunk of data to the end of the file.
func (t *File) Send(data []byte) (int, error) {
	return t.file.Write(data)
}

// Close closes the underlying file's file descriptor.
func (t *File) Close() error {
	return t.file.Close()
}
