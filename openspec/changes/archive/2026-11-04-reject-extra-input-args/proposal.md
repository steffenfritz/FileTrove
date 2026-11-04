# Proposal

## Why

Fixes #188. When a user runs `ftrove -i dir/*`, the shell expands the glob, `-i` takes
only the first match and every other path is silently ignored. The user gets a
drastically incomplete scan without any warning, which violates the project principle
that results must never be silently incomplete. Even a glob that ftrove accepted would
not be equivalent to scanning `dir/`, because shell globs do not match hidden files.

## What Changes

- `ftrove` aborts with exit code 1 when positional arguments remain after flag parsing.
  The check runs before anything else happens, so no session, database row, log file
  or debug package is created.
- The error message lists the unexpected arguments, states that `-i` expects exactly
  one directory, and tells the user to pass the directory itself (e.g.
  `ftrove -i corpus/`), noting that a shell glob drops hidden files.
- The check applies to every mode (scan, resume, install, exports, listings, `-v`),
  since no mode accepts positional arguments.
- README: note under "Running a scan" that `-i` takes one directory and that globs
  must not be used.

Compatibility: no DB schema or export format change. CLI flags are unchanged. Command
lines that contain stray positional arguments, which were ignored before, now fail;
this is intended, as those invocations produced incomplete or unintended results.

Out of scope:
- Rejecting a scan without `-i` (separate change; candidate 11 in
  `openspec/open-questions.md`)
- Accepting multiple input directories or expanding globs internally
- The missing records in run 2 of #188 (4/4 files processed, but only the session in
  the JSONL export); needs separate investigation
- Printing the scan summary to the console (#185, #189)

## Capabilities

### New Capabilities

_None._

### Modified Capabilities

- `cli`: adds a requirement that unconsumed positional arguments are rejected before
  any work is done.

## Impact

- `cmd/ftrove/main.go`: argument check directly after flag parsing
- `cmd/ftrove`: unexported helper for the check and its message, with the package's
  first unit tests; the public library API of package `filetrove` is unchanged
- `Taskfile.yml`: `task test` also runs the tests in `cmd/ftrove`
- `README.md`: usage note
