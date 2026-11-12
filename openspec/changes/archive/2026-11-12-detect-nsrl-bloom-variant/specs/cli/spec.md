## ADDED Requirements

### Requirement: NSRL bloom filter selection
For a scan or resume, `ftrove` SHALL select the NSRL bloom filter in `db/` as follows.
With `--nsrl-variant <variant>`, it SHALL use `nsrl-<variant>.bloom`, else `nsrl.bloom`,
and SHALL NOT use a different variant. Without `--nsrl-variant`, it SHALL use the first
existing file of `nsrl-all.bloom`, `nsrl-mobile.bloom`, `nsrl-modern.bloom` and
`nsrl.bloom`.

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

### Requirement: Warning when NSRL checks are disabled
If no NSRL bloom filter can be loaded, the scan SHALL continue without NSRL checks, and
`ftrove` SHALL log a warning and also print it to stderr, independent of `-V`, stating
that NSRL checks are disabled and why.

#### Scenario: No bloom filter on the console
- **WHEN** no bloom filter is present in `db/`
- **AND** the user runs `ftrove -i <dir>` without `-V`
- **THEN** a warning naming the bloom filter files looked for is printed to stderr
- **AND** the session records `none` as NSRL version
- **AND** the exit code is 0

#### Scenario: Bloom filter cannot be loaded
- **WHEN** the selected bloom filter exists but cannot be read or decoded
- **THEN** a warning naming the file and the error is logged and printed to stderr
- **AND** the session records `none` as NSRL version
- **AND** the scan continues
