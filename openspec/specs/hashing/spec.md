# hashing Specification

## Purpose
Defines which cryptographic hashes FileTrove computes for every file and how they are
represented, since hash values serve as fixity evidence in archival and forensic use.

## Requirements

### Requirement: Five hash algorithms per file
For every indexed file FileTrove SHALL compute MD5, SHA-1, SHA-256, SHA-512 and
BLAKE2b-512 over the complete file content.

#### Scenario: Regular file
- **WHEN** a readable regular file is indexed
- **THEN** its record contains an MD5, SHA-1, SHA-256, SHA-512 and BLAKE2b-512 value

### Requirement: Single read for all hashes
All five hashes of a file SHALL be computed from a single sequential read of the file,
so that every hash describes the same byte sequence.

#### Scenario: Hashes describe identical content
- **WHEN** a file is indexed
- **THEN** all five hash values are derived from one read of the file

### Requirement: Lowercase hexadecimal representation
Hash values SHALL be stored and exported as lowercase hexadecimal strings without
prefix or separators.

#### Scenario: Empty file
- **WHEN** an empty file is indexed
- **THEN** its MD5 value is `d41d8cd98f00b204e9800998ecf8427e`
- **AND** its SHA-1 value is `da39a3ee5e6b4b0d3255bfef95601890afd80709`
- **AND** its SHA-256 value is `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`
