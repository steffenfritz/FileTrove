# cli Specification

## Purpose
Defines how users invoke `ftrove`: installation, scanning, resuming a scan, and the
working-directory layout the binary relies on. Export and listing modes are specified
in the `export` capability.

## Requirements

### Requirement: Working directory layout
`ftrove` SHALL resolve its database, signature file, NSRL bloom filter and log file
relative to the current working directory: `db/filetrove.db`, `db/siegfried.sig`,
`db/nsrl-<variant>.bloom` (fallback `db/nsrl.bloom`) and `logs/filetrove.log`.

#### Scenario: Scan from the installation directory
- **WHEN** the user runs `ftrove -i <dir>` from a directory created by `ftrove --install`
- **THEN** results are written to `db/filetrove.db` and log output to `logs/filetrove.log` below that directory

### Requirement: Installation
`ftrove --install <path>` SHALL create `<path>/db` and `<path>/logs`, create an empty
FileTrove database `<path>/db/filetrove.db` that records the FileTrove version and the
install date, and obtain `siegfried.sig` unless it is already present. It SHALL then try
to provide the NSRL bloom filter for the selected `--nsrl-variant` (`modern`, `mobile`,
`all`; default `all`) by copying it from a known location or downloading it. After
installation the program SHALL exit without scanning.

#### Scenario: Successful installation
- **WHEN** the user runs `ftrove --install .` with network access
- **THEN** `db/filetrove.db`, `db/siegfried.sig`, `db/nsrl-all.bloom` and `logs/` exist
- **AND** the exit code is 0

#### Scenario: Siegfried signature cannot be obtained
- **WHEN** `siegfried.sig` is not present and cannot be downloaded
- **THEN** an error explaining that `siegfried.sig` must be copied into `db/` is logged
- **AND** the exit code is 1

#### Scenario: NSRL bloom filter cannot be obtained
- **WHEN** the bloom filter can be neither copied nor downloaded
- **THEN** the installer prints how to build or copy it manually and that scanning works without it
- **AND** the installation still succeeds

#### Scenario: Unknown NSRL variant
- **WHEN** the user passes `--nsrl-variant foo`
- **THEN** the installer prints that the variant is unknown and falls back to `all`

### Requirement: Version output
`ftrove -v` SHALL print the banner and license text and exit with code 0 without
touching the database.

#### Scenario: Print version
- **WHEN** the user runs `ftrove -v`
- **THEN** banner and license are printed and the exit code is 0

### Requirement: Scan session
A scan (`ftrove -i <dir>`) SHALL create exactly one new session with a fresh UUIDv4,
record its start time (RFC 3339), the absolute path of `<dir>` as mountpoint, the
optional project name (`-p`) and archivist name (`-a`), the enabled optional features,
and the versions of FileTrove, siegfried, Go and the NSRL bloom filter. When the scan
finishes, the end time (RFC 3339) SHALL be recorded.

#### Scenario: Completed scan
- **WHEN** the user runs `ftrove -i /media/evidence -p "Case 42" -a "J. Smith"` and the scan completes
- **THEN** the database contains one new session with project "Case 42", archivist "J. Smith", mountpoint `/media/evidence`, a start time and an end time

#### Scenario: Scan without NSRL bloom filter
- **WHEN** no bloom filter is present in `db/`
- **THEN** a warning that NSRL checks are disabled is logged
- **AND** the session records `none` as NSRL version
- **AND** the scan continues

### Requirement: Scan summary
At the end of a scan `ftrove` SHALL log the scanned path, the session UUID, the number
of indexed files, the number of indexed directories, the number of files found in the
NSRL bloom filter, the duration, and the command to export the session as TSV.

#### Scenario: Summary in the log file
- **WHEN** a scan completes
- **THEN** `logs/filetrove.log` contains the session UUID and the file, directory and NSRL counts

### Requirement: Log output
After the installation check succeeds, log messages SHALL be appended to
`logs/filetrove.log`. With `-V` they SHALL additionally be written to stdout.

#### Scenario: Verbose scan
- **WHEN** the user runs `ftrove -V -i <dir>`
- **THEN** log messages appear both on stdout and in `logs/filetrove.log`

### Requirement: Optional scan features
The flags `-e` (EXIF for JPEG/TIFF), `-x` (extended attributes), `-y <rules>` (YARA-X
rules) and `-d <json>` (Dublin Core metadata for the session) SHALL enable the
corresponding collection, and the session SHALL record each enabled feature as `True`.

#### Scenario: YARA rules cannot be compiled
- **WHEN** the user passes `-y` with a rule file that does not compile
- **THEN** an error is logged and the exit code is 1

#### Scenario: Dublin Core file cannot be read
- **WHEN** the user passes `-d` with a file that is not valid Dublin Core JSON
- **THEN** an error is logged and the exit code is 1

### Requirement: Time zone for timestamps
With `-z <IANA zone>` file and directory timestamps SHALL be converted to that zone;
without it the local time zone SHALL be used.

#### Scenario: Invalid time zone
- **WHEN** the user passes `-z Mars/Olympus`
- **THEN** an error with an example of a valid zone is logged and the exit code is 1

### Requirement: Resume an aborted session
`ftrove -r <session uuid>` SHALL continue the given session: it SHALL reuse the stored
mountpoint (`-i` is not required), keep the session UUID, rebuild the file list and
process only the files that follow the last file recorded for that session. File and
NSRL counts in the summary SHALL include the files processed before the abort.

#### Scenario: Resume after abort
- **WHEN** a scan was aborted after recording file N and the user runs `ftrove -r <uuid>`
- **THEN** processing starts with the file following N in walk order
- **AND** no new session is created

#### Scenario: Nothing left to process
- **WHEN** the last recorded file is the last file of the walk
- **THEN** `ftrove` logs that no files are left to process and exits with code 0

#### Scenario: Last recorded file no longer present
- **WHEN** the last recorded file is not part of the rebuilt file list
- **THEN** an error naming that file is logged and the exit code is 1
