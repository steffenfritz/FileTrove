# Spec Delta

## ADDED Requirements

### Requirement: Reject unconsumed positional arguments
`ftrove` SHALL abort with exit code 1 if any positional arguments remain after flag
parsing, regardless of the requested mode. The check SHALL happen before any other
action, so that no session, database entry, log file or debug package is created.

#### Scenario: Shell-expanded glob passed to -i
- **WHEN** the user runs `ftrove -i corpus/*` and the shell expands it to `corpus/a corpus/b corpus/c`
- **THEN** ftrove exits with code 1
- **AND** no new session exists in `db/filetrove.db`

#### Scenario: Stray argument in an export
- **WHEN** the user runs `ftrove -j <uuid> extra`
- **THEN** ftrove exits with code 1 and writes nothing to stdout

#### Scenario: Valid invocation
- **WHEN** the user runs `ftrove -i corpus/`
- **THEN** the scan runs as before

### Requirement: Guidance for unexpected arguments
When ftrove rejects positional arguments, it SHALL print to stderr the list of unexpected
arguments, that `-i` expects exactly one directory, and that a shell glob does not match
hidden files. If `-i` was given, it SHALL suggest scanning the parent directory of the
`-i` value instead.

#### Scenario: Error message after glob expansion
- **WHEN** the user runs `ftrove -i corpus/*` and the shell expands it to `corpus/a corpus/b corpus/c`
- **THEN** stderr names `corpus/b` and `corpus/c` as unexpected arguments
- **AND** stderr suggests `ftrove -i corpus` and mentions that hidden files would be missed by a glob

#### Scenario: Positional argument without -i
- **WHEN** the user runs `ftrove corpus/`
- **THEN** stderr names `corpus/` as unexpected argument and explains that the directory is passed with `-i`
