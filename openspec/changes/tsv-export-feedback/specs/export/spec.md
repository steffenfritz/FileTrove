# Spec Delta

## ADDED Requirements

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
print file paths to stdout in this case.

#### Scenario: Failure after some files were written
- **WHEN** `<uuid>_session.tsv` and `<uuid>_files.tsv` were written and writing `<uuid>_directories.tsv` fails
- **THEN** stderr contains the error, states that the export is incomplete and lists the two written files
- **AND** stdout is empty and the exit code is 1

#### Scenario: Failure before any file was written
- **WHEN** `<uuid>_session.tsv` cannot be created, e.g. because the working directory is not writable
- **THEN** stderr contains the error and states that no TSV file was written
- **AND** the exit code is 1
