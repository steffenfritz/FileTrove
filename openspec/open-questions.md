# Open Questions and Change Candidates

Not part of the specs. Items here are either undecided behavior (to be clarified before
they are specified) or known deviations between code and the project principles in
`config.yaml` (candidates for `/opsx:propose`). Remove an item once it is decided or
archived as a change.

## To clarify

### Stored file paths: relative or absolute?
`files.filepath` and `directories.dirpath` are stored as walked, i.e. prefixed with the
`-i` argument exactly as given. A relative `-i data` yields relative paths, while
`sessionsmd.mountpoint` is always absolute (`cmd/ftrove/main.go:102`, `:504`).
Affects the database, all exports (PREMIS `contentLocation`), `hierarchy` values and
resume. Undecided which behavior is intended.

### Timestamp format
Timestamps are stored using Go's `time.Time.String()` (e.g.
`2026-11-04 09:12:00 +0100 CET`), not RFC 3339. Deferred; not specified yet.

## Change candidates

| # | Deviation | Location |
|---|-----------|----------|
| 1 | Hash failure still inserts the file with empty hash values; only logged | `cmd/ftrove/main.go:512` |
| 2 | Unreadable regular files land in both the file list and the skipped list | `filewalk.go:33` |
| 3 | Missing ctime/btime stored as zero time `0001-01-01 …` instead of empty | `cmd/ftrove/main.go:544` |
| 4 | Entropy `0` for files > 1 GB, indistinguishable from empty files | `entropy.go:41` |
| 5 | Without a bloom filter every file is `filensrl=FALSE` instead of "not checked" | `cmd/ftrove/main.go:559` |
| 6 | TSV export writes floats with float32 precision in E-notation | `db.go:419` |
| 7 | `-t` with an unknown session UUID panics (index out of range) | `cmd/ftrove/main.go:263` |
| 8 | DB/binary compatibility check only runs when an NSRL bloom filter exists | `install.go:186` |
| 9 | Exact version match makes every new release reject existing databases, even without schema change; should be based on a schema version (see database spec, schema stability) | `version.go:28` |
| 10 | Invalid `-z` or YARA compile error leaves a session without end time | `cmd/ftrove/main.go:349` |
| 11 | Unconsumed positional arguments are ignored silently; missing `-i` is not rejected (#188) | `cmd/ftrove/main.go:51` |
| 12 | Resume does not restore the session's feature flags | `cmd/ftrove/main.go:400` |
| 13 | TSV export requires a full installation, JSONL and PREMIS do not | `cmd/ftrove/main.go:223` |
