# Proposal

## Why

Fixes #185. `ftrove -t <uuid>` writes its TSV files silently. Before the export starts,
the logger is redirected to `logs/filetrove.log`, so the console shows only the banner
and "Redirecting logs to logs/filetrove.log". Users cannot tell whether the export
worked or where the files are, and the redirect notice makes it look like something
went wrong. Errors are hidden in the same way: a failed export exits with code 1 but
prints nothing to the console, and any files written before the failure stay behind
with no warning. That breaks the project principle that results must never be silently
incomplete.

## What Changes

- After a successful TSV export, `ftrove` prints the absolute path of every TSV file it
  wrote to stdout, one per line, in the order the files were written. Nothing else is
  written to stdout (without `-V`), so the output can be piped, e.g. into `xargs`.
- In addition, a short header naming the session and the number of files is printed to
  stderr, where the banner already goes.
- If the TSV export fails, `ftrove` prints the error to stderr, says the export is
  incomplete, and lists the files written so far. The exit code stays 1.
- The installation check that runs before the TSV export and before every scan prints
  its messages (missing signature file, database or NSRL bloom filter, version
  mismatch) to stderr instead of stdout. Without this, an installation without a bloom
  filter puts `ERROR: nsrl bloom filter does not exist.` in front of the paths on stdout.
- README: "Viewing results" shows the new output and an example of piping it.

Compatibility: no DB schema change. The contents and names of the TSV files are
unchanged. CLI flags are unchanged. The only new output of `-t` is on stdout and
stderr, where it printed nothing before. The installation check messages move from
stdout to stderr. On a terminal they look the same, but scripts that read them from
stdout must read stderr instead. With `-V`, log lines still also go to stdout
(see the `cli` spec, Log output), so the output is only cleanly pipeable without `-V`.

Out of scope:
- Panic when `-t` gets an unknown session UUID (open question 7)
- TSV export requiring a full installation, unlike JSONL and PREMIS (open question 13)
- Deleting partially written files after a failure
- Showing the report or session ID at scan time (#189)
- Changes to the TSV documentation beyond the new console output, and changing the
  tutorial to use JSONL (alternatives listed in #185)

## Capabilities

### New Capabilities

_None._

### Modified Capabilities

- `export`: adds requirements for console feedback of the TSV export, with the written file
  paths on stdout on success and the error plus the files already written on stderr on
  failure.

## Impact

- `cmd/ftrove/main.go`: the TSV export block records each file it writes and reports
  success or failure on the console
- `cmd/ftrove`: unexported helpers for file names and console reporting, with
  table-driven tests; the public API of package `filetrove` (`Export*TSV`) is unchanged
- `install.go`: `CheckInstall` writes its messages to stderr, with a test
- `README.md`: "Viewing results" section
