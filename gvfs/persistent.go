package gvfs

import "github.com/spf13/afero"

// Just pure file operations, no knowledge of "Pages" or "Frontmatter"
func WriteFileSafely(fs afero.Fs, path string, content []byte) error {
	// MkdirAll, Create, Write, Close logic goes here
	return nil
}
