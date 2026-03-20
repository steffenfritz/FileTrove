# Building FileTrove from Source

FileTrove requires CGO because it links against the [YARA-X](https://virustotal.github.io/yara-x/) C library. This guide walks you through setting up a complete development environment on **macOS** and **Debian/Ubuntu**.

---

## macOS

### 1. Install Homebrew

If you don't have Homebrew yet:

```sh
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

### 2. Install the Task build tool

```sh
brew install go-task/tap/go-task
```

### 3. Install remaining dependencies and build YARA-X

From the repository root, run:

```sh
cd cmd/ftrove
task setup_mac
```

This will:
- Install Go, Rust, Zig, SQLite, and pkg-config via Homebrew
- Install `cargo-c`
- Clone YARA-X v1.14.0 into `$HOME/yara-x` and build its C library into `$HOME/yara_install`
- Write the required CGO environment variables to `cmd/ftrove/.yara_env` (auto-loaded by subsequent `task` invocations)

### 4. Build

```sh
task build
```

---

## Debian / Ubuntu (latest stable or LTS)

### 1. Install the Task build tool

```sh
curl -fsSL https://taskfile.dev/install.sh | sh -s -- -b ~/.local/bin
```

Ensure `~/.local/bin` is on your `PATH`, then verify with `task --version`.

### 2. Install remaining dependencies and build YARA-X

From the repository root, run:

```sh
cd cmd/ftrove
task setup_linux
```

This will:
- Install build tools, pkg-config, sqlite3, and Go via apt (Go 1.21+ supports `GOTOOLCHAIN=auto`, so the version in `go.mod` is fetched automatically)
- Install Rust via rustup and `cargo-c`
- Clone YARA-X v1.14.0 into `$HOME/yara-x` and build its C library into `$HOME/yara_install`
- Write the required CGO environment variables to `cmd/ftrove/.yara_env` (auto-loaded by subsequent `task` invocations)

### 3. Build

```sh
task build
```

---

## Building the NSRL Bloom filter

FileTrove uses a Bloom filter for NSRL lookups. If you cloned the repo with Git LFS, `db/nsrl.bloom` is already present. Otherwise, build it from upstream NIST data (requires `sqlite3`, `curl`, `unzip`, and a built `admftrove`):

```sh
task nsrl:build-modern    # Modern OS software only (~30-45 MB)
task nsrl:build-mobile    # Modern + Android + iOS (~50-65 MB)
task nsrl:build-all       # All subsets including legacy (~80-110 MB)
```

Check whether your local bloom file matches the configured upstream version:

```sh
task nsrl:check
```

To update, bump `NSRL_VERSION` in `Taskfile.nsrl.yml` to the latest [NIST RDS release](https://www.nist.gov/itl/ssd/software-quality-group/national-software-reference-library-nsrl/nsrl-download/current-rds), run `task nsrl:clean`, and rebuild.

---

## Running tests

From the repository root:

```sh
go test -v ./...
```

For a full integration test (build + install + scan + DB inspection):

```sh
cd cmd/ftrove
task changetest
```
