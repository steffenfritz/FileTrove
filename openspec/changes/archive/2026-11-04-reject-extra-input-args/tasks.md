# Tasks

## 1. Argument check

- [x] 1.1 Extend the `test` task in `Taskfile.yml` from `go test -v` to `go test -v . ./cmd/ftrove` so tests of the `ftrove` binary package run as well; verify `task test` still passes and lists package `github.com/steffenfritz/FileTrove/cmd/ftrove`
- [x] 1.2 Add an unexported helper in `cmd/ftrove` (e.g. `checkUnexpectedArgs` in `args.go`) that takes the remaining positional arguments and the `-i` value and returns nil for no arguments, otherwise an error whose message lists the unexpected arguments, states that `-i` expects exactly one directory, mentions that shell globs miss hidden files, and suggests `ftrove -i <parent of -i>` when `-i` is set; verify with table-driven tests in `cmd/ftrove/args_test.go` covering no args, glob expansion (`corpus/a` + `corpus/b corpus/c` → suggests `corpus`), args without `-i`, and a trailing-slash `-i` value
- [x] 1.3 Call the helper in `cmd/ftrove/main.go` directly after `flag.Parse()`, before the debug package, banner and any database access; print the message to stderr and exit with code 1; verify with `task test` and by building `ftrove` and running `ftrove -i testdata/*` (exit code 1, no new row in `sessionsmd`) and `ftrove -i testdata/` (scan runs)

## 2. Documentation

- [x] 2.1 Add a note under "Running a scan" in `README.md` that `-i` takes exactly one directory, that it is scanned recursively including hidden files, and that globs such as `dir/*` are rejected; verify the documented command matches the implemented behavior

## Workflow follow-up

- Run `/opsx:verify`, then `/opsx:archive` to merge the delta into `openspec/specs/cli/spec.md`.
- Remove the #188 part of candidate 11 from `openspec/open-questions.md` on archive; the missing `-i` part stays open.
