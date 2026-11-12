# Tasks

## 1. Bloom filter selection

- [x] 1.1 Add `cmd/ftrove/nsrlselect.go` with an unexported `selectBloomPath(dbDir, variant string, explicit bool) (string, bool)` that returns the bloom filter path according to the cli spec (explicit: `nsrl-<variant>.bloom`, then `nsrl.bloom`; not explicit: first existing of `nsrl-all.bloom`, `nsrl-mobile.bloom`, `nsrl-modern.bloom`, `nsrl.bloom`) and whether a file was found; verify with table-driven tests in `cmd/ftrove/nsrlselect_test.go` using a temporary `db/` covering: only modern, modern + all, explicit modern with modern + all, explicit mobile with only modern (not found), only legacy with and without explicit variant, empty directory
- [x] 1.2 In `cmd/ftrove/main.go`, replace the hard-coded path at the bloom loading block with `selectBloomPath("db", *nsrlVariant, flag.CommandLine.Changed("nsrl-variant"))`; keep the existing `LoadNSRL` call and session NSRL version handling; verify with `task test`

## 2. Console warning and help text

- [x] 2.1 When no bloom filter is found or `ft.LoadNSRL` fails, keep the `logger.Warn` and additionally print `WARNING: NSRL bloom filter not available, NSRL checks disabled: <reason>` to stderr, independent of `-V`; for "not found" the reason names the file(s) looked for; verify by building `ftrove` and scanning `testdata` in an installation without a bloom filter (warning on stderr, exit 0, session NSRL version `none` in `ftrove -l`/session export) and with a dummy `db/nsrl-modern.bloom` (warning names the decode error)
- [x] 2.2 Change the `--nsrl-variant` help text to cover both uses, e.g. "NSRL bloom filter variant (modern, mobile, all): downloaded during install; used for a scan instead of the installed one."; verify with `ftrove --help`

## 3. End-to-end check

- [x] 3.1 Reproduce #183 with the fix: install into a temporary directory with `--nsrl-variant modern` (real `nsrl-modern.bloom`), run `ftrove -i testdata` without the flag; verify no NSRL warning, NSRL count in `logs/filetrove.log` > 0, and the session's NSRL version ends in `-modern`

## 4. Documentation

- [x] 4.1 Update the NSRL section in `README.md` and the bundle/NSRL part of `BUILDING.md`: a scan uses the installed bloom filter automatically (order `all`, `mobile`, `modern`, then legacy `nsrl.bloom`), `--nsrl-variant` on a scan selects a specific one, and a missing filter is reported on stderr; verify the described commands match the behavior from 2.1 and 3.1
- [x] 4.2 Check that the cli spec delta matches the implementation; verify with `openspec validate detect-nsrl-bloom-variant --strict`

## Workflow follow-up

- Run `/opsx:verify`, then `/opsx:archive` to merge the delta into `openspec/specs/cli/spec.md`.
- Reference #183 in the PR so the issue closes on merge.
- After PR #193 is merged, update open question 14 in `openspec/open-questions.md`: the scan side is resolved; only the dist bundle/`--install` variant mismatch remains.
