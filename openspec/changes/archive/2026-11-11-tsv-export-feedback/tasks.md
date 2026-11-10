# Tasks

## 1. Console report for the TSV export

- [x] 1.1 Add `cmd/ftrove/tsvexport.go` with unexported `tsvFileName(uuid, table string) string`, `reportTSVExport(stdout, stderr io.Writer, uuid string, files []string)` (absolute paths one per line on stdout, header with UUID and file count on stderr) and `reportTSVExportFailure(stderr io.Writer, uuid, failed string, files []string, err error)` (error, failed file possibly partial, "export incomplete", absolute paths of completed files or "no TSV file was written"); verify with table-driven tests in `cmd/ftrove/tsvexport_test.go` covering three files, seven files, failure after two files and failure before the first file, checking stdout and stderr separately
- [x] 1.2 In the TSV block of `cmd/ftrove/main.go`, append `tsvFileName(uuid, <table>)` after each successful `ft.Export*TSV` call, call `reportTSVExportFailure` with `os.Stderr` before every `os.Exit(1)` in that block, and call `reportTSVExport` with `os.Stdout`/`os.Stderr` before the final `return`; keep the existing `logger` calls; verify with `task test`, then build `ftrove` and run `ftrove -t <uuid>` on a test session (paths on stdout, exit 0), `ftrove -t <uuid> 2>/dev/null | xargs ls -l` (all listed files exist), and the export in a read-only directory (error on stderr, empty stdout, exit 1)
- [x] 1.3 In `CheckInstall` in `install.go`, write all messages to `os.Stderr` instead of stdout; verify with a test in `install_test.go` that runs `CheckInstall` in a temporary directory without a bloom filter and checks that stdout is empty and stderr names the missing bloom filter, then rerun `ftrove -t <uuid> 2>/dev/null | xargs ls -l` in an installation without a bloom filter (only the TSV paths reach `ls`)

## 2. Documentation

- [x] 2.1 Update "Viewing results" in `README.md`: show example output of `./ftrove -t <uuid>`, an example of piping it (e.g. `./ftrove -t <uuid> | xargs wc -l`), and a note that `-V` also writes log lines to stdout; verify the documented commands produce the documented output
- [x] 2.2 Check that the export spec delta matches the implemented output (open questions 7 and 13 stay in `openspec/open-questions.md`); verify with `openspec validate tsv-export-feedback --strict`

## Workflow follow-up

- Run `/opsx:verify`, then `/opsx:archive` to merge the delta into `openspec/specs/export/spec.md`.
- Reference #185 in the PR so the issue closes on merge.
