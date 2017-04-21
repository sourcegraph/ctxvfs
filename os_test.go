// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ctxvfs_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/mapfs"
)

func TestOpen(t *testing.T) {

	tmpDir, err := ioutil.TempDir("", "vfs-localfs")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	if err = os.MkdirAll(tmpDir, 0700); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(tmpDir, "foo")
	if err = ioutil.WriteFile(file, []byte("bar"), 0600); err != nil {
		t.Fatal(err)
	}

	fs := vfs.NameSpace{}
	fs.Bind("/", vfs.OS("/"), "/", vfs.BindAfter)

	_, err = fs.Open(file)
	if err != nil {
		t.Errorf("Cannot read local file (%s)", err)
	}

	fs = vfs.NameSpace{}
	fs.Bind("/", vfs.OS(tmpDir), "/", vfs.BindAfter)

	_, err = fs.Open("foo")
	if err != nil {
		t.Errorf("Cannot read local file (%s)", err)
	}

	_, err = fs.Open("/foo")
	if err != nil {
		t.Errorf("Cannot read local file (%s)", err)
	}


}

func TestOpenOverlay(t *testing.T) {

	tmpDir, err := ioutil.TempDir("", "vfs-localfs")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	if err = os.MkdirAll(tmpDir, 0700); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(tmpDir, "foo")
	if err = ioutil.WriteFile(file, []byte("bar"), 0600); err != nil {
		t.Fatal(err)
	}

	mfs := mapfs.New(map[string]string{file: "c"})
	fs := vfs.NameSpace{}
	fs.Bind("/", mfs, "/", vfs.BindBefore)
	fs.Bind("/", vfs.OS("/"), "/", vfs.BindAfter)

	_, err = fs.Open(file)
	if err != nil {
		t.Errorf("Cannot read overlayed file (%s)", err)
	}
}
