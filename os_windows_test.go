// +build windows

package ctxvfs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWindowsMount(t *testing.T) {

	tmpDir, err := ioutil.TempDir("", "vfs-localfs")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	if err = os.MkdirAll(tmpDir, 0700); err != nil {
		t.Fatal(err)
	}
	tmpFile := filepath.Join(tmpDir, "foo")
	if err = ioutil.WriteFile(tmpFile, []byte("bar"), 0600); err != nil {
		t.Fatal(err)
	}

	fStat, _ := os.Stat(tmpFile)

	tmpFile = filepath.ToSlash(tmpFile)
	tmpDir = filepath.ToSlash(tmpDir)

	volname := filepath.VolumeName(tmpFile)

	var (
		empty       = ""           // - empty root, OS("")
		volume      = volname      // - volume name, OS("C:")
		slashvolume = "/" + volume // - volume name prefixed with slash, OS("/C:")
		dir         = tmpDir       // - directory, OS("C:/foo/bar")
		dirslash    = tmpDir + "/" // - directory followed by slash, OS("C:/foo/bar/")
		slashdir    = "/" + tmpDir // - directory prefixed with slash, OS("/C:/foo/bar")
	)

	var (
		file           = tmpFile                             // absolute file path, C:/foo/bar
		slashfile      = "/" + tmpFile                       // absolute file path prefixed with slash, /C:/foo/bar
		relfile        = strings.TrimPrefix(tmpFile, volume) // path w/o volume name, /foo/bar
		shortfile      = "foo"                               // relative path, foo
		slashshortfile = "/foo"                              // relative path prefixed with slash, /foo
	)

	// keys are mount points,
	// values indicate if OS.Stat(path) should succeed for a given path
	tests := map[string]map[string]bool{
		empty: {
			file:           true,
			slashfile:      true,
			relfile:        false,
			shortfile:      false,
			slashshortfile: false,
		},
		volume: {
			file:           false,
			slashfile:      false,
			relfile:        true,
			shortfile:      false,
			slashshortfile: false,
		},
		slashvolume: {
			file:           false,
			slashfile:      false,
			relfile:        true,
			shortfile:      false,
			slashshortfile: false,
		},
		dir: {
			file:           false,
			slashfile:      false,
			relfile:        false,
			shortfile:      true,
			slashshortfile: true,
		},
		slashdir: {
			file:           false,
			slashfile:      false,
			relfile:        false,
			shortfile:      true,
			slashshortfile: true,
		},
		dirslash: {
			file:           false,
			slashfile:      false,
			relfile:        false,
			shortfile:      true,
			slashshortfile: true,
		},
	}

	for mount, expectations := range tests {
		fs := OS(mount)
		for path, expected := range expectations {
			stat, err := fs.Stat(nil, path)
			actual := err == nil
			if expected != actual {
				var expectation string
				if actual {
					expectation = "not fail"
				} else {
					expectation = "fail"
				}
				t.Errorf("stat `%s` on `%s` expected to %s", path, mount, expectation)
				continue
			}
			if actual && !os.SameFile(fStat, stat) {
				t.Errorf("stat `%s` on `%s` points to wrong file", path, mount)
			}

		}
	}
}
