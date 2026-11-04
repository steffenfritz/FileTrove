# database Specification

## Purpose
Defines the SQLite database as FileTrove's primary result store and the guarantees
other tools (exports, webftrove, third-party queries) can rely on. The authoritative
column list is `database_schema.dbml`.

## Requirements

### Requirement: Single SQLite database
All scan results SHALL be stored in one SQLite database, `db/filetrove.db`, created by
`ftrove --install`. Its tables SHALL be `filetrove`, `sessionsmd`, `files`,
`directories`, `dublincore`, `exif`, `yara`, `xattr` and `ntfsads`, with the columns
documented in `database_schema.dbml`.

#### Scenario: Fresh installation
- **WHEN** `ftrove --install .` completes
- **THEN** `db/filetrove.db` contains all nine tables and no sessions

### Requirement: Database version record
The `filetrove` table SHALL contain the FileTrove version string of the binary that
created the database (`<version>+<git commit>` for builds via `task`) and the
installation date.

#### Scenario: Version after installation
- **WHEN** a binary built from VERSION `1.0.0-BETA-10` at commit `7f23486` installs a database
- **THEN** `SELECT version FROM filetrove` returns `1.0.0-BETA-10+7f23486`

### Requirement: Schema stability from release 1.0.0
Starting with release 1.0.0, the database schema SHALL remain stable within a major
version. Any change to tables or columns SHALL only be released with a new major
version.

#### Scenario: Schema change after 1.0.0
- **WHEN** a change adds, removes or renames a table or column after 1.0.0 has been released
- **THEN** it is released no earlier than the next major version (e.g. 2.0.0)

#### Scenario: Minor or patch release
- **WHEN** FileTrove 1.x.y is released
- **THEN** databases created by any earlier 1.x release remain usable without migration

### Requirement: Sessions group all results
Every row in `files`, `directories`, `dublincore`, `exif`, `yara`, `xattr` and
`ntfsads` SHALL reference its session through `sessionuuid`, which equals
`sessionsmd.uuid` of exactly one session.

#### Scenario: Query a session's files
- **WHEN** a consumer runs `SELECT * FROM files WHERE sessionuuid=<uuid>`
- **THEN** exactly the files indexed in that session are returned

### Requirement: Identifiers per run
Each file, directory, EXIF, YARA, xattr and Dublin Core entry SHALL get a fresh UUIDv4
when it is written. These identifiers are unique within the database but SHALL NOT be
assumed stable across separate scans of the same data.

#### Scenario: Scanning the same directory twice
- **WHEN** the same directory is scanned in two sessions
- **THEN** the same file has different `fileuuid` values in the two sessions

### Requirement: Feature flags in the session
`sessionsmd.exifflag`, `dublincoreflag`, `yaraflag` and `xattrflag` SHALL be `True` when
the corresponding feature was enabled for the session and empty otherwise.
`yarasource` SHALL contain the rule file path given with `-y`.

#### Scenario: EXIF enabled
- **WHEN** a session is scanned with `-e`
- **THEN** its `exifflag` is `True` and its `yaraflag` is empty

### Requirement: Hierarchy depth
`files.hierarchy` and `directories.hierarchy` SHALL contain the number of path
separators in the stored path.

#### Scenario: Nested file
- **WHEN** a file is stored with the path `data/a/b.txt` on a system using `/`
- **THEN** its `hierarchy` value is 2

### Requirement: Read-only access by webftrove
`webftrove` SHALL open the database in read-only mode and SHALL NOT modify it.

#### Scenario: Browsing results
- **WHEN** a user browses a database with `webftrove --db db/filetrove.db`
- **THEN** the database file is unchanged afterwards
