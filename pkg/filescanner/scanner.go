package filescanner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gallyamow/go-fias-exporter/internal/model"
)

type Filter struct {
	MinSize     int64
	IncludeExts []string
	ExcludeExts []string
	IncludeDirs []string
	ExcludeDirs []string
}

// ScanDir scans the directory and returns a list of files.
// IO-bound task, no reason to use goroutines.
func ScanDir(root string, filter Filter) ([]model.FileInfo, error) {
	var files []model.FileInfo

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "skip %s: %v\n", path, err)
			return nil
		}

		if d.IsDir() {
			if isDirExcluded(d, filter) {
				return fs.SkipDir
			}
			return nil
		}

		// symlinks?

		info, err := d.Info()
		if err != nil {
			return nil
		}

		if isFileExcluded(info, filter) {
			return nil
		}

		files = append(files, model.FileInfo{
			Path: path,
			Size: info.Size(),
		})

		return nil
	})

	return files, err
}

func isFileExcluded(info fs.FileInfo, filter Filter) bool {
	// size
	if filter.MinSize != 0 {
		if info.Size() < filter.MinSize {
			return true
		}
	}

	// extension
	ext := strings.ToLower(filepath.Ext(info.Name()))

	if ext != "" {
		if !slices.Contains(filter.IncludeExts, ext) {
			return true
		}
		if slices.Contains(filter.ExcludeExts, ext) {
			return true
		}
	}

	return false
}

func isDirExcluded(d fs.DirEntry, filter Filter) bool {
	dir := strings.ToLower(d.Name())

	if !slices.Contains(filter.IncludeDirs, dir) {
		return true
	}

	if slices.Contains(filter.ExcludeDirs, dir) {
		return true
	}

	// regexp?

	return false
}
