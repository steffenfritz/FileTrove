## MODIFIED Requirements

### Requirement: Working directory layout
`ftrove` SHALL resolve its database, signature file, NSRL bloom filter and log file
relative to the current working directory: `db/filetrove.db`, `db/siegfried.sig`,
`db/nsrl-<variant>.bloom` (fallback `db/nsrl.bloom`) and `logs/filetrove.log`.

For a scan or resume, `ftrove` SHALL select the NSRL bloom filter as follows:
- If `--nsrl-variant <variant>` is given, it SHALL use `db/nsrl-<variant>.bloom`, and if
  that file does not exist, `db/nsrl.bloom`. It SHALL NOT use a different variant.
- If `--nsrl-variant` is not given, it SHALL use the first existing file of
  `db/nsrl-all.bloom`, `db/nsrl-mobile.bloom`, `db/nsrl-modern.bloom` and
  `db/nsrl.bloom`, in this order.

#### Scenario: Scan from the installation directory
- **WHEN** the user runs `ftrove -i <dir>` from a directory created by `ftrove --install`
- **THEN** results are written to `db/filetrove.db` and log output to `logs/filetrove.log` below that directory

#### Scenario: Installed with a non-default variant
- **WHEN** the installation was created with `ftrove --install . --nsrl-variant modern` and `db/` contains only `nsrl-modern.bloom`
- **AND** the user runs `ftrove -i <dir>` without `--nsrl-variant`
- **THEN** NSRL checks run against `db/nsrl-modern.bloom`
- **AND** the session records the NSRL version of that filter

#### Scenario: Several variants installed
- **WHEN** `db/` contains `nsrl-modern.bloom` and `nsrl-all.bloom`
- **AND** the user runs `ftrove -i <dir>` without `--nsrl-variant`
- **THEN** NSRL checks run against `db/nsrl-all.bloom`

#### Scenario: Explicit variant
- **WHEN** `db/` contains `nsrl-modern.bloom` and `nsrl-all.bloom`
- **AND** the user runs `ftrove -i <dir> --nsrl-variant modern`
- **THEN** NSRL checks run against `db/nsrl-modern.bloom`

#### Scenario: Explicit variant not installed
- **WHEN** `db/` contains only `nsrl-modern.bloom`
- **AND** the user runs `ftrove -i <dir> --nsrl-variant mobile`
- **THEN** the scan runs without NSRL checks
- **AND** a warning naming the missing bloom filter is printed to stderr

#### Scenario: Legacy bloom filter
- **WHEN** `db/` contains only `nsrl.bloom`
- **AND** the user runs `ftrove -i <dir>` with or without `--nsrl-variant`
- **THEN** NSRL checks run against `db/nsrl.bloom`

### Requirement: Scan session
A scan (`ftrove -i <dir>`) SHALL create exactly one new session with a fresh UUIDv4,
record its start time (RFC 3339), the absolute path of `<dir>` as mountpoint, the
optional project name (`-p`) and archivist name (`-a`), the enabled optional features,
and the versions of FileTrove, siegfried, Go and the NSRL bloom filter. When the scan
finishes, the end time (RFC 3339) SHALL be recorded.

If no NSRL bloom filter can be loaded, the scan SHALL continue without NSRL checks, and
`ftrove` SHALL log a warning and also print it to stderr, independent of `-V`, stating
that NSRL checks are disabled and why.

#### Scenario: Completed scan
- **WHEN** the user runs `ftrove -i /media/evidence -p "Case 42" -a "J. Smith"` and the scan completes
- **THEN** the database contains one new session with project "Case 42", archivist "J. Smith", mountpoint `/media/evidence`, a start time and an end time

#### Scenario: Scan without NSRL bloom filter
- **WHEN** no bloom filter is present in `db/`
- **THEN** a warning that NSRL checks are disabled is logged
- **AND** the same warning is printed to stderr, also without `-V`
- **AND** the session records `none` as NSRL version
- **AND** the scan continues

#### Scenario: Bloom filter cannot be loaded
- **WHEN** the selected bloom filter exists but cannot be read or decoded
- **THEN** a warning naming the file and the error is logged and printed to stderr
- **AND** the session records `none` as NSRL version
- **AND** the scan continues
