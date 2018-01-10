// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ctxvfs

import (
	"io/ioutil"
	"strings"
	"os"
	"path/filepath"
	"testing"
)

func TestOpenLocal(t *testing.T) {

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

	fs := NameSpace{}
	fs.Bind("/", OS("/"), "/", BindAfter)

	_, err = fs.Open(nil, file)
	if err != nil {
		t.Errorf("Cannot read local file (%s)", err)
	}

	fs = NameSpace{}
	fs.Bind("/", OS(tmpDir), "/", BindAfter)

	_, err = fs.Open(nil, "foo")
	if err != nil {
		t.Errorf("Cannot read local file (%s)", err)
	}

	_, err = fs.Open(nil, "/foo")
	if err != nil {
		t.Errorf("Cannot read local file (%s)", err)
	}

	fs = NameSpace{}
	fs.Bind("", OS(""), "", BindAfter)

	_, err = fs.Open(nil, file)
	if err != nil {
		t.Errorf("Cannot read local file (%s)", err)
	}

	_, err = fs.Open(nil, "/"+file)
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

	mapfs := Map(map[string][]byte{
		strings.TrimPrefix(file, "/"): []byte("qux"),
	})

	fs := NameSpace{}
	fs.Bind("", mapfs, "", BindBefore)
	fs.Bind("", OS(""), "", BindAfter)

	stream, err := fs.Open(nil, file)
	if err != nil {
		t.Errorf("Cannot read overlayed file (%s)", err)
	}
	content, err := ioutil.ReadAll(stream)
	if err != nil {
		t.Errorf("Cannot read overlayed file (%s)", err)
	}

	if string(content) != "qux" {
		t.Errorf("Cannot read overlayed file, expected `qux`, got `%s`", content)
	}

}
