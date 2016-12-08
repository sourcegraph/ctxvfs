package ctxvfs

import (
	"context"
	"reflect"
	"sort"
	"testing"
)

func TestListAllFiles(t *testing.T) {
	want := []string{
		"/foo/bar.txt",
		"/foo/bar/three.txt",
		"/other-top.txt",
		"/top.txt",
	}
	m := map[string][]byte{}
	for _, p := range want {
		m[p[1:]] = []byte("a")
	}
	fs := Map(m)

	got, err := ListAllFiles(context.Background(), fs)
	if err != nil {
		t.Fatal(err)
	}

	sort.Strings(got)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got != want\ngot:  %v\nwant: %v", got, want)
	}
}
