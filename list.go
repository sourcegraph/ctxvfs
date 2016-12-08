package ctxvfs

import "context"

// ListAllFileser is type that implements ListAllFiles
type ListAllFileser interface {
	// ListAllFiles returns a slice of all file paths on a VFS. This is
	// only files (excludes directories). This interface exists for
	// implementors to provide a faster way to list all files than walking
	// the file tree with ReadDir/etc.
	ListAllFiles(ctx context.Context) ([]string, error)
}

// ListAllFiles returns a slice of all file paths (excludes directories).
func ListAllFiles(ctx context.Context, fs FileSystem) ([]string, error) {
	if l, ok := fs.(ListAllFileser); ok {
		return l.ListAllFiles(ctx)
	}
	var filenames []string
	w := Walk(ctx, "/", fs)
	for w.Step() {
		if err := w.Err(); err != nil {
			return nil, err
		}
		fi := w.Stat()
		if fi.Mode().IsRegular() {
			filenames = append(filenames, w.Path())
		}
	}
	return filenames, nil
}
