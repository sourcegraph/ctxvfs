// +build windows

package ctxvfs

import (
	"testing"
)

func TestWindowsMount(t *testing.T) {

	// keys are mount points,
	// values contain expected results for osFS.resolve(key)
	tests := map[string]map[string]string{
		"": {
			"C:/folder/file":  "C:\\folder\\file",
			"D:/folder/file":  "D:\\folder\\file",
			"/C:/folder/file": "C:\\folder\\file",
			"/folder/file":    "",
			"file":            "",
			"/file":           "",
		},
		"C:": {
			"C:/folder/file":  "C:\\folder\\file",
			"/C:/folder/file": "C:\\folder\\file",
			"/folder/file":    "C:\\folder\\file",
			"file":            "C:\\file",
			"/file":           "C:\\file",
			"D:/folder/file":  "",
		},
		"/C:": {
			"C:/folder/file":  "C:\\folder\\file",
			"/C:/folder/file": "C:\\folder\\file",
			"/folder/file":    "C:\\folder\\file",
			"file":            "C:\\file",
			"/file":           "C:\\file",
			"D:/folder/file":  "",
		},
		"C:/folder": {
			"C:/folder/file":           "C:\\folder\\file",
			"/C:/folder/file":          "C:\\folder\\file",
			"/folder/file":             "C:\\folder\\folder\\file",
			"file":                     "C:\\folder\\file",
			"/file":                    "C:\\folder\\file",
			"C:/anoherfolder/file":     "",
			"C:\\anoherfolder\\file":   "",
			"/C:\\anoherfolder/file":   "",
			"foo/bar/../../file":       "C:\\folder\\file",
			"foo/bar\\\\\\..\\../file": "C:\\folder\\file",
			"c:/Folder/fiLe":           "c:\\Folder\\fiLe",
		},
		"/C:/folder": {
			"C:/folder/file":         "C:\\folder\\file",
			"/C:/folder/file":        "C:\\folder\\file",
			"/folder/file":           "C:\\folder\\folder\\file",
			"file":                   "C:\\folder\\file",
			"/file":                  "C:\\folder\\file",
			"C:/anoherfolder/file":   "",
			"C:\\anoherfolder/file":  "",
			"/C:\\anoherfolder/file": "",
		},
		"C:/folder/": {
			"C:/folder/file":         "C:\\folder\\file",
			"/C:/folder/file":        "C:\\folder\\file",
			"/folder/file":           "C:\\folder\\folder\\file",
			"file":                   "C:\\folder\\file",
			"/file":                  "C:\\folder\\file",
			"C:/anoherfolder/file":   "",
			"C:\\anoherfolder/file":  "",
			"/C:\\anoherfolder/file": "",
		},
	}

	for mount, expectations := range tests {
		fs := osFS(mount)
		for path, expected := range expectations {
			actual, _ := fs.resolve(path)
			if actual != expected {
				t.Errorf("expected `%s`, got `%s` while resolving `%s` on `%s`", expected, actual, path, mount)
			}
		}
	}
}
