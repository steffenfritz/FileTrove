# Proposal

## Why

Fixes #183. A scan looks for the NSRL bloom filter under `db/nsrl-<variant>.bloom`, where
`<variant>` comes from `--nsrl-variant` (default `all`), and falls back to
`db/nsrl.bloom`. After `ftrove --install <dir> --nsrl-variant modern`, the installation
holds `db/nsrl-modern.bloom`, but a scan without the flag looks for `nsrl-all.bloom` and
`nsrl.bloom`, finds neither, and runs without NSRL checks. The only trace is a warning
in `logs/filetrove.log`; without `-V` nothing appears on the console. The installation
check does not help either, because it accepts any `db/nsrl-*.bloom` as present. Every
file is then recorded with `filensrl=FALSE`, which breaks the project principle that
results must never be silently incomplete. The `--help` text describes
`--nsrl-variant` as an install-only option, so users do not discover that passing it
to the scan works around the problem.

## What Changes

- When `--nsrl-variant` is not given, a scan (`-i`) or resume (`-r`) uses the bloom
  filter that is installed in `db/`: if exactly one `db/nsrl-<variant>.bloom` exists for
  a known variant, that file is used. If several exist, the most comprehensive one is
  chosen in the fixed order `all`, `mobile`, `modern`. If none exists, the legacy
  `db/nsrl.bloom` is used as before.
- When `--nsrl-variant` is given, the scan uses `db/nsrl-<variant>.bloom` and otherwise
  the legacy `db/nsrl.bloom`, as today. It does not silently switch to another variant.
- When a scan runs without NSRL checks (no bloom filter found, or the file cannot be
  loaded), `ftrove` prints a warning to stderr in addition to the log entry, naming the
  reason. The scan continues and the session records `none` as NSRL version, as today.
- The `--help` text of `--nsrl-variant` says that the flag selects the variant to
  download during install and the variant to use during a scan.
- README / BUILDING.md: describe how the scan picks the bloom filter.

Compatibility: no DB schema change, no export format change. The `--nsrl-variant` flag
keeps its name and default; only its effect on a scan without the flag changes.
Installations with a single bloom filter or with `nsrl-all.bloom` behave as before,
except that a missing filter is now also reported on stderr. Installations that hold
only `nsrl-modern.bloom` or `nsrl-mobile.bloom` now get NSRL checks where they silently
had none before, so `filensrl` and the NSRL count change from all `FALSE`/0 to actual
results. The session's NSRL version (which includes the variant, e.g.
`2026.03.1-modern`) shows which filter was used.

Out of scope:
- `filensrl=FALSE` instead of "not checked" when no bloom filter is loaded
  (open question 5); this would change the DB contents and exports
- The database compatibility check only running when a bloom filter exists
  (open question 8)
- The mismatch between `NSRL_VARIANT` in `task dist:bundle` and `--nsrl-variant` at
  install time (dist taskfile, tracked separately); this change only makes the scan
  find whatever variant was installed
- A new flag or config file to persist the chosen variant

## Capabilities

### New Capabilities

_None._

### Modified Capabilities

- `cli`: the "Working directory layout" requirement defines how a scan selects the NSRL
  bloom filter when several or non-default variants are installed, and the "Scan
  session" requirement reports a scan without NSRL checks on stderr.

## Impact

- `cmd/ftrove/main.go`: bloom filter selection before scan/resume; console warning when
  NSRL checks are disabled; `--nsrl-variant` help text
- `nsrl.go` (or a helper in `cmd/ftrove`): a small, testable function that resolves the
  bloom filter path from the files in `db/` and the optional explicit variant, with
  table-driven tests
- `README.md`, `BUILDING.md`: NSRL section
