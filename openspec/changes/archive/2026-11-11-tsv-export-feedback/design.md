# Design

## Context

The TSV export in `cmd/ftrove/main.go` runs after `CheckInstall` and after the logger
has been switched to `logs/filetrove.log` (plus stdout with `-V`). It calls up to seven
`ft.Export*TSV(uuid)` functions in `db.go`. Each one creates `<uuid>_<table>.tsv` in the
working directory and returns only an error. The banner and the logger used before the
switch both write to stderr. With `-V`, stdout carries log lines.

## Goals / Non-Goals

**Goals:**
- Report success and failure on the console without changing the library API or the
  TSV files.
- Keep the report logic out of `main()` so it can be unit tested with buffers.

**Non-Goals:**
- Restructuring the export block or the order of the checks in `main()`
- Cleaning up partially written files

## Decisions

**Paths on stdout, header and errors on stderr.** Stdout holds only data, as it does for
`-j` and `-P`, so `ftrove -t <uuid> | xargs …` works. The human-readable header goes to
stderr next to the banner. Alternative: one block of text on stdout. Rejected, because
it cannot be piped without filtering.

**Build file names in `cmd/ftrove`, keep `ft.Export*TSV` unchanged.** An unexported
helper `tsvFileName(uuid, table)` returns `<uuid>_<table>.tsv`, matching what `db.go`
creates. After each export call succeeds, main appends that name to a list.
Alternative: have each `Export*TSV` return the path it wrote. That is cleaner, but it
changes seven public signatures and the tests in `db_test.go` for a reporting concern.
The only callers are in `main.go`, so the change could be made later if the duplicated
naming becomes a problem.

**Absolute paths via `filepath.Abs`.** The files are created relative to the working
directory, so `filepath.Abs` gives the right absolute path. If `Abs` fails (the working
directory was removed), the relative name is printed instead of failing an export that
did succeed.

**Report helpers take `io.Writer`.** `reportTSVExport(stdout, stderr, uuid, files)` and
`reportTSVExportFailure(stderr, uuid, failed, files, err)` write the console output. Main
passes `os.Stdout` and `os.Stderr`, and tests pass buffers. Logging through `logger`
stays as it is, so the log file keeps its current entries.

**Installation check prints to stderr.** `CheckInstall` in `install.go` runs before
the TSV export and prints its messages with `fmt.Println`. In an installation without
an NSRL bloom filter, this puts an error line on stdout ahead of the paths. All its
messages move to `os.Stderr`, the same stream the banner and the early logger use.
Alternative: skip the bloom filter check for `-t`, since the export does not need it.
Rejected, because the other messages (missing database, version mismatch) would still
go to stdout, and which checks run in which mode is open question 13.

**Failure message names the failed file.** The stderr message names the file that was
being written when the error occurred, in addition to the files that were completed.
That file may exist in a partially written state, and the message says so.

## Risks / Trade-offs

- [File names in `db.go` and `cmd/ftrove` drift apart] → `tsvFileName` is the single
  place in `cmd/ftrove`. A comment there points to `db.go`. A mismatch would show up in
  the README check in tasks 2.1/2.2, since the printed paths would not exist.
- [`-V` mixes log lines into stdout] → Documented in the README. Changing where `-V`
  writes is a `cli` spec change and out of scope.
- [The panic on an unknown UUID (open question 7) happens before any report] →
  Unchanged by this change. It stays listed in `openspec/open-questions.md`.
