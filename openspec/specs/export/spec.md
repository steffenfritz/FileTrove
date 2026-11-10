# export Specification

## Purpose
Defines how scan results leave the database: session listings on the console and the
TSV, JSONL and PREMIS v3 exports consumed by other archival and forensic tools.

## Requirements

### Requirement: List sessions
`ftrove -l` SHALL print a table of all sessions with row id, UUID, start time, end time,
project, archivist name and mountpoint. `ftrove -L <uuid>` SHALL print UUID, project,
archivist, mountpoint, file count and NSRL count of one session.

#### Scenario: Find the UUID for an export
- **WHEN** the user runs `ftrove -l`
- **THEN** every session in the database is listed with its UUID

### Requirement: TSV export
`ftrove -t <uuid>` SHALL write the session to tab-separated files in the current
working directory: `<uuid>_session.tsv`, `<uuid>_files.tsv` and
`<uuid>_directories.tsv`, plus `<uuid>_exif.tsv`, `<uuid>_dublincore.tsv`,
`<uuid>_yara.tsv` and `<uuid>_xattr.tsv` when the respective feature flag of the
session is `True`. Each file SHALL start with a header row of the database column
names; NULL values SHALL be written as empty fields.

#### Scenario: Session without optional features
- **WHEN** the user exports a session that was scanned without `-e`, `-d`, `-y` and `-x`
- **THEN** exactly the session, files and directories TSV files are written

#### Scenario: Session with EXIF
- **WHEN** the user exports a session that was scanned with `-e`
- **THEN** `<uuid>_exif.tsv` is written in addition

### Requirement: TSV export reports written files
After a successful `ftrove -t <uuid>`, `ftrove` SHALL print the absolute path of every
TSV file it wrote to stdout, one path per line, in the order the files were written,
and exit with code 0. It SHALL also print a line to stderr that names the session UUID
and the number of files written. Without `-V`, nothing else SHALL be written to stdout.

#### Scenario: Paths on stdout
- **WHEN** the user runs `ftrove -t <uuid>` in `/work` for a session scanned without optional features
- **THEN** stdout contains exactly three lines: `/work/<uuid>_session.tsv`, `/work/<uuid>_files.tsv` and `/work/<uuid>_directories.tsv`
- **AND** the exit code is 0

#### Scenario: Optional files are listed
- **WHEN** the user exports a session that was scanned with `-e` and `-x`
- **THEN** stdout also lists `<uuid>_exif.tsv` and `<uuid>_xattr.tsv` as absolute paths

#### Scenario: Pipe the output
- **WHEN** the user runs `ftrove -t <uuid> | xargs wc -l`
- **THEN** `wc` receives only the TSV file paths, and banner and header are shown on the terminal through stderr

#### Scenario: Installation without NSRL bloom filter
- **WHEN** the user runs `ftrove -t <uuid>` in an installation without an NSRL bloom filter
- **THEN** stdout contains only the TSV file paths
- **AND** the message that the bloom filter is missing appears on stderr

#### Scenario: Summary on stderr
- **WHEN** the user exports a session and four TSV files are written
- **THEN** stderr contains a line naming the session UUID and the count 4

### Requirement: TSV export reports failures on the console
If writing any TSV file fails, `ftrove` SHALL print to stderr the error, a statement
that the export is incomplete, and the absolute paths of the TSV files written before
the failure, and exit with code 1. This SHALL apply with and without `-V`. It SHALL NOT
print the list of TSV file paths to stdout in this case.

#### Scenario: Failure after some files were written
- **WHEN** `<uuid>_session.tsv` and `<uuid>_files.tsv` were written and writing `<uuid>_directories.tsv` fails
- **THEN** stderr contains the error, states that the export is incomplete and lists the two written files
- **AND** stdout is empty and the exit code is 1

#### Scenario: Failure before any file was written
- **WHEN** `<uuid>_session.tsv` cannot be created, e.g. because the working directory is not writable
- **THEN** stderr contains the error and states that no TSV file was written
- **AND** the exit code is 1

### Requirement: JSONL export
`ftrove -j <uuid>` SHALL write the session to stdout as JSON Lines. Every line SHALL be
an object `{"type": <type>, "payload": <record>}`. Records SHALL appear in the order
session, file, directory, exif, dublincore, yara, xattr, ntfsads; types without rows
for the session SHALL be omitted. Log messages SHALL NOT be written to stdout.

#### Scenario: Filter files with jq
- **WHEN** the user runs `ftrove -j <uuid> | jq 'select(.type == "file")'`
- **THEN** one object per indexed file is returned

#### Scenario: Unknown session
- **WHEN** the user exports a UUID that does not exist
- **THEN** an error is logged and the exit code is 1

### Requirement: PREMIS v3 export
`ftrove -P <uuid>` SHALL write one PREMIS v3 XML document to stdout containing a
software agent for FileTrove with its version, a person agent when an archivist name is
recorded, one `ingestion` event identified by the session UUID with the session start
time, and one `premis:file` object per indexed file.

#### Scenario: Session with archivist
- **WHEN** the user exports a session scanned with `-a "J. Smith"`
- **THEN** the document contains a software agent, a person agent named "J. Smith", one ingestion event and one object per file

#### Scenario: Unknown session
- **WHEN** the user exports a UUID that does not exist
- **THEN** no XML is written, an error is logged and the exit code is 1

### Requirement: PREMIS file objects
Each PREMIS file object SHALL carry the file UUID, one fixity entry per non-empty hash
(`MD5`, `SHA-1`, `SHA-256`, `SHA-512`, `BLAKE2b-512`), the size, the format name
(falling back to the MIME type), the format version if known, the PRONOM identifier as
`specification` registry entry if known, the original file name, the file path as
content location, and a link to the session's event.

#### Scenario: File with PRONOM identification
- **WHEN** a session contains a PDF identified as `fmt/276`
- **THEN** its object has a format registry entry with name `PRONOM` and key `fmt/276`

### Requirement: Exports without full installation
JSONL and PREMIS exports SHALL work with only `db/filetrove.db` present, without
`siegfried.sig` or an NSRL bloom filter.

#### Scenario: Export on an analysis machine
- **WHEN** only `db/filetrove.db` is copied to another machine
- **THEN** `ftrove -j <uuid>` and `ftrove -P <uuid>` succeed there
