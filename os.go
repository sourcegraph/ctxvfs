// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ctxvfs

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	pathpkg "path"
	"path/filepath"
)

// OS returns an implementation of FileSystem reading from the
// tree rooted at root.  Recording a root is convenient everywhere
// but necessary on Windows, because the slash-separated path
// passed to Open has no way to specify a drive letter.  Using a root
// lets code refer to OS(`c:\`), OS(`d:\`) and so on.
//
// TODO(sqs): The ctx parameter in the FileSystem methods is currently
// ignored.
func OS(root string) FileSystem {
	return osFS(root)
}

type osFS string

func (root osFS) String() string { return "os(" + string(root) + ")" }

func (root osFS) resolve(path string) string {
	// We special case "" to mean the OS root filesystem. We do this so on
	// windows we can have a VFS which allows access to /C:/foo, /D:/foo,
	// etc.
	if root == "" {
		return path
	}

	// Clean the path so that it cannot possibly begin with ../.
	// If it did, the result of filepath.Join would be outside the
	// tree rooted at root.  We probably won't ever see a path
	// with .. in it, but be safe anyway.
	path = pathpkg.Clean("/" + path)

	return filepath.Join(string(root), path)
}

func (root osFS) Open(ctx context.Context, path string) (ReadSeekCloser, error) {
	f, err := os.Open(root.resolve(path))
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}
	if fi.IsDir() {
		f.Close()
		return nil, fmt.Errorf("Open: %s is a directory", path)
	}
	return f, nil
}

func (root osFS) Lstat(ctx context.Context, path string) (os.FileInfo, error) {
	return os.Lstat(root.resolve(path))
}

func (root osFS) Stat(ctx context.Context, path string) (os.FileInfo, error) {
	return os.Stat(root.resolve(path))
}

func (root osFS) ReadDir(ctx context.Context, path string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(root.resolve(path)) // is sorted
}
