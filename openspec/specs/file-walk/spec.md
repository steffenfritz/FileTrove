# file-walk Specification

## Purpose
Defines which filesystem entries a scan visits, in which order, and how entries that
cannot or must not be processed are reported.

## Requirements

### Requirement: Recursive walk in lexical order
A scan SHALL walk the input directory recursively and visit entries in lexical order
within each directory, so that repeated walks of an unchanged tree yield the same
sequence. Resuming a session relies on this order.

#### Scenario: Nested directories
- **WHEN** the input directory contains `b.txt`, `a/` and `a/c.txt`
- **THEN** the walk order is `a`, `a/c.txt`, `b.txt`

### Requirement: Regular files are indexed
Every regular file reached by the walk SHALL be processed and recorded as one entry
of the session.

#### Scenario: Regular file
- **WHEN** the input directory contains a readable regular file
- **THEN** the session contains one file record for it

### Requirement: Directories are indexed
Every directory reached by the walk, including the input directory itself, SHALL be
recorded as one directory entry of the session.

#### Scenario: Input directory
- **WHEN** the user scans `<dir>` containing one subdirectory
- **THEN** the session contains directory records for `<dir>` and the subdirectory

### Requirement: Symlinks and special files are skipped and reported
Symbolic links SHALL NOT be followed. Symbolic links and special files (sockets, named
pipes, device nodes) SHALL NOT be indexed as files, and each SHALL be logged as a
skipped entry with its path.

#### Scenario: Symlink in the input tree
- **WHEN** the input directory contains a symlink to a file outside the tree
- **THEN** the symlink target is not indexed
- **AND** the log contains a warning naming the symlink path as skipped

### Requirement: Inaccessible directories do not abort the walk
If a directory cannot be read (e.g. permission denied, stale network mount), the walk
SHALL continue with the remaining entries and SHALL log the path as skipped.

#### Scenario: Unreadable subdirectory
- **WHEN** a subdirectory has no read permission for the scanning user
- **THEN** the scan completes with all other entries
- **AND** the log contains a warning naming the subdirectory as skipped

### Requirement: Filesystem boundaries are crossed
The walk SHALL descend into mount points below the input directory (including network
shares); it does not restrict itself to one device.

#### Scenario: Mounted volume below the input directory
- **WHEN** a different filesystem is mounted below the input directory
- **THEN** its files are indexed as part of the session
