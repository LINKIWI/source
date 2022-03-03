package disk

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/xattr"
)

// Client is a client for simple file operations.
type Client struct {
	root string
}

// New creates a Client rooted at the specified directory.
func New(root string) (*Client, error) {
	resolved, err := filepath.EvalSymlinks(root)
	if err != nil {
		return nil, err
	}

	return &Client{resolved}, nil
}

// StatFile stats a file.
func (c *Client) StatFile(path string) (*Stat, error) {
	resolved := filepath.Join(c.root, path)
	attributes := make(map[string]string)

	stat, err := os.Stat(resolved)
	if err != nil {
		return nil, err
	}

	attrNames, err := xattr.List(resolved)
	if err != nil {
		return nil, err
	}

	for _, attr := range attrNames {
		value, err := xattr.Get(resolved, attr)
		if err != nil {
			return nil, err
		}

		attributes[attr] = string(value)
	}

	return &Stat{
		Path:         path,
		BaseName:     stat.Name(),
		Size:         uint64(stat.Size()),
		Mode:         stat.Mode(),
		ModifiedTime: stat.ModTime(),
		Attributes:   attributes,
	}, nil
}

// ListFiles lists the contents of a directory with stat information of all children.
func (c *Client) ListFiles(path string, recursive bool) ([]*Stat, error) {
	var stats []*Stat

	file, err := os.Open(filepath.Join(c.root, path))
	if err != nil {
		return nil, err
	}

	entries, err := file.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		stat, err := c.StatFile(filepath.Join(path, entry))
		if err != nil {
			return nil, err
		}

		switch {
		case stat.Mode.IsDir() && recursive:
			children, err := c.ListFiles(filepath.Join(path, entry), recursive)
			if err != nil {
				return nil, err
			}

			stats = append(stats, children...)
		case stat.Mode.IsRegular():
			stats = append(stats, stat)
		case stat.Mode.IsDir():
			stat.Size = 0
			stats = append(stats, stat)
		}
	}

	return stats, nil
}

// ReadFile opens a file for reading. It is expected that the caller closes the file.
func (c *Client) ReadFile(path string) (io.ReadSeekCloser, *Stat, error) {
	file, err := os.Open(filepath.Join(c.root, path))
	if err != nil {
		return nil, nil, err
	}

	stat, err := c.StatFile(path)
	if err != nil {
		return nil, nil, err
	}

	return file, stat, nil
}

// WriteFile opens a file for writing. It is expected that the caller closes the file.
func (c *Client) WriteFile(path string) (WriteSeekCloser, error) {
	return os.OpenFile(filepath.Join(c.root, path), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
}

// WriteAttributes writes a map of attributes to an existing file using xattr syscalls.
func (c *Client) WriteAttributes(path string, attributes map[string]string) error {
	file, err := os.OpenFile(filepath.Join(c.root, path), os.O_WRONLY, 0)
	if err != nil {
		return err
	}

	defer file.Close()

	for attr, value := range attributes {
		if err := xattr.FSet(file, attr, []byte(value)); err != nil {
			return err
		}
	}

	return nil
}

// DeleteFile deletes a file.
func (c *Client) DeleteFile(path string) error {
	return os.Remove(filepath.Join(c.root, path))
}

// CreateDirectory creates a directory tree at the specified path, including all parents as needed.
func (c *Client) CreateDirectory(path string) error {
	return os.MkdirAll(filepath.Join(c.root, path), 0775)
}
