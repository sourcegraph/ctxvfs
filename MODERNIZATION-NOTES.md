# Modernization notes

Judgment calls made while converting this repository to Go modules / Go 1.26:

- **Pinned to a deprecated dependency: `golang.org/x/tools/godoc v0.1.0-deprecated`.**
  `strip_context.go` needs `golang.org/x/tools/godoc/vfs`. That package was moved out of
  the main `golang.org/x/tools` module into its own module, which upstream publishes only
  as `v0.1.0-deprecated`; there is no maintained successor. It was not inlined under
  `internal/` because `StripContext` returns the exported `vfs.FileSystem` and
  `RootType` returns the exported named type `vfs.RootType` — a local copy would be a
  distinct type, so the returned value would no longer satisfy godoc's `vfs.FileSystem`
  and every caller that passes it to `vfs.NameSpace` would break.

- **No Windows CI coverage.** `appveyor.yml` was removed and the new workflow runs on
  `ubuntu-latest` only, so `os_windows.go` and `os_windows_test.go` (`TestWindowsMount`)
  are no longer built or run anywhere in CI — the workflow's `go build`/`go vet` steps use
  the host `GOOS`. Both files do still compile and vet cleanly under
  `GOOS=windows go build ./...` / `GOOS=windows go vet ./...`, verified locally, but that
  is not enforced automatically.

- **`osFS.ReadDir` now stats each entry.** `ioutil.ReadDir` was replaced with `os.ReadDir`
  plus a `DirEntry.Info()` call per entry to keep returning `[]os.FileInfo`. Ordering and
  the `Lstat`-based info are unchanged, but this costs one extra syscall per entry versus
  returning `[]os.DirEntry`, which would have been an API break.
